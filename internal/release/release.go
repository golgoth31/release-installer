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
				URL:    "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
				File:   "ri_{{ .Version }}_SHA256SUMS.txt",
				Format: "sha256",
			},
			Available: Available{
				OS: OS{ //nolint:go-lint
					Linux: "linux",
				},
				Arch: Arch{ //nolint:go-lint
					Amd64: "amd64",
				},
			},
		},
	}

	if file != "myself" {
		yamlData.SetConfigType("yaml") // or viper.SetConfigType("YAML")
		yamlData.SetConfigFile(fmt.Sprintf("%s/%s.yaml", releasePath, file))

		if err := yamlData.ReadInConfig(); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		r := &Release{}

		if err := yamlData.Unmarshal(r); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		return r
	}

	return myself
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
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	for _, val := range release {
		out = append(out, val.GetTagName())
	}

	return out
}
