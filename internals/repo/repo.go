package repo

import (
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
