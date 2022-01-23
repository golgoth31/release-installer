package cmd

import (
	"fmt"
	"runtime"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var installConfig *viper.Viper

var installCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "install [release]",
	Short: "Install one release",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.PersistentFlags().GetBool("force")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		version, err := cmd.PersistentFlags().GetString("version")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst := release.New(conf, args[0], version)

		inst.Rel.Spec.Arch, err = cmd.PersistentFlags().GetString("arch")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst.Rel.Spec.Os, err = cmd.PersistentFlags().GetString("os")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst.Rel.Spec.Path, err = cmd.PersistentFlags().GetString("path")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		inst.Rel.Spec.Default, err = cmd.PersistentFlags().GetBool("default")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}

		out.JumpLine()
		out.StepTitle(fmt.Sprintf("Installing \"%s\"", args[0]))
		out.JumpLine()
		logger.StdLog.Info().Msg("Requested:")
		logger.StdLog.Info().Msgf("  Version: %s", inst.Rel.Spec.GetVersion())
		logger.StdLog.Info().Msgf("  OS:      %s", inst.Rel.Spec.GetOs())
		logger.StdLog.Info().Msgf("  Arch:    %s", inst.Rel.Spec.GetArch())
		logger.StdLog.Info().Msgf("  Default: %t", inst.Rel.Spec.GetDefault())
		logger.StdLog.Info().Msgf("  Path:    %s", inst.Rel.Spec.GetPath())
		out.JumpLine()

		inst.Install(force)

		out.JumpLine()
		out.Success(
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

	// if err := installConfig.BindPFlag("spec.version", installCmd.PersistentFlags().Lookup("version")); err != nil {
	// 	logger.StdLog.Fatal().Err(err).Msg("")
	// }

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
