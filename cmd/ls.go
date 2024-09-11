package cmd

import (
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
			// trkNames := ls(viper.GetString("filename"), viper.GetBool("all"))
			Gpx{Filepath: viper.GetString("filename")}.
				Ls(viper.GetBool("all")).
				Print(viper.GetBool("all"))
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
