package cmd

import (
	"fmt"
	"strings"

	"github.com/pakket-project/pakket/internals/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

// search
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all installed packages",
	Aliases: []string{"ls"},
	Example: "pakket list",
	Run: func(cmd *cobra.Command, args []string) {
		var packages []string
		for _, v := range config.LockFile.Packages {
			packages = append(packages, v.Name)
		}

		fmt.Printf("Installed packages:\n\n%s\n", strings.Join(packages, "\n"))
	},
}
