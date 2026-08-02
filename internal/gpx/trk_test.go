package gpx

import (
	"math"
	"testing"

	"github.com/jple/gpx-cli/internal/geo"
	// . "github.com/jple/gpx-cli/internal/gpx"
)

const gpxtestset = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx>
  <trk>
    <trkseg>
      <trkpt lat="45" lon="5"><ele>0</ele></trkpt>
      <trkpt lat="45.09" lon="5"><ele>100</ele></trkpt>
      <trkpt lat="45.18" lon="5"><ele>200</ele></trkpt>
      <trkpt lat="45.27" lon="5"><ele>300</ele></trkpt>

      <trkpt lat="45.18" lon="5"><ele>200</ele></trkpt>
      <trkpt lat="45.09" lon="5"><ele>100</ele></trkpt>
      <trkpt lat="45" lon="5"><ele>0</ele></trkpt>
    </trkseg>
  </trk>
</gpx>
`

func TestTotalDistance(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))

	const want float64 = 60.0
	var have float64

	have = math.Round(gpx.Trks[0].AllTrkpts().TotalDistance())
	if want != have {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}
}

func TestTotalAscentDescent(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))

	const rollingWindow = 3
	trkpts := gpx.Trks[0].AllTrkpts()
	var want float64 = math.Round(1000 * geo.TotalAscent(geo.Rolling(trkpts.Elevations(), rollingWindow, geo.Mean)))
	var have float64

	have = math.Round(1000 * trkpts.TotalAscent(rollingWindow))
	if have != want {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}

	want = math.Round(1000 * geo.TotalDescent(geo.Rolling(trkpts.Elevations(), rollingWindow, geo.Mean)))
	have = math.Round(1000 * gpx.Trks[0].AllTrkpts().TotalDescent(rollingWindow))
	if have != want {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}
}

func TestGetDistanceFromTo(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))
	from := 1
	to := 4
	have := math.Round(gpx.Trks[0].AllTrkpts().TotalDistanceBetweenIndex(from, to))
	want := float64(to-from) * 10

	if have != want {
		t.Fatalf("Have %v, want %v", have, want)
	}
}

func TestCumulativeDistances(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))

	var want []float64
	var have []float64

	// Calculate have
	have = gpx.Trks[0].AllTrkpts().CumulativeDistances()

	// Calculate want
	for _, trk := range gpx.Trks {
		trkpts := trk.AllTrkpts()
		for i, trkpt := range trkpts {
			if i == 0 {
				want = append(want, 0.0)
			} else {
				want = append(want, want[i-1]+geo.Dist(trkpt.Coord, trkpts[i-1].Coord))
			}
		}
	}

	// Compare
	if len(have) != len(want) {
		t.Fatalf("Length differs from have (%v) and want (%v)\n", len(have), len(want))
	}
	for i := 0; i < len(have); i++ {
		if have[i] != want[i] {
			t.Fatalf("Value at index %v differs from have (%v) with want (%v)\n", i, have[i], want[i])
		}
	}

}
