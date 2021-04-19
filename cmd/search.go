package cmd

import (
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
	Use:     "search package",
	Short:   "Search for a specific package",
	Args:    cobra.MinimumNArgs(1),
	Example: "stew search wget",
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]

		spinner, err := yacspin.New(util.SpinnerConf)
		if err != nil {
			panic(err)
		}

		spinner.Start()

		spinner.Message("Searching for package " + color.CyanString(packageName))
		pkg := pkg.GetPackageData(packageName) // Get package
		if err != nil {
			spinner.StopFailMessage("Cannot find package " + color.CyanString(packageName))
			spinner.StopFail()
			os.Exit(1)
		}
		spinner.StopMessage("Found package " + color.CyanString(packageName) + ":\n")
		spinner.Stop()

		fmt.Printf("Name: %s\n", pkg.Package.Name)
		fmt.Printf("Description: %s\n", pkg.Package.Description)
		fmt.Printf("Version: %s\n", pkg.Package.Version)
		fmt.Printf("Homepage: %s\n", pkg.Package.Homepage)
	},
}
