package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/reference"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultVersionToShow = 5
)

var listCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "list [release name]",
	Short: "List available releases or versions",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var files, list []string

		installed, err := cmd.Flags().GetBool("installed")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		noFormat, err := cmd.Flags().GetBool("noformat")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		number, err := cmd.Flags().GetInt("number")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		yamlData := viper.New()
		yamlData.SetConfigType("yaml")

		if len(args) == 0 {
			// List all available references
			ref := reference.Reference{} //nolint:exhaustivestruct

			out.JumpLine()
			if installed {
				out.StepTitle("Installed releases")
			} else {
				out.StepTitle("Available releases")
			}

			out.JumpLine()

			if err := filepath.Walk(conf.Reference.Path, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}

				return nil
			}); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}

			for _, file := range files {
				logger.StdLog.Debug().Msgf("Reading file: %s", file)
				ref.File = file
				if err := ref.Load(); err != nil {
					logger.StdLog.Fatal().Err(err).Msg("")
				}
				rel := release.New(conf, ref.Ref.Metadata.GetName(), "")
				defaultVal, err := rel.GetDefault()
				if err == nil {
					logger.SuccessLog.Info().Msgf(
						"%s (%s)",
						ref.Ref.Metadata.GetName(),
						defaultVal,
					)
				} else {
					if !installed {
						logger.StdLog.Info().Msg(ref.Ref.Metadata.GetName())
					}
				}
			}
			out.JumpLine()
		} else {
			rel := release.New(conf, args[0], "")
			defaultVal, err := rel.GetDefault()
			if err != nil {
				logger.StdLog.Debug().Err(err).Msgf("Unable to get default file")
			}
			if installed {
				out.JumpLine()
				out.StepTitle(fmt.Sprintf("Installed versions for release \"%s\"", args[0]))
				out.JumpLine()

				logger.StdLog.Debug().Msg(rel.InstallDir)
				_, err := os.Stat(rel.InstallDir)
				if err != nil {
					logger.StdLog.Fatal().Msg("Not installed")
				}
				if err = filepath.Walk(rel.InstallDir, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && info.Name() != "default" {
						logger.StdLog.Debug().Msg(path)
						files = append(files, path)
					}

					return nil
				}); err != nil {
					logger.StdLog.Fatal().Err(err).Msg("")
				}

				for _, file := range files {
					rel.File = file
					if err := rel.Load(); err != nil {
						logger.StdLog.Fatal().Err(err).Msg("")
					}

					if defaultVal == rel.Rel.Spec.GetVersion() {
						logger.SuccessLog.Info().Msg(rel.Rel.Spec.GetVersion())
					} else {
						logger.StdLog.Info().Msg(rel.Rel.Spec.GetVersion())
					}
				}
				out.JumpLine()
			} else {
				if !noFormat {
					out.JumpLine()
					out.StepTitle(fmt.Sprintf("Available versions for release \"%s\"", args[0]))
					out.JumpLine()
				}

				ref := reference.New(conf, args[0])

				list = ref.ListVersions(number)

				for i := 0; i < len(list); i++ {
					if !noFormat {
						if defaultVal == list[i] {
							logger.SuccessLog.Info().Msg(list[i])
						} else {
							logger.StdLog.Info().Msg(list[i])
						}
					} else {
						logger.StepLog.Info().Msg(list[i])
					}
				}
				if !noFormat {
					out.JumpLine()
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("installed", "i", false, "Show installed releases only")
	listCmd.PersistentFlags().IntP("number", "n", defaultVersionToShow, "Number of releases or versions to show")
	listCmd.PersistentFlags().Bool("noformat", false, "remove formating")
}
