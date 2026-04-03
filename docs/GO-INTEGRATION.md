# Go — CI & Docker Reference

> You already know multi-stage Docker and GitHub Actions. Go makes both simpler because it compiles to a single static binary with zero runtime dependencies.

---

## Why Go's Docker Story Is Different

In Bun/Node, your image needs the runtime installed to execute your code. In Go, you compile first and the output is a self-contained binary. The image only needs to hold that one file.

This means:

- Final image can be `scratch` (literally empty) or `distroless` (minimal base, better for debugging)
- No `node_modules`, no `bun install`, no package manager in production
- Final image size: typically 5–20MB vs 100–300MB for a Bun app

---

## Dockerfile

The standard pattern. Two stages: build the binary, copy it into a minimal image.

```dockerfile
# ---- Stage 1: Build ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency manifests first (layer cache — same trick as your bun install layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and compile
COPY . .

# CGO_ENABLED=0 — static binary, no C dependencies
# GOOS=linux — target OS (important if building on Mac)
# -o /app/server — output binary name and location
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go


# ---- Stage 2: Runtime ----
FROM scratch

# Copy only the binary from builder
COPY --from=builder /app/server /server

# If you need TLS (HTTPS outbound calls), you need CA certs — scratch has none
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 3000

CMD ["/server"]
```

### `scratch` vs `distroless` vs `alpine`

| Base | Size | Shell access | Use when |
|---|---|---|---|
| `scratch` | ~0MB | None | Purely serving HTTP, no outbound TLS |
| `gcr.io/distroless/static` | ~2MB | None | Want `scratch` but with CA certs baked in |
| `alpine` | ~5MB | sh | Need to debug in container, run shell commands |

If your app makes HTTPS calls to external services (e.g. Redis, third-party APIs), use `distroless/static` — it includes CA certs so TLS handshakes work. `scratch` has nothing.

### If you need timezone data

```dockerfile
FROM golang:1.24-alpine AS builder
# ...same as above

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/server /server
CMD ["/server"]
```

---

## docker-compose.yml

```yaml
services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      PORT: 3000
```

No volume mounts for `node_modules` to worry about. If you need a dev setup with hot reload via `air`, that's a separate compose target:

```yaml
services:
  api:
    image: golang:1.24-alpine
    working_dir: /app
    volumes:
      - .:/app
    ports:
      - "3000:3000"
    command: sh -c "go install github.com/air-verse/air@latest && air"
```

---

## GitHub Actions CI

```yaml
name: CI

on:
  push:
    branches: ['**']
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'        # pin your version
          cache: true               # caches go module downloads automatically

      - name: Install dependencies
        run: go mod download

      - name: Type check / vet
        run: go vet ./...           # Go's built-in static analysis, equivalent to tsc --noEmit

      - name: Run tests
        run: go test ./...

      - name: Run tests with race detector
        run: go test -race ./...    # catches concurrent map writes and similar bugs
```

That's it. No separate lint install, no Biome setup. `go vet` ships with Go.

### With coverage (equivalent to your `--coverage` flag)

```yaml
      - name: Run tests with coverage
        run: go test -coverprofile=coverage.out ./...

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v5
        with:
          files: coverage.out
```

### With build verification

```yaml
      - name: Build binary
        run: go build -o /dev/null ./...   # /dev/null discards output, just verifies it compiles
```

---

## The `-race` Flag

Worth calling out separately because it's unique to Go.

```bash
go test -race ./...
```

Go's race detector instruments your binary to catch concurrent access to shared memory — things like two goroutines writing to the same map without a mutex. It's not a linter, it's a runtime check that runs during your test suite.

In CI: run it on pushes to main/dev branches. It adds ~20% overhead, acceptable for CI but too slow for local watch mode.

---

## `.dockerignore`

```
.git
.github
*.md
*_test.go
tmp/          # air's hot reload output directory
```

Go doesn't need `node_modules` in the ignore list obviously — but test files (`*_test.go`) and hot reload artifacts (`tmp/`) are what you'd typically exclude.

---

## Full CI + Docker Build Matrix

If you want to verify the Docker image builds in CI too:

```yaml
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build Docker image
        uses: docker/build-push-action@v6
        with:
          context: .
          push: false     # just build, don't push (for PR checks)
          tags: myapp:test
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

`cache-from/cache-to: type=gha` — GitHub Actions cache for Docker layer caching. Same concept as your existing Bun CI layer cache trick.

---

## Comparison With Your Current Setup

| | Bun + Hono | Go |
|---|---|---|
| Runtime in Docker | ✓ (oven/bun base) | ✗ (binary only) |
| Final image size | ~80–150MB | ~5–20MB |
| CI install step | `bun install` | `go mod download` |
| Type check in CI | `bunx tsc --noEmit` | `go vet ./...` |
| Test command | `bun test` | `go test ./...` |
| Race condition detection | Manual | `go test -race ./...` |
| Hot reload | `bun --watch` | `air` |
| Multi-stage build | Same pattern | Same pattern, smaller result |