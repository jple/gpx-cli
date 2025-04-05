package cmd

import (
	"fmt"
	"strconv"

	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateClosestTrkptsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "closest_trkpts [lat] [lon]",
		Short: "Get closest trkpt from specified lat/lon point",
		Long:  `Get closest trkpt from specified lat/lon point`,
		Run: func(cmd *cobra.Command, args []string) {
			lat, _ := strconv.ParseFloat(args[0], 64)
			lon, _ := strconv.ParseFloat(args[1], 64)

			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.ParseFile(gpx.Filepath)

			closest := gpx.GetClosestTrkpts(Pos{Lat: lat, Lon: lon})
			for _, trkpt := range closest {
				fmt.Printf("%+v\n", trkpt)
			}
			// fmt.Printf("gpx: %+v\n", gpx.Trk[8].Trkseg[0].Trkpt[654])
			// closest[0].Lat = 44.156
			// fmt.Printf("%+v\n", closest[0])
			// fmt.Printf("gpx: %+v\n", gpx.Trk[8].Trkseg[0].Trkpt[654])
		},
	}

	return cmd
}
