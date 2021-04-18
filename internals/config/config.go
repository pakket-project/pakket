package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"github.com/stewproject/stew/util"
)

type RepositoriesMetadata struct {
	Name         string `toml:"name"`
	Path         string `toml:"path"`
	PackagesPath string `toml:"packagesPath"`
	GitURL       string `toml:"gitUrl"`
}

type Repositories struct {
	Locations []RepositoriesMetadata `toml:"locations" mapstructure:"locations"`
}

type ConfigStruct struct {
	Repositories Repositories `toml:"repositories"`
}

var (
	Config ConfigStruct
)

func GetConfig() ConfigStruct {
	// viper.SetDefault("ContentDir", "content")
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

func AddRepo(repoMetadata RepositoriesMetadata) error {
	GetConfig()

	Config.Repositories.Locations = append(Config.Repositories.Locations, repoMetadata)

	config, err := toml.Marshal(&Config)
	if err != nil {
		return err
	}

	os.WriteFile(util.ConfigFile, config, 0666)
	return nil
}
