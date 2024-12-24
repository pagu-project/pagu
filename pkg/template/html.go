package template

import (
	"bytes"
	"html/template"
)

func ExecuteHTML(tmpl string, keyValue any) (string, error) {
	buf := bytes.Buffer{}
	tp, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	if err := tp.Execute(&buf, keyValue); err != nil {
		return "", err
	}

	return buf.String(), nil
}
