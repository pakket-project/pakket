package cmd

import (
	"github.com/spf13/cobra"
	repoCmd "github.com/stewproject/stew/cmd/repo"
)

func init() {
	rootCmd.AddCommand(repoRoot)
	repoRoot.AddCommand(repoCmd.AddCmd)    // add
	repoRoot.AddCommand(repoCmd.ListCmd)   // list
	repoRoot.AddCommand(repoCmd.DeleteCmd) // delete
	repoRoot.AddCommand(repoCmd.SyncCmd)   // sync
}

var repoRoot = &cobra.Command{
	Use: "repo",
}
