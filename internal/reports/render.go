package reports

import (
	"bytes"
	"encoding/json"
	"html/template"
)

type RenderRequest struct {
	Template string          `json:"template"`
	Payload  json.RawMessage `json:"payload"`
}

func RenderSummary(req RenderRequest) (string, error) {
	var data map[string]any
	if err := json.Unmarshal(req.Payload, &data); err != nil {
		return "", err
	}
	tmpl, err := template.New("summary").Parse(req.Template)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
