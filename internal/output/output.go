package output

import (
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/logrusorgru/aurora/v3"
)

func (o *Output) StepTitle(str string) {
	logger.StepLog.Error().Msgf(" %v", aurora.Bold(str).Underline())
}
