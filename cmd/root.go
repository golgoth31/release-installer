// Package cmd ...
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golgoth31/release-installer/configs"
	internalConfig "github.com/golgoth31/release-installer/internal/config"
	defaultConfig "github.com/golgoth31/release-installer/internal/default"
	"github.com/golgoth31/release-installer/internal/migration"
	"github.com/golgoth31/release-installer/pkg/config"
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/output"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	out        output.Output
	conf       *config.Config
	cmdVersion string
	cmdForce   bool
	cmdPath    string
	rootCmd    = &cobra.Command{ //nolint:exhaustivestruct
		Use:   "ri [command]",
		Short: "A tool to download and install binaries",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

const (
	dirPerms os.FileMode = 0750
)

// Execute executes the root command.
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

	// if err := viper.BindPFlag("cfgFile", rootCmd.PersistentFlags().Lookup("config")); err != nil {
	// 	log.Fatal().Err(err).Msg("")
	// }
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

		if err := os.Mkdir(homedir, dirPerms); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("")
		}
	}

	defaultConfig.SetDefault(homedir)

	if err := viper.WriteConfigAs(homedir + "/release-installer.yaml"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	conf = internalConfig.Load()

	data, err := ioutil.ReadFile(homedir + "/version")
	if err != nil {
		logger.StdLog.Debug().Err(err).Msg("Reading version file")
	}

	if string(data) != configs.Version {
		if err := migration.Migrate(homedir, string(data), conf); err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Unable to migrate")
		}

		f, err := os.Create(homedir + "/version")
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Unable to create version file")
		}

		defer func() {
			if ferr := f.Close(); ferr != nil {
				logger.StdLog.Fatal().Err(ferr).Msg("Failed to close version file")
			}
		}()

		_, err = f.WriteString(configs.Version)
		if err != nil {
			logger.StdLog.Fatal().Err(err).Msg("Unable to write version file")
		}
	}

}
