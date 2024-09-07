package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreateDistCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dist [lat1] [lon1] [lat2] [lon2]",
		Short: "Calculate distance (km) based on GPS coordinates",
		Long: `Calculate distance (km) based on GPS coordinates
Expects 4 float positional arguments (lat1, lon1, lat2, lon2)
representing 2 GPS coordinates`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}

	return cmd
}
