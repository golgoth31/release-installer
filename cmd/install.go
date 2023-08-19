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

			out.StepTitle(fmt.Sprintf("Installing %q", args[0]), 1)
			out.Info("Requested:")
			out.Info(fmt.Sprintf("  Version: %s", inst.Rel.Spec.GetVersion()))
			out.Info(fmt.Sprintf("  OS:      %s", inst.Rel.Spec.GetOs()))
			out.Info(fmt.Sprintf("  Arch:    %s", inst.Rel.Spec.GetArch()))
			out.Info(fmt.Sprintf("  Default: %t", inst.Rel.Spec.GetDefault()))
			out.Info(fmt.Sprintf("  Path:    %s", inst.Rel.Spec.GetPath()))

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
