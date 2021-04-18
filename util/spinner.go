package util

import (
	"time"

	"github.com/theckman/yacspin"
)

var (
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
)
