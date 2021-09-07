package pkg

import (
	"fmt"

	"github.com/stewproject/stew/util/style"
)

type PackageNotFoundError struct {
	Package string
}

func (err PackageNotFoundError) Error() string {
	return fmt.Sprintf("package %s not found", err.Package)
}

type VersionNotFoundError struct {
	Package string
	Version string
}

func (err VersionNotFoundError) Error() string {
	return fmt.Sprintf("version %s of package %s not found", err.Version, err.Package)
}

// invalid sha256 hash (thrown if remote hash differs from repository hash)
type InvalidChecksum struct {
	Repository string
}

func (err InvalidChecksum) Error() string {
	return fmt.Sprintf("Cannot validate checksum. It is possible somebody has tampered with the file, or that you are the victim of a MITM-attack.\nContact the repository maintainer if you believe this is an error. (%s)", style.Repo.Render(err.Repository))
}
