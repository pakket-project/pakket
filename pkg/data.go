package pkg

import (
	"github.com/pelletier/go-toml/v2"
)

type PackageDefinition struct {
	Package Package `toml:"package"`
}

type Package struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Version     string `toml:"version"`
	License     string `toml:"license"`
	Homepage    string `toml:"homepage"`
	Homebrew    bool   `toml:"homebrew,omitempty"`
}

type Dependencies struct {
	Dependencies         []string `toml:"dependencies,multiline,omitempty"`
	BuildDependencies    []string `toml:"buildDependencies,multiline,omitempty"`
	OptionalDependencies []string `toml:"optionalDependencies,multiline,omitempty"`
}

type PlatformData struct {
	Checksum string `toml:"checksum,omitempty"`
}

type VersionMetadata struct {
	Url          string       `toml:"url"`
	Dependencies Dependencies `toml:"dependencies,omitempty"`
	Intel        PlatformData `toml:"intel,omitempty"`
	Silicon      PlatformData `toml:"silicon,omitempty"`
}

// Convert data to version metadata
func ParseVersion(versionMetadata []byte) (VersionMetadata, error) {
	metadata := VersionMetadata{Intel: PlatformData{Checksum: ""}, Silicon: PlatformData{Checksum: ""}}
	err := toml.Unmarshal(versionMetadata, &metadata)

	return metadata, err
}

// Convert data to pkg definition
func ParsePackage(packageDefinition []byte) (PackageDefinition, error) {
	var def PackageDefinition
	err := toml.Unmarshal(packageDefinition, &def)

	return def, err
}
