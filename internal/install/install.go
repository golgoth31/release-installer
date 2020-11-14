// Package install ...
package install

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/mitchellh/go-homedir"

	"gopkg.in/yaml.v2"

	"github.com/golgoth31/release-installer/internal/output"
	"github.com/golgoth31/release-installer/internal/progressbar"
	"github.com/golgoth31/release-installer/internal/release"
	getter "github.com/hashicorp/go-getter"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/viper"
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
	treleaseURL := template.Must(template.New("releaseURL").Parse(releaseData.Spec.URL))
	if err := treleaseURL.Execute(&releaseURL, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating release URL")
		revertError = err
	}

	treleaseFileName := template.Must(template.New("releaseFileName").Parse(releaseData.Spec.File.Src))
	if err := treleaseFileName.Execute(&releaseFileName, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating release file name")
		revertError = err
	}

	tchecksumURL := template.Must(template.New("checksumURL").Parse(releaseData.Spec.Checksum.URL))
	if err := tchecksumURL.Execute(&checksumURL, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating checksum URL")
		revertError = err
	}

	tchecksumFileName := template.Must(template.New("checksumFileName").Parse(releaseData.Spec.Checksum.File))
	if err := tchecksumFileName.Execute(&checksumFileName, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating checksum file name")
		revertError = err
	}

	tbinaryPath := template.Must(template.New("binaryPath").Parse(releaseData.Spec.File.BinaryPath))
	if err := tbinaryPath.Execute(&binaryPath, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating checksum file name")
		revertError = err
	}

	return releaseURL, releaseFileName, checksumURL, checksumFileName, binaryPath, revertError
}

func (i *Install) IsInstalled(name string) bool {
	installPath := fmt.Sprintf(
		"%s/%s/%s",
		viper.GetString("homedir"),
		viper.GetString("install.dir"),
		name,
	)

	_, err := os.Stat(installPath)
	if err != nil {
		return false
	}

	return true
}

// func (i *Install) IsDefault(name string) (string, error) {

// }

func (i *Install) GetDefault(name string) (string, error) {
	installed := i.IsInstalled(name)
	if installed {
		installPath := fmt.Sprintf(
			"%s/%s/%s",
			viper.GetString("homedir"),
			viper.GetString("install.dir"),
			name,
		)
		defaultFile := fmt.Sprintf(
			"%s/%s",
			installPath,
			"default",
		)

		data, err := ioutil.ReadFile(defaultFile)
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Reading default file")
		}

		return string(data), nil
	}

	return "", fmt.Errorf("Not installed")
}

func (i *Install) saveConfig() {
	installPath := fmt.Sprintf(
		"%s/%s/%s",
		viper.GetString("homedir"),
		viper.GetString("install.dir"),
		i.Metadata.Release,
	)

	if _, err := os.Stat(installPath); err != nil {
		if err = os.MkdirAll(installPath, 0750); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
	}

	saving, err := yaml.Marshal(i)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	saveData := viper.New()

	saveData.SetConfigType("yaml")

	err = saveData.ReadConfig(bytes.NewBuffer(saving))
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := saveData.WriteConfigAs(
		installPath + "/" + i.Spec.Version + ".yaml",
	); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}

func (i *Install) saveDefault() {
	installPath := fmt.Sprintf(
		"%s/%s/%s",
		viper.GetString("homedir"),
		viper.GetString("install.dir"),
		i.Metadata.Release,
	)

	f, err := os.Create(installPath + "/default")
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
	defer f.Close() //nolint: errcheck,gosec

	_, err = f.WriteString(i.Spec.Version)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}

func (i *Install) removeConfig(revertError error) {
	installPath := fmt.Sprintf(
		"%s/%s/%s",
		viper.GetString("homedir"),
		viper.GetString("install.dir"),
		i.Metadata.Release,
	)

	if err := os.Remove(installPath + "/" + i.Spec.Version + ".yaml"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if revertError != nil {
		logger.StdLog.Fatal().Err(revertError).Msg("")
	}
}

// Install ...
func (i *Install) Install() { //nolint: funlen
	// define getter opts
	var err error

	i.Spec.Path, err = homedir.Expand(i.Spec.Path)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	link := i.Spec.Path + "/" + releaseData.Spec.File.Binary
	file := link + "_" + i.Spec.Version

	i.saveConfig()

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

	if checksumURL.String() != "" && checksumFileName.String() != "" {
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

	fmt.Println()
	logger.SuccessLog.Info().Msgf("File saved as: %s", file)

	if i.Spec.Default {
		fmt.Println()
		logger.StdLog.Info().Msgf("Creating symlink: %s\n", link)

		_, err = os.Stat(link)
		if err == nil {
			if err = os.Remove(link); err != nil {
				i.removeConfig(err)
			}
		}

		if err = os.Symlink(file, link); err != nil {
			i.removeConfig(err)
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

	if err = os.Chmod(dst, 0750); err != nil { //nolint: gosec
		i.removeConfig(err)
	}

	return nil
}
