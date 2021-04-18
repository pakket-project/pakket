package repo

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/util"
)

type Repository struct {
	Name         string   `toml:"name"`
	PackagesPath string   `toml:"packagesPath"`
	Maintainers  []string `toml:"maintainers"`
}

type ConfigStruct struct {
	Repository Repository `toml:"repository"`
}

func GetConfig(repo string) *ConfigStruct {
	data, err := os.ReadFile(util.RepoPath + "/" + repo)
	if err != nil {
		panic(err)
	}

	var def ConfigStruct
	err = toml.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	return &def
}

func GetConfigFromData(data []byte) *ConfigStruct {
	var def ConfigStruct
	err := toml.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	return &def
}

func AddRepo(url string) (*git.Repository, error) {
	// Clone repository to temp dir
	repo, err := git.PlainClone(util.TempRepoPath, false, &git.CloneOptions{
		URL:   url,
		Depth: 1,
	})
	if err != nil {
		return repo, err
	}

	// Get config
	configData, err := os.ReadFile(path.Join(util.TempRepoPath, "config.toml"))
	if err != nil {
		util.RemoveFolder(util.TempRepoPath)
		return repo, err
	}
	config := GetConfigFromData(configData)

	// Path to repo
	repoPath := path.Join(util.RepoPath, config.Repository.Name)

	// Check if name contains subfolder
	if subfolder := strings.Contains(config.Repository.Name, "/"); subfolder {
		paths := strings.Split(config.Repository.Name, "/")
		if len(paths) > 2 {
			return repo, errors.New("repository name can only contain one slash (/)")
		}
		// create directory
		err = os.Mkdir(path.Join(util.RepoPath, paths[0]), 0777)

		if err != nil {
			if os.IsExist(err) {
			} else {
				util.RemoveFolder(util.TempRepoPath)
				return repo, err
			}
		}
	}

	err = os.Rename(util.TempRepoPath, repoPath)
	if err != nil {
		util.RemoveFolder(util.TempRepoPath)
		return repo, err
	}

	return repo, err
}
