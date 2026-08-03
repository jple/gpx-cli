package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jple/gpx-cli/internal/gpx"
)

func CreateInfoCmd() *cobra.Command {
	var detail BoolValue = false
	var speed FloatValue = 4.5
	var trkId IntValue = -1
	var from, to StringValue
	var kind StringValue = "string"

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
			Description: "Details about i-th trk. Value -1 will display all trk stat",
		}, {
			Name: "from", Shortname: "a", DefaultValue: &from,
			Description: "Name from",
		}, {
			Name: "to", Shortname: "b", DefaultValue: &to,
			Description: "Name to",
		}, {
			Name: "kind", Shortname: "k", DefaultValue: &kind,
			Description: "Export type: string, csv, html",
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
			g := gpx.Gpx{}
			g.Parse(viper.GetString("filename"))

			trkid := viper.GetInt("trk-id")
			from := viper.GetString("from")
			to := viper.GetString("to")

			switch {
			// Defaut case : User wants to display stat of all trks
			default:
				gpx.ExportStatsTrks(os.Stdout,
					g.Trks,
					viper.GetFloat64("speed"),
					detail.Value(),
					viper.GetString("kind"),
					nil)

			// Case 1 : User wants to display stat between two named points
			case from != "" && to != "":
				trkpts, err := gpx.GpxTrkptsBetweenNames(
					g,
					viper.GetString("from"),
					viper.GetString("to"),
				)
				if err != nil {
					panic(err)
				}

				gpx.ExportStatsTrks(os.Stdout,
					[]gpx.Trk{
						gpx.Trk{
							Trksegs: []gpx.Trkseg{
								gpx.Trkseg{trkpts}}}},
					viper.GetFloat64("speed"),
					true, // always show detail here
					viper.GetString("kind"),
					nil)

			// Case 2 : User wants to display stat of a selected trk
			case trkid >= 0:
				gpx.ExportStatsTrks(os.Stdout,
					[]gpx.Trk{g.Trks[trkid]},
					viper.GetFloat64("speed"),
					true, // always show detail here
					viper.GetString("kind"),
					nil)

			}

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
