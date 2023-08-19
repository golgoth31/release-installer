// Package output ...
package output

import (
	"github.com/pterm/pterm"
)

// StepTitle string.
func (o *Output) StepTitle(str string, level int) {
	pterm.DefaultSection.WithLevel(level).Println(str)
}

// SuccessTitle string.
func (o *Output) SuccessTitle(str string) {
	pterm.DefaultSection.WithStyle(&pterm.ThemeDefault.SuccessMessageStyle).Print(str)
}

// JumpLine string.
func (o *Output) JumpLine() {
	pterm.DefaultBasicText.Print("\n")
}

// NoFormat string.
func (o *Output) NoFormat(str string) {
	pterm.DefaultBasicText.Printf("%s\n", str)
}

// Info string.
func (o *Output) Info(str string) {
	pref := pterm.Prefix{
		Text:  "\u00BB",
		Style: &pterm.ThemeDefault.InfoPrefixStyle,
	}
	pterm.Info.WithPrefix(pref).Println(str)
}

// Success string.
func (o *Output) Success(str string) {
	pref := pterm.Prefix{
		Text:  "\u221A",
		Style: &pterm.ThemeDefault.SuccessPrefixStyle,
	}
	pterm.Success.WithPrefix(pref).Println(str)
}

// SuccessString string.
func (o *Output) SuccessString(str string) string {
	pref := pterm.Prefix{
		Text:  "\u221A",
		Style: &pterm.ThemeDefault.SuccessPrefixStyle,
	}

	return pterm.Success.WithPrefix(pref).Sprint(str)
}

// Warn string.
func (o *Output) Warn(str string) {
	pref := pterm.Prefix{
		Text:  "\u26A0",
		Style: &pterm.ThemeDefault.WarningPrefixStyle,
	}
	pterm.Warning.WithPrefix(pref).Println(str)
}
