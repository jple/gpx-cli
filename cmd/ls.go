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

	flags := []FlagConfig{
		{
			Name: "all", Shortname: "a", DefaultValue: false,
			Description: "Include trkpt names",
		},
	}

	for _, f := range flags {
		cmd.Flags().BoolP(f.Name, f.Shortname, f.DefaultValue.(bool), f.Description)
		cmd.Flags().Lookup(f.Name).NoOptDefVal = "true"
		viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name))
	}

	return cmd
}