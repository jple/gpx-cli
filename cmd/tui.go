package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	tui "github.com/jple/gpx-cli/tui"

	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateTuiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Visualize GPX in TUI mode",
		Run: func(cmd *cobra.Command, args []string) {

			var gpx Gpx
			if viper.GetString("filename") == "" {
				fmt.Println("No GPX file loaded")
				os.Exit(1)
			} else {
				gpx.ParseFile(viper.GetString("filename"))
				gpx.SetVitesse(4.5)
			}

			var m tui.GpxTui = tui.GpxTui{
				GpxSummary: gpx.GetInfo(true),
				Gpx:        gpx,
			}
			p := tea.NewProgram(m)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}

		},
	}

	return cmd
}
