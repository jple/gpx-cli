package core

import (
	"math"
)

type Pt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele float64 `xml:"ele,omitempty"`
}

func Dist(p1 Pt, p2 Pt) float64 {
	// output: km
	var R float64 = 6371

	theta2 := DegToRad(p2.Lat)
	theta1 := DegToRad(p1.Lat)
	phi2 := DegToRad(p2.Lon)
	phi1 := DegToRad(p1.Lon)

	h := Haversin(theta2 - theta1)
	h += math.Cos(theta1) * math.Cos(theta2) * Haversin(phi2-phi1)
	out := R * Ahaversin(h)
	return out

}

func DiffElevation(p1 Pt, p2 Pt) float64 {
	return p2.Ele - p1.Ele
}

type Wpt struct {
	Pt

	// NOTE: innerxml to prevent escaping (more readable, less secure :/)
	Name *string `xml:"name,omitempty"`
	// Name *string `xml:",innerxml"` // NOTE: parse not working !
	Type *string `xml:"type,omitempty"`
	Cmt  *string `xml:"cmt,omitempty"`
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

// TODO: create generics for AddName
func (trkpt *Trkpt) SetName(name string) {
	trkpt.Name = &name
}
func (trkpt *Trkpt) SetElevation(ele float64) {
	trkpt.Ele = ele
}

// Returns a []float64 containing a specific calculation between each trkpt (first value is 0)
func (trkpts Trkpts) MapDiff(f func(prev, curr Pt) float64) []float64 {
	var res []float64
	if len(trkpts) == 0 {
		return res
	}

	prev := trkpts[0].Pt
	for _, trkpt := range trkpts {
		res = append(res, f(prev, trkpt.Pt))
		prev = trkpt.Pt
	}
	return res
}

// Returns all trkpt Ele
func (trkpts Trkpts) GetElevations() []float64 {
	getCurrentEle := func(prev, curr Pt) float64 {
		return curr.Ele
	}
	return trkpts.MapDiff(getCurrentEle)
}

// Returns distance between each trkpt (first value set to 0)
func (trkpts Trkpts) GetCumulatedDistances() []float64 {
	var cumdist float64
	getCumDist := func(prev, curr Pt) float64 {
		cumdist += Dist(prev, curr)
		return cumdist
	}
	return trkpts.MapDiff(getCumDist)
}

func (trkpts Trkpts) GetTotalDistance() float64 {
	return Sum(trkpts.MapDiff(Dist))
}

func checkIndex(i, j, n int) {
	if i < 0 {
		panic("i must be > 0")
	} else if i > j {
		panic("i must be <= j")
	} else if j >= n {
		panic("j must be < n")
	}
}

// Calculate cumulated distance between two index of trk
func (trkpts Trkpts) GetTotalDistanceFromTo(i, j int) float64 {
	checkIndex(i, j, len(trkpts))
	return trkpts[i : j+1].GetTotalDistance()
}

func (trkpts Trkpts) GetTotalAscent() float64 {
	// return CumAscent(trkpts.GetElevations())
	//NOTE: testing
	return CumAscent(Rolling(trkpts.GetElevations(), 10, Mean))
	// return CumAscent(TrendSummary(Rolling(trkpts.GetElevations(), 5, Mean), 2).GetValues())
}
func (trkpts Trkpts) GetTotalDescent() float64 {
	// return CumDescent(trkpts.GetElevations())
	//NOTE: testing
	return CumDescent(Rolling(trkpts.GetElevations(), 10, Mean))
	// return CumDescent(TrendSummary(Rolling(trkpts.GetElevations(), 5, Mean), 2).GetValues())
}

func (trkpts Trkpts) FindName(name string) int {
	for i, trkpt := range trkpts {
		if trkpt.Name != nil && *trkpt.Name == name {
			return i
		}
	}
	panic(name + " not found in trkpts names")
	return -1
}

func (trkpts Trkpts) GetTotalDistanceFromToName(from, to string) float64 {
	i := trkpts.FindName(from)
	j := trkpts.FindName(to)
	return trkpts.GetTotalDistanceFromTo(i, j)
}

func (trkpts Trkpts) GetSummary(vitessePlat float64) TrkptsSummary {
	trkptsSummary := TrkptsSummary{
		From: "start",
		To:   "end",

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
