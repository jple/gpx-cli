package core

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Graph struct {
	XYs  plotter.XYs
	Name string

	LineStyle  *draw.LineStyle
	GlyphStyle *draw.GlyphStyle
}

func CreateGraph(series any, name string, lstyle *draw.LineStyle, ptstyle *draw.GlyphStyle) Graph {
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
	case XY:
		xys = CreatePlotterXYs(s)
	case plotter.XYs:
		xys = s
	}

	return Graph{
		XYs:        xys,
		Name:       name,
		LineStyle:  lstyle,
		GlyphStyle: ptstyle,
	}
}

func NewPlot(title, xlabel, ylabel string, graphs []Graph) *plot.Plot {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel
	p.Add(plotter.NewGrid())

	for _, graph := range graphs {
		l, pt, err := plotter.NewLinePoints(graph.XYs)
		if err != nil {
			panic(err)
		}

		// Style line
		if graph.LineStyle != nil {
			l.LineStyle = *graph.LineStyle
		} else {
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		}
		// Style dots
		if graph.GlyphStyle != nil {
			pt.GlyphStyle = *graph.GlyphStyle
		} else {
			pt.Shape = draw.PyramidGlyph{}
		}

		p.Add(l, pt)
		p.Legend.Add(graph.Name, l, pt)
	}

	return p
}
