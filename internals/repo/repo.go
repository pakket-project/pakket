package repo

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

type Repository struct {
	Name         string   `toml:"name"`
	PackagesPath string   `toml:"packagesPath"`
	Maintainers  []string `toml:"maintainers"`
}

type Metadata struct {
	Repository Repository `toml:"repository"`
}

// Get metadata of repository
func GetMetadataFromRepo(repo string) *Metadata {
	data, err := os.ReadFile(path.Join(util.RepoPath, repo, "metadata.toml"))
	if err != nil {
		panic(err)
	}

	var def Metadata
	err = toml.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	return &def
}

// Unmarshal metadata
func UnmarshalMetadata(data []byte) *Metadata {
	var def Metadata
	err := toml.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	return &def
}

// Add repository
func AddRepo(gitURL string) error {
	// Clone repository to temp dir
	_, err := git.PlainClone(util.TempRepoPath, false, &git.CloneOptions{
		URL: gitURL,
		// SingleBranch: true,
		// Depth: 1, // TODO: cloning with depth & pulling seems broken (https://github.com/go-git/go-git/issues/305)
	})
	if err != nil {
		return err
	}

	// Get metadata
	Metadata, err := os.ReadFile(path.Join(util.TempRepoPath, "metadata.toml"))
	if err != nil {
		util.RemoveFolder(util.TempRepoPath)
		return err
	}
	metadata := UnmarshalMetadata(Metadata)

	// Path to repo
	repoPath := path.Join(util.RepoPath, metadata.Repository.Name)

	// Check if name contains subfolder
	if subfolder := strings.Contains(metadata.Repository.Name, "/"); subfolder {
		paths := strings.Split(metadata.Repository.Name, "/")
		if len(paths) > 2 {
			return errors.New("repository name can only contain one slash (/)")
		}
		// create directory
		err = os.Mkdir(path.Join(util.RepoPath, paths[0]), 0777)

		if err != nil {
			if os.IsExist(err) {
			} else {
				util.RemoveFolder(util.TempRepoPath)
				return err
			}
		}
	}

	// Move temp repo path to /usr/local/stew/repositories
	err = os.Rename(util.TempRepoPath, repoPath)
	if err != nil {
		util.RemoveFolder(util.TempRepoPath)
		return err
	}

	// Add to config
	err = config.AddRepo(config.RepositoriesMetadata{Name: metadata.Repository.Name, Path: repoPath, PackagesPath: metadata.Repository.PackagesPath, GitURL: gitURL, IsGit: true})
	if err != nil {
		return err
	}

	return err
}
