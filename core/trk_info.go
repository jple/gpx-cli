package core

import (
	"slices"
)

// In Trk perspective, Trkpts is to contain Trkpt between two Trkpt.Name
// ie. from a the first Trkpt or a Trkpt.Name
// up to the last Trkpt without Trkpt.Name
func (trk Trk) GetTrkpts() Trkpts {
	trkpts := Trkpts{}
	for _, trkseg := range trk.Trksegs {
		trkpts = slices.Concat(trkpts, trkseg.Trkpts)
	}
	return trkpts
}

// GetListTrkptsPerName creates ListTrkpts ([]Trkpts, ie. [][]Trkpt)
// where each Trkpts is a []Trkpt from named item to next named item (included)
// NOTE: the next named item is included in order to correctly calculate summary of Trkpts
// Without it, the calculation would skip the calculation between ListTrkpts[i].Trkpts[-1] and
// ListTrkpts[I+].Trkpts[0]
func (trk Trk) GetListTrkptsPerName() ListTrkpts {
	currentTrkpts := Trkpts{}
	listTrkpts := ListTrkpts{}

	for i, trkpt := range trk.GetTrkpts() {
		currentTrkpts = append(currentTrkpts, trkpt)

		// Reach a new named pt (which is not first element)
		// This is the end of the current "section" (currentTrkpts)
		if trkpt.Name != nil && i > 0 {
			// Append "section" to listTrkpts
			listTrkpts = append(listTrkpts, currentTrkpts)
			// Prepare the new section, with the named pt
			currentTrkpts = Trkpts{trkpt}
		}
	}
	// Add last "section" (trkpts) into the list
	listTrkpts = append(listTrkpts, currentTrkpts)
	return listTrkpts
}

func (trk Trk) GetInfoWholeTrack(vitessePlat float64) TrkptsSummary {
	return trk.GetTrkpts().GetSummary(vitessePlat)
}

func (trk Trk) GetInfoPerSection(vitessePlat float64) []TrkptsSummary {
	listTrkpts := trk.GetListTrkptsPerName()
	listTrkptsSummary := []TrkptsSummary{}

	// Calculate summary for each "section" (trkpts)
	for i, trkpts := range listTrkpts {
		if len(trkpts) == 0 {
			continue
		}

		// ============= Calculate geo info ============================
		trkptsSummary := trkpts.GetSummary(vitessePlat)
		// NOTE: by design, GetListTrkptsPerName is adding a last additional named item in trkpts
		// (see implementation). So NPoints calculation needs to be corrected:
		if i < len(listTrkpts)-1 {
			trkptsSummary.NPoints -= 1
		}

		// ============= Calculate values: From and To ============================
		// Set From with trk.Name, or the first Trkpts name (depending on which available)
		if i == 0 && trk.Name != "" {
			trkptsSummary.From = trk.Name
		}
		if trkpts[0].Name != nil {
			trkptsSummary.From = *trkpts[0].Name
		}

		// Set To with the last trkpts item (which is by design, a named element)
		if i < len(listTrkpts)-1 && len(trkpts) > 0 {
			trkptsSummary.To = *trkpts[len(trkpts)-1].Name
		}

		// ============= Update trkSummary.ListTrkptsSummary ============================
		listTrkptsSummary = append(listTrkptsSummary, trkptsSummary)

	}

	return listTrkptsSummary
}

func (trk Trk) GetInfo(trkid int, vitessePlat float64) TrkSummary {
	return TrkSummary{
		ListTrkptsSummary: trk.GetInfoPerSection(vitessePlat),
		Track:             trk.GetInfoWholeTrack(vitessePlat),
	}
}
