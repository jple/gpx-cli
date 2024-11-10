package cmd

import (
	"image/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gonum.org/v1/plot"

	. "github.com/jple/gpx-cli/core"
)

func CreatePlotCmd() *cobra.Command {
	var trkId IntValue = 0
	flagsConf := []FlagConfig{
		{
			Name: "trk-id", Shortname: "i", DefaultValue: &trkId,
			Description: "Plot elevations about i-th trk.",
		},
	}

	cmd := &cobra.Command{
		Use:   "plot",
		Short: "Plot Trk into png file",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.ParseFile(gpx.Filepath)
			trk := gpx.Trk[viper.GetInt("trk-id")]

			rollElev := trk.GetRollElevations(5, Mean)
			rollDist := trk.GetRollDistances(5, Mean)

			xys := ToPlotterXYs(rollDist, rollElev)

			// =============
			var plotCt PlotContent
			plotCt.Title = trk.Name
			plotCt.X = plot.Axis{}
			plotCt.X.Label.Text = "Distance (km)"
			plotCt.Y.Label.Text = "Elevation (m)"
			plotCt.Graphs = append(plotCt.Graphs,
				CreateGraph(
					xys,
					"elev / dist",
					color.RGBA{R: 255, A: 255}),
			)
			Plot2(plotCt, "elev_dist.png")

			// =============
			plotCt.Title = trk.Name
			plotCt.X = plot.Axis{}
			plotCt.X.Label.Text = "index"
			plotCt.Y.Label.Text = "Elevation (m)"
			plotCt.Graphs = nil
			plotCt.Graphs = append(plotCt.Graphs,
				CreateGraph(
					rollElev,
					"rollmean",
					color.RGBA{R: 255, A: 255}),
				CreateGraph(
					VariationSummary(rollElev),
					"estimation",
					color.RGBA{B: 255, A: 255}),
			)
			Plot2(plotCt, "points.png")
		},
	}
	initFlags(cmd, flagsConf)

	return cmd
}
