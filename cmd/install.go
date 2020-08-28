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
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		if err := viper.Unmarshal(inst); err != nil {
			log.Fatal("Failed unmarshal to install")
		}
		log.Debugf("Default value: %v", viper.GetBool("default"))
		log.Debugf("Default value: %v", viper.GetBool("spec.default"))

		log.Infof("Installing release \"%s\" in \"%s\" with values:", rel, inst.Spec.Path)
		fmt.Printf("Version: %s\n", inst.Spec.Version)
		fmt.Printf("OS: %s\n", inst.Spec.Os)
		fmt.Printf("Arch: %s\n", inst.Spec.Arch)
		fmt.Printf("Default: %v\n", inst.Spec.Default)
		fmt.Println()

		inst.Install()

		log.Info("Release installed")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringP("os", "o", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("os")
	viper.BindPFlag("spec.os", installCmd.PersistentFlags().Lookup("os"))

	installCmd.PersistentFlags().StringP("arch", "a", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("arch")
	viper.BindPFlag("spec.arch", installCmd.PersistentFlags().Lookup("arch"))

	installCmd.PersistentFlags().StringP("version", "v", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("version")
	viper.BindPFlag("spec.version", installCmd.PersistentFlags().Lookup("version"))

	installCmd.PersistentFlags().StringP("path", "p", "", "Destination to install file in, should set in your \"$PATH\"")
	installCmd.MarkPersistentFlagRequired("path")
	viper.BindPFlag("spec.path", installCmd.PersistentFlags().Lookup("path"))

	installCmd.PersistentFlags().StringP("release", "r", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("release")
	viper.BindPFlag("metadata.release", installCmd.PersistentFlags().Lookup("release"))

	installCmd.PersistentFlags().BoolP("default", "d", false, "Set this install as default")
	viper.BindPFlag("spec.default", installCmd.PersistentFlags().Lookup("default"))
}
