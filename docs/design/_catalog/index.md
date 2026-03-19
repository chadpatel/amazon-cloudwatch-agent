# Design Catalog

This catalog indexes design documents for the Amazon CloudWatch Agent. Each entry links to a design doc that describes the motivation, approach, and implementation details for a feature or subsystem.

## How to Use

1. Check the category pages below to find existing designs related to your work area.
2. Read relevant design docs before starting implementation.
3. When creating a new feature or making a significant change, add a design doc and register it in the appropriate category.

Source of truth hierarchy: **Code > docs/arch/ > READMEs > docs/design/_catalog/**

Always verify catalog entries against actual code. When code and docs disagree, code is authoritative.

## Categories

| Category | Description | Architecture Docs |
|---|---|---|
| [Metrics & Collection](metrics-collection.md) | System metrics, custom metrics (StatsD/collectd), Prometheus, GPU, NVMe, and Container Insights collection | [Domain Model](../../arch/domain-model.md) |
| [Logs & Traces](logs-traces.md) | Log file tailing, CloudWatch Logs shipping, log filtering, X-Ray and OTel trace collection | — |
| [Config & Translation](config-translation.md) | JSON config schema, config translation engine, config download, config wizard | [Configuration](../../arch/configuration.md), [Config Translation](../../arch/config-translation.md) |
| [Plugins & Pipeline](plugins-pipeline.md) | OTel pipeline components, Telegraf adapter, processors, exporters, extensions | [Plugins & Pipeline](../../arch/plugins-pipeline.md) |
| [Platform & Packaging](platform-packaging.md) | Multi-OS support, RPM/DEB/ZIP/tar packaging, Docker builds, agent lifecycle | [Deployment](../../arch/deployment.md) |
| [Agent Infrastructure](agent-infrastructure.md) | Health reporting, entity store, build system, testing, CLI tools | [Scope](../../arch/SCOPE.md), [Build & Tooling](../../arch/build-and-tooling.md), [Project Structure](../../arch/project-structure.md), [Testing](../../arch/testing.md), [Error Handling](../../arch/error-handling.md), [Telemetry](../../arch/telemetry.md), [Concurrency](../../arch/concurrency.md), [Recipes](../../arch/recipes.md) |
| [Container & Kubernetes](container-kubernetes.md) | Container Insights, ECS/EKS decorators, Kubernetes metadata | — |
| [Telemetry Pipeline](telemetry-pipeline.md) | Receivers, processors, exporters, OTel collector runtime | [Plugins & Pipeline](../../arch/plugins-pipeline.md) |

## Adding New Designs

Add entries to the appropriate category file using this format:

```markdown
| Title | `path/to/design.md` | Status | One-line summary |
```

### Status Values

| Status | Meaning |
|---|---|
| `Draft` | Design is in progress, not yet reviewed |
| `Proposed` | Design is ready for review |
| `Accepted` | Design is approved for implementation |
| `Implemented` | Design has been implemented in code |
| `Deprecated` | Design is no longer relevant |
