package core

import (
	"math"
	"testing"
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

func TestGetTotalDistance(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))

	const want float64 = 60.0
	var have float64

	have = math.Round(gpx.Trks[0].GetTrkpts().GetTotalDistance())
	if want != have {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}
}

func TestCumAscentDescent(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))

	const want float64 = 300
	var have float64

	have = CumAscent(gpx.Trks[0].GetTrkpts().GetElevations())
	if want != have {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}

	have = CumDescent(gpx.Trks[0].GetTrkpts().GetElevations())
	if want != -have {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}
}

func TestGetDistanceFromTo(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(gpxtestset))
	from := 1
	to := 4
	have := math.Round(gpx.Trks[0].GetTrkpts().GetTotalDistanceFromTo(from, to))
	want := float64(to-from) * 10

	if have != want {
		t.Fatalf("Have %v, want %v", have, want)
	}
}

// TODO: WIP
// func TestGetCumulatedDistances(t *testing.T) {
// 	gpx := Gpx{}
// 	gpx.Parse([]byte(gpxtestset))

// 	// const want float64 = 60.0
// 	var have []float64

// 	have = gpx.Trks[0].GetTrkpts().GetCumulatedDistances()
// 	fmt.Println(have)
// 	// if want != have {
// 	// 	t.Fatalf("Have: %v, Want: %v", have, want)
// 	// }
// }
