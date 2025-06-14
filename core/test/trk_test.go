package core

import (
	"testing"

	. "github.com/jple/gpx-cli/core"
)

func TestFlattenTrkpts(t *testing.T) {
	gpx := Gpx{}
	gpx.ParseFile("data/npoints.gpx")
	trk := gpx.Trk[0]
	trkpts := trk.GetFlattenTrkpts()
	have := len(trkpts)
	want := 9
	if want != have {
		t.Fatalf("Have: %v, Want: %v", have, want)
	}
}

func TestSections(t *testing.T) {
	gpx := Gpx{}
	gpx.ParseFile("data/npoints.gpx")
	sections := gpx.Trk[0].GetSections()
	want := 3
	have := len(sections)
	if want != have {
		t.Fatalf("Total number of sections: Have: %v, Want: %v", have, want)
	}
	for i, s := range sections {
		want := 3
		have := len(s)
		if want != have {
			t.Fatalf("In section %v, Have: %v, Want: %v", i, have, want)
		}
	}
}
