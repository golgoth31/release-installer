package config

import (
	"fmt"

	common_proto "github.com/golgoth31/release-installer/pkg/proto/common"
	reference_proto "github.com/golgoth31/release-installer/pkg/proto/reference"
)

func (c *Config) SetMyself() reference_proto.Reference {
	myself := reference_proto.Reference{
		ApiVersion: "release/v1",
		Kind:       "Release",
		Metadata: &common_proto.Metadata{
			Name: "release-installer",
			Web:  c.RepoURL,
		},
		Spec: &reference_proto.Spec{
			Repo: &reference_proto.Repo{
				Type:  "github",
				Name:  "release-installer",
				Owner: "golgoth31",
			},
			File: &reference_proto.File{
				Url:        fmt.Sprintf("%s/releases/download/{{ .Version }}", c.RepoURL),
				Src:        "ri_{{ .Version }}_{{ .Os }}_{{ .Arch }}",
				BinaryPath: ".",
				Binary:     "ri",
				Mode:       "file",
				Link:       "",
			},
			Checksum: &reference_proto.Checksum{
				Check:  true,
				Url:    fmt.Sprintf("%s/releases/download/{{ .Version }}", c.RepoURL),
				File:   "ri_{{ .Version }}_SHA256SUMS.txt",
				Format: "sha256",
			},
			Available: &reference_proto.Available{
				Os: &reference_proto.Os{
					Linux:   "linux",
					Darwin:  "darwin",
					Windows: "windows",
				},
				Arch: &reference_proto.Arch{
					Amd64: "amd64",
					Arm64: "arm64",
					Arm:   "armv7",
					I386:  "i386",
				},
			},
		},
	}

	return myself
}
