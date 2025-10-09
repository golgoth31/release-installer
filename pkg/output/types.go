package output

import (
	"github.com/charmbracelet/lipgloss"
)

// Output provides styled console output functionality.
type Output struct {
	styles *Styles
}

// Styles contains all the styling definitions for output.
type Styles struct {
	Section       lipgloss.Style
	Success       lipgloss.Style
	Info          lipgloss.Style
	Warning       lipgloss.Style
	Error         lipgloss.Style
	Prefix        lipgloss.Style
	SuccessPrefix lipgloss.Style
	InfoPrefix    lipgloss.Style
	WarningPrefix lipgloss.Style
}

// New creates a new Output instance with default styles.
func New() *Output {
	return &Output{
		styles: newStyles(),
	}
}

// newStyles creates default styles for output.
func newStyles() *Styles {
	return &Styles{
		Section: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")), // Blue
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")), // Green
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")), // Yellow
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")), // Bright Yellow
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")), // Red
		Prefix: lipgloss.NewStyle().
			Bold(true),
		SuccessPrefix: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true),
		InfoPrefix: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true),
		WarningPrefix: lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true),
	}
}
