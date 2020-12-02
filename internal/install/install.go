// Package install ...
package install

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/Masterminds/sprig/v3"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/output"
	"github.com/golgoth31/release-installer/internal/progressbar"
	"github.com/golgoth31/release-installer/internal/release"
	getter "github.com/hashicorp/go-getter"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	releaseData        *release.Release
	defaultProgressBar getter.ProgressTracker = &progressbar.ProgressBar{}
	out                output.Output
)

func (i *Install) templates() (
	releaseURL bytes.Buffer,
	releaseFileName bytes.Buffer,
	checksumURL bytes.Buffer,
	checksumFileName bytes.Buffer,
	binaryPath bytes.Buffer,
	revertError error) {
	revertError = nil
	// template strings
	treleaseURL := template.Must(template.New("releaseURL").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.URL))
	if err := treleaseURL.Execute(&releaseURL, i.Spec); err != nil {
		revertError = err
	}

	treleaseFileName := template.Must(template.New("releaseFileName").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.Src))
	if err := treleaseFileName.Execute(&releaseFileName, i.Spec); err != nil {
		revertError = err
	}

	tchecksumURL := template.Must(template.New("checksumURL").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.Checksum.URL))
	if err := tchecksumURL.Execute(&checksumURL, i.Spec); err != nil {
		revertError = err
	}

	tchecksumFileName := template.Must(template.New("checksumFileName").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.Checksum.File))
	if err := tchecksumFileName.Execute(&checksumFileName, i.Spec); err != nil {
		revertError = err
	}

	tbinaryPath := template.Must(template.New("binaryPath").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.BinaryPath))
	if err := tbinaryPath.Execute(&binaryPath, i.Spec); err != nil {
		revertError = err
	}

	return releaseURL, releaseFileName, checksumURL, checksumFileName, binaryPath, revertError
}

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

	_, err := os.Stat(versionPath)
	if err != nil {
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

		return "", err
	}

	logger.StdLog.Debug().Msgf("default data: %s", data)

	return string(data), nil
}

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

func (i *Install) SaveConfig() {
	installDir, versionFile, _ := i.Paths()

	if _, err := os.Stat(installDir); err != nil {
		if err = os.MkdirAll(installDir, 0o750); err != nil {
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
	defer f.Close() //nolint: errcheck,gosec

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
	i.Get()
	_, versionFile, _ := i.Paths()

	link := i.Spec.Path + "/" + releaseData.Spec.File.Binary
	file := link + "_" + i.Spec.Version

	logger.StdLog.Debug().Msgf("Remove binary: %s", file)

	if err := os.Remove(file); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	logger.StdLog.Debug().Msgf("Remove yaml: %s", versionFile)

	if err := os.Remove(versionFile); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}

// Install ...
func (i *Install) Install(force bool) { //nolint:go-lint
	// define getter opts
	var err error

	i.Spec.Arch = i.setRealValues("Arch")
	i.Spec.Os = i.setRealValues("Os")

	logger.StdLog.Debug().Msgf("Release Arch: %s", i.Spec.Arch)
	logger.StdLog.Debug().Msgf("Release Os: %s", i.Spec.Os)

	i.Spec.Path, err = homedir.Expand(i.Spec.Path)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	link := i.Spec.Path + "/" + releaseData.Spec.File.Binary
	file := link + "_" + i.Spec.Version

	if !i.IsInstalled() || force {
		releaseURL, releaseFileName, checksumURL, checksumFileName, binaryPath, revertError := i.templates()
		if revertError != nil {
			i.removeConfig(revertError)
		}

		var srcFile string

		switch releaseData.Spec.File.Mode {
		case "file":
			srcFile = releaseFileName.String()
		case "archive":
			srcFile = releaseData.Spec.File.Binary
		default:
			i.removeConfig(fmt.Errorf("Unknown release mode"))
		}

		downURL := fmt.Sprintf(
			"%s/%s",
			releaseURL.String(),
			releaseFileName.String(),
		)
		getterDownURL := downURL

		out.StepTitle("Release files")
		fmt.Println()

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

		fmt.Println()
		out.StepTitle("Downloading files")
		fmt.Println()

		pwd, err := os.Getwd()
		if err != nil {
			i.removeConfig(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Build the client
		opts := []getter.ClientOption{}
		opts = append(opts, getter.WithProgress(defaultProgressBar))
		client := &getter.Client{ //nolint:go-lint
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

		fmt.Println()
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
		fmt.Println()
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
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	if err = out.Sync(); err != nil {
		return err
	}

	if err = os.Chmod(dst, 0o750); err != nil { //nolint: gosec
		i.removeConfig(err)
	}

	return nil
}
