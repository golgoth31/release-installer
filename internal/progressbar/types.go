package progressbar

import (
	"io"
	"sync"
)

type ProgressBar struct {
	// lock everything below
	lock sync.Mutex
}

type readCloser struct {
	io.Reader
	close func() error
}
