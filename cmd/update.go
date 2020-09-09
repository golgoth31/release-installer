/*
Copyright © 2020 David Sabatie <david.sabatie@notrenet.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
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
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		// Get the working directory for the repository
		w, err := r.Worktree()
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
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
			logger.StdLog.Fatal().Err(err).Msg("")
		} else {
			logger.SuccessLog.Info().Msg("Done")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
