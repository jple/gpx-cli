package tui

import (
	"errors"
	"fmt"
	"os"

	"github.com/jple/gpx-cli/internal/gpx"
	"github.com/jple/gpx-cli/internal/summary"

	tea "github.com/charmbracelet/bubbletea"
	sym "github.com/jple/text-symbol"
)

type GpxTui struct {
	GpxSummary summary.GpxSummary
	cursor     int
	Gpx        gpx.Gpx
	PrintInfo  bool
}

func (m GpxTui) Init() tea.Cmd {
	return nil
}

// TODO: redondunt name !
type Section struct {
	summary.TrkptsSummary
	TrkId int
}

func (m GpxTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	flatSpeed := 4.5
	var cursorMax int
	var TrkSections []Section
	// Note: cursor is only going through TrkSections, not track name
	for i, trkSummary := range m.GpxSummary.TrkSummaries {
		cursorMax += len(trkSummary.PerSection)
		for _, section := range trkSummary.PerSection {
			TrkSections = append(TrkSections,
				Section{TrkptsSummary: section, TrkId: i},
			)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// =========== exit =============
		case "ctrl+c", "q":
			return m, tea.Quit

		// =========== move cursor =============
		case "up":
			if m.cursor > 0 {
				m.cursor -= 1
			}
			return m, nil
		case "down":
			if m.cursor < cursorMax {
				m.cursor += 1
			}
			return m, nil

		// =========== action =============
		case "left":
			selectedTrkId := TrkSections[m.cursor].TrkId
			if selectedTrkId > 0 && selectedTrkId < len(m.Gpx.Trks) {
				m.Gpx = m.Gpx.Merge(selectedTrkId-1, selectedTrkId)
				m.GpxSummary = m.Gpx.Summarize(flatSpeed)
			}
			return m, nil
		case "right":
			m.Gpx.SplitAtName(TrkSections[m.cursor].To)
			m.GpxSummary = m.Gpx.Summarize(flatSpeed)
			return m, nil
		case "s":
			filename := "tata/0.gpx"
			for i := 0; FileExists(filename); i++ {
				filename = fmt.Sprintf("tata/%v.gpx", i)
			}
			m.Gpx.Save(filename)
		}
	}

	return m, nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

func (m GpxTui) View() string {
	var str string

	var TrkSections []summary.TrkptsSummary
	// Note: cursor is only going through TrkSections, not track name
	for _, trkSummary := range m.GpxSummary.TrkSummaries {
		for _, section := range trkSummary.PerSection {
			TrkSections = append(TrkSections, section)
		}
	}

	var k int
	for _, trkSummary := range m.GpxSummary.TrkSummaries {
		str += fmt.Sprintf("%v: %v", sym.Underline("Etape"), sym.Green(trkSummary.TrkName))
		str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
			trkSummary.AllSections.Len(),
			sym.ArrowIconLeftRight(), trkSummary.AllSections.Distance,
			sym.UpAndDown(), trkSummary.AllSections.TotalAscent, trkSummary.AllSections.TotalDescent,
			sym.ArrowWaveRight(), trkSummary.AllSections.DistanceEffort,
			sym.StopWatch(), trkSummary.AllSections.DurationHour, trkSummary.AllSections.DurationMin)
		// str += "\n"

		for _, sectionSummary := range trkSummary.PerSection {
			if m.cursor == k {
				// NOTE: background never applies to tab to terminal behaviour
				str += sym.BgBrightGreen(sectionSummary.ToString())
			} else {
				str += sectionSummary.ToString()
			}
			k += 1
		}
	}

	str += "Press 'ctrl-c' or 'q' to exit..."

	return str
	// }

}
