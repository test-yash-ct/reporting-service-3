package reports

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"

	"github.com/healthops/reporting-service/internal/obs"
)

type RenderRequest struct {
	Template string          `json:"template"`
	Payload  json.RawMessage `json:"payload"`
}

func RenderSummary(ctx context.Context, req RenderRequest) (string, error) {
	_ = obs.RequestIDFromContext(ctx)
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
