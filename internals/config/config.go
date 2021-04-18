package config

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/util"
)

// repositories = [
//   {name = "stew/core", path = "/usr/local/stew/repositories/core-packages", packagesPath = "/packages"},
// ]

type RepositoriesMetadata struct {
	Name         string `toml:"name"`
	Path         string `toml:"path"`
	PackagesPath string `toml:"packagesPath"`
}

type ConfigStruct struct {
	Repositories []RepositoriesMetadata `toml:"repositories"`
}

var (
	Config ConfigStruct
)

func GetConfig() ConfigStruct {
	data, err := os.ReadFile(util.ConfigFile)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}

	return Config
}
