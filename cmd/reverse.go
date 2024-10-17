package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateReverseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reverse",
		Short: "Reverse GPX elements",
		Long:  "Reverse trk and trkwpt in GPX file.",
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.SetVitesse(4.5)
			gpx.ParseFile(gpx.Filepath)

			trkId := viper.GetInt("trk_id")
			if trkId == -1 {
				fmt.Println("Reverse all gpx content")
				gpx.Reverse()
			} else {
				fmt.Println("Reverse trk n°", trkId)
				gpx.Trk[trkId].Reverse()
			}

			fmt.Println("Save to", viper.GetString("output"))
			gpx.Save(viper.GetString("output"))
		},
	}

	var output StringValue = "out.gpx"
	var trkId IntValue = -1
	initFlags(cmd, []FlagConfig{
		{
			Name: "output", Shortname: "o", DefaultValue: &output,
			Description:    "Output filename (default: out.gpx)",
			PersistentFlag: BoolPointer(true),
		},
		{
			Name: "trk_id", Shortname: "t", DefaultValue: &trkId,
			Description: "Trk id to reverse. If -1, applies to all trk. (example: -t 2)",
		},
	})

	return cmd
}
