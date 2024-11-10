package cmd

import (
	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateLsCmd() *cobra.Command {
	var all BoolValue = false
	flagsConf := []FlagConfig{
		{
			Name: "all", Shortname: "a", DefaultValue: &all,
			Description: "Include trkpt names",
			NoOptDefVal: StringPointer("true"),
		},
	}

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List all trk names",
		// TraverseChildren: true,

		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			Gpx{Filepath: viper.GetString("filename")}.
				Ls(viper.GetBool("all")).
				Print(viper.GetBool("all"))
		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
