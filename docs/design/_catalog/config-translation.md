# Config & Translation

Design documents for the JSON config schema, config translation engine (JSON → OTel YAML/TOML), config download from SSM/S3, and the interactive config wizard.

## Designs

| Title | Path | Status | Summary |
|---|---|---|---|
| JSON Config Schema | `translator/config/schema.json` | 🚀 Implemented | JSON Schema draft-04 defining agent/metrics/logs/traces sections with validation constraints. Predates design practice — no design doc available. (Jun 2020) |
| Config Translation Engine | `translator/translate/translate.go`, `translator/cmdutil/translatorutil.go` | 🚀 Implemented | Legacy Telegraf rule-based translation from JSON config to TOML output (CWAgent.conf). Uses `Rule` interface with `ApplyRule()` and `init()` self-registration. Predates design practice — no design doc available. (Jun 2020) |
| Config Merge Logic | `translator/jsonconfig/mergeJsonConfig.go` | 🚀 Implemented | Multi-file JSON config merge with default/append/remove modes and conflict resolution rules. Predates design practice — no design doc available. (Jun 2020) |
| OTel Translation Rules | `translator/translate/otel/translate_otel.go` | 🚀 Implemented | Generic `Translator[C, ID]` interface for translating JSON config to OTel collector YAML. Pipeline-based registration with 10 pipeline types. Predates design practice — no design doc available. (Oct 2022) |
| Config Downloader (SSM/S3) | `cmd/config-downloader/` | 🚀 Implemented | CLI tool to fetch agent config from SSM Parameter Store or S3, with multi-config support and default/append/remove modes. Predates design practice — no design doc available. (Jun 2020) |
| Config Wizard | `cmd/amazon-cloudwatch-agent-config-wizard/` | 🚀 Implemented | Interactive CLI wizard for generating agent JSON config via guided prompts. Predates design practice — no design doc available. (Jun 2020) |
