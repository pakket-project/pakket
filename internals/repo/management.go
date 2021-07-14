package repo

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

// Add repository
func Add(gitURL string) (metadata *Metadata, err error) {
	os.RemoveAll(util.DownloadPath)

	// Clone repository to temp dir
	_, err = git.PlainClone(util.DownloadPath, false, &git.CloneOptions{
		URL: gitURL,
		// SingleBranch: true,
		// Depth: 1, // TODO: cloning with depth & pulling seems broken (https://github.com/go-git/go-git/issues/305)
	})
	if err != nil {
		return nil, err
	}

	// Check if repository is valid Stew repository
	if exists := util.DoesPathExist(path.Join(util.DownloadPath, "metadata.toml")); !exists {
		return nil, UnvalidRepository{GitLink: gitURL}
	}

	// Get metadata
	metadatabytes, err := os.ReadFile(path.Join(util.DownloadPath, "metadata.toml"))
	if err != nil {
		os.RemoveAll(util.DownloadPath)
		return nil, err
	}
	metadata = UnmarshalMetadata(metadatabytes)

	// Path to repo
	repoPath := path.Join(util.RepoPath, fmt.Sprintf("%s/%s", metadata.Repository.Author, metadata.Repository.Name))

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

	// create directory
	err = os.Mkdir(path.Join(util.RepoPath, metadata.Repository.Author), 0770)

	if err != nil {
		if os.IsExist(err) {
		} else {
			os.RemoveAll(util.DownloadPath)
			return metadata, err
		}
	}

	// Move temp repo path to /usr/local/stew/repositories
	err = os.Rename(util.DownloadPath, repoPath)
	if err != nil {
		os.RemoveAll(util.DownloadPath)
		return metadata, err
	}

	// delete temp repo
	os.RemoveAll(util.DownloadPath)

	// Add to config
	err = config.AddRepo(config.RepositoriesMetadata{
		Name:         metadata.Repository.Name,
		Author:       metadata.Repository.Author,
		Path:         repoPath,
		PackagesPath: metadata.Repository.PackagesPath,
		GitURL:       gitURL,
		IsGit:        true,
	})
	if err != nil {
		return metadata, err
	}

	return metadata, err
}

func Delete(configIndex int) error {
	repo := config.Config.Repositories.Locations[configIndex]
	err := os.RemoveAll(repo.Path)
	if err != nil {
		return err
	}

	aPath := strings.Split(repo.Path, "/")
	authorPath := strings.Join(aPath[:len(aPath)-1], "/")

	empty, err := util.IsEmpty(authorPath)
	if err != nil {
		return err
	}

	if empty {
		err = os.Remove(authorPath)
		if err != nil {
			return err
		}
	}
	err = config.DelRepo(configIndex)
	if err != nil {
		return err
	}

	return nil
}
