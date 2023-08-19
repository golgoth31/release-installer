package installv1

// Install struct.
type Install struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

// Metadata struct.
type Metadata struct {
	Release string `json:"release"`
}

// Spec struct.
type Spec struct {
	Version string `json:"version"`
	Os      string `json:"os"`
	Arch    string `json:"arch"`
	Path    string `json:"path"`
	Default bool   `json:"default"`
}
