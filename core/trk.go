package core

import (
	"fmt"
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

		eleDiff := DiffElevation(p_prev, p)
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

		d += Dist(p_prev, p)
		p_prev = p
	}

	(*trk).Distance = d
	return d

}

// func (trk *Trk) calcDenivEffort() {
// 	(*trk).DenivPosEffort = (*trk).DenivPos / 100 // 100m D+ = 1km plat
// 	(*trk).DenivNegEffort = (*trk).DenivNeg / 300 // 300m D- = 1km plat
// }

func (trk *Trk) CalcDistanceEffort() {
	(*trk).DistanceEffort = CalcDistanceEffort(
		(*trk).Distance,
		(*trk).DenivPos,
		(*trk).DenivNeg,
	)
}

func (trk *Trk) setVitesse(v float64) {
	(*trk).Vitesse = v
}

func (trk *Trk) calcDuration() {
	(*trk).Duration = (*trk).DistanceEffort / (*trk).Vitesse
	(*trk).DurationHour, (*trk).DurationMin = FloatToHourMin((*trk).Duration)
}

func (trk *Trk) calcAll() {
	(*trk).calcDistance()
	(*trk).calcDeniv()
	// (*trk).calcDenivEffort()
	(*trk).CalcDistanceEffort()
	(*trk).calcDuration()
}

func (trk Trk) Info(ascii_format ...bool) {
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
