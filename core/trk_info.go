package core

import (
	"slices"
)

// Trkpts is to contain Trkpt between two Trkpt.Name
// ie. from a the first Trkpt or a Trkpt.Name
// up to the last Trkpt without Trkpt.Name
// TODO/refacto: create generic Trkpts type, and move math_trkpts.go into methods

func (trk Trk) GetFlattenTrkpts() Trkpts {
	trkpts := Trkpts{}
	for _, trkseg := range trk.Trkseg {
		trkpts = slices.Concat(trkpts, trkseg.Trkpt)
	}
	return trkpts
}

func (trk Trk) GetListTrkpts() ListTrkpts {
	return trk.GetFlattenTrkptsSplitName()
}

// GetFlattenTrkptsSplitName creates ListTrkpts ([][]Trkpt) containing Trkpts ([]Trkpt until name (excluded))
func (trk Trk) GetFlattenTrkptsSplitName() ListTrkpts {
	sections := ListTrkpts{}
	s := Trkpts{}
	trkpts := trk.GetFlattenTrkpts()
	for i, trkpt := range trkpts {
		if i == 0 || trkpt.Name == nil {
			// Update current trkpts
			s = append(s, trkpt)
		} else { // End trkpts
			// Update list of sections, starting new trkpts
			sections = append(sections, s)
			s = Trkpts{trkpt}
		}
	}
	sections = append(sections, s) // add last trkpts
	return sections
}

// WIP: refacto GetInfo
func (trk Trk) GetInfo2(trkid int, vitessePlat float64, detail bool) TrkSummary {
	listTrkpts := trk.GetListTrkpts()
	trkSummary := TrkSummary{Id: trkid, Name: trk.Name}
	sectionSummary := SectionSummary{
		TrkId:       trkid,
		TrkName:     trk.Name,
		VitessePlat: vitessePlat,

		// Cumulative values between "From" and "To"
		// NPoints
		// Distance       float64
		// DenivPos       float64
		// DenivNeg       float64
		// DistanceEffort float64
		// DurationHour   int8
		// DurationMin    int8
	}
	for _, trkpts := range listTrkpts {
		sectionSummary.NPoints = len(trkpts)
		sectionSummary.Distance = trkpts.GetTotalDistance()
		sectionSummary.DenivPos = trkpts.GetTotalAscent()
		sectionSummary.DenivNeg = trkpts.GetTotalDescent()

		trkSummary.ListSectionSummary = append(trkSummary.ListSectionSummary, sectionSummary)
	}
	return trkSummary
}
