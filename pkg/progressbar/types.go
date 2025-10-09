package progressbar

import (
	"io"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// ProgressBar implements the getter.ProgressTracker interface for file downloads.
type ProgressBar struct {
	mu sync.Mutex
}

// ProgressReader wraps an io.Reader to track progress.
type ProgressReader struct {
	io.Reader
	progress   progress.Model
	total      int64
	current    int64
	lastUpdate time.Time
	mu         sync.Mutex
}

// ProgressWriter wraps an io.Writer to track progress.
type ProgressWriter struct {
	io.Writer
	progress   progress.Model
	total      int64
	current    int64
	lastUpdate time.Time
	mu         sync.Mutex
}

// ProgressConfig holds configuration for progress tracking.
type ProgressConfig struct {
	Width     int
	ShowSpeed bool
	ShowETA   bool
	Style     lipgloss.Style
}

// DefaultProgressConfig returns default configuration for progress bars.
func DefaultProgressConfig() ProgressConfig {
	return ProgressConfig{
		Width:     80,
		ShowSpeed: true,
		ShowETA:   true,
		Style:     lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
	}
}
