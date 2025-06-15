package core

import "math"

type Pt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele float64 `xml:"ele,omitempty"`
}

type Wpt struct {
	Pt

	Name *string `xml:"name"`
	Type *string `xml:"type"`
	Cmt  *string `xml:"cmt"`
}

type Trkpt struct {
	Wpt

	Extensions *struct {
		TrkExtension struct {
			Visugpx string `xml:"visugpx,attr,omitempty"`
			Node    int    `xml:"node,omitempty"`
		} `xml:"TrkExtension,omitempty"`
	} `xml:"extensions,omitempty"`
}
type Trkpts []Trkpt
type ListTrkpts []Trkpts

// TODO: create generics for AddName
func (trkpt *Trkpt) AddName(name string) {
	trkpt.Name = &name
}
func (trkpt *Trkpt) AddElevation(ele float64) {
	trkpt.Ele = ele
}

// Returns a []float64 containing a specific calculation on each trkpt
func (trkpts Trkpts) Map(calculation func(posPrev, pos Pt) float64) []float64 {
	var res []float64
	if len(trkpts) == 0 {
		return res
	}

	posPrev := Pt{
		Lon: trkpts[0].Lon,
		Lat: trkpts[0].Lat,
		Ele: trkpts[0].Ele,
	}
	for _, trkpt := range trkpts {
		pos := Pt{
			Lon: trkpt.Lon,
			Lat: trkpt.Lat,
			Ele: trkpt.Ele,
		}
		res = append(res, calculation(posPrev, pos))
		posPrev = pos
	}
	return res
}

// Same as Map, but posPrev is replaced by p0
func (trkpts Trkpts) Map0(calculation func(unusedPos, pos Pt) float64) []float64 {
	var res []float64
	pos0 := Pt{
		Lon: trkpts[0].Lon,
		Lat: trkpts[0].Lat,
		Ele: trkpts[0].Ele,
	}
	for _, trkpt := range trkpts {
		pos := Pt{
			Lon: trkpt.Lon,
			Lat: trkpt.Lat,
			Ele: trkpt.Ele,
		}
		res = append(res, calculation(pos0, pos))
	}
	return res
}

// Returns all trkpt Ele
func (trkpts Trkpts) GetElevations() []float64 {
	getCurrentEle := func(posPrev, pos Pt) float64 {
		return pos.Ele
	}
	return trkpts.Map(getCurrentEle)
}

// Returns distance between each trkpt (first value set to 0)
func (trkpts Trkpts) GetDistances() []float64 {
	return trkpts.Map(Dist)
}

// Returns distance between each trkpt (first value set to 0)
func (trkpts Trkpts) GetCumulatedDistances() []float64 {
	return trkpts.Map0(Dist)
}

func (trkpts Trkpts) GetTotalDistance() float64 {
	var d float64 = 0
	dists := trkpts.GetDistances()
	for _, dist := range dists {
		d += dist
	}
	return d
}

func (trkpts Trkpts) GetDiffElevations() []float64 {
	return trkpts.Map(DiffElevation)
}
func (trkpts Trkpts) GetAscents() []float64 {
	diffElevations := trkpts.GetDiffElevations()
	var out []float64
	for _, diffEle := range diffElevations {
		out = append(out, math.Max(diffEle, 0))
	}
	return out
}
func (trkpts Trkpts) GetDescents() []float64 {
	diffElevations := trkpts.GetDiffElevations()
	var out []float64
	for _, diffEle := range diffElevations {
		out = append(out, math.Min(diffEle, 0))
	}
	return out
}

func (trkpts Trkpts) GetTotalAscent() float64 {
	ascents := trkpts.GetAscents()
	var out float64
	for _, v := range ascents {
		out += v
	}
	return out
}
func (trkpts Trkpts) GetTotalDescent() float64 {
	descents := trkpts.GetDescents()
	var out float64
	for _, v := range descents {
		out += v
	}
	return out
}

// func (trkpts Trkpts) GetFrom() string {
// 	return trkpts.From
// }
// func (trkpts *Trkpts) SetFrom(s string) {
// 	trkpts.From = s
// }
// func (trkpts Trkpts) GetTo() string {
// 	return trkpts.To
// }
// func (trkpts *Trkpts) SetTo(s string) {
// 	trkpts.To = s
// }

func (trkpts Trkpts) GetSummary(vitessePlat float64) TrkptsSummary {
	trkptsSummary := TrkptsSummary{
		VitessePlat: vitessePlat,
		From:        "start",
		To:          "end",

		NPoints:  len(trkpts),
		Distance: trkpts.GetTotalDistance(),
		DenivPos: trkpts.GetTotalAscent(),
		DenivNeg: trkpts.GetTotalDescent(),
	}

	// Set calculation value
	trkptsSummary.DistanceEffort = CalcDistanceEffort(
		trkptsSummary.Distance,
		trkptsSummary.DenivPos,
		trkptsSummary.DenivNeg)
	_, trkptsSummary.DurationHour, trkptsSummary.DurationMin =
		CalcDuration(
			trkptsSummary.DistanceEffort,
			vitessePlat)

	return trkptsSummary
}
