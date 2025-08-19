package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateInfoCmd() *cobra.Command {
	var detail BoolValue = false
	var speed FloatValue = 4.5
	var trkId IntValue = -1
	var from, to StringValue

	flagsConf := []FlagConfig{
		{
			Name: "detail", Shortname: "d", DefaultValue: &detail,
			Description: "Details on trk",
			NoOptDefVal: StringPointer("true"),
		}, {
			Name: "speed", Shortname: "s", DefaultValue: &speed,
			Description: "Hiking speed on flat (km/h)",
		}, {
			Name: "trk-id", Shortname: "t", DefaultValue: &trkId,
			Description: "Details about i-th trk. Value -1 will display all trk summary",
		}, {
			Name: "from", Shortname: "a", DefaultValue: &from,
			Description: "Name from",
		}, {
			Name: "to", Shortname: "b", DefaultValue: &to,
			Description: "Name to",
		},
	}

	cmd := &cobra.Command{
		Use:   "info",
		Short: "General info on the tracks",
		Long: `Example: 
	Get info on each trk: 
	gpx-cli -f file.gpx info

	same, with intermediate trkname: 
	gpx-cli -f file.gpx info -d

	Get info on a specific trk id:
	gpx-cli -f file.gpx info -t 1

	Get info between two named trkpt:
	gpx-cli -f file.gpx info --from trkptname1 --to trkptname2
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))

			var printArgs PrintArgs = PrintArgs{AsciiFormat: true}

			trkid := viper.GetInt("trk-id")
			from := viper.GetString("from")
			to := viper.GetString("to")
			if from != "" && to != "" {
				// TODO: pouvoir afficher les name intermédiaires entre from et to 
				// avec un argument detail
				var printArgs PrintArgs = PrintArgs{AsciiFormat: true, PrintFrom: true}
				// NOTE: print car GpxSummary ne print pas
				// TODO: retirer le print de tous les ToString, et faire un print à l'utilisation
				fmt.Println(
					gpx.GetInfoBetweenName(
						viper.GetString("from"),
						viper.GetString("to"),
						viper.GetFloat64("speed"),
					).ToString(printArgs))
				return
			}
			if trkid > -1 {
				printArgs.PrintFrom = true

				fmt.Printf("[%v] ", trkid)
				fmt.Println(
					gpx.
						Trks[trkid].
						GetInfo(trkid, viper.GetFloat64("speed")).
						ToString(printArgs))
				return
			}
			printArgs.PrintFrom = bool(detail)
			gpx.GetInfo(viper.GetFloat64("speed")).ToString(printArgs)

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
