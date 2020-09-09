// Package config ...
package config

import "github.com/spf13/viper"

// SetDefault ...
func SetDefault(homedir string) {
	viper.SetDefault("homedir", homedir)
	viper.SetDefault("releases.dir", "releases")
	viper.SetDefault("releases.git", "https://github.com/golgoth31/release-installer-definitions.git")
	viper.SetDefault("install.dir", "install")
	viper.SetDefault("binaryPath", "~/bin")
	viper.SetDefault("arch", "amd64")
	viper.SetDefault("os", "linux")
}
