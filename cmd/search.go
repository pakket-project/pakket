package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
	"github.com/stewproject/stew/util/style"
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

		spinner.Message(fmt.Sprintf("Searching for package %s", style.Pkg.Render(packageName)))

		pkgData, pkgPath, repo, err := pkg.GetPackageMetadata(packageName) // Get package
		if _, ok := err.(pkg.PackageNotFoundError); ok {
			spinner.StopFailMessage(fmt.Sprintf("Cannot find package %s\n", style.Pkg.Render(packageName)))
			spinner.StopFail()
			return
		} else if err != nil {
			spinner.StopFailMessage(fmt.Sprintf("Error while searching package: %s", err.Error()))
			spinner.StopFail()
			return
		}

		var version string

		// If version argument is supplied
		if len(args) == 2 {
			version = args[1]
		} else {
			version = pkgData.Package.Version
		}

		versionData, err := pkg.GetPackageVersion(packageName, *pkgPath, version)
		if err != nil {
			spinner.StopFailMessage(fmt.Sprintf("Cannot find version %s", version))
			spinner.StopFail()
			return
		}

		// TODO
		// Rosetta support icon
		// if versionData.Binaries.SupportsRosetta {
		// 	supportsRosetta = style.Success.Render("✓")
		// } else {
		// 	supportsRosetta = style.Error.Render("✗")
		// }

		spinner.StopMessage(fmt.Sprintf("Found package %s:\n", style.Pkg.Render(packageName)))
		spinner.Stop()

		// Print package information
		fmt.Printf("Description: %s\n", pkgData.Package.Description)
		fmt.Printf("Latest version: %s\n", pkgData.Package.Version)
		fmt.Printf("Available versions: %s\n", strings.Join(pkgData.Package.AvailableVersions, ", "))
		fmt.Printf("Repository: %s\n", style.Repo.Render(repo))
		fmt.Printf("Homepage: %s\n\n", style.Link.Render(pkgData.Package.Homepage))

		// TODO
		// Binaries
		// if len(versionData.Binaries.Intel) > 0 || len(versionData.Binaries.Silicon) > 0 {
		// 	fmt.Println("This package has Intel & Silicon binaries available")
		// } else if versionData.Silicon {
		// 	fmt.Println("This package has Silicon binaries available")
		// } else {
		// 	fmt.Println("This package has Intel binaries available")
		// 	fmt.Printf("Rosetta support: %s\n", supportsRosetta)
		// }

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
