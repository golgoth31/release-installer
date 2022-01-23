package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/golgoth31/release-installer/configs"
	"github.com/golgoth31/release-installer/pkg/install"
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "update",
	Short: "Update my binary and the releases definitions",
	Run: func(cmd *cobra.Command, args []string) {
		rel := release.New("myself")
		list := rel.ListVersions(1)

		force, err := cmd.PersistentFlags().GetBool("force")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		if configs.Version != list[0] || force {
			out.StepTitle("Updating ri binary")
			out.JumpLine()

			path, errFlag := cmd.PersistentFlags().GetString("path")
			if errFlag != nil {
				logger.StdLog.Fatal().Err(errFlag).Msg("")
			}

			inst := install.NewInstall("myself")
			inst.Spec.Arch = runtime.GOARCH
			inst.Spec.Os = runtime.GOOS
			inst.Spec.Path = path
			inst.Spec.Version = list[0]
			inst.Spec.Default = true
			inst.Install(force)
		} else {
			out.StepTitle("No need to update ri binary")
		}

		out.JumpLine()
		out.StepTitle("Updating releases definitions")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		releasePath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("releases.dir"),
		)

		// Build the client
		pwd, err := os.Getwd()
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		opts := []getter.ClientOption{}
		client := &getter.Client{ //nolint:exhaustivestruct
			Ctx: ctx,
			Src: "https://github.com/golgoth31/release-installer-definitions/" +
				"releases/download/latest/ri-releases-definitions.tar.gz",
			Dst:     releasePath,
			Pwd:     pwd,
			Mode:    getter.ClientModeAny,
			Options: opts,
		}

		if err = client.Get(); err != nil {
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
