package gpx

import "github.com/jple/gpx-cli/internal/summary"

func (trk Trk) SummaryAllSections(flatSpeed float64) summary.TrkptsSummary {
	return trk.AllTrkpts().Summary(flatSpeed)
}

func (trk Trk) SummarizePerSection(flatSpeed float64) []summary.TrkptsSummary {
	TrkSections := trk.SplitByTrkptName()
	summaries := []summary.TrkptsSummary{}

	// Calculate summary for each "section" (trkpts)
	for i, trkpts := range TrkSections {
		if len(trkpts) == 0 {
			continue
		}

		// ============= Calculate geo info ============================
		summary := trkpts.Summary(flatSpeed)
		// NOTE: by design, SplitByTrkptName is adding a last additional named item in trkpts
		// (see implementation). So N calculation needs to be corrected:
		if i < len(TrkSections)-1 {
			summary.N -= 1
		}

		// ============= Calculate values: From and To ============================
		// Set From with trk.Name, or the first Trkpts name (depending on which available)
		if i == 0 && trk.Name != "" {
			summary.From = trk.Name
		}
		if trkpts[0].Name != nil {
			summary.From = *trkpts[0].Name
		}

		// Set To with the last trkpts item (which is by design, a named element)
		if i < len(TrkSections)-1 && len(trkpts) > 0 {
			summary.To = *trkpts[len(trkpts)-1].Name
		}

		// ============= Update trkSummary.PerSection ============================
		summaries = append(summaries, summary)

	}

	return summaries
}

func (trk Trk) Summarize(trkid int, flatSpeed float64) summary.TrkSummary {
	return summary.TrkSummary{
		PerSection:  trk.SummarizePerSection(flatSpeed),
		AllSections: trk.SummaryAllSections(flatSpeed),
	}
}
