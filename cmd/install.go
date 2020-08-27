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

		inst.Download()
		fmt.Println(inst.ApiVersion)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	installCmd.PersistentFlags().StringP("os", "o", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("os")
	viper.BindPFlag("spec.os", installCmd.PersistentFlags().Lookup("os"))
	installCmd.PersistentFlags().StringP("arch", "a", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("arch")
	viper.BindPFlag("spec.arch", installCmd.PersistentFlags().Lookup("arch"))
	installCmd.PersistentFlags().StringP("version", "v", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("version")
	viper.BindPFlag("spec.version", installCmd.PersistentFlags().Lookup("version"))
	installCmd.PersistentFlags().StringP("path", "p", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("path")
	viper.BindPFlag("spec.path", installCmd.PersistentFlags().Lookup("path"))
	installCmd.PersistentFlags().StringP("release", "r", "", "A help for foo")
	installCmd.MarkPersistentFlagRequired("release")
	viper.BindPFlag("metadata.release", installCmd.PersistentFlags().Lookup("release"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
