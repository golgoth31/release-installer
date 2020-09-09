package cmd

import (
	"fmt"

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
}
