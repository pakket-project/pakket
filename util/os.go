package util

import (
	"os"
	"path"
)

var (
	StewPath     = "/usr/local/stew"
	ConfigFile   = "/usr/local/stew/config.toml"
	RepoPath     = "/usr/local/stew/repositories"
	Binpath      = "/usr/local/bin"
	TempRepoPath = "/var/tmp/stew/repo"
)

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
