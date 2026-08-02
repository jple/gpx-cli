package cmd

import (
	"fmt"
	"os"

	. "github.com/jple/gpx-cli/internal/gpx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: merge with GetInfoCmd
func CreateHTMLSummaryCmd() *cobra.Command {
	var speed FloatValue = 4.5
	flagsConf := []FlagConfig{
		{
			Name: "speed", Shortname: "s", DefaultValue: &speed,
			Description: "Hiking speed on flat (km/h)",
		},
	}

	cmd := &cobra.Command{
		Use:   "html-summary",
		Short: "List all trk names",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd, flagsConf)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gpx := Gpx{}
			content := gpx.Parse(viper.GetString("filename")).
				Summarize(viper.GetFloat64("speed")).ToHTML()

			filepath := viper.GetString("output")
			if filepath == "out.gpx" { // default value
				filepath = "gpx_summary.html"
			}
			// Write content to file
			file, err := os.Create(filepath)
			if err != nil {
				fmt.Println("Error creating XML file:", err)
				return
			}

			_, err = file.WriteString(content)
			if err != nil {
				fmt.Println("Error writing to HTML file:", err)
				return
			}

		},
	}

	initFlags(cmd, flagsConf)

	return cmd
}
