package pkg

import (
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

type PackageNotFoundError struct {
	Package string
}

func (pkg PackageNotFoundError) Error() string {
	return fmt.Sprintf("package %s not found", pkg.Package)
}

// Search all repositories for specific package
func GetPackageData(packageName string) (*Definition, error) {
	for i := 0; i < len(config.Config.Repositories.Locations); i++ {
		repo := config.Config.Repositories.Locations[i]
		packagePath := repo.Path + repo.PackagesPath + "/" + packageName

		if exists := util.DoesPathExist(packagePath); !exists {
			continue
		}

		data, err := os.ReadFile(path.Join(packagePath, "definition.toml"))
		if err != nil {
			panic(err)
		}
		var def Definition
		err = toml.Unmarshal(data, &def)
		if err != nil {
			panic(err)
		}
		return &def
	}
	return &Definition{} //TODO
}
