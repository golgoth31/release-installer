// Package defaultconfig ...
package defaultconfig

import (
	"github.com/spf13/viper"
)

// SetDefault config values.
func SetDefault(homedir string) {
	viper.SetDefault("homedir", homedir)
	viper.SetDefault("references.dir", "references")
	viper.SetDefault("references.repo", "https://github.com/golgoth31/release-installer-definitions")
	viper.SetDefault("repo", "https://github.com/golgoth31/release-installer")
	viper.SetDefault("releases.dir", "releases")
	viper.SetDefault("binary.dir", "~/bin")
}
