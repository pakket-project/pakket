package pkg

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

func GetPackageData(packageName string) *Definition {
	for i := 0; i < len(config.Config.Repositories); i++ {
		repo := config.Config.Repositories[i]
		packagePath := repo.Path + repo.PackagesPath + "/" + packageName

		if exists := util.DoesPathExist(packagePath); !exists {
			continue
		}

		data, err := os.ReadFile(packagePath + "/definition.toml")
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
