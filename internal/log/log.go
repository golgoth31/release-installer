// Package log ...
package log

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v3"
	"github.com/rs/zerolog"
)

// Logger ...
var StdLog zerolog.Logger
var StepLog zerolog.Logger

func OkStatus() string {
	return fmt.Sprint(aurora.Green("\u221A")) // √
}
func DebugStatus() string {
	return fmt.Sprint(aurora.White("\u25CC")) // √
}

func WarnStatus() string {
	return fmt.Sprint(aurora.Yellow("\u26A0"))
}

func ErrorStatus() string {
	return fmt.Sprint(aurora.Red("\u274C")) // ×
}

func FatalStatus() string {
	return fmt.Sprint(aurora.Red("\u2620")) // ×
}

func InfoStatus() string {
	return fmt.Sprint(aurora.Blue("\u00BB")) // ×
}

// SetLogger ...
func Initialize() {
	stdOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	stdOutput.FormatLevel = func(i interface{}) string {
		var level string
		switch i {
		case "debug":
			level = DebugStatus()
		case "info":
			level = InfoStatus()
		case "warning":
			level = WarnStatus()
		case "error":
			level = ErrorStatus()
		case "fatal":
			level = FatalStatus()
		}
		return level
	}
	stdOutput.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	StdLog = zerolog.New(stdOutput)

	stepOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	stepOutput.FormatLevel = func(i interface{}) string {
		return ""
	}
	stepOutput.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	StepLog = zerolog.New(stepOutput)

}
