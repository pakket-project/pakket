package cmd

import (
	"github.com/spf13/cobra"
)

type Pakket struct {
	Version string
}

var (
	pakket = Pakket{Version: "v0.0.1"}
	yes    bool
)

var rootCmd = &cobra.Command{
	Use:   "pakket",
	Short: "pakket is a package manager for macOS. Contribute: https://github.com/pakket-project/pakket",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
