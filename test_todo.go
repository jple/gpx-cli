package main

// // TODO: move to test
// func TestTrkptsLen() {
// 	flatSpeed := 4.5
// 	gpx := Gpx{}
// 	// gpx.Parse("core/test/data/split.gpx")
// 	gpx.Parse("core/test/data/npoints.gpx")
// 	gpxSummary := gpx.Stats(flatSpeed)

// 	var countingTrkpt = func(trk Trk) int {
// 		n := 0
// 		for _, trkseg := range trk.Trksegs {
// 			n += len(trkseg.Trkpts)
// 		}
// 		return n
// 	}

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

// func testDist2Point() {
// 	filename := "src/test-info-d.gpx"
// 	from := "a"
// 	to := "c"
// 	speed := 4.5

// 	gpx := Gpx{}
// 	gpx.Parse(filename)

// 	fmt.Println(gpx.StatsBetweenTrkptsIndex(0, 4, speed).ToString())

// 	fmt.Println("============")
// 	if from == "" || to == "" {
// 		fmt.Println("from and to must be filled")
// 	} else {
// 		fmt.Println(gpx.StatsBetweenTrkptsNames(from, to, speed).ToString())
// 	}

// }
