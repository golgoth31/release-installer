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
	Repo      Repo
	File      File
	Checksum  Checksum
	Available Available
}

// File ...
type File struct {
	URL        string
	Src        string
	Binary     string
	Mode       string
	BinaryPath string
}

// Checksum ...
type Checksum struct {
	Check  bool
	URL    string
	File   string
	Format string
}

// Repo ...
type Repo struct {
	Type  string
	Name  string
	Owner string
}

// Available ...
type Available struct {
	Os   Os
	Arch Arch
}

// OS ...
type Os struct {
	Linux   string `json:"linux,omitempty"`
	Windows string `json:"windows,omitempty"`
	Darwin  string `json:"darwin,omitempty"`
}

// Arch ...
type Arch struct {
	I386  string `json:"i386,omitempty"`
	Amd64 string `json:"amd64,omitempty"`
	Arm64 string `json:"arm64,omitempty"`
	Arm   string `json:"arm,omitempty"`
}
