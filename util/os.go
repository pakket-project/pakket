package util

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	// either "intel" or "silicon"
	Arch string
)

// Check if path exists
func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	exists := os.IsNotExist(err)

	return !exists
}

// Get macOS version.
// TODO: find a better way to do this
func GetVersion() (ver string) {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		panic(err)
	}
	ver = string(out)

	if strings.HasPrefix(ver, "11") {
		return "big_sur"
	} else if strings.HasPrefix(ver, "10.15") {
		return "catalina"
	} else if strings.HasPrefix(ver, "10.14") {
		return "mojave"
	} else {
		panic("unsupported macOS version")
	}
}

// Check if directory is empty
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
