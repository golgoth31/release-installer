package cmd

import (
	"context"
	"fmt"
	"os"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/progressbar"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultProgressBar getter.ProgressTracker = &progressbar.ProgressBar{}
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize release-installer",
	Run: func(cmd *cobra.Command, args []string) {
		out.StepTitle("Initializing release-installer")

		installPath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("install.dir"),
		)

		if _, err := os.Stat(installPath); err != nil {
			if err = os.Mkdir(installPath, 0750); err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}
		}
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
		client := &getter.Client{
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

		logger.SuccessLog.Info().Msg("Done")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
