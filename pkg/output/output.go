// Package output ...
package output

import (
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/logrusorgru/aurora/v3"
)

// StepTitle ...
func (o *Output) StepTitle(str string) {
	logger.StepLog.Info().Msgf(" %v", aurora.Bold(str).Underline())
}

// SuccessTitle ...
func (o *Output) SuccessTitle(str string) {
	logger.SuccessLog.Info().Msgf("%v", aurora.Bold(str).Underline())
}

// JumpLine ...
func (o *Output) JumpLine() {
	logger.StepLog.Info().Msg("")
}

// NoFormat ...
func (o *Output) NoFormat(str string) {
	logger.StepLog.Info().Msg(str)
}

// Info ...
func (o *Output) Info(str string) {
	logger.StdLog.Info().Msg(str)
}

// Success ...
func (o *Output) Success(str string) {
	logger.SuccessLog.Info().Msgf(str)
}
