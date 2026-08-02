package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/internal/gpx"
	"github.com/jple/gpx-cli/internal/plot"
	"github.com/jple/gpx-cli/internal/trends"
)

func CreateTermPlotCmd() *cobra.Command {
	var trkId IntValue = 0
	flagsConf := []FlagConfig{
		{
			Name: "trk-id", Shortname: "t", DefaultValue: &trkId,
			Description: "Details about i-th trk. Value 0 will display all trk summary",
		},
	}

	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Display trends of a trk in the terminal",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.Parse(viper.GetString("filename"))
			trk := gpx.Trks[viper.GetInt("trk-id")]

			rawXY := plot.XY{trk.AllTrkpts().CumulativeDistances(), trk.AllTrkpts().Elevations()}
			rollXY := rawXY.CenterRollY(5)
			varSum := trends.BuildTrends(rollXY.Y, 30)

			var prevInd int
			var prevVal float64
			for k, v := range varSum {
				ind := v.Index
				val := v.Value
				dist := trk.AllTrkpts().TotalDistanceBetweenIndex(prevInd, ind)
				pct := (val - prevVal) / (dist * 1000) * 100
				if k > 0 {
					fmt.Printf("%.0f m\t--(%0.2f km)-->\t%.0f m \t(%.0f %%)\n", prevVal, dist, val, pct)
				}
				prevInd = ind
				prevVal = val
			}

			fmt.Println(varSum.ToString())

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
