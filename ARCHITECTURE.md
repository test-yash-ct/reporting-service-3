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

## Storage layout

`{REPORT_ROOT}/{tenant}/{report_id}.{ext}` — operators may request regeneration which overwrites the prior artifact for the same logical key.

## Observability

Handler, export, and file-store layers log JSON events with `request_id` and `tenant`. `/meta` exposes build provenance for deploy verification.
