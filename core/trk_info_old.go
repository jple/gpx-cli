package core

import (
	"fmt"
	"math"
)

// TODO: NPoints buggy. see GetInfo on test/data/split.gpx
// 1. create [][]Trkpt containing []Trkpt until name (excluded)
// 2. calculate NPoints, ... on each []Trkpt

func (trk Trk) GetInfo_old(trkid int, vitessePlat float64, detail bool) TrkSummary {
	var distance, denivPos, denivNeg float64 = 0, 0, 0
	var p_prev Pos

	n := 0
	from, trkName := "start", trk.Name
	trkSummary := TrkSummary{Name: trk.Name}

	n_seg := len(trk.Trkseg)
	var trackDuration, sectDuration float64 = 0, 0
	for i, trkseg := range trk.Trkseg {
		n_pt := len(trkseg.Trkpt)
		// WIP: fix when trkseg is empty
		if n_pt == 0 {
			fmt.Println("hell")
			continue
		}

		fmt.Printf("seg ")
		for j, trkpt := range trkseg.Trkpt {
			fmt.Printf("pt ")
			p := Pos{
				Lat: trkpt.Lat,
				Lon: trkpt.Lon,
				Ele: trkpt.Ele,
			}
			n += 1
			fmt.Printf("+1 ")

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
			// fmt.Printf("+%v ", eleDiff)
			fmt.Printf("+%v ", denivPos)

			var section TrkptsSummary
			// Create trkpts from this iteration, if any of these conditions are fulfilled
			// - detail is asked, and current trkpt contains a name
			// - last trkpt of the trk
			isLastTrkpt := (i == n_seg-1) && (j == n_pt-1)
			if (detail && trkpt.Name != nil) || isLastTrkpt {
				fmt.Printf("NS\n")
				section = TrkptsSummary{
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
				sectDuration, section.DurationHour, section.DurationMin = CalcDuration(section.DistanceEffort, vitessePlat)

				// Reset cumulative values
				distance = 0
				denivPos = 0
				denivNeg = 0
				n = 0

				trackDuration += sectDuration
			}

			if detail && trkpt.Name != nil {
				section.To = *trkpt.Name

				// Set value for next iteration
				from = *trkpt.Name
			}

			// if isLastTrkpt {
			// 	section.To = "end"
			// }

			if (detail && trkpt.Name != nil) || isLastTrkpt {
				// Append new section
				trkSummary.ListTrkptsSummary = append(trkSummary.ListTrkptsSummary, section)

				// Update whole track values
				trkSummary.Track = TrkptsSummary{
					TrkName:     trkName,
					VitessePlat: vitessePlat,
					From:        trkName,
					To:          "end",

					NPoints:        trkSummary.Track.NPoints + section.NPoints,
					Distance:       trkSummary.Track.Distance + section.Distance,
					DenivPos:       trkSummary.Track.DenivPos + section.DenivPos,
					DenivNeg:       trkSummary.Track.DenivNeg + section.DenivNeg,
					DistanceEffort: trkSummary.Track.DistanceEffort + section.DistanceEffort,

					DurationHour: trkSummary.Track.DurationHour + section.DurationHour,
					DurationMin:  trkSummary.Track.DurationMin + section.DurationMin,
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
