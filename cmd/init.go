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
		out.StepTitle("Initializing release-installer")

		out.JumpLine()
		out.StepTitle("Create Directory")

		if _, err := os.Stat(conf.Release.Path); err != nil {
			if err = os.Mkdir(conf.Release.Path, dirPerms); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}
		}

		out.Success("Done")
		out.JumpLine()
		out.StepTitle("Downloading references")

		if err := utils.Download(
			fmt.Sprintf(
				"%s%s/%s",
				"https://",
				viper.GetString("references.repo"),
				"releases/download/latest/ri-releases-definitions.tar.gz",
			),
			conf.Reference.Path,
			false,
		); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		out.Success("Done")
		out.JumpLine()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
