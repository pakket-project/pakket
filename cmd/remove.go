package cmd

import (
	"fmt"
	"strings"

	"github.com/pakket-project/pakket/internals/pkg"
	"github.com/pakket-project/pakket/util"
	"github.com/pakket-project/pakket/util/style"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&yes, "yes", "y", false, "skips all confirmation prompts")
}

var (
	pkgsToRemove []string
)

// remove package
var removeCmd = &cobra.Command{
	Use:     "remove package",
	Short:   "Remove packages",
	Args:    cobra.MinimumNArgs(1),
	Example: "pakket remove golang wget python",
	Run: func(cmd *cobra.Command, args []string) {
		keys := make(map[string]bool)
		for _, v := range args {
			name := v

			// check for dulicates, skip if duplicate
			if _, value := keys[name]; value {
				continue
			}
			keys[name] = true

			// check if package is already installed (lockfile)
			// config.LockFile.Packages

			pkgsToRemove = append(pkgsToRemove, name)

			// totalSize += pkgData.BinSize
		}

		fmt.Printf("Packages to remove: %s (%d)\n", strings.Join(pkgs, ", "), len(pkgs))
		fmt.Printf("Total size removing: %s\n", util.ByteToString(totalSize))

		if !yes {
			yes = util.Confirm("\nDo you want to continue?")
		}

		if yes {
			for _, v := range pkgsToRemove {
				err := pkg.RemovePackage(v)

				if err != nil {
					fmt.Printf("\n%s: %s\n", style.Error.Render("Error"), err.Error())
				} else {
					fmt.Printf("removed %s", v)
				}
			}
		}
	},
}
