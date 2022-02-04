package cmd

import (
	"errors"
	"fmt"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/spf13/cobra"
)

var (
	cmdPurge bool
	cleanCmd = &cobra.Command{ //nolint:exhaustivestruct
		Use:   "clean",
		Short: "Remove a specific release version",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			out.JumpLine()
			inst := release.New(conf, args[0], cmdVersion)
			if cmdVersion == "" {
				out.StepTitle(fmt.Sprintf("Cleaning \"%s\"", args[0]))
				out.JumpLine()

				versions, _ := inst.List()
				for _, version := range versions {
					purgeDefault := false
					logger.StdLog.Info().Msgf("Remove files for version %s", version)

					curRel := release.New(conf, args[0], version)
					if err := curRel.Load(); err != nil {
						logger.StdLog.Fatal().Err(err).Msg("")
					}

					if cmdPurge {
						purgeDefault = curRel.IsDefault()
					}

					if err := curRel.Delete(purgeDefault); err != nil {
						if !errors.Is(err, release.ErrIsDefault) {
							logger.StdLog.Fatal().Err(err).Msg("")
						} else {
							logger.StdLog.Warn().Msg("Not deleted as default version")
						}
					} else {
						logger.SuccessLog.Info().Msg("Done")
					}
				}
			} else if inst.IsInstalled() {
				out.StepTitle(fmt.Sprintf("Cleaning version \"%s\" from release \"%s\"", cmdVersion, args[0]))
				out.JumpLine()
				if err := inst.Delete(false); err != nil {
					if !errors.Is(err, release.ErrIsDefault) {
						logger.StdLog.Fatal().Err(err).Msg("")
					} else {
						logger.StdLog.Warn().Msg("Not deleted as default version")
					}
				} else {
					logger.SuccessLog.Info().Msg("Done")
				}
			} else {
				logger.StdLog.Info().Msg("Release not installed")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().StringVarP(&cmdVersion, "version", "v", "", "Release version to clean (only this one)")

	// if err := cleanCmd.MarkPersistentFlagRequired("version"); err != nil {
	// 	logger.StdLog.Fatal().Err(err).Msg("")
	// }

	cleanCmd.PersistentFlags().BoolVarP(
		&cmdPurge,
		"purge",
		"p",
		false,
		"Purge all the binaries of the release (default one included)",
	)
}
