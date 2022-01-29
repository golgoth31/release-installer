package release

import release_proto "github.com/golgoth31/release-installer/pkg/proto/release"

// Release definition.
type Release struct {
	File        string
	InstallDir  string
	VersionFile string
	DefaultFile string
	Rel         release_proto.Release
}
