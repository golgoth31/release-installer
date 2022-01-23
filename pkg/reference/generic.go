package reference

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golgoth31/release-installer/pkg/config"
	logger "github.com/golgoth31/release-installer/pkg/log"
	common_proto "github.com/golgoth31/release-installer/pkg/proto/common"
	reference_proto "github.com/golgoth31/release-installer/pkg/proto/reference"
)

// New is used to create a reference instance.
func New(conf *config.Config, name string) *Reference {
	ref := &Reference{} //nolint:exhaustivestruct

	if name == "myself" {
		ref.Ref = conf.SetMyself()

		return ref
	}

	ref.File = fmt.Sprintf("%s/%s.yaml", conf.Reference.Path, name)
	ref.Ref = reference_proto.Reference{ //nolint:exhaustivestruct
		Metadata: &common_proto.Metadata{ //nolint:exhaustivestruct
			Name: name,
		},
	}

	if err := ref.Load(); err != nil {
		logger.StdLog.Fatal().Err(err).Msgf("Unable to load data for reference: %s", name)
	}

	return ref
}

func List(conf *config.Config) ([]string, error) {
	var files []string

	if err := filepath.Walk(conf.Reference.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	}); err != nil {
		return []string{}, err
	}

	return files, nil
}
