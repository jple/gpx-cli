<<<<<<< HEAD
//go:build exclude

// TODO: reset after summary refacto

=======
>>>>>>> refacto/simplify_summary_struct
package tui

import (
	"errors"
	"fmt"
<<<<<<< HEAD
	"os"

	"github.com/jple/gpx-cli/internal/gpx"
<<<<<<< HEAD
	"github.com/jple/gpx-cli/internal/summary"
=======
	"github.com/jple/gpx-cli/internal/gpx/summary"
>>>>>>> refacto/reorg_export_stats

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
<<<<<<< HEAD
	summary.TrkptsSummary
=======
	summary.GeoStatistic
>>>>>>> refacto/reorg_export_stats
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
<<<<<<< HEAD
				Section{TrkptsSummary: section, TrkId: i},
=======
				Section{GeoStatistic: section, TrkId: i},
>>>>>>> refacto/reorg_export_stats
			)
		}
	}
=======
	"io"
	"os"
	"path"
	"strings"

	g "github.com/jple/gpx-cli/internal/gpx"
	"github.com/jple/gpx-cli/internal/gpx/geostat"

	tea "github.com/charmbracelet/bubbletea"
)

const tmpTuiDir string = "tmpTui"

type GpxTui struct {
	cursor, cursorMax int

	flatspeed float64
	Gpx       g.Gpx
	TrkStats  []TrkStat
}

type TrkStat struct {
	geostat.GeoStatistic
	TrkId     int
	isSection bool
}

func (m *GpxTui) computeStats() {
	// reset
	m.TrkStats = nil

	// update TrkStats
	g.ExportStatsTrks(io.Discard, m.Gpx.Trks, m.flatspeed, true, "string", // always export detail, in string mode
		func(stat geostat.GeoStatistic, trkid int, isSection bool) {
			m.TrkStats = append(m.TrkStats, TrkStat{stat, trkid, isSection})
		})

	m.cursorMax = max(0, len(m.TrkStats)-1)

	// NOTE: this prevent cursor being out of range when merging the last line
	if m.cursor > m.cursorMax {
		m.cursor = m.cursorMax
	}
}

func (m *GpxTui) Init() tea.Cmd {
	m.flatspeed = 4.5
	m.computeStats()
	return nil
}

func (m *GpxTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
>>>>>>> refacto/simplify_summary_struct

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
<<<<<<< HEAD
			}
			return m, nil
		case "down":
			if m.cursor < cursorMax {
				m.cursor += 1
=======
			} else {
				m.cursor = m.cursorMax
			}
			return m, nil
		case "down":
			if m.cursor < m.cursorMax {
				m.cursor++
			} else {
				m.cursor = 0
>>>>>>> refacto/simplify_summary_struct
			}
			return m, nil

		// =========== action =============
		case "left":
<<<<<<< HEAD
			selectedTrkId := TrkSections[m.cursor].TrkId
			if selectedTrkId > 0 && selectedTrkId < len(m.Gpx.Trks) {
				m.Gpx = m.Gpx.Merge(selectedTrkId-1, selectedTrkId)
<<<<<<< HEAD
				m.GpxSummary = m.Gpx.Summarize(flatSpeed)
=======
				m.GpxSummary = m.Gpx.Stats(flatSpeed)
>>>>>>> refacto/reorg_export_stats
			}
			return m, nil
		case "right":
			m.Gpx.SplitAtName(TrkSections[m.cursor].To)
<<<<<<< HEAD
			m.GpxSummary = m.Gpx.Summarize(flatSpeed)
=======
			m.GpxSummary = m.Gpx.Stats(flatSpeed)
>>>>>>> refacto/reorg_export_stats
			return m, nil
		case "s":
			filename := "tata/0.gpx"
			for i := 0; FileExists(filename); i++ {
				filename = fmt.Sprintf("tata/%v.gpx", i)
=======
			selectedTrkId := m.TrkStats[m.cursor].TrkId
			if selectedTrkId > 0 {
				m.Gpx = m.Gpx.Merge(selectedTrkId-1, selectedTrkId)
			}
			m.computeStats()
			return m, nil
		case "right":
			stat := m.TrkStats[m.cursor]
			if stat.isSection {
				m.Gpx.SplitAtName(stat.To)
			}
			m.computeStats()
			return m, nil
		case "s":
			filename := path.Join(tmpTuiDir, "0.gpx")
			for i := 0; FileExists(filename); i++ {
				filename = path.Join(tmpTuiDir, fmt.Sprintf("%v.gpx", i))
>>>>>>> refacto/simplify_summary_struct
			}
			m.Gpx.Save(filename)
		}
	}

	return m, nil
}

<<<<<<< HEAD
=======
func (m *GpxTui) View() string {
	b := new(strings.Builder)
	b.WriteString("Actions list:\n")
	b.WriteString("- up/down : navigation\n")
	b.WriteString("- left/right : merge/split trk\n")
	b.WriteString("- s : save current gpx in folder '" + tmpTuiDir + "'\n")
	b.WriteString("Press 'ctrl-c' or 'q' to exit...\n")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("cursor: %v/%v\n", m.cursor, m.cursorMax))

	w := StringWriter{m.cursor, 0, b}
	g.ExportStatsTrks(&w, m.Gpx.Trks, m.flatspeed, true, "string", nil)
	return w.builder.String()

}

>>>>>>> refacto/simplify_summary_struct
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}
<<<<<<< HEAD

func (m GpxTui) View() string {
	var str string

<<<<<<< HEAD
	var TrkSections []summary.TrkptsSummary
=======
	var TrkSections []summary.GeoStatistic
>>>>>>> refacto/reorg_export_stats
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
=======
>>>>>>> refacto/simplify_summary_struct
