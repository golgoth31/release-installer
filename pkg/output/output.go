package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StepTitle prints a section title with the specified level.
func (o *Output) StepTitle(str string, level int) {
	indent := strings.Repeat("  ", level-1)
	styled := o.styles.Section.Render(fmt.Sprintf("%s%s", indent, str))
	fmt.Fprintln(os.Stdout, styled)
}

// SuccessTitle prints a success title.
func (o *Output) SuccessTitle(str string) {
	styled := o.styles.Success.Render(str)
	fmt.Fprint(os.Stdout, styled)
}

// JumpLine prints a newline.
func (o *Output) JumpLine() {
	fmt.Fprint(os.Stdout, "\n")
}

// NoFormat prints text without any formatting.
func (o *Output) NoFormat(str string) {
	fmt.Fprintf(os.Stdout, "%s\n", str)
}

// Info prints an info message with a prefix.
func (o *Output) Info(str string) {
	prefix := o.styles.InfoPrefix.Render("»")
	styled := o.styles.Info.Render(str)
	fmt.Fprintf(os.Stdout, "%s %s\n", prefix, styled)
}

// Success prints a success message with a prefix.
func (o *Output) Success(str string) {
	prefix := o.styles.SuccessPrefix.Render("✓")
	styled := o.styles.Success.Render(str)
	fmt.Fprintf(os.Stdout, "%s %s\n", prefix, styled)
}

// SuccessString returns a formatted success string without printing.
func (o *Output) SuccessString(str string) string {
	prefix := o.styles.SuccessPrefix.Render("✓")
	styled := o.styles.Success.Render(str)
	return fmt.Sprintf("%s %s", prefix, styled)
}

// Warn prints a warning message with a prefix.
func (o *Output) Warn(str string) {
	prefix := o.styles.WarningPrefix.Render("⚠")
	styled := o.styles.Warning.Render(str)
	fmt.Fprintf(os.Stdout, "%s %s\n", prefix, styled)
}

// Error prints an error message with a prefix.
func (o *Output) Error(str string) {
	prefix := o.styles.Error.Render("✗")
	styled := o.styles.Error.Render(str)
	fmt.Fprintf(os.Stdout, "%s %s\n", prefix, styled)
}

// Fatal prints a fatal error message and exits.
func (o *Output) Fatal(str string) {
	o.Error(str)
	os.Exit(1)
}

// WithStyle applies a custom style to the output.
func (o *Output) WithStyle(style lipgloss.Style) *Output {
	return &Output{
		styles: &Styles{
			Section:       style,
			Success:       o.styles.Success,
			Info:          o.styles.Info,
			Warning:       o.styles.Warning,
			Error:         o.styles.Error,
			Prefix:        o.styles.Prefix,
			SuccessPrefix: o.styles.SuccessPrefix,
			InfoPrefix:    o.styles.InfoPrefix,
			WarningPrefix: o.styles.WarningPrefix,
		},
	}
}
