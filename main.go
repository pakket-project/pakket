package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/stewproject/stew/cmd"
	"github.com/stewproject/stew/internals/config"
)

func main() {
	// Error if not running MacOS
	if runtime.GOOS != "darwin" {
		fmt.Println("You must be on MacOS to run Stew!")
		os.Exit(1)
	}

	// Error if running Apple Silicon
	if runtime.GOARCH == "arm64" {
		fmt.Println("Apple Silicon is not yet supported.")
		os.Exit(1)
	}

	config.GetConfig()

	cmd.Execute()
}
