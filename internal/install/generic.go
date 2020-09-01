package install

import (
	"github.com/golgoth31/release-installer/internal/release"
	"github.com/spf13/viper"
)

// NewInstall ...
func NewInstall(rel string) *Install {
	yamlData = viper.New()
	releaseData = release.New(rel)

	return &Install{APIVersion: "release/v1", Kind: "Install"}
}
