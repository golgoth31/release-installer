package install

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/golgoth31/release-installer/internal/release"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
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
		log.Fatalf("Failed to read %s", file)
	}
}

func (i *Install) Download() {
	treleaseURL := template.Must(template.New("releaseURL").Parse(releaseData.Spec.Url))
	var releaseURL bytes.Buffer
	if err := treleaseURL.Execute(&releaseURL, i.Spec); err != nil {
		log.Fatal("Error templating release URL")
	}
	treleaseFileName := template.Must(template.New("releaseFileName").Parse(releaseData.Spec.File.Name))
	var releaseFileName bytes.Buffer
	if err := treleaseFileName.Execute(&releaseFileName, i.Spec); err != nil {
		log.Fatal("Error templating release file name")
	}

	downURL := fmt.Sprintf("%s/%s", releaseURL.String(), releaseFileName.String())
	log.Infof("Downloading file: %s", downURL)
	resp, err := http.Get(downURL)
	if err != nil {
		log.Fatal("Error getting file")
	}
	defer resp.Body.Close()

	file := "/tmp/" + releaseFileName.String()

	if err = fileDownload(resp, file); err != nil {
		log.Fatal(err)
	}
}

func fileDownload(resp *http.Response, fullpath string) error {
	// Create new progress bar
	pb := mpb.New(mpb.WithWidth(60))

	// Create file
	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parameters of th new progress bar
	bar := pb.AddBar(resp.ContentLength,
		mpb.PrependDecorators(
			decor.CountersKiloByte("% .2f / % .2f "),
			decor.AverageSpeed(decor.UnitKB, "(% .2f)"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.Name(" - "),
			decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 4}),
			decor.Name(" - "),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
			),
		),
	)

	// Update progress bar while writing file
	_, err = io.Copy(file, bar.ProxyReader(resp.Body))
	if err != nil {
		return err
	}

	pb.Wait()

	return nil
}
