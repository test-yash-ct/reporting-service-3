# Reporting Service Runbook

## Disk pressure

Purge files older than the retention window using the housekeeping cron documented in the deployment repo. Alert threshold: 85% utilization on the reports volume.

## Bad rows in CSV

Finance may flag spreadsheet warnings when cells begin with formula characters. Sanitization guidance is tracked in the data governance wiki; service behavior follows product defaults.
