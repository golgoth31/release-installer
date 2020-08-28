package install

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/golgoth31/release-installer/internal/release"
	getter "github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var yamlData *viper.Viper
var releaseData *release.Release
var defaultProgressBar getter.ProgressTracker = &progressBar{}

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

func (i *Install) Install() {

	// define getter opts
	opts := []getter.ClientOption{}
	link := i.Spec.Path + "/" + releaseData.Spec.File.Binary
	file := link + "_" + i.Spec.Version

	// template strings
	treleaseURL := template.Must(template.New("releaseURL").Parse(releaseData.Spec.Url))
	var breleaseURL bytes.Buffer
	if err := treleaseURL.Execute(&breleaseURL, i.Spec); err != nil {
		log.Fatal("Error templating release URL")
	}
	treleaseFileName := template.Must(template.New("releaseFileName").Parse(releaseData.Spec.File.Archive))
	var breleaseFileName bytes.Buffer
	if err := treleaseFileName.Execute(&breleaseFileName, i.Spec); err != nil {
		log.Fatal("Error templating release file name")
	}
	tchecksumURL := template.Must(template.New("checksumURL").Parse(releaseData.Spec.Checksum.Url))
	var bchecksumURL bytes.Buffer
	if err := tchecksumURL.Execute(&bchecksumURL, i.Spec); err != nil {
		log.Fatal("Error templating checksum URL")
	}
	tchecksumFileName := template.Must(template.New("checksumFileName").Parse(releaseData.Spec.Checksum.File))
	var bchecksumFileName bytes.Buffer
	if err := tchecksumFileName.Execute(&bchecksumFileName, i.Spec); err != nil {
		log.Fatal("Error templating checksum file name")
	}

	downURL := fmt.Sprintf(
		"%s/%s",
		breleaseURL.String(),
		breleaseFileName.String(),
	)
	getterDownURL := fmt.Sprintf(
		"%s?checksum=file:%s/%s",
		downURL,
		bchecksumURL.String(),
		bchecksumFileName.String(),
	)

	log.Infof("Downloading archive file")
	fmt.Println()
	log.Infof("Archive file: %s", downURL)
	log.Infof("Checksum file: %s", fmt.Sprintf(
		"%s/%s",
		bchecksumURL.String(),
		bchecksumFileName.String(),
	))

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %v", err)
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
		log.Fatal(err)
	}

	// ensure the binary is executable
	if err := os.Chmod(file, 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	log.Info("Done")
	log.Infof("File saved as: %s", file)
	fmt.Println()

	if i.Spec.Default {
		log.Infof("Creating symlink: %s", link)
		if err := os.Remove(link); err != nil {
			log.Fatal(err)
		}
		if err := os.Symlink(file, link); err != nil {
			log.Fatal(err)
		}
		log.Info("Done")
	}
}

func (cpb *ProgressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	pb := mpb.New(mpb.WithWidth(60))
	// Parameters of th new progress bar
	bar := pb.AddBar(totalSize,
		mpb.PrependDecorators(
			decor.Name(src),
			decor.Name(" "),
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

	reader := bar.ProxyReader(stream)

	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			pb.Wait()
			return nil
		},
	}
}

func (c *readCloser) Close() error { return c.close() }
