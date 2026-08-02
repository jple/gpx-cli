package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/jple/gpx-cli/internal/geo"
	. "github.com/jple/gpx-cli/internal/gpx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getClosestTrkpts returns slice of pointers to Trkpt that are closest to p
// This function is used to add name to trkpt in-place without the need to
// specify Trk, Trkseg, Trkpt id
func getClosestTrkpts(gpx Gpx, p geo.Coord) []*Trkpt {
	var trkpts []*Trkpt
	var minDist float64
	// var ind struct{ i, j, k int }

	for i, _ := range gpx.Trks {
		for j, _ := range gpx.Trks[i].Trksegs {
			for k, trkpt := range gpx.Trks[i].Trksegs[j].Trkpts {
				if i == 0 && j == 0 && k == 0 {
					minDist = geo.Dist(p, trkpt.Coord)
				}

				d := geo.Dist(p, trkpt.Coord)

				if d == minDist {
					trkpts = append(trkpts, &gpx.Trks[i].Trksegs[j].Trkpts[k])
				} else if d < minDist {
					// Using index to prevent copy value to keep correct address
					trkpts = []*Trkpt{&gpx.Trks[i].Trksegs[j].Trkpts[k]}
					// ind = struct{ i, j, k int }{i, j, k}
					minDist = d
				}
			}
		}
	}

	return trkpts
}

func CreateAddNameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-name [name] [lat] [lon]",
		Short: "Add name to closest trkpts",
		Long:  `Add name to closest trkpts`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				fmt.Printf("Expected 3 args, got %v\n", len(args))
				for _, a := range args {
					fmt.Println(a)
				}
			}

			name := args[0]
			lat, _ := strconv.ParseFloat(args[1], 64)
			lon, _ := strconv.ParseFloat(args[2], 64)

			gpx := Gpx{}
			gpx.Parse(viper.GetString("filename"))

			p := geo.Coord{Lat: lat, Lon: lon}
			// TODO: add test if already exists
			gpx.AddWpt(Wpt{Name: &name, Coord: p})
			closest := getClosestTrkpts(gpx, p)

			for i, _ := range closest {
				// Confirmation if existing name
				if n := closest[i].Name; n != nil {
					scanner := bufio.NewScanner(os.Stdin)
					fmt.Printf("Overwrite existing name (%v) ?[y/n] ", *n)
					for scanner.Scan() {
						yn := scanner.Text()
						if yn == "y" {
							fmt.Println("Replacing name to", name)
							closest[i].SetName(name)
							break
						} else if scanner.Text() == "n" {
							break
						}
					}
				} else {
					closest[i].SetName(name)
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
