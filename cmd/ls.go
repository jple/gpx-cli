package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List all trk names",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(viper.Get("all"))
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
