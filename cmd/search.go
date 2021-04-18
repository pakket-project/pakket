package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/packages"
	"github.com/theckman/yacspin"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search formula",
	Short: "Search for a specific formula.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you must provide a formula to search for")
		}
		return nil
	},
	Example: "stew search wget",
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]
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
		spinner, err := yacspin.New(cfg)
		if err != nil {
			panic(err)
		}

		spinner.Start()

		spinner.Message("Searching for formula " + color.CyanString(packageName))
		def := packages.GetPackageData(packageName)
		if err != nil {
			spinner.StopFailMessage("Cannot find formula " + color.CyanString(packageName))
			spinner.StopFail()
			os.Exit(1)
		}
		spinner.StopMessage("Found formula " + color.CyanString(packageName) + ":\n")
		spinner.Stop()

		fmt.Printf("Name: %s\n", def.Package.Name)
		fmt.Printf("Description: %s\n", def.Package.Description)
		fmt.Printf("Version: %s\n", "TBD")
		fmt.Printf("Homepage: %s\n", def.Package.Homepage)
	},
}
