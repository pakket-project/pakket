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
		outputName = "pakket-silicon"

		fmt.Println("Building pakket for Silicon architecture...")
	} else if arch == IntelArch {
		outputName = "pakket-intel"

		fmt.Println("Building pakket for Intel architecture...")
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

func Install() {
	if runtime.GOARCH == "arm64" {
		os.Rename("build/pakket-silicon", "/usr/local/bin/pakket")
	} else if runtime.GOARCH == "amd64" {
		os.Rename("build/pakket-intel", "/usr/local/bin/pakket")
	}
}
