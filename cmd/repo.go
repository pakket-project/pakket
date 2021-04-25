package cmd

import (
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

		spinner.Start()
		spinner.Message("Adding repositories...")

		var addedRepos []string
		for i := range args {
			repoLink := args[i]

			// Fetch repo, add to config
			metadata, err := repo.AddRepo(repoLink)

			// If repo already exists, but is not defined in the config
			if _, ok := err.(repo.UndefinedRepositoryAlreadyExistsError); ok {
				util.PrintSpinnerMsg(spinner, "Error while adding repository: "+style.Error.Render("Repository "+metadata.Repository.Name+"already exists, but is not defined in the config, so we're adding it."))
				config.AddRepo(config.RepositoriesMetadata{
					Name:         metadata.Repository.Name,
					Path:         path.Join(util.RepoPath, metadata.Repository.Name),
					PackagesPath: metadata.Repository.PackagesPath,
					GitURL:       repoLink,
					IsGit:        true,
				})
				continue
			} else if err != nil {
				util.PrintSpinnerMsg(spinner, "Error while adding repository: "+style.Error.Render(err.Error()))
				continue
			}

			addedRepos = append(addedRepos, metadata.Repository.Name)
		}

		spinner.StopMessage("Successfully added the repositories " + style.Pkg.Render(strings.Join(addedRepos, ", ")) + ".")
		spinner.Stop()
	},
}
