package pkg

import (
	"fmt"
)

type PackageNotFoundError struct {
	Package string
}

func (pkg PackageNotFoundError) Error() string {
	return fmt.Sprintf("package %s not found", pkg.Package)
}

type VersionNotFoundError struct {
	Package string
	Version string
}

func (pkg VersionNotFoundError) Error() string {
	return fmt.Sprintf("version %s of package %s not found", pkg.Version, pkg.Package)
}
