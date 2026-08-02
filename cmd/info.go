package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/internal/gpx"
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
			gpx.Parse(viper.GetString("filename"))

			trkid := viper.GetInt("trk-id")
			from := viper.GetString("from")
			to := viper.GetString("to")
			if from != "" && to != "" {
				// TODO: pouvoir afficher les name intermédiaires entre from et to avec un argument detail
				fmt.Printf(
					gpx.SummarizeBetweenTrkptsNames(
						viper.GetString("from"),
						viper.GetString("to"),
						viper.GetFloat64("speed"),
					).ToString())
				return
			}
			if trkid > -1 {
				fmt.Printf("[%v] ", trkid)
				fmt.Println(
					gpx.
						Trks[trkid].
						Summarize(trkid, viper.GetFloat64("speed")).
						ToString(detail.Value()))
				return
			}
			fmt.Println(gpx.Summarize(viper.GetFloat64("speed")).ToString(detail.Value()))

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
