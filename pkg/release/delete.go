// Package release ...
package release

import (
	"errors"
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

// Delete release.
func (r *Release) Delete(purge bool) error {
	var link string

	if err := r.Load(); err != nil {
		return err
	}

	if !purge && r.IsDefault() {
		return ErrIsDefault
	}

	binaryFile, err := utils.TemplateStringRelease(referenceData.Ref.Spec.File.GetBinary(), &r.Rel)
	if err != nil {
		return err
	}

	if referenceData.Ref.Spec.File.GetLink() == "" {
		link = r.Rel.Spec.GetPath() + "/" + binaryFile
	} else {
		link = r.Rel.Spec.GetPath() + "/" + referenceData.Ref.Spec.File.GetLink()
	}

	file := link + "_" + r.Rel.Spec.GetVersion()

	if err := os.Remove(file); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	logger.StdLog.Debug().Msgf("Remove binary file: %s", file)

	if err := os.Remove(r.VersionFile); err != nil {
		return err
	}

	logger.StdLog.Debug().Msgf("Remove yaml manifest: %s", r.VersionFile)

	if purge {
		if err := os.Remove(link); err != nil {
			return err
		}

		logger.StdLog.Debug().Msgf("Remove default link: %s", link)

		if err := os.Remove(r.DefaultFile); err != nil {
			return err
		}

		logger.StdLog.Debug().Msgf("Remove default file: %s", r.DefaultFile)
	}

	return nil
}
