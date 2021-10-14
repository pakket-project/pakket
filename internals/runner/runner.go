package runner

import (
	"os"
	"os/exec"

	"github.com/pakket-project/pakket/util"
)

var (
	scriptEnv []string
)

func RunScript(script string, env ...string) (err error) {
	// make sure permissions are set correctly
	err = os.Chmod(script, 0755)
	if err != nil {
		return err
	}

	// run the script
	cmd := exec.Command("bash", "-euxo", "pipefail", script)

	// set the environment variables
	scriptEnv = append(scriptEnv, os.Environ()...)
	scriptEnv = append(scriptEnv, env...)
	//these variables dont exist in the non-builder runner
	// scriptEnv = append(scriptEnv, "PAKKET_PKG_PATH="+util.TmpPkgPath, "PAKKET_SRC_DIR="+util.TmpSrcPath)
	scriptEnv = append(scriptEnv, "PAKKET_ARCH="+util.Arch)
	cmd.Env = scriptEnv

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
