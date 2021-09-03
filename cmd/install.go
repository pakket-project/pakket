package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
	"github.com/stewproject/stew/util/style"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&yes, "yes", "y", false, "skips all confirmation prompts")
}

var (
	pkgsToInstall []pkg.PkgData
	pkgs          []string
	// errors        []error
	totalSize int64
	yes       bool
)

// repo add
var installCmd = &cobra.Command{
	Use:     "install package[@version]",
	Short:   "Install packages",
	Args:    cobra.MinimumNArgs(1),
	Example: "stew install golang wget python@3.9",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			p := strings.Split(v, "@")
			if len(p) == 1 {
				p = append(p, "latest")
			}
			pkgName := p[0] // package name
			version := p[1] // version

			pkgData, err := pkg.GetPackage(pkgName, version)
			if err != nil {
				panic(err)
			}

			pkgs = append(pkgs, fmt.Sprintf("%s-%s", pkgData.PkgDef.Package.Name, pkgData.Version))
			pkgsToInstall = append(pkgsToInstall, pkgData)
			// totalSize += pkgData.BinSize
		}

		fmt.Printf("Packages: %s (%d)\n", strings.Join(pkgs, ", "), len(pkgs))
		fmt.Printf("Total download size: %s\n", util.ByteToString(totalSize))

		if !yes {
			yes = util.Confirm("\nDo you want to continue?")
		}

		if yes {
			for _, v := range pkgsToInstall {
				err := pkg.InstallPackage(v)
				if err != nil {
					fmt.Printf("\n%s: %s\n", style.Error.Render("Error"), err.Error())
				} else {
					fmt.Printf("\nInstalled %s", v.PkgDef.Package.Name)
				}
			}
		}
	},
}
