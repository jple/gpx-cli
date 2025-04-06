package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateTermPlotCmd() *cobra.Command {
	var trkId IntValue = 0
	flagsConf := []FlagConfig{
		{
			Name: "trk-id", Shortname: "i", DefaultValue: &trkId,
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

			// Printer(varSum)
			varSum.PrintTrends()
		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
