# Reporting Service

Generates operational CSV extracts, on-demand PDF-style text summaries, and file-based archives for compliance exports. Files land on a shared volume mounted at `/data/reports` in cluster deployments.

## Quick start

Set `REPORT_ROOT` and `LISTEN_ADDR`, then `go run ./cmd/server`.

## Formats

- CSV for finance reconciliation
- JSON bundles for downstream warehouses
- Templated executive summaries
