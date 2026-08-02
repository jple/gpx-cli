package gpx

import "github.com/jple/gpx-cli/internal/geo"

type Wpt struct {
	geo.Coord

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

/* NOTE(live thinking):
If we change this struct to :

type Trkpts []struct {
	Index int
	Trkpt
}

TrkptValue could be deleted
*/

// TODO: create generics for AddName
// TODO: remove ? set is not do idiomatic
func (trkpt *Trkpt) SetName(name string) *Trkpt {
	trkpt.Name = &name
	return trkpt
}

// ----------------------- query ----------------------
// Returns all trkpt Elevation
func (trkpts Trkpts) Elevations() []float64 {
	elevations := make([]float64, len(trkpts))
	for i, trkpt := range trkpts {
		elevations[i] = trkpt.Elevation
	}
	return elevations
}

// Returns distances between each successive trkpt
func (trkpts Trkpts) Distances() []float64 {
	distances := make([]float64, len(trkpts))
	for i := range trkpts {
		if i == 0 {
			continue
		}
		distances[i] = geo.Dist(trkpts[i].Coord, trkpts[i-1].Coord)
	}
	return distances
}

// Same as Distances(), but each value is cumulated to the previous one
func (trkpts Trkpts) CumulativeDistances() []float64 {
	distances := trkpts.Distances()
	for i := range distances {
		if i == 0 {
			continue
		}
		distances[i] += distances[i-1]
	}
	return distances
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
<<<<<<< HEAD
=======

const RollingWindowSize = 10

func checkIndex(i, j, n int) {
	if i < 0 {
		panic("i must be > 0")
	} else if i > j {
		panic("i must be <= j")
	} else if j >= n {
		panic("j must be < n")
	}
}

func (trkpts Trkpts) TotalDistance() float64 {
	return geo.Sum(trkpts.Distances())
}

// Calculate cumulated distance between two index of trk
func (trkpts Trkpts) TotalDistanceBetweenIndex(i, j int) float64 {
	checkIndex(i, j, len(trkpts))
	return trkpts[i : j+1].TotalDistance()
}

// TODO: to be include in cmd
func (trkpts Trkpts) TotalDistanceBetweenNames(from, to string) float64 {
	i := trkpts.FindName(from)
	j := trkpts.FindName(to)
	return trkpts.TotalDistanceBetweenIndex(i, j)
}

func (trkpts Trkpts) TotalAscent(rollingWindowSize int) float64 {
	return geo.TotalAscent(geo.Rolling(trkpts.Elevations(), rollingWindowSize, geo.Mean))
}
func (trkpts Trkpts) TotalDescent(rollingWindowSize int) float64 {
	return geo.TotalDescent(geo.Rolling(trkpts.Elevations(), rollingWindowSize, geo.Mean))
}
>>>>>>> refacto/reorg_export_stats
