package cmd

import (
	"fmt"
	"image/color"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/jple/gpx-cli/internal/geo"
	. "github.com/jple/gpx-cli/internal/gpx"
	"github.com/jple/gpx-cli/internal/plot"
)

func CreatePlotCmd() *cobra.Command {
	var trkId IntValue = 0
	flagsConf := []FlagConfig{
		{
			Name: "trk-id", Shortname: "t", DefaultValue: &trkId,
			Description: "Plot elevations about i-th trk.",
		},
	}

	cmd := &cobra.Command{
		Use:   "plot",
		Short: "Plot Trk into png file",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
			if !viper.IsSet("output") {
				viper.Set("output", "profil.png")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.Parse(viper.GetString("filename"))
			trk := gpx.Trks[viper.GetInt("trk-id")]

			rawXY := plot.XY{trk.AllTrkpts().CumulativeDistances(), trk.AllTrkpts().Elevations()}
			rollXY := rawXY.CenterRollY(5)
			trendXY := rollXY.TrendY(2)

			// Create graph
			graphs := []plot.Graph{
				plot.CreateGraph(
					rawXY,
					"raw",
					&draw.LineStyle{
						Width: vg.Points(5),
						Color: color.RGBA{B: 255, A: 255},
					},
					&draw.GlyphStyle{},
				),

				plot.CreateGraph(
					rollXY,
					"rolling mean",
					&draw.LineStyle{
						Width:  vg.Points(3),
						Dashes: []vg.Length{vg.Points(20), vg.Points(20)},
						Color:  color.RGBA{255, 200, 200, 255},
					},
					&draw.GlyphStyle{},
				),

				plot.CreateGraph(
					trendXY,
					"trend",
					&draw.LineStyle{
						Width:  vg.Points(2),
						Dashes: []vg.Length{vg.Points(5), vg.Points(3)},
						Color:  color.RGBA{R: 155, G: 155, A: 255},
					},
					&draw.GlyphStyle{
						Color:  color.RGBA{R: 155, G: 155, A: 255},
						Radius: 5,
						Shape:  draw.PyramidGlyph{},
					},
				),
			}

			// Calculate data for subtitle
			XYs := []struct {
				Name string
				This plot.XY
			}{
				{"raw", rawXY},
				{"roll", rollXY},
				{"trend", trendXY},
			}
			var subtitle string
			for _, xy := range XYs {
				subtitle += fmt.Sprintf("[%v] min: %.0f | max: %.0f | ascents: %.0f | descents: %.0f\n",
					xy.Name, slices.Min(xy.This.Y), slices.Max(xy.This.Y), geo.TotalAscent(xy.This.Y), geo.TotalDescent(xy.This.Y))
			}
			p := plot.NewPlot(
				trk.Name+"\n"+subtitle,
				"Distance (km)", "Elevation (m)",
				graphs)
			if err := p.Save(8*vg.Inch, 6*vg.Inch, viper.GetString("output")); err != nil {
				panic(err)
			}
		},
	}
	initFlags(cmd, flagsConf)

	return cmd
}
