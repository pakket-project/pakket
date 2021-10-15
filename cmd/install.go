package cmd

import (
	"fmt"
	"strings"

	"github.com/pakket-project/pakket/internals/config"
	"github.com/pakket-project/pakket/internals/pkg"
	"github.com/pakket-project/pakket/util"
	"github.com/pakket-project/pakket/util/style"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&yes, "yes", "y", false, "skips all confirmation prompts")
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "force")
}

var (
	pkgsToInstall []pkg.PkgData
	pkgs          []string
	// errors        []error
	totalSize int64
)

// repo add
var installCmd = &cobra.Command{
	Use:     "install package[@version]",
	Short:   "Install packages",
	Args:    cobra.MinimumNArgs(1),
	Example: "pakket install golang wget python@3.9",
	Run: func(cmd *cobra.Command, args []string) {
		keys := make(map[string]bool)
		for _, v := range args {
			p := strings.Split(v, "@")
			name := p[0]

			// check for dulicates, skip if duplicate
			if _, value := keys[name]; value {
				continue
			}
			keys[name] = true

			// check if package is already installed (lockfile)
			if v, ok := config.LockFile.Packages[name]; ok {
				fmt.Printf("%s is already installed\n", v.Name)
				continue
			}

			var version *string
			if len(p) > 1 {
				version = &p[1]
			} else {
				version = nil
			}

			pkgData, err := pkg.GetPackage(name, version)
			if err != nil {
				panic(err)
			}

			pkgs = append(pkgs, fmt.Sprintf("%s-%s", pkgData.PkgDef.Package.Name, pkgData.Version))
			pkgsToInstall = append(pkgsToInstall, *pkgData)
			totalSize += pkgData.BinSize
		}

		if len(pkgs) > 0 {

			fmt.Printf("Packages: %s (%d)\n", strings.Join(pkgs, ", "), len(pkgs))
			fmt.Printf("Total download size: %s\n", util.ByteToString(totalSize))

			if !yes {
				yes = util.Confirm("\nDo you want to continue?")
			}

			if yes {
				for _, v := range pkgsToInstall {
					err := pkg.InstallPackage(v, force, yes)
					if err != nil {
						fmt.Printf("\n%s: %s\n", style.Error.Render("Error"), err.Error())
					} else {
						fmt.Printf("Installed %s\n", v.PkgDef.Package.Name)
					}
				}
			}
		}
	},
}
