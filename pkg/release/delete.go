// Package release ...
package release

import (
	"os"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/utils"
)

func (r *Release) removeConfig(revertError error) {
	if err := os.Remove(r.VersionFile); err != nil {
		logger.StdLog.Debug().Err(err).Msg("")
	}

	if revertError != nil {
		logger.StdLog.Fatal().Err(revertError).Msg("")
	}
}

// Delete ...
func (r *Release) Delete() {
	var link string

	if err := r.Load(); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	binaryFile, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.File.GetBinary(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	if referenceData.Ref.Spec.File.GetLink() == "" {
		link = r.Rel.Spec.GetPath() + "/" + binaryFile
	} else {
		link = r.Rel.Spec.GetPath() + "/" + referenceData.Ref.Spec.File.GetLink()
	}

	file := link + "_" + r.Rel.Spec.GetVersion()

	if err := os.Remove(file); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.SuccessLog.Info().Msgf("Remove binary file: %s", file)

	if err := os.Remove(r.VersionFile); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.SuccessLog.Info().Msgf("Remove yaml manifest: %s", r.VersionFile)
}
