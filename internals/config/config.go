package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
)

var (
	// path that houses config, lockfile, etc
	ConfigPath = "/etc/pakket"
	// Path to main pakket config
	ConfigFile = path.Join(ConfigPath, "config.toml")
	// Path to lockfile w/ currently installed packages
	LockfilePath = path.Join(ConfigPath, "lockfile.toml")
	// Path to local database
	LocalPath = path.Join(ConfigPath, "local")
)

type Mirror struct {
	URL  string `toml:"url"`
	Name string `toml:"name"`
}

type Paths struct {
	// absolute path to temp downloads directory (ie /var/tmp/pakket)
	Downloads string `toml:"downloads"`
	// installation prefix. must be "/usr/local", other values are currently not supported
	Prefix string `toml:"prefix"`
}

type ConfigStruct struct {
	Mirrors   []Mirror `toml:"mirrors"`
	Community bool     `toml:"community"`
	Paths     Paths    `toml:"paths"`
}

var (
	// Config
	C ConfigStruct
)

func GetConfig() (err error) {
	var file []byte

	file, err = os.ReadFile(ConfigFile)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(ConfigFile)

		if err != nil {
			return err
		}

		file, err = os.ReadFile(ConfigFile)
		if err != nil {
			return err
		}
	}

	err = toml.Unmarshal(file, &C)

	if C.Paths.Prefix != "/usr/local" {
		fmt.Println("prefix must be /usr/local, changing it for you...")

		C.Paths.Prefix = "/usr/local"
		err := WriteConfig()
		if err != nil {
			return err
		}
	}

	if len(C.Mirrors) == 0 {
		fmt.Println("no mirrors, automatically adding one...")

		C.Mirrors = append(C.Mirrors, Mirror{URL: "https://core.pakket.sh", Name: "Main Pakket mirror"})
		err := WriteConfig()
		if err != nil {
			return err
		}
	}

	if C.Paths.Downloads == "" {
		fmt.Println("download directory not set, automatically setting it...")

		C.Paths.Downloads = "/var/tmp/pakket"
		err := WriteConfig()
		if err != nil {
			return err
		}
	}
	return err
}

func WriteConfig() (err error) {
	newConfig, err := toml.Marshal(&C)
	if err != nil {
		return err
	}

	err = os.WriteFile(ConfigFile, newConfig, 0666)
	return err
}

// // Add repository to config
// func AddMirror(mirror Mirror) error {
// 	Config.Mirrors = append(Config.Mirrors, mirror)

// 	config, err := toml.Marshal(&Config)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(util.ConfigFile, config, 0660)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Delete mirror
// func DelMirror(configIndex int) error {
// 	Config.Mirrors = append(Config.Mirrors[:configIndex], Config.Mirrors[configIndex+1:]...)

// 	config, err := toml.Marshal(&Config)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(util.ConfigFile, config, 0660)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
