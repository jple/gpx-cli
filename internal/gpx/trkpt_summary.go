package gpx

import (
	"github.com/jple/gpx-cli/internal/geo"
	"github.com/jple/gpx-cli/internal/summary"
)

const RollingWindowSize = 10

func checkIndex(i, j, n int) {
	if i < 0 {
		panic("i must be > 0")
	} else if i > j {
		panic("i must be <= j")
	} else if j >= n {
		panic("j must be < n")
	}
}

func (trkpts Trkpts) TotalDistance() float64 {
	return geo.Sum(trkpts.Distances())
}

// Calculate cumulated distance between two index of trk
func (trkpts Trkpts) TotalDistanceBetweenIndex(i, j int) float64 {
	checkIndex(i, j, len(trkpts))
	return trkpts[i : j+1].TotalDistance()
}

// TODO: to be include in cmd
func (trkpts Trkpts) TotalDistanceBetweenNames(from, to string) float64 {
	i := trkpts.FindName(from)
	j := trkpts.FindName(to)
	return trkpts.TotalDistanceBetweenIndex(i, j)
}

func (trkpts Trkpts) TotalAscent(rollingWindowSize int) float64 {
	return geo.TotalAscent(geo.Rolling(trkpts.Elevations(), rollingWindowSize, geo.Mean))
}
func (trkpts Trkpts) TotalDescent(rollingWindowSize int) float64 {
	return geo.TotalDescent(geo.Rolling(trkpts.Elevations(), rollingWindowSize, geo.Mean))
}

func (trkpts Trkpts) Summary(flatSpeed float64) summary.TrkptsSummary {
	trkptsSummary := summary.TrkptsSummary{
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
