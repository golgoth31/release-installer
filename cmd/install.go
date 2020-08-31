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
	"fmt"

	"github.com/golgoth31/release-installer/internal/install"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/output"
	"github.com/logrusorgru/aurora"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var out output.Output
var installConfig *viper.Viper

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rel, _ := cmd.Flags().GetString("release")

		inst := install.NewInstall(rel)
		if err := installConfig.Unmarshal(inst); err != nil {
			logger.StdLog.Fatal().Msg("Failed unmarshal to install")
		}

		fmt.Println()
		out.StepTitle(fmt.Sprintf("Installing release \"%s\"", rel))
		fmt.Println()
		logger.StdLog.Info().Msgf("Version: %s", inst.Spec.Version)
		logger.StdLog.Info().Msgf("OS:      %s", inst.Spec.Os)
		logger.StdLog.Info().Msgf("Arch:    %s", inst.Spec.Arch)
		logger.StdLog.Info().Msgf("Default: %t", inst.Spec.Default)
		logger.StdLog.Info().Msgf("Path:    %s", inst.Spec.Path)
		fmt.Println()

		inst.Install()

		fmt.Println()
		fmt.Printf(" %v", aurora.Bold(logger.OkStatus()))
		out.StepTitle(
			"Release installed",
		)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installConfig = viper.New()

	installCmd.PersistentFlags().StringP("os", "o", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("os")
	installConfig.BindPFlag("spec.os", installCmd.PersistentFlags().Lookup("os"))

	installCmd.PersistentFlags().StringP("arch", "a", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("arch")
	installConfig.BindPFlag("spec.arch", installCmd.PersistentFlags().Lookup("arch"))

	installCmd.PersistentFlags().StringP("version", "v", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("version")
	installConfig.BindPFlag("spec.version", installCmd.PersistentFlags().Lookup("version"))

	installCmd.PersistentFlags().StringP("path", "p", "", "Destination to install file in, should set in your \"$PATH\"")
	installCmd.MarkPersistentFlagRequired("path")
	installConfig.BindPFlag("spec.path", installCmd.PersistentFlags().Lookup("path"))

	installCmd.PersistentFlags().StringP("release", "r", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("release")
	installConfig.BindPFlag("metadata.release", installCmd.PersistentFlags().Lookup("release"))

	installCmd.PersistentFlags().BoolP("default", "d", false, "Set this install as default")
	installConfig.BindPFlag("spec.default", installCmd.PersistentFlags().Lookup("default"))
}
