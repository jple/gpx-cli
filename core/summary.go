package core

import (
	"fmt"

	sym "github.com/jple/text-symbol"
)

// TkrptsSummary summaries either
// - the whole trk
// - a section between a trkpt name and the next one (no matter trkseg)
type TrkptsSummary struct {
	// Tkr.Name or Trkpt.Name
	From string `json:",omitempty"`
	To   string `json:",omitempty"`

	// Cumulative values between "From" and "To"
	NPoints        int     `json:",omitempty"`
	Distance       float64 `json:",omitempty"`
	DenivPos       float64 `json:",omitempty"`
	DenivNeg       float64 `json:",omitempty"`
	DistanceEffort float64 `json:",omitempty"`
	DurationHour   int8    `json:",omitempty"`
	DurationMin    int8    `json:",omitempty"`
}

type TrkSummary struct {
	ListTrkptsSummary []TrkptsSummary
	Track             TrkptsSummary
}

type GpxSummary struct {
	VitessePlat float64 `json:",omitempty"`

	Trks []struct {
		Name string
		TrkSummary
	}
}

// TODO/refacto: rename or delete. Poor readability
type PrintArgs struct {
	PrintFrom   bool
	AsciiFormat bool
	Silent      bool
}

func (gpxSummary GpxSummary) ToString(args PrintArgs) string {
	var str string
	for i, trk := range gpxSummary.Trks {
		str += fmt.Sprintf("[%v] %v: %v", i, sym.Underline("Etape"), sym.Green(trk.Name))
		str += trk.TrkSummary.ToString(args)
	}
	// TODO: remove this
	if !args.Silent {
		fmt.Printf(str)
	}
	return str
}

func (trkSummary TrkSummary) ToString(args PrintArgs) string {
	var str string

	str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
		trkSummary.Track.NPoints,
		sym.ArrowIconLeftRight(), trkSummary.Track.Distance,
		sym.UpAndDown(), trkSummary.Track.DenivPos, trkSummary.Track.DenivNeg,
		sym.ArrowWaveRight(), trkSummary.Track.DistanceEffort,
		sym.StopWatch(), trkSummary.Track.DurationHour, trkSummary.Track.DurationMin)

	for _, sectionInfo := range trkSummary.ListTrkptsSummary {
		str += sectionInfo.ToString(args)
	}
	// TODO/rename: poor readability. this parameter is actually used to print details or not
	if args.PrintFrom {
		str += "\n"
	}

	return str
}

func (s TrkptsSummary) ToString(args PrintArgs) string {
	var str string
	// TODO/rename: poor readability. this parameter is actually used to print details or not
	if args.PrintFrom {
		if args.AsciiFormat {
			// str += fmt.Sprintf("      --> %v", sym.Green(s.To))
			str += fmt.Sprintf("      %v --> %v", sym.Green(s.From), sym.Green(s.To))
		} else {
			str += fmt.Sprintf("      --> %v", s.To)
			// str += fmt.Sprintf("      %v --> %v", s.From, s.To)
		}

		// str += fmt.Sprintf("\t(%v %vh%02d)\n",
		// 	sym.StopWatch(), s.DurationHour, s.DurationMin)
		str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
			s.NPoints,
			sym.ArrowIconLeftRight(), s.Distance,
			sym.UpAndDown(), s.DenivPos, s.DenivNeg,
			sym.ArrowWaveRight(), s.DistanceEffort,
			sym.StopWatch(), s.DurationHour, s.DurationMin)

	}

	return str
}
