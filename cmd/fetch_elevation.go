package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/jple/gpx-cli/core"
	"github.com/jple/gpx-cli/ign"
)

func CreateFetchElevationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch-elevation",
		Short: "Fetch elevation from IGN API, and fill GPX",
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))

			var pts ign.Points
			// for _, trk := range gpx.Trks {
			for i, trk := range gpx.Trks {
				// for _, trkseg := range trk.Trksegs {
				for j, trkseg := range trk.Trksegs {
					// TODO: asynchronous fetch to prevent 2 loops
					c := 0
					chunkSize := 2000
					for trkchunk := range slices.Chunk(trkseg.Trkpts, chunkSize) {
						pts = nil
						for _, trkpt := range trkchunk {
							pts = append(pts, ign.Point{Lat: trkpt.Lat, Lon: trkpt.Lon})
						}

						queryAgain := true
						retry := 0
						var elevations []float64
						for queryAgain {
							if retry > 5 {
								panic("Max retry exceeded to query API")
							}
							elevations = ign.PostElevations(pts)
							if len(elevations) > 0 {
								queryAgain = false
							} else {
								time.Sleep(250 * time.Millisecond) // Prevent server overwhelm
								retry++
								fmt.Println("retry", retry)
							}
						}

						for k, _ := range trkchunk {
							gpx.Trks[i].Trksegs[j].Trkpts[c*chunkSize+k].Ele = elevations[k]
						}
						c++
					}
				}
			}

			if viper.GetBool("inplace") {
				gpx.Save(viper.GetString("filename"))
			} else {
				gpx.Save(viper.GetString("output"))
			}

		},
	}

	return cmd
}
