package pkg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cavaliercoder/grab"
	"github.com/pakket-project/pakket/internals/config"
	"github.com/pakket-project/pakket/internals/repo"
	"github.com/pakket-project/pakket/util"
	"github.com/pakket-project/pakket/util/style"
)

func (pkg *PkgData) HandleScript(name string, savePath string, yes bool) (err error) {
	path := fmt.Sprintf("%s/%s.bash", savePath, name)

	exists, err := pkg.downloadScript(name+".bash", path)
	if err != nil {
		return err
	}

	if exists {
		err := pkg.runScript(name, savePath, yes)

		if err != nil {
			return err
		}
	}

	return nil
}

func (pkg *PkgData) downloadScript(name string, savePath string) (exists bool, err error) {
	url := fmt.Sprintf("%s/%s/%s/%s", repo.CoreRepositoryURL, pkg.PkgDef.Package.Name, pkg.PkgDef.Package.Version, name)
	resp, err := grab.Get(savePath, url)

	if !(resp.HTTPResponse.StatusCode == 404 || resp.HTTPResponse.StatusCode == 200) {
		return false, err
	}

	return resp.HTTPResponse.StatusCode == 200, nil
}

func (pkg *PkgData) runScript(name string, savePath string, yes bool) (err error) {
	url := style.Link.Render(fmt.Sprintf("%s/%s.bash", pkg.RepoURL, name))
	fmt.Printf("\nPackage %s has a %s script: %s \n", pkg.PkgDef.Package.Name, name, url)

	if !yes {
		yes = util.Confirm(fmt.Sprintf("Allow package %s to run a %s script?", pkg.PkgDef.Package.Name, name))
		fmt.Print("\n")
	}

	if !yes {
		fmt.Printf("Not running %s script. Please note this may cause errors when using the package.\n", name)
		return nil
	}

	script := fmt.Sprintf("%s/%s.bash", savePath, name)

	// make sure permissions are set correctly
	err = os.Chmod(script, 0755)
	if err != nil {
		return err
	}

	// run the script
	cmd := exec.Command("bash", "-euxo", "pipefail", script)

	// set the environment variables
	var scriptEnv []string

	scriptEnv = append(scriptEnv, os.Environ()...)
	scriptEnv = append(scriptEnv, "PAKKET_PREFIX="+config.C.Paths.Prefix)
	cmd.Env = scriptEnv

	cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout

	return cmd.Run()
}
