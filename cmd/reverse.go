package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateReverseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reverse",
		Short: "Reverse GPX elements",
		Long:  "Reverse trk and trkwpt in GPX file.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}

	flags := []FlagConfig{
		{
			Name: "trk_id", Shortname: "t", DefaultValue: "all",
			Description: "Trk id to reverse. (example: -t 2)",
		},
	}

	for _, f := range flags {
		cmd.Flags().StringP(f.Name, f.Shortname, f.DefaultValue.(string), f.Description)
		viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name))
	}

	return cmd
}
