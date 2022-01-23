package install

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig/v3"
)

func (i *Install) templates() (
	releaseURL bytes.Buffer,
	releaseFileName bytes.Buffer,
	checksumURL bytes.Buffer,
	checksumFileName bytes.Buffer,
	binaryPath bytes.Buffer,
	binaryFile bytes.Buffer,
	revertError error) {
	revertError = nil
	// template strings
	treleaseURL := template.Must(
		template.New("releaseURL").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.URL),
	)
	if err := treleaseURL.Execute(&releaseURL, i.Spec); err != nil {
		revertError = err
	}

	treleaseFileName := template.Must(
		template.New("releaseFileName").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.Src),
	)
	if err := treleaseFileName.Execute(&releaseFileName, i.Spec); err != nil {
		revertError = err
	}

	tchecksumURL := template.Must(
		template.New("checksumURL").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.Checksum.URL),
	)
	if err := tchecksumURL.Execute(&checksumURL, i.Spec); err != nil {
		revertError = err
	}

	tchecksumFileName := template.Must(
		template.New("checksumFileName").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.Checksum.File),
	)
	if err := tchecksumFileName.Execute(&checksumFileName, i.Spec); err != nil {
		revertError = err
	}

	tbinaryPath := template.Must(
		template.New("binaryPath").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.BinaryPath),
	)
	if err := tbinaryPath.Execute(&binaryPath, i.Spec); err != nil {
		revertError = err
	}

	tbinaryFile := template.Must(
		template.New("binaryFile").Funcs(sprig.FuncMap()).Parse(releaseData.Spec.File.Binary),
	)
	if err := tbinaryFile.Execute(&binaryFile, i.Spec); err != nil {
		revertError = err
	}

	return releaseURL, releaseFileName, checksumURL, checksumFileName, binaryPath, binaryFile, revertError
}
