// Package output ...
package output

import (
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/logrusorgru/aurora/v3"
)

// StepTitle ...
func (o *Output) StepTitle(str string) {
	logger.StepLog.Info().Msgf(" %v", aurora.Bold(str).Underline())
}

// Success ...
func (o *Output) Success(str string) {
	logger.StepLog.Info().Msgf(" %v %v", aurora.Bold(logger.OkStatus()), aurora.Bold(str).Underline())
}
