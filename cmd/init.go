package cmd

import (
	"fmt"
	"os"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "init",
	Short: "Initialize release-installer",
	Run: func(cmd *cobra.Command, args []string) {
		out.StepTitle("Initializing release-installer", level1)
		out.StepTitle("Create Directory", level2)

		if _, err := os.Stat(conf.Release.Path); err != nil {
			if err = os.Mkdir(conf.Release.Path, dirPerms); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}
		}

		out.Success("Done")
		out.StepTitle("Downloading references", level2)

		if err := utils.Download(
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("references.repo"),
				"releases/download/latest/ri-releases-definitions.tar.gz",
			),
			conf.Reference.Path,
			true,
		); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		out.JumpLine()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
