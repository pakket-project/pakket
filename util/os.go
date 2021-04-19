package util

import (
	"os"
	"path"
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

// First removes all files in the directory and then remove the directory itself.
func RemoveFolder(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(path.Join(dir, name))
		if err != nil {
			return err
		}

	}
	err = os.Remove(dir)
	if err != nil {
		return err
	}

	return nil
}
