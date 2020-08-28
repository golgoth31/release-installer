package install

import (
	"io"
	"sync"
)

type Install struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Release string `json:"release"`
}

type Spec struct {
	Version string `json:"version"`
	Os      string `json:"os"`
	Arch    string `json:"arch"`
	Default bool   `json:"default"`
	Path    string `json:"path"`
}

type progressBar struct {
	// lock everything below
	lock sync.Mutex
}

type readCloser struct {
	io.Reader
	close func() error
}
