package pkg

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/go-vela/archiver/v3"
	"github.com/pakket-project/pakket/internals/config"
	"github.com/pakket-project/pakket/internals/errors"
	"github.com/pakket-project/pakket/util"
)

// download and unarchive package.
func DownloadPackage(pkg PkgData, savePath string) (err error) {
	err = os.MkdirAll(util.DownloadPath, 0770)
	if err != nil {
		return
	}

	// Download tar
	resp, err := grab.Get(util.DownloadPath, pkg.TarURL)
	defer os.RemoveAll(resp.Filename)
	if err != nil {
		return err
	}

	fileData, err := os.ReadFile(resp.Filename)
	if err != nil {
		return err
	}

	downloadChecksum := fmt.Sprintf("%x", sha256.Sum256(fileData))

	if downloadChecksum != pkg.PlfData.Checksum {
		return errors.InvalidChecksum{
			Package: pkg.PkgDef.Package.Name,
		}
	}

	// unarchive
	err = archiver.Unarchive(resp.Filename, savePath)
	if err != nil {
		return err
	}

	return err
}

func InstallPackage(pkg PkgData, force bool, yes bool) (err error) {
	// check if package is already installed
	if v, ok := config.LockFile.Packages[pkg.PkgDef.Package.Name]; ok {
		return fmt.Errorf("%s is already installed", v.Name)
	}

	savePath := path.Join(util.DownloadPath, pkg.PkgDef.Package.Name)

	// run preinstall script
	err = HandleScript("preinstall", pkg, savePath, yes)
	if err != nil {
		return err
	}

	err = DownloadPackage(pkg, savePath) // Download package, save tar to tarPath
	defer os.RemoveAll(savePath)
	if err != nil {
		return err
	}

	filesToMove := make(map[string]string)
	var finalPaths []string

	err = filepath.Walk(savePath, func(filePath string, f os.FileInfo, err error) error {
		if filePath == savePath {
			return err
		}

		localPath := strings.ReplaceAll(filePath, path.Join(savePath, pkg.PkgDef.Package.Name), "")

		if localPath == "" || localPath == "/info.toml" {
			return err
		}

		// if path is a dir, return
		if f.IsDir() {
			return err
		}

		filesToMove[filePath] = path.Join("/", "usr", "local", localPath)

		return err
	})
	if err != nil {
		return err
	}

	always := false
	for oldPath, newPath := range filesToMove {
		var exists bool
		var confirm bool

		exists = util.DoesPathExist(newPath)
		if exists && !force {
			if !always {
				confirm, always = util.DestructiveConfirm(fmt.Sprintf("File %s already exists. Overwrite?", newPath), true)
			}
			err := os.Remove(newPath)
			if err != nil {
				return err
			}
		}

		if exists && force {
			err := os.Remove(newPath)
			if err != nil {
				return err
			}
		}

		if (!exists) || (exists && confirm) || (exists && force) || (exists && always) {
			err = os.MkdirAll(path.Dir(newPath), 0755)
			if err != nil {
				return err
			}

			finalPaths = append(finalPaths, newPath)
			err = os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
		}
	}

	// add to lockfile
	err = config.AddPkgToLockfile(config.LockfileMetadata{Name: pkg.PkgDef.Package.Name, Version: pkg.Version, Checksum: pkg.PlfData.Checksum, Repository: pkg.Repository, Files: finalPaths})
	if err != nil {
		return err
	}

	//run postinstall script
	err = HandleScript("postinstall", pkg, savePath, yes)
	if err != nil {
		return err
	}

	return nil
}

func RemovePackage(pkg string) (err error) {
	// remove from lockfile
	lockfile, err := config.RemovePkgFromLockfile(pkg)
	if err != nil {
		return err
	}

	for _, path := range lockfile.Files {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
