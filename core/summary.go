package core

import (
	"fmt"

	sym "github.com/jple/text-symbol"
)

type Pos struct {
	Lat float64
	Lon float64
	Ele float64

	Name string
}

// Section summaries either
// - the whole trk
// - a section between a trkpt name and the next one (no matter trkseg)
type SectionInfo struct {

	// TODO: Dupplicates on TrkSummary
	TrkId       int
	TrkName     string
	VitessePlat float64

	// TrkptName
	From         string
	FromTrksegId *int
	FromTrkptId  *int
	To           string

	// Cumulative values between "From" and "To"
	NPoints        int
	Distance       float64
	DenivPos       float64
	DenivNeg       float64
	DistanceEffort float64
	DurationHour   int8
	DurationMin    int8
}

type TrkSummary struct {
	Id      int
	Name    string
	Section []SectionInfo
	Track   SectionInfo
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
	str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
		trkSummary.Track.NPoints,
		sym.ArrowIconLeftRight(), trkSummary.Track.Distance,
		sym.UpAndDown(), trkSummary.Track.DenivPos, trkSummary.Track.DenivNeg,
		sym.ArrowWaveRight(), trkSummary.Track.DistanceEffort,
		sym.StopWatch(), trkSummary.Track.DurationHour, trkSummary.Track.DurationMin)

	// TODO: rename PrintFrom
	// this parameter is actually used to print details or not
	// if args.PrintFrom {
	// 	str += "\n"
	// }
	for _, sectionInfo := range trkSummary.Section {
		str += sectionInfo.ToString(args)
	}

	return str
}

func (s SectionInfo) ToString(args PrintArgs) string {
	var str string
	// TODO: rename PrintFrom
	// this parameter is actually used to print details or not
	if args.PrintFrom {
		if args.AsciiFormat {
			// str += fmt.Sprintf("      --> %v", sym.Green(s.To))
			str += fmt.Sprintf("      %v --> %v", sym.Green(s.From), sym.Green(s.To))
		} else {
			str += fmt.Sprintf("      --> %v", s.To)
		}

		str += fmt.Sprintf("\t(%v %vh%02d)\n",
			sym.StopWatch(), s.DurationHour, s.DurationMin)
		// str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
		// 	s.NPoints,
		// 	sym.ArrowIconLeftRight(), s.Distance,
		// 	sym.UpAndDown(), s.DenivPos, s.DenivNeg,
		// 	sym.ArrowWaveRight(), s.DistanceEffort,
		// 	sym.StopWatch(), s.DurationHour, s.DurationMin)

		str += "\n"
	}

	return str
}
