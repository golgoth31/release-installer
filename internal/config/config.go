package config

import (
	"fmt"

	"github.com/golgoth31/release-installer/pkg/config"
	"github.com/spf13/viper"
)

func Load() *config.Config {
	c := config.Config{
		Release: &config.Release{
			Path: fmt.Sprintf(
				"%s/%s",
				viper.GetString("homedir"),
				viper.GetString("releases.dir"),
			),
			APIVersion: "release/v1",
			Kind:       "Release",
		},
		Reference: &config.Reference{
			Path: fmt.Sprintf(
				"%s/%s",
				viper.GetString("homedir"),
				viper.GetString("references.dir"),
			),
		},
		RepoURL: fmt.Sprintf(
			"https://%s",
			viper.GetString("references.repo"),
		),
	}

	return &c
}
