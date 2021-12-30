package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var (
	LockFile LockfileStruct
)

type LockfileMetadata struct {
	Name       string `toml:"name"`
	Version    string `toml:"version"`
	Checksum   string `toml:"checksum"`
	Repository string `toml:"repository"`
}

type LockfileStruct struct {
	Packages map[string]LockfileMetadata `toml:"packages" mapstructure:"packages"`
}

func GetLockfile() (err error) {
	file, err := os.ReadFile(LockfilePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(LockfilePath)
		if err != nil {
			return err
		}
		file, err = os.ReadFile(LockfilePath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = toml.Unmarshal(file, &LockFile)

	return err
}

// Add package information to lockfile
func AddPkgToLockfile(metadata LockfileMetadata, files []string) (err error) {
	if len(LockFile.Packages) == 0 {
		LockFile.Packages = make(map[string]LockfileMetadata)
	}

	LockFile.Packages[metadata.Name] = metadata

	newLockfile, err := toml.Marshal(&LockFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(LockfilePath, newLockfile, 0666)
	if err != nil {
		return fmt.Errorf("failed to write lockfile: %s", err)
	}

	err = os.MkdirAll(path.Join(LocalPath, fmt.Sprintf("%s-%s", metadata.Name, metadata.Version)), 0770)
	if err != nil {
		return err
	}

	filesFile, err := os.Create(path.Join(LocalPath, fmt.Sprintf("%s-%s", metadata.Name, metadata.Version), "files"))
	if err != nil {
		return err
	}

	_, err = filesFile.WriteString(strings.Join(files, "\n"))
	if err != nil {
		return err
	}

	return err
}

// Remove package information from lockfile
func RemovePkgFromLockfile(name string) (lockfileData LockfileMetadata, files []string, err error) {
	if _, ok := LockFile.Packages[name]; !ok {
		return LockfileMetadata{}, nil, errors.New("package not found")
	}
	lockfile := LockFile.Packages[name]
	delete(LockFile.Packages, name)

	newLockfile, err := toml.Marshal(&LockFile)
	if err != nil {
		return LockfileMetadata{}, nil, err
	}

	err = os.WriteFile(LockfilePath, newLockfile, 0666)
	if err != nil {
		return LockfileMetadata{}, nil, err
	}

	// delete files from local db
	fileBytes, err := os.ReadFile(path.Join(LocalPath, fmt.Sprintf("%s-%s", lockfile.Name, lockfile.Version), "files"))
	if err != nil {
		return LockfileMetadata{}, nil, err
	}

	files = strings.Split(string(fileBytes), "\n")

	err = os.RemoveAll(path.Join(LocalPath, fmt.Sprintf("%s-%s", lockfile.Name, lockfile.Version)))

	return lockfile, files, err
}
