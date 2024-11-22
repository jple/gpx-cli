package core

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type PlotContent struct {
	Title  string
	X      plot.Axis
	Y      plot.Axis
	Graphs []Graph
}

type Graph struct {
	XYs   plotter.XYs
	Name  string
	Color color.RGBA
}

func ToPlotterXYs(X []float64, Y []float64) plotter.XYs {
	if len(X) != len(Y) {
		panic("X and Y don't have the same length")
	}

	var xys plotter.XYs
	n := len(X)
	for i := 0; i < n; i++ {
		xys = append(xys, plotter.XY{
			X: X[i], Y: Y[i]})
	}
	return xys

}
func CreateGraph(series any, name string, color color.RGBA) Graph {
	var xys plotter.XYs

	switch s := series.(type) {
	case []float64:
		for k, v := range s {
			xys = append(xys, plotter.XY{float64(k), v})
		}
	case Series:
		for _, v := range s {
			xys = append(xys, plotter.XY{float64(v.Index), v.Value})
		}
	case plotter.XYs:
		xys = s
	}

	return Graph{
		XYs:  xys,
		Name: name, Color: color,
	}
}

func Plot2(plotCt PlotContent, save_filename string) {
	p := plot.New()

	p.Title.Text = plotCt.Title
	p.X.Label.Text = plotCt.X.Label.Text
	p.Y.Label.Text = plotCt.Y.Label.Text
	p.Add(plotter.NewGrid())

	blank := color.RGBA{0, 0, 0, 0}

	for _, graph := range plotCt.Graphs {
		l, pt, err := plotter.NewLinePoints(graph.XYs)
		if err != nil {
			panic(err)
		}

		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		pt.Shape = draw.PyramidGlyph{}

		if graph.Color != blank {
			l.LineStyle.Color = graph.Color
			pt.Color = graph.Color
		} else {
			l.LineStyle.Color = color.RGBA{B: 255, A: 255}
			pt.Color = color.RGBA{B: 255, A: 255}
		}

		// lpLine, lpPoints, err := plotter.NewLinePoints(linePointsData)
		// if err != nil {
		// 	panic(err)
		// }
		// lpLine.Color = color.RGBA{G: 255, A: 255}
		// lpPoints.Shape = draw.PyramidGlyph{}
		// lpPoints.Color = color.RGBA{R: 255, A: 255}

		p.Add(l, pt)
		p.Legend.Add(graph.Name, l, pt)

		// err := plotutil.AddLinePoints(p, names[i], pts)
		// if err != nil {
		// 	panic(err)
		// }

	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, save_filename); err != nil {
		panic(err)
	}
}
