package config

import (
	"fmt"

	"github.com/spf13/viper"
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

type Repositories struct {
	Locations []RepositoriesMetadata `toml:"locations"`
}

type ConfigStruct struct {
	Repositories Repositories `toml:"repositories"`
}

var (
	Config ConfigStruct
)

func GetConfig() ConfigStruct {
	viper.SetDefault("ContentDir", "content")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(util.StewPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}

	return Config
}
