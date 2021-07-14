package config

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/util"
)

type LockfileMetadata struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Url     string `toml:"url"`
	Sha256  string `toml:"sha256"`
}

type LockfileStruct struct {
	Packages []LockfileMetadata `toml:"packages" mapstructure:"packages"`
}

func readLockfile() (lockfile *LockfileStruct, err error) {
	lockFile, err := os.ReadFile(util.LockfilePath)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(lockFile, &lockfile)

	return lockfile, err
}

// Clone & add repository to config
func AddPkgToLockfile(metadata LockfileMetadata) (err error) {
	lockfile, err := readLockfile()
	if err != nil {
		return err
	}
	lockfile.Packages = append(lockfile.Packages, metadata)

	newLockfile, err := toml.Marshal(&lockfile)
	if err != nil {
		return err
	}

	err = os.WriteFile(util.LockfilePath, newLockfile, 0666)
	return err
}
