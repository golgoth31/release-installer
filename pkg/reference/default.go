package reference

var myself = &Reference{
	APIVersion: "release/v1",
	Kind:       "Release",
	Metadata: Metadata{
		Name: "release-installer",
		Web:  "https://github.com/golgoth31/release-installer",
	},
	Spec: Spec{
		Repo: Repo{
			Type:  "github",
			Name:  "release-installer",
			Owner: "golgoth31",
		},
		File: File{
			URL:        "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
			Src:        "ri_{{ .Version }}_{{ .Os }}_{{ .Arch }}",
			BinaryPath: ".",
			Binary:     "ri",
			Mode:       "file",
			Link:       "",
		},
		Checksum: Checksum{
			Check:  true,
			URL:    "https://github.com/golgoth31/release-installer/releases/download/{{ .Version }}",
			File:   "ri_{{ .Version }}_SHA256SUMS.txt",
			Format: "sha256",
		},
		Available: Available{
			Os: Os{
				Linux:   "linux",
				Darwin:  "darwin",
				Windows: "windows",
			},
			Arch: Arch{
				Amd64: "amd64",
				Arm64: "arm64",
				Arm:   "armv7",
			},
		},
	},
}
