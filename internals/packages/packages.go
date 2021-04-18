package packages

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/stewproject/stew/util"
)

func GetPackageData(packageName string) Definition {
	packageDir := util.RepoPath + "/core-packages/packages/" + packageName
	data, err := os.ReadFile(packageDir + "/definition.toml")
	if err != nil {
		panic(err)
	}

	var def Definition // TODO: needed? pointer
	err = toml.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	return def
}
