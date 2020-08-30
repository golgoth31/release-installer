package install

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/golgoth31/release-installer/internal/output"
	"github.com/golgoth31/release-installer/internal/release"
	getter "github.com/hashicorp/go-getter"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/viper"
)

var yamlData *viper.Viper
var releaseData *release.Release
var defaultProgressBar getter.ProgressTracker = &progressBar{}

var out output.Output

func (i *Install) LoadYaml(file string) {
	yamlData.SetConfigType("yaml")
	yamlData.SetConfigFile(file)

	if err := yamlData.ReadInConfig(); err != nil {
		// out.Status(out.FatalStatus(), fmt.Sprintf("Failed to read %s", file))
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
func (i *Install) templates() (
	releaseURL bytes.Buffer,
	releaseFileName bytes.Buffer,
	checksumURL bytes.Buffer,
	checksumFileName bytes.Buffer) {
	// template strings
	treleaseURL := template.Must(template.New("releaseURL").Parse(releaseData.Spec.Url))
	if err := treleaseURL.Execute(&releaseURL, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating release URL")
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	treleaseFileName := template.Must(template.New("releaseFileName").Parse(releaseData.Spec.File.Archive))
	if err := treleaseFileName.Execute(&releaseFileName, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating release file name")
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	tchecksumURL := template.Must(template.New("checksumURL").Parse(releaseData.Spec.Checksum.Url))
	if err := tchecksumURL.Execute(&checksumURL, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating checksum URL")
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	tchecksumFileName := template.Must(template.New("checksumFileName").Parse(releaseData.Spec.Checksum.File))
	if err := tchecksumFileName.Execute(&checksumFileName, i.Spec); err != nil {
		// out.Status(out.FatalStatus(), "Error templating checksum file name")
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	return releaseURL, releaseFileName, checksumURL, checksumFileName
}

func (i *Install) Install() {

	// define getter opts
	opts := []getter.ClientOption{}
	link := i.Spec.Path + "/" + releaseData.Spec.File.Binary
	file := link + "_" + i.Spec.Version

	releaseURL, releaseFileName, checksumURL, checksumFileName := i.templates()

	downURL := fmt.Sprintf(
		"%s/%s",
		releaseURL.String(),
		releaseFileName.String(),
	)
	getterDownURL := fmt.Sprintf(
		"%s?checksum=file:%s/%s",
		downURL,
		checksumURL.String(),
		checksumFileName.String(),
	)

	out.StepTitle("Release files")
	fmt.Println()
	logger.StdLog.Info().Msgf("Checksum file: %s/%s",
		checksumURL.String(),
		checksumFileName.String(),
	)
	logger.StdLog.Info().Msgf("Archive file:  %s", downURL)

	fmt.Println()
	out.StepTitle("Downloading files")
	fmt.Println()

	pwd, err := os.Getwd()
	if err != nil {
		logger.StdLog.Fatal().Msgf("Error getting wd: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Build the client
	opts = append(opts, getter.WithProgress(defaultProgressBar))
	client := &getter.Client{
		Ctx:     ctx,
		Src:     getterDownURL,
		Dst:     file,
		Pwd:     pwd,
		Mode:    getter.ClientModeFile,
		Options: opts,
	}

	if err := client.Get(); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	// ensure the binary is executable
	if err := os.Chmod(file, 0755); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	fmt.Println()
	logger.StdLog.Info().Msgf("File saved as: %s", file)

	if i.Spec.Default {
		logger.StdLog.Info().Msgf("Creating symlink: %s", link)
		_, err := os.Stat(link)
		if err == nil {
			if err := os.Remove(link); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}
		}
		if err := os.Symlink(file, link); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
		logger.StdLog.Info().Msgf("Done")
	}
}
