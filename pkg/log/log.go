// Package log ...
package log

import (
	"os"

	"github.com/logrusorgru/aurora/v3"
	"github.com/rs/zerolog"
)

var (
	// StdLog logger.
	StdLog zerolog.Logger

	// StepLog logger.
	StepLog zerolog.Logger

	// SuccessLog logger.
	SuccessLog zerolog.Logger
)

// OkStatus string.
func OkStatus() string {
	return aurora.Green("\u221A").String() // √
}

// DebugStatus string.
func DebugStatus() string {
	return aurora.White("\u25CC").String() // ◌
}

// WarnStatus string.
func WarnStatus() string {
	return aurora.Yellow("\u26A0").String() // ⚠
}

// ErrorStatus string.
func ErrorStatus() string {
	return aurora.Red("\u2A2F").String() // ⨯
}

// FatalStatus string.
func FatalStatus() string {
	return aurora.Red("\u2620").String() // ☠
}

// InfoStatus string.
func InfoStatus() string {
	return aurora.Blue("\u00BB").String() // »
}

// Initialize logger.
func Initialize() {
	stdOutput := zerolog.ConsoleWriter{Out: os.Stdout} //nolint:exhaustivestruct
	stdOutput.FormatLevel = func(i any) string {
		var level string

		switch i {
		case "debug":
			level = DebugStatus()
		case "info":
			level = InfoStatus()
		case "warn":
			level = WarnStatus()
		case "error":
			level = ErrorStatus()
		case "fatal":
			level = FatalStatus()
		}

		return level
	}
	stdOutput.FormatTimestamp = func(i any) string {
		return ""
	}
	StdLog = zerolog.New(stdOutput)

	stepOutput := zerolog.ConsoleWriter{Out: os.Stdout} //nolint:exhaustivestruct
	stepOutput.FormatLevel = func(i any) string {
		return ""
	}
	stepOutput.FormatTimestamp = func(i any) string {
		return ""
	}
	StepLog = zerolog.New(stepOutput)

	successOutput := zerolog.ConsoleWriter{Out: os.Stdout} //nolint:exhaustivestruct
	successOutput.FormatLevel = func(i any) string {
		return OkStatus()
	}
	successOutput.FormatTimestamp = func(i any) string {
		return ""
	}
	SuccessLog = zerolog.New(successOutput)
}
