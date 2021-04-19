package cmd

import (
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/internals/repo"
	"github.com/stewproject/stew/util"
	"github.com/theckman/yacspin"
)

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(addCmd)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Main repository management command. See subcommands for details",
}

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

			err := repo.AddRepo(repoLink)
			if err != nil {
				spinner.StopFailMessage("error while adding repository")
				spinner.StopFail()
			}

			var metadata *repo.Metadata
			for b := range config.Config.Repositories.Locations {
				if config.Config.Repositories.Locations[b].GitURL == repoLink {
					metadata = repo.GetMetadataFromRepo(config.Config.Repositories.Locations[b].Name)
					break
				}
			}

			addedRepos = append(addedRepos, metadata.Repository.Name)
		}
		spinner.StopMessage("Successfully added the repositories " + color.CyanString(strings.Join(addedRepos, ", ")))
		spinner.Stop()
	},
}
