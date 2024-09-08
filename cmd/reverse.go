package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

	initFlags(cmd, []FlagConfig{
		{
			Name: "trk_id", Shortname: "t", DefaultValue: "all",
			Description: "Trk id to reverse. (example: -t 2)",
		},
	})

	return cmd
}
