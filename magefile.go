// +build mage

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type (
	Build   mg.Namespace
	Install mg.Namespace
)

var name = "stew"
var buildDir = "build"

// var Default = Build

type BuildOptions struct {
	Static bool
}

const (
	SiliconArch = "arm64"
	IntelArch   = "amd64"
)

func build(arch string, opts BuildOptions) error {
	env := map[string]string{
		"GOOS":   "darwin",
		"GOARCH": arch,
	}

	// Build static binary
	if opts.Static {
		env["CGO_ENABLED"] = "0"
	}

	var outputName string

	if arch == SiliconArch {
		outputName = "stew-silicon"

		fmt.Println("Building Stew for Silicon architecture...")
	} else if arch == IntelArch {
		outputName = "stew-intel"

		fmt.Println("Building Stew for Intel architecture...")
	}

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
	return build(IntelArch, BuildOptions{})
}

func (Build) Silicon() error {
	return build(SiliconArch, BuildOptions{})
}

func (Build) All() {
	mg.Deps(Build.Intel, Build.Silicon)
}

func Clean() {
	sh.Rm("build")
}

func (Install) Intel() {
	os.Rename("build/stew-intel", "/usr/local/bin/stew")
}
