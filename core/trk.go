package core

import (
	"fmt"
)

type Pos struct {
	Lat float64
	Lon float64
	Ele float64
}

func (trk *Trk) calcTopograph() (float64, float64, float64) {
	var distance, denivPos, denivNeg float64 = 0, 0, 0

	var p_prev Pos
	for i, trkpt := range (*trk).Trkseg.Trkpt {
		p := Pos{
			Lat: trkpt.Lat,
			Lon: trkpt.Lon,
			Ele: trkpt.Ele,
		}
		if i == 0 {
			p_prev = p
			continue
		}

		eleDiff := DiffElevation(p_prev, p)
		if eleDiff > 1 {
			denivPos += eleDiff
		} else if eleDiff < -1 {
			denivNeg += eleDiff
		}

		distance += Dist(p_prev, p)
		p_prev = p
	}

	(*trk).DenivPos = denivPos
	(*trk).DenivNeg = denivNeg
	(*trk).Distance = distance

	return distance, denivPos, denivNeg
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
	(*trk).Duration = (*trk).DistanceEffort / (*trk).Vitesse
	(*trk).DurationHour, (*trk).DurationMin = FloatToHourMin((*trk).Duration)
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
}

func (trk *Trk) CalcAll() TrkSummary {
	(*trk).calcTopograph()
	// (*trk).Distance()
	// (*trk).calcDeniv()
	(*trk).convertToEffortMetrics()
	(*trk).calcDuration()

	return TrkSummary{
		From:    trk.Name,
		To:      "",
		NPoints: len(trk.Trkseg.Trkpt),

		VitessePlat:    trk.Vitesse,
		Distance:       trk.Distance,
		DenivPos:       trk.DenivPos,
		DenivNeg:       trk.DenivNeg,
		DistanceEffort: trk.DistanceEffort,
		DurationHour:   trk.DurationHour,
		DurationMin:    trk.DurationMin,
	}
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

	fmt.Printf("Vitesse sur plat:       %.0f km/h\n", summary.VitessePlat)
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

	fmt.Printf("Vitesse sur plat:       %.0f km/h\n", trk.Vitesse)
	fmt.Printf("Temps parcours estimé:  %vh%v\n", trk.DurationHour, trk.DurationMin)

}
