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

			// ascii_format, _ := strconv.ParseBool(args[0])
			trkNames := ls(viper.Get("filename").(string), viper.GetBool("all"))
			trkNames.Print(viper.GetBool("all"))
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

type TrkName struct {
	TrkName    string
	TrkptNames []string
}
type TrkNames []TrkName

func ls(gpxFilename string, all bool) TrkNames {
	var gpx Gpx
	gpx.ParseFile(gpxFilename)

	var out TrkNames
	for i, trk := range gpx.Trk {
		out = append(out, TrkName{TrkName: trk.Name})

		if all {
			trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
			for _, trkpt := range trkpts {
				if trkpt.Name != nil {
					out[i].TrkptNames = append(out[i].TrkptNames, *trkpt.Name)
				}
			}
		}
	}

	return out
}

func (trkNames TrkNames) Print(all bool, ascii_format ...bool) {

	for i, trkName := range trkNames {
		if len(ascii_format) > 0 && !ascii_format[0] {
			fmt.Printf("[%v] %v\n", i, trkName.TrkName)
		} else {
			fmt.Printf("[%v] \u001b[1;32m%v\u001b[22;0m\n", i, trkName.TrkName)
		}
		if all {
			for _, trkptName := range trkName.TrkptNames {
				fmt.Println(trkptName)
			}
			fmt.Println()
		}
	}
}
