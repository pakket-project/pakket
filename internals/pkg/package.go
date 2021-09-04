package pkg

// for use with GetPackage(). Contains all data needed to install a package.
type PkgData struct {
	PkgDef     PackageDefinition
	VerData    VersionMetadata
	PlfData    PlatformData
	PkgUrl     string
	Version    string
	Repository string
}

// One function to get all information needed to install a package. Version should be "latest" for latest version. binSize is the size of the tarball in bytes.
func GetPackage(name string, version string) (*PkgData, error) {
	return nil, nil
}
