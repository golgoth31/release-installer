package cmd

import (
	"fmt"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/golgoth31/release-installer/internal/release"
	"github.com/spf13/cobra"
)

// releasesCmd represents the releases command.
var releaseCmd = &cobra.Command{ //nolint:go-lint
	Use:   "release [release name]",
	Short: "List available version of a release",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		number, err := cmd.Flags().GetInt("number")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
		rel := release.New(args[0])

		fmt.Println()
		out.StepTitle(fmt.Sprintf("Available versions for release \"%s\"", rel.Metadata.Name))
		fmt.Println()

		list := rel.ListVersions(number)
		for i := 0; i < len(list); i++ {
			logger.StdLog.Info().Msg(list[i])
		}

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.PersistentFlags().IntP("number", "n", 5, "Number of versions to show")
}
