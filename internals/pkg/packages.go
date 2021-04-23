package pkg

import (
	"fmt"
	"os"
	"path"

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
func GetPackageData(packageName string) (*PackageDefinition, error) {
	for i := 0; i < len(config.Config.Repositories.Locations); i++ {
		repo := config.Config.Repositories.Locations[i]
		packagePath := path.Join(repo.Path, repo.PackagesPath, packageName)

		if exists := util.DoesPathExist(packagePath); !exists {
			continue
		}
		fmt.Println(path.Join(packagePath, "package.toml"))
		data, err := os.ReadFile(path.Join(packagePath, "package.toml"))
		if err != nil {
			panic(err)
		}

		def, err := ParsePackage(data)
		if err != nil {
			panic(err)
		}

		return &def, nil
	}

	return nil, PackageNotFoundError{Package: packageName}
}
