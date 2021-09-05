package pkg

import (
	"fmt"
	"io"
	"net/http"

	"github.com/stewproject/stew/internals/errors"
	"github.com/stewproject/stew/internals/repo"
	"github.com/stewproject/stew/util"
)

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
// Probaly going to clean this up and seperate it in seperate functions later.
func GetPackage(pkgName string, pkgVersion *string) (pkgData *PkgData, err error) {
	// search core repository
	resp, err := http.Get(fmt.Sprintf("%s/%s/package.toml", repo.CoreRepositoryURL, pkgName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check if response is 200
	if resp.StatusCode != 200 {
		return nil, errors.PackageNotFoundError{Package: pkgName}
	}

	// found, get package definition
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pkgDef, err := ParsePackage(body)
	if err != nil {
		return nil, err
	}

	var version string

	// get version metadata
	if pkgVersion == nil {
		// latest version
		version = pkgDef.Package.Version
	} else {
		version = *pkgVersion
	}

	resp, err = http.Get(fmt.Sprintf("%s/%s/%s/metadata.toml", repo.CoreRepositoryURL, pkgName, version))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	verData, err := ParseVersion(body)

	var plfData PlatformData

	// get platform data
	if util.Arch == "silicon" {
		plfData = verData.Silicon
	} else if util.Arch == "intel" {
		plfData = verData.Intel
	}

	pkgUrl := fmt.Sprintf("%s/%s/%s/%s-%s-%s.tar.xz", repo.CorePackagesURL, pkgName, version, pkgName, version, util.Arch)

	return &PkgData{PkgDef: pkgDef, VerData: verData, PlfData: plfData, Repository: "core", Version: version, PkgUrl: pkgUrl}, err
}
