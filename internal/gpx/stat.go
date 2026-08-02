package gpx

import (
	"io"
	"slices"

	"github.com/jple/gpx-cli/internal/geo"
	"github.com/jple/gpx-cli/internal/gpx/geostat"
)

// TODO: convert these into simple function
func (trkpts Trkpts) Stats(flatSpeed float64) geostat.GeoStatistic {
	stat := geostat.GeoStatistic{
		From: "start",
		To:   "end",

		N:            len(trkpts),
		Distance:     trkpts.TotalDistance(),
		TotalAscent:  trkpts.TotalAscent(RollingWindowSize),
		TotalDescent: trkpts.TotalDescent(RollingWindowSize),
	}

	// Add name to stat
	if trkpts[0].Name != nil {
		stat.Name = *trkpts[0].Name
		stat.From = *trkpts[0].Name
	}
	// TODO: stat.To should be set to nextTrk.Name if not existsd
	if trkpts[len(trkpts)-1].Name != nil {
		stat.To = *trkpts[len(trkpts)-1].Name
	}

	// Set calculation value
	stat.DistanceEffort = geo.CalcDistanceEffort(
		stat.Distance,
		stat.TotalAscent,
		stat.TotalDescent)
	_, stat.DurationHour, stat.DurationMin =
		geo.CalcDuration(
			stat.DistanceEffort,
			flatSpeed)

	return stat
}

func (trk Trk) Stats(flatSpeed float64) geostat.GeoStatistic {
	trkpts := trk.AllTrkpts()
	stat := trkpts.Stats(flatSpeed)

	if trkpts[0].Name != nil {
		stat.Name = *trkpts[0].Name
		stat.From = *trkpts[0].Name
	} else {
		stat.Name = trk.Name
		stat.From = trk.Name
	}

	// TODO: stat.To should be set to nextTrk.Name if not existsd
	if trkpts[len(trkpts)-1].Name != nil {
		stat.To = *trkpts[len(trkpts)-1].Name
	}

	return stat
}

func (trk Trk) StatsPerSection(flatSpeed float64) []geostat.GeoStatistic {
	stats := []geostat.GeoStatistic{}

	// Calculate stat for each "section" (trkpts)
	TrkSections := trk.SplitByTrkptName()
	for i, trkpts := range TrkSections {
		if len(trkpts) == 0 {
			continue
		}

		// ============= Calculate geo info ============================
		stat := trkpts.Stats(flatSpeed)
		// NOTE: by design, SplitByTrkptName is adding a last additional named item in trkpts
		// (see implementation). So N calculation needs to be corrected:
		if i < len(TrkSections)-1 {
			stat.N -= 1
		}

		// ============= Calculate values: From and To ============================
		// Set From with trk.Name, or the first Trkpts name (depending on which available)
		if i == 0 && trk.Name != "" {
			stat.From = trk.Name
		}
		if trkpts[0].Name != nil {
			stat.From = *trkpts[0].Name
		}

		// Set To with the last trkpts item (which is by design, a named element)
		if i < len(TrkSections)-1 && len(trkpts) > 0 {
			stat.To = *trkpts[len(trkpts)-1].Name
		}

		// ============= Update trkStat.PerSection ============================
		stats = append(stats, stat)

	}

	return stats
}

func ExportStatsTrks(w io.Writer, trks []Trk, flatspeed float64, detail bool, kind string,
	callback func(stat geostat.GeoStatistic, trkid int, isSection bool)) {

	exportOpt := geostat.Option{
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
		if callback != nil {
			callback(stat, i, false)
		}

		if detail {
			statsAllSections := trk.StatsPerSection(flatspeed)
			statsUsefulSections := slices.DeleteFunc(statsAllSections,
				func(stat geostat.GeoStatistic) bool {
					return stat.N <= 1
				})

			if len(statsUsefulSections) > 1 {
				for _, stat := range statsUsefulSections {
					stat.WriteValues(w, nil, kind, exportOpt, geostat.Option{Indent: 4})
					if callback != nil {
						callback(stat, i, true)
					}
				}
			}
		}

		if i == len(trks)-1 {
			stat.WriteFooter(w, kind)
		}
	}
}
