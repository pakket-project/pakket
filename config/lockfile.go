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
	Lockfile LockfileStruct
)

type LockfileMetadata struct {
	Name       string `toml:"name"`
	Version    string `toml:"version"`
	Checksum   string `toml:"checksum"`
	Repository string `toml:"repository"`
}

func NewMetadata(name, version, checksum, repository string) LockfileMetadata {
	return LockfileMetadata{
		Name:       name,
		Version:    version,
		Checksum:   checksum,
		Repository: repository,
	}
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

	err = toml.Unmarshal(file, &Lockfile)

	return err
}

// Add package information to lockfile
func (lock *LockfileStruct) Add(metadata LockfileMetadata, files []string) (err error) {
	if len(lock.Packages) == 0 {
		lock.Packages = make(map[string]LockfileMetadata)
	}

	lock.Packages[metadata.Name] = metadata

	newLockfile, err := toml.Marshal(&lock)
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
func (lock *LockfileStruct) Remove(name string) (lockfileData *LockfileMetadata, files []string, err error) {
	if !lock.Exists(name) {
		return nil, nil, errors.New("package not found")
	}

	lockfile := lock.Packages[name]
	delete(lock.Packages, name)

	newLockfile, err := toml.Marshal(&lock)
	if err != nil {
		return nil, nil, err
	}

	err = os.WriteFile(LockfilePath, newLockfile, 0666)
	if err != nil {
		return nil, nil, err
	}

	// delete files from local db
	fileBytes, err := os.ReadFile(path.Join(LocalPath, fmt.Sprintf("%s-%s", lockfile.Name, lockfile.Version), "files"))
	if err != nil {
		return nil, nil, err
	}

	files = strings.Split(string(fileBytes), "\n")

	err = os.RemoveAll(path.Join(LocalPath, fmt.Sprintf("%s-%s", lockfile.Name, lockfile.Version)))

	return &lockfile, files, err
}

func (lock *LockfileStruct) Exists(pkgName string) bool {
	_, ok := lock.Packages[pkgName]
	return ok
}
