package install

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/release"
	"github.com/spf13/viper"
)

var yamlData *viper.Viper
var releaseData *release.Release

func NewInstall(rel string) *Install {
	yamlData = viper.New()
	releaseData = release.New(rel)
	return &Install{ApiVersion: "release/v1", Kind: "Install"}
}

func (i *Install) LoadYaml(file string) {
	yamlData.SetConfigType("yaml")
	yamlData.SetConfigFile(file)

	if err := yamlData.ReadInConfig(); err != nil {
		log.Logger.Fatal().Err(err).Msgf("Failed to read %s", file)
	}
}

func (i *Install) Download() {
	treleaseURL := template.Must(template.New("releaseURL").Parse(releaseData.Spec.Url))
	var releaseURL bytes.Buffer
	if err := treleaseURL.Execute(&releaseURL, i.Spec); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error templating release URL")
	}
	treleaseFileName := template.Must(template.New("releaseFileName").Parse(releaseData.Spec.File.Name))
	var releaseFileName bytes.Buffer
	if err := treleaseFileName.Execute(&releaseFileName, i.Spec); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error templating release file name")
	}

	downURL := fmt.Sprintf("%s/%s", releaseURL.String(), releaseFileName.String())
	log.Logger.Info().Msgf("Downloading file: %s", downURL)
	resp, err := http.Get(downURL)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Error getting file")
	}
	defer resp.Body.Close()

	// Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Error reading body")
	}

	fmt.Printf("%s\n", body)
}
