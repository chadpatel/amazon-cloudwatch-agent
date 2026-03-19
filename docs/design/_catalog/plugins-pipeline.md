# Plugins & Pipeline

Design documents for OTel pipeline components, the Telegraf adapter receiver, processors (ec2tagger, k8sdecorator, awsentity, awsapplicationsignals), exporters (cloudwatch, cloudwatchlogs, EMF, X-Ray), and extensions (agenthealth, entitystore, k8smetadata, server).

## Designs

| Title | Path | Status | Summary |
|---|---|---|---|
| Telegraf Adapter Receiver | `receiver/adapter/` | 🚀 Implemented | Bridges Telegraf input plugins into OTel receiver pipeline via `AdaptedReceiver` with start/scrape/shutdown lifecycle. Predates design practice — no design doc available. (Oct 2022) |
| EC2 Tagger Processor | `plugins/processors/ec2tagger/` | 🚀 Implemented | Decorates metrics with EC2 instance tags and EBS volume metadata. Background refresh with jittered polling. Predates design practice — no design doc available. (Jul 2020) |
| Kubernetes Decorator | `plugins/processors/k8sdecorator/` | 🚀 Implemented | Telegraf-style processor decorating metrics with Kubernetes pod/node metadata via k8s API watchers. Predates design practice — no design doc available. (Jun 2020) |
| ECS Decorator | `plugins/processors/ecsdecorator/` | 🚀 Implemented | Telegraf-style processor decorating metrics with ECS container metadata from cgroup and task metadata endpoint. Predates design practice — no design doc available. (Jun 2020) |
| AWS Entity Processor | `plugins/processors/awsentity/` | 🚀 Implemented | Injects AWS entity attributes (service name, environment, platform) into telemetry for CloudWatch entity association. (Jul 2024) |
| Application Signals Processor | `plugins/processors/awsapplicationsignals/` | 🚀 Implemented | Processes Application Signals metrics and traces with rules engine, cardinality control, metric normalization, and resolver-based attribute enrichment. (May 2024) |
| GPU Attributes Processor | `plugins/processors/gpuattributes/` | 🚀 Implemented | Normalizes GPU metric attributes for Neuron and NVIDIA sources with aggregation and modifier sub-packages. (Mar 2024) |
| Kueue Attributes Processor | `plugins/processors/kueueattributes/` | 🚀 Implemented | Injects Kueue workload queue attributes into Container Insights metrics. Minimal OTel processor (factory + config + processor). (Nov 2024) |
| Prometheus Adapter Processor | `plugins/processors/prometheusadapter/` | 🚀 Implemented | Adapts Prometheus metric types with identity tracking and delta-to-cumulative calculation. (Jun 2025) |
| Rollup Processor | `processor/rollupprocessor/` | 🚀 Implemented | Aggregates metrics by rolling up dimension combinations with configurable attribute keys and cache. (Oct 2024) |
| CloudWatch Exporter | `plugins/outputs/cloudwatch/` | 🚀 Implemented | Dual Telegraf Output + OTel exporter for CloudWatch PutMetricData with multi-stage channel pipeline, duration aggregation, and adaptive publish interval. Predates design practice — no design doc available. (Jun 2020) |
| CloudWatch Logs Exporter | `plugins/outputs/cloudwatchlogs/` | 🚀 Implemented | Ships logs to CloudWatch Logs PutLogEvents with internal pusher (batch, retry, target management, auto-create log groups/streams). Predates design practice — no design doc available. (Jun 2020) |
