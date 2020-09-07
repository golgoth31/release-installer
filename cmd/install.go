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

var (
	out           output.Output
	installConfig *viper.Viper
)

// installCmd represents the install command.
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rel, err := cmd.Flags().GetString("release")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst := install.NewInstall(rel)
		if err := installConfig.Unmarshal(inst); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
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

	if err := installCmd.MarkPersistentFlagRequired("os"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := installConfig.BindPFlag("spec.os", installCmd.PersistentFlags().Lookup("os")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP("arch", "a", "", "A help for foo")

	if err := installCmd.MarkPersistentFlagRequired("arch"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := installConfig.BindPFlag("spec.arch", installCmd.PersistentFlags().Lookup("arch")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP("version", "v", "", "A help for foo")

	if err := installCmd.MarkPersistentFlagRequired("version"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := installConfig.BindPFlag("spec.version", installCmd.PersistentFlags().Lookup("version")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP("path", "p", "", "Destination to install file in, should set in your \"$PATH\"")

	if err := installCmd.MarkPersistentFlagRequired("path"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := installConfig.BindPFlag("spec.path", installCmd.PersistentFlags().Lookup("path")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP("release", "r", "", "A help for foo")

	if err := installCmd.MarkPersistentFlagRequired("release"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := installConfig.BindPFlag("metadata.release", installCmd.PersistentFlags().Lookup("release")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().BoolP("default", "d", false, "Set this install as default")

	if err := installConfig.BindPFlag("spec.default", installCmd.PersistentFlags().Lookup("default")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
