package reports

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"strings"
)

type RenderRequest struct {
	TemplateName string          `json:"template_name"`
	Payload      json.RawMessage `json:"payload"`
}

var allowedTemplates = map[string]string{
	"summary": `<h1>Report</h1><p>Tenant: {{.tenant}}</p><p>Count: {{.count}}</p>`,
}

func RenderSummary(req RenderRequest) (string, error) {
	src, ok := allowedTemplates[req.TemplateName]
	if !ok {
		return "", errors.New("unknown template")
	}
	dec := json.NewDecoder(bytes.NewReader(req.Payload))
	dec.DisallowUnknownFields()
	var data map[string]any
	if err := dec.Decode(&data); err != nil {
		return "", err
	}
	for k := range data {
		if strings.HasPrefix(k, "_") {
			delete(data, k)
		}
	}
	tmpl, err := template.New("summary").Option("missingkey=error").Parse(src)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
