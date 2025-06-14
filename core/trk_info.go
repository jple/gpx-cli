package core

import (
	"slices"
)

// Trkpts is to contain Trkpt between two Trkpt.Name
// ie. from a the first Trkpt or a Trkpt.Name
// up to the last Trkpt without Trkpt.Name

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
func (trk Trk) GetInfo(trkid int, vitessePlat float64, detail bool) TrkSummary {
	listTrkpts := trk.GetListTrkpts()
	trkSummary := TrkSummary{Id: trkid, Name: trk.Name}

	var trackDuration float64
	for i, trkpts := range listTrkpts {
		// ============= Calculation geo info ============================
		trkptsSummary := trkpts.GetSummary(vitessePlat)
		trkptsSummary.TrkId = trkid
		trkptsSummary.TrkName = trk.Name

		// ============= Calculation From and To ============================
		// Set From with trk.Name, or first Trkpts name (depending on which available)
		if i == 0 && trk.Name != "" {
			trkptsSummary.From = trk.Name
		}
		if trkpts[0].Name != nil {
			trkptsSummary.From = *trkpts[0].Name
		}

		// Set To with next Trkpts name
		if i < len(listTrkpts)-1 && len(trkpts) > 0 {
			nextTrkpts := listTrkpts[i+1]
			trkptsSummary.To = *nextTrkpts[0].Name
		}

		// ============= Update trkSummary.ListTrkptsSummary ============================
		if detail {
			trkSummary.ListTrkptsSummary = append(trkSummary.ListTrkptsSummary, trkptsSummary)
		}

		// ============= Update trkSummary.Track ============================
		sectionDuration, _, _ := CalcDuration(trkptsSummary.DistanceEffort, vitessePlat)
		trackDuration += sectionDuration
		trkSummary.Track = TrkptsSummary{
			TrkName:     trk.Name,
			VitessePlat: vitessePlat,
			From:        trk.Name,
			To:          "end",

			NPoints:        trkSummary.Track.NPoints + trkptsSummary.NPoints,
			Distance:       trkSummary.Track.Distance + trkptsSummary.Distance,
			DenivPos:       trkSummary.Track.DenivPos + trkptsSummary.DenivPos,
			DenivNeg:       trkSummary.Track.DenivNeg + trkptsSummary.DenivNeg,
			DistanceEffort: trkSummary.Track.DistanceEffort + trkptsSummary.DistanceEffort,

			DurationHour: trkSummary.Track.DurationHour + trkptsSummary.DurationHour,
			DurationMin:  trkSummary.Track.DurationMin + trkptsSummary.DurationMin,
		}
		trkSummary.Track.DurationHour, trkSummary.Track.DurationMin = FloatToHourMin(trackDuration)
	}

	return trkSummary
}
