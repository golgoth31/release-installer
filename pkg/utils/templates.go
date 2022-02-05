package utils

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	release_proto "github.com/golgoth31/release-installer/pkg/proto/release"
)

func TemplateStringRelease(input string, spec *release_proto.Release) (string, error) {
	templateBuffer := bytes.Buffer{}
	templateString := template.Must(
		template.
			New("templateString").
			Funcs(sprig.FuncMap()).
			Parse(input),
	)

	if err := templateString.Execute(&templateBuffer, spec.Spec); err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}
