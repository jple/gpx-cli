package core

import (
	"fmt"

	sym "github.com/jple/text_symbol"
)

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

func (gpxSummary GpxSummary) ToString(args PrintArgs) string {
	var str string
	for i, trkSummary := range gpxSummary {
		str += fmt.Sprintf("[%v] ", i)
		str += trkSummary.ToString(args)
	}
	// TODO: remove this
	if !args.Silent {
		fmt.Printf(str)
	}
	return str
}

func (trkSummary TrkSummary) ToString(args PrintArgs) string {
	var str string

	trkName := trkSummary.Name
	str += fmt.Sprintf("%v: %v", sym.Underline("Etape"), sym.Green(trkName))
	// TODO: rename PrintFrom
	// this parameter is actually used to print details or not
	if args.PrintFrom {
		str += "\n"
	}
	for _, sectionInfo := range trkSummary.Section {
		str += sectionInfo.ToString(args) + ""
	}

	return str
}

func (s SectionInfo) ToString(args PrintArgs) string {
	var str string
	if args.PrintFrom {
		if args.AsciiFormat {
			str += fmt.Sprintf("\t--> %v", sym.Green(s.To))

		} else {
			str += fmt.Sprintf("%v --> %v", s.To)
		}
	}

	// fmt.Printf("Number of points:       %v\n", s.NPoints)
	// fmt.Printf("Distance:               %.1f km\n", s.Distance)
	// fmt.Printf("D+/D-:                  %.0f m / %.0f m\n", s.DenivPos, s.DenivNeg)
	// fmt.Printf("Distance effort:        %.1f km\n", s.DistanceEffort)
	// fmt.Printf("Vitesse sur plat:       %.1f km/h\n", s.VitessePlat)
	// fmt.Printf("Temps parcours estimé:  %vh%v\n", s.DurationHour, s.DurationMin)

	str += fmt.Sprintf("\t(%v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
		sym.ArrowIconLeftRight(), s.Distance,
		sym.UpAndDown(), s.DenivPos, s.DenivNeg,
		sym.ArrowWaveRight(), s.DistanceEffort,
		sym.StopWatch(), s.DurationHour, s.DurationMin)

	// str += fmt.Sprintf("Number of points:       %v\n", s.NPoints)
	// str += fmt.Sprintf("Distance:               %.1f km\n", s.Distance)
	// str += fmt.Sprintf("D+/D-:                  %.0f m / %.0f m\n", s.DenivPos, s.DenivNeg)
	// str += fmt.Sprintf("Distance effort:        %.1f km\n", s.DistanceEffort)
	// str += fmt.Sprintf("Vitesse sur plat:       %.1f km/h\n", s.VitessePlat)
	// str += fmt.Sprintf("Temps parcours estimé:  %vh%v\n", s.DurationHour, s.DurationMin)

	return str
}
