// Package release ...
package release

import (
	"fmt"

	logger "github.com/golgoth31/release-installer/internal/log"
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
		},
		Spec: Spec{
			URL: "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
			File: File{
				Src:        "ri-{{ .Os }}-{{ .Arch }}",
				BinaryPath: ".",
				Binary:     "ri",
				Mode:       "file",
			},
			Checksum: Checksum{
				URL:    "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
				File:   "ri_{{ .Version }}_SHA256SUMS",
				Format: "sha256",
			},
			Available: Available{
				OS: OS{
					Linux: "linux",
				},
				Arch: Arch{
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
