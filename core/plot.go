package core

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func CreateXYs(xy [2][]float64) plotter.XYs {
	// pts := make(plotter.XYs, len(xy))
	pts := make(plotter.XYs, len(xy[0]))
	for i, _ := range pts {
		// pts[i].X = float64(i)
		pts[i].X = xy[0][i]
		pts[i].Y = xy[1][i]
	}

	return pts
}
func Plot(xys [][2][]float64, names []string, colors ...[]color.RGBA) {

	if len(xys) != len(names) {
		fmt.Printf("ERROR: len(xys) (%v) != len(names) (%v)\n", len(xys), len(names))
		os.Exit(1)
	}
	if len(colors) > 0 && len(xys) != len(colors[0]) {
		fmt.Printf("ERROR: len(xys) (%v) != len(colors) (%v)\n", len(xys), len(colors))
		os.Exit(1)
	}

	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Elevation (m)"
	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// pts := CreateXYs(y)
	for i, _ := range xys {
		pts := CreateXYs(xys[i])

		l, pt, err := plotter.NewLinePoints(pts)
		if err != nil {
			panic(err)
		}
		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		pt.Shape = draw.PyramidGlyph{}
		if len(colors) > 0 {
			l.LineStyle.Color = colors[0][i]
			pt.Color = colors[0][i]
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
		p.Legend.Add(names[i], l, pt)

		// err := plotutil.AddLinePoints(p, names[i], pts)
		// if err != nil {
		// 	panic(err)
		// }

	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func Plot2() {
	rand.Seed(int64(0))

	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err := plotutil.AddLinePoints(p,
		"First", randomPoints(15),
		"Second", randomPoints(15),
		"Third", randomPoints(20))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}
