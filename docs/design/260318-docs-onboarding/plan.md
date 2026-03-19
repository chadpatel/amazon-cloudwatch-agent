# Architecture Documentation Onboarding

Create comprehensive architecture documentation for the Amazon CloudWatch Agent repository.

## Documentation Standard

- **Do document**: Purpose, setup, key concepts, public API, non-obvious patterns, code examples from the actual codebase
- **Don't document**: Implementation details readable from code, every function/type, internal mechanics
- Code should be self-documenting; docs explain the "why" and "how to use", not the "what"
- Every claim must be verified against actual source code
- Include mermaid diagrams showing relationships between components
- Code examples must be real code from this repo, not hypothetical
- **Decision Trees**: When a subsystem has multiple approaches to the same goal (identified in Pattern Detection), the arch doc must include a decision tree or comparison table explaining when to use each approach — not just describe them in isolation
- **Testing in every recipe**: Each atomic recipe must end with a "#### Testing" subsection showing the corresponding test pattern — how to test the thing the recipe creates. Use the test boilerplate from the Phase 2 Test Recipe Candidates if available, otherwise derive from existing test files for similar patterns.

## Reference Files

- `AGENTS.md` — Project conventions and existing patterns table
- `docs/arch/overview.md` — Existing architecture overview
- `docs/design/_catalog/index.md` — Catalog structure and categories
- `README.md` — Project overview
- `Makefile` — Build system targets
- `.golangci.yml` — Linter configuration
- `service/defaultcomponents/components.go` — OTel component registration
- `plugins/plugins.go` — Telegraf plugin registration
- `translator/cmdutil/translatorutil.go` — Config translation orchestration
- `translator/config/schema.json` — JSON config schema
- `receiver/adapter/receiver.go` — Telegraf-to-OTel adapter
- `cfg/aws/credentials.go` — AWS credential chain
- `cfg/commonconfig/common-config.toml` — Common agent config
- `cfg/envconfig/` — Environment variable config

## File Changes

| File | Change | Task |
|------|--------|------|
| `docs/arch/SCOPE.md` | New | 1.1 |
| `docs/arch/build-and-tooling.md` | New | 1.2 |
| `docs/arch/project-structure.md` | New | 1.3 |
| `docs/arch/configuration.md` | New | 1.4 |
| `docs/arch/config-translation.md` | New | 1.5 |
| `docs/arch/plugins-pipeline.md` | New | 1.6 |
| `docs/arch/domain-model.md` | New | 1.7 |
| `docs/arch/testing.md` | New | 1.8 |
| `docs/arch/concurrency.md` | New | 1.9 |
| `docs/arch/error-handling.md` | New | 1.10 |
| `docs/arch/telemetry.md` | New | 1.11 |
| `docs/arch/deployment.md` | New | 1.12 |
| `docs/design/260318-docs-onboarding/worklog.md` | New (Pattern Inventory) | 2.1 |
| `docs/arch/recipes.md` | New | 3.1 |
| `AGENTS.md` | Update patterns table | 4.1 |
| `docs/arch/overview.md` | Update with new subsystem links | 4.2 |
| `docs/design/_catalog/index.md` | Update categories with arch doc links | 4.3 |
| `docs/design/_catalog/agent-infrastructure.md` | Add catalog entry for docs-onboarding + backfill | 5.1 |
| `docs/design/_catalog/metrics-collection.md` | Backfill existing feature entries | 5.2 |
| `docs/design/_catalog/logs-traces.md` | Backfill existing feature entries | 5.3 |
| `docs/design/_catalog/config-translation.md` | Backfill existing feature entries | 5.4 |
| `docs/design/_catalog/plugins-pipeline.md` | Backfill existing feature entries | 5.5 |
| `docs/design/_catalog/platform-packaging.md` | Backfill existing feature entries | 5.6 |
| `docs/design/_catalog/container-kubernetes.md` | Backfill existing feature entries | 5.7 |
| `docs/design/_catalog/telemetry-pipeline.md` | Backfill existing feature entries | 5.8 |
| `docs/design/260318-docs-onboarding/worklog.md` | Append validation results | 6.1 |

## Tasks

### Phase 1: Create Architecture Docs

- [x] **1.1** Create `docs/arch/SCOPE.md` — Define what this repo owns and does not own. Read `AGENTS.md`, `README.md`, `go.mod`, and `docs/arch/overview.md`. In-scope: CW agent binary, config translation, Telegraf plugin adapter, OTel pipeline components, system/custom/container metrics collection, log tailing, trace collection, config wizard, config downloader, workload discovery, packaging (RPM/DEB/ZIP/tar/Docker). Out-of-scope: CloudWatch service backend, OTel Collector core framework (upstream), Telegraf core (upstream), AWS SDK, CloudWatch Logs service, X-Ray service backend, SSM Parameter Store, Container Insights backend processing, EKS/ECS control plane. Include boundary notes about the agent's relationship to upstream OTel Collector and Telegraf.

- [x] **1.2** Create `docs/arch/build-and-tooling.md` — Document the build system and developer tooling. Read `Makefile`, `.golangci.yml`, `.github/workflows/PR-build.yml`, `.github/workflows/PR-test.yml`, `.github/workflows/test-build.yml`. Document: Make targets (build, release, test, lint, fmt, clean, docker-build), cross-compilation (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64), CGO settings per platform, LDFLAGS for version injection, golangci-lint v2 config (gosec, misspell, nonamedreturns, revive), goimports with local prefix, addlicense for MIT headers, impi for import ordering, secret scanning in check_secrets target, CI workflow structure. Include mermaid diagram of build pipeline.

- [x] **1.3** Create `docs/arch/project-structure.md` — Document the Go module layout and package organization. Read top-level directory listing, `go.mod`, `cmd/` binaries, and key package entry points. Document: the 7 binary entry points in `cmd/`, the plugin system (`plugins/inputs/`, `plugins/outputs/`, `plugins/processors/`), the OTel component layout (`receiver/`, `processor/`, `extension/`), the translator package hierarchy, `internal/` shared packages, `cfg/` configuration packages, `service/` wiring, `tool/` CLI utilities, `packaging/` platform-specific packaging. Include mermaid diagram showing package dependency flow. Explain the dual nature: Telegraf-style plugins registered in `plugins/plugins.go` vs OTel components registered in `service/defaultcomponents/components.go`.

- [x] **1.4** Create `docs/arch/configuration.md` — Document the configuration management system. Read `translator/config/schema.json`, `translator/config/schema.go`, `translator/config/defaultConfig.go`, `translator/config/os.go`, `translator/config/mode.go`, `translator/config/envconst.go`, `cfg/commonconfig/`, `cfg/envconfig/`, `cfg/aws/credentials.go`. Document: JSON config schema structure (agent, metrics, logs, traces sections), schema validation, common config TOML format, environment variable overrides, AWS credential chain setup, OS-specific config handling, config modes (ec2, onPrem, auto). Include decision tree for when to use JSON config vs common config vs env vars.

- [x] **1.5** Create `docs/arch/config-translation.md` — Document the config translation engine. Read `translator/cmdutil/translatorutil.go`, `translator/translate/translate.go`, `translator/rule.go`, `translator/processRuleToApply.go`, `translator/processNoRuleToApply.go`, `translator/processDefaultConfig.go`, `translator/processRuleToApplyAndMerge.go`, `translator/translate/otel/translate_otel.go`, `translator/tocwconfig/tocwconfig_test.go` (for sample configs), `translator/jsonconfig/mergeJsonConfig.go`, `translator/context/context.go`, `translator/registerrules/register_rules.go`. Document: the two-phase translation (legacy Telegraf rules + OTel pipeline rules), the Rule interface pattern, how JSON config sections map to OTel pipeline components, the merge logic for multi-file configs, the translation context, sample config test pattern (JSON→YAML/TOML/conf golden files). Include mermaid diagram of translation flow from JSON input to OTel YAML output.

- [x] **1.6** Create `docs/arch/plugins-pipeline.md` — Document the OTel pipeline and plugin system. Read `service/defaultcomponents/components.go`, `plugins/plugins.go`, `receiver/adapter/receiver.go`, `plugins/outputs/cloudwatch/cloudwatch.go`, `plugins/outputs/cloudwatch/factory.go`, `plugins/outputs/cloudwatchlogs/cloudwatchlogs.go`, `plugins/processors/ec2tagger/ec2tagger.go`, `plugins/processors/ec2tagger/factory.go`, `plugins/processors/awsentity/processor.go`, `plugins/processors/awsentity/factory.go`, `plugins/processors/awsapplicationsignals/processor.go`, `plugins/processors/awsapplicationsignals/factory.go`, `plugins/processors/gpuattributes/processor.go`, `plugins/processors/kueueattributes/processor.go`, `plugins/processors/k8sdecorator/k8sdecorator.go`, `plugins/processors/ecsdecorator/ecsdecorator.go`, `plugins/processors/prometheusadapter/processor.go`, `receiver/awsnvmereceiver/`, `receiver/systemmetricsreceiver/`, `extension/agenthealth/extension.go`, `extension/entitystore/extension.go`, `extension/k8smetadata/extension.go`, `extension/server/extension.go`, `processor/rollupprocessor/`. Document: the OTel component factory pattern (Config, Factory, NewFactory), the Telegraf adapter bridge, how Telegraf input plugins become OTel receivers, the processor chain, exporter implementations, extension lifecycle. Include mermaid diagram of a complete pipeline from receiver to exporter. Include decision tree: when to create a Telegraf-style plugin vs a native OTel component.

- [x] **1.7** Create `docs/arch/domain-model.md` — Document shared types, metric types, and constants. Read `metric/distribution/`, `metric/`, `internal/metric/metric.go`, `internal/entity/entity.go`, `internal/containerinsightscommon/const.go`, `internal/containerinsightscommon/k8sconst.go`, `internal/containerinsightscommon/ecsconst.go`, `internal/logscommon/const.go`, `internal/constants/constants.go`, `internal/structuredlogscommon/common.go`, `internal/cloudwatch/unit.go`, `cfg/agentinfo/`. Document: metric distribution types, CloudWatch unit mapping, entity model, Container Insights metric constants, structured logs common types. Include key type definitions with code examples.

- [x] **1.8** Create `docs/arch/testing.md` — Document the testing infrastructure and patterns. Read `Makefile` (test targets), `translator/tocwconfig/tocwconfig_test.go` (golden file pattern), `plugins/outputs/cloudwatch/cloudwatch_test.go`, `plugins/processors/ec2tagger/ec2tagger_test.go`, `plugins/processors/awsentity/processor_test.go`, `extension/entitystore/extension_test.go`, `extension/server/extension_test.go`, `internal/state/manager_test.go`, `tool/testutil/`, `internal/detector/detectortest/`, `internal/state/statetest/`. Document: unit test patterns, golden file testing (JSON→YAML/TOML/conf sample configs), test data directories, test helper packages, race detector testing (`make test-data-race` and the PKG_WITH_DATA_RACE exclusion list), integration test workflows in `.github/workflows/`, factory test pattern for OTel components. Include decision tree: when to use golden file tests vs unit tests vs integration tests.

- [x] **1.9** Create `docs/arch/concurrency.md` — Document concurrency patterns. Read `internal/publisher/publisher.go`, `internal/publisher/nonblockingfifoqueue.go`, `internal/mapWithExpiry/mapWIthExpiry.go`, `internal/tls/certWatcher.go`, `plugins/inputs/logfile/logfile.go` (tailer goroutines), `plugins/inputs/logfile/tailersrc.go`, `plugins/outputs/cloudwatch/cloudwatch.go` (aggregator), `Makefile` (PKG_WITH_DATA_RACE list). Document: the non-blocking publisher pattern, map-with-expiry for TTL caches, TLS certificate file watching, log file tailing goroutines, CloudWatch metric aggregation, known data race packages and mitigation strategy. Include mermaid diagram of publisher pub/sub flow.

- [x] **1.10** Create `docs/arch/error-handling.md` — Document error handling patterns. Read `internal/retryer/imdsretryer.go`, `internal/retryer/logthrottle.go`, `extension/entitystore/retryer.go`, `internal/httpclient/httpclient.go`, `internal/clientutil/client.go`, `plugins/outputs/cloudwatch/cloudwatch.go` (error handling in Write), `plugins/outputs/cloudwatchlogs/cloudwatchlogs.go`. Document: IMDS retry with backoff, log-throttled retries, entity store retry pattern, HTTP client error handling, CloudWatch API error handling (throttling, invalid parameter), error propagation patterns in the OTel pipeline. Include decision tree: when to retry vs fail fast vs log-and-continue.

- [x] **1.11** Create `docs/arch/telemetry.md` — Document agent self-observability. Read `extension/agenthealth/extension.go`, `extension/agenthealth/config.go`, `extension/agenthealth/handler/stats/`, `extension/agenthealth/handler/useragent/`, `extension/agenthealth/metadata/metadata.go`, `profiler/profiler.go`, `logger/logger.go`, `logs/logs.go`, `cfg/agentinfo/`. Document: agent health extension (self-reporting to CloudWatch), stats collection handlers, user-agent string construction, profiler for performance data, logger configuration (lumberjack rotation), usage_data config option. Include mermaid diagram of health data flow.

- [x] **1.12** Create `docs/arch/deployment.md` — Document packaging and deployment. Read `packaging/linux/amazon-cloudwatch-agent.spec`, `packaging/linux/amazon-cloudwatch-agent.conf`, `packaging/dependencies/amazon-cloudwatch-agent-ctl`, `packaging/dependencies/amazon-cloudwatch-agent.service`, `packaging/windows/`, `packaging/darwin/`, `packaging/README.md`, `amazon-cloudwatch-container-insights/cloudwatch-agent-dockerfile/`, `Makefile` (package-* targets), `Tools/src/`. Document: RPM/DEB packaging flow, Windows ZIP packaging, macOS tar packaging, Docker multi-arch builds, systemd service file, agent control script (amazon-cloudwatch-agent-ctl), installation paths per OS, JMX metrics JAR bundling. Include mermaid diagram of the release pipeline.

### Phase 2: Pattern Detection

- [x] **2.1** Create `docs/design/260318-docs-onboarding/worklog.md` with a Pattern Inventory. Perform structural analysis of the codebase:

  **Feature/module directories:** Enumerate all top-level plugin directories under `plugins/inputs/`, `plugins/outputs/`, `plugins/processors/`, `receiver/`, `extension/`, `translator/translate/otel/receiver/`, `translator/translate/otel/processor/`, `translator/translate/otel/exporter/`, `translator/translate/otel/extension/`, `translator/translate/otel/pipeline/`, `internal/detector/`, `cmd/`.

  **Repeated directory structures:** Compare internal structures of plugin directories. Look for the OTel component pattern: `factory.go`, `config.go`, `processor.go`/`extension.go`/`receiver.go`, `*_test.go`, `internal/`, `testdata/`, `README.md`. Compare Telegraf plugin structures. Compare translator rule structures.

  **Central registration files:** Identify `service/defaultcomponents/components.go`, `plugins/plugins.go`, `translator/registerrules/register_rules.go`, and any other files imported/modified by 3+ features.

  **Mutually exclusive approaches:** Detect: (a) Telegraf-style plugins vs native OTel components, (b) legacy translation rules (`translator/translate/metrics/`, `translator/translate/logs/`) vs OTel translation rules (`translator/translate/otel/`), (c) OS-specific file suffixes (`_unix.go`, `_windows.go`, `_darwin.go`, `_notwindows.go`) vs build tags.

  **Test infrastructure:** Scan for test patterns: golden file tests in `translator/tocwconfig/sampleConfig/`, factory test pattern (`factory_test.go`), testdata directories, test helper packages (`tool/testutil/`, `internal/detector/detectortest/`, `internal/state/statetest/`), mock patterns.

  **Failure modes:** Search for error logging, retry patterns, sentinel returns, custom error types across subsystems.

  **Extension point ranking:** Count implementations of key interfaces: OTel component factories (receiver.Factory, processor.Factory, exporter.Factory, extension.Factory), Telegraf input/output/processor plugins, translator Rule interface, detector types.

  **File naming conventions:** Glob for `factory.go`, `config.go`, `*_test.go`, `*_unix.go`, `*_windows.go`, `*_darwin.go`, `*_notwindows.go`, `rule*.go`, `README.md` patterns across directories.

  **Feature complexity tiers:** Classify plugins/components by complexity: (a) simple single-file plugins (e.g., kueueattributes), (b) standard plugins with config+factory+processor+tests (e.g., ec2tagger, gpuattributes), (c) complex multi-package plugins (e.g., awsapplicationsignals, awsentity, cloudwatch output, logfile input), (d) full subsystems (e.g., translator, extension/entitystore).

  Output all findings in the worklog with one subsection per detection category. No empty sections.

### Phase 3: Create Recipes

- [x] **3.1** Create `docs/arch/recipes.md` — Using the Pattern Inventory from Phase 2, create recipes containing:

  **(1) Atomic Recipes:**
  - Add a new Telegraf input plugin
  - Add a new OTel processor (with factory pattern)
  - Add a new OTel extension
  - Add a new OTel receiver
  - Add a new OTel exporter
  - Add a translation rule for a new JSON config field
  - Add an OTel translation rule (translator/translate/otel/)
  - Add a new workload detector
  - Add a new CLI binary to cmd/
  - Add OS-specific code with build constraints
  - Add a golden file translation test

  Each atomic recipe must end with a `#### Testing` subsection.

  **(2) Composite Recipes** (end-to-end workflows):
  - Add a new metric type end-to-end (input plugin → translation rule → OTel pipeline → CloudWatch exporter)
  - Add a new log source end-to-end
  - Add Container Insights support for a new resource type

  **(3) Registration Points Checklist:**
  - `service/defaultcomponents/components.go` — Register OTel component factory
  - `plugins/plugins.go` — Register Telegraf plugin
  - `translator/registerrules/register_rules.go` — Register translation rules
  - `translator/translate/otel/translate_otel.go` — Register OTel translation
  - `Makefile` — Add new binary targets if adding a cmd/

  **(4) Common Failure Modes:**
  - IMDS metadata fetch failures (retryer)
  - CloudWatch API throttling
  - Config translation validation errors
  - Plugin initialization failures
  - Data race conditions (known packages)

  **(5) Decision Trees:**
  - Telegraf plugin vs native OTel component
  - Legacy translation rule vs OTel translation rule
  - OS-specific files vs build tags
  - Golden file test vs unit test vs integration test

  **(6) Recipe Coverage Matrix:**
  Table mapping each arch doc to its corresponding recipe(s).

  **(7) Test-Writing Recipes:**
  - How to write a factory test for an OTel component
  - How to write a golden file translation test
  - How to write a processor unit test
  - How to add a package to the data race exclusion list

  **(8) Boilerplate Templates:**
  Top 5 extension points with fill-in-the-blanks templates:
  - OTel Processor (factory.go + config.go + processor.go)
  - OTel Extension (factory.go + config.go + extension.go)
  - OTel Receiver (factory.go + config.go + receiver.go)
  - Telegraf Input Plugin
  - Translation Rule (translator/translate/otel/)

### Phase 4: Update References

- [x] **4.1** Update `AGENTS.md` — Add entries to the Codebase Patterns table for each new architecture doc created in Phase 1. Add entries like: `| Config translation | docs/arch/config-translation.md |`, etc.

- [x] **4.2** Update `docs/arch/overview.md` — Add a "Related Architecture Docs" section at the bottom linking to all subsystem docs created in Phase 1. If any new subsystems were discovered during documentation that aren't in the overview, add them to the appropriate section.

- [x] **4.3** Update `docs/design/_catalog/index.md` — Add a "Current Docs" column or section to the categories table linking to the architecture docs. Add a new category row if needed for cross-cutting concerns (build, testing, deployment).

### Phase 5: Catalog Registration

- [x] **5.1** Update `docs/design/_catalog/agent-infrastructure.md` — Add catalog entry for this docs-onboarding design: `| Architecture Documentation | docs/design/260318-docs-onboarding/plan.md | 🟡 In Progress | Comprehensive architecture docs, recipes, and pattern inventory |`. Also backfill entries for existing features: agenthealth extension, entitystore extension, server extension, config wizard, config downloader, workload discovery. For each, scan `git log --oneline --diff-filter=A -- <path>` to find approximate addition date. Use format from catalog. Set status to `🚀 Implemented`. Add note "Predates design practice — no design doc available" where no design doc exists.

- [x] **5.2** Update `docs/design/_catalog/metrics-collection.md` — Backfill entries for: system metrics (cpu, mem, disk, net via Telegraf), StatsD input, collectd input, Prometheus input, nvidia_smi input, win_perf_counters input, NVMe receiver, system metrics receiver, JMX metrics. Scan git log for dates. Set status to `🚀 Implemented`.

- [x] **5.3** Update `docs/design/_catalog/logs-traces.md` — Backfill entries for: logfile input plugin, CloudWatch Logs output, windows_event_log input, log filtering, X-Ray trace collection, OTel trace pipeline. Scan git log for dates.

- [x] **5.4** Update `docs/design/_catalog/config-translation.md` — Backfill entries for: JSON config schema, config translation engine, config merge logic, OTel translation rules, config downloader (SSM/S3), config wizard. Scan git log for dates.

- [x] **5.5** Update `docs/design/_catalog/plugins-pipeline.md` — Backfill entries for: Telegraf adapter receiver, ec2tagger processor, k8sdecorator processor, ecsdecorator processor, awsentity processor, awsapplicationsignals processor, gpuattributes processor, kueueattributes processor, prometheusadapter processor, rollup processor, CloudWatch exporter, CloudWatch Logs exporter. Scan git log for dates.

- [x] **5.6** Update `docs/design/_catalog/platform-packaging.md` — Backfill entries for: RPM packaging, DEB packaging, Windows ZIP packaging, macOS tar packaging, Docker multi-arch builds, systemd service, agent-ctl script. Scan git log for dates.

- [x] **5.7** Update `docs/design/_catalog/container-kubernetes.md` — Backfill entries for: Container Insights base, k8sdecorator, ecsdecorator, k8smetadata extension, Kueue support, GPU Container Insights. Scan git log for dates.

- [x] **5.8** Update `docs/design/_catalog/telemetry-pipeline.md` — Backfill entries for: OTel collector runtime, default components registration, config provider, plugin registry. Scan git log for dates.

### Phase 6: Validation

- [x] **6.1** Update `docs/design/260318-docs-onboarding/worklog.md` — Append validation results:
  (a) For each architecture doc, verify `recipes.md` has at least one recipe covering its primary use case. Output the Recipe Coverage Matrix.
  (b) Verify every Decision Points table entry has a corresponding decision tree in `recipes.md`.
  (c) Verify top 5 Boilerplate Candidates all have templates.
  (d) Verify `SCOPE.md` has at least 3 in-scope and 1 out-of-scope entries.
  (e) Run `git log --diff-filter=A --name-only` to identify most frequently created file types. Verify top 5 new-file patterns each have a recipe. Flag any without.
  (f) For every Pattern Inventory entry above threshold, verify `recipes.md` has a corresponding recipe. Create any missing recipes by appending to `recipes.md` before completing validation.
  Update the catalog entry status from `🟡 In Progress` to `✅ Done` in `docs/design/_catalog/agent-infrastructure.md`.
