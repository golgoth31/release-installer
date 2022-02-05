package cmd

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/golgoth31/release-installer/configs"
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/reference"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/golgoth31/release-installer/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "update",
	Short: "Update my binary and the releases definitions",
	Run: func(cmd *cobra.Command, args []string) {
		ref := reference.New(conf, "myself")
		list := ref.ListVersions(1)

		out.Info(fmt.Sprintf("Current ri version: %s", configs.Version))
		out.JumpLine()

		if configs.Version != list[0] || cmdForce {
			out.StepTitle(fmt.Sprintf("Updating ri binary from %s to %s", configs.Version, list[0]))
			out.JumpLine()

			inst := release.New(conf, "myself", list[0])
			inst.Rel.Spec.Arch = runtime.GOARCH
			inst.Rel.Spec.Os = runtime.GOOS
			inst.Rel.Spec.Default = true
			inst.Rel.Spec.Path = cmdPath

			inst.Install(cmdForce)

			path, err := homedir.Expand(inst.Rel.Spec.Path)
			if err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}

			if err := syscall.Exec(fmt.Sprintf("%s/ri_%s", path, list[0]), os.Args, os.Environ()); err != nil {
				logger.StdLog.Fatal().Err(err).Msgf("%s", fmt.Sprintf("%s/ri_%s", path, list[0]))
			}
		} else {
			out.StepTitle("No need to update ri binary")
		}

		out.JumpLine()
		out.StepTitle("Updating reference definitions")
		if err := utils.Download(
			fmt.Sprintf(
				"%s/%s",
				conf.RepoURL,
				"releases/download/latest/ri-releases-definitions.tar.gz",
			),
			conf.Reference.Path,
			true,
		); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		out.JumpLine()
		logger.SuccessLog.Info().Msg("Done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(
		&cmdPath,
		"path",
		"p",
		"~/bin",
		"Destination to install binary in, should be set in your \"$PATH\"",
	)

	updateCmd.PersistentFlags().BoolVarP(&cmdForce, "force", "f", false, "Force update")
}
