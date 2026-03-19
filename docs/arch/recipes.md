# Recipes

Step-by-step guides for common development tasks. Each recipe references real code patterns from this codebase.

Related docs: [plugins-pipeline.md](plugins-pipeline.md), [config-translation.md](config-translation.md), [testing.md](testing.md), [project-structure.md](project-structure.md)

---

## Table of Contents

- [Atomic Recipes](#atomic-recipes)
  - [Add a New Telegraf Input Plugin](#add-a-new-telegraf-input-plugin)
  - [Add a New OTel Processor](#add-a-new-otel-processor)
  - [Add a New OTel Extension](#add-a-new-otel-extension)
  - [Add a New OTel Receiver](#add-a-new-otel-receiver)
  - [Add a New OTel Exporter](#add-a-new-otel-exporter)
  - [Add a Translation Rule for a New JSON Config Field](#add-a-translation-rule-for-a-new-json-config-field)
  - [Add an OTel Translation Rule](#add-an-otel-translation-rule)
  - [Add a New Workload Detector](#add-a-new-workload-detector)
  - [Add a New CLI Binary](#add-a-new-cli-binary)
  - [Add OS-Specific Code with Build Constraints](#add-os-specific-code-with-build-constraints)
  - [Add a Golden File Translation Test](#add-a-golden-file-translation-test)
- [Composite Recipes](#composite-recipes)
  - [Add a New Metric Type End-to-End](#add-a-new-metric-type-end-to-end)
  - [Add a New Log Source End-to-End](#add-a-new-log-source-end-to-end)
  - [Add Container Insights Support for a New Resource Type](#add-container-insights-support-for-a-new-resource-type)
- [Registration Points Checklist](#registration-points-checklist)
- [Common Failure Modes](#common-failure-modes)
- [Decision Trees](#decision-trees)
- [Recipe Coverage Matrix](#recipe-coverage-matrix)
- [Test-Writing Recipes](#test-writing-recipes)
- [Boilerplate Templates](#boilerplate-templates)

---

## Atomic Recipes

### Add a New Telegraf Input Plugin

> **When to use**: Only for wrapping an existing upstream Telegraf plugin. For new data sources, prefer a [native OTel receiver](#add-a-new-otel-receiver). See [Decision Trees](#telegraf-plugin-vs-native-otel-component).

1. Create `plugins/inputs/<name>/<name>.go` implementing `telegraf.Input`:

```go
package myinput

import (
    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
)

type MyInput struct {
    ServiceAddress string
}

func (m *MyInput) SampleConfig() string { return "" }
func (m *MyInput) Gather(acc telegraf.Accumulator) error {
    // Collect metrics and call acc.AddFields()/acc.AddGauge()
    return nil
}

func init() {
    inputs.Add("myinput", func() telegraf.Input {
        return &MyInput{ServiceAddress: ":8080"}
    })
}
```

Reference: [`plugins/inputs/statsd/statsd.go`](../../plugins/inputs/statsd/statsd.go) lines 615–627.

2. Register via blank import in `plugins/plugins.go`:

```go
_ "github.com/aws/amazon-cloudwatch-agent/plugins/inputs/myinput"
```

3. The adapter receiver (`receiver/adapter/`) automatically wraps Telegraf inputs as OTel receivers with type prefix `telegraf_`.

#### Testing

Write a unit test for `Gather()`:

```go
func TestGather(t *testing.T) {
    input := &MyInput{ServiceAddress: ":0"}
    acc := &testutil.Accumulator{}
    require.NoError(t, input.Gather(acc))
    assert.True(t, acc.HasMeasurement("my_metric"))
}
```

---

### Add a New OTel Processor

> **When to use**: Preferred for all new data processing components. See [`plugins/processors/kueueattributes/`](../../plugins/processors/kueueattributes/) as the minimal reference.

1. Create three files in `plugins/processors/<name>/`:

**config.go** — Configuration struct:

```go
package myprocessor

import "go.opentelemetry.io/collector/component"

type Config struct {
    // Add config fields with mapstructure tags
}

var _ component.Config = (*Config)(nil)

func (cfg *Config) Validate() error { return nil }
```

**factory.go** — Factory registration:

```go
package myprocessor

import (
    "context"
    "fmt"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/processor"
    "go.opentelemetry.io/collector/processor/processorhelper"
)

var TypeStr, _ = component.NewType("myprocessor")

func NewFactory() processor.Factory {
    return processor.NewFactory(
        TypeStr,
        func() component.Config { return &Config{} },
        processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha))
}

func createMetricsProcessor(
    ctx context.Context,
    set processor.Settings,
    cfg component.Config,
    next consumer.Metrics,
) (processor.Metrics, error) {
    c, ok := cfg.(*Config)
    if !ok {
        return nil, fmt.Errorf("configuration parsing error")
    }
    p := &myProcessor{config: c, logger: set.Logger}
    return processorhelper.NewMetrics(ctx, set, cfg, next, p.processMetrics,
        processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}))
}
```

**processor.go** — Processing logic:

```go
package myprocessor

import (
    "context"
    "go.opentelemetry.io/collector/pdata/pmetric"
    "go.uber.org/zap"
)

type myProcessor struct {
    config *Config
    logger *zap.Logger
}

func (p *myProcessor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
    // Transform metrics here
    return md, nil
}
```

2. Register in `service/defaultcomponents/components.go`:

```go
import "github.com/aws/amazon-cloudwatch-agent/plugins/processors/myprocessor"
// In Factories(), add to the processors list:
myprocessor.NewFactory(),
```

#### Testing

Create `factory_test.go` following the standard pattern:

```go
func TestCreateDefaultConfig(t *testing.T) {
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    assert.NotNil(t, cfg)
    assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateProcessor(t *testing.T) {
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    set := processortest.NewNopSettings(component.MustNewType("myprocessor"))

    mp, err := factory.CreateMetrics(context.Background(), set, cfg, consumertest.NewNop())
    assert.NoError(t, err)
    assert.NotNil(t, mp)

    tp, err := factory.CreateTraces(context.Background(), set, cfg, consumertest.NewNop())
    assert.Equal(t, err, pipeline.ErrSignalNotSupported)
    assert.Nil(t, tp)
}
```

Reference: [`plugins/processors/kueueattributes/factory_test.go`](../../plugins/processors/kueueattributes/factory_test.go).

---

### Add a New OTel Extension

> **When to use**: For cross-cutting services (health, metadata stores, HTTP servers) that don't process telemetry data directly. See [`extension/server/`](../../extension/server/) as the minimal reference.

1. Create three files in `extension/<name>/`:

**config.go**:

```go
package myextension

import "go.opentelemetry.io/collector/component"

type Config struct {
    ListenAddress string `mapstructure:"listen_addr"`
}

var _ component.Config = (*Config)(nil)
```

**factory.go**:

```go
package myextension

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/extension"
)

var TypeStr, _ = component.NewType("myextension")

func NewFactory() extension.Factory {
    return extension.NewFactory(TypeStr,
        func() component.Config { return &Config{} },
        createExtension,
        component.StabilityLevelAlpha)
}

func createExtension(_ context.Context, set extension.Settings, cfg component.Config) (extension.Extension, error) {
    return &myExtension{config: cfg.(*Config), logger: set.Logger}, nil
}
```

**extension.go**:

```go
package myextension

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/extension"
    "go.uber.org/zap"
)

type myExtension struct {
    config *Config
    logger *zap.Logger
}

var _ extension.Extension = (*myExtension)(nil)

func (e *myExtension) Start(_ context.Context, _ component.Host) error { return nil }
func (e *myExtension) Shutdown(_ context.Context) error { return nil }
```

2. Register in `service/defaultcomponents/components.go` under extensions.

#### Testing

```go
func TestFactory(t *testing.T) {
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    assert.NoError(t, componenttest.CheckConfigStruct(cfg))

    ext, err := factory.Create(context.Background(), extensiontest.NewNopSettings(TypeStr), cfg)
    assert.NoError(t, err)
    assert.NotNil(t, ext)
}
```

Reference: [`extension/agenthealth/factory_test.go`](../../extension/agenthealth/factory_test.go).

---

### Add a New OTel Receiver

> **When to use**: For new data collection sources. Use `scraperhelper` for polling-based metrics. See [`receiver/awsnvmereceiver/`](../../receiver/awsnvmereceiver/) as reference.

1. Create files in `receiver/<name>/`:

**config.go**:

```go
package myreceiver

import (
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/scraper/scraperhelper"
)

type Config struct {
    scraperhelper.ControllerConfig `mapstructure:",squash"`
}

var _ component.Config = (*Config)(nil)
```

**factory.go**:

```go
package myreceiver

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver"
    otelscraper "go.opentelemetry.io/collector/scraper"
    "go.opentelemetry.io/collector/scraper/scraperhelper"
)

var TypeStr, _ = component.NewType("myreceiver")

func NewFactory() receiver.Factory {
    return receiver.NewFactory(TypeStr,
        func() component.Config {
            return &Config{ControllerConfig: scraperhelper.NewDefaultControllerConfig()}
        },
        receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelAlpha))
}

func createMetricsReceiver(_ context.Context, set receiver.Settings, cfg component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
    c := cfg.(*Config)
    s := newScraper(c, set.Logger)
    scraper, err := otelscraper.NewMetrics(s.scrape, otelscraper.WithStart(s.start))
    if err != nil {
        return nil, err
    }
    return scraperhelper.NewMetricsController(&c.ControllerConfig, set, consumer,
        scraperhelper.AddScraper(TypeStr, scraper))
}
```

**scraper.go**:

```go
package myreceiver

import (
    "context"
    "go.opentelemetry.io/collector/pdata/pmetric"
    "go.uber.org/zap"
)

type scraper struct {
    config *Config
    logger *zap.Logger
}

func newScraper(cfg *Config, logger *zap.Logger) *scraper {
    return &scraper{config: cfg, logger: logger}
}

func (s *scraper) start(_ context.Context, _ component.Host) error { return nil }

func (s *scraper) scrape(_ context.Context) (pmetric.Metrics, error) {
    md := pmetric.NewMetrics()
    // Populate metrics
    return md, nil
}
```

2. Register in `service/defaultcomponents/components.go` under receivers.

#### Testing

```go
func TestFactory(t *testing.T) {
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    assert.NoError(t, componenttest.CheckConfigStruct(cfg))

    r, err := factory.CreateMetrics(context.Background(), receivertest.NewNopSettings(TypeStr), cfg, consumertest.NewNop())
    assert.NoError(t, err)
    assert.NotNil(t, r)
}
```

Reference: [`receiver/awsnvmereceiver/factory_test.go`](../../receiver/awsnvmereceiver/factory_test.go).

---

### Add a New OTel Exporter

> **When to use**: For new telemetry destinations. See [`plugins/outputs/cloudwatch/factory.go`](../../plugins/outputs/cloudwatch/factory.go) as reference.

1. Create files in a new exporter package:

**factory.go**:

```go
package myexporter

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/exporter"
    "go.opentelemetry.io/collector/exporter/exporterhelper"
)

var TypeStr, _ = component.NewType("myexporter")

func NewFactory() exporter.Factory {
    return exporter.NewFactory(TypeStr,
        func() component.Config { return &Config{} },
        exporter.WithMetrics(createMetricsExporter, component.StabilityLevelAlpha))
}

func createMetricsExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Metrics, error) {
    c := cfg.(*Config)
    e := &myExporter{config: c, logger: set.Logger}
    return exporterhelper.NewMetrics(ctx, set, cfg, e.pushMetrics)
}
```

2. Register in `service/defaultcomponents/components.go` under exporters.

#### Testing

```go
func TestFactory(t *testing.T) {
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    assert.NoError(t, componenttest.CheckConfigStruct(cfg))

    e, err := factory.CreateMetrics(context.Background(), exportertest.NewNopSettings(TypeStr), cfg)
    assert.NoError(t, err)
    assert.NotNil(t, e)
}
```

---

### Add a Translation Rule for a New JSON Config Field

> **When to use**: Only for legacy Telegraf TOML output. For OTel YAML output, use [Add an OTel Translation Rule](#add-an-otel-translation-rule). See [Decision Trees](#legacy-translation-rule-vs-otel-translation-rule).

1. Create a rule file in the appropriate translator directory (e.g., `translator/translate/metrics/metrics_collect/<plugin>/`):

```go
package myplugin

import "github.com/aws/amazon-cloudwatch-agent/translator"

type MyField struct{}

func (m *MyField) ApplyRule(input interface{}) (string, interface{}) {
    // Extract field from JSON config map
    im := input.(map[string]interface{})
    result := map[string]interface{}{}
    if val, ok := im["my_field"]; ok {
        result["my_field"] = val
    }
    return "my_field", result
}

func init() {
    m := new(MyField)
    translator.RegisterRule("my_field", m)
}
```

2. Register via blank import in `translator/registerrules/register_rules.go` (or the appropriate OS-specific file):

```go
_ "github.com/aws/amazon-cloudwatch-agent/translator/translate/metrics/metrics_collect/myplugin"
```

#### Testing

Add a golden file test — see [Add a Golden File Translation Test](#add-a-golden-file-translation-test).

---

### Add an OTel Translation Rule

> **When to use**: Preferred for all new config translation. Maps JSON config fields to OTel component configuration. See [`translator/translate/otel/processor/kueue/translator.go`](../../translator/translate/otel/processor/kueue/translator.go) as the minimal reference.

1. Create `translator/translate/otel/<type>/<name>/translator.go`:

```go
package mycomponent

import (
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/confmap"
    "go.opentelemetry.io/collector/processor"

    "github.com/aws/amazon-cloudwatch-agent/plugins/processors/myprocessor"
    "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

type translator struct {
    name    string
    factory processor.Factory
}

var _ common.ComponentTranslator = (*translator)(nil)

func NewTranslator() common.ComponentTranslator {
    return NewTranslatorWithName("")
}

func NewTranslatorWithName(name string) common.ComponentTranslator {
    return &translator{name, myprocessor.NewFactory()}
}

func (t *translator) ID() component.ID {
    return component.NewIDWithName(t.factory.Type(), t.name)
}

func (t *translator) Translate(conf *confmap.Conf) (component.Config, error) {
    cfg := t.factory.CreateDefaultConfig().(*myprocessor.Config)
    // Read from conf and populate cfg fields
    return cfg, nil
}
```

2. Add the translator to the appropriate pipeline translator (e.g., in `translator/translate/otel/pipeline/host/translator.go`):

```go
import "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/processor/mycomponent"
// In the Translate() method, add to processors:
translators.Processors.Set(mycomponent.NewTranslator())
```

#### Testing

```go
func TestTranslate(t *testing.T) {
    tt := NewTranslator()
    conf := confmap.NewFromStringMap(map[string]interface{}{
        "metrics": map[string]interface{}{
            "metrics_collected": map[string]interface{}{
                "mymetric": map[string]interface{}{},
            },
        },
    })
    cfg, err := tt.Translate(conf)
    require.NoError(t, err)
    assert.NotNil(t, cfg)
}
```

Reference: [`translator/translate/otel/exporter/awscloudwatch/translator_test.go`](../../translator/translate/otel/exporter/awscloudwatch/translator_test.go).

---

### Add a New Workload Detector

> **When to use**: For detecting new application types in the workload discovery system. See [`internal/detector/tomcat/tomcat.go`](../../internal/detector/tomcat/tomcat.go) as the minimal reference.

1. Create `internal/detector/<name>/<name>.go` implementing `detector.ProcessDetector`:

```go
package myapp

import (
    "context"
    "log/slog"
    "github.com/aws/amazon-cloudwatch-agent/internal/detector"
)

type myDetector struct {
    logger *slog.Logger
}

func NewDetector(logger *slog.Logger) detector.ProcessDetector {
    return &myDetector{logger: logger}
}

func (d *myDetector) Detect(ctx context.Context, process detector.Process) (*detector.Metadata, error) {
    exe, err := process.ExeWithContext(ctx)
    if err != nil {
        return nil, err
    }
    // Check if process matches
    if !isMyApp(exe) {
        return nil, detector.ErrIncompatibleDetector
    }
    return &detector.Metadata{
        Categories: []detector.Category{"myapp"},
        Name:       "detected-name",
    }, nil
}
```

2. Register in the parent detector (e.g., `internal/detector/java/java.go` for JVM-based apps) or in `cmd/workload-discovery/discovery.go`:

```go
processDetectors: []detector.ProcessDetector{
    myapp.NewDetector(logger),
},
```

#### Testing

Use `detectortest.MockProcess`:

```go
func TestDetect(t *testing.T) {
    d := NewDetector(slog.Default())
    proc := &detectortest.MockProcess{}
    proc.On("ExeWithContext", mock.Anything).Return("/usr/bin/myapp", nil)
    proc.On("PID").Return(int32(1234))

    md, err := d.Detect(context.Background(), proc)
    require.NoError(t, err)
    assert.Equal(t, "myapp", string(md.Categories[0]))
}
```

Reference: [`internal/detector/detectortest/mock.go`](../../internal/detector/detectortest/mock.go).

---

### Add a New CLI Binary

> **When to use**: For standalone tools. See [`cmd/workload-discovery/`](../../cmd/workload-discovery/) as reference.

1. Create `cmd/<name>/main.go`:

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    flag.Parse()
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}

func run() error {
    // Implementation
    return nil
}
```

2. Add build target in `Makefile` if needed (follow existing `cmd/` patterns).

#### Testing

Test the `run()` function directly:

```go
func TestRun(t *testing.T) {
    err := run()
    assert.NoError(t, err)
}
```

---

### Add OS-Specific Code with Build Constraints

> **When to use**: When behavior differs across operating systems. Prefer file suffixes for entire-file specialization; use build tags only for complex multi-platform exclusion.

**Pattern A — File suffixes** (preferred):

```
mypackage/
├── impl.go           # Shared code
├── impl_linux.go     # Linux-specific
├── impl_windows.go   # Windows-specific
└── impl_unix.go      # Linux + Darwin (Go implicit constraint)
```

**Pattern B — Build tags** (for negation/complex constraints):

```go
//go:build !linux && !windows

package mypackage

// Stub implementation for unsupported platforms
```

Reference: [`internal/detector/nvidia/nvidia_stub.go`](../../internal/detector/nvidia/nvidia_stub.go) uses `//go:build !linux && !windows`.

**Pattern C — Non-Windows stubs** (for Windows-only plugins):

```
myplugin/
├── myplugin_windows.go      # Real implementation
└── myplugin_notwindows.go   # No-op stub
```

Reference: [`plugins/inputs/win_perf_counters/`](../../plugins/inputs/win_perf_counters/).

#### Testing

Create matching `*_<os>_test.go` files. Run cross-platform CI via `make build` which compiles for all 5 platform/arch combos.

---

### Add a Golden File Translation Test

> **When to use**: When adding or modifying JSON config translation. Tests the full JSON → TOML + YAML pipeline.

1. Create input JSON in `translator/tocwconfig/sampleConfig/<name>_<os>.json`:

```json
{
    "metrics": {
        "metrics_collected": {
            "mymetric": {}
        }
    }
}
```

2. Generate expected outputs by running the translator, then save as golden files:
   - `sampleConfig/<name>_<os>.conf` (TOML output)
   - `sampleConfig/<name>_<os>.yaml` (YAML output)

3. Add test case in `translator/tocwconfig/tocwconfig_test.go`:

```go
func TestMyMetric(t *testing.T) {
    resetContext(t)
    checkTranslation(t, "mymetric_linux", "linux", nil, "")
}
```

The `checkTranslation` helper loads the JSON, runs both translation paths, and compares against the `.conf` and `.yaml` golden files.

#### Testing

This recipe *is* a test recipe. Run with:

```bash
go test ./translator/tocwconfig/ -run TestMyMetric
```

To update golden files after intentional changes, run the translator and overwrite the expected output files.

---

## Composite Recipes

### Add a New Metric Type End-to-End

Complete workflow: input plugin → translation rule → OTel pipeline → CloudWatch exporter.

1. **Create the OTel processor** (or receiver) — follow [Add a New OTel Processor](#add-a-new-otel-processor) or [Add a New OTel Receiver](#add-a-new-otel-receiver)

2. **Register the component** in `service/defaultcomponents/components.go`

3. **Create the OTel translation rule** — follow [Add an OTel Translation Rule](#add-an-otel-translation-rule)

4. **Wire into a pipeline translator** — either add to an existing pipeline (e.g., `translator/translate/otel/pipeline/host/translator.go`) or create a new pipeline translator:

```go
// translator/translate/otel/pipeline/mymetric/translator.go
package mymetric

import (
    "go.opentelemetry.io/collector/confmap"
    "go.opentelemetry.io/collector/pipeline"
    "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
    "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/receiver/myreceiver"
    "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/exporter/awscloudwatch"
)

type translator struct{}

var _ common.PipelineTranslator = (*translator)(nil)

func NewTranslator() common.PipelineTranslator { return &translator{} }

func (t *translator) ID() pipeline.ID {
    return pipeline.NewIDWithName(pipeline.SignalMetrics, "mymetric")
}

func (t *translator) Translate(conf *confmap.Conf) (*common.ComponentTranslators, error) {
    if conf == nil || !conf.IsSet(common.ConfigKey(common.MetricsKey, common.MetricsCollectedKey, "mymetric")) {
        return nil, &common.MissingKeyError{ID: t.ID(), JsonKey: "metrics::metrics_collected::mymetric"}
    }
    return &common.ComponentTranslators{
        Receivers:  common.NewTranslatorMap(myreceiver.NewTranslator()),
        Processors: common.NewTranslatorMap[component.Config, component.ID](),
        Exporters:  common.NewTranslatorMap(awscloudwatch.NewTranslator()),
        Extensions: common.NewTranslatorMap[component.Config, component.ID](),
    }, nil
}
```

5. **Register the pipeline** in `translator/translate/otel/translate_otel.go`:

```go
import "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/pipeline/mymetric"
// In init() or Translate():
RegisterPipeline(mymetric.NewTranslator())
```

6. **Add golden file tests** — follow [Add a Golden File Translation Test](#add-a-golden-file-translation-test)

7. **Update JSON schema** if adding new config keys in `translator/config/schema.json`

---

### Add a New Log Source End-to-End

1. **Define JSON config structure** — add to `logs.logs_collected` in `translator/config/schema.json`

2. **Create OTel translation rule** mapping JSON config to a receiver (typically `filelogreceiver` or `tcplogreceiver`):

```go
// translator/translate/otel/receiver/mylogsource/translator.go
func (t *translator) Translate(conf *confmap.Conf) (component.Config, error) {
    // Map JSON config fields to filelog receiver config
    cfg := t.factory.CreateDefaultConfig()
    // Configure include paths, multiline, encoding, etc.
    return cfg, nil
}
```

3. **Create pipeline translator** under `translator/translate/otel/pipeline/` that wires receiver → processors → `awscloudwatchlogs` exporter

4. **Register pipeline** in `translator/translate/otel/translate_otel.go`

5. **Add golden file tests** with sample JSON config and expected YAML output

---

### Add Container Insights Support for a New Resource Type

1. **Add constants** in `internal/containerinsightscommon/`:
   - Metric names in `const.go`
   - Kubernetes attribute keys in `k8sconst.go` (if K8s-related)
   - Type constant with `MetricName()` mapping

2. **Create or extend a processor** to decorate metrics with the new resource type's metadata. Follow [Add a New OTel Processor](#add-a-new-otel-processor).

3. **Add translation rules** to wire the processor into the Container Insights pipeline (`translator/translate/otel/pipeline/containerinsights/`)

4. **Register in `service/defaultcomponents/components.go`** if new component

5. **Add golden file tests** with Container Insights JSON config

Reference: Kueue support added via `plugins/processors/kueueattributes/` + `translator/translate/otel/processor/kueue/` + `translator/translate/otel/pipeline/containerinsights/`.

---

## Registration Points Checklist

When adding a new component, check all applicable registration points:

| Registration Point | When to Update | File |
|---|---|---|
| OTel component factory | Any new receiver, processor, exporter, or extension | `service/defaultcomponents/components.go` |
| Telegraf plugin | New Telegraf-style input/output/processor (legacy) | `plugins/plugins.go` |
| Legacy translation rules | New TOML config field (legacy) | `translator/registerrules/register_rules.go` |
| OTel translation | New OTel pipeline or component translator | `translator/translate/otel/translate_otel.go` |
| Pipeline translator | New pipeline type | Pipeline's `init()` or explicit `RegisterPipeline()` |
| Makefile | New binary in `cmd/` | `Makefile` |
| JSON schema | New user-facing config field | `translator/config/schema.json` |
| `go.mod` | New external dependency | `go.mod` (run `go mod tidy`) |

---

## Common Failure Modes

### IMDS Metadata Fetch Failures

**Symptom**: Agent fails to detect EC2 instance metadata, falls back to `onPrem` mode.

**Mitigation**: `internal/retryer/imdsretryer.go` extends the SDK `DefaultRetryer` to retry `EC2MetadataError`. Configure retries via `IMDS_NUMBER_RETRY` env var (default: 1).

```go
// Retry pattern from internal/retryer/imdsretryer.go
func (r IMDSRetryer) ShouldRetry(req *request.Request) bool {
    if r.DefaultRetryer.ShouldRetry(req) {
        return true
    }
    // Also retry EC2MetadataError
    var ec2Err *ec2metadata.EC2MetadataError
    return errors.As(req.Error, &ec2Err)
}
```

### CloudWatch API Throttling

**Symptom**: `LimitExceededFault` or `ThrottlingException` in logs.

**Mitigation**: `plugins/outputs/cloudwatch/cloudwatch.go` retries up to 5 times with backoff. The `LogThrottleRetryer` aggregates throttle events and logs at most once per minute. Adaptive publish interval halves when buffer fills.

### Config Translation Validation Errors

**Symptom**: Agent panics at startup with schema validation error.

**Mitigation**: `translator/isValid.go` validates JSON against embedded schema. Errors include the JSON path. Fix the JSON config to match `translator/config/schema.json`.

### Plugin Initialization Failures

**Symptom**: Component `Start()` returns error or background goroutine fails silently.

**Mitigation**: Extensions log errors and set `ready` atomic flag. Check agent logs for `"failed to"` messages. Entity store and k8smetadata use retry loops with `done` channel for graceful shutdown.

### Data Race Conditions

**Symptom**: Test failures with `-race` flag in known packages.

**Mitigation**: 9 packages are excluded from race detection in `Makefile` (`PKG_WITH_DATA_RACE`). When adding concurrency to these packages, be aware of existing races. When adding a new package with known races, add it to the exclusion list.

---

## Decision Trees

### Telegraf Plugin vs Native OTel Component

```
Need to add a new data source or processor?
│
├─ Wrapping an existing upstream Telegraf plugin?
│  └─ YES → Telegraf input + blank import in plugins/plugins.go
│           (adapter receiver wraps it automatically)
│
├─ Need logs or traces support?
│  └─ YES → Native OTel component (adapter only supports metrics)
│
├─ New code from scratch?
│  └─ YES → Native OTel component (preferred for all new code)
│
└─ Modifying existing Telegraf-style plugin?
   └─ Keep as Telegraf if changes are minimal;
      consider migration to OTel for major rewrites
```

### Legacy Translation Rule vs OTel Translation Rule

```
Need to translate a new JSON config field?
│
├─ Field only affects TOML output (CWAgent.conf)?
│  └─ YES → Legacy rule in translator/translate/
│           (rare — most new features use OTel path)
│
├─ Field affects OTel YAML output?
│  └─ YES → OTel translation rule in translator/translate/otel/
│
└─ Field affects both outputs?
   └─ Create both: legacy rule for TOML + OTel translator for YAML
      (both paths run during translation)
```

### OS-Specific Files vs Build Tags

```
Need platform-specific code?
│
├─ Entire file is OS-specific?
│  └─ YES → Use file suffix: _linux.go, _windows.go, _unix.go, _darwin.go
│
├─ Need negation (e.g., "everything except Windows")?
│  └─ YES → Use _notwindows.go suffix OR //go:build !windows tag
│
├─ Complex multi-platform exclusion?
│  └─ YES → Use //go:build tag (e.g., //go:build !linux && !windows)
│           Reference: internal/detector/nvidia/nvidia_stub.go
│
└─ Only a few lines differ?
   └─ Consider runtime check (runtime.GOOS) if simple enough;
      otherwise use file suffixes
```

### Golden File Test vs Unit Test vs Integration Test

```
What are you testing?
│
├─ JSON config → OTel YAML/TOML translation?
│  └─ Golden file test in translator/tocwconfig/
│
├─ OTel component factory creates valid config?
│  └─ Factory test (factory_test.go pattern)
│
├─ Individual function or method behavior?
│  └─ Unit test in the same package
│
├─ OTel translator maps config correctly?
│  └─ Translator test with confmap.Conf input
│
├─ Full agent behavior with real AWS services?
│  └─ Integration test in amazon-cloudwatch-agent-test repo
│
└─ Concurrency correctness?
   └─ Unit test with -race flag
      (check PKG_WITH_DATA_RACE if failures are known)
```

---

## Recipe Coverage Matrix

| Architecture Doc | Primary Recipe(s) |
|---|---|
| [SCOPE.md](SCOPE.md) | (Reference doc — no recipe needed) |
| [build-and-tooling.md](build-and-tooling.md) | [Add a New CLI Binary](#add-a-new-cli-binary) |
| [project-structure.md](project-structure.md) | [Add a New OTel Processor](#add-a-new-otel-processor), [Add a New Telegraf Input Plugin](#add-a-new-telegraf-input-plugin) |
| [configuration.md](configuration.md) | [Add a Translation Rule](#add-a-translation-rule-for-a-new-json-config-field), [Add an OTel Translation Rule](#add-an-otel-translation-rule) |
| [config-translation.md](config-translation.md) | [Add an OTel Translation Rule](#add-an-otel-translation-rule), [Add a Golden File Translation Test](#add-a-golden-file-translation-test) |
| [plugins-pipeline.md](plugins-pipeline.md) | [Add a New OTel Processor](#add-a-new-otel-processor), [Add a New OTel Receiver](#add-a-new-otel-receiver), [Add a New OTel Extension](#add-a-new-otel-extension), [Add a New OTel Exporter](#add-a-new-otel-exporter) |
| [domain-model.md](domain-model.md) | [Add a New Metric Type End-to-End](#add-a-new-metric-type-end-to-end) |
| [testing.md](testing.md) | [Test-Writing Recipes](#test-writing-recipes), [Add a Golden File Translation Test](#add-a-golden-file-translation-test) |
| [concurrency.md](concurrency.md) | [Common Failure Modes: Data Race Conditions](#data-race-conditions) |
| [error-handling.md](error-handling.md) | [Common Failure Modes](#common-failure-modes) |
| [telemetry.md](telemetry.md) | [Add a New OTel Extension](#add-a-new-otel-extension) |
| [deployment.md](deployment.md) | [Add a New CLI Binary](#add-a-new-cli-binary), [Add OS-Specific Code](#add-os-specific-code-with-build-constraints) |

---

## Test-Writing Recipes

### How to Write a Factory Test for an OTel Component

Every OTel component (processor, receiver, extension, exporter) should have a `factory_test.go`. The pattern is consistent across all 16 existing factory tests.

```go
// factory_test.go
package mycomponent

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/component/componenttest"
    "go.opentelemetry.io/collector/consumer/consumertest"
    "go.opentelemetry.io/collector/pipeline"
    "go.opentelemetry.io/collector/processor/processortest"
)

func TestCreateDefaultConfig(t *testing.T) {
    factory := NewFactory()
    require.NotNil(t, factory)
    cfg := factory.CreateDefaultConfig()
    assert.NotNil(t, cfg)
    assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateProcessor(t *testing.T) {
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    set := processortest.NewNopSettings(component.MustNewType("mycomponent"))

    // Test supported signal
    mp, err := factory.CreateMetrics(context.Background(), set, cfg, consumertest.NewNop())
    assert.NoError(t, err)
    assert.NotNil(t, mp)

    // Test unsupported signals return ErrSignalNotSupported
    tp, err := factory.CreateTraces(context.Background(), set, cfg, consumertest.NewNop())
    assert.Equal(t, err, pipeline.ErrSignalNotSupported)
    assert.Nil(t, tp)

    lp, err := factory.CreateLogs(context.Background(), set, cfg, consumertest.NewNop())
    assert.Equal(t, err, pipeline.ErrSignalNotSupported)
    assert.Nil(t, lp)
}
```

Key helpers:
- `componenttest.CheckConfigStruct(cfg)` — validates config struct tags
- `processortest.NewNopSettings` / `receivertest.NewNopSettings` / `extensiontest.NewNopSettings` — create test settings
- `consumertest.NewNop()` — no-op consumer for downstream
- `pipeline.ErrSignalNotSupported` — expected error for unsupported signal types

### How to Write a Golden File Translation Test

1. Create the JSON input file: `translator/tocwconfig/sampleConfig/myfeature_linux.json`
2. Run the translator to generate expected outputs, save as `.conf` and `.yaml`
3. Add the test function:

```go
func TestMyFeature(t *testing.T) {
    resetContext(t)
    checkTranslation(t, "myfeature_linux", "linux", nil, "")
}
```

`checkTranslation` does:
- Loads `sampleConfig/myfeature_linux.json`
- Runs legacy translation → compares with `sampleConfig/myfeature_linux.conf`
- Runs OTel translation → compares with `sampleConfig/myfeature_linux.yaml`

For tests needing environment variables:

```go
func TestMyFeatureWithEnv(t *testing.T) {
    resetContext(t)
    expectedEnvVars := map[string]string{"MY_VAR": "value"}
    checkTranslation(t, "myfeature_linux", "linux", expectedEnvVars, "")
}
```

For platform-specific golden files, use the `appendString` parameter to select a suffix.

### How to Write a Processor Unit Test

Test the `processMetrics` function directly with constructed `pmetric.Metrics`:

```go
func TestProcessMetrics(t *testing.T) {
    cfg := &Config{}
    logger, _ := zap.NewDevelopment()
    p := &myProcessor{config: cfg, logger: logger}

    md := pmetric.NewMetrics()
    rm := md.ResourceMetrics().AppendEmpty()
    sm := rm.ScopeMetrics().AppendEmpty()
    m := sm.Metrics().AppendEmpty()
    m.SetName("test_metric")
    m.SetEmptyGauge().DataPoints().AppendEmpty().SetDoubleValue(42.0)

    result, err := p.processMetrics(context.Background(), md)
    require.NoError(t, err)
    // Assert on result metrics
    assert.Equal(t, 1, result.MetricCount())
}
```

For processors that depend on extensions (e.g., entity store), use function injection:

```go
// In processor.go:
var getEntityStore = entitystore.GetEntityStore

// In processor_test.go:
getEntityStore = func() *entitystore.EntityStore { return mockStore }
```

Reference: [`plugins/processors/awsentity/processor_test.go`](../../plugins/processors/awsentity/processor_test.go).

### How to Add a Package to the Data Race Exclusion List

If a package has known data races that cannot be immediately fixed:

1. Add the package path to `PKG_WITH_DATA_RACE` in `Makefile`:

```makefile
PKG_WITH_DATA_RACE = \
    ... \
    github.com/aws/amazon-cloudwatch-agent/my/package
```

2. The `test-data-race` target automatically excludes these packages from `go test -race`.

3. Document the race condition in a comment or issue for future resolution.

> **Note**: This is a mitigation, not a fix. Prefer fixing the race condition when possible.

---

## Boilerplate Templates

### OTel Processor (factory.go + config.go + processor.go)

```
plugins/processors/{name}/
├── config.go
├── factory.go
├── processor.go
└── factory_test.go
```

**config.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import "go.opentelemetry.io/collector/component"

type Config struct {
    // TODO: Add fields with mapstructure tags
}

var _ component.Config = (*Config)(nil)

func (cfg *Config) Validate() error { return nil }
```

**factory.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "context"
    "fmt"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/processor"
    "go.opentelemetry.io/collector/processor/processorhelper"
)

var (
    TypeStr, _ = component.NewType("{name}")
    processorCapabilities = consumer.Capabilities{MutatesData: true}
)

func NewFactory() processor.Factory {
    return processor.NewFactory(TypeStr,
        func() component.Config { return &Config{} },
        processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha))
}

func createMetricsProcessor(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
    c, ok := cfg.(*Config)
    if !ok {
        return nil, fmt.Errorf("configuration parsing error")
    }
    return processorhelper.NewMetrics(ctx, set, cfg, next,
        newProcessor(c, set.Logger).processMetrics,
        processorhelper.WithCapabilities(processorCapabilities))
}
```

**processor.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "context"

    "go.opentelemetry.io/collector/pdata/pmetric"
    "go.uber.org/zap"
)

type {name}Processor struct {
    config *Config
    logger *zap.Logger
}

func newProcessor(cfg *Config, logger *zap.Logger) *{name}Processor {
    return &{name}Processor{config: cfg, logger: logger}
}

func (p *{name}Processor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
    // TODO: Implement metric processing
    return md, nil
}
```

### OTel Extension (factory.go + config.go + extension.go)

```
extension/{name}/
├── config.go
├── factory.go
├── extension.go
└── factory_test.go
```

**config.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import "go.opentelemetry.io/collector/component"

type Config struct {
    // TODO: Add fields with mapstructure tags
}

var _ component.Config = (*Config)(nil)
```

**factory.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/extension"
)

var TypeStr, _ = component.NewType("{name}")

func NewFactory() extension.Factory {
    return extension.NewFactory(TypeStr,
        func() component.Config { return &Config{} },
        func(_ context.Context, set extension.Settings, cfg component.Config) (extension.Extension, error) {
            return &{name}Extension{config: cfg.(*Config), logger: set.Logger}, nil
        },
        component.StabilityLevelAlpha)
}
```

**extension.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/extension"
    "go.uber.org/zap"
)

type {name}Extension struct {
    config *Config
    logger *zap.Logger
}

var _ extension.Extension = (*{name}Extension)(nil)

func (e *{name}Extension) Start(_ context.Context, _ component.Host) error {
    // TODO: Initialize resources
    return nil
}

func (e *{name}Extension) Shutdown(_ context.Context) error {
    // TODO: Clean up resources
    return nil
}
```

### OTel Receiver (factory.go + config.go + scraper.go)

```
receiver/{name}/
├── config.go
├── factory.go
├── scraper.go
└── factory_test.go
```

**config.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/scraper/scraperhelper"
)

type Config struct {
    scraperhelper.ControllerConfig `mapstructure:",squash"`
    // TODO: Add fields
}

var _ component.Config = (*Config)(nil)
```

**factory.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver"
    otelscraper "go.opentelemetry.io/collector/scraper"
    "go.opentelemetry.io/collector/scraper/scraperhelper"
)

var TypeStr, _ = component.NewType("{name}")

func NewFactory() receiver.Factory {
    return receiver.NewFactory(TypeStr,
        func() component.Config {
            return &Config{ControllerConfig: scraperhelper.NewDefaultControllerConfig()}
        },
        receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelAlpha))
}

func createMetricsReceiver(_ context.Context, set receiver.Settings, cfg component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
    c := cfg.(*Config)
    s := newScraper(c, set.Logger)
    scraper, err := otelscraper.NewMetrics(s.scrape, otelscraper.WithStart(s.start), otelscraper.WithShutdown(s.shutdown))
    if err != nil {
        return nil, err
    }
    return scraperhelper.NewMetricsController(&c.ControllerConfig, set, consumer,
        scraperhelper.AddScraper(TypeStr, scraper))
}
```

**scraper.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/pdata/pmetric"
    "go.uber.org/zap"
)

type scraper struct {
    config *Config
    logger *zap.Logger
}

func newScraper(cfg *Config, logger *zap.Logger) *scraper {
    return &scraper{config: cfg, logger: logger}
}

func (s *scraper) start(_ context.Context, _ component.Host) error { return nil }
func (s *scraper) shutdown(_ context.Context) error { return nil }

func (s *scraper) scrape(_ context.Context) (pmetric.Metrics, error) {
    md := pmetric.NewMetrics()
    // TODO: Populate metrics
    return md, nil
}
```

### Telegraf Input Plugin

```
plugins/inputs/{name}/
├── {name}.go
└── {name}_test.go
```

**{name}.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
)

type {Name} struct {
    // TODO: Add config fields
}

func (p *{Name}) SampleConfig() string { return "" }

func (p *{Name}) Gather(acc telegraf.Accumulator) error {
    // TODO: Collect metrics
    // acc.AddFields("measurement", fields, tags)
    return nil
}

func init() {
    inputs.Add("{name}", func() telegraf.Input {
        return &{Name}{}
    })
}
```

### OTel Translation Rule

```
translator/translate/otel/{type}/{name}/
├── translator.go
└── translator_test.go
```

**translator.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/confmap"
    "go.opentelemetry.io/collector/{type}"

    "github.com/aws/amazon-cloudwatch-agent/{component_package}"
    "github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

type translator struct {
    name    string
    factory {type}.Factory
}

var _ common.ComponentTranslator = (*translator)(nil)

func NewTranslator() common.ComponentTranslator {
    return NewTranslatorWithName("")
}

func NewTranslatorWithName(name string) common.ComponentTranslator {
    return &translator{name, {component_package}.NewFactory()}
}

func (t *translator) ID() component.ID {
    return component.NewIDWithName(t.factory.Type(), t.name)
}

func (t *translator) Translate(conf *confmap.Conf) (component.Config, error) {
    cfg := t.factory.CreateDefaultConfig()
    // TODO: Read from conf and populate cfg
    return cfg, nil
}
```

**translator_test.go**:
```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package {name}

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.opentelemetry.io/collector/confmap"
)

func TestTranslate(t *testing.T) {
    tt := NewTranslator()
    conf := confmap.NewFromStringMap(map[string]interface{}{
        // TODO: Add test config
    })
    cfg, err := tt.Translate(conf)
    require.NoError(t, err)
    assert.NotNil(t, cfg)
}
```
