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
	var flags *pflag.FlagSet

	if f.PersistentFlag != nil && *f.PersistentFlag {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	// Create Flag or PersistentFlag
	flags.VarP(
		f.DefaultValue, f.Name, f.Shortname, f.Description)

	// Set flag without value (eg. --verbose)
	if f.NoOptDefVal != nil {
		flags.Lookup(f.Name).NoOptDefVal = *f.NoOptDefVal
	}
}

func bindFlag(cmd *cobra.Command, f FlagConfig) {
	if f.PersistentFlag != nil && *f.PersistentFlag {
		viper.BindPFlag(f.Name, cmd.PersistentFlags().Lookup(f.Name))
	} else {
		viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name))
	}
}
func bindFlags(cmd *cobra.Command, flagsConf []FlagConfig) {
	for _, f := range flagsConf {
		bindFlag(cmd, f)
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
