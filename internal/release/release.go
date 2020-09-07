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
