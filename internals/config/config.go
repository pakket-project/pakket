package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"github.com/stewproject/stew/util"
)

type Mirror struct {
	URL  string `toml:"url"`
	Name string `toml:"name"`
}

type ConfigStruct struct {
	Mirrors   []Mirror `toml:"mirrors"`
	Community bool     `toml:"community"`
}

var (
	// Loaded config
	Config   ConfigStruct
	LockFile LockfileStruct
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
func AddMirror(mirror Mirror) error {
	Config.Mirrors = append(Config.Mirrors, mirror)

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
func DelRepo(configIndex int) error {
	Config.Mirrors = append(Config.Mirrors[:configIndex], Config.Mirrors[configIndex+1:]...)

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
