package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/golgoth31/release-installer/configs"
	"github.com/golgoth31/release-installer/internal/release"

	"github.com/golgoth31/release-installer/internal/install"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command.
var updateCmd = &cobra.Command{ //nolint:go-lint
	Use:   "update",
	Short: "Update the releases definitions",
	Run: func(cmd *cobra.Command, args []string) {
		rel := release.New("myself")
		list := rel.ListVersions(1)

		if configs.Version == list[0] {
			out.StepTitle("No need to update ri binary")
		} else {
			out.StepTitle("Updating ri binary")
			fmt.Println()

			path, err := cmd.PersistentFlags().GetString("path")
			if err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}

			inst := install.NewInstall("myself")
			inst.Spec.Arch = runtime.GOARCH
			inst.Spec.Os = runtime.GOOS
			inst.Spec.Path = path
			inst.Spec.Version = list[0]
			inst.Spec.Default = true
			inst.Install()
		}

		fmt.Println()
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
		// opts = append(opts, getter.WithProgress(defaultProgressBar))
		client := &getter.Client{ //nolint:go-lint
			Ctx:     ctx,
			Src:     "https://github.com/golgoth31/release-installer-definitions/releases/download/latest/ri-releases-definitions.tar.gz",
			Dst:     releasePath,
			Pwd:     pwd,
			Mode:    getter.ClientModeAny,
			Options: opts,
		}

		if err = client.Get(); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		fmt.Println()
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
}
