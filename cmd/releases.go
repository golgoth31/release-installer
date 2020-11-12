package cmd

import (
	"fmt"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/spf13/cobra"
)

// releasesCmd represents the releases command.
var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "List available releases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("releases called")
	},
}

func init() {
	rootCmd.AddCommand(releasesCmd)

	releasesCmd.PersistentFlags().StringP("version", "v", "", "A help for foo")

	if err := releasesCmd.MarkPersistentFlagRequired("version"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
