package cmd

import (
	"github.com/spf13/cobra"
)

type Stew struct {
	Version string
}

var (
	stew = Stew{Version: "v0.0.1"}
)

var rootCmd = &cobra.Command{
	Use:   "stew",
	Short: "Stew is a package manager for macOS. Contribute: https://github.com/stewproject/stew",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
