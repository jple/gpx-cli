package geo

import "math"

const EarthRadius = 6371.0 // km

type Coord struct {
	Lat       float64 `xml:"lat,attr"`
	Lon       float64 `xml:"lon,attr"`
	Elevation float64 `xml:"ele,omitempty"`
}

// Dist returns disntace between p1 and p2 in km
func Dist(p1 Coord, p2 Coord) float64 {
	theta2 := DegToRad(p2.Lat)
	theta1 := DegToRad(p1.Lat)
	phi2 := DegToRad(p2.Lon)
	phi1 := DegToRad(p1.Lon)

	h := Haversin(theta2 - theta1)
	h += math.Cos(theta1) * math.Cos(theta2) * Haversin(phi2-phi1)
	out := EarthRadius * Ahaversin(h)
	return out
}
