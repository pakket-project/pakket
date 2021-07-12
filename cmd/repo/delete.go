package repoCmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/internals/repo"
)

var ()

var DeleteCmd = &cobra.Command{
	Use:     "delete name/repo",
	Short:   "Delete a repository.",
	Example: "stew repo delete stew/core",
	Args:    cobra.MinimumNArgs(1),
	// validate args
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for i := range args {
			if len(strings.Split(args[i], "/")) != 2 {
				return errors.New("arguments must be formatted like <author>/<name>")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var foundRepo bool
		for i := range args {
			splittedArgs := strings.Split(args[i], "/")
			author := splittedArgs[0]
			name := splittedArgs[1]

			for i, v := range config.Config.Repositories.Locations {
				if v.Author == author && v.Name == name {
					foundRepo = true
					err := repo.Delete(i)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		if !foundRepo {
			fmt.Println("No repositories found.")
		}
	},
}
