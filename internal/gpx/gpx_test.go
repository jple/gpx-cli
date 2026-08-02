package gpx

import (
	"fmt"
	"math"
	"testing"
)

const findtrkpts_testset = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx>
	<metadata>
		<name>GR Belledonne Sud</name>
	</metadata>
	<trk>
		<name>Le Recoin</name>
		<trkseg>
			<trkpt lat="45" lon="5"> <name>a</name></trkpt>
            <trkpt lat="45.09" lon="5"></trkpt>
			<trkpt lat="45.18" lon="5"> <name>b</name></trkpt>
            <trkpt lat="45.27" lon="5"></trkpt>
			<trkpt lat="45.36" lon="5"> <name>c</name></trkpt>
        </trkseg>
    </trk>

	<trk>
		<name>Deuxieme</name>
		<trkseg>
			<trkpt lat="45.45" lon="5"><name>d</name></trkpt>
			<trkpt lat="45.54" lon="5"></trkpt>
			<trkpt lat="45.63" lon="5"><name>e</name></trkpt>
			<trkpt lat="45.72" lon="5"></trkpt>
			<trkpt lat="45.81" lon="5"><name>f</name></trkpt>
        </trkseg>
    </trk>

    <trk>
		<name>Trois</name>
		<trkseg>
			<trkpt lat="45.90" lon="5"><name>g</name></trkpt>
			<trkpt lat="45.99" lon="5"></trkpt>
        </trkseg>
    </trk>
</gpx>
`

func TestFindTrkptsId(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(findtrkpts_testset))

	var tests = []struct {
		trkId, segId, ptId int
		want               int
	}{
		{0, 0, 0, 0},
		{0, 0, 2, 2},
		{1, 0, 2, 7},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("inputs (%v, %v, %v)", test.trkId, test.segId, test.ptId)
		t.Run(testname, func(t *testing.T) {
			have := gpx.FindTrkptsId(test.trkId, test.segId, test.ptId)
			// fmt.Printf("Have pt : %v\nWant pt : %v\n",
			// 	gpx.AllTrkpts()[have],
			// 	gpx.Trks[test.trkId].Trksegs[test.segId].Trkpts[test.ptId],
			// )
			if have != test.want {
				t.Errorf("Have %v, test.want %v", have, test.want)
			}
		})
	}
}

<<<<<<< HEAD
func TestSummarizeBetweenTrkptsNames(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(findtrkpts_testset))

	have := int(math.Floor(gpx.SummarizeBetweenTrkptsNames("b", "e", 4.5).Distance))
=======
func TestStatsBetweenTrkptsNames(t *testing.T) {
	gpx := Gpx{}
	gpx.Parse([]byte(findtrkpts_testset))

	have := int(math.Floor(gpx.StatsBetweenTrkptsNames("b", "e", 4.5).Distance))
>>>>>>> refacto/reorg_export_stats
	want := 50
	if have != 50 {
		t.Fatalf("Have distance %v, want %v\n", have, want)
	}
}
