// Package release ...
package release

import (
	"errors"
	"fmt"
	"os"
	"strings"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/output"
	"github.com/golgoth31/release-installer/pkg/reference"
	"github.com/golgoth31/release-installer/pkg/utils"
	"github.com/mitchellh/go-homedir"
)

var (
	referenceData   *reference.Reference
	out             output.Output
	errUnknownMode  = errors.New("unknown release mode")
	errUnknownArch  = errors.New("unknown arch")
	errUnknownOs    = errors.New("unknown os")
	errUnknownField = errors.New("unknown field")
)

const (
	dirPerms  os.FileMode = 0750
	filePerms os.FileMode = 0600
)

// setRealValues extract the arch name as given by the reference into standard one
func (r *Release) setRealValues(field string) (string, error) {
	switch field {
	case "Arch":
		switch strings.ToLower(r.Rel.Spec.GetArch()) {
		case "amd64":
			return referenceData.Ref.Spec.Available.Arch.GetAmd64(), nil
		case "arm64":
			return referenceData.Ref.Spec.Available.Arch.GetArm64(), nil
		case "arm":
			return referenceData.Ref.Spec.Available.Arch.GetArm(), nil
		case "386":
			return referenceData.Ref.Spec.Available.Arch.GetI386(), nil
		default:
			logger.StdLog.Debug().Msgf("Release install file: %s", r.VersionFile)

			return "", fmt.Errorf("%w", errUnknownArch)
		}
	case "Os":
		switch strings.ToLower(r.Rel.Spec.GetOs()) {
		case "linux":
			return referenceData.Ref.Spec.Available.Os.GetLinux(), nil
		case "windows":
			return referenceData.Ref.Spec.Available.Os.GetWindows(), nil
		case "darwin":
			return referenceData.Ref.Spec.Available.Os.GetDarwin(), nil
		default:
			return "", fmt.Errorf("%w", errUnknownOs)
		}
	}

	return "", fmt.Errorf("%w", errUnknownField)
}

// IsInstalled checks if a release is installed.
func (r *Release) IsInstalled() bool {
	logger.StdLog.Debug().Msgf("Release install file: %s", r.VersionFile)

	if _, err := os.Stat(r.VersionFile); err != nil {
		return false
	}

	return true
}

// Install ...
func (r *Release) Install(force bool) { //nolint:go-lint
	var (
		err  error
		link string
	)

	r.Rel.Spec.Arch, err = r.setRealValues("Arch")
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	r.Rel.Spec.Os, err = r.setRealValues("Os")
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.StdLog.Debug().Msgf("Release Arch: %s", r.Rel.Spec.Arch)
	logger.StdLog.Debug().Msgf("Release Os: %s", r.Rel.Spec.Os)

	r.Rel.Spec.Path, err = homedir.Expand(r.Rel.Spec.Path)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	releaseURL, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.File.GetUrl(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	releaseFileName, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.File.GetSrc(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	checksumURL, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.Checksum.GetUrl(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	checksumFileName, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.Checksum.GetFile(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	binaryPath, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.File.GetBinaryPath(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	binaryFile, revertError := utils.TemplateStringRelease(referenceData.Ref.Spec.File.GetBinary(), &r.Rel)
	if revertError != nil {
		r.removeConfig(revertError)
	}

	if referenceData.Ref.Spec.File.GetLink() == "" {
		link = r.Rel.Spec.Path + "/" + binaryFile
	} else {
		link = r.Rel.Spec.Path + "/" + referenceData.Ref.Spec.File.GetLink()
	}

	file := link + "_" + r.Rel.Spec.GetVersion()

	if !r.IsInstalled() || force {
		var srcFile string

		switch referenceData.Ref.Spec.File.GetMode() {
		case "file":
			srcFile = releaseFileName
		case "archive":
			srcFile = binaryFile
		default:
			r.removeConfig(fmt.Errorf("%w", errUnknownMode))
		}

		downURL := fmt.Sprintf(
			"%s/%s",
			releaseURL,
			releaseFileName,
		)
		getterDownURL := downURL

		out.StepTitle("Release files")
		out.JumpLine()

		if referenceData.Ref.Spec.Checksum.GetCheck() {
			getterDownURL = fmt.Sprintf(
				"%s?checksum=file:%s/%s",
				downURL,
				checksumURL,
				checksumFileName,
			)

			logger.StdLog.Info().Msgf("Checksum file: %s/%s",
				checksumURL,
				checksumFileName,
			)
		}

		logger.StdLog.Info().Msgf("Archive file:  %s", downURL)

		out.JumpLine()
		out.StepTitle("Downloading files")
		out.JumpLine()

		if errDownload := utils.Download(
			getterDownURL,
			"/tmp",
			true,
		); errDownload != nil {
			r.removeConfig(errDownload)
		}

		// Move binary file to requested path
		if err = utils.MoveFile(
			fmt.Sprintf("/tmp/%s/%s", binaryPath, srcFile),
			file,
			dirPerms,
		); err != nil {
			r.removeConfig(err)
		}

		out.JumpLine()
		logger.SuccessLog.Info().Msgf("File saved as: %s", file)
	} else {
		out.StepTitle("This version is already installed")
	}

	_, err = r.GetDefault()
	if err != nil {
		logger.StdLog.Debug().Msgf("No default for release: %s\n", r.Rel.Metadata.GetName())
		r.Rel.Spec.Default = true
		// This should be done only if the current version is declared as default
		// } else {
		// 	if defaultVer != r.Rel.Spec.Version {
		// 		curDefInst := NewInstall(i.Metadata.Release)
		// 		curDefInst.Spec.Version = defaultVer
		// 		curDefInst.Get()
		// 		curDefInst.Spec.Default = false
		// 		curDefInst.SaveConfig()
		// 	}
	}

	r.SaveConfig()

	if r.Rel.Spec.GetDefault() {
		out.JumpLine()
		logger.StdLog.Info().Msgf("Creating symlink: %s\n", link)

		_, err := os.Stat(link)
		if err == nil {
			if err = os.Remove(link); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("Unable to remove symlink")
			}
		} else {
			logger.StdLog.Debug().Msgf("file not found: %s\n", r.Rel.Metadata.GetName())
		}

		if err = os.Symlink(file, link); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Unable to create symlink")
		}

		r.saveDefault()

		logger.SuccessLog.Info().Msgf("Done")
	}
}
