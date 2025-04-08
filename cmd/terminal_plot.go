package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
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
		Use:   "term-plot",
		Short: "WIP new command",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))
			trk := gpx.Trk[viper.GetInt("trk-id")]

			rollmean := trk.GetRollElevations(5, Mean)
			varSum := TrendSummary(rollmean)

			var prevInd int
			var prevVal float64
			for k, v := range varSum {
				ind := v.Index
				val := v.Value
				dist := trk.GetDistanceFromTo(prevInd, ind)
				pct := (val - prevVal) / (dist * 1000) * 100
				if k > 0 {
					fmt.Printf("%.0f m\t--(%0.2f km)-->\t%.0f m \t(%.0f %%)\n", prevVal, dist, val, pct)
				}
				prevInd = ind
				prevVal = val

			}

			varSum.PrintTrends()

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
