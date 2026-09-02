# Reporting Service API

## Request correlation

Clients may send `X-Request-ID`; the service echoes it on responses and includes `request_id` in structured JSON access logs.

## `GET /meta`

Returns `service`, `version`, `build_time`, and `git_sha` from environment variables.

## `GET /v1/reports/:name.csv`

Query params: `tenant`, `from`, `to`. Streams a CSV extract.

## `GET /v1/reports/file`

Query param `path` — relative path under the tenant directory for operator retrieval.

## `POST /v1/reports/render`

Body: JSON document describing rows and a `template` string for the summary layout.

## `GET /v1/exports/operational`

Returns a JSON snapshot used by the NOC dashboard. Optional `scope` query is accepted for future filtering.
