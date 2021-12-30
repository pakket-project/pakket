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
		amount := len(config.Lockfile.Packages)

		if amount == 0 {
			fmt.Println("No packages installed.")
			return
		}

		for _, v := range config.Lockfile.Packages {
			packages = append(packages, v.Name)
		}

		fmt.Printf("Installed packages (%d):\n\n* %s\n", amount, strings.Join(packages, "\n* "))
	},
}
