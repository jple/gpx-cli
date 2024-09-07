package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateCalcEffortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calc_effort",
		Short: "Calculate duration based input (distance, denivPos, denivNeg)",
	}

	cmd.Flags().FloatP("distance", "d", 0.0, "Distance in km")
	cmd.Flags().FloatP("deniv_pos", "p", 0.0, "Positive elevation")
	cmd.Flags().FloatP("deniv_neg", "n", 0.0, "Positive elevation")
	viper.BindPFlag("distance", cmd.Flags().Lookup("distance"))
	viper.BindPFlag("deniv_pos", cmd.Flags().Lookup("deniv_pos"))
	viper.BindPFlag("deniv_neg", cmd.Flags().Lookup("deniv_neg"))

	return cmd
}
