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
	cmd := &cobra.Command{
		Use:   "info",
		Short: "General info on the tracks",
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.SetVitesse(4.5)
			gpx.ParseFile(gpx.Filepath)

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

	var trkId IntValue = -1
	initFlags(cmd, []FlagConfig{
		{
			Name: "trk-id", Shortname: "i", DefaultValue: &trkId,
			Description: "Details about i-th trk. Value -1 will display all trk summary",
		},
	})

	return cmd
}
