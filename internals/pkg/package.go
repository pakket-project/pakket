package pkg

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pakket-project/pakket/internals/errors"
	"github.com/pakket-project/pakket/internals/repo"
	"github.com/pakket-project/pakket/util"
)

// for use with GetPackage(). Contains all data needed to install a package.
type PkgData struct {
	PkgDef     PackageDefinition
	VerData    VersionMetadata
	PlfData    PlatformData
	TarURL     string
	RepoURL    string
	Version    string
	Repository string
	BinSize    int64
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
	arch := util.Arch
	if arch == "silicon" {
		plfData = verData.Silicon
	} else if arch == "intel" {
		plfData = verData.Intel
	}

	pkgUrl := fmt.Sprintf("%s/%s/%s/%s-%s-%s.tar.xz", repo.CorePackagesURL, pkgName, version, pkgName, version, arch)
	pkgRepoUrl := fmt.Sprintf("%s/%s/%s", repo.CoreRepositoryURL, pkgName, version)

	// get pkg size
	size, err := GetPackageSize(pkgUrl)

	return &PkgData{PkgDef: pkgDef, VerData: verData, PlfData: plfData, Repository: "core", Version: version, TarURL: pkgUrl, RepoURL: pkgRepoUrl, BinSize: size}, err
}

func GetPackageSize(url string) (bytes int64, err error) {
	resp, err := http.Head(url)
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
