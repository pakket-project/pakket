package util

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// TODO: these need to be defined in the config
var (
	// Path that houses config & repositores
	PakketPath = "/etc/pakket"
	// Path to main pakket config
	ConfigFile = path.Join(PakketPath, "config.toml")
	// Path to lockfile w/ currently installed packages
	LockfilePath = path.Join(PakketPath, "lockfile.toml")
	// Path to pakket repositories
	RepoPath = path.Join(PakketPath, "repositories")
	// Path to temporary downloads/clones
	DownloadPath = "/var/tmp/pakket"
	// Prefix
	Prefix = "/usr/local"
)

var (
	// is either "silicon" or "intel"
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
