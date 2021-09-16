package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of pakket",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pakket - %s\n", pakket.Version)
	},
}
