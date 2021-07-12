package repo

import (
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/util"
)

type Repository struct {
	Name         string   `toml:"name"`
	Author       string   `toml:"author"`
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
