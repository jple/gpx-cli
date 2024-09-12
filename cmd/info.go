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
		Short: "General info on the track",
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.SetVitesse(4.5)
			gpxSummary := gpx.GetInfo(viper.GetBool("detail"), true)
			gpxSummary.Print()
		},
	}

	initFlags(cmd, []FlagConfig{
		{
			Name: "detail", Shortname: "d", DefaultValue: false,
			Description: "Details info for each portion having name in trkpt",
			NoOptDefVal: StringPointer("true"),
		},
	})

	cmd.AddCommand(CreateInfoDetailCmd())

	return cmd
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
