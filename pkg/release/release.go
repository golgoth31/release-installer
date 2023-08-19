package release

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golgoth31/release-installer/pkg/config"
	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

// Paths returns various path.
func (r *Release) paths(conf *config.Config) {
	r.InstallDir = fmt.Sprintf(
		"%s/%s",
		conf.Release.Path,
		r.Rel.Metadata.GetName(),
	)

	r.VersionFile = fmt.Sprintf(
		"%s/%s.yaml",
		r.InstallDir,
		r.Rel.Spec.GetVersion(),
	)

	r.DefaultFile = fmt.Sprintf(
		"%s/%s",
		r.InstallDir,
		"default",
	)
}

// GetDefault gets default version installed.
func (r *Release) GetDefault() (string, error) {
	data, err := os.ReadFile(r.DefaultFile)
	if err != nil {
		logger.StdLog.Debug().Err(err).Msg("Reading default file")

		return "", fmt.Errorf("%w", err)
	}

	logger.StdLog.Debug().Msgf("default data: %s", data)

	return string(data), nil
}

// IsDefault checks if current release is the default installed.
func (r *Release) IsDefault() bool {
	def, _ := r.GetDefault()

	return def == r.Rel.Spec.GetVersion()
}

// Load release data from yaml manifest.
func (r *Release) Load() error {
	jsonData, err := utils.Load(r.VersionFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if errUnmarshall := protojson.Unmarshal(jsonData, &r.Rel); errUnmarshall != nil {
		return fmt.Errorf("%w", errUnmarshall)
	}

	return nil
}

func (r *Release) List() ([]string, error) {
	var verions []string

	if _, err := os.Stat(r.InstallDir); err != nil {
		return []string{}, fmt.Errorf("release not installed")
	}

	if err := filepath.Walk(r.InstallDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != "default" {
			version := strings.Split(info.Name(), ".yaml")
			logger.StdLog.Debug().Msg(version[0])
			verions = append(verions, version[0])
		}

		return nil
	}); err != nil {
		return []string{}, err
	}

	return verions, nil
}
