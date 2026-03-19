# Prometheus Histogram Support — CloudWatch PMD Path

## Goal

Enable Prometheus classic histograms to flow correctly through the CloudWatch PutMetricData (PMD) path. The OTel prometheus receiver migration (branch `rebase-prom-migration`) already makes histograms visible to the pipeline, but the CloudWatch output plugin converts them incorrectly. This plan fixes the conversion and adds tests proving both backward compatibility and new histogram support.

## Context

- GitHub Issue: https://github.com/aws/amazon-cloudwatch-agent/issues/593
- Branch: `rebase-prom-migration` (OTel prometheus receiver migration, rebased onto main)
- The EMF/CloudWatch Logs path already handles histograms correctly via the `awsemfexporter`
- The AMP path already handles histograms correctly via `prometheusremotewrite`
- The CloudWatch PMD path is broken: `ConvertFromOtel` in SEH1/Regular distributions treats bucket boundaries as literal values instead of bucket edges

### Why the Current Conversion Is Wrong

OTel histograms have `ExplicitBounds = [5, 10]` and `BucketCounts = [600, 10, 387]` meaning:
- 600 observations ≤ 5
- 10 observations in (5, 10]
- 387 observations > 10

The current `SEH1Distribution.ConvertFromOtel` does `sd.buckets[int16(k)] = float64(v)` — it casts the boundary to `int16` and uses it as a bucket key. This produces garbage data.

## Prerequisites

The contrib fork cherry-pick (Phase 0, task 0.1) must succeed before the go.mod update (task 0.2) can reference the new commit hash. After task 0.2, Ralph should run `go mod tidy && go mod vendor` to update vendored dependencies before proceeding to Phase 1.

## Reference Files

- `/local/home/patchad/workplace/cwa/opentelemetry-collector-contrib/receiver/prometheusreceiver/config.go` — Contrib fork prometheus receiver config, where EscapedCaptureGroupOne needs to be added
- `plugins/outputs/cloudwatch/convert_otel.go` — Current (broken) histogram conversion via `ClassicDistribution.ConvertFromOtel`
- `plugins/outputs/cloudwatch/aggregator.go` — Aggregation logic, `aggregationDatum` struct
- `plugins/outputs/cloudwatch/cloudwatch.go` — `BuildMetricDatum`, `buildMetricDatumDist` pattern to follow
- `plugins/outputs/cloudwatch/util.go` — `setNewDistributionFunc`, `payload` sizing
- `metric/distribution/distribution.go` — `Distribution`, `ClassicDistribution`, `ExponentialDistribution` interfaces
- `metric/distribution/seh1/seh1_distribution.go` — SEH1 distribution (existing, broken for histograms)
- `metric/distribution/exph/exph.go` — Exponential histogram distribution (already on branch)
- `plugins/processors/prometheusadapter/processor.go` — Passes histograms through with attribute decoration
- `plugins/outputs/cloudwatch/convert_otel_test.go` — Existing test patterns for `createTestHistogram`

## Code Patterns

### Histogram Conversion Algorithm

The conversion maps OTel histogram buckets to CloudWatch values/counts pairs. Each outer bucket is subdivided into inner buckets (up to 10) using density-aware distribution. This algorithm was written by dricross in the contrib fork (commit `3b43f3b3b3` on `contributing/classichistograms`).

Key types:

```go
// CloudWatchHistogram holds the converted histogram data ready for PutMetricData.
type CloudWatchHistogram struct {
    maximum     float64
    minimum     float64
    sampleCount float64
    sum         float64
    values      []float64
    counts      []float64
}

func (h *CloudWatchHistogram) ValuesAndCounts() ([]float64, []float64) {
    return h.values, h.counts
}

func (h *CloudWatchHistogram) Maximum() float64  { return h.maximum }
func (h *CloudWatchHistogram) Minimum() float64  { return h.minimum }
func (h *CloudWatchHistogram) Sum() float64       { return h.sum }
func (h *CloudWatchHistogram) SampleCount() float64 { return h.sampleCount }
```

### Histogram Validation

```go
func CheckValidity(dp pmetric.HistogramDataPoint) error {
    // Check: len(BucketCounts) == len(ExplicitBounds) + 1
    // Check: min <= max
    // Check: no NaN/Inf in bounds, min, max, sum
    // Check: boundaries in ascending order
}
```

### Histogram Aggregation

When multiple histogram data points arrive in the same aggregation window with identical boundaries, merge them:

```go
func AddHistograms(state, dp pmetric.HistogramDataPoint) {
    // If boundaries differ: replace state with dp (can't merge different schemas)
    // If boundaries match: add bucket counts element-wise, sum counts, merge min/max/sum
}
```

### aggregationDatum with Histogram Field

```go
type aggregationDatum struct {
    cloudwatch.MetricDatum
    aggregationInterval time.Duration
    distribution        distribution.Distribution
    histogram           *pmetric.HistogramDataPoint  // NEW: for native OTel histograms
    entity              cloudwatch.Entity
}
```

### buildMetricDatumHist Pattern

Follow the existing `buildMetricDatumDist` pattern but use `ConvertOTelToCloudWatch` instead of `Distribution.ValuesAndCounts()`:

```go
func (c *CloudWatch) buildMetricDatumHist(metric *aggregationDatum, dimensionsList [][]*cloudwatch.Dimension) []*cloudwatch.MetricDatum {
    cwhist := ConvertOTelToCloudWatch(*metric.histogram)
    values, counts := cwhist.ValuesAndCounts()
    // Split into batches of MaxValuesPerDatum
    // For each batch, create datum with StatisticSet + Values + Counts
    // Expand across dimensionsList (same as buildMetricDatumDist)
}
```

### convertOtelHistogramDataPoints — Fixed Version

```go
func convertOtelHistogramDataPoints(
    dataPoints pmetric.HistogramDataPointSlice,
    name string, unit string, scale float64,
    entity cloudwatch.Entity,
) []*aggregationDatum {
    datums := make([]*aggregationDatum, 0, dataPoints.Len())
    for i := 0; i < dataPoints.Len(); i++ {
        dp := dataPoints.At(i)
        if err := CheckValidity(dp); err != nil {
            log.Printf("W! dropping invalid histogram datapoint for metric %s: %v", name, err)
            continue
        }
        // ... build aggregationDatum with histogram field set (NOT distribution) ...
        ad.histogram = &dp
        datums = append(datums, &ad)
    }
    return datums
}
```

### Backward Compat Test Pattern

```go
func TestConvertOtelMetric_GaugeUnchanged(t *testing.T) {
    // Create gauge metric with known values
    // Convert via ConvertOtelMetric
    // Assert dimensions, values, timestamps match expected
}
```

## Tasks

### Phase 0: Unblock Build — Contrib Fork + go.mod

- [x] **0.1** Update `/local/home/patchad/workplace/cwa/opentelemetry-collector-contrib/receiver/prometheusreceiver/config.go` — Cherry-pick the `EscapedCaptureGroupOne` constant and `preprocessPrometheusConfig` function from commit `f8d9489763` on the `contributing/replace-telegraf-prom-plugin` branch. Add `const EscapedCaptureGroupOne = "__capture_group_1__"` after the imports. Add `preprocessPrometheusConfig(in)` call at the start of `unmarshalYAML`. Add the `preprocessPrometheusConfig` function that recursively walks the config map and replaces any `"replacement"` field value equal to `EscapedCaptureGroupOne` with `"$1"`. Reference: run `git show f8d9489763 -- receiver/prometheusreceiver/config.go` in `/local/home/patchad/workplace/cwa/opentelemetry-collector-contrib` to see the exact diff.

- [x] **0.2** Update `go.mod` — Change the `prometheusreceiver` replace directive to point to the local contrib fork checkout using a filesystem replace: `github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver => /local/home/patchad/workplace/cwa/opentelemetry-collector-contrib/receiver/prometheusreceiver`. Then run `go mod tidy && go mod vendor` and verify the build compiles by running `make build`. If the build fails due to other stale module versions, update those replace directives similarly to use local paths.

### Phase 1: Histogram Conversion Logic

- [x] **1.1** Create `plugins/outputs/cloudwatch/histogram.go` — Add `CloudWatchHistogram` type with `ValuesAndCounts`/`Maximum`/`Minimum`/`Sum`/`SampleCount` methods, `ConvertOTelToCloudWatch(dp pmetric.HistogramDataPoint) *CloudWatchHistogram` function implementing the bucket-to-values/counts mapping algorithm (subdivide each outer bucket into up to 10 inner buckets using density-aware distribution: inverse quadratic when density decreasing, uniform when stable, quadratic when increasing), `CheckValidity(dp pmetric.HistogramDataPoint) error` function that validates bucket count length matches bounds+1, min≤max, no NaN/Inf, ascending boundaries, and `AddHistograms(state, dp pmetric.HistogramDataPoint)` function that merges two histogram data points with matching boundaries by adding bucket counts element-wise and merging min/max/sum/count. Port the algorithm from contrib fork commit `3b43f3b3b3` on `contributing/classichistograms` branch of `amazon-contributing/opentelemetry-collector-contrib`.

- [x] **1.2** Create `plugins/outputs/cloudwatch/histogram_test.go` — Add unit tests for `CheckValidity` (valid histogram, mismatched bucket/bounds lengths, NaN values, non-ascending boundaries, min>max), `ConvertOTelToCloudWatch` (single bucket no boundaries, standard 3-bucket histogram with bounds `[5,10]` and counts `[600,10,387]`, empty histogram, histogram without min/max), and `AddHistograms` (matching boundaries merge, different boundaries reset, sum/count accumulation). Use `pmetric.NewHistogramDataPoint()` to construct test data.

### Phase 2: Wire Into Pipeline

- [x] **2.1** Update `plugins/outputs/cloudwatch/aggregator.go` — Add `histogram *pmetric.HistogramDataPoint` field to `aggregationDatum` struct. In `durationAggregator.aggregating()`, update the first-entry logic: if `m.histogram != nil`, store it directly without creating a distribution. Update the existing-entry logic: if `m.histogram != nil`, call `AddHistograms(*aggregatedMetric.histogram, *m.histogram)` to merge; else if `m.distribution != nil`, call `AddDistribution`; else add scalar entry. Add import for `go.opentelemetry.io/collector/pdata/pmetric`.

- [x] **2.2** Update `plugins/outputs/cloudwatch/convert_otel.go` — In `convertOtelHistogramDataPoints`, replace the current `ClassicDistribution.ConvertFromOtel` approach: instead of creating a distribution, call `CheckValidity(dp)` and skip invalid data points with a warning log, then set `ad.histogram = &dp` on the aggregation datum. Remove the `distribution.NewClassicDistribution()` and `classic.ConvertFromOtel(dp, unit)` calls from this function. Keep `convertOtelExponentialHistogramDataPoints` unchanged (it correctly uses `exph.NewExponentialDistribution()`).

- [x] **2.3** Update `plugins/outputs/cloudwatch/cloudwatch.go` — Add `buildMetricDatumHist(metric *aggregationDatum, dimensionsList [][]*cloudwatch.Dimension) []*cloudwatch.MetricDatum` method on `CloudWatch`. It should call `ConvertOTelToCloudWatch(*metric.histogram)` to get values/counts, then split into batches of `c.config.MaxValuesPerDatum`, creating a `StatisticSet` (max, min, sum, sampleCount) and `MetricDatum` for each batch × dimension combination. Only assign `Sum` on the first batch to avoid double-counting. In `BuildMetricDatum`, add a check: if `metric.histogram != nil`, call `buildMetricDatumHist`; else if `metric.distribution != nil`, call `buildMetricDatumDist`; else handle scalar (existing logic).

### Phase 3: Tests

- [x] **3.1** Update `plugins/outputs/cloudwatch/convert_otel_test.go` — Add `TestConvertOtelHistogramDataPoints_Valid` (3-bucket histogram produces aggregationDatum with histogram field set, not distribution), `TestConvertOtelHistogramDataPoints_Invalid` (mismatched bounds/counts is dropped with warning), `TestConvertOtelMetric_GaugeUnchanged` (gauge metrics still produce correct scalar datums — backward compat), `TestConvertOtelMetric_SumUnchanged` (sum/counter metrics still produce correct scalar datums — backward compat), `TestConvertOtelMetric_ExponentialHistogramUnchanged` (exponential histograms still use exph distribution — backward compat). Use existing `createTestHistogram` helper, updating its test data to use valid OTel histogram format (bounds `[5,10]`, counts `[600,10,387]` — 3 counts for 2 bounds).

- [x] **3.2** Update `plugins/outputs/cloudwatch/cloudwatch_test.go` — Add `TestBuildMetricDatumHist` (histogram datum has correct Values, Counts, StatisticValues with Max/Min/Sum/SampleCount), `TestBuildMetricDatumHist_SplitLargeHistogram` (histogram with >MaxValuesPerDatum values/counts is split correctly, Sum only on first batch), `TestBuildMetricDatumHist_WithRollup` (histogram expands across rollup dimensions), `TestBuildMetricDatum_BackwardCompat_Distribution` (distribution metrics still produce correct datums via buildMetricDatumDist), `TestBuildMetricDatum_BackwardCompat_Scalar` (scalar metrics still produce correct datums).

### Phase 4: Catalog

- [x] **4.1** Update `docs/design/_catalog/metrics-collection.md` — Change status of "Prometheus Histogram Support" entry from `🟡 In Progress` to `✅ Done`

## File Changes

| File | Change | Task |
|------|--------|------|
| `/local/home/patchad/workplace/cwa/opentelemetry-collector-contrib/receiver/prometheusreceiver/config.go` | Add EscapedCaptureGroupOne constant + preprocessPrometheusConfig | 0.1 |
| `go.mod` | Point prometheusreceiver to local contrib fork, go mod tidy + vendor | 0.2 |
| `plugins/outputs/cloudwatch/histogram.go` | New — conversion algorithm, validation, aggregation helper | 1.1 |
| `plugins/outputs/cloudwatch/histogram_test.go` | New — unit tests for conversion/validation/aggregation | 1.2 |
| `plugins/outputs/cloudwatch/aggregator.go` | Add histogram field to aggregationDatum, histogram-aware aggregation | 2.1 |
| `plugins/outputs/cloudwatch/convert_otel.go` | Use histogram field instead of broken ClassicDistribution.ConvertFromOtel | 2.2 |
| `plugins/outputs/cloudwatch/cloudwatch.go` | Add buildMetricDatumHist, update BuildMetricDatum routing | 2.3 |
| `plugins/outputs/cloudwatch/convert_otel_test.go` | Add histogram + backward compat tests | 3.1 |
| `plugins/outputs/cloudwatch/cloudwatch_test.go` | Add histogram datum + backward compat tests | 3.2 |
| `docs/design/_catalog/metrics-collection.md` | Update status to ✅ Done | 4.1 |
