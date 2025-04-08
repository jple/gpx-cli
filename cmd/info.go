package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func StringPointer(s string) *string {
	return &s
}

func CreateInfoCmd() *cobra.Command {
	var trkId IntValue = -1
	flagsConf := []FlagConfig{
		{
			Name: "trk-id", Shortname: "t", DefaultValue: &trkId,
			Description: "Details about i-th trk. Value -1 will display all trk summary",
		},
	}

	cmd := &cobra.Command{
		Use:   "info",
		Short: "General info on the tracks",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))
			gpx.SetVitesse(4.5)

			var printArgs PrintArgs = PrintArgs{AsciiFormat: true}

			if viper.GetInt("trk-id") != -1 {
				printArgs.PrintFrom = true

				fmt.Printf("[%v] ", viper.GetInt("trk-id"))
				gpx.
					Trk[viper.GetInt("trk-id")].
					GetInfo(gpx.Extensions.Vitesse, true).
					Print(printArgs)
			} else {
				printArgs.PrintFrom = false

				gpx.
					GetInfo(true).
					Print(printArgs)
			}

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
