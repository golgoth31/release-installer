package release

import (
	"github.com/golgoth31/release-installer/pkg/config"
	common_proto "github.com/golgoth31/release-installer/pkg/proto/common"
	release_proto "github.com/golgoth31/release-installer/pkg/proto/release"
	"github.com/golgoth31/release-installer/pkg/reference"
)

func New(conf *config.Config, name string, version string) *Release {
	referenceData = reference.New(conf, name)

	rel := &Release{
		Rel: release_proto.Release{
			ApiVersion: conf.Release.APIVersion,
			Kind:       conf.Release.Kind,
			Metadata: &common_proto.Metadata{
				Name: name,
			},
			Spec: &release_proto.Spec{
				Version: version,
			},
		},
	}
	rel.paths(conf)

	return rel
}
