package release

import (
	"fmt"
	"io/ioutil"

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
	data, err := ioutil.ReadFile(r.DefaultFile)
	if err != nil {
		logger.StdLog.Debug().Err(err).Msg("Reading default file")

		return "", fmt.Errorf("%w", err)
	}

	logger.StdLog.Debug().Msgf("default data: %s", data)

	return string(data), nil
}

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
