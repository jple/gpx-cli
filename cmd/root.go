package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gpx-cli",
		Short: "GPX File utility",
	}
)

func init() {
	rootCmd.AddCommand(reverse)
	rootCmd.AddCommand(calc_effort)
}
