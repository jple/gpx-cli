package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "gpx-cli",
		Short: "GPX File utility",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("test: ", viper.Get("filename"))
		// },
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	initFlags(rootCmd, []FlagConfig{
		{
			Name: "filename", Shortname: "f", DefaultValue: "",
			Description:    "GPX filename to load",
			PersistentFlag: BoolPointer(true),
		},
	})

	rootCmd.AddCommand(CreateReverseCmd())
	rootCmd.AddCommand(CreateCalcEffortCmd())
	rootCmd.AddCommand(CreateDistCmd())
	rootCmd.AddCommand(CreateInfoCmd())
	rootCmd.AddCommand(CreateLsCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
