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
	URL       string
	File      File
	Checksum  Checksum
	Available Available
}

// File ...
type File struct {
	Src        string
	Binary     string
	Mode       string
	BinaryPath string
}

// Checksum ...
type Checksum struct {
	URL    string
	File   string
	Format string
}

type Available struct {
	OS   OS
	Arch Arch
}

type OS struct {
	Linux   string
	Windows string
	Darwin  string
}

type Arch struct {
	i386  string
	amd64 string
	arm64 string
	arm   string
}
