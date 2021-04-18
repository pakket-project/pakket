package util

import (
	"os"
)

var (
	StewPath   = "/usr/local/stew"
	ConfigFile = "/usr/local/stew/config.toml"
	RepoPath   = "/usr/local/stew/repositories"
	Binpath    = "/usr/local/bin"
)

func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	exists := os.IsNotExist(err)

	return !exists
}
