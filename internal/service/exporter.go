package service

import (
	"context"
	"log"

	"github.com/healthops/reporting-service/internal/events"
	"github.com/healthops/reporting-service/internal/obs"
)

// Exporter records a successful operational export as an integration event.
type Exporter struct {
	Outbox events.Outbox
}

func NewExporter(outbox events.Outbox) *Exporter {
	if outbox == nil {
		outbox = events.Nop{}
	}
	return &Exporter{Outbox: outbox}
}

func (e *Exporter) RecordOperationalExport(ctx context.Context, tenantID, exportID string) {
	requestID := obs.RequestIDFromContext(ctx)
	ev := events.New(events.TypeReportExported, tenantID, requestID, map[string]any{
		"export_id": exportID,
	})
	if err := e.Outbox.Append(ctx, ev); err != nil {
		log.Printf("outbox append failed request_id=%s event_type=%s", requestID, ev.EventType)
	}
}
