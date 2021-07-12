package repo

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

// Add repository
func Add(gitURL string) (metadata *Metadata, err error) {
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

