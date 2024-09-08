package core

import (
	"fmt"
	"math"
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

func (trk *Trk) calcTopograph(detail bool) []TrkSummary {
	var distance, denivPos, denivNeg float64 = 0, 0, 0

	var p_prev Pos
	pointName_prev := "start"
	n := 0
	var trkSummary []TrkSummary

	for i, trkpt := range (*trk).Trkseg.Trkpt {
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
					VitessePlat:    trk.Vitesse,
					Distance:       distance,
					DenivPos:       denivPos,
					DenivNeg:       denivNeg,
					DistanceEffort: CalcDistanceEffort(distance, denivPos, denivNeg),
				}
				_, x.DurationHour, x.DurationMin = CalcDuration(x.DistanceEffort, trk.Vitesse)

				trkSummary = append(trkSummary, x)

				pointName_prev = *trkpt.Name

				distance = 0
				denivPos = 0
				denivNeg = 0
				n = 0
			}
		}

		if i == len((*trk).Trkseg.Trkpt)-1 {
			x := TrkSummary{
				From:           pointName_prev,
				To:             "end",
				NPoints:        n,
				VitessePlat:    trk.Vitesse,
				Distance:       distance,
				DenivPos:       denivPos,
				DenivNeg:       denivNeg,
				DistanceEffort: CalcDistanceEffort(distance, denivPos, denivNeg),
			}
			_, x.DurationHour, x.DurationMin = CalcDuration(x.DistanceEffort, trk.Vitesse)

			trkSummary = append(trkSummary, x)
		}

	}

	(*trk).DenivPos = denivPos
	(*trk).DenivNeg = denivNeg
	(*trk).Distance = distance

	return trkSummary
}

func (trk *Trk) convertToEffortMetrics() {
	(*trk).DistanceEffort = CalcDistanceEffort(
		(*trk).Distance,
		(*trk).DenivPos,
		(*trk).DenivNeg,
	)
}

func (trk *Trk) SetVitesse(v float64) {
	(*trk).Vitesse = v
}

func (trk *Trk) calcDuration() {
	// (*trk).Duration = (*trk).DistanceEffort / (*trk).Vitesse
	// (*trk).DurationHour, (*trk).DurationMin = FloatToHourMin((*trk).Duration)
	(*trk).Duration, (*trk).DurationHour, (*trk).DurationMin = CalcDuration(
		(*trk).DistanceEffort, (*trk).Vitesse,
	)
}

func (trk Trk) CalcAll(detail bool) []TrkSummary {
	trkSummary := trk.calcTopograph(detail)
	trk.convertToEffortMetrics()
	trk.calcDuration()

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

func (trk Trk) PrintInfo(ascii_format ...bool) {
	if len(ascii_format) > 0 {
		if ascii_format[0] {
			fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trk.Name)
		} else {
			fmt.Printf("Track name: %v\n", trk.Name)
		}
	} else {
		fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trk.Name)
	}

	fmt.Println("Number of points: ", len(trk.Trkseg.Trkpt))

	fmt.Printf("Distance:               %.1f km\n", trk.Distance)

	fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", trk.DenivPos, trk.DenivNeg)

	fmt.Printf("Distance effort:        %.1f km\n", trk.DistanceEffort)

	fmt.Printf("Vitesse sur plat:       %.1f km/h\n", trk.Vitesse)
	fmt.Printf("Temps parcours estimé:  %vh%v\n", trk.DurationHour, trk.DurationMin)

}
