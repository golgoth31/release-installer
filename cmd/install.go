package cmd

import (
	"fmt"
	"runtime"

	"github.com/golgoth31/release-installer/internal/install"
	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/output"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	out           output.Output
	installConfig *viper.Viper
)

// installCmd represents the install command.
var installCmd = &cobra.Command{ //nolint:go-lint
	Use:   "install [release]",
	Short: "Install one release",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rel := args[0]

		force, err := cmd.PersistentFlags().GetBool("force")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		// installConfig.Set("metadata.release", rel)
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

		inst.Install(force)

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

	installCmd.PersistentFlags().StringP("os", "o", runtime.GOOS, "Release binary OS")

	if err := installConfig.BindPFlag("spec.os", installCmd.PersistentFlags().Lookup("os")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP("arch", "a", runtime.GOARCH, "Release binary architecture")

	if err := installConfig.BindPFlag("spec.arch", installCmd.PersistentFlags().Lookup("arch")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP("version", "v", "", "Release version")

	if err := installCmd.MarkPersistentFlagRequired("version"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	if err := installConfig.BindPFlag("spec.version", installCmd.PersistentFlags().Lookup("version")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().StringP(
		"path",
		"p",
		"~/bin",
		"Destination to install binary in, should be set in your \"$PATH\"",
	)

	if err := installConfig.BindPFlag("spec.path", installCmd.PersistentFlags().Lookup("path")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().BoolP("default", "d", false, "Set this install as default")

	if err := installConfig.BindPFlag("spec.default", installCmd.PersistentFlags().Lookup("default")); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	installCmd.PersistentFlags().BoolP("force", "f", false, "Force release install")
}
