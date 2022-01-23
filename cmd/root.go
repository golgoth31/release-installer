// Package cmd ...
package cmd

import (
	"fmt"
	"os"

	"github.com/golgoth31/release-installer/internal/config"
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/output"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	out           output.Output
	installConfig *viper.Viper
)

var rootCmd = &cobra.Command{ //nolint:exhaustivestruct
	Use:   "ri [command]",
	Short: "A tool to download and install binaries",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.StdLog.Err(err).Msg("")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.release-installer/release-installer.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var homedir string

	logLevel, err := rootCmd.Flags().GetBool("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	if logLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	logger.Initialize()
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
		homedir = fmt.Sprintf("%s/.release-installer", home)
		// Search config in home directory with name ".release-installer" (without extension).
		viper.AddConfigPath(homedir)
		viper.SetConfigName("release-installer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		logger.StdLog.Debug().Msgf("Reading config file: %s", viper.ConfigFileUsed())
	} else {
		logger.StdLog.Debug().Msg("Config file not found, saving default")

		if err := os.Mkdir(homedir, 0750); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
	}

	config.SetDefault(homedir)

	if err := viper.WriteConfigAs(homedir + "/release-installer.yaml"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
