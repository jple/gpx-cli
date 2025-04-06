package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	. "github.com/jple/gpx-cli/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateAddNameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add_name [name] [lat] [lon]",
		Short: "Add name to closest trkpts",
		Long:  `Add name to closest trkpts`,
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			lat, _ := strconv.ParseFloat(args[1], 64)
			lon, _ := strconv.ParseFloat(args[2], 64)

			gpx := Gpx{Filepath: viper.GetString("filename")}
			gpx.ParseFile(gpx.Filepath)

			closest := gpx.GetClosestTrkpts(Pos{Lat: lat, Lon: lon})

			// Add name to matches trkpts
			for i, _ := range closest {
				// Confirmation if existing name
				if n := closest[i].Name; n != nil {
					scanner := bufio.NewScanner(os.Stdin)
					fmt.Printf("Overwrite existing name (%v) ?[y/n] ", *n)
					for scanner.Scan() {
						yn := scanner.Text()
						if yn == "y" {
							fmt.Println("Replacing name to", name)
							closest[i].AddName(name)
							break
						} else if scanner.Text() == "n" {
							break
						}
					}
				} else {
					closest[i].AddName(name)
				}
			}

			if viper.GetBool("inplace") {
				gpx.Save(viper.GetString("filename"))
			} else {
				gpx.Save(viper.GetString("output"))
			}
		},
	}

	return cmd
}
