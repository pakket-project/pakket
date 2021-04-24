package pkg

import (
	"github.com/pelletier/go-toml"
)

type PackageDefinition struct {
	Package Package `toml:"package"`
}

type Package struct {
	Name              string   `toml:"name"`
	Description       string   `toml:"description"`
	Latest            string   `toml:"version"`
	License           string   `toml:"license"`
	Homepage          string   `toml:"homepage"`
	Url               string   `toml:"url"`
	Sha256            string   `toml:"sha256"`
	AvailableVersions []string `toml:"availableVersions"`
}

type Dependencies struct {
	Dependencies         []string `toml:"dependencies"`
	BuildDependencies    []string `toml:"buildDependencies"`
	OptionalDependencies []string `toml:"optionalDependencies"`
}

type BinaryMetadata struct {
	Url               string   `toml:"url"`
	Binpath           string   `toml:"binPath"`
	Sha256            string   `toml:"sha256"`
	SupportedVersions []string `toml:"supportedVersions"`
}

type Binaries struct {
	SupportsRosetta bool             `toml:"supportsRosetta"`
	Intel           []BinaryMetadata `toml:"intel"`
	Silicon         []BinaryMetadata `toml:"silicon"`
}

type VersionMetadata struct {
	Dependencies Dependencies `toml:"dependencies"`
	Binaries     Binaries     `toml:"binaries"`
}

// Convert data to package definition
func ParseVersion(versionMetadata []byte) (VersionMetadata, error) {
	var metadata VersionMetadata
	err := toml.Unmarshal(versionMetadata, &metadata)

	return metadata, err
}

func ParsePackage(packageDefinition []byte) (PackageDefinition, error) {
	var def PackageDefinition
	err := toml.Unmarshal(packageDefinition, &def)

	return def, err
}
