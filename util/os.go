package util

import (
	"os"
)

var (
	// Path that houses config & repositores
	StewPath = "/usr/local/stew"
	// Path to main Stew config
	ConfigFile = "/usr/local/stew/config.toml"
	// Path to Stew repositories
	RepoPath = "/usr/local/stew/repositories"
	// Path to binary installation path
	Binpath = "/usr/local/bin"
	// Path to temporary clone repositories
	TempRepoPath = "/var/tmp/stew/repo"
)

// Check if path exists
func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	exists := os.IsNotExist(err)

	return !exists
}
