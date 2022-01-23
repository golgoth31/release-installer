// Package defaultConfig ...
package defaultConfig

import (
	"github.com/spf13/viper"
)

// SetDefault ...
func SetDefault(homedir string) {
	viper.SetDefault("homedir", homedir)
	viper.SetDefault("references.dir", "references")
	viper.SetDefault("references.repo", "github.com/golgoth31/release-installer-definitions")
	viper.SetDefault("releases.dir", "releases")
	viper.SetDefault("binary.dir", "~/bin")
}
