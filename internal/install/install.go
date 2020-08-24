package install

import (
	"fmt"

	"github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/viper"
)

var yamlData *viper.Viper

func NewInstall() Install {
	yamlData = viper.New()
	return Install{}
}

// func (i *Install) GenerateYaml() {
// 	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

// 	viper.ReadConfig(bytes.NewBuffer(yamlExample))

// 	viper.Get("name") // this would be "steve"
// }
func (i *Install) LoadYaml() {
	yamlData.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	yamlData.SetConfigFile("exemple.yaml")

	if err := yamlData.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", yamlData.ConfigFileUsed())
	}
	log.Logger.Info().Msg(yamlData.GetString("metadata.name"))
	yamlData.SafeWriteConfigAs("test.yaml")
}
