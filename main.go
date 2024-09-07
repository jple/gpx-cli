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

	vitessePlat := 4.5
	if len(os.Args) > 1 {
		if os.Args[1] == "dist" {
			l1, _ := strconv.ParseFloat(os.Args[2], 8)
			L1, _ := strconv.ParseFloat(os.Args[3], 8)
			l2, _ := strconv.ParseFloat(os.Args[4], 8)
			L2, _ := strconv.ParseFloat(os.Args[5], 8)

			cli_dist(l1, L1, l2, L2)
			os.Exit(0)
		}

		if os.Args[1] == "calc_effort" {
			km, _ := strconv.ParseFloat(os.Args[2], 64)
			denivPos, _ := strconv.ParseFloat(os.Args[3], 64)
			denivNeg, _ := strconv.ParseFloat(os.Args[4], 64)

			cli_calc_effort(km, denivPos, denivNeg, vitessePlat)
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
					from           string
					to             string
					distance       float64
					denivPos       float64
					denivNeg       float64
					distanceEffort float64
					durationH      int8
					durationM      int8
				}

				var p_prev Pos
				var d float64 = 0
				var denivPos float64 = 0
				var denivNeg float64 = 0

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
					eleDiff := diffElevation(p_prev, p)
					if eleDiff > 1 {
						denivPos += eleDiff
					} else if eleDiff < -1 {
						denivNeg += eleDiff
					}

					p_prev = p

					if trkpt.Name != nil {
						x := struct {
							from           string
							to             string
							distance       float64
							denivPos       float64
							denivNeg       float64
							distanceEffort float64
							durationH      int8
							durationM      int8
						}{
							from:           pointName_prev,
							to:             *trkpt.Name,
							distance:       d,
							denivPos:       denivPos,
							denivNeg:       denivNeg,
							distanceEffort: calcDistanceEffort(d, denivPos, denivNeg),
							durationH:      0,
							durationM:      0,
						}
						duration := x.distanceEffort / vitessePlat
						x.durationH, x.durationM = floatToHourMin(duration)

						l = append(l, x)

						pointName_prev = *trkpt.Name
						d = 0
						denivPos = 0
						denivNeg = 0
					}
				}

				for _, p := range l {
					// fmt.Printf("from: %v \nto: %v \ndistance: %.1f", p.from, p.to, p.distance)
					fmt.Printf("%v \u2192 %v \ndistance: %.1f km\n", p.from, p.to, p.distance)
					fmt.Printf("deniv: +%.0f / -%.0f\n", p.denivPos, math.Abs(p.denivNeg))
					fmt.Printf("distance effort: %.1f km\n", p.distanceEffort)
					fmt.Printf("Durée: %vh%v", p.durationH, p.durationM)

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
				trk.setVitesse(vitessePlat)

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
			trk.setVitesse(vitessePlat)

			trk.calcAll()
			trk.Info()

			fmt.Println()
		}
	}
}
