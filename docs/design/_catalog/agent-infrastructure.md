# Agent Infrastructure

Designs related to health reporting, entity store, packaging, and CLI tools.

## Designs

| Title | Document | Status | Description |
|---|---|---|---|
| Architecture Documentation | `docs/design/260318-docs-onboarding/plan.md` | ✅ Done | Comprehensive architecture docs, recipes, and pattern inventory |
| Agent Health Extension | — | 🚀 Implemented | Agent self-health reporting via `extension/agenthealth/` (Oct 2023). Predates design practice — no design doc available |
| Entity Store Extension | — | 🚀 Implemented | EC2/EKS entity metadata store via `extension/entitystore/` (Jul 2024). Predates design practice — no design doc available |
| Server Extension | — | 🚀 Implemented | Local HTTP server for status/config via `extension/server/` (Sep 2024). Predates design practice — no design doc available |
| Config Wizard | — | 🚀 Implemented | Interactive config wizard via `cmd/amazon-cloudwatch-agent-config-wizard/` (Jun 2020). Predates design practice — no design doc available |
| Config Downloader | — | 🚀 Implemented | SSM/S3 config fetcher via `cmd/config-downloader/` (Jun 2020). Predates design practice — no design doc available |
| Workload Discovery | — | 🚀 Implemented | Workload detection binary via `cmd/workload-discovery/` (Sep 2025). Predates design practice — no design doc available |
