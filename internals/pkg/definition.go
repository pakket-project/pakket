package pkg

import (
	"github.com/pelletier/go-toml"
)

type Package struct {
	Name         string               `toml:"name"`
	Description  string               `toml:"description"`
	Version      string               `toml:"version"`
	License      string               `toml:"license"`
	Homepage     string               `toml:"homepage"`
	Url          string               `toml:"url"`
	Sha256       string               `toml:"sha256"`
	Dependencies []DependencyMetadata `toml:"dependencies"`
}

type DependencyMetadata struct {
	Name     string `toml:"name"`
	Optional bool   `toml:"optional"`
	Version  string `toml:"version"`
	Build    bool   `toml:"build"`
}

type BinariesMetadata struct {
	Url       string `toml:"url"`
	BinPath   string `toml:"binPath"`
	Version   string `toml:"version"`
	Sha256    string `toml:"sha256"`
	Available bool   `toml:"available"`
}

type Binaries struct {
	BigSur_arm BinariesMetadata `toml:"big_sur_arm"`
	BigSur     BinariesMetadata `toml:"big_sur"`
	Catalina   BinariesMetadata `toml:"catalina"`
	Mojave     BinariesMetadata `toml:"mojave"`
}

type Definition struct {
	Package  Package  `toml:"package"`
	Binaries Binaries `toml:"binaries"`
}

func ParseDefinition(definition []byte) (Definition, error) {
	var def Definition
	err := toml.Unmarshal(definition, &def)

	return def, err
}
