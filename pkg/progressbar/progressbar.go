package progressbar

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
)

// New creates a new ProgressBar instance.
func New() *ProgressBar {
	return &ProgressBar{}
}

// TrackProgress creates a progress tracker for file downloads.
// This implements the getter.ProgressTracker interface.
func (pb *ProgressBar) TrackProgress(
	src string,
	currentSize,
	totalSize int64,
	stream io.ReadCloser,
) io.ReadCloser {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	// Create a progress model with proper settings
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)
	// Initialize the progress bar
	prog.SetPercent(0.0)

	// Create a progress reader
	pr := &ProgressReader{
		Reader:     stream,
		progress:   prog,
		total:      totalSize,
		current:    currentSize,
		lastUpdate: time.Now(),
	}

	return pr
}

// Read implements io.Reader interface and updates progress.
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)

	pr.mu.Lock()
	pr.current += int64(n)

	// Update progress every 100ms to avoid too frequent updates
	if time.Since(pr.lastUpdate) > 100*time.Millisecond {
		pr.updateProgress()
		pr.lastUpdate = time.Now()
	}
	pr.mu.Unlock()

	return n, err
}

// Close implements io.Closer interface.
func (pr *ProgressReader) Close() error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	// Final progress update
	pr.updateProgress()
	fmt.Println() // New line after progress bar

	if closer, ok := pr.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// updateProgress renders the current progress.
func (pr *ProgressReader) updateProgress() {
	if pr.total <= 0 {
		return
	}

	percent := float64(pr.current) / float64(pr.total)
	pr.progress.SetPercent(percent)

	// Create progress bar display - use simple bar for reliability
	bar := createSimpleBar(percent, 50)

	// Add file info
	info := fmt.Sprintf(" %s", formatBytes(pr.current))
	if pr.total > 0 {
		info += fmt.Sprintf("/%s", formatBytes(pr.total))
	}

	// Add percentage
	info += fmt.Sprintf(" (%.1f%%)", percent*100)

	// Combine and display - ensure we have enough space and clear the line
	display := fmt.Sprintf("\r\033[K%s%s", bar, info)
	fmt.Print(display)
	os.Stdout.Sync()
}

// NewWriter creates a progress writer for uploads or writes.
func NewWriter(w io.Writer, total int64) *ProgressWriter {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)
	// Initialize the progress bar
	prog.SetPercent(0.0)

	return &ProgressWriter{
		Writer:     w,
		progress:   prog,
		total:      total,
		current:    0,
		lastUpdate: time.Now(),
	}
}

// Write implements io.Writer interface and updates progress.
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)

	pw.mu.Lock()
	pw.current += int64(n)

	// Update progress every 100ms
	if time.Since(pw.lastUpdate) > 100*time.Millisecond {
		pw.updateProgress()
		pw.lastUpdate = time.Now()
	}
	pw.mu.Unlock()

	return n, err
}

// Close implements io.Closer interface.
func (pw *ProgressWriter) Close() error {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	// Final progress update
	pw.updateProgress()
	fmt.Println() // New line after progress bar

	if closer, ok := pw.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// updateProgress renders the current progress for writer.
func (pw *ProgressWriter) updateProgress() {
	if pw.total <= 0 {
		return
	}

	percent := float64(pw.current) / float64(pw.total)
	pw.progress.SetPercent(percent)

	// Create progress bar display - use simple bar for reliability
	bar := createSimpleBar(percent, 50)

	// Add file info
	info := fmt.Sprintf(" %s", formatBytes(pw.current))
	if pw.total > 0 {
		info += fmt.Sprintf("/%s", formatBytes(pw.total))
	}

	// Add percentage
	info += fmt.Sprintf(" (%.1f%%)", percent*100)

	// Combine and display - ensure we have enough space and clear the line
	display := fmt.Sprintf("\r\033[K%s%s", bar, info)
	fmt.Print(display)
	os.Stdout.Sync()
}

// createSimpleBar creates a simple text-based progress bar.
func createSimpleBar(percent float64, width int) string {
	// Ensure percent is between 0 and 1
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	filled := int(float64(width) * percent)
	// bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	bar := strings.Repeat("=", filled) + strings.Repeat("_", width-filled)
	return fmt.Sprintf("[%s]", bar)
}

// formatBytes formats bytes into human readable format.
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// SimpleProgress creates a simple progress display without bubbles.
func SimpleProgress(current, total int64, label string) {
	if total <= 0 {
		return
	}

	percent := float64(current) / float64(total)
	barWidth := 30
	filled := int(float64(barWidth) * percent)

	bar := strings.Repeat("=", filled) + strings.Repeat("_", barWidth-filled)

	display := fmt.Sprintf("\r%s [%s] %.1f%% %s/%s",
		label, bar, percent*100, formatBytes(current), formatBytes(total))

	fmt.Print(display)
	os.Stdout.Sync()
}
