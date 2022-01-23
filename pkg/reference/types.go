package reference

<<<<<<< HEAD
// Reference ...
type Reference struct {
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
	Link       string
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

// Os ...
type Os struct {
	Linux   string `json:"linux,omitempty"`
	Windows string `json:"windows,omitempty"`
	Darwin  string `json:"darwin,omitempty"`
}

// Arch ...
type Arch struct {
	Amd64 string `json:"amd64,omitempty"`
	Arm64 string `json:"arm64,omitempty"`
	Arm   string `json:"arm,omitempty"`
=======
import reference_proto "github.com/golgoth31/release-installer/pkg/proto/reference"

// Reference structure.
type Reference struct {
	File string
	Ref  reference_proto.Reference
>>>>>>> ce9013b (Feat: V2)
}
