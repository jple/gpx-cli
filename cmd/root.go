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

type FlagConfig struct {
	Name         string
	Shortname    string
	DefaultValue any
	Description  string
}

func init() {
	rootCmd.AddCommand(CreateReverseCmd())
	rootCmd.AddCommand(CreateCalcEffortCmd())
	rootCmd.AddCommand(CreateDistCmd())
	rootCmd.AddCommand(CreateInfoCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
