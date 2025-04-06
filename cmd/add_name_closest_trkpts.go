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
		Use:   "add-name [name] [lat] [lon]",
		Short: "Add name to closest trkpts",
		Long:  `Add name to closest trkpts`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				fmt.Printf("Expected 3 args, got %v\n", len(args))
				for _, a := range args {
					fmt.Println(a)
				}
			}
			name := args[0]
			lat, _ := strconv.ParseFloat(args[1], 64)
			lon, _ := strconv.ParseFloat(args[2], 64)

			gpx := Gpx{}
			gpx.ParseFile(viper.GetString("filename"))

			closest := gpx.GetClosestTrkpts(Pos{Lat: lat, Lon: lon})

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
