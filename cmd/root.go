package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "gpx-cli",
		Short: "GPX File utility",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	var s StringValue = ""
	flagsConf := []FlagConfig{
		{
			Name: "filename", Shortname: "f",
			// DefaultValue: "",
			DefaultValue:   &s,
			Description:    "GPX filename to load",
			PersistentFlag: BoolPointer(true),
		},
	}
	initFlags(rootCmd, flagsConf)
	bindFlags(rootCmd, flagsConf)

	rootCmd.AddCommand(CreateReverseCmd())
	rootCmd.AddCommand(CreateCalcEffortCmd())
	rootCmd.AddCommand(CreateDistCmd())
	rootCmd.AddCommand(CreateInfoCmd())
	rootCmd.AddCommand(CreateLsCmd())
	rootCmd.AddCommand(CreatePlotCmd())
	rootCmd.AddCommand(CreateTermPlotCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
