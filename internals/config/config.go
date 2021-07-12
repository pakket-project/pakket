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
	Author       string `toml:"author"`
	Path         string `toml:"path"`
	PackagesPath string `toml:"packagesPath"`
	GitURL       string `toml:"gitUrl"`
	IsGit        bool   `toml:"isGit"`
}

type Repositories struct {
	Locations []RepositoriesMetadata `toml:"locations" mapstructure:"locations"`
}

type ConfigStruct struct {
	Repositories Repositories `toml:"repositories"`
}

var (
	// Loaded config
	Config ConfigStruct
)

// Get main Stew config
func GetConfig() ConfigStruct {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(util.StewPath)

	err := viper.ReadInConfig()
	// If config is not found
	if _, errNotExists := err.(viper.ConfigFileNotFoundError); errNotExists {
		// Check if StewPath exists, if not, creates directory
		if exists := util.DoesPathExist(util.StewPath); !exists {
			panic(fmt.Errorf("%s doesn't exist", util.StewPath))
		}

		_, err = os.Create(util.ConfigFile) // create config file
		if err != nil {
			panic(err)
		}

		err = viper.ReadInConfig() // read file again
		if err != nil {
			panic(fmt.Errorf("fatal error reading config file: %s", err))
		}
	} else if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}

	return Config
}

// Add repository to config
func AddRepo(repoMetadata RepositoriesMetadata) error {
	Config.Repositories.Locations = append(Config.Repositories.Locations, repoMetadata)

	config, err := toml.Marshal(&Config)
	if err != nil {
		return err
	}

	err = os.WriteFile(util.ConfigFile, config, 0660)
	if err != nil {
		return err
	}

	return nil
}

// Delete repository
func DelRepo(repoMetadata RepositoriesMetadata) error {
	for i := range Config.Repositories.Locations {
		if repoMetadata == Config.Repositories.Locations[i] {
			Config.Repositories.Locations = append(Config.Repositories.Locations[:i], Config.Repositories.Locations[i+1:]...)
		}

	}

	config, err := toml.Marshal(&Config)
	if err != nil {
		return err
	}

	err = os.WriteFile(util.ConfigFile, config, 0660)
	if err != nil {
		return err
	}

	return nil
}
