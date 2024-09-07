package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

func (gpx *Gpx) ParseFile(gpxFilename string) {
	data, _ := os.ReadFile(gpxFilename)
	if err := xml.Unmarshal(data, &gpx); err != nil {
		fmt.Println(err)
	}
}

func cli_dist(l1 float64, L1 float64, l2 float64, L2 float64) float64 {
	// Usage : go run . dist 45.23 15.18 45.42 15.26

	p1 := Pos{Lat: l1, Lon: L1, Ele: float64(0)}
	p2 := Pos{Lat: l2, Lon: L2, Ele: float64(0)}
	fmt.Printf("%.2f km\n", dist(p1, p2))
	return dist(p1, p2)
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
