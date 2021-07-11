package repoCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all repositories.",
	Example: "stew repo list",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range config.Config.Repositories.Locations {
			fmt.Printf("%s\n", v.Name)
		}

		fmt.Printf("\nTotal of %v repositories found.\n", len(config.Config.Repositories.Locations))
	},
}
