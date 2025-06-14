package core

// TODO: this function should replace Trk method
func GetEachTrkptDistance(trkpts []Trkpt) []float64 {
	var dists []float64
	posPrev := Pos{
		Lon: trkpts[0].Lon,
		Lat: trkpts[0].Lat,
		Ele: trkpts[0].Ele,
	}
	for _, trkpt := range trkpts {
		pos := Pos{
			Lon: trkpt.Lon,
			Lat: trkpt.Lat,
			Ele: trkpt.Ele,
		}
		dists = append(dists, Dist(posPrev, pos))
		posPrev = pos
	}
	return dists
}
func GetCumDistance(trkpts []Trkpt) float64 {
	var d float64 = 0
	dists := GetEachTrkptDistance(trkpts)
	for _, dist := range dists {
		d += dist
	}
	return d
}
