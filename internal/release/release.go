// Package release ...
package release

import (
	"context"
	"fmt"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"
)

var yamlData *viper.Viper

// New ...
func New(rel string) *Release {
	yamlData = viper.New()

	return loadYaml(rel)
}

func loadYaml(file string) *Release {
	releasePath := fmt.Sprintf(
		"%s/%s",
		viper.GetString("homedir"),
		viper.GetString("releases.dir"),
	)

	myself := &Release{
		APIVersion: "release/v1",
		Kind:       "Release",
		Metadata: Metadata{
			Name: "release-installer",
			Web:  "https://github.com/golgoth31/release-installer",
		},
		Spec: Spec{
			Repo: Repo{
				Type:  "github",
				Name:  "release-installer",
				Owner: "golgoth31",
			},
			File: File{
				URL:        "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
				Src:        "ri_{{ .Version }}_{{ .Os }}_{{ .Arch }}",
				BinaryPath: ".",
				Binary:     "ri",
				Mode:       "file",
			},
			Checksum: Checksum{
				Check:  true,
				URL:    "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
				File:   "ri_{{ .Version }}_SHA256SUMS.txt",
				Format: "sha256",
			},
			Available: Available{
				Os: Os{
					Linux:  "linux",
					Darwin: "darwin",
				},
				Arch: Arch{
					Amd64: "amd64",
					Arm64: "arm64",
					Arm:   "armv7",
				},
			},
		},
	}

	if file == "myself" {
		return myself
	}

	yamlData.SetConfigType("yaml")
	yamlData.SetConfigFile(fmt.Sprintf("%s/%s.yaml", releasePath, file))

	if err := yamlData.ReadInConfig(); err != nil {
		logger.StdLog.Debug().Err(err).Msg("Unable to read release definition")
		logger.StdLog.Fatal().Msg("Unknown release")
	}

	r := &Release{}

	if err := yamlData.Unmarshal(r); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	return r
}

// ListVersions ...
func (r *Release) ListVersions(num int) []string {
	ctx := context.Background()
	client := github.NewClient(nil)
	opts := &github.ListOptions{
		Page:    1,
		PerPage: num,
	}
	out := []string{}

	release, _, err := client.Repositories.ListReleases(ctx, r.Spec.Repo.Owner, r.Spec.Repo.Name, opts)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to get version list")
	}

	for _, val := range release {
		out = append(out, val.GetTagName())
	}

	return out
}
