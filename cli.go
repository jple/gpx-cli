package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

func cli_dist() float64 {
	// Usage : go run . dist 45.23 15.18 45.42 15.26
	l1, _ := strconv.ParseFloat(os.Args[2], 8)
	L1, _ := strconv.ParseFloat(os.Args[3], 8)
	l2, _ := strconv.ParseFloat(os.Args[4], 8)
	L2, _ := strconv.ParseFloat(os.Args[5], 8)

	fmt.Printf("l1: %v \n L1: %v \n l2: %v \n L2: %v \n", l1, L1, l2, L2)
	p1 := Pos{Lat: l1, Lon: L1, Ele: float64(0)}
	p2 := Pos{Lat: l2, Lon: L2, Ele: float64(0)}
	return dist(p1, p2)
}

func (gpx *Gpx) ParseFile(gpxFilename string) {
	data, _ := os.ReadFile(gpxFilename)
	if err := xml.Unmarshal(data, &gpx); err != nil {
		fmt.Println(err)
	}
}

func ls(gpxFilename string) []string {
	var gpx Gpx
	gpx.ParseFile(gpxFilename)

	var out []string
	for _, trk := range gpx.Trk {
		out = append(out, trk.Name)
	}

	return out
}
