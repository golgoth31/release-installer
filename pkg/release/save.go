// Package release ...
package release

import (
	"fmt"
	"os"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"sigs.k8s.io/yaml"
)

// SaveConfig saves configuration of installed release
func (r *Release) SaveConfig() {
	if _, err := os.Stat(r.InstallDir); err != nil {
		if err = os.MkdirAll(r.InstallDir, dirPerms); err != nil {
			logger.StdLog.Fatal().Err(err).Msgf("Unable to create directory: %s", r.InstallDir)
		}
	}

	if err := r.Write(); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}
}

func (r *Release) saveDefault() {
	f, err := os.Create(r.DefaultFile)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to create file")
	}

	defer func() {
		if ferr := f.Close(); ferr != nil {
			logger.StdLog.Fatal().Err(ferr).Msg("Failed to close file")
		}
	}()

	_, err = f.WriteString(r.Rel.Spec.GetVersion())
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to write file")
	}
}

func (r *Release) Write() error {
	out, errOut := yaml.Marshal(&r.Rel)
	if errOut != nil {
		return fmt.Errorf("%w", errOut)
	}

	if err := os.WriteFile(r.VersionFile, out, 0600); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
