package pkg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cavaliercoder/grab"
	"github.com/pakket-project/pakket/internals/repo"
	"github.com/pakket-project/pakket/util"
	"github.com/pakket-project/pakket/util/style"
)

func executeScript(script string) (err error) {
	// make sure permissions are set correctly
	err = os.Chmod(script, 0755)
	if err != nil {
		return err
	}

	// run the script
	cmd := exec.Command("bash", "-euxo", "pipefail", script)

	// set the environment variables
	cmd.Env = os.Environ()

	cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout

	return cmd.Run()
}

func HandleScript(name string, pkg PkgData, savePath string) (err error) {
	path := fmt.Sprintf("%s/%s.bash", savePath, name)

	exists, err := downloadScript(name+".bash", pkg, path)
	if err != nil {
		return err
	}

	if exists {
		err := runScript(name, pkg, savePath)

		if err != nil {
			return err
		}
	}

	return nil
}

func downloadScript(name string, pkg PkgData, savePath string) (exists bool, err error) {
	url := fmt.Sprintf("%s/%s/%s/%s", repo.CoreRepositoryURL, pkg.PkgDef.Package.Name, pkg.PkgDef.Package.Version, name)
	resp, err := grab.Get(savePath, url)

	if !(resp.HTTPResponse.StatusCode == 404 || resp.HTTPResponse.StatusCode == 200) {
		return false, err
	}

	return resp.HTTPResponse.StatusCode == 200, nil
}

func runScript(name string, pkg PkgData, savePath string) (err error) {
	url := style.Link.Render(fmt.Sprintf("%s/%s.bash", pkg.RepoURL, name))
	fmt.Printf("\nPackage %s has a %s script: %s \n", pkg.PkgDef.Package.Name, name, url)

	yes := util.Confirm(fmt.Sprintf("Allow package %s to run a %s script?", pkg.PkgDef.Package.Name, name))
	fmt.Print("\n")

	if !yes {
		fmt.Printf("Not running %s script. Please note this may cause errors when using the package.\n", name)
		return nil
	}

	return executeScript(fmt.Sprintf("%s/%s.bash", savePath, name))
}
