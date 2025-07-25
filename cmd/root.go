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

	var infile, outfile StringValue = "", "out.gpx"
	var inplace BoolValue = false
	flagsConf := []FlagConfig{
		{
			Name: "filename", Shortname: "f",
			DefaultValue:   &infile,
			Description:    "GPX filename to load",
			PersistentFlag: BoolPointer(true),
		},
		// TODO: these params are not always useful. Must move
		{
			Name: "output", Shortname: "o",
			DefaultValue:   &outfile,
			Description:    "Path to save GPX",
			PersistentFlag: BoolPointer(true),
		}, {
			Name: "inplace", Shortname: "x",
			DefaultValue:   &inplace,
			Description:    "Save inplace",
			PersistentFlag: BoolPointer(true),
			NoOptDefVal:    StringPointer("true"),
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
	rootCmd.AddCommand(CreateAddNameCmd())
	rootCmd.AddCommand(CreateFetchElevationCmd())
	rootCmd.AddCommand(CreateSplitCmd())
	rootCmd.AddCommand(CreateMergeCmd())
	rootCmd.AddCommand(CreateTuiCmd())
	rootCmd.AddCommand(CreateColorCmd())

	rootCmd.AddCommand(CreateTestCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
