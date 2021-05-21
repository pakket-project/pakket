package repo

import (
	"errors"
	"fmt"
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

// Repository already exists
type RepositoryAlreadyExistsError struct {
	Repository string
}

func (repo RepositoryAlreadyExistsError) Error() string {
	return fmt.Sprintf("repository %s already exists", repo.Repository)
}

// Repository already exists, not defined in config
type UndefinedRepositoryAlreadyExistsError struct {
	Repository string
}

func (repo UndefinedRepositoryAlreadyExistsError) Error() string {
	return fmt.Sprintf("repository %s already exists, but is not defined in the config", repo.Repository)
}

// Get metadata of repository
func GetMetadataFromRepo(repo string) Metadata {
	data, err := os.ReadFile(path.Join(util.RepoPath, repo, "metadata.toml"))
	if err != nil {
		panic(err)
	}

	var def Metadata
	err = toml.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	return def
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
func AddRepo(gitURL string) (metadata *Metadata, err error) {
	// Clone repository to temp dir
	_, err = git.PlainClone(util.DownloadPath, false, &git.CloneOptions{
		URL: gitURL,
		// SingleBranch: true,
		// Depth: 1, // TODO: cloning with depth & pulling seems broken (https://github.com/go-git/go-git/issues/305)
	})
	if err != nil {
		return nil, err
	}

	// Get metadata
	metadatabytes, err := os.ReadFile(path.Join(util.DownloadPath, "metadata.toml"))
	if err != nil {
		os.RemoveAll(util.DownloadPath)
		return nil, err
	}
	metadata = UnmarshalMetadata(metadatabytes)

	// Path to repo
	repoPath := path.Join(util.RepoPath, metadata.Repository.Name)

	// Check if already exists
	if exists := util.DoesPathExist(repoPath); exists {
		os.RemoveAll(util.DownloadPath)

		// Check if repo is defined in config
		for i := range config.Config.Repositories.Locations {
			if config.Config.Repositories.Locations[i].Name == metadata.Repository.Name {
				// already defined in config
				return metadata, RepositoryAlreadyExistsError{Repository: metadata.Repository.Name}
			}
		}
		// not defined in config
		return metadata, UndefinedRepositoryAlreadyExistsError{Repository: metadata.Repository.Name}
	}

	// Check if name contains subfolder
	if subfolder := strings.Contains(metadata.Repository.Name, "/"); subfolder {
		paths := strings.Split(metadata.Repository.Name, "/")
		if len(paths) > 2 {
			return metadata, errors.New("repository name can only contain one slash (/)")
		}
		// create directory
		err = os.Mkdir(path.Join(util.RepoPath, paths[0]), 0777)

		if err != nil {
			if os.IsExist(err) {
			} else {
				os.RemoveAll(util.DownloadPath)
				return metadata, err
			}
		}
	}

	// Move temp repo path to /usr/local/stew/repositories
	err = os.Rename(util.DownloadPath, repoPath)
	if err != nil {
		os.RemoveAll(util.DownloadPath)
		return metadata, err
	}

	// Add to config
	err = config.AddRepo(config.RepositoriesMetadata{Name: metadata.Repository.Name, Path: repoPath, PackagesPath: metadata.Repository.PackagesPath, GitURL: gitURL, IsGit: true})
	if err != nil {
		return metadata, err
	}

	return metadata, err
}
