package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "General info on the track",
		Run: func(cmd *cobra.Command, args []string) {
			Info(
				viper.GetString("filename"), 4.5, false, true,
			)
		},
	}

	cmd.AddCommand(CreateInfoDetailCmd())

	return cmd
}

func Info(gpxFilename string, vitessePlat float64, detail bool, ascii_format bool) {
	var gpx Gpx
	gpx.ParseFile(gpxFilename)

	for i, trk := range gpx.Trk {
		trk.SetVitesse(vitessePlat)

		summary := trk.CalcAll()
		fmt.Printf("[%v] ", i)
		summary.Print()
		// trk.PrintInfo(ascii_format)

		fmt.Println()
	}

}

func CreateInfoDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Detailed info on the selected track",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}

	initFlags(cmd, []FlagConfig{
		{
			Name: "trk_id", Shortname: "t", DefaultValue: int8(0),
			Description: "Trk id to reverse. (example: -t 2)",
		},
	})

	return cmd
}
