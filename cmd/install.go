package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/pkg"
	"github.com/stewproject/stew/util"
	"github.com/theckman/yacspin"
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
		spinner, _ := yacspin.New(util.SpinnerConf)
		spinner.Start()

		for _, v := range args {
			p := strings.Split(v, "@")
			if len(p) == 1 {
				p = append(p, "latest")
			}
			pkgName := p[0]
			ver := p[1]

			spinner.Message(fmt.Sprintf("Installing %s (%s)", pkgName, ver))
			err := pkg.InstallPackage(pkgName, ver)
			if _, ok := err.(pkg.PackageNotFoundError); ok {
				errors = append(errors, err)
			}
		}
		fmt.Println(errors)
		spinner.StopMessage(fmt.Sprintf("Succesfully installed %s", strings.Join(args, ", ")))
		spinner.Stop()
	},
}
