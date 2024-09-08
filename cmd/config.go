package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FlagConfig struct {
	Name         string
	Shortname    string
	DefaultValue any
	Description  string

	PersistentFlag *bool
}

func BoolPointer(b bool) *bool {
	return &b
}

func initFlags(cmd *cobra.Command, flags []FlagConfig) {
	for _, f := range flags {

		// Create flags
		switch f.DefaultValue.(type) {
		case string:
			if f.PersistentFlag != nil {
				if *f.PersistentFlag {
					cmd.PersistentFlags().StringP(
						f.Name, f.Shortname, f.DefaultValue.(string), f.Description)
					break
				}
			}

			cmd.Flags().StringP(f.Name, f.Shortname, f.DefaultValue.(string), f.Description)
			break
		case float64:
			if f.PersistentFlag != nil {
				if *f.PersistentFlag {
					cmd.PersistentFlags().Float64P(
						f.Name, f.Shortname, f.DefaultValue.(float64), f.Description)
					break
				}
			}
			cmd.Flags().Float64P(f.Name, f.Shortname, f.DefaultValue.(float64), f.Description)
			break
		case bool:
			if f.PersistentFlag != nil {
				if *f.PersistentFlag {
					cmd.PersistentFlags().BoolP(
						f.Name, f.Shortname, f.DefaultValue.(bool), f.Description)
					break
				}
			}
			cmd.Flags().BoolP(f.Name, f.Shortname, f.DefaultValue.(bool), f.Description)
			break
		}

		// Bind flags to viper
		if f.PersistentFlag != nil {
			if *f.PersistentFlag {
				viper.BindPFlag(f.Name, cmd.PersistentFlags().Lookup(f.Name))
				continue
			}
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
