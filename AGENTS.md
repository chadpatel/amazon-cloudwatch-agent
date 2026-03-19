# Amazon CloudWatch Agent - AI Agent Guide

For AI Coding Assistants: This file provides essential context for understanding and working with this codebase. For detailed docs, see docs/arch/.

## Key Information

| Field | Value |
|---|---|
| Package | `github.com/aws/amazon-cloudwatch-agent` |
| Version | See `CWAGENT_VERSION` file |
| Language | Go 1.25.7 |
| Build System | Make |
| Test Framework | `go test` |
| Lint | golangci-lint v2 (`.golangci.yml`) |

## Architecture

```
┌─────────────────────────────────────────────────────┐
│              CloudWatch Agent (OTel Collector)       │
│                                                     │
│  ┌──────────┐  ┌────────────┐  ┌─────────────────┐ │
│  │ Receivers │→ │ Processors │→ │    Exporters    │ │
│  │ (inputs)  │  │            │  │ (outputs)       │ │
│  └──────────┘  └────────────┘  └─────────────────┘ │
│       ↑                                             │
│  ┌──────────┐  ┌────────────┐  ┌─────────────────┐ │
│  │ Telegraf  │  │ Extensions │  │   Translator    │ │
│  │ Adapter   │  │ (health,   │  │ (JSON→OTel     │ │
│  │ Receiver  │  │  entity,   │  │  config)        │ │
│  │           │  │  k8s, srv) │  │                 │ │
│  └──────────┘  └────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────┘
```

The agent is built on the OpenTelemetry Collector framework and uses Telegraf plugins via an adapter receiver. User-facing JSON configuration is translated into OTel collector YAML by the translator package.

## Core Features

1. **System Metrics Collection** — Collects in-guest metrics from EC2, on-premises, and hybrid servers (Linux, Windows, macOS)
2. **Log Collection** — Tails log files and ships to CloudWatch Logs with filtering support
3. **Custom Metrics** — Ingests StatsD and collectd protocol metrics
4. **Trace Collection** — Collects OpenTelemetry and X-Ray traces
5. **Container Insights** — Kubernetes and ECS container monitoring with metadata decoration
6. **Config Translation** — Converts user JSON config to OTel collector pipeline configuration
7. **Agent Health** — Self-reports health and performance data to CloudWatch
8. **Workload Discovery** — Detects running workloads (JVM, Tomcat, Kafka, NVIDIA GPU)

## Quick Start

```bash
# Build all platforms
make build

# Build + package (RPM, DEB, ZIP, tar.gz)
make release

# Run tests
make test

# Run tests with race detector
make test-data-race

# Lint (installs golangci-lint if needed)
make lint

# Format code (goimports + go fmt)
make fmt

# Docker build (linux amd64+arm64)
make docker-build

# Clean build artifacts
make clean
```

## Project Structure

```
cmd/                          # Binary entry points
├── amazon-cloudwatch-agent/  #   Main agent binary
├── start-amazon-cloudwatch-agent/ # Agent launcher
├── config-translator/        #   JSON→OTel config translator CLI
├── config-downloader/        #   SSM/S3 config fetcher
├── amazon-cloudwatch-agent-config-wizard/ # Interactive config wizard
├── workload-discovery/       #   Workload detection binary
└── xray-migration/           #   X-Ray daemon migration tool
plugins/                      # Telegraf-style plugins
├── inputs/                   #   Data sources (logfile, statsd, prometheus, nvidia_smi, win_perf_counters)
├── outputs/                  #   Destinations (cloudwatch, cloudwatchlogs)
└── processors/               #   Data transforms (ec2tagger, k8sdecorator, ecsdecorator, awsentity, awsapplicationsignals, gpuattributes)
receiver/                     # OTel receivers
├── adapter/                  #   Wraps Telegraf plugins as OTel receivers
├── awsnvmereceiver/          #   NVMe disk metrics
└── systemmetricsreceiver/    #   CPU, memory, disk, ethtool, JVM metrics
extension/                    # OTel extensions
├── agenthealth/              #   Agent health reporting
├── entitystore/              #   EC2/EKS entity metadata
├── k8smetadata/              #   Kubernetes pod/node metadata
└── server/                   #   Local HTTP server for status/config
translator/                   # Config translation engine
├── translate/                #   JSON config → OTel pipeline rules
├── tocwconfig/               #   Final OTel config assembly
├── jsonconfig/               #   JSON config merge logic
├── config/                   #   Schema, defaults, OS config
└── cmdutil/                  #   CLI translation utilities
internal/                     # Shared internal packages
├── detector/                 #   Workload detection (JVM, Tomcat, Kafka, NVIDIA)
├── k8sCommon/                #   Kubernetes client utilities
├── containerinsightscommon/  #   Container Insights constants
├── state/                    #   State management
├── retryer/                  #   IMDS retry logic
└── ...                       #   cloudwatch, tls, publisher, util, etc.
cfg/                          # Configuration
├── aws/                      #   AWS credentials and SDK config
├── commonconfig/             #   Common agent config (TOML)
└── envconfig/                #   Environment variable config
service/                      # OTel service wiring
├── defaultcomponents/        #   Component factory registration
├── configprovider/           #   Config provider flags
└── registry/                 #   Plugin registry
processor/                    # OTel processors
└── rollupprocessor/          #   Metric rollup
metric/                       # Metric types and distribution
tool/                         # CLI tools and utilities
packaging/                    # RPM, DEB, Windows, macOS packaging
```

## Development Workflow

1. Read `AGENTS.md`, `README.md`, and `docs/arch/overview.md` before starting work
2. Read design docs in `docs/design/` for feature context
3. Read architecture docs in `docs/arch/` before working in any area
4. Respect `.gitignore` — never commit build artifacts or `coverage.txt`

### Build Commands

```bash
make build          # Full multi-platform build
make test           # Unit tests with coverage
make lint           # golangci-lint v2 + license check + import order
make fmt            # goimports + go fmt
make clean          # Remove build/ and release/
```

### Commit Format

Use conventional commits with a `---` separator and `Prompt:` prefix, wrapped at 72 columns:

```
feat(translator): add support for new metric type

Implement translation rules for the new histogram metric type
in the OTel pipeline configuration.

---
Prompt: Add histogram metric translation support
```

Never use `--no-verify` on commits.

## Steering Rules

Engineering best practices and Go language guardrails are in `.kiro/steering/rules/`. Agents should read and follow them before making changes.

## Codebase Patterns

| Pattern | Documentation |
|---|---|
| Architecture overview | `docs/arch/overview.md` |
| Scope and boundaries | `docs/arch/SCOPE.md` |
| Build system and tooling | `docs/arch/build-and-tooling.md` |
| Project structure | `docs/arch/project-structure.md` |
| Configuration management | `docs/arch/configuration.md` |
| Config translation | `docs/arch/config-translation.md` |
| Plugins and OTel pipeline | `docs/arch/plugins-pipeline.md` |
| Domain model and types | `docs/arch/domain-model.md` |
| Testing infrastructure | `docs/arch/testing.md` |
| Concurrency patterns | `docs/arch/concurrency.md` |
| Error handling | `docs/arch/error-handling.md` |
| Telemetry and observability | `docs/arch/telemetry.md` |
| Deployment and packaging | `docs/arch/deployment.md` |
| Recipes and boilerplate | `docs/arch/recipes.md` |

## Refactoring Guidelines

- Update all consumers when changing interfaces or function signatures
- Do not leave deprecated aliases — migrate all call sites in the same change
- Run `make test` and `make lint` after refactoring

## Design Catalog

Source of truth hierarchy: **Code > docs/arch/ > READMEs > docs/design/_catalog/**

- The catalog in `docs/design/_catalog/` indexes design documents
- Always verify catalog entries against actual code
- When code and docs disagree, code is authoritative
