//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type (
	Build mg.Namespace
)

var name = "pakket"
var buildDir = "build"

// var Default = Build

type BuildOptions struct {
	Static bool
}

func build(arch string, opts BuildOptions) error {
	env := map[string]string{
		"GOOS":   "darwin",
		"GOARCH": arch,
	}

	// Build static binary
	if opts.Static {
		env["CGO_ENABLED"] = "0"
	}

	outputName := fmt.Sprintf("pakket-%s", arch)

	fmt.Printf("Building pakket for %s architecture...\n", arch)

	outputPath := path.Join(
		buildDir,
		outputName,
	)

	args := []string{
		"build",
		"-o",
		outputPath,
		"-ldflags",
		"-s -w",
		".",
	}

	return sh.RunWith(
		env, "go", args...,
	)
}

func (Build) Intel() error {
	return build("amd64", BuildOptions{})
}

func (Build) Silicon() error {
	return build("arm64", BuildOptions{})
}

func (Build) All() {
	mg.Deps(Build.Intel, Build.Silicon)
}

func Clean() {
	sh.Rm("build")
}

func Install() {
	os.Rename(fmt.Sprintf("build/pakket-%s", runtime.GOARCH), "/usr/local/bin/pakket")
}
