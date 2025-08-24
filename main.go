package main

import (
	"encoding/json"
	"fmt"

	"github.com/jple/gpx-cli/cmd"
	. "github.com/jple/gpx-cli/core"
)

func prettyprint(in any) string {
	j, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		// ./main.go:16:14: non-constant format string in call to log.Fatalf
		// log.Fatalf(err.Error())
	}
	return string(j)
}

// TODO: move to test
func TestNPoints() {
	vitessePlat := 4.5
	gpx := Gpx{}
	// gpx.ParseFile("core/test/data/split.gpx")
	gpx.ParseFile("core/test/data/npoints.gpx")
	gpxSummary := gpx.GetInfo(vitessePlat)

	var countingTrkpt = func(trk Trk) int {
		n := 0
		for _, trkseg := range trk.Trksegs {
			n += len(trkseg.Trkpts)
		}
		return n
	}

	var sumTrkNPoints = func(trkSummary TrkSummary) int {
		var n int
		for _, section := range trkSummary.ListTrkptsSummary {
			n += section.NPoints
		}
		return n
	}

	var trackNPoints = func(trkSummary TrkSummary) int {
		return trkSummary.Track.NPoints
	}

	want := countingTrkpt(gpx.Trks[0])
	have := sumTrkNPoints(gpxSummary.Trks[0].TrkSummary)
	have2 := trackNPoints(gpxSummary.Trks[0].TrkSummary)
	if have != want {
	}

	fmt.Printf("have2: %v\n", have2)
	fmt.Printf("have : %v trkpts\nwants : %v trkSummary.NPoints\n", have, want)

	for _, section := range gpxSummary.Trks[0].ListTrkptsSummary {
		fmt.Println(section.From)
		fmt.Println(section.NPoints)
		fmt.Println(section.DenivPos)
		// fmt.Printf("seg: %v, pt: %v\n", *section.FromTrksegId, *section.FromTrkptId)
	}
}

func test() {
	gpx := Gpx{}
	gpx.ParseFile("core/test/data/npoints.gpx")
	trkSummary := gpx.Trks[0].GetInfo(0, 4.5)
	fmt.Printf("%+v\n", trkSummary)
	fmt.Println(prettyprint(trkSummary))
	// fmt.Println(trkSummary.ToString(PrintArgs{PrintFrom: true}))
}

func testDist2Point() {
	filename := "src/test-info-d.gpx"
	from := "a"
	to := "c"
	speed := 4.5

	gpx := Gpx{}
	gpx.ParseFile(filename)

	var printArgs PrintArgs = PrintArgs{AsciiFormat: true, PrintFrom: true}
	fmt.Println(gpx.GetInfoBetweenTrkptsId(0, 4, speed).ToString(printArgs))

	fmt.Println("============")
	if from == "" || to == "" {
		fmt.Println("from and to must be filled")
	} else {
		fmt.Println(gpx.GetInfoBetweenName(from, to, speed).ToString(printArgs))
	}

}

func main() {
	// test()
	// TestNPoints()
	// sym.ShowUnicode()

	cmd.Execute()

	// testDist2Point()

}
