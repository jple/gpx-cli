package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreateInfoDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Detailed info on the selected track",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}

	initFlags(cmd, []FlagConfig{
		{
			Name: "trk_id", Shortname: "t", DefaultValue: int8(0),
			Description: "Trk id to reverse. (example: -t 2)",
		},
	})

	return cmd
}

func CreateInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "General info on the track",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}

	cmd.AddCommand(CreateInfoDetailCmd())

	return cmd
}
