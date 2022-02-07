package cmd

import (
	"fmt"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/reference"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/spf13/cobra"
)

const (
	defaultVersionToShow = 5
)

var (
	cmdInstalled bool
	cmdNumber    int
	cmdNoFormat  bool
	listCmd      = &cobra.Command{ //nolint:exhaustivestruct
		Use:   "list [release name]",
		Short: "List available releases or versions",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var list []string

			if len(args) == 0 {
				// List all available references
				ref := reference.Reference{} //nolint:exhaustivestruct

				if cmdInstalled {
					out.StepTitle("Installed releases", 1)
				} else {
					out.StepTitle("Available releases", 1)
				}

				files, err := reference.List(conf)
				if err != nil {
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
						out.Success(
							fmt.Sprintf(
								"%s (%s)",
								ref.Ref.Metadata.GetName(),
								defaultVal,
							),
						)
					} else {
						if !cmdInstalled {
							out.Info(ref.Ref.Metadata.GetName())
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
				if cmdInstalled {
					out.StepTitle(
						fmt.Sprintf(
							"Installed versions for release \"%s\"",
							args[0],
						),
						1,
					)

					logger.StdLog.Debug().Msg(rel.InstallDir)
					versions, err := rel.List()
					if err != nil {
						logger.StdLog.Fatal().Err(err).Msg("")
					}

					for _, version := range versions {
						curRel := release.New(conf, args[0], version)
						if err := curRel.Load(); err != nil {
							logger.StdLog.Fatal().Err(err).Msg("")
						}

						if curRel.IsDefault() {
							out.Success(curRel.Rel.Spec.GetVersion())
						} else {
							out.Info(curRel.Rel.Spec.GetVersion())
						}
					}
					out.JumpLine()
				} else {
					if !cmdNoFormat {
						out.StepTitle(fmt.Sprintf("Available versions for release \"%s\"", args[0]), 1)
					}

					ref := reference.New(conf, args[0])

					list = ref.ListVersions(cmdNumber)

					for i := 0; i < len(list); i++ {
						if !cmdNoFormat {
							if defaultVal == list[i] {
								out.Success(list[i])
							} else {
								out.Info(list[i])
							}
						} else {
							out.NoFormat(list[i])
						}
					}
					if !cmdNoFormat {
						out.JumpLine()
					}
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolVarP(&cmdInstalled, "installed", "i", false, "Show installed releases only")
	listCmd.PersistentFlags().IntVarP(
		&cmdNumber,
		"number",
		"n",
		defaultVersionToShow,
		"Number of releases or versions to show",
	)
	listCmd.PersistentFlags().BoolVar(&cmdNoFormat, "noformat", false, "remove formating")
}
