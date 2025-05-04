package cmd

import (
	"fmt"
	"strconv"

	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
)

func CreateCalcEffortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calc-effort [dist] [denivPos] [denivNeg] [[vitessePlat]]",
		Short: "Calculate duration based input (distance, denivPos, denivNegvitessePlat). vitessePlat is 6 km/h if not set",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			dist, _ := strconv.ParseFloat(args[0], 64)
			dPos, _ := strconv.ParseFloat(args[1], 64)
			dNeg, _ := strconv.ParseFloat(args[2], 64)
			var vitessePlat float64
			if len(args) > 3 {
				v, _ := strconv.ParseFloat(args[3], 64)
				vitessePlat = v
			} else {
				vitessePlat = 6.0
			}

			fmt.Println("Input:")
			fmt.Println("distance", dist, "km")
			fmt.Println("denivele positif: ", dPos, "m")
			fmt.Println("denivele negatif:", dNeg, "m")
			fmt.Println("vitesse plat", vitessePlat, "km/h")
			fmt.Println("")

			fmt.Println("Output:")
			distEffort := CalcDistanceEffort(dist, dPos, dNeg)
			fmt.Println("Distance effort : ", distEffort, "km")
			_, h, m := CalcDuration(distEffort, vitessePlat)
			fmt.Println("Duration : ", h, "h", m)
		},
	}

	return cmd
}
