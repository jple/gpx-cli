package core

import (
	"slices"
)

// Trkpts is to contain Trkpt between two Trkpt.Name
// ie. from a the first Trkpt or a Trkpt.Name
// up to the last Trkpt without Trkpt.Name
// TODO/refacto: create generic Trkpts type, and move math_trkpts.go into methods

func (trk Trk) GetTrkpts() Trkpts {
	trkpts := Trkpts{}
	for _, trkseg := range trk.Trkseg {
		trkpts = slices.Concat(trkpts, trkseg.Trkpt)
	}
	return trkpts
}

// This is just an alias to GetListTrkptsPerName
func (trk Trk) GetListTrkpts() ListTrkpts {
	return trk.GetListTrkptsPerName()
}

// GetListTrkptsPerName creates ListTrkpts ([][]Trkpt) containing Trkpts ([]Trkpt until name (excluded))
func (trk Trk) GetListTrkptsPerName() ListTrkpts {
	trkpts := Trkpts{}
	listTrkpts := ListTrkpts{}

	for i, trkpt := range trk.GetTrkpts() {
		if i == 0 || trkpt.Name == nil {
			// Update current trkpts
			trkpts = append(trkpts, trkpt)
		} else { // End trkpts, or trkpt has Name
			// Update list of Trkpts, starting new trkpts
			listTrkpts = append(listTrkpts, trkpts)
			trkpts = Trkpts{trkpt}
		}
	}
	listTrkpts = append(listTrkpts, trkpts) // add last trkpts
	return listTrkpts
}

// WIP: refacto GetInfo
// NOTE: detail is not used. Should be removed ?
func (trk Trk) GetInfo2(trkid int, vitessePlat float64, detail bool) TrkSummary {
	listTrkpts := trk.GetListTrkpts()
	trkSummary := TrkSummary{Id: trkid, Name: trk.Name}

	var sectionDuration, trackDuration float64
	for i, trkpts := range listTrkpts {
		// ============= Calculation geo info ============================
		sectionSummary := SectionSummary{
			TrkId:       trkid,
			TrkName:     trk.Name,
			VitessePlat: vitessePlat,
			From:        "start",
			To:          "end"}

		// Set calculation value
		sectionSummary.NPoints = len(trkpts)
		sectionSummary.Distance = trkpts.GetTotalDistance()
		sectionSummary.DenivPos = trkpts.GetTotalAscent()
		sectionSummary.DenivNeg = trkpts.GetTotalDescent()
		sectionSummary.DistanceEffort = CalcDistanceEffort(
			sectionSummary.Distance,
			sectionSummary.DenivPos,
			sectionSummary.DenivNeg)
		sectionDuration, sectionSummary.DurationHour, sectionSummary.DurationMin =
			CalcDuration(
				sectionSummary.DistanceEffort,
				vitessePlat)

		// ============= Calculation From and To ============================
		// Set From with trk.Name, or first Trkpts name (depending on which available)
		if i == 0 && trk.Name != "" {
			sectionSummary.From = trk.Name
		}
		sectionSummary.From = *trkpts[0].Name

		// Set To with next Trkpts name
		if i < len(listTrkpts)-1 && len(trkpts) > 0 {
			nextTrkpts := listTrkpts[i+1]
			sectionSummary.To = *nextTrkpts[0].Name
		}

		// ============= Update trkSummary.ListSectionSummary ============================
		trkSummary.ListSectionSummary = append(trkSummary.ListSectionSummary, sectionSummary)

		// ============= Update trkSummary.Track ============================
		trackDuration += sectionDuration
		trkSummary.Track = SectionSummary{
			TrkName:     trk.Name,
			VitessePlat: vitessePlat,
			From:        trk.Name,
			To:          "end",

			NPoints:        trkSummary.Track.NPoints + sectionSummary.NPoints,
			Distance:       trkSummary.Track.Distance + sectionSummary.Distance,
			DenivPos:       trkSummary.Track.DenivPos + sectionSummary.DenivPos,
			DenivNeg:       trkSummary.Track.DenivNeg + sectionSummary.DenivNeg,
			DistanceEffort: trkSummary.Track.DistanceEffort + sectionSummary.DistanceEffort,

			DurationHour: trkSummary.Track.DurationHour + sectionSummary.DurationHour,
			DurationMin:  trkSummary.Track.DurationMin + sectionSummary.DurationMin,
		}
		trkSummary.Track.DurationHour, trkSummary.Track.DurationMin = FloatToHourMin(trackDuration)
	}

	return trkSummary
}
