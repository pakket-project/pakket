package pkg

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

// for use with GetPackage(). Contains all data needed to install a package.
type PkgData struct {
	PkgDef     PackageDefinition
	PkgPath    string
	VerData    VersionMetadata
	Version    string
	BinData    BinaryMetadata
	BinSize    int64
	Repository string
}

// One function to get all information needed to install a package. Version should be "latest" for latest version. binSize is the size of the tarball in bytes.
func GetPackage(pkgName string, pkgVer string) (PkgData, error) {
	pkgDef, pkgPath, repo, err := GetPackageMetadata(pkgName) // get package metadata
	if err != nil {
		return PkgData{BinSize: 0}, err
	}

	var version string
	if pkgVer == "latest" {
		version = pkgDef.Package.Version
	} else {
		version = pkgVer
	}

	verData, err := GetPackageVersion(pkgName, *pkgPath, version)
	if err != nil {
		return PkgData{PkgDef: *pkgDef, PkgPath: *pkgPath, BinSize: 0, Repository: repo}, err
	}

	binData := GetBinaryMetadata(*verData)
	binSize, err := GetPackageSize(*binData)
	if err != nil {
		return PkgData{PkgDef: *pkgDef, PkgPath: *pkgPath, VerData: *verData, BinData: *binData, BinSize: 0, Repository: repo}, err
	}

	return PkgData{PkgDef: *pkgDef, PkgPath: *pkgPath, VerData: *verData, Version: version, BinData: *binData, BinSize: binSize, Repository: repo}, err
}

// Search all repositories for specific package
func GetPackageMetadata(packageName string) (pkgDef *PackageDefinition, pkgPath *string, repository string, err error) {
	for i := 0; i < len(config.Config.Repositories.Locations); i++ {
		repo := config.Config.Repositories.Locations[i]
		repository = fmt.Sprintf("%s/%s", repo.Author, repo.Name)
		packagePath := path.Join(repo.Path, repo.PackagesPath, packageName)

		if exists := util.DoesPathExist(packagePath); !exists {
			continue
		}

		data, err := os.ReadFile(path.Join(packagePath, "package.toml"))
		if err != nil {
			return nil, &packagePath, repository, err
		}

		def, err := ParsePackage(data)
		if err != nil {
			return &def, &packagePath, repository, err
		}

		return &def, &packagePath, repository, nil
	}

	return nil, nil, repository, PackageNotFoundError{Package: packageName}
}

func GetPackageVersion(pkgName, pkgPath, version string) (*VersionMetadata, error) {
	data, err := os.ReadFile(path.Join(pkgPath, version, "metadata.toml"))

	if os.IsNotExist(err) {
		return nil, VersionNotFoundError{Package: pkgName, Version: version}
	} else if err != nil {
		return nil, err
	}

	metadata, err := ParseVersion(data)
	if err != nil {
		return &metadata, err
	}

	return &metadata, nil
}

func GetBinaryMetadata(ver VersionMetadata) (binary *BinaryMetadata) {
	macVer := util.GetVersion() // Get macOS version

	// Get binary information
	if runtime.GOARCH == silicon { // Apple Silicon
		for _, b := range ver.Binaries.Silicon {
			for _, v := range b.SupportedVersions {
				if v == macVer || v == "all" {
					binary = &b
				}
			}
		}
	} else if runtime.GOARCH == intel { // Intel
		for _, b := range ver.Binaries.Intel {
			for _, v := range b.SupportedVersions {
				if v == macVer || v == "all" {
					binary = &b
				}
			}
		}
	}

	return binary
}

func GetPackageSize(binary BinaryMetadata) (bytes int64, err error) {
	resp, err := http.Head(binary.Url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	bytes, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
