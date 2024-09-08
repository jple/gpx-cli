package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "gpx-cli",
		Short: "GPX File utility",
	}
)

type FlagConfig struct {
	Name         string
	Shortname    string
	DefaultValue any
	ValueType    *any
	Description  string
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("gpx", "i", "", "GPX File input")

	rootCmd.AddCommand(CreateReverseCmd())
	rootCmd.AddCommand(CreateCalcEffortCmd())
	rootCmd.AddCommand(CreateDistCmd())
	rootCmd.AddCommand(CreateInfoCmd())
	rootCmd.AddCommand(CreateLsCmd())
}

func initFlags(cmd *cobra.Command, flags []FlagConfig) {
	for _, f := range flags {
		switch f.DefaultValue.(type) {
		case string:
			cmd.Flags().StringP(f.Name, f.Shortname, f.DefaultValue.(string), f.Description)
			break
		case float64:
			cmd.Flags().Float64P(f.Name, f.Shortname, f.DefaultValue.(float64), f.Description)
			break
		case bool:
			cmd.Flags().BoolP(f.Name, f.Shortname, f.DefaultValue.(bool), f.Description)
			break
		}
		viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name))
	}

}
func initConfig() {
	// Load config from file
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gpxrc")
	}

	// Load/Overwrite config from env var
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() error {
	return rootCmd.Execute()
}
