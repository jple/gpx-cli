package cmd

import (
	"fmt"
	"strconv"

	"github.com/jple/gpx-cli/internal/geo"

	"github.com/spf13/cobra"
)

func CreateCalcEffortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calc-effort [dist] [denivPos] [denivNeg] [[flatSpeed]]",
		Short: "Calculate duration based input (distance, denivPos, denivNegflatSpeed). flatSpeed is 6 km/h if not set",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			dist, _ := strconv.ParseFloat(args[0], 64)
			dPos, _ := strconv.ParseFloat(args[1], 64)
			dNeg, _ := strconv.ParseFloat(args[2], 64)
			var flatSpeed float64
			if len(args) > 3 {
				v, _ := strconv.ParseFloat(args[3], 64)
				flatSpeed = v
			} else {
				flatSpeed = 6.0
			}

			fmt.Println("Input:")
			fmt.Println("distance", dist, "km")
			fmt.Println("denivele positif: ", dPos, "m")
			fmt.Println("denivele negatif:", dNeg, "m")
			fmt.Println("vitesse plat", flatSpeed, "km/h")
			fmt.Println("")

			fmt.Println("Output:")
			distEffort := geo.CalcDistanceEffort(dist, dPos, dNeg)
			fmt.Println("Distance effort : ", distEffort, "km")
			_, h, m := geo.CalcDuration(distEffort, flatSpeed)
			fmt.Println("Duration : ", h, "h", m)
		},
	}

	return cmd
}
