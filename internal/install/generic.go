package install

import (
	"github.com/golgoth31/release-installer/internal/release"
)

// NewInstall ...
func NewInstall(rel string) *Install {
	releaseData = release.New(rel)

	return &Install{
		APIVersion: "install/v1",
		Kind:       "Install",
	}
}
