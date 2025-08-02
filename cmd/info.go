package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateInfoCmd() *cobra.Command {
	var trkId IntValue = -1
	var detail BoolValue = false
	var speed FloatValue = 4.5

	flagsConf := []FlagConfig{
		{
			Name: "trk-id", Shortname: "t", DefaultValue: &trkId,
			Description: "Details about i-th trk. Value -1 will display all trk summary",
		},
		{
			Name: "detail", Shortname: "d", DefaultValue: &detail,
			Description: "Details on trk",
			NoOptDefVal: StringPointer("true"),
		},
		{
			Name: "speed", Shortname: "s", DefaultValue: &speed,
			Description: "Hiking speed on flat (km/h)",
		},
	}

	cmd := &cobra.Command{
		Use:   "info",
		Short: "General info on the tracks",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))

			var printArgs PrintArgs = PrintArgs{AsciiFormat: true}

			trkid := viper.GetInt("trk-id")
			if trkid == -1 {
				printArgs.PrintFrom = bool(detail)
				gpx.GetInfo(viper.GetFloat64("speed")).ToString(printArgs)
			} else {
				printArgs.PrintFrom = true

				fmt.Printf("[%v] ", trkid)
				gpx.
					Trks[trkid].
					GetInfo(trkid, viper.GetFloat64("speed")).
					ToString(printArgs)
			}

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
