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

func GetBinaryMetadata(ver VersionMetadata) (binary BinaryMetadata) {
	macVer := util.GetVersion() // Get macOS version

	// Get binary information
	if runtime.GOARCH == silicon { // Apple Silicon
		for _, b := range ver.Binaries.Silicon {
			for _, v := range b.SupportedVersions {
				if v == macVer || v == "all" {
					binary = b
				}
			}
		}
	} else if runtime.GOARCH == intel { // Intel
		for _, b := range ver.Binaries.Intel {
			for _, v := range b.SupportedVersions {
				if v == macVer || v == "all" {
					binary = b
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
