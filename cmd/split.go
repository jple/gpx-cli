package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateSplitTrkCmd() *cobra.Command {
	var name, point StringValue
	flagsConf := []FlagConfig{
		{
			Name: "name", Shortname: "n", DefaultValue: &name,
			Description: "Trkpt name to split Trk",
		}, {
			Name: "point", Shortname: "p", DefaultValue: &point,
			Description: "(closest) point to split Trk. Format 'lat,lon' (eg. 45.123,3.564)",
		},
	}
	cmd := &cobra.Command{
		Use:   "split-trk",
		Short: "Split Trk",
		Long:  `Split Trk at specific name or closest trkpt`,
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))

			fmt.Println("Before split")
			gpx.Ls(true).Print(true)

			// Split on name input
			if name := viper.GetString("name"); name != "" {
				gpx = gpx.SplitAtName(name)
			}

			// Split on closest point (lat, lon) input
			if point := viper.GetString("point"); point != "" {
				// TODO: rework gpx.GetClosestTrkpts to output i, j, k
				fmt.Println("ok")
				return
			}

			// TODO: to remove. For test only
			fmt.Println("====================")
			fmt.Println("After split")
			gpx.Ls(true).Print(true)

			// Save
			fmt.Println("====================")
			gpx.Save(viper.GetString("output"))
		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
