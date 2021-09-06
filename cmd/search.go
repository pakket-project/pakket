package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/errors"
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

		var version *string
		if len(args) > 1 {
			version = &args[1]
		} else {
			version = nil
		}

		spinner.Message(fmt.Sprintf("Searching for package %s", style.Pkg.Render(packageName)))

		pkgData, err := pkg.GetPackage(packageName, version) // Get package
		if _, ok := err.(errors.PackageNotFoundError); ok {
			spinner.StopFailMessage(fmt.Sprintf("Cannot find package %s\n", style.Pkg.Render(packageName)))
			spinner.StopFail()
			return
		} else if err != nil {
			spinner.StopFailMessage(fmt.Sprintf("Error while searching package: %s", err.Error()))
			spinner.StopFail()
			return
		}

		//Rosetta support icon
		var supportsRosetta string
		if pkgData.VerData.SupportsRosetta {
			supportsRosetta = style.Success.Render("✓")
		} else {
			supportsRosetta = style.Error.Render("✗")
		}

		spinner.StopMessage(fmt.Sprintf("Found package %s:\n", style.Pkg.Render(packageName)))
		spinner.Stop()

		// Print package information
		fmt.Printf("Description: %s\n", pkgData.PkgDef.Package.Description)
		fmt.Printf("Latest version: %s\n", pkgData.PkgDef.Package.Version)
		fmt.Printf("Available versions: %s\n", strings.Join(pkgData.PkgDef.Package.AvailableVersions, ", "))
		fmt.Printf("Homepage: %s\n\n", style.Link.Render(pkgData.PkgDef.Package.Homepage))

		intelPackage := style.Success.Render("✓")
		siliconPackage := style.Success.Render("✓")

		// Packages
		if pkgData.VerData.Intel.Hash == "" {
			// no intel package
			intelPackage = style.Error.Render("✗")
		}
		fmt.Printf("Intel: %s\n", intelPackage)

		if pkgData.VerData.Silicon.Hash == "" {
			// no silicon package
			siliconPackage = style.Error.Render("✗")
			fmt.Printf("Apple Silicon: %s\n", siliconPackage)
			fmt.Printf("Rosetta: %s\n", supportsRosetta)
		} else {
			fmt.Printf("Apple Silicon: %s\n", siliconPackage)
		}

		// Dependencies (if -d flag)
		if showDependencies {
			fmt.Print("\n")
			if len(pkgData.VerData.Dependencies.Dependencies) > 0 {
				fmt.Printf("Dependencies: %s\n", strings.Join(pkgData.VerData.Dependencies.Dependencies, ", "))
			}

			if len(pkgData.VerData.Dependencies.BuildDependencies) > 0 {
				fmt.Printf("Build dependencies: %s\n", strings.Join(pkgData.VerData.Dependencies.BuildDependencies, ", "))
			}

			if len(pkgData.VerData.Dependencies.OptionalDependencies) > 0 {
				fmt.Printf("Optional dependencies: %s\n", strings.Join(pkgData.VerData.Dependencies.OptionalDependencies, ", "))
			}
		}

	},
}
