package cmd

import (
	"fmt"
	"slices"

	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List all trk names",

		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {

			ls(viper.Get("filename").(string), viper.Get("all").(bool))
			// for i, trackName := range ls(viper.Get("filename").(string)) {
			// 	fmt.Printf("%v: \u001b[1;32m%v\u001b[22;0m\n", i, trackName)
			// }
		},
	}

	initFlags(cmd, []FlagConfig{
		{
			Name: "all", Shortname: "a", DefaultValue: false,
			Description: "Include trkpt names",
		},
	})

	return cmd
}

func ls(gpxFilename string, all bool) []string {
	var gpx Gpx
	gpx.ParseFile(gpxFilename)

	var out []string
	for i, trk := range gpx.Trk {
		s := fmt.Sprintf("%v: \u001b[1;32m%v\u001b[22;0m\n", i, trk.Name)
		out = append(out, trk.Name)
		fmt.Printf(s)

		if all {
			trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
			for _, trkpt := range trkpts {
				if trkpt.Name != nil {
					s := *trkpt.Name
					out = append(out, s)
					fmt.Println(s)
				}
			}
			fmt.Println()
		}
	}

	return out
}
