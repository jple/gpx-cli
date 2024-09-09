package core

import (
	"fmt"
	"math"
	"slices"
)

type Pos struct {
	Lat float64
	Lon float64
	Ele float64
}

type TrkSummary struct {
	From           string
	To             string
	NPoints        int
	VitessePlat    float64
	Distance       float64
	DenivPos       float64
	DenivNeg       float64
	DistanceEffort float64
	DurationHour   int8
	DurationMin    int8
	Toto           string
}

func (trk *Trk) calcTopograph(vitessePlat float64, detail bool) []TrkSummary {
	var distance, denivPos, denivNeg float64 = 0, 0, 0

	var p_prev Pos
	pointName_prev := "start"
	n := 0
	var trkSummary []TrkSummary

	trkpts := slices.Concat((*trk).Trkseg)[0].Trkpt
	for i, trkpt := range trkpts {
		p := Pos{
			Lat: trkpt.Lat,
			Lon: trkpt.Lon,
			Ele: trkpt.Ele,
		}
		if i == 0 {
			p_prev = p
			if trkpt.Name != nil {
				pointName_prev = *trkpt.Name
			}
			continue
		}

		eleDiff := DiffElevation(p_prev, p)
		denivPos += math.Max(eleDiff, 1)
		denivNeg += math.Min(eleDiff, -1)

		distance += Dist(p_prev, p)
		n += 1
		p_prev = p

		if detail {
			if trkpt.Name != nil {
				x := TrkSummary{
					From:           pointName_prev,
					To:             *trkpt.Name,
					NPoints:        n,
					VitessePlat:    vitessePlat,
					Distance:       distance,
					DenivPos:       denivPos,
					DenivNeg:       denivNeg,
					DistanceEffort: CalcDistanceEffort(distance, denivPos, denivNeg),
				}
				_, x.DurationHour, x.DurationMin = CalcDuration(x.DistanceEffort, vitessePlat)

				trkSummary = append(trkSummary, x)

				pointName_prev = *trkpt.Name

				distance = 0
				denivPos = 0
				denivNeg = 0
				n = 0
			}
		}

		if i == len(trkpts)-1 {
			x := TrkSummary{
				From:           pointName_prev,
				To:             "end",
				NPoints:        n,
				VitessePlat:    vitessePlat,
				Distance:       distance,
				DenivPos:       denivPos,
				DenivNeg:       denivNeg,
				DistanceEffort: CalcDistanceEffort(distance, denivPos, denivNeg),
			}
			_, x.DurationHour, x.DurationMin = CalcDuration(x.DistanceEffort, vitessePlat)

			trkSummary = append(trkSummary, x)
		}

	}

	(*trk).Extensions.DenivPos = denivPos
	(*trk).Extensions.DenivNeg = denivNeg
	(*trk).Extensions.Distance = distance

	return trkSummary
}

func (trk *Trk) convertToEffortMetrics() {
	(*trk).Extensions.DistanceEffort = CalcDistanceEffort(
		(*trk).Extensions.Distance,
		(*trk).Extensions.DenivPos,
		(*trk).Extensions.DenivNeg,
	)
}

func (trk *Trk) calcDuration(vitessePlat float64) {
	(*trk).Extensions.Duration, (*trk).Extensions.DurationHour, (*trk).Extensions.DurationMin = CalcDuration(
		(*trk).Extensions.DistanceEffort, vitessePlat,
	)
}

func (trk Trk) CalcAll(vitessePlat float64, detail bool) []TrkSummary {
	trkSummary := trk.calcTopograph(vitessePlat, detail)
	trk.convertToEffortMetrics()
	trk.calcDuration(vitessePlat)

	return trkSummary
}

func (summary TrkSummary) Print(ascii_format ...bool) {
	if len(ascii_format) > 0 {
		if ascii_format[0] {
			fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", summary.From)
		} else {
			fmt.Printf("Track name: %v\n", summary.From)
		}
	} else {
		fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", summary.From)
	}

	fmt.Printf("Number of points:       %v\n", summary.NPoints)

	fmt.Printf("Distance:               %.1f km\n", summary.Distance)

	fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", summary.DenivPos, summary.DenivNeg)

	fmt.Printf("Distance effort:        %.1f km\n", summary.DistanceEffort)

	fmt.Printf("Vitesse sur plat:       %.1f km/h\n", summary.VitessePlat)
	fmt.Printf("Temps parcours estimé:  %vh%v\n", summary.DurationHour, summary.DurationMin)
}

func (trk Trk) PrintInfo(vitessePlat float64, ascii_format ...bool) {
	if len(ascii_format) > 0 {
		if ascii_format[0] {
			fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trk.Name)
		} else {
			fmt.Printf("Track name: %v\n", trk.Name)
		}
	} else {
		fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trk.Name)
	}

	trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
	fmt.Println("Number of points: ", len(trkpts))

	fmt.Printf("Distance:               %.1f km\n", trk.Extensions.Distance)

	fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", trk.Extensions.DenivPos, trk.Extensions.DenivNeg)

	fmt.Printf("Distance effort:        %.1f km\n", trk.Extensions.DistanceEffort)

	fmt.Printf("Vitesse sur plat:       %.1f km/h\n", vitessePlat)
	fmt.Printf("Temps parcours estimé:  %vh%v\n", trk.Extensions.DurationHour, trk.Extensions.DurationMin)

}
