package core

import (
	"testing"

	. "github.com/jple/gpx-cli/core"
)

func TestSplit(t *testing.T) {
	gpx := Gpx{}
	gpx.ParseFile("data/split.gpx")
	// gpx.GetInfo(true).ToString(PrintArgs{})

	gpx1 := gpx.SplitAtName("first")
	gpx2 := gpx.SplitAtName("between")
	gpx3 := gpx.SplitAtName("last")

	satisfy := func(prefixText string, have, want int) {
		if have != want {
			t.Fatalf("%v expect value %v, but have %v", prefixText, have, want)
		}
	}

	satisfy("len(gpx1.Trk)", len(gpx1.Trk), len(gpx.Trk))
	satisfy("len(gpx2.Trk)", len(gpx2.Trk), len(gpx.Trk)+1)
	satisfy("len(gpx3.Trk)", len(gpx3.Trk), len(gpx.Trk)+1)

	satisfy("len(gpx2.Trk[0].Trkseg)", len(gpx2.Trk[0].Trkseg), 1)
	satisfy("len(gpx3.Trk[0].Trkseg)", len(gpx3.Trk[0].Trkseg), len(gpx.Trk[0].Trkseg)-1)
}
