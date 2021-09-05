package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/stewproject/stew/cmd"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util"
)

func main() {
	// Error if not running MacOS
	if runtime.GOOS != "darwin" {
		fmt.Println("You must be on MacOS to run Stew!")
		os.Exit(1)
	}

	if runtime.GOARCH == "arm64" {
		util.Arch = "silicon"
	} else if runtime.GOARCH == "amd64" {
		util.Arch = "intel"
	} else {
		fmt.Println("Unsupported architecture! Stew only runs on Intel and Apple Silicon based Macs.")
		os.Exit(1)
	}

	config.GetConfig() // Get config

	cmd.Execute()
}
