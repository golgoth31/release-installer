package release

type Release struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   Metadata
	Spec       Spec
}

type Metadata struct {
	Name string
	Web  string
}

type Spec struct {
	Url      string
	File     File
	Checksum Checksum
}

type File struct {
	Name       string
	BinaryPath string
}

type Checksum struct {
	Url    string
	File   string
	Format string
}