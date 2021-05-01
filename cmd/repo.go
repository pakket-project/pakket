package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/internals/repo"
	"github.com/stewproject/stew/util"
	"github.com/stewproject/stew/util/style"
	"github.com/theckman/yacspin"
)

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(addCmd)
}

// repo
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Main repository management command. See subcommands for details",
}

// repo add
var addCmd = &cobra.Command{
	Use:     "add git-links",
	Short:   "Add repositories. Supports multiple links",
	Args:    cobra.MinimumNArgs(1),
	Example: "stew repo add https://github.com/stewproject/packages https://github.com/stingalleman/stew-repository",
	Run: func(cmd *cobra.Command, args []string) {
		spinner, _ := yacspin.New(util.SpinnerConf)

		_, err := os.Stat(util.RepoPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Repo path %s was not found: creating it\n ", util.RepoPath)
				err := os.Mkdir(util.RepoPath, fs.ModeDir)
				if err != nil {
					fmt.Printf("failed to create %s: %s\n", util.RepoPath, style.Error.Render(err.Error()))
					os.Exit(1)
				}
			} else {
				// other stat errors
				fmt.Printf("failed to stat %s: %s\n", util.RepoPath, style.Error.Render(err.Error()))
				os.Exit(1)
			}
		}

		spinner.Start()
		spinner.Message("Adding repositories...")

		var addedRepos []string
		for i := range args {
			repoLink := args[i]

			// Fetch repo, add to config
			metadata, err := repo.AddRepo(repoLink)

			// If repo already exists, but is not defined in the config
			if _, ok := err.(repo.UndefinedRepositoryAlreadyExistsError); ok {
				util.PrintSpinnerMsg(spinner, fmt.Sprintf("Error while adding repository: %s", style.Error.Render(fmt.Sprintf("Repository %s already exists, but is not defined in the config, so Stew is adding it.", metadata.Repository.Name))))

				config.AddRepo(config.RepositoriesMetadata{
					Name:         metadata.Repository.Name,
					Path:         path.Join(util.RepoPath, metadata.Repository.Name),
					PackagesPath: metadata.Repository.PackagesPath,
					GitURL:       repoLink,
					IsGit:        true,
				})
			} else if err != nil {
				util.PrintSpinnerMsg(spinner, fmt.Sprintf("Error while adding repository: %s", style.Error.Render(err.Error())))
				continue
			}
			addedRepos = append(addedRepos, metadata.Repository.Name)
		}

		if len(addedRepos) <= 0 {
			spinner.StopFailMessage("No new repositories added.")
			spinner.StopFail()
		} else if len(addedRepos) == 1 {
			spinner.StopMessage(fmt.Sprintf("Successfully added the repository %s.", style.Pkg.Render(strings.Join(addedRepos, ", "))))
		} else {
			spinner.StopMessage(fmt.Sprintf("Successfully added the repositories %s.", style.Pkg.Render(strings.Join(addedRepos, ", ")))) // TODO: make , white
		}
		spinner.Stop()
	},
}
