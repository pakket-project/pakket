package errors

import (
	"fmt"
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

// invalid checksum (thrown if remote hash differs from repository hash)
type InvalidChecksum struct {
	Mirror string
}

func (err InvalidChecksum) Error() string {
	return fmt.Sprintf("Cannot validate checksum. It is possible somebody has tampered with the file on the mirror, or that you are the victim of a MITM-attack.\nMirror: p%s", err.Mirror)
}
