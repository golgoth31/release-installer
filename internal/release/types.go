package release

// Release ...
type Release struct {
	APIVersion string `json:"apiVersion"`
	Kind       string
	Metadata   Metadata
	Spec       Spec
}

// Metadata ...
type Metadata struct {
	Name string
	Web  string
}

// Spec ...
type Spec struct {
	URL      string
	File     File
	Checksum Checksum
}

// File ...
type File struct {
	Archive    string
	Binary     string
	BinaryPath string
}

// Checksum ...
type Checksum struct {
	URL    string
	File   string
	Format string
}
