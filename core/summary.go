package core

import "fmt"

type Pos struct {
	Lat float64
	Lon float64
	Ele float64

	Name string
}

type SectionInfo struct {
	TrkName        string
	From           string
	To             string
	NPoints        int
	VitessePlat    float64
	Distance       float64
	DenivPos       float64
	DenivNeg       float64
	DistanceEffort float64
	DurationHour   int8
	DurationMin    int8
}

type TrkSummary struct {
	Name    string
	Section []SectionInfo
}

type GpxSummary []TrkSummary

type PrintArgs struct {
	PrintFrom   bool
	AsciiFormat bool
}

func (s SectionInfo) Print(args PrintArgs) {
	if args.PrintFrom {
		if args.AsciiFormat {
			fmt.Printf("\u001b[4mSection:\u001b[24m \u001b[1;32m%v --> %v\u001b[22;0m\n", s.From, s.To)
		} else {
			fmt.Printf("Section: %v --> %v\n", s.From, s.To)
		}
	}

	fmt.Printf("Number of points:       %v\n", s.NPoints)
	fmt.Printf("Distance:               %.1f km\n", s.Distance)
	fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", s.DenivPos, s.DenivNeg)
	fmt.Printf("Distance effort:        %.1f km\n", s.DistanceEffort)
	fmt.Printf("Vitesse sur plat:       %.1f km/h\n", s.VitessePlat)
	fmt.Printf("Temps parcours estimé:  %vh%v\n", s.DurationHour, s.DurationMin)
}

func (trkSummary TrkSummary) Print(args PrintArgs) {
	trkName := trkSummary.Name
	fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trkName)
	fmt.Println()
	for _, sectionInfo := range trkSummary.Section {
		sectionInfo.Print(args)
		fmt.Println()
	}
}

func (gpxSummary GpxSummary) Print(printArgs PrintArgs) {
	for i, trkSummary := range gpxSummary {
		fmt.Printf("[%v] ", i)
		trkSummary.Print(printArgs)
	}
}
