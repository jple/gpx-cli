package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
)

func CreateCalcInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dist",
		Short: "Calculate info (km, elev, km_e, duration) between two named trkpt",
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))
			// gpx.SetVitesse(4.5)

			// TODO: WIP

		},
	}

	return cmd
}
