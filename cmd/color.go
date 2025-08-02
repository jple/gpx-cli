package cmd

import (
	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateColorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-color",
		Short: "Add color to each trk",

		Run: func(cmd *cobra.Command, args []string) {
			filename := viper.GetString("filename")

			gpx := Gpx{}
			gpx.ParseFile(filename).
				AddColor().Save(filename)
			// AddColor().Save("out.gpx")
		},
	}

	return cmd
}
