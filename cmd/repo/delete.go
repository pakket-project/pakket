package repoCmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/internals/repo"
	"github.com/stewproject/stew/util"
	"github.com/stewproject/stew/util/style"
)

var DeleteCmd = &cobra.Command{
	Use:     "delete name/repo",
	Short:   "Delete a repository",
	Example: "stew repo delete stew/core",
	Args:    cobra.ExactArgs(1),
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
		splittedArgs := strings.Split(args[0], "/")
		author := splittedArgs[0]
		name := splittedArgs[1]

		for i, v := range config.Config.Repositories.Locations {
			if v.Author == author && v.Name == name {
				foundRepo = true

				if confirm := util.DestructiveConfirm(fmt.Sprintf("Do you really want to delete the repository %s/%s?", style.Repo.Render(v.Author), style.Repo.Render(v.Name))); confirm {
					fmt.Printf("Deleted %s/%s.\n", v.Author, v.Name)
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
