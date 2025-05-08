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
	Silent      bool
}

func (s SectionInfo) Print(args PrintArgs) string {
	var str string
	if args.PrintFrom {
		if args.AsciiFormat {
			str += fmt.Sprintf("\u001b[4mSection:\u001b[24m \u001b[1;32m%v --> %v\u001b[22;0m\n", s.From, s.To)
			// fmt.Printf("\u001b[4mSection:\u001b[24m \u001b[1;32m%v --> %v\u001b[22;0m\n", s.From, s.To)
		} else {
			str += fmt.Sprintf("Section: %v --> %v\n", s.From, s.To)
			// fmt.Printf("Section: %v --> %v\n", s.From, s.To)
		}
	}

	// fmt.Printf("Number of points:       %v\n", s.NPoints)
	// fmt.Printf("Distance:               %.1f km\n", s.Distance)
	// fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", s.DenivPos, s.DenivNeg)
	// fmt.Printf("Distance effort:        %.1f km\n", s.DistanceEffort)
	// fmt.Printf("Vitesse sur plat:       %.1f km/h\n", s.VitessePlat)
	// fmt.Printf("Temps parcours estimé:  %vh%v\n", s.DurationHour, s.DurationMin)

	str += fmt.Sprintf("Number of points:       %v\n", s.NPoints)
	str += fmt.Sprintf("Distance:               %.1f km\n", s.Distance)
	str += fmt.Sprintf("D+/D-:                  %.0f m / %.0f m\n", s.DenivPos, s.DenivNeg)
	str += fmt.Sprintf("Distance effort:        %.1f km\n", s.DistanceEffort)
	str += fmt.Sprintf("Vitesse sur plat:       %.1f km/h\n", s.VitessePlat)
	str += fmt.Sprintf("Temps parcours estimé:  %vh%v\n", s.DurationHour, s.DurationMin)

	return str
}

func (trkSummary TrkSummary) Print(args PrintArgs) string {
	var str string

	trkName := trkSummary.Name
	str += fmt.Sprintf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n\n", trkName)
	// fmt.Printf("\u001b[4mTrack name:\u001b[24m \u001b[1;32m%v\u001b[22;0m\n", trkName)
	// fmt.Println()
	for _, sectionInfo := range trkSummary.Section {
		str += sectionInfo.Print(args) + "\n"
		// fmt.Println()
	}

	return str
}

func (gpxSummary GpxSummary) Print(args PrintArgs) string {
	var str string
	for i, trkSummary := range gpxSummary {
		str += fmt.Sprintf("[%v] ", i)
		str += trkSummary.Print(args)
	}
	if !args.Silent {
		fmt.Printf(str)
	}
	return str

}
