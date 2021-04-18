package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
	"github.com/theckman/yacspin"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pull the latest repositories",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := yacspin.Config{
			Frequency:         50 * time.Millisecond,
			HideCursor:        true,
			ColorAll:          false,
			CharSet:           yacspin.CharSets[14],
			Suffix:            " ",
			SuffixAutoColon:   false,
			StopCharacter:     "✓",
			StopFailCharacter: "✗",
			StopColors:        []string{"fgGreen"},
			StopFailColors:    []string{"fgRed"},
			Colors:            []string{"fgCyan"},
		}
		spinner, _ := yacspin.New(cfg)

		spinner.Start()
		spinner.Message("Syncing taps...")

		for i := 0; i < len(config.Config.Repositories); i++ {
			repo := config.Config.Repositories[i]

			if exist := util.DoesPathExist(repo.Path + "/.git"); !exist {
				spinner.StopCharacter("")
				spinner.Stop()
				fmt.Printf("Repository %s is not a git repo, skipping...\n", color.CyanString(repo.Name))
				spinner.Start()
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

			spinner.Message("syncing " + repo.Name + "...")
			err = tree.Pull(&git.PullOptions{RemoteName: "origin"})
			if err.Error() == "already up-to-date" {
			} else if err != nil {
				panic(fmt.Errorf("error while pulling git repo (%s) %s", repo.Name, err))
			}
		}

		spinner.StopCharacter("✓")
		spinner.StopMessage("Successfully synced!")
		spinner.Stop()
	},
}
