package main

import (
	"fmt"
	"math"
)

type Pos struct {
	Lat float64
	Lon float64
	Ele float64
}

func (trk *Trk) calcDeniv() (float64, float64) {
	var denivPos, denivNeg float64 = 0, 0

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

		eleDiff := diffElevation(p_prev, p)
		if eleDiff > 1 {
			denivPos += eleDiff
		} else if eleDiff < -1 {
			denivNeg += eleDiff
		}

		p_prev = p
	}

	(*trk).DenivPos = denivPos
	(*trk).DenivNeg = denivNeg

	return denivPos, denivNeg
}

func (trk *Trk) calcDistance() float64 {
	var d float64 = 0

	var p_prev Pos
	for i, trkpt := range (*trk).Trkseg.Trkpt {
		p := Pos{
			Lat: trkpt.Lat,
			Lon: trkpt.Lon,
			Ele: trkpt.Ele,
		}
		if i == 0 {
			p_prev = p
		}

		d += dist(p_prev, p)
		p_prev = p
	}

	(*trk).Distance = d
	return d

}

func (trk *Trk) calcDenivEffort() {
	(*trk).DenivPosEffort = (*trk).DenivPos / 100 // 100m D+ = 1km plat
	(*trk).DenivNegEffort = (*trk).DenivNeg / 300 // 300m D- = 1km plat
}

func (trk *Trk) calcDistanceEffort() {
	(*trk).DistanceEffort = (*trk).Distance + (*trk).DenivPosEffort + math.Abs((*trk).DenivNegEffort)
}

func (trk *Trk) setVitesse(v float64) {
	(*trk).Vitesse = v
}

func (trk *Trk) calcDuration() {
	(*trk).Duration = (*trk).DistanceEffort / (*trk).Vitesse
	(*trk).DurationHour = math.Floor((*trk).Duration)
	(*trk).DurationMin = ((*trk).Duration - (*trk).DurationHour) * 60
}

func (trk *Trk) calcAll() {
	(*trk).calcDistance()
	(*trk).calcDeniv()
	(*trk).calcDenivEffort()
	(*trk).calcDistanceEffort()
	(*trk).calcDuration()
}

func (trk Trk) Info() {
	fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trk.Name)
	// fmt.Printf("Track name: %v\n", trk.Name)
	fmt.Println("Number of points: ", len(trk.Trkseg.Trkpt))

	fmt.Printf("Distance:               %.1f km\n", trk.Distance)

	fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", trk.DenivPos, trk.DenivNeg)

	fmt.Printf("Distance effort:        %.1f km\n", trk.DistanceEffort)

	fmt.Printf("Vitesse sur plat:       %.0f km/h\n", trk.Vitesse)
	fmt.Printf("Temps parcours estimé:  %.0fh%.0f\n", trk.DurationHour, trk.DurationMin)

}
