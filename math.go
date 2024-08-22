package main

import "math"

func Haversin(t float64) float64 {
	return math.Pow(math.Sin(t/2), 2)
}

func Ahaversin(t float64) float64 {
	return 2 * math.Asin(math.Sqrt(t))
}

func degToRad(deg float64) float64 {
	return 2 * math.Pi / 360 * deg
}

func dist(p1 Pos, p2 Pos) float64 {
	// output: km
	var R float64 = 6371

	theta2 := degToRad(p2.Lat)
	theta1 := degToRad(p1.Lat)
	phi2 := degToRad(p2.Lon)
	phi1 := degToRad(p1.Lon)

	h := Haversin(theta2 - theta1)
	h += math.Cos(theta1) * math.Cos(theta2) * Haversin(phi2-phi1)
	out := R * Ahaversin(h)
	return out

}

func diffElevation(p1 Pos, p2 Pos) float64 {
	return p2.Ele - p1.Ele
}

func calcDistanceEffort(distance float64, denivPos float64, denivNeg float64) float64 {
	denivPosEffort := denivPos / 100.0  // 100m D+ = 1km plat
	denivNegEffort := denivNeg / 300.0  // 300m D- = 1km plat
	return distance + denivPosEffort + math.Abs(denivNegEffort)
}