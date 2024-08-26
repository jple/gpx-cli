package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
)

func main() {
	gpxFilename := "./src/my.gpx"
	// gpxFilename := "./out/toto.xml"
	// gpxFilename := "./src/toto_inverse.gpx"
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
			km, _ := strconv.ParseFloat(os.Args[2], 64)
			denivPos, _ := strconv.ParseFloat(os.Args[3], 64)
			denivNeg, _ := strconv.ParseFloat(os.Args[4], 64)
			vitessePlat := 4.0

			km_effort := km + (denivPos / 100.0) + math.Abs(denivNeg/300.0)
			duration := km_effort / vitessePlat
			h, m := floatToHourMin(duration)
			fmt.Println(h, "h", m, "m")
			os.Exit(0)
		}

		if os.Args[1] == "ls" {
			if len(os.Args) > 2 {
				if os.Args[2] == "-a" {
					var gpx Gpx
					gpx.ParseFile(gpxFilename)

					for _, trk := range gpx.Trk {
						fmt.Printf("\u001b[1;32m%v\u001b[22;0m\n", trk.Name)
						for _, trkpt := range trk.Trkseg.Trkpt {
							if trkpt.Name != nil {
								fmt.Println(*trkpt.Name)
							}
						}
						fmt.Println()
					}
				}
			} else {
				for i, trackName := range ls(gpxFilename) {
					fmt.Println(i, ":", trackName)
				}
			}
		}

		// Revert trk and save to out/toto.xml
		// Usage : reverse i  	# i-ieth trk, output trk
		// Usage : reverse all  # all tracks, output gpx
		if os.Args[1] == "reverse" {
			var gpx Gpx
			gpx.ParseFile(gpxFilename)

			if len(os.Args) > 2 {
				if os.Args[2] == "all" {
					slices.Reverse(gpx.Trk)
					for _, trk := range gpx.Trk {
						slices.Reverse(trk.Trkseg.Trkpt)
					}

					xmlFile, _ := os.Create("out/toto.xml")
					encoder := xml.NewEncoder(xmlFile)
					encoder.Indent("", "\t")
					if err := encoder.Encode(gpx); err != nil {
						fmt.Printf("error: %v\n", err)
					}

					os.Exit(0)
				}

				if i, err := strconv.ParseInt(os.Args[2], 0, 0); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					slices.Reverse(gpx.Trk[i].Trkseg.Trkpt)

					xmlFile, _ := os.Create("out/toto.xml")
					encoder := xml.NewEncoder(xmlFile)
					encoder.Indent("", "\t")
					if err := encoder.Encode(&gpx.Trk[i]); err != nil {
						fmt.Printf("error: %v\n", err)
					}
				}
			}

		}

		if os.Args[1] == "info-detail" {
			if len(os.Args) < 2 {
				fmt.Printf("Error, int missing parameter to select trk")
				os.Exit(1)
			}

			var gpx Gpx
			gpx.ParseFile(gpxFilename)

			if i, err := strconv.ParseInt(os.Args[2], 0, 0); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				trk := gpx.Trk[i]
				fmt.Printf("\u001b[1;32m%v\u001b[22;0m\n", trk.Name)
				pointName_prev := "start"

				var l []struct {
					from     string
					to       string
					distance float64
				}

				var p_prev Pos
				var d float64 = 0
				for i, trkpt := range trk.Trkseg.Trkpt {
					p := Pos{
						Lat: trkpt.Lat,
						Lon: trkpt.Lon,
						Ele: trkpt.Ele,
					}
					if i == 0 {
						p_prev = p
					}

					d += dist(p_prev, p)
					p_prev = p

					if trkpt.Name != nil {
						x := struct {
							from     string
							to       string
							distance float64
						}{
							from:     pointName_prev,
							to:       *trkpt.Name,
							distance: d,
						}
						l = append(l, x)

						pointName_prev = *trkpt.Name
						d = 0
					}
				}

				for _, p := range l {
					// fmt.Printf("from: %v \nto: %v \ndistance: %.1f", p.from, p.to, p.distance)
					fmt.Printf("%v \u2192 %v \ndistance: %.1f", p.from, p.to, p.distance)
					fmt.Println()
					fmt.Println()
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

			for i, trk := range gpx.Trk {
				trk.setVitesse(4.0)

				trk.calcAll()
				fmt.Printf("%v- ", i)
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
