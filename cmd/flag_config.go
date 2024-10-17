package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagConfig struct {
	Name         string
	Shortname    string
	DefaultValue pflag.Value
	Description  string

	PersistentFlag *bool
	NoOptDefVal    *string
}

func BoolPointer(b bool) *bool {
	return &b
}

func createFlag(cmd *cobra.Command, f FlagConfig) {
	// Create PersistentFlag
	if f.PersistentFlag != nil && *f.PersistentFlag {
		cmd.PersistentFlags().VarP(
			f.DefaultValue, f.Name, f.Shortname, f.Description)

		// Set flag without value (eg. --verbose)
		if f.NoOptDefVal != nil {
			cmd.PersistentFlags().Lookup(f.Name).NoOptDefVal = *f.NoOptDefVal
		}

		// Bind Flag
		viper.BindPFlag(f.Name, cmd.PersistentFlags().Lookup(f.Name))
	} else
	// Create "local" flag
	{
		cmd.Flags().VarP(
			f.DefaultValue, f.Name, f.Shortname, f.Description)

		// Set flag without value (eg. --verbose)
		if f.NoOptDefVal != nil {
			cmd.Flags().Lookup(f.Name).NoOptDefVal = *f.NoOptDefVal
		}

		// Bind Flag
		viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name))
	}
}

func initFlags(cmd *cobra.Command, flags []FlagConfig) {
	for _, f := range flags {
		createFlag(cmd, f)
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
