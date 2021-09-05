package pkg

import (
	"github.com/pelletier/go-toml/v2"
)

type PackageDefinition struct {
	Package Package `toml:"package"`
}

type Package struct {
	Name              string   `toml:"name"`
	Description       string   `toml:"description"`
	Version           string   `toml:"version"`
	License           string   `toml:"license"`
	Homepage          string   `toml:"homepage"`
	Head              string   `toml:"head"`
	AvailableVersions []string `toml:"availableVersions"`
}

type Dependencies struct {
	Dependencies         []string `toml:"dependencies"`
	BuildDependencies    []string `toml:"buildDependencies"`
	OptionalDependencies []string `toml:"optionalDependencies"`
}

type PlatformData struct {
	Hash string `toml:"hash"`
}

type VersionMetadata struct {
	Url             string       `toml:"url"`
	Dependencies    Dependencies `toml:"dependencies"`
	SupportsRosetta bool         `toml:"supportsRosetta"`
	Intel           PlatformData `toml:"intel"`
	Silicon         PlatformData `toml:"silicon"`
}

// Convert data to version metadata
func ParseVersion(versionMetadata []byte) (VersionMetadata, error) {
	var metadata VersionMetadata
	err := toml.Unmarshal(versionMetadata, &metadata)

	return metadata, err
}

// Convert data to pkg definition
func ParsePackage(packageDefinition []byte) (PackageDefinition, error) {
	var def PackageDefinition
	err := toml.Unmarshal(packageDefinition, &def)

	return def, err
}
