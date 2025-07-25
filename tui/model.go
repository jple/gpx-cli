package tui

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jple/gpx-cli/core"
	sym "github.com/jple/text-symbol"
)

type GpxTui struct {
	GpxSummary core.GpxSummary
	cursor     int
	Gpx        core.Gpx
	PrintInfo  bool
}

func (m GpxTui) Init() tea.Cmd {
	return nil
}

func (m GpxTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cursorMax int
	var sections []core.TrkptsSummary
	// Note: cursor is only going through sections, not track name
	for _, trkSummary := range m.GpxSummary {
		cursorMax += len(trkSummary.ListTrkptsSummary)
		for _, section := range trkSummary.ListTrkptsSummary {
			sections = append(sections, section)
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
			selectedTrkId := sections[m.cursor].TrkId
			if selectedTrkId > 0 && selectedTrkId < len(m.Gpx.Trks) {
				m.Gpx = m.Gpx.Merge(selectedTrkId-1, selectedTrkId)
				m.GpxSummary = m.Gpx.GetInfo()
			}
			return m, nil
		case "right":
			m.Gpx = m.Gpx.SplitAtName(sections[m.cursor].To)
			m.GpxSummary = m.Gpx.GetInfo()
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

	var sections []core.TrkptsSummary
	// Note: cursor is only going through sections, not track name
	for _, trkSummary := range m.GpxSummary {
		for _, section := range trkSummary.ListTrkptsSummary {
			sections = append(sections, section)
		}
	}

	var k int
	for _, trkSummary := range m.GpxSummary {
		str += fmt.Sprintf("%v: %v", sym.Underline("Etape"), sym.Green(trkSummary.Name))
		str += fmt.Sprintf("\t(%v pts, %v %.0fkm, %v +%.0fm/%.0fm | %v %.0fkm_e, %v %vh%02d)\n",
			trkSummary.Track.NPoints,
			sym.ArrowIconLeftRight(), trkSummary.Track.Distance,
			sym.UpAndDown(), trkSummary.Track.DenivPos, trkSummary.Track.DenivNeg,
			sym.ArrowWaveRight(), trkSummary.Track.DistanceEffort,
			sym.StopWatch(), trkSummary.Track.DurationHour, trkSummary.Track.DurationMin)
		// str += "\n"

		for _, sectionInfo := range trkSummary.ListTrkptsSummary {
			if m.cursor == k {
				// NOTE: background never applies to tab to terminal behaviour
				str += sym.BgBrightGreen(sectionInfo.ToString(core.PrintArgs{PrintFrom: true, AsciiFormat: true}))
			} else {
				str += sectionInfo.ToString(core.PrintArgs{PrintFrom: true, AsciiFormat: true})
			}
			k += 1
		}
	}

	str += "Press 'ctrl-c' or 'q' to exit..."

	return str
	// }

}
