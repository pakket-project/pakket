package pkg

import (
	sha "crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/grab"
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
		"/bin":   "/usr/local/bin",
		"/share": "/usr/local/share",
		// Completions
		// "/completions/fish": "",
		// "/completions/zsh":  "/usr/share/zsh/5.3/functions",
		// Manuals
		// "/man/man1": "/usr/local/share/man/man1",
		// "/man/man2": "/usr/local/share/man/man2",
		// "/man/man3": "/usr/local/share/man/man3",
		// "/man/man4": "/usr/local/share/man/man4",
		// "/man/man5": "/usr/local/share/man/man5",
		// "/man/man6": "/usr/local/share/man/man6",
		// "/man/man7": "/usr/local/share/man/man7",
		// "/man/man8": "/usr/local/share/man/man8",
		// "/man/man9": "/usr/local/share/man/man9",
		// "/man/mann": "/usr/local/share/man/mann",
	}
)

// download and unarchive package.
func DownloadPackage(pkg PkgData, savePath string) (err error) {
	err = os.MkdirAll(util.DownloadPath, 0770)
	if err != nil {
		return
	}

	// Download tar
	resp, err := grab.Get(util.DownloadPath, pkg.PkgUrl)
	defer os.RemoveAll(resp.Filename)
	if err != nil {
		return
	}

	fileData, err := os.ReadFile(resp.Filename)
	if err != nil {
		return err
	}

	hashBytes := sha.Sum256(fileData)
	downloadHash := hex.EncodeToString(hashBytes[:])

	if downloadHash != pkg.PlfData.hash {
		return InvalidHash{
			Repository: pkg.Repository,
		}
	}

	// unarchive
	err = archiver.Unarchive(resp.Filename, savePath)
	if err != nil {
		return
	}

	return err
}

func InstallPackage(pkg PkgData) (err error) {
	savePath := path.Join(util.DownloadPath, pkg.PkgDef.Package.Name)

	err = DownloadPackage(pkg, savePath) // Download package, save tar to tarPath
	defer os.RemoveAll(savePath)
	if err != nil {
		return err
	}

	filesToMove := make(map[string]string)

	err = filepath.Walk(savePath, func(filePath string, f os.FileInfo, err error) error {
		if filePath == savePath {
			return err
		}

		localPath := strings.ReplaceAll(filePath, path.Join(savePath, pkg.PkgDef.Package.Name), "")

		if localPath == "" || localPath == "/info.toml" {
			return err
		}

		// if file is in /share directory
		if strings.Split(localPath, "/")[1] == "share" {
			if f.IsDir() {
				return err
			}

			filesToMove[filePath] = path.Join("/", "usr", "local", localPath)
			return err
		}

		// if path is a file, return
		if !f.IsDir() {
			return err
		}

		if finalPath, ok := locations[localPath]; ok {
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

	for oldPath, newPath := range filesToMove {
		var exists bool
		var confirm bool

		exists = util.DoesPathExist(newPath)
		if exists {
			confirm = util.DestructiveConfirm(fmt.Sprintf("File %s already exists. Overwrite?", newPath))
			err := os.Remove(newPath)
			if err != nil {
				return err
			}
		}

		if (!exists) || (exists && confirm) {
			err = os.MkdirAll(path.Dir(newPath), 0755)
			if err != nil {
				return err
			}

			err = os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
		}
	}

	// add to lockfile
	err = config.AddPkgToLockfile(config.LockfileMetadata{Name: pkg.PkgDef.Package.Name, Version: pkg.Version, Sha256: pkg.PlfData.hash, Repository: pkg.Repository})
	if err != nil {
		return err
	}

	return nil
}
