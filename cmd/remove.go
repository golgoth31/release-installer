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
	"github.com/golgoth31/release-installer/internal/install"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a specific release version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rel := args[0]

		ver, err := cmd.PersistentFlags().GetString("version")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst := install.NewInstall(rel)
		inst.Spec.Version = ver
		if inst.IsInstalled() {
			def, err := inst.GetDefault()
			if err != nil {
				logger.StdLog.Fatal().Err(err).Msg("")
			}

			if ver == def {
				logger.StdLog.Fatal().Msg("This version is the default, can't be removed")
			}

			inst.Delete()
		} else {
			logger.StdLog.Fatal().Msg("Release not installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.PersistentFlags().StringP("version", "v", "", "Release version")

	if err := removeCmd.MarkPersistentFlagRequired("version"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
