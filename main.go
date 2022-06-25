package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pakket-project/pakket/cmd"
	"github.com/pakket-project/pakket/config"
)

func main() {
	// Error if not running MacOS
	if runtime.GOOS != "darwin" {
		fmt.Println("You must be on MacOS to run pakket!")
		os.Exit(1)
	}

	if !(runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64") {
		fmt.Println("Unsupported architecture! Pakket only runs on Intel and Apple Silicon based Macs.")
		os.Exit(1)
	}

	config.GetConfig()   // Get config
	config.GetLockfile() // Get lockfile

	cmd.Execute()
}
