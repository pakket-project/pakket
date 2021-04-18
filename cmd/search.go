package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
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

		spinner, err := yacspin.New(util.SpinnerConf)
		if err != nil {
			panic(err)
		}

		spinner.Start()

		spinner.Message("Searching for formula " + color.CyanString(packageName))
		def := pkg.GetPackageData(packageName)
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
