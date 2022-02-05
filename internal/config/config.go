// Package config convert some viper config to internal struct.
package config

import (
	"fmt"

	"github.com/golgoth31/release-installer/pkg/config"
	"github.com/spf13/viper"
)

// Load loads config data from viper.
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
		RepoURL: viper.GetString("references.repo"),
	}

	return &c
}
