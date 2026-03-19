# Telemetry Pipeline

Designs related to receivers, processors, exporters, and the OTel collector runtime.

## Designs

| Title | Document | Status | Description |
|---|---|---|---|
| OTel Collector Runtime | `cmd/amazon-cloudwatch-agent/`, `cmd/amazon-cloudwatch-agent/internal/` | 🚀 Implemented | Main agent binary built on `otelcol.Collector`; wires config provider, default components, and Telegraf plugin checks into the OTel service lifecycle. Predates design practice — no design doc available |
| Default Components Registration | `service/defaultcomponents/components.go` | 🚀 Implemented | Central `Factories()` function registering all OTel component factories (17 receivers, 25 processors, 7 exporters, 11 extensions). Predates design practice — no design doc available |
| Config Provider | `service/configprovider/` | 🚀 Implemented | `GetSettings()` builds `otelcol.ConfigProviderSettings` from URI list; `flags.go` defines `--config` flag for OTel config file paths. Predates design practice — no design doc available |
| Plugin Registry | `service/registry/` | 🚀 Implemented | Dynamic component registration via `Register()`/`Options()` with `WithReceiver`/`WithProcessor`/`WithExporter`/`WithExtension` option functions; supplements static `defaultcomponents` registration. Predates design practice — no design doc available |
