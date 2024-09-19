package cmd

import (
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
		Short: "General info on the track",
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.SetVitesse(4.5)
			gpx.ParseFile(gpx.Filepath)

			var summary GpxSummary
			var printArgs PrintArgs = PrintArgs{AsciiFormat: true}
			if viper.Get("trk-id") != -1 {
				summary = append(summary, gpx.Trk[viper.GetInt("trk-id")].GetInfo(gpx.Extensions.Vitesse, true))
				printArgs.PrintFrom = true
			} else {
				summary = gpx.GetInfo(true)
				printArgs.PrintFrom = false
			}

			for _, s := range summary {
				s.Print(printArgs)
			}
		},
	}

	initFlags(cmd, []FlagConfig{
		{
			Name: "detail", Shortname: "d", DefaultValue: false,
			Description: "Details info for each portion having name in trkpt",
			NoOptDefVal: StringPointer("true"),
		},
		{
			Name: "trk-id", Shortname: "i", DefaultValue: int8(-1),
			Description: "Details about i-th trk",
		},
	})

	return cmd
}
