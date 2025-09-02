package core

import (
	"testing"

	. "github.com/jple/gpx-cli/core"
)

func TestTrkpts(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse("data/npoints.gpx")
	trk := gpx.Trks[0]
	trkpts := trk.GetTrkpts()
	have := len(trkpts)
	want := 9
	if want != have {
		t.Fatalf("Have: %v, Want: %v\n", have, want)
	}
}

func TestListTrkpts(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse("data/npoints.gpx")
	sections := gpx.Trks[0].GetListTrkptsPerName()
	want := 3
	have := len(sections)
	if want != have {
		t.Fatalf("Total number of sections: Have: %v, Want: %v\n", have, want)
	}
	for i, s := range sections {
		want := 3 + 1
		if i == len(sections)-1 {
			want = 3
		}

		have := len(s)
		if want != have {
			t.Fatalf("In section %v, Have: %v pts, Want: %v pts\n", i, have, want)
		}
	}
}
