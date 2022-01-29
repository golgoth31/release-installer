// Package log ...
package log

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v3"
	"github.com/rs/zerolog"
)

var (
	// StdLog ...
	StdLog zerolog.Logger

	// StepLog ...
	StepLog zerolog.Logger

	// SuccessLog ...
	SuccessLog zerolog.Logger
)

// OkStatus ...
func OkStatus() string {
	return fmt.Sprint(aurora.Green("\u221A")) // √
}

// DebugStatus ...
func DebugStatus() string {
	return fmt.Sprint(aurora.White("\u25CC")) // √
}

// WarnStatus ...
func WarnStatus() string {
	return fmt.Sprint(aurora.Yellow("\u26A0"))
}

// ErrorStatus ...
func ErrorStatus() string {
	return fmt.Sprint(aurora.Red("\u274C")) // ×
}

// FatalStatus ...
func FatalStatus() string {
	return fmt.Sprint(aurora.Red("\u2620")) // ×
}

// InfoStatus ...
func InfoStatus() string {
	return fmt.Sprint(aurora.Blue("\u00BB")) // ×
}

// Initialize ...
func Initialize() {
	stdOutput := zerolog.ConsoleWriter{Out: os.Stdout} //nolint:exhaustivestruct
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

	stepOutput := zerolog.ConsoleWriter{Out: os.Stdout} //nolint:exhaustivestruct
	stepOutput.FormatLevel = func(i interface{}) string {
		return ""
	}
	stepOutput.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	StepLog = zerolog.New(stepOutput)

	successOutput := zerolog.ConsoleWriter{Out: os.Stdout} //nolint:exhaustivestruct
	successOutput.FormatLevel = func(i interface{}) string {
		return OkStatus()
	}
	successOutput.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	SuccessLog = zerolog.New(successOutput)
}
