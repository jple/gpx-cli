package main

import (
	"github.com/jple/gpx-cli/cmd"
<<<<<<< HEAD
	. "github.com/jple/gpx-cli/internal/gpx"
	"github.com/jple/gpx-cli/internal/summary"
=======
>>>>>>> refacto/reorg_export_stats
)

// func prettyprint(in any) string {
// 	j, err := json.MarshalIndent(in, "", "  ")
// 	if err != nil {
// 		// ./main.go:16:14: non-constant format string in call to log.Fatalf
// 		// log.Fatalf(err.Error())
// 	}
// 	return string(j)
// }

<<<<<<< HEAD
// TODO: move to test
func TestTrkptsLen() {
	flatSpeed := 4.5
	gpx := Gpx{}
	// gpx.Parse("core/test/data/split.gpx")
	gpx.Parse("core/test/data/npoints.gpx")
	gpxSummary := gpx.Summarize(flatSpeed)
=======
// // TODO: move to test
// func TestTrkptsLen() {
// 	flatSpeed := 4.5
// 	gpx := Gpx{}
// 	// gpx.Parse("core/test/data/split.gpx")
// 	gpx.Parse("core/test/data/npoints.gpx")
// 	gpxSummary := gpx.Stats(flatSpeed)
>>>>>>> refacto/reorg_export_stats

// 	var countingTrkpt = func(trk Trk) int {
// 		n := 0
// 		for _, trkseg := range trk.Trksegs {
// 			n += len(trkseg.Trkpts)
// 		}
// 		return n
// 	}

<<<<<<< HEAD
	var sumTrkptsLen = func(trkSummary summary.TrkSummary) int {
		var n int
		for _, section := range trkSummary.PerSection {
			n += section.Len()
		}
		return n
	}

	var trackLen = func(trkSummary summary.TrkSummary) int {
		return trkSummary.AllSections.Len()
	}

	want := countingTrkpt(gpx.Trks[0])
	have := sumTrkptsLen(gpxSummary.TrkSummaries[0].TrkSummary)
	have2 := trackLen(gpxSummary.TrkSummaries[0].TrkSummary)
	if have != want {
	}

	fmt.Printf("have2: %v\n", have2)
	fmt.Printf("have : %v trkpts\nwants : %v trkSummary.Len() \n", have, want)

	for _, section := range gpxSummary.TrkSummaries[0].PerSection {
		fmt.Println(section.From)
		fmt.Println(section.Len())
		fmt.Println(section.TotalAscent)
		// fmt.Printf("seg: %v, pt: %v\n", *section.FromTrksegId, *section.FromTrkptId)
	}
}

func test() {
	gpx := Gpx{}
	gpx.Parse("core/test/data/npoints.gpx")
	trkSummary := gpx.Trks[0].Summarize(0, 4.5)
	fmt.Printf("%+v\n", trkSummary)
	fmt.Println(prettyprint(trkSummary))
	// fmt.Println(trkSummary.ToString(summary.PrintArgs{ShowFromTo: true}))
}
=======
// 	var sumTrkptsLen = func(trkSummary summary.TrkSummary) int {
// 		var n int
// 		for _, section := range trkSummary.PerSection {
// 			n += section.Len()
// 		}
// 		return n
// 	}

// 	var trackLen = func(trkSummary summary.TrkSummary) int {
// 		return trkSummary.AllSections.Len()
// 	}

// 	want := countingTrkpt(gpx.Trks[0])
// 	have := sumTrkptsLen(gpxSummary.TrkSummaries[0].TrkSummary)
// 	have2 := trackLen(gpxSummary.TrkSummaries[0].TrkSummary)
// 	if have != want {
// 	}

// 	fmt.Printf("have2: %v\n", have2)
// 	fmt.Printf("have : %v trkpts\nwants : %v trkSummary.Len() \n", have, want)

// 	for _, section := range gpxSummary.TrkSummaries[0].PerSection {
// 		fmt.Println(section.From)
// 		fmt.Println(section.Len())
// 		fmt.Println(section.TotalAscent)
// 		// fmt.Printf("seg: %v, pt: %v\n", *section.FromTrksegId, *section.FromTrkptId)
// 	}
// }

// func test() {
// 	gpx := Gpx{}
// 	gpx.Parse("core/test/data/npoints.gpx")
// 	trkSummary := gpx.Trks[0].Stats(0, 4.5)
// 	fmt.Printf("%+v\n", trkSummary)
// 	fmt.Println(prettyprint(trkSummary))
// 	// fmt.Println(trkSummary.ToString(summary.PrintArgs{ShowFromTo: true}))
// }
>>>>>>> refacto/reorg_export_stats

// func testDist2Point() {
// 	filename := "src/test-info-d.gpx"
// 	from := "a"
// 	to := "c"
// 	speed := 4.5

// 	gpx := Gpx{}
// 	gpx.Parse(filename)

<<<<<<< HEAD
	fmt.Println(gpx.SummarizeBetweenTrkptsIndex(0, 4, speed).ToString())

	fmt.Println("============")
	if from == "" || to == "" {
		fmt.Println("from and to must be filled")
	} else {
		fmt.Println(gpx.SummarizeBetweenTrkptsNames(from, to, speed).ToString())
	}
=======
// 	fmt.Println(gpx.StatsBetweenTrkptsIndex(0, 4, speed).ToString())

// 	fmt.Println("============")
// 	if from == "" || to == "" {
// 		fmt.Println("from and to must be filled")
// 	} else {
// 		fmt.Println(gpx.StatsBetweenTrkptsNames(from, to, speed).ToString())
// 	}
>>>>>>> refacto/reorg_export_stats

// }

func main() {
	// test()
	// TestTrkptsLen()
	// sym.ShowUnicode()

	cmd.Execute()

	// testDist2Point()

<<<<<<< HEAD
	// gpx := Gpx{}
	// gpx.Parse("testfile/gr54-oisans-argentiere.gpx")
	// // csv := gpx.Summarize(4.5).ToCsv()
	// csv := gpx.Summarize(4.5).ToHTML()
=======
	// // csv := gpx.Stats(4.5).ToCsv()
	// csv := gpx.Stats(4.5).ToHTML()
>>>>>>> refacto/reorg_export_stats
	// fmt.Println(csv)

}
