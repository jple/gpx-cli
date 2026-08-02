package gpx

import "github.com/jple/gpx-cli/internal/summary"

func (gpx Gpx) Summarize(flatSpeed float64) summary.GpxSummary {
	gpxSummary := summary.GpxSummary{FlatSpeed: flatSpeed}
	for i, trk := range gpx.Trks {
		trkSummary := trk.Summarize(i, flatSpeed)
		gpxSummary.TrkSummaries = append(gpxSummary.TrkSummaries, struct {
			TrkName string
			summary.TrkSummary
		}{trk.Name, trkSummary})
	}
	return gpxSummary
}

func (gpx Gpx) SummarizeBetweenTrkptsIndex(i1, i2 int, flatSpeed float64) summary.TrkptsSummary {
	trkpts := gpx.AllTrkpts()
	summary := trkpts[i1 : i2+1].Summary(flatSpeed)

	if trkpts[i1].Name != nil {
		summary.From = *trkpts[i1].Name
	}
	if trkpts[i2].Name != nil {
		summary.To = *trkpts[i2].Name
	}

	return summary
}

func (gpx Gpx) SummarizeBetweenTrkptsNames(name1, name2 string, flatSpeed float64) summary.TrkptsSummary {
	i1, err := gpx.FindTrkptsIdByName(name1)
	check(err)
	i2, err := gpx.FindTrkptsIdByName(name2)
	check(err)
	return gpx.SummarizeBetweenTrkptsIndex(i1, i2, flatSpeed)
}
