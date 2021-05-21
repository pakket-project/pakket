package util

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	// Path that houses config & repositores
	StewPath = "/etc/stew"
	// Path to main Stew config
	ConfigFile = path.Join(StewPath, "config.toml")
	// Path to Stew repositories
	RepoPath = path.Join(StewPath, "repositories")
	// Path to temporary downloads/clones
	TempRepoPath = path.Join(StewPath, "downloads")
)

// Check if path exists
func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	exists := os.IsNotExist(err)

	return !exists
}

// Get macOS version.
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
