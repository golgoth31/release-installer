package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golgoth31/release-installer/internal/install"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/release"
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
		releasePath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("releases.dir"),
		)
		installPath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("install.dir"),
		)
		// inst := install.Install{}
		yamlData := viper.New()
		yamlData.SetConfigType("yaml")

		if len(args) == 0 {
			inst := install.Install{} //nolint:exhaustivestruct
			logger.JumpLine()
			if installed {
				out.StepTitle("Installed releases")
			} else {
				out.StepTitle("Available releases")
			}

			logger.JumpLine()

			if err := filepath.Walk(releasePath, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}

				return nil
			}); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
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
			logger.JumpLine()
		} else {
			if installed {
				logger.JumpLine()
				out.StepTitle(fmt.Sprintf("Installed versions for release \"%s\"", args[0]))
				logger.JumpLine()

				inst := install.NewInstall(args[0])
				instRelPath := fmt.Sprintf(
					"%s/%s",
					installPath,
					args[0],
				)
				logger.StdLog.Debug().Msg(instRelPath)
				_, err := os.Stat(instRelPath)
				if err != nil {
					logger.StdLog.Fatal().Msg("Not installed")
				}
				if err = filepath.Walk(instRelPath, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && info.Name() != "default" {
						logger.StdLog.Debug().Msg(path)
						files = append(files, path)
					}

					return nil
				}); err != nil {
					logger.StdLog.Fatal().Err(err).Msg("")
				}

				for _, file := range files {
					yamlData.SetConfigFile(file)

					if err := yamlData.ReadInConfig(); err != nil {
						logger.StdLog.Fatal().Err(err).Msgf("Unable to read file: %s", file)
					}

					if err := yamlData.Unmarshal(inst); err != nil {
						logger.StdLog.Fatal().Err(err).Msg("Unable to load yaml data")
					}

					defaultVal, err := inst.GetDefault()
					if err != nil {
						logger.StdLog.Fatal().Err(err).Msgf("Unable to get default file")
					}

					if defaultVal == inst.Spec.Version {
						logger.SuccessLog.Info().Msg(inst.Spec.Version)
					} else {
						logger.StdLog.Info().Msg(inst.Spec.Version)
					}
				}
				logger.JumpLine()
			} else {
				inst := install.NewInstall(args[0])
				rel := release.New(args[0])
				if !noFormat {
					logger.JumpLine()
					out.StepTitle(fmt.Sprintf("Available versions for release \"%s\"", rel.Metadata.Name))
					logger.JumpLine()
				}

				list = rel.ListVersions(number)

				for i := 0; i < len(list); i++ {
					defaultVal, err := inst.GetDefault()
					if err != nil {
						logger.StdLog.Debug().Err(err).Msgf("Unable to get default file")
					}

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
					logger.JumpLine()
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
