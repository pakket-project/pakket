package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
	"github.com/theckman/yacspin"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

// search
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
		pkgData, _, err := pkg.GetPackageMetadata(packageName) // Get package
		if _, ok := err.(pkg.PackageNotFoundError); ok {
			spinner.StopFailMessage("Cannot find package " + color.CyanString(packageName) + "\n")
			spinner.StopFail()
			return
		} else if err != nil {
			spinner.StopFailMessage("Error while searching package\n" + err.Error())
			spinner.StopFail()
			return
		}

		spinner.StopMessage("Found package " + color.CyanString(packageName) + ":\n")
		spinner.Stop()

		fmt.Printf("Name: %s\n", pkgData.Package.Name)
		fmt.Printf("Description: %s\n", pkgData.Package.Description)
		fmt.Printf("Latest version: %s\n", pkgData.Package.Latest)
		fmt.Printf("Homepage: %s\n", pkgData.Package.Homepage)
	},
}
