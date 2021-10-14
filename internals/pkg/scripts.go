package pkg

import (
	"fmt"

	"github.com/cavaliercoder/grab"
	"github.com/pakket-project/pakket/internals/repo"
	"github.com/pakket-project/pakket/internals/runner"
	"github.com/pakket-project/pakket/util"
	"github.com/pakket-project/pakket/util/style"
)

func HandleScript(name string, pkg PkgData, savePath string) (err error) {
	path := fmt.Sprintf("%s/%s.bash", savePath, name)

	exists, err := DownloadScript(name+".bash", pkg, path)
	if err != nil {
		return err
	}

	if exists {
		err := RunScript(name, pkg, savePath)

		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadScript(name string, pkg PkgData, savePath string) (exists bool, err error) {
	url := fmt.Sprintf("%s/%s/%s/%s", repo.CoreRepositoryURL, pkg.PkgDef.Package.Name, pkg.PkgDef.Package.Version, name)
	resp, err := grab.Get(savePath, url)

	if !(resp.HTTPResponse.StatusCode == 404 || resp.HTTPResponse.StatusCode == 200) {
		return false, err
	}

	return resp.HTTPResponse.StatusCode == 200, nil
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
