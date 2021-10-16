package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

var (
	LockFile LockfileStruct
)

type LockfileMetadata struct {
	Name       string   `toml:"name"`
	Version    string   `toml:"version"`
	Checksum   string   `toml:"checksum"`
	Repository string   `toml:"repository"`
	Files      []string `toml:"files"`
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
func AddPkgToLockfile(metadata LockfileMetadata) (err error) {
	if len(LockFile.Packages) == 0 {
		LockFile.Packages = make(map[string]LockfileMetadata)
	}

	LockFile.Packages[metadata.Name] = metadata

	newLockfile, err := toml.Marshal(&LockFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(LockfilePath, newLockfile, 0666)
	return err
}

// Remove package information from lockfile
func RemovePkgFromLockfile(name string) (lockfileData LockfileMetadata, err error) {
	if _, ok := LockFile.Packages[name]; !ok {
		return LockfileMetadata{}, errors.New("package not found")
	}
	lockfile := LockFile.Packages[name]
	delete(LockFile.Packages, name)

	newLockfile, err := toml.Marshal(&LockFile)
	if err != nil {
		return LockfileMetadata{}, err
	}

	err = os.WriteFile(LockfilePath, newLockfile, 0666)
	return lockfile, err
}
