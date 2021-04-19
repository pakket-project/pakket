package util

import (
	"fmt"
	"time"

	"github.com/theckman/yacspin"
)

var (
	// Spinner config
	SpinnerConf = yacspin.Config{
		Frequency:         50 * time.Millisecond,
		HideCursor:        true,
		ColorAll:          false,
		CharSet:           yacspin.CharSets[14],
		Suffix:            " ",
		SuffixAutoColon:   false,
		StopCharacter:     "✓",
		StopFailCharacter: "✗",
		StopColors:        []string{"fgGreen"},
		StopFailColors:    []string{"fgRed"},
		Colors:            []string{"fgCyan"},
	}
	// Spinner = &yacspin.Spinner{}
)

// Stops spinner to print a message. Uses Println
func PrintSpinnerMsg(s *yacspin.Spinner, msg string) {
	s.StopCharacter("")
	s.Stop()
	fmt.Println(msg)
	s.Start()
	s.StopCharacter(SpinnerConf.StopCharacter)
}
