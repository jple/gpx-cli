package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	gpxFilename := "./my.gpx"
	_, err := os.Open(gpxFilename)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "dist" {
			fmt.Println(cli_dist())
			os.Exit(0)
		}

		if os.Args[1] == "ls" {
			for i, trackName := range ls(gpxFilename) {
				fmt.Println(i, ":", trackName)
			}
		}

		if os.Args[1] == "test" {
			var gpx Gpx
			gpx.ParseFile(gpxFilename)

			for _, trk := range gpx.Trk {
				for _, trkpt := range trk.Trkseg.Trkpt {
					if trkpt.Name != nil {
						fmt.Println(*trkpt.Name)
					}
				}
			}

		}
		if os.Args[1] == "info" {

			var gpx Gpx
			gpx.ParseFile(gpxFilename)

			for _, trk := range gpx.Trk {
				trk.setVitesse(4.0)

				trk.calcAll()
				trk.Info()

				fmt.Println()
			}

		}

		// for _, trkseg := range trk.Trkseg.Trkpt {
		// 	fmt.Println(trkseg)
		// }

		// fmt.Println(gpx.Trk[0].Trkseg.Trkpt[0])
		// fmt.Println(gpx.Trk[0].Trkseg.Len())
	}

	if len(os.Args) == 1 {

		var gpx Gpx
		gpx.ParseFile(gpxFilename)

		// fmt.Println(reflect.TypeOf(gpx.Trk))
		// for _, trk := range gpx.Trk {
		// 	fmt.Println(trk.Name)
		// }
		// trk := gpx.Trk[0]
		for _, trk := range gpx.Trk {
			trk.setVitesse(4.0)

			trk.calcAll()
			trk.Info()

			fmt.Println()
		}
	}
}
