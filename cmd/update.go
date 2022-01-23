package cmd

import (
	"fmt"

	"github.com/golgoth31/release-installer/configs"
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/reference"
	"github.com/golgoth31/release-installer/pkg/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "update",
	Short: "Update my binary and the releases definitions",
	Run: func(cmd *cobra.Command, args []string) {
		ref := reference.New(conf, "myself")
		list := ref.ListVersions(1)

		force, err := cmd.PersistentFlags().GetBool("force")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		if configs.Version != list[0] || force {
			out.StepTitle("Updating ri binary")
			out.JumpLine()

			// path, errFlag := cmd.PersistentFlags().GetString("path")
			// if errFlag != nil {
			// 	logger.StdLog.Fatal().Err(errFlag).Msg("")
			// }

			// inst := release.NewInstall("myself")
			// inst.Spec.Arch = runtime.GOARCH
			// inst.Spec.Os = runtime.GOOS
			// inst.Spec.Path = path
			// inst.Spec.Version = list[0]
			// inst.Spec.Default = true
			// inst.Install(force)
		} else {
			out.StepTitle("No need to update ri binary")
		}

		out.JumpLine()
		out.StepTitle("Updating releases definitions")
		if err := utils.Download(
			fmt.Sprintf(
				"%s/%s",
				conf.RepoURL,
				"releases/download/latest/ri-releases-definitions.tar.gz",
			),
			conf.Reference.Path,
			false,
		); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		out.JumpLine()
		logger.SuccessLog.Info().Msg("Done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringP(
		"path",
		"p",
		"~/bin",
		"Destination to install binary in, should be set in your \"$PATH\"",
	)

	updateCmd.PersistentFlags().BoolP("force", "f", false, "Force update")
}
