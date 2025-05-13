package tui

import (
	"fmt"

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
	var sections []core.SectionInfo
	// Note: cursor is only going through sections, not track name
	for _, trksummary := range m.GpxSummary {
		cursorMax += len(trksummary.Section)
		for _, section := range trksummary.Section {
			sections = append(sections, section)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// // Return to default view with any key
		// if m.PrintInfo {
		// 	m.PrintInfo = false
		// 	return m, nil
		// }

		switch msg.String() {

		// =========== exit =============
		case "ctrl+c", "q":
			return m, tea.Quit

		// =========== move cursor =============
		case "up":
			// if m.PrintInfo {
			// 	return m, nil
			// }
			if m.cursor > 0 {
				m.cursor -= 1
			}
			return m, nil
		case "down":
			// if m.PrintInfo {
			// 	return m, nil
			// }
			if m.cursor < cursorMax {
				m.cursor += 1
			}
			return m, nil

		// =========== action =============
		case "m":
			selectedTrkId := sections[m.cursor].TrkId
			if selectedTrkId > 0 && selectedTrkId < len(m.Gpx.Trk) {
				m.Gpx = m.Gpx.Merge(selectedTrkId-1, selectedTrkId)
				m.GpxSummary = m.Gpx.GetInfo(true)
			}
			return m, nil
		case "s":
			// selectedTrkname := m.TrknameList[m.cursor]
			// if selectedTrkname.IsTrkpt() {
			// 	m.Gpx = m.Gpx.SplitAtName(*selectedTrkname.TrkptName)
			// 	m.TrknameList = m.Gpx.Ls(true)
			// }

			m.Gpx = m.Gpx.SplitAtName(sections[m.cursor].To)
			m.GpxSummary = m.Gpx.GetInfo(true)
			return m, nil
			// case "i":
			// 	m.PrintInfo = true
			// 	return m, nil
		}
	}

	return m, nil
}

func (m GpxTui) View() string {
	// if m.PrintInfo {
	// 	var printArgs core.PrintArgs = core.PrintArgs{AsciiFormat: true, Silent: true}
	// 	return m.Gpx.GetInfo(true).Print(printArgs)
	// } else {

	var str string

	var sections []core.SectionInfo
	// Note: cursor is only going through sections, not track name
	for _, trksummary := range m.GpxSummary {
		for _, section := range trksummary.Section {
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
		str += "\n"

		for _, sectionInfo := range trkSummary.Section {
			if m.cursor == k {
				// str += ">>> "
				str += fmt.Sprintf(">>> **c:%v** ", m.cursor)
			}
			str += sectionInfo.ToString(core.PrintArgs{PrintFrom: true, AsciiFormat: true})
			k += 1
		}
	}

	// for i, trkname := range m.TrknameList {
	// 	if m.cursor == i {
	// 		s += ">>> "
	// 		// s += fmt.Sprintf(">>> **c:%v** ", m.cursor)
	// 	}
	// 	if trkname.IsTrk() {

	// 		s += fmt.Sprintf("\u001b[1;32m%v\u001b[22;0m\n", trkname.Name)
	// 		// s += fmt.Sprintf("%v\n", trkname.Name)
	// 	}
	// 	if trkname.IsTrkpt() {
	// 		s += fmt.Sprintf("    %v\n", *trkname.TrkptName)
	// 	}
	// }

	str += "Press 'ctrl-c' or 'q' to exit..."
	return str
	// }

}
