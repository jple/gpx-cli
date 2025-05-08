package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jple/gpx-cli/core"
)

type GpxTui struct {
	TrknameList core.TrknameList
	cursor      int
	Gpx         core.Gpx
	PrintInfo   bool
}

func (m GpxTui) Init() tea.Cmd {
	return nil
}

func (m GpxTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cursorMax int = len(m.TrknameList) - 1

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Return to default view with any key
		if m.PrintInfo {
			m.PrintInfo = false
			return m, nil
		}

		switch msg.String() {

		// =========== exit =============
		case "ctrl+c", "q":
			return m, tea.Quit

		// =========== move cursor =============
		case "up":
			if m.PrintInfo {
				return m, nil
			}
			if m.cursor > 0 {
				m.cursor -= 1
			}
			return m, nil
		case "down":
			if m.PrintInfo {
				return m, nil
			}
			if m.cursor < cursorMax {
				m.cursor += 1
			}
			return m, nil

		// =========== action =============
		case "m":
			selectedTrkId := m.TrknameList[m.cursor].Id
			m.Gpx = m.Gpx.Merge(selectedTrkId, selectedTrkId+1)
			m.TrknameList = m.Gpx.Ls(true)
			return m, nil
		case "s":
			selectedTrkname := m.TrknameList[m.cursor]
			if selectedTrkname.IsTrkpt() {
				m.Gpx = m.Gpx.SplitAtName(*selectedTrkname.TrkptName)
				m.TrknameList = m.Gpx.Ls(true)
			}
			return m, nil
		case "i":
			m.PrintInfo = true
			return m, nil
		}
	}

	return m, nil
}

func (m GpxTui) View() string {
	if m.PrintInfo {
		var printArgs core.PrintArgs = core.PrintArgs{AsciiFormat: true, Silent: true}
		return m.Gpx.GetInfo(true).Print(printArgs)
	} else {
		var s string
		for i, trkname := range m.TrknameList {
			if m.cursor == i {
				s += ">>> "
				// s += fmt.Sprintf(">>> **c:%v** ", m.cursor)
			}
			if trkname.IsTrk() {

				s += fmt.Sprintf("\u001b[1;32m%v\u001b[22;0m\n", trkname.Name)
				// s += fmt.Sprintf("%v\n", trkname.Name)
			}
			if trkname.IsTrkpt() {
				s += fmt.Sprintf("    %v\n", *trkname.TrkptName)
			}
		}

		s += "Press 'ctrl-c' or 'q' to exit..."
		return s
	}

}
