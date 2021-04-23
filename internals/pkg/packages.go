package pkg

import (
	"fmt"
	"os"
	"path"

	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

type PackageNotFoundError struct {
	Package string
}

func (pkg PackageNotFoundError) Error() string {
	return fmt.Sprintf("package %s not found", pkg.Package)
}

type VersionNotFoundError struct {
	Package string
	Version string
}

func (pkg VersionNotFoundError) Error() string {
	return fmt.Sprintf("version %s of package %s not found", pkg.Version, pkg.Package)
}

// Search all repositories for specific package
func GetPackageMetadata(packageName string) (pkgDef *PackageDefinition, pkgPath *string, err error) {
	for i := 0; i < len(config.Config.Repositories.Locations); i++ {
		repo := config.Config.Repositories.Locations[i]
		packagePath := path.Join(repo.Path, repo.PackagesPath, packageName)

		if exists := util.DoesPathExist(packagePath); !exists {
			continue
		}

		data, err := os.ReadFile(path.Join(packagePath, "package.toml"))
		if err != nil {
			return nil, &packagePath, err
		}

		def, err := ParsePackage(data)
		if err != nil {
			return &def, &packagePath, err
		}

		return &def, &packagePath, nil
	}

	return nil, nil, PackageNotFoundError{Package: packageName}
}

func GetPackageVersion(Package string, version string) (*VersionMetadata, error) {
	_, pkgPath, err := GetPackageMetadata(Package)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path.Join(*pkgPath, version, "metadata.toml"))

	if os.IsNotExist(err) {
		return nil, VersionNotFoundError{Package: Package, Version: version}
	} else if err != nil {
		return nil, err
	}

	metadata, err := ParseVersion(data)
	if err != nil {
		return &metadata, err
	}

	return &metadata, nil
}
