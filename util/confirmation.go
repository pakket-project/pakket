package util

import (
	"fmt"
	"strings"

	"github.com/pakket-project/pakket/util/style"
)

// Destrutive confirm prompt. Default is no.
func DestructiveConfirm(prompt string) bool {
	var confirmation string

	fmt.Printf("%s [y/N] ", style.Error.Render(prompt))
	fmt.Scanf("%s", &confirmation)
	confirmation = strings.ToLower(confirmation)

	return confirmation == "y" || confirmation == "yes"
}

// Confirm prompt. Default is yes.
func Confirm(prompt string) bool {
	var confirmation string

	fmt.Printf("%s [Y/n] ", prompt)
	fmt.Scanf("%s", &confirmation)
	confirmation = strings.ToLower(confirmation)

	return confirmation == "" || confirmation == "y" || confirmation == "yes"
}
