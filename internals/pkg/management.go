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
	err = os.MkdirAll(config.C.Paths.Downloads, 0770)
	if err != nil {
		fmt.Println("no download path specified in config")
		return
	}

	// Download tar
	resp, err := grab.Get(config.C.Paths.Downloads, pkg.TarURL)
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
	if config.Lockfile.Exists(pkg.PkgDef.Package.Name) {
		return fmt.Errorf("%s is already installed", pkg.PkgDef.Package.Name)
	}

	savePath := path.Join(config.C.Paths.Downloads, pkg.PkgDef.Package.Name)

	// install dependencies
	for _, dep := range pkg.VerData.Dependencies.Dependencies {
		var name string
		var version *string

		if strings.Contains(dep, "@") {
			splitted := strings.Split(dep, "@")
			name = splitted[0]
			version = &splitted[1]
		} else {
			name = dep
			version = nil
		}

		pkgData, err := GetPackage(name, version)
		if err != nil {
			fmt.Printf("error while installing %s: %s\n", dep, err.Error())
			continue
		}

		// TODO: dont force yes
		err = InstallPackage(*pkgData, false, true)
		if err != nil {
			fmt.Printf("error while installing %s: %s\n", dep, err.Error())
			continue
		}

		fmt.Printf("installed dependency %s@%s\n", pkgData.PkgDef.Package.Name, pkgData.PkgDef.Package.Version)
	}

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

		filesToMove[filePath] = path.Join(config.C.Paths.Prefix, localPath)

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
	err = config.Lockfile.Add(config.NewMetadata(pkg.PkgDef.Package.Name, pkg.Version, pkg.PlfData.Checksum, pkg.Repository), finalPaths)
	if err != nil {
		return err
	}

	// run postinstall script
	err = HandleScript("postinstall", pkg, savePath, yes)
	if err != nil {
		return err
	}

	return nil
}

func RemovePackage(pkg string) (err error) {
	// remove from lockfile
	_, files, err := config.Lockfile.Remove(pkg)
	if err != nil {
		return err
	}

	for _, path := range files {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
