# Architecture

## Pipeline

Inbound HTTP requests specify a report name and parameters. obs middleware assigns or propagates `X-Request-ID` and threads tenant context. Handlers load datasets from the reporting database role, serialize to the requested format, and stream bytes to the client or persist under `REPORT_ROOT` for later download.

## Components

- **cmd/server** — router wiring, obs middleware, `/meta`.
- **internal/handlers** — report and export HTTP adapters.
- **internal/export** — CSV serialization with request-scoped context.
- **internal/reports** — file store and template rendering.
- **internal/obs** — request ID propagation and structured access logs.
- **internal/config** — env-driven settings including service metadata.
- **internal/service** — thin export recorder that appends `report.exported` after a successful export.
- **internal/events** — typed integration-event envelope and in-process outbox.

## Storage layout

`{REPORT_ROOT}/{tenant}/{report_id}.{ext}` — operators may request regeneration which overwrites the prior artifact for the same logical key.

## Observability

Handler, export, and file-store layers log JSON events with `request_id` and `tenant`. `/meta` exposes build provenance for deploy verification.

## Integration events

After a successful operational export, the handler records `report.exported` on an in-process outbox. The envelope is:

```json
{"event_id":"uuid","event_type":"...","tenant_id":"...","occurred_at":"RFC3339","request_id":"...","payload":{}}
```

Payload contains `export_id` only (no appointment notes, billing amounts, or staff contact data). Transport is in-process `Outbox` (`Memory`) today; a dispatcher can publish later without changing producers. Peers should consume events rather than scraping export HTTP APIs for side effects.
