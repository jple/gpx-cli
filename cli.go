package main

import (
	"fmt"

	. "github.com/jple/gpx-cli/core"
)

func cli_dist(l1 float64, L1 float64, l2 float64, L2 float64) float64 {
	// Usage : go run . dist 45.23 15.18 45.42 15.26

	p1 := Pos{Lat: l1, Lon: L1, Ele: float64(0)}
	p2 := Pos{Lat: l2, Lon: L2, Ele: float64(0)}
	fmt.Printf("%.2f km\n", Dist(p1, p2))
	return Dist(p1, p2)
}

func cli_calc_effort(km float64, denivPos float64, denivNeg float64, vitessePlat float64) {
	km_effort := CalcDistanceEffort(km, denivPos, denivNeg)
	duration := km_effort / vitessePlat
	h, m := FloatToHourMin(duration)
	fmt.Println(h, "h", m, "m")
}
