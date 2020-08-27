package release

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var yamlData *viper.Viper

func New(rel string) *Release {
	yamlData = viper.New()
	return loadYaml(rel)
}
func loadYaml(file string) *Release {
	yamlData.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	yamlData.SetConfigFile(fmt.Sprintf("releases/%s.yaml", file))

	if err := yamlData.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read %s", file)
	}

	r := &Release{}
	err := yamlData.Unmarshal(r)
	if err != nil {
		log.Fatal("unable to decode into release struct")
	}
	return r
}
