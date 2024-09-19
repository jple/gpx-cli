package cmd

import (
	"fmt"
	"strconv"

	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
)

func CreateDistCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dist [lat1] [lon1] [lat2] [lon2]",
		Short: "Calculate distance (km) based on GPS coordinates",
		Long: `Calculate distance (km) based on GPS coordinates
Expects 4 float positional arguments (lat1, lon1, lat2, lon2)
representing 2 GPS coordinates`,
		Run: func(cmd *cobra.Command, args []string) {
			lat1, _ := strconv.ParseFloat(args[0], 64)
			lon1, _ := strconv.ParseFloat(args[1], 64)
			lat2, _ := strconv.ParseFloat(args[2], 64)
			lon2, _ := strconv.ParseFloat(args[3], 64)

			d := Dist(
				Pos{Lat: lat1, Lon: lon1},
				Pos{Lat: lat2, Lon: lon2},
			)

			fmt.Printf("Distance: %.1f km\n", d)

		},
	}

	return cmd
}
