package config

// Config is the main internal config.
type Config struct {
	RepoURL   string
	Release   *Release
	Reference *Reference
}

// Release is the release config.
type Release struct {
	Path       string
	APIVersion string
	Kind       string
}

// Reference is the reference config.
type Reference struct {
	Path string
}
