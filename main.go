package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jple/gpx-cli/cmd"
	. "github.com/jple/gpx-cli/core"
)

func prettyprint(in any) string {
	j, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(j)
}

func TestNPoints() {
	gpx := Gpx{}
	// gpx.ParseFile("core/test/data/split.gpx")
	gpx.ParseFile("core/test/data/npoints.gpx")
	gpxSummary := gpx.GetInfo()

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
	have := sumTrkNPoints(gpxSummary[0])
	have2 := trackNPoints(gpxSummary[0])
	if have != want {
	}

	fmt.Printf("have2: %v\n", have2)
	fmt.Printf("have : %v trkpts\nwants : %v trkSummary.NPoints\n", have, want)

	for _, section := range gpxSummary[0].ListTrkptsSummary {
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

func main() {
	// test()
	cmd.Execute()
	// TestNPoints()
	// sym.ShowUnicode()
}
