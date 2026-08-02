package summary

import (
	"fmt"

	sym "github.com/jple/text-symbol"
)

type GpxSummary struct {
	FlatSpeed float64 `json:",omitempty"`

	TrkSummaries []struct {
		TrkName string
		TrkSummary
	}
}

type TrkSummary struct {
	PerSection  []TrkptsSummary
	AllSections TrkptsSummary
}

// TrkptsSummary summaries either
// - the whole trk
// - a section between a trkpt name and the next one (no matter trkseg)
type TrkptsSummary struct {
	// Trk.Name or Trkpt.Name
	From string `json:",omitempty"`
	To   string `json:",omitempty"`

	// Cumulative values between "From" and "To"
	N              int     `json:",omitempty"` // N is the trkpts count
	Distance       float64 `json:",omitempty"`
	TotalAscent    float64 `json:",omitempty"`
	TotalDescent   float64 `json:",omitempty"`
	DistanceEffort float64 `json:",omitempty"`
	DurationHour   int8    `json:",omitempty"`
	DurationMin    int8    `json:",omitempty"`
}

func (s TrkptsSummary) Len() int {
	return s.N
}

func (gpxSummary GpxSummary) ToString(detail bool) string {
	var str string
	for i, trkSummaries := range gpxSummary.TrkSummaries {
		str += fmt.Sprintf("[%v] %v: %v", i, sym.Underline("Etape"), sym.Green(trkSummaries.TrkName))
		str += trkSummaries.TrkSummary.ToString(detail)
	}
	return str
}

func (trkSummary TrkSummary) ToString(detail bool) string {
	var str string

	str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
		trkSummary.AllSections.Len(),
		sym.ArrowIconLeftRight(), trkSummary.AllSections.Distance,
		sym.UpAndDown(), trkSummary.AllSections.TotalAscent, trkSummary.AllSections.TotalDescent,
		sym.ArrowWaveRight(), trkSummary.AllSections.DistanceEffort,
		sym.StopWatch(), trkSummary.AllSections.DurationHour, trkSummary.AllSections.DurationMin)

	if detail {
		for _, sectionSummary := range trkSummary.PerSection {
			str += sectionSummary.ToString()
		}
		// NOTE: additional line for readability between trks
		// if detail {
		str += "\n"
	}

	return str
}

func (s TrkptsSummary) ToString() string {
	str := fmt.Sprintf("      %v --> %v", sym.Green(s.From), sym.Green(s.To))

	// str += fmt.Sprintf("\t(%v %vh%02d)\n",
	// 	sym.StopWatch(), s.DurationHour, s.DurationMin)
	str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
		s.Len(),
		sym.ArrowIconLeftRight(), s.Distance,
		sym.UpAndDown(), s.TotalAscent, s.TotalDescent,
		sym.ArrowWaveRight(), s.DistanceEffort,
		sym.StopWatch(), s.DurationHour, s.DurationMin)
	return str
}
