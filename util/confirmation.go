package util

import (
	"fmt"
	"strings"

	"github.com/pakket-project/pakket/util/style"
)

// Destrutive confirm prompt. Default is no.
func DestructiveConfirm(prompt string, always bool) (yes bool, alwaysYes bool) {
	var confirmation string
	var confirmPrompt string

	if always {
		confirmPrompt = fmt.Sprintf("%s [(y)es/(a)lways/(N)o] ", style.Error.Render(prompt))
	} else {
		confirmPrompt = fmt.Sprintf("%s [(y)es/(a)lways/(N)o] ", style.Error.Render(prompt))
	}

	fmt.Print(confirmPrompt)
	fmt.Scanf("%s", &confirmation)
	confirmation = strings.ToLower(confirmation)
	return confirmation == "y" || confirmation == "yes", confirmation == "a" || confirmation == "always"
}

// Confirm prompt. Default is yes.
func Confirm(prompt string) bool {
	var confirmation string

	fmt.Printf("%s [Y/n] ", prompt)
	fmt.Scanf("%s", &confirmation)
	confirmation = strings.ToLower(confirmation)

	return confirmation == "" || confirmation == "y" || confirmation == "yes"
}
