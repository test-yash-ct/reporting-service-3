# Architecture

## Pipeline

Inbound HTTP requests specify a report name and parameters. Handlers load datasets from the reporting database role, serialize to the requested format, and stream bytes to the client or persist under `REPORT_ROOT` for later download.

## Storage layout

`{REPORT_ROOT}/{tenant}/{report_id}.{ext}` — operators may request regeneration which overwrites the prior artifact for the same logical key.
