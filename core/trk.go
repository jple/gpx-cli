package core

import (
	"math"
	"slices"
	"strconv"
)

func (trk Trk) GetInfo(vitessePlat float64, detail bool) TrkSummary {
	var distance, denivPos, denivNeg float64 = 0, 0, 0

	var p_prev Pos
	var from string
	var trkName string = trk.Name

	n := 0
	var trkSummary TrkSummary
	trkSummary.Name = trk.Name

	trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
	for i, trkpt := range trkpts {
		p := Pos{
			Lat: trkpt.Lat,
			Lon: trkpt.Lon,
			Ele: trkpt.Ele,
		}

		if i == 0 {
			p_prev = p
			if detail {
				from = "start"
				if trkpt.Name != nil {
					from = *trkpt.Name
				}
			} else {
				from = trk.Name
			}
			continue
		}

		eleDiff := DiffElevation(p_prev, p)
		denivPos += math.Max(eleDiff, 0)
		denivNeg += math.Min(eleDiff, 0)

		distance += Dist(p_prev, p)
		n += 1

		var x SectionInfo
		if (detail && trkpt.Name != nil) || (i == len(trkpts)-1) {
			x = SectionInfo{
				TrkName:        trkName,
				From:           from,
				NPoints:        n,
				VitessePlat:    vitessePlat,
				Distance:       distance,
				DenivPos:       denivPos,
				DenivNeg:       denivNeg,
				DistanceEffort: CalcDistanceEffort(distance, denivPos, denivNeg),
			}
			_, x.DurationHour, x.DurationMin = CalcDuration(x.DistanceEffort, vitessePlat)
		}

		if detail && trkpt.Name != nil {
			x.To = *trkpt.Name
			from = *trkpt.Name

			distance = 0
			denivPos = 0
			denivNeg = 0
			n = 0
		}

		if i == len(trkpts)-1 {
			x.To = "end"
		}

		if (detail && trkpt.Name != nil) || (i == len(trkpts)-1) {
			trkSummary.Section = append(trkSummary.Section, x)
		}

		p_prev = p
	}

	// trk.Extensions.DenivPos = denivPos
	// trk.Extensions.DenivNeg = denivNeg
	// trk.Extensions.Distance = distance

	return trkSummary
}

func (trk Trk) GetLonLat() ([]string, []string) {
	var lons, lats []string

	trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
	for _, trkpt := range trkpts {
		lons = append(lons, strconv.FormatFloat(trkpt.Lon, 'f', -1, 64))
		lats = append(lats, strconv.FormatFloat(trkpt.Lat, 'f', -1, 64))
	}

	return lons, lats
}

func (p_trk *Trk) Reverse() Trk {
	trk := *p_trk

	slices.Reverse(trk.Trkseg)
	for _, trkseg := range trk.Trkseg {
		slices.Reverse(trkseg.Trkpt)
	}
	return trk
}

func (trk Trk) GetElevations() []float64 {
	var trkpts []Trkpt = slices.Concat(trk.Trkseg)[0].Trkpt
	var elevs []float64
	for _, trkpt := range trkpts {
		elevs = append(elevs, trkpt.Ele)
	}
	return elevs
}
func (trk Trk) GetDistances() []float64 {
	var trkpts []Trkpt = slices.Concat(trk.Trkseg)[0].Trkpt
	var dists []float64
	posInit := Pos{
		Lon: trkpts[0].Lon,
		Lat: trkpts[0].Lat,
		Ele: trkpts[0].Ele,
	}
	var pos Pos
	for _, trkpt := range trkpts {
		pos = Pos{
			Lon: trkpt.Lon,
			Lat: trkpt.Lat,
			Ele: trkpt.Ele,
		}
		dists = append(dists,
			Dist(posInit, pos))
	}
	return dists
}

func (trk Trk) GetRollElevations(winSize int, calc RollCalc) []float64 {
	return Rolling(trk.GetElevations(), winSize, calc)
}
func (trk Trk) GetRollDistances(winSize int, calc RollCalc) []float64 {
	return Rolling(trk.GetDistances(), winSize, calc)
}

// func (trk Trk) Plot(filename string) {
// 	var elevs []float64 = trk.GetElevations()
// 	var rollmean []float64 = Rolling(elevs, 5, Mean)

// 	var elevsSummary []float64
// 	for _, v := range VariationSummary(rollmean) {
// 		elevsSummary = append(elevsSummary, v.Value)
// 	}

// 	var xys [][2][]float64
// 	var xs, ys []float64
// 	for i, v := range rollmean {
// 		ys = append(ys, v)
// 		xs = append(xs, float64(i))
// 	}
// 	xys = append(xys, [2][]float64{xs, ys})
// 	xs, ys = nil, nil
// 	for _, v := range VariationSummary(rollmean) {
// 		ys = append(ys, v.Value)
// 		xs = append(xs, float64(v.Index))
// 	}
// 	xys = append(xys, [2][]float64{xs, ys})

// 	// ==============

// 	// ys := [][]float64{elevsSummary, rollmean}
// 	// names := []string{"Estimation", "rollmean"}
// 	names := []string{"rollmean", "estimation"}
// 	colors := []color.RGBA{
// 		color.RGBA{R: 255, A: 255},
// 		color.RGBA{B: 255, A: 255},
// 	}
// 	// Plot(xys, names, colors)

// }
