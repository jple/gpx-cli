package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateCalcEffortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calc_effort",
		Short: "Calculate duration based input (distance, denivPos, denivNeg)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}

	flags := []FlagConfig{
		{
			Name: "distance", Shortname: "d", DefaultValue: 0.0,
			Description: "Distance in km",
		},
		{
			Name: "deniv_pos", Shortname: "p", DefaultValue: 0.0,
			Description: "Positive elevation",
		},
		{
			Name: "deniv_neg", Shortname: "n", DefaultValue: 0.0,
			Description: "Negative elevation",
		},
	}

	for _, f := range flags {
		cmd.Flags().Float64P(f.Name, f.Shortname, f.DefaultValue.(float64), f.Description)
		viper.BindPFlag(f.Name, cmd.Flags().Lookup(f.Name))
	}

	return cmd
}
