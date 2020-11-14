package cmd

import (
	"fmt"

	"github.com/golgoth31/release-installer/internal/release"
	"github.com/spf13/cobra"
)

// releasesCmd represents the releases command.
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "List available release",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rel := release.New(args[0])
		fmt.Printf("%v", rel)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// releasesCmd.PersistentFlags().StringP("version", "v", "", "A help for foo")

	// if err := releasesCmd.MarkPersistentFlagRequired("version"); err != nil {
	// 	logger.StdLog.Fatal().Err(err).Msg("")
	// }
}
