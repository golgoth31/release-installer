// Package install ...
package install

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/output"
	"github.com/golgoth31/release-installer/pkg/progressbar"
	"github.com/golgoth31/release-installer/pkg/release"
	getter "github.com/hashicorp/go-getter"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	releaseData        *release.Release
	defaultProgressBar getter.ProgressTracker = &progressbar.ProgressBar{}
	out                output.Output
	errUnknown         = errors.New("Unknown release mode")
)

const (
	dirPerms os.FileMode = 0750
)

// SetArch ...
func (i *Install) setRealValues(field string) (val string) {
	var relAvailable, instSpec interface{}

	switch field {
	case "Arch":
		relAvailable = releaseData.Spec.Available.Arch
		instSpec = i.Spec.Arch
	case "Os":
		relAvailable = releaseData.Spec.Available.Os
		instSpec = i.Spec.Os
	}

	lowerVal := strings.Title(strings.ToLower(fmt.Sprintf("%v", instSpec)))
	realData := reflect.ValueOf(relAvailable).FieldByName(lowerVal)

	if fmt.Sprintf("%s", realData) != "" {
		val = fmt.Sprintf("%s", realData)
	}

	return val
}

// IsInstalled checks if a release is installed.
func (i *Install) IsInstalled() bool {
	_, versionPath, _ := i.Paths()

	logger.StdLog.Debug().Msgf("Release install file: %s", versionPath)

	if _, err := os.Stat(versionPath); err != nil {
		return false
	}

	return true
}

// GetDefault gets default version installed.
func (i *Install) GetDefault() (string, error) {
	_, _, defaultFile := i.Paths()

	data, err := ioutil.ReadFile(defaultFile)
	if err != nil {
		logger.StdLog.Debug().Err(err).Msg("Reading default file")

		return "", fmt.Errorf("%w", err)
	}

	logger.StdLog.Debug().Msgf("default data: %s", data)

	return string(data), nil
}

// Paths returns various path.
func (i *Install) Paths() (installDir string, versionFile string, defaultFile string) {
	installDir = fmt.Sprintf(
		"%s/%s/%s",
		viper.GetString("homedir"),
		viper.GetString("install.dir"),
		i.Metadata.Release,
	)

	versionFile = fmt.Sprintf(
		"%s/%s.yaml",
		installDir,
		i.Spec.Version,
	)

	defaultFile = fmt.Sprintf(
		"%s/%s",
		installDir,
		"default",
	)

	return
}

// SaveConfig saves configuration of installed release
func (i *Install) SaveConfig() {
	installDir, versionFile, _ := i.Paths()

	if _, err := os.Stat(installDir); err != nil {
		if err = os.MkdirAll(installDir, dirPerms); err != nil {
			logger.StdLog.Fatal().Err(err).Msgf("Unable to create directory: %s", installDir)
		}
	}

	saving, err := yaml.Marshal(i)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	saveData := viper.New()

	saveData.SetConfigType("yaml")

	if err := saveData.ReadConfig(bytes.NewBuffer(saving)); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := saveData.WriteConfigAs(versionFile); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}

func (i *Install) saveDefault() {
	_, _, defaultFile := i.Paths()

	f, err := os.Create(defaultFile)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to create file")
	}

	defer func() {
		if ferr := f.Close(); ferr != nil {
			logger.StdLog.Fatal().Err(ferr).Msg("Failed to close file")
		}
	}()

	_, err = f.WriteString(i.Spec.Version)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to write file")
	}
}

func (i *Install) removeConfig(revertError error) {
	_, versionFile, _ := i.Paths()

	if err := os.Remove(versionFile); err != nil {
		logger.StdLog.Debug().Err(err).Msg("")
	}

	if revertError != nil {
		logger.StdLog.Fatal().Err(revertError).Msg("")
	}
}

// Get ...
func (i *Install) Get() {
	_, versionFile, _ := i.Paths()

	vip := viper.New()
	vip.SetConfigFile(versionFile)

	if err := vip.ReadInConfig(); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := vip.Unmarshal(i); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.StdLog.Debug().Msgf("Binary path: %s", i.Spec.Path)
}

// Delete ...
func (i *Install) Delete() {
	var link string

	i.Get()
	_, versionFile, _ := i.Paths()

	_, _, _, _, _, binaryFile, revertError := i.templates()
	if revertError != nil {
		i.removeConfig(revertError)
	}

	if releaseData.Spec.File.Link == "" {
		link = i.Spec.Path + "/" + binaryFile.String()
	} else {
		link = i.Spec.Path + "/" + releaseData.Spec.File.Link
	}

	file := link + "_" + i.Spec.Version

	if err := os.Remove(file); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.SuccessLog.Info().Msgf("Remove binary file: %s", file)

	if err := os.Remove(versionFile); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.SuccessLog.Info().Msgf("Remove yaml manifest: %s", versionFile)
}

// Install ...
func (i *Install) Install(force bool) { //nolint:go-lint
	// define getter opts
	var err error

	var link string

	i.Spec.Arch = i.setRealValues("Arch")
	i.Spec.Os = i.setRealValues("Os")

	logger.StdLog.Debug().Msgf("Release Arch: %s", i.Spec.Arch)
	logger.StdLog.Debug().Msgf("Release Os: %s", i.Spec.Os)

	i.Spec.Path, err = homedir.Expand(i.Spec.Path)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	releaseURL, releaseFileName, checksumURL, checksumFileName, binaryPath, binaryFile, revertError := i.templates()
	if revertError != nil {
		i.removeConfig(revertError)
	}

	if releaseData.Spec.File.Link == "" {
		link = i.Spec.Path + "/" + binaryFile.String()
	} else {
		link = i.Spec.Path + "/" + releaseData.Spec.File.Link
	}

	file := link + "_" + i.Spec.Version

	if !i.IsInstalled() || force {
		var (
			srcFile string
		)

		switch releaseData.Spec.File.Mode {
		case "file":
			srcFile = releaseFileName.String()
		case "archive":
			srcFile = binaryFile.String()
		default:
			i.removeConfig(fmt.Errorf("%w", errUnknown))
		}

		downURL := fmt.Sprintf(
			"%s/%s",
			releaseURL.String(),
			releaseFileName.String(),
		)
		getterDownURL := downURL

		out.StepTitle("Release files")
		out.JumpLine()

		if releaseData.Spec.Checksum.Check {
			getterDownURL = fmt.Sprintf(
				"%s?checksum=file:%s/%s",
				downURL,
				checksumURL.String(),
				checksumFileName.String(),
			)

			logger.StdLog.Info().Msgf("Checksum file: %s/%s",
				checksumURL.String(),
				checksumFileName.String(),
			)
		}

		logger.StdLog.Info().Msgf("Archive file:  %s", downURL)

		out.JumpLine()
		out.StepTitle("Downloading files")
		out.JumpLine()

		pwd, osErr := os.Getwd()
		if osErr != nil {
			i.removeConfig(osErr)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Build the client
		opts := []getter.ClientOption{}
		opts = append(opts, getter.WithProgress(defaultProgressBar))
		client := &getter.Client{
			Ctx:     ctx,
			Src:     getterDownURL,
			Dst:     "/tmp",
			Pwd:     pwd,
			Mode:    getter.ClientModeAny,
			Options: opts,
		}

		if err = client.Get(); err != nil {
			i.removeConfig(err)
		}

		// Move binary file to requested path
		if err = i.moveFile(
			fmt.Sprintf("/tmp/%s/%s", binaryPath.String(), srcFile),
			file,
		); err != nil {
			i.removeConfig(err)
		}

		out.JumpLine()
		logger.SuccessLog.Info().Msgf("File saved as: %s", file)
	} else {
		out.StepTitle("This version is already installed")
	}

	defaultVer, err := i.GetDefault()
	if err != nil {
		logger.StdLog.Debug().Msgf("No default for release: %s\n", i.Metadata.Release)
		i.Spec.Default = true
	} else {
		if defaultVer != i.Spec.Version {
			curDefInst := NewInstall(i.Metadata.Release)
			curDefInst.Spec.Version = defaultVer
			curDefInst.Get()
			curDefInst.Spec.Default = false
			curDefInst.SaveConfig()
		}
	}

	i.SaveConfig()

	if i.Spec.Default {
		out.JumpLine()
		logger.StdLog.Info().Msgf("Creating symlink: %s\n", link)

		_, err := os.Stat(link)
		if err == nil {
			if err = os.Remove(link); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("Unable to remove symlink")
			}
		} else {
			logger.StdLog.Debug().Msgf("file not found: %s\n", i.Metadata.Release)
		}

		if err = os.Symlink(file, link); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Unable to create symlink")
		}

		i.saveDefault()

		logger.SuccessLog.Info().Msgf("Done")
	}
}

func (i *Install) moveFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		if ferr := in.Close(); ferr != nil {
			logger.StdLog.Fatal().Err(ferr).Msg("Failed to close file")
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		if ferr := out.Close(); ferr != nil {
			logger.StdLog.Fatal().Err(ferr).Msg("Failed to close file")
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = out.Sync(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = os.Chmod(dst, dirPerms); err != nil {
		i.removeConfig(err)
	}

	if err = os.Remove(src); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
