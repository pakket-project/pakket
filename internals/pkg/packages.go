package pkg

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

var (
	silicon = "arm64"
	intel   = "amd64"
)

// Manuals
// 1 - General commands
// 2 - System calls
// 3 - Subroutines
// 4 - Special files
// 5 - File formats
// 6 - Games
// 7 - Macros and conventions
// 8 - Maintenance commands
// 9 - Kernel interface
// n - New commands
var (
	locations = map[string]string{
		// Binaries
		"/bin": "/usr/local/bin",
		// Completions
		// "/completions/fish": "",
		// "/completions/zsh":  "/usr/share/zsh/5.3/functions",
		// Manuals
		"/man/man1": "/usr/local/share/man/man1",
		"/man/man2": "/usr/local/share/man/man2",
		"/man/man3": "/usr/local/share/man/man3",
		"/man/man4": "/usr/local/share/man/man4",
		"/man/man5": "/usr/local/share/man/man5",
		"/man/man6": "/usr/local/share/man/man6",
		"/man/man7": "/usr/local/share/man/man7",
		"/man/man8": "/usr/local/share/man/man8",
		"/man/man9": "/usr/local/share/man/man9",
		"/man/mann": "/usr/local/share/man/mann",
	}
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

func DownloadPackage(url string) (tarPath string, err error) {
	err = os.MkdirAll(util.DownloadPath, 0774)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return "", err
	}
	tarPath = path.Join(util.DownloadPath, params["filename"])
	// Create the file
	out, err := os.Create(tarPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return tarPath, err
}

func InstallPackage(pkgName, version string) (err error) {
	fmt.Println("Getting metadata...")
	pkg, pkgPath, err := GetPackageMetadata(pkgName) // Get package metadata
	if err != nil {
		return err
	}

	// Get package version metadata
	var ver *VersionMetadata
	if version == "latest" {
		ver, err = GetPackageVersion(pkgName, *pkgPath, pkg.Package.Version)
		if err != nil {
			return err
		}
	} else {
		ver, err = GetPackageVersion(pkgName, *pkgPath, version)
		if err != nil {
			return err
		}
	}

	fmt.Println("Getting correct version...")
	var binary BinaryMetadata
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

	fmt.Println("Downloading tarball...")
	tarPath, err := DownloadPackage(binary.Url) // Download package, save tar to tarPath
	if err != nil {
		return err
	}

	fmt.Println("Untarring...")
	// Unarchive tarball
	savePath := path.Join(util.DownloadPath, pkg.Package.Name)
	err = archiver.Unarchive(tarPath, savePath)
	if err != nil {
		return err
	}
	defer os.RemoveAll(savePath)
	err = os.RemoveAll(tarPath)
	if err != nil {
		return err
	}

	filesToMove := make(map[string]string)

	fmt.Println("Checking files...")
	err = filepath.Walk(savePath, func(filePath string, f os.FileInfo, err error) error {
		if !f.IsDir() || filePath == savePath {
			return err
		}

		localpath := strings.ReplaceAll(filePath, path.Join(savePath, pkg.Package.Name), "")

		if localpath == "" || localpath == "info.toml" {
			return err
		}
		if finalPath, ok := locations[localpath]; ok {
			err = filepath.Walk(filePath, func(secondPath string, f os.FileInfo, err error) error {
				if secondPath == filePath {
					return err
				}

				filesToMove[secondPath] = path.Join(finalPath, f.Name())
				return err
			})
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}

	fmt.Println("Moving files...")
	for oldPath, newPath := range filesToMove {
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
	}
	fmt.Println("Done")
	return nil
}
