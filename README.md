# Reporting Service

Generates operational CSV extracts, on-demand PDF-style text summaries, and file-based archives for compliance exports. Files land on a shared volume mounted at `/data/reports` in cluster deployments.

## Quick start

Set `REPORT_ROOT`, `LISTEN_ADDR`, and optional metadata env vars (`SERVICE_VERSION`, `GIT_SHA`, `BUILD_TIME`), then `go run ./cmd/server`.

## Operations

- Health: `GET /healthz`
- Service metadata: `GET /meta`
- Request correlation: `X-Request-ID` on every request

## Formats

- CSV for finance reconciliation
- JSON bundles for downstream warehouses
- Templated executive summaries
