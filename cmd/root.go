// Package cmd ...
package cmd

import (
	"fmt"
	"os"
	"strings"

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
	out        = output.New()
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
	level1   int         = 1
	level2   int         = 2
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
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
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			logger.StdLog.Fatal().Err(errHomeDir).Msg("")
		}
		homedir = fmt.Sprintf("%s/.release-installer", home)
		// Search config in home directory with name ".release-installer" (without extension).
		viper.AddConfigPath(homedir)
		viper.SetConfigName("release-installer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if errViperRead := viper.ReadInConfig(); errViperRead == nil {
		logger.StdLog.Debug().Msgf("Reading config file: %s", viper.ConfigFileUsed())
	} else {
		logger.StdLog.Debug().Msg("Config file not found, saving default")

		if errMkdir := os.Mkdir(homedir, dirPerms); errMkdir != nil {
			if os.IsNotExist(errMkdir) {
				logger.StdLog.Fatal().Err(errMkdir).Msg("")
			}
		}
	}

	defaultConfig.SetDefault(homedir)

	conf = internalConfig.Load()

	data, err := os.ReadFile(homedir + "/version")
	if err != nil {
		logger.StdLog.Debug().Err(err).Msg("Reading version file")
	}

	if strings.TrimSpace(string(data)) != configs.Version {
		if err := migration.Migrate(homedir, strings.TrimSpace(string(data)), conf); err != nil {
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

	if err := viper.WriteConfigAs(homedir + "/release-installer.yaml"); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}
