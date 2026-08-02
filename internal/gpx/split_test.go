package gpx

import (
	"testing"
)

const split_trkpts_testset = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx>
	<trk>
		<trkseg>
			<trkpt lat="45" lon="5"><name>first</name></trkpt>
        </trkseg>
    </trk>
	<trk>
		<trkseg>
			<trkpt lat="45.45" lon="5"></trkpt>
        </trkseg>
		<trkseg>
			<trkpt lat="45.45" lon="5"><name>between</name></trkpt>
        </trkseg>
    </trk>
    <trk>
		<trkseg>
			<trkpt lat="45.45" lon="5"></trkpt>
			<trkpt lat="45.90" lon="5"><name>last</name></trkpt>
        </trkseg>
		<trkseg>
			<trkpt lat="45.45" lon="5"></trkpt>
        </trkseg>
    </trk>
</gpx>
`

func TestSplitAtName(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(split_trkpts_testset))
	// gpx.Parse("data/split.gpx")
	// gpx.Summarize(0.0).ToString(PrintArgs{})

	gpx1 := gpx
	gpx2 := gpx
	gpx3 := gpx

	gpx1.SplitAtName("first")
	gpx2.SplitAtName("between")
	gpx3.SplitAtName("last")

	satisfy := func(prefixText string, have, want int) {
		if have != want {
			t.Fatalf("%v expect value %v, but have %v\n", prefixText, want, have)
		}
	}

	satisfy("len(gpx1.Trks)", len(gpx1.Trks), len(gpx.Trks))   // no split done
	satisfy("len(gpx2.Trks)", len(gpx2.Trks), len(gpx.Trks)+1) // 1 split done
	satisfy("len(gpx3.Trks)", len(gpx3.Trks), len(gpx.Trks)+1) // 1 split done

	// every trk must have exactly 1 trkseg
	satisfy("len(gpx2.Trks[1].Trksegs)", len(gpx2.Trks[1].Trksegs), 1)
	satisfy("len(gpx2.Trks[2].Trksegs)", len(gpx2.Trks[2].Trksegs), 1)
	satisfy("len(gpx2.Trks[3].Trksegs)", len(gpx2.Trks[2].Trksegs), 1)
	satisfy("len(gpx3.Trks[1].Trksegs)", len(gpx3.Trks[2].Trksegs), 1)
	satisfy("len(gpx3.Trks[2].Trksegs)", len(gpx3.Trks[2].Trksegs), 1)
	satisfy("len(gpx3.Trks[3].Trksegs)", len(gpx3.Trks[2].Trksegs), 1)

	// split on gpx2 creates 2 trk with 1 trkpt
	satisfy("len(gpx2.Trks[1].AllTrkpts())", len(gpx2.Trks[1].AllTrkpts()), 1)
	satisfy("len(gpx2.Trks[2].AllTrkpts())", len(gpx2.Trks[2].AllTrkpts()), 1)

	// split on gpx3 creates 2 trk with 1 trkpt & 2 trkpt
	satisfy("len(gpx3.Trks[2].AllTrkpts())", len(gpx3.Trks[2].AllTrkpts()), 1)
	satisfy("len(gpx3.Trks[3].AllTrkpts())", len(gpx3.Trks[3].AllTrkpts()), 2)

}
