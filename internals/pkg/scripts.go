package pkg

import (
	"fmt"

	"github.com/cavaliercoder/grab"
	"github.com/pakket-project/pakket/internals/repo"
	"github.com/pakket-project/pakket/internals/runner"
	"github.com/pakket-project/pakket/util"
	"github.com/pakket-project/pakket/util/style"
)

func DownloadScript(name string, pkg PkgData, savePath string) (err error) {
	url := fmt.Sprintf("%s/%s/%s/%s", repo.CoreRepositoryURL, pkg.PkgDef.Package.Name, pkg.PkgDef.Package.Version, name)
	_, err = grab.Get(savePath, url)

	return err
}

func RunScript(name string, pkg PkgData, savePath string) (err error) {
	url := style.Link.Render(fmt.Sprintf("%s/%s.bash", pkg.PkgRepoUrl, name))
	fmt.Printf("\nPackage %s has a %s script: %s \n", pkg.PkgDef.Package.Name, name, url)

	yes := util.Confirm(fmt.Sprintf("Allow package %s to run a %s script?", pkg.PkgDef.Package.Name, name))
	fmt.Print("\n")

	if !yes {
		fmt.Printf("Not running %s script. Please note this may cause errors when using the package.\n", name)
		return nil
	}

	return runner.RunScript(fmt.Sprintf("%s/%s.bash", savePath, name))
}
