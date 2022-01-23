package cmd

import (
	"fmt"
	"runtime"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/release"
	"github.com/spf13/cobra"
)

var (
	cmdArch    string
	cmdOs      string
	cmdDefault bool
	installCmd = &cobra.Command{ //nolint:exhaustivestruct
		Use:   "install [release]",
		Short: "Install one release",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			inst := release.New(conf, args[0], cmdVersion)
			inst.Rel.Spec.Arch = cmdArch
			inst.Rel.Spec.Os = cmdOs
			inst.Rel.Spec.Path = cmdPath
			inst.Rel.Spec.Default = cmdDefault

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

			inst.Install(cmdForce)

			out.JumpLine()
			out.Success(
				"Release installed",
			)
		},
	}
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringVarP(&cmdOs, "os", "o", runtime.GOOS, "Release binary OS")
	installCmd.PersistentFlags().StringVarP(&cmdArch, "arch", "a", runtime.GOARCH, "Release binary architecture")
	installCmd.PersistentFlags().StringVarP(&cmdVersion, "version", "v", "", "Release version")
	installCmd.PersistentFlags().StringVarP(
		&cmdPath,
		"path",
		"p",
		"~/bin",
		"Destination to install binary in, should be set in your \"$PATH\"",
	)
	installCmd.PersistentFlags().BoolVarP(&cmdDefault, "default", "d", false, "Set this install as default")
	installCmd.PersistentFlags().BoolVarP(&cmdForce, "force", "f", false, "Force release install")

	if err := installCmd.MarkPersistentFlagRequired("version"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
