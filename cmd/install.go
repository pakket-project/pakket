package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
)

func init() {
	rootCmd.AddCommand(installCmd)
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
			totalSize += pkgData.BinSize
		}

		fmt.Printf("Packages: %s (%d)\n", strings.Join(pkgs, ", "), len(pkgs))
		fmt.Printf("Total download size: %s\n", util.ByteToString(totalSize))
		if confirm := util.Confirm("\nDo you want to continue?"); confirm {
			for _, v := range pkgsToInstall {
				err := pkg.InstallPackage(v.PkgDef, v.BinData)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Printf("Installed %s", v.PkgDef.Package.Name)
				}
			}
		}
	},
}
