package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
	"github.com/theckman/yacspin"
)

var showDependencies bool

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&showDependencies, "dependencies", "d", false, "show dependencies")
}

// search
var searchCmd = &cobra.Command{
	Use:     "search [package] <version>",
	Short:   "Search for a specific package",
	Aliases: []string{"info"},
	Args:    cobra.MinimumNArgs(1),
	Example: "stew search wget",
	Run: func(cmd *cobra.Command, args []string) {
		spinner, err := yacspin.New(util.SpinnerConf)
		if err != nil {
			panic(err)
		}
		spinner.Start()

		packageName := args[0]

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

		var version string

		// If version argument is supplied
		if len(args) == 2 {
			version = args[1]
		} else {
			version = pkgData.Package.Latest
		}

		versionData, err := pkg.GetPackageVersion(packageName, version)
		if err != nil {
			spinner.StopFailMessage("Cannot find version " + version)
			spinner.StopFail()
			return
		}

		var supportsRosetta string

		// Rosetta support icon
		if versionData.Binaries.SupportsRosetta {
			supportsRosetta = color.GreenString("✓")
		} else {
			supportsRosetta = color.RedString("✗")
		}

		spinner.StopMessage("Found package " + color.Bold(color.CyanString(packageName)) + " (" + color.MagentaString(version) + ")" + ":")
		spinner.Stop()

		// Print package information
		fmt.Printf("Name: %s\n", pkgData.Package.Name)
		fmt.Printf("Description: %s\n", pkgData.Package.Description)
		fmt.Printf("Latest version: %s\n", pkgData.Package.Latest)
		fmt.Printf("Available versions: %s\n", strings.Join(pkgData.Package.AvailableVersions, ", "))
		fmt.Printf("Homepage: %s\n\n", pkgData.Package.Homepage)

		// Binaries
		if len(versionData.Binaries.Intel) > 0 || len(versionData.Binaries.Silicon) > 0 {
			fmt.Println("This package has Intel & Silicon binaries available")
		} else if len(versionData.Binaries.Silicon) > 0 {
			fmt.Println("This package has Silicon binaries available")
		} else {
			fmt.Println("This package has Intel binaries available")
			fmt.Printf("Rosetta support: %s\n", supportsRosetta)
		}

		// Dependencies (if -d flag)
		if showDependencies {
			fmt.Print("\n")
			if len(versionData.Dependencies.Dependencies) > 0 {
				fmt.Printf("Dependencies: %s\n", strings.Join(versionData.Dependencies.Dependencies, ", "))
			}

			if len(versionData.Dependencies.BuildDependencies) > 0 {
				fmt.Printf("Build dependencies: %s\n", strings.Join(versionData.Dependencies.BuildDependencies, ", "))
			}

			if len(versionData.Dependencies.OptionalDependencies) > 0 {
				fmt.Printf("Optional dependencies: %s\n", strings.Join(versionData.Dependencies.OptionalDependencies, ", "))
			}
		}

	},
}
