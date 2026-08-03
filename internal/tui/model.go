package tui

import (
	"errors"
	"fmt"
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
			} else {
				m.cursor = m.cursorMax
			}
			return m, nil
		case "down":
			if m.cursor < m.cursorMax {
				m.cursor++
			} else {
				m.cursor = 0
			}
			return m, nil

		// =========== action =============
		case "left":
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
			}
			m.Gpx.Save(filename)
		}
	}

	return m, nil
}

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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}
