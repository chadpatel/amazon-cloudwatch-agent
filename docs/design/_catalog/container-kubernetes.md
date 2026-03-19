# Container & Kubernetes

Designs related to Container Insights, ECS/EKS decorators, and Kubernetes metadata.

## Designs

| Title | Document | Status | Description |
|---|---|---|---|
| Container Insights Base | `internal/containerinsightscommon/`, `translator/translate/otel/pipeline/containerinsights/`, `translator/translate/otel/receiver/awscontainerinsight/` | 🚀 Implemented (Jun 2020 / Oct 2022) | Shared constants, metric names, K8s/ECS attribute keys; OTel pipeline and receiver translation for Container Insights. Predates design practice — no design doc available |
| Kubernetes Decorator Processor | `plugins/processors/k8sdecorator/` | 🚀 Implemented (Jun 2020) | Telegraf-style processor decorating metrics with K8s pod/node metadata from API server. Predates design practice — no design doc available |
| ECS Decorator Processor | `plugins/processors/ecsdecorator/` | 🚀 Implemented (Jun 2020) | Telegraf-style processor decorating metrics with ECS task/container metadata from cgroup and ECS metadata endpoint. Predates design practice — no design doc available |
| K8s Metadata Extension | `extension/k8smetadata/` | 🚀 Implemented (Mar 2025) | OTel extension providing Kubernetes pod and node metadata to processors via singleton accessor. Predates design practice — no design doc available |
| Kueue Support | `plugins/processors/kueueattributes/`, `translator/translate/otel/receiver/awscontainerinsightskueue/`, `translator/translate/otel/processor/kueue/` | 🚀 Implemented (Nov 2024) | Kueue workload queue attribute injection and Container Insights Kueue receiver/processor translation. Predates design practice — no design doc available |
| GPU Container Insights | `plugins/processors/gpuattributes/`, `translator/translate/otel/processor/gpu/` | 🚀 Implemented (Mar 2024) | GPU metric attribute normalization for NVIDIA and Neuron accelerators in Container Insights. Predates design practice — no design doc available |
| Container Insights JMX | `translator/translate/otel/pipeline/containerinsightsjmx/` | 🚀 Implemented (Oct 2024) | OTel pipeline translation for JMX metrics collection within Container Insights. Predates design practice — no design doc available |
