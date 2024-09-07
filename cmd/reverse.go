package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateReverseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reverse",
		Short: "Reverse whole GPX elements",
		Long:  "Reverse all trk and trkwpt in GPX file",
	}

	cmd.Flags().StringP("trk_id", "t", "all", "Trk Id to reverse. Defaults to all")
	viper.BindPFlag("trk_id", cmd.Flags().Lookup("trk_id"))

	return cmd
}
