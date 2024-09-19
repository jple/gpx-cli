package core

import (
	"math"
	"slices"
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
		denivPos += math.Max(eleDiff, 1)
		denivNeg += math.Min(eleDiff, -1)

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
