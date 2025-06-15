package core

import (
	"fmt"
	"slices"
	"strconv"
)

type Trkseg struct {
	Trkpts Trkpts `xml:"trkpt"`
}
type Trk struct {
	Name    string   `xml:"name,omitempty"`
	Trksegs []Trkseg `xml:"trkseg"`

	// Should be optional
	Extensions struct {
		DenivPos float64 `xml:"DenivPos,omitempty"`
		DenivNeg float64 `xml:"DenivNeg,omitempty"`
		Distance float64 `xml:"Distance,omitempty"`
		// Conversion du denivele positif/negatif en km effort
		DenivPosEffort float64 `xml:"DenivPosEffort,omitempty"`
		DenivNegEffort float64 `xml:"DenivNegEffort,omitempty"`

		// Distance équivalente sur plat en incluant le dénivelé 
		DistanceEffort float64 `xml:"DistanceEffort,omitempty"`

		// Estimation de temps de marche
		Duration     float64 `xml:"Duration,omitempty"`
		DurationHour int8    `xml:"DurationHour,omitempty"`
		DurationMin  int8    `xml:"DurationMin,omitempty"`

		Line struct {
			Xmlns string `xml:"xmlns,attr,omitempty"`

			Color      string `xml:"color,omitempty"`
			Opacity    string `xml:"opacity,omitempty"`
			Weight     string `xml:"Weight,omitempty"`
			Width      int    `xml:"width,omitempty"`
			Linecap    string `xml:"linecap,omitempty"`
			Linejoin   string `xml:"linejoin,omitempty"`
			Dasharray  *int   `xml:"dasharray,omitempty"`
			Dashoffset int    `xml:"dashoffset,omitempty"`

			Extensions *struct {
				Jonction int `xml:"jonction,omitempty"`
			} `xml:"extensions,omitempty"`
		} `xml:"line,omitzero"`
	} `xml:"extensions,omitempty"`
}

func (trk Trk) GetLonLat() ([]string, []string) {
	var lons, lats []string

	// trkpts := slices.Concat(trk.Trksegs)[0].Trkpts
	var trkpts []Trkpt
	for _, trkseg := range trk.Trksegs {
		trkpts = slices.Concat(trkpts, trkseg.Trkpts)
	}
	for _, trkpt := range trkpts {
		lons = append(lons, strconv.FormatFloat(trkpt.Lon, 'f', -1, 64))
		lats = append(lats, strconv.FormatFloat(trkpt.Lat, 'f', -1, 64))
	}

	return lons, lats
}

func (p_trk *Trk) Reverse() Trk {
	trk := *p_trk

	slices.Reverse(trk.Trksegs)
	for _, trkseg := range trk.Trksegs {
		slices.Reverse(trkseg.Trkpts)
	}
	return trk
}

func (trk Trk) GetElevations() []float64 {
	trkpts := trk.GetTrkpts()
	return trkpts.GetElevations()
}

// Calculate cumulated distance between two index of trk
// TODO/refacto: move into Trkpts
func (trk Trk) GetDistanceFromTo(i, j int) float64 {
	if i >= j {
		fmt.Println("i must be < j")
		return 0.0
	}
	// var trkpts []Trkpt = slices.Concat(trk.Trksegs)[0].Trkpts
	var trkpts []Trkpt
	for _, trkseg := range trk.Trksegs {
		trkpts = slices.Concat(trkpts, trkseg.Trkpts)
	}
	var dist float64
	posPrev := Pt{
		Lon: trkpts[i].Lon,
		Lat: trkpts[i].Lat,
		Ele: trkpts[i].Ele,
	}
	for k, trkpt := range trkpts {
		if k <= i {
			continue
		}
		if k >= j {
			break
		}

		pos := Pt{
			Lon: trkpt.Lon,
			Lat: trkpt.Lat,
			Ele: trkpt.Ele,
		}
		dist += Dist(posPrev, pos)
		posPrev = pos
	}
	return dist
}

// Calculate cumulated distance for each trkpt
// (distance between trkpt[0] and trkpt[i])
func (trk Trk) GetCumulatedDistances() []float64 {
	trkpts := trk.GetTrkpts()
	return trkpts.GetCumulatedDistances()
}

func (trk Trk) GetRollElevations(winSize int, calc RollCalc) []float64 {
	return Rolling(trk.GetElevations(), winSize, calc)
}
func (trk Trk) GetRollDistances(winSize int, calc RollCalc) []float64 {
	return Rolling(trk.GetCumulatedDistances(), winSize, calc)
}

func (trk *Trk) AddName(name string) {
	trk.Name = name
}
