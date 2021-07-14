package repoCmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
	"github.com/theckman/yacspin"
)

var SyncCmd = &cobra.Command{
	Use:   "repo sync",
	Short: "Pull the latest repositories",
	Run: func(cmd *cobra.Command, args []string) {
		numRepos := len(config.Config.Repositories.Locations)
		if numRepos == 0 {
			fmt.Println("You haven't added any repositories yet. Use stew repository add <repository> to add a repository first.")
			os.Exit(1)
		}
		spinner, _ := yacspin.New(util.SpinnerConf)

		spinner.Start()
		spinner.Message("Syncing taps...")

		for i := 0; i < numRepos; i++ {
			repo := config.Config.Repositories.Locations[i]

			if !config.Config.Repositories.Locations[i].IsGit {
				// not a git repo, not syncing
				continue
			}

			r, err := git.PlainOpen(repo.Path)
			if err != nil {
				panic(fmt.Errorf("error while opening git repo (%s) %s", repo.Name, err))
			}
			tree, err := r.Worktree()
			if err != nil {
				panic(fmt.Errorf("error while opening git repo (%s) %s", repo.Name, err))
			}

			spinner.Message(fmt.Sprintf("syncing %s...", repo.Name))
			err = tree.Pull(&git.PullOptions{})
			if err == git.NoErrAlreadyUpToDate {
				// do nothing
			} else if err != nil {
				panic(fmt.Errorf("error while pulling git repo (%s) %s", repo.Name, err))
			}
		}

		spinner.StopCharacter("✓")
		spinner.StopMessage("Successfully synced!")
		spinner.Stop()
	},
}
