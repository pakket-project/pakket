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
	errors []error
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
			pkgName := p[0]
			ver := p[1]
			pkgDef, pkgPath, err := pkg.GetPackageMetadata(pkgName)
			if err != nil {
				panic(err)
			}
			vers, err := pkg.GetPackageVersion(pkgName, *pkgPath, pkgDef.Package.Version)
			if err != nil {
				panic(err)
			}
			binary := pkg.GetBinaryMetadata(*vers)
			size, err := pkg.GetPackageSize(binary)
			if err != nil {
				panic(err)
			}
			if confirm := util.Confirm(fmt.Sprintf("Continue? %s", util.ByteToString(size))); confirm {
				err := pkg.InstallPackage(pkgName, ver)
				if _, ok := err.(pkg.PackageNotFoundError); ok {
					errors = append(errors, err)
				}
			}
		}
		if len(errors) > 0 {
			fmt.Println(errors)
		}

		fmt.Printf("installed %s\n", strings.Join(args, ", "))
	},
}
