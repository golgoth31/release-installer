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
	Linux   string `json:"linux,omitempty"`
	Windows string `json:"windows,omitempty"`
	Darwin  string `json:"darwin,omitempty"`
}

type Arch struct {
	I386  string `json:"i386,omitempty"`
	Amd64 string `json:"amd64,omitempty"`
	Arm64 string `json:"arm64,omitempty"`
	Arm   string `json:"arm,omitempty"`
}
