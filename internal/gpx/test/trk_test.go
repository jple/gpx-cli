package core

import (
	"testing"
	// . "github.com/jple/gpx-cli/internal/gpx"
)

func TestTrkpts(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse("data/npoints.gpx")
	trk := gpx.Trks[0]
	trkpts := trk.AllTrkpts()
	have := len(trkpts)
	want := 9
	if want != have {
		t.Fatalf("Have: %v, Want: %v\n", have, want)
	}
}

func TestTrkSections(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse("data/npoints.gpx")
	TrkSections := gpx.Trks[0].SplitByTrkptName()
	want := 3
	have := len(TrkSections)
	if want != have {
		t.Fatalf("Total number of TrkSections: Have: %v, Want: %v\n", have, want)
	}
	for i, s := range TrkSections {
		want := 3 + 1
		if i == len(TrkSections)-1 {
			want = 3
		}

		have := len(s)
		if want != have {
			t.Fatalf("In section %v, Have: %v pts, Want: %v pts\n", i, have, want)
		}
	}
}
