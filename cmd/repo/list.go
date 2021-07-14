package repoCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util/style"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all repositories.",
	Example: "stew repo list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Config.Repositories.Locations) > 0 {
			for _, v := range config.Config.Repositories.Locations {
				fmt.Printf("%s/%s (%s)\n", style.Repo.Render(v.Author), style.Repo.Render(v.Name), style.Link.Render(v.GitURL))
			}

			fmt.Printf("\nTotal of %v repositories found.\n", len(config.Config.Repositories.Locations))
		} else {
			fmt.Println("No repositories found. Add the core repository: \"stew repo add https://github.com/stewproject/packages\"")
		}
	},
}
