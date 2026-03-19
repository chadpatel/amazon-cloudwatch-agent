# Development Commands

## Build

```bash
# Full multi-platform build (Linux amd64+arm64, Windows amd64, macOS amd64+arm64)
make build

# Build + package (RPM, DEB, ZIP, tar.gz)
make release

# Docker build (linux amd64+arm64)
make docker-build

# Docker build single arch
make docker-build-amd64
make docker-build-arm64

# Build only linux amd64 for quick Docker iteration
make build-for-docker-amd64

# Clean build artifacts
make clean
```

## Test

```bash
# Unit tests with coverage
make test
# Equivalent to: CGO_ENABLED=0 go test -timeout 15m -coverprofile coverage.txt -failfast ./...

# Unit tests with race detector (excludes known-racy packages)
make test-data-race
# Equivalent to: CGO_ENABLED=1 go test -timeout 15m -race -parallel 4 ./...

# Run tests for a single package
CGO_ENABLED=0 go test -timeout 15m -failfast ./plugins/outputs/cloudwatch/...
```

## Lint

```bash
# Full lint (golangci-lint v2.4.0 + license check + import order)
make lint

# Quick lint (license check + import order only, no golangci-lint)
make simple-lint

# License header check only
make checklicense

# Import order check only
make impi
```

Linters enabled in `.golangci.yml`: gosec, misspell, nonamedreturns, revive.

## Format

```bash
# Format code (goimports + go fmt + license headers)
make fmt

# Format shell scripts
make fmt-sh
```

## Dev Server

Not applicable — this is a CLI agent binary, not a server application.
