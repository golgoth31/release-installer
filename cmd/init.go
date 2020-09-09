package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		opt := &git.CloneOptions{
			Depth: 1,
			URL:   viper.GetString("releases.git"),
		}
		releasePath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("releases.dir"),
		)
		logger.StdLog.Info().Msgf("Cloning from %s", viper.GetString("releases.git"))
		logger.StdLog.Info().Msgf("Cloning in %s", releasePath)
		_, err := git.PlainClone(releasePath, false, opt)
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
		logger.SuccessLog.Info().Msg("Done")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
