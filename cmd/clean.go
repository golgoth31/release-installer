package cmd

import (
	"fmt"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "clean",
	Short: "Remove a specific release version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rel := args[0]

		out.JumpLine()
		out.StepTitle(fmt.Sprintf("Removing \"%s\"", rel))
		out.JumpLine()

		ver, err := cmd.PersistentFlags().GetString("version")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst := release.New(conf, rel, ver)
		if inst.IsInstalled() {
			def, err := inst.GetDefault()
			if err != nil {
				logger.StdLog.Fatal().Err(err).Msg("No default available")
			}

			if ver == def {
				logger.StdLog.Fatal().Msg("This version is the default, can't be removed")
			}

			inst.Delete()
		} else {
			logger.StdLog.Info().Msg("Release not installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().StringP("version", "v", "", "Release version to clean (only this one)")

	if err := cleanCmd.MarkPersistentFlagRequired("version"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	cleanCmd.PersistentFlags().BoolP("purge", "p", false, "Purge all the binries of the release (default one included)")
}
