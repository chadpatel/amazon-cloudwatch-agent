# Logs & Traces

Design documents for log file tailing, CloudWatch Logs shipping, log filtering, X-Ray trace collection, and OpenTelemetry trace pipelines.

## Designs

| Title | Path | Status | Summary |
|---|---|---|---|
| Log File Input Plugin | `plugins/inputs/logfile/` | 🚀 Implemented (Jun 2020) | Tails log files with glob patterns, multiline support, offset tracking, and backpressure handling. Predates design practice — no design doc available |
| CloudWatch Logs Output | `plugins/outputs/cloudwatchlogs/` | 🚀 Implemented (Jun 2020) | Ships logs to CloudWatch Logs via PutLogEvents with retry, batching, and auto log group/stream creation. Predates design practice — no design doc available |
| Windows Event Log Input | `plugins/inputs/windows_event_log/` | 🚀 Implemented (Jun 2020) | Collects Windows event logs via wineventlog API. Predates design practice — no design doc available |
| Log Filtering | `plugins/inputs/logfile/logfilter*.go` | 🚀 Implemented (Jan 2022) | Filters log lines by include/exclude patterns before shipping. Predates design practice — no design doc available |
| EMF Logs Pipeline | `translator/translate/otel/pipeline/emf_logs/` | 🚀 Implemented (Feb 2023) | OTel pipeline for Embedded Metric Format (EMF) structured logs. Predates design practice — no design doc available |
| X-Ray Trace Pipeline | `translator/translate/otel/pipeline/xray/`, `translator/translate/otel/exporter/awsxray/` | 🚀 Implemented (Mar 2023) | OTel pipeline for X-Ray trace collection and export. Predates design practice — no design doc available |
| OTel Trace Receiver (OTLP) | `translator/translate/otel/receiver/otlp/` | 🚀 Implemented (Apr 2023) | OTLP receiver translation for traces and metrics ingestion. Predates design practice — no design doc available |
| Application Signals Traces | `translator/translate/otel/pipeline/applicationsignals/` | 🚀 Implemented (May 2024) | OTel pipeline for Application Signals trace processing with rules, resolver, and cardinality control. Predates design practice — no design doc available |
| CloudWatch Logs Exporter (OTel) | `translator/translate/otel/exporter/awscloudwatchlogs/` | 🚀 Implemented (Nov 2024) | OTel exporter translation for CloudWatch Logs via tcplog/udplog receivers. Predates design practice — no design doc available |
| TCP/UDP Log Receivers | `translator/translate/otel/receiver/tcplog/`, `translator/translate/otel/receiver/udplog/` | 🚀 Implemented (Nov 2024) | OTel receiver translations for TCP and UDP log ingestion. Predates design practice — no design doc available |
