package main

import (
	"fmt"
	"log"
	"os"
	"math"
	"strconv"
)

func main() {
	gpxFilename := "./src/my.gpx"
	_, err := os.Open(gpxFilename)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "dist" {
			fmt.Println(cli_dist())
			os.Exit(0)
		}

		if os.Args[1] == "calc_effort" {
			// WIP
			km, _ := strconv.ParseFloat(os.Args[2], 64)
			denivPos, _ := strconv.ParseFloat(os.Args[3], 64)
			denivNeg, _ := strconv.ParseFloat(os.Args[4], 64)
			vitessePlat := 4.0

			km_effort := km + (denivPos / 100.0) + math.Abs(denivNeg / 300.0)
			duration := km_effort / vitessePlat
			h, m := floatToHourMin(duration) 
			fmt.Println(h, "h", m, "m")
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
			var ascii_format bool
			if len(os.Args) > 2 {
				ascii_format, _ = strconv.ParseBool(os.Args[2])
			} else {
				ascii_format = true
			}

			var gpx Gpx
			gpx.ParseFile(gpxFilename)

			for _, trk := range gpx.Trk {
				trk.setVitesse(4.0)

				trk.calcAll()
				trk.Info(ascii_format)

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
