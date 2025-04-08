package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateSplitCmd() *cobra.Command {
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
		Use:   "split",
		Short: "Split Trk",
		Long:  `Split Trk at specific name or closest trkpt`,
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			fmt.Println(gpx)
			// gpx.ParseFile(viper.GetString("filename"))

			if name := viper.GetString("name"); name != "" {
				found := false
			out:
				for i, trk := range gpx.Trk {
					for j, trkseg := range trk.Trkseg {
						for k, trkpt := range trkseg.Trkpt {
							if *trkpt.Name == name {
								found = true

								// TODO: create this
								fmt.Println(i, j, k)
								// gpx.Split(i, j, k)
								break out
							}
						}
					}
				}

				if !found {
					fmt.Println("Name (", name, ") not found in gpx")
				}
				fmt.Println("ok")
				return
			}
			if point := viper.GetString("point"); point != "" {
				// TODO: rework gpx.GetClosestTrkpts to output i, j, k
				fmt.Println("ok")
				return
			}
		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
