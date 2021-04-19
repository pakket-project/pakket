package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/stewproject/stew/cmd"
	"github.com/stewproject/stew/internals/config"
)

func main() {
	// Error if not running MacOS
	if runtime.GOOS != "darwin" {
		fmt.Println("You must be on MacOS to run Stew!")
		os.Exit(1)
	}

	// Check if running on Intel or Apple Silicon
	r, err := syscall.Sysctl("sysctl.proc_translated")
	if err != nil && err.Error() == "no such file or directory" {
		config.GetConfig()

		cmd.Execute()
	}

	if r == "\x00\x00\x00" || r == "\x01\x00\x00" {
		fmt.Println("Looks like you're running Stew on Apple Silicon! We don't support this yet.")
	}
}
