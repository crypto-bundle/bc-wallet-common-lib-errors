# Docker — Test + Lint Container

## Purpose

Run `go test` and `make lint` in an isolated, reproducible container. Useful for local parity with CI environments.

## Prerequisites

- Docker Engine 20.10+
- Internet access during first `docker build` (to download Go toolchain + golangci-lint + tinyerrors dep)

## Quick Start

```bash
# Build the test image (takes ~30–60s on first run)
docker build -f Dockerfile.test -t bc-wallet-errors:test .

# Run tests + lint (default CMD)
docker run --rm bc-wallet-errors:test

# Override default command
docker run --rm bc-wallet-errors:test go test -v ./pkg/errformatter/...
```

## Image Details

| Property | Value |
|---|---|
| Base image | `golang:1.27-alpine` |
| Linter | golangci-lint v1.64.6 |
| Default CMD | `go test ./pkg/errformatter/... && make lint` |
| Exposed ports | None (test-only, no server) |

## Notes & Limitations

- **No scratch variant** — this library has an external dependency (`tinyerrors`), so a minimal scratch image wouldn't prove zero-dependency deployment. A scratch image can be added later if needed.
- **Requires network** on first build to download Go deps and golangci-lint. Subsequent builds reuse layer cache.
- **Not for production runtime** — this image contains the full Go compiler and is meant only for test/lint execution.
- **No docker-compose.yml** — deferred until multi-container workflows are required.
