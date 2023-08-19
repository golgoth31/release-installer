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
			inst := release.New(conf, args[0], cmdVersion)
			if cmdVersion == "" {
				out.StepTitle(fmt.Sprintf("Cleaning %q", args[0]), level1)

				versions, _ := inst.List()
				for _, version := range versions {
					purgeDefault := false
					out.Info(
						fmt.Sprintf(
							"Remove files for version %s",
							version,
						),
					)

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
							out.Warn("Not deleted as default version")
						}
					} else {
						out.Success("Done")
					}
				}
			} else if inst.IsInstalled() {
				out.StepTitle(
					fmt.Sprintf(
						"Cleaning version %q from release %q",
						cmdVersion,
						args[0],
					),
					level1,
				)
				if err := inst.Delete(false); err != nil {
					if !errors.Is(err, release.ErrIsDefault) {
						logger.StdLog.Fatal().Err(err).Msg("")
					} else {
						out.Warn("Not deleted as default version")
					}
				} else {
					out.Success("Done")
				}
			} else {
				out.Info("Release not installed")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().StringVarP(&cmdVersion, "version", "v", "", "Release version to clean (only this one)")

	cleanCmd.PersistentFlags().BoolVarP(
		&cmdPurge,
		"purge",
		"p",
		false,
		"Purge all the binaries of the release (default one included)",
	)
}
