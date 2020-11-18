package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golgoth31/release-installer/internal/install"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// releasesCmd represents the releases command.
var releasesCmd = &cobra.Command{ //nolint:go-lint
	Use:   "releases",
	Short: "List available releases",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installed, err := cmd.Flags().GetBool("installed")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		fmt.Println()
		if installed {
			out.StepTitle("Installed releases")
		} else {
			out.StepTitle("Available releases")
		}

		fmt.Println()

		var files []string
		releasePath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("releases.dir"),
		)

		inst := install.Install{}
		yamlData := viper.New()
		yamlData.SetConfigType("yaml") // or viper.SetConfigType("YAML")

		err = filepath.Walk(releasePath, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}

			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			yamlData.SetConfigFile(file)

			if err := yamlData.ReadInConfig(); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}
			inst.Metadata.Release = yamlData.GetString("metadata.name")
			defaultVal, err := inst.GetDefault()
			if err == nil {
				logger.SuccessLog.Info().Msgf(
					"%s (%s)",
					yamlData.GetString("metadata.name"),
					defaultVal,
				)
			} else {
				if !installed {
					logger.StdLog.Info().Msg(yamlData.GetString("metadata.name"))
				}
			}
		}

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(releasesCmd)

	releasesCmd.PersistentFlags().BoolP("installed", "i", false, "A help for foo")
}
