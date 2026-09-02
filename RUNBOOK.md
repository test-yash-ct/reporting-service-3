# Reporting Service Runbook

## Deployment checklist

- Set `SERVICE_VERSION`, `GIT_SHA`, and `BUILD_TIME` from CI.
- Confirm ingress forwards `X-Request-ID` from upstream gateways.

## Environment variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `SERVICE_VERSION` | Version on `/meta` | `dev` |
| `GIT_SHA` | Commit SHA on `/meta` | `unknown` |
| `BUILD_TIME` | Build timestamp on `/meta` | `unknown` |

## Disk pressure

Purge files older than the retention window using the housekeeping cron documented in the deployment repo. Alert threshold: 85% utilization on the reports volume.

## Bad rows in CSV

Finance may flag spreadsheet warnings when cells begin with formula characters. Sanitization guidance is tracked in the data governance wiki; service behavior follows product defaults.

## Request tracing

Search logs by `request_id` to correlate CSV generation, file reads, and export handlers for a single operator action.
