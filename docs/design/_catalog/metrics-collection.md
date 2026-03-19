# Metrics & Collection

Design documents for system metrics, custom metrics (StatsD/collectd), Prometheus scraping, GPU metrics, NVMe disk metrics, and Container Insights collection.

## Designs

| Title | Document | Status | Description |
|---|---|---|---|
| System Metrics (Telegraf) | — | 🚀 Implemented | CPU, memory, disk, diskio, net, swap, processes, procstat, ethtool via upstream Telegraf input plugins in `plugins/plugins.go` (Jun 2020). Predates design practice — no design doc available |
| Windows Performance Counters | — | 🚀 Implemented | Windows perf counter collection via `plugins/inputs/win_perf_counters/` (Jun 2020). Predates design practice — no design doc available |
| StatsD Input | — | 🚀 Implemented | StatsD protocol metric ingestion via `plugins/inputs/statsd/` (Jun 2020). Predates design practice — no design doc available |
| collectd Input | — | 🚀 Implemented | collectd protocol metric ingestion via upstream Telegraf `socket_listener` plugin + `tool/processors/collectd/` (Jun 2020). Predates design practice — no design doc available |
| Prometheus Input | — | 🚀 Implemented | Prometheus scraping via `plugins/inputs/prometheus/` with calculator, filter, handler, and target allocator (Apr 2023). Predates design practice — no design doc available |
| NVIDIA GPU Metrics | — | 🚀 Implemented | NVIDIA SMI GPU metrics via `plugins/inputs/nvidia_smi/` with schema v11/v12 support (Feb 2024). Predates design practice — no design doc available |
| JMX Metrics | — | 🚀 Implemented | JVM metrics via JMX with bundled JAR, OTel JMX receiver translation in `translator/translate/otel/receiver/jmx/` (May 2024). Predates design practice — no design doc available |
| NVMe Receiver | — | 🚀 Implemented | NVMe disk metrics via native OTel receiver `receiver/awsnvmereceiver/` (Aug 2025). Predates design practice — no design doc available |
| System Metrics Receiver | — | 🚀 Implemented | Native OTel receiver for CPU, memory, disk, ethtool, JVM metrics via `receiver/systemmetricsreceiver/` with per-subsystem Linux scrapers (Mar 2026). Predates design practice — no design doc available |
| Prometheus Histogram Support | `docs/design/260319-prometheus-histograms/plan.md` | ✅ Done | Classic histogram conversion for CloudWatch PMD path using density-aware bucket mapping |
