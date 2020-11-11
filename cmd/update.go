package cmd

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the releases definitions",
	Run: func(cmd *cobra.Command, args []string) {
		out.StepTitle("Updating releases definitions")
		releasePath := fmt.Sprintf(
			"%s/%s",
			viper.GetString("homedir"),
			viper.GetString("releases.dir"),
		)
		logger.StdLog.Info().Msgf("Cloning from %s", viper.GetString("releases.git"))
		logger.StdLog.Info().Msgf("Cloning in %s", releasePath)
		r, err := git.PlainOpen(releasePath)
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Clone")
		}

		// Get the working directory for the repository
		w, err := r.Worktree()
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Worktree")
		}

		// Pull the latest changes from the origin remote and merge into the current branch
		err = w.Pull(&git.PullOptions{
			Depth:      1,
			RemoteName: "origin",
			Force:      true,
		})
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			logger.SuccessLog.Info().Msg("Already up-to-date")
		} else if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Pull")
		} else {
			logger.SuccessLog.Info().Msg("Done")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
