package gpx

import (
	"io"
	"slices"

	"github.com/jple/gpx-cli/internal/geo"
	"github.com/jple/gpx-cli/internal/gpx/summary"
)

// TODO: convert these into simple function
func (trkpts Trkpts) Stats(flatSpeed float64) summary.GeoStatistic {
	trkptsSummary := summary.GeoStatistic{
		From: "start",
		To:   "end",

		N:            len(trkpts),
		Distance:     trkpts.TotalDistance(),
		TotalAscent:  trkpts.TotalAscent(RollingWindowSize),
		TotalDescent: trkpts.TotalDescent(RollingWindowSize),
	}

	// Set calculation value
	trkptsSummary.DistanceEffort = geo.CalcDistanceEffort(
		trkptsSummary.Distance,
		trkptsSummary.TotalAscent,
		trkptsSummary.TotalDescent)
	_, trkptsSummary.DurationHour, trkptsSummary.DurationMin =
		geo.CalcDuration(
			trkptsSummary.DistanceEffort,
			flatSpeed)

	return trkptsSummary
}

func (trk Trk) Stats(flatSpeed float64) summary.GeoStatistic {
	trkpts := trk.AllTrkpts()
	summary := trkpts.Stats(flatSpeed)

	if trkpts[0].Name != nil {
		summary.Name = *trkpts[0].Name
		summary.From = *trkpts[0].Name
	} else {
		summary.Name = trk.Name
		summary.From = trk.Name
	}

	// TODO: summary.To should be set to nextTrk.Name if not existsd
	if trkpts[len(trkpts)-1].Name != nil {
		summary.To = *trkpts[len(trkpts)-1].Name
	}

	return summary
}

func (trk Trk) StatsPerSection(flatSpeed float64) []summary.GeoStatistic {
	summaries := []summary.GeoStatistic{}

	// Calculate summary for each "section" (trkpts)
	TrkSections := trk.SplitByTrkptName()
	for i, trkpts := range TrkSections {
		if len(trkpts) == 0 {
			continue
		}

		// ============= Calculate geo info ============================
		summary := trkpts.Stats(flatSpeed)
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
func ExportStatsTrks(w io.Writer, trks []Trk, flatspeed float64, detail bool, kind string) {
	exportOpt := summary.Option{
		WithDestination:   true,
		WithDistance:      true,
		WithAscentDescent: true,
		WithDuration:      true,
	}
	for i, trk := range trks {
		stat := trk.Stats(flatspeed)

		if i == 0 {
			stat.WriteHeader(w, kind, exportOpt)
		}

		stat.WriteValues(w, &i, kind, exportOpt)

		if detail {
			statsAllSections := trk.StatsPerSection(flatspeed)
			statsUsefulSections := slices.DeleteFunc(statsAllSections,
				func(stat summary.GeoStatistic) bool {
					return stat.N <= 1
				})

			if len(statsUsefulSections) > 1 {
				for _, stat := range statsUsefulSections {
					stat.WriteValues(w, nil, kind, exportOpt, summary.Option{Indent: 4})
				}
			}
		}

		if i == len(trks)-1 {
			stat.WriteFooter(w, kind)
		}
	}
}
