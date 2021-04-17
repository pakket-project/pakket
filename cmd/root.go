package cmd

import (
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run stew help for usage")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
