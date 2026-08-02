package plot

import (
	"math"
	"slices"

	"github.com/jple/gpx-cli/internal/geo"
	"github.com/jple/gpx-cli/internal/trends"
	"gonum.org/v1/plot/plotter"
)

// XY is representing X and Y values to be plot on a Graph
type XY struct {
	X, Y []float64
}

func (xy XY) CheckLen() {
	if len(xy.X) != len(xy.Y) {
		// TODO: shoudl return an error
		panic("X and Y don't have the same length")
	}
}
func (xy XY) GetLen() int {
	return len(xy.X)
}

// NOTE: X and Y does not have the same length
func (xy XY) RollY(winsize int) XY {
	return XY{
		X: xy.X,
		Y: geo.Rolling(xy.Y, winsize, geo.Mean),
	}
}

func (xy XY) CenterRollY(winsize int) XY {
	// recenter adds elements to make Y same length as X
	recenter := func(s []float64, winsize int) []float64 {
		half := (float64(winsize) - 1) / 2.0
		nbef := int(math.Round(half))
		naft := int(math.Floor(half))

		return slices.Concat(
			slices.Repeat([]float64{s[0]}, nbef),
			s,
			slices.Repeat([]float64{s[len(s)-1]}, naft),
		)
	}

	return XY{
		X: xy.X,
		Y: recenter(geo.Rolling(xy.Y, winsize, geo.Mean), winsize),
	}
}

func (xy XY) TrendY(threshold float64) XY {
	var xtrend, ytrend []float64
	for _, trend := range trends.BuildTrends(xy.Y, threshold) {
		xtrend = append(xtrend, xy.X[trend.Index])
		ytrend = append(ytrend, trend.Value)
	}

	return XY{xtrend, ytrend}
}

func CreatePlotterXYs(xy XY) plotter.XYs {
	xy.CheckLen()

	var out plotter.XYs
	n := xy.GetLen()
	for i := 0; i < n; i++ {
		out = append(
			out,
			plotter.XY{X: xy.X[i], Y: xy.Y[i]})
	}
	return out
}
