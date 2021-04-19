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
	Short: "Print the version number of Stew",
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Stew Package Manager - %s\n", stew.Version)
	},
}
