package release

import (
	"github.com/golgoth31/release-installer/pkg/config"
	common_proto "github.com/golgoth31/release-installer/pkg/proto/common"
	release_proto "github.com/golgoth31/release-installer/pkg/proto/release"
	"github.com/golgoth31/release-installer/pkg/reference"
)

// New is used to create a release instance.
func New(conf *config.Config, name, version string) *Release {
	referenceData = reference.New(conf, name)

	rel := &Release{ //nolint:exhaustivestruct
		Rel: release_proto.Release{
			ApiVersion: conf.Release.APIVersion,
			Kind:       conf.Release.Kind,
			Metadata: &common_proto.Metadata{ //nolint:exhaustivestruct
				Name: name,
			},
			Spec: &release_proto.Spec{ //nolint:exhaustivestruct
				Version: version,
				Binary:  referenceData.Ref.Spec.File.GetBinary(),
			},
		},
	}
	rel.paths(conf)

	return rel
}
