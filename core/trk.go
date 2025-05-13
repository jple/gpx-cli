package core

import (
	"fmt"
	"math"
	"slices"
	"strconv"
)

func (trk Trk) GetInfo(trkid int, vitessePlat float64, detail bool) TrkSummary {
	var distance, denivPos, denivNeg float64 = 0, 0, 0
	var p_prev Pos

	n := 0
	from, trkName := "start", trk.Name
	trkSummary := TrkSummary{Name: trk.Name}

	n_seg := len(trk.Trkseg)
	var trackDuration, sectDuration float64 = 0, 0
	for i, trkseg := range trk.Trkseg {
		n_pt := len(trkseg.Trkpt)
		for j, trkpt := range trkseg.Trkpt {
			p := Pos{
				Lat: trkpt.Lat,
				Lon: trkpt.Lon,
				Ele: trkpt.Ele,
			}

			// Init variables on first trkpt of trk
			if i == 0 && j == 0 {
				p_prev = p
				if detail {
					if trkpt.Name != nil {
						from = *trkpt.Name
					}
				} else {
					from = trk.Name
				}
				continue
			}

			// Update cumulative values
			eleDiff := DiffElevation(p_prev, p)
			denivPos += math.Max(eleDiff, 0)
			denivNeg += math.Min(eleDiff, 0)
			distance += Dist(p_prev, p)
			n += 1

			var x SectionInfo
			// Create section from this iteration, if any of these conditions are fulfilled
			// - detail is asked, and current trkpt contains a name
			// - last trkpt of the trk
			isLastTrkpt := (i == n_seg-1) && (j == n_pt-1)
			if (detail && trkpt.Name != nil) || isLastTrkpt {
				x = SectionInfo{
					TrkId:       trkid,
					TrkName:     trkName,
					VitessePlat: vitessePlat,

					From:         from,
					To:           "end", // is updated later if required
					FromTrksegId: &i,
					FromTrkptId:  &j,

					NPoints:        n,
					Distance:       distance,
					DenivPos:       denivPos,
					DenivNeg:       denivNeg,
					DistanceEffort: CalcDistanceEffort(distance, denivPos, denivNeg),
				}
				sectDuration, x.DurationHour, x.DurationMin = CalcDuration(x.DistanceEffort, vitessePlat)

				// Reset cumulative values
				distance = 0
				denivPos = 0
				denivNeg = 0
				n = 0

				trackDuration += sectDuration
			}

			if detail && trkpt.Name != nil {
				x.To = *trkpt.Name

				// Set value for next iteration
				from = *trkpt.Name
			}

			// if isLastTrkpt {
			// 	x.To = "end"
			// }

			// if (detail && trkpt.Name != nil) || isLastTrkpt {
			// 	trkSummary.Section = append(trkSummary.Section, x)
			// }

			if (detail && trkpt.Name != nil) || isLastTrkpt {
				// Append new section
				trkSummary.Section = append(trkSummary.Section, x)

				// Update whole track values
				trkSummary.Track = SectionInfo{
					TrkName:     trkName,
					VitessePlat: vitessePlat,
					From:        trkName,
					To:          "end",

					NPoints:        trkSummary.Track.NPoints + x.NPoints,
					Distance:       trkSummary.Track.Distance + x.Distance,
					DenivPos:       trkSummary.Track.DenivPos + x.DenivPos,
					DenivNeg:       trkSummary.Track.DenivNeg + x.DenivNeg,
					DistanceEffort: trkSummary.Track.DistanceEffort + x.DistanceEffort,

					DurationHour: trkSummary.Track.DurationHour + x.DurationHour,
					DurationMin:  trkSummary.Track.DurationMin + x.DurationMin,
				}

				trkSummary.Track.DurationHour, trkSummary.Track.DurationMin = FloatToHourMin(trackDuration)
			}

			// Reset previous pos
			p_prev = p
		}
	}

	// trk.Extensions.DenivPos = denivPos
	// trk.Extensions.DenivNeg = denivNeg
	// trk.Extensions.Distance = distance

	return trkSummary
}

func (trk Trk) GetLonLat() ([]string, []string) {
	var lons, lats []string

	// trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
	var trkpts []Trkpt
	for _, trkseg := range trk.Trkseg {
		trkpts = slices.Concat(trkpts, trkseg.Trkpt)
	}
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
	// var trkpts []Trkpt = slices.Concat(trk.Trkseg)[0].Trkpt
	var trkpts []Trkpt
	for _, trkseg := range trk.Trkseg {
		trkpts = slices.Concat(trkpts, trkseg.Trkpt)
	}
	var elevs []float64
	for _, trkpt := range trkpts {
		elevs = append(elevs, trkpt.Ele)
	}
	return elevs
}

// Calculate cumulated distance between two index of trk
func (trk Trk) GetDistanceFromTo(i, j int) float64 {
	if i >= j {
		fmt.Println("i must be < j")
		return 0.0
	}
	// var trkpts []Trkpt = slices.Concat(trk.Trkseg)[0].Trkpt
	var trkpts []Trkpt
	for _, trkseg := range trk.Trkseg {
		trkpts = slices.Concat(trkpts, trkseg.Trkpt)
	}
	var dist float64
	posPrev := Pos{
		Lon: trkpts[i].Lon,
		Lat: trkpts[i].Lat,
		Ele: trkpts[i].Ele,
	}
	for k, trkpt := range trkpts {
		if k <= i {
			continue
		}
		if k >= j {
			break
		}

		pos := Pos{
			Lon: trkpt.Lon,
			Lat: trkpt.Lat,
			Ele: trkpt.Ele,
		}
		dist += Dist(posPrev, pos)
		posPrev = pos
	}
	return dist
}

// Caculate cumulated distance for each trkpt
// (distance between trkpt[0] and trkpt[i])
func (trk Trk) GetDistanceEachTrkpts() []float64 {
	// var trkpts []Trkpt = slices.Concat(trk.Trkseg)[0].Trkpt
	var trkpts []Trkpt
	for _, trkseg := range trk.Trkseg {
		trkpts = slices.Concat(trkpts, trkseg.Trkpt)
	}
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
	return Rolling(trk.GetDistanceEachTrkpts(), winSize, calc)
}

func (trk *Trk) AddName(name string) {
	trk.Name = name
}

// TODO: create generics for AddName
func (trkpt *Trkpt) AddName(name string) {
	trkpt.Name = &name
}
func (trkpt *Trkpt) AddElevation(ele float64) {
	trkpt.Ele = ele
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
