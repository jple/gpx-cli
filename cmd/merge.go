package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateMergeCmd() *cobra.Command {
	var trkId IntValue = -1
	flagsConf := []FlagConfig{
		{
			Name: "trk_id", Shortname: "t", DefaultValue: &trkId,
			Description: "Trk id to merge with the next one.",
		},
	}
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge Trk",
		Long:  `Merge Trk at id with the next one`,
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			trkId := viper.GetInt("trk_id")
			if trkId == -1 {
				fmt.Println("Argument trk_id must be defined")
				return
			}

			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))

			if trkId == len(gpx.Trks)-1 {
				fmt.Printf("The chosen trkId (%v) is the last one\n", trkId)
				fmt.Println("Nothing to do")
				return
			}
			if trkId >= len(gpx.Trks)-1 {
				fmt.Printf("The chosen trkId (%v) is greater than the number of trk ()\n", trkId, len(gpx.Trks))
				fmt.Println("Nothing to do")
				return
			}

			fmt.Println("---------------")
			fmt.Println("Before merge")
			fmt.Println("---------------")
			gpx.Ls(true).Print(true)

			gpx = gpx.Merge(trkId, trkId+1)

			// TODO: to remove. For test only
			fmt.Println("---------------")
			fmt.Println("After merge")
			fmt.Println("---------------")
			gpx.Ls(true).Print(true)

			// Save
			fmt.Println("====================")
			gpx.Save(viper.GetString("output"))

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
