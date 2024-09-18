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

	Name string
}

type SectionInfo struct {
	TrkName        string
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

type TrkSummary []SectionInfo

func (trk Trk) GetInfo(vitessePlat float64, detail bool) TrkSummary {
	var distance, denivPos, denivNeg float64 = 0, 0, 0

	var p_prev Pos
	var from string
	var trkName string = trk.Name

	n := 0
	var trkSummary TrkSummary

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
			trkSummary = append(trkSummary, x)
		}

		p_prev = p
	}

	// trk.Extensions.DenivPos = denivPos
	// trk.Extensions.DenivNeg = denivNeg
	// trk.Extensions.Distance = distance

	return trkSummary
}

func (s SectionInfo) Print(ascii_format ...bool) {
	if len(ascii_format) > 0 {
		if ascii_format[0] {
			fmt.Printf("\u001b[4mFrom:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", s.From)
		} else {
			fmt.Printf("From: %v\n", s.From)
		}
	} else {
		fmt.Printf("\u001b[4mFrom:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", s.From)
	}

	fmt.Printf("Number of points:       %v\n", s.NPoints)
	fmt.Printf("Distance:               %.1f km\n", s.Distance)
	fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", s.DenivPos, s.DenivNeg)
	fmt.Printf("Distance effort:        %.1f km\n", s.DistanceEffort)
	fmt.Printf("Vitesse sur plat:       %.1f km/h\n", s.VitessePlat)
	fmt.Printf("Temps parcours estimé:  %vh%v\n", s.DurationHour, s.DurationMin)
}

func (trkSummary TrkSummary) Print() {
	trkName := trkSummary[0].TrkName
	fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trkName)
	fmt.Println()
	for _, sectionInfo := range trkSummary {
		sectionInfo.Print()
		fmt.Println()
	}
}
