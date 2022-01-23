package config

type Config struct {
	RepoURL   string
	Release   *Release
	Reference *Reference
}

type Release struct {
	Path       string
	APIVersion string
	Kind       string
}

type Reference struct {
	Path string
}
