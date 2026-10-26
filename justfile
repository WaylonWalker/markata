# justfile for markata
# https://just.systems/man/en/

# Default recipe - show available commands
default:
    @just --list

# Project variables
project := "markata"
main := "./cmd/markata"
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
commit := `git rev-parse --short HEAD 2>/dev/null || echo "none"`
date := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
ldflags := "-s -w -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Version=" + version + " -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Commit=" + commit + " -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Date=" + date

# ─────────────────────────────────────────────────────────────────────────────
# Development
# ─────────────────────────────────────────────────────────────────────────────

# Build the binary for local development
build:
    go build -ldflags '{{ldflags}}' -o {{project}} {{main}}

# Explicit canonical build alias
build-markata:
    go build -ldflags '{{ldflags}}' -o markata ./cmd/markata

# Build with race detector (slower but catches race conditions)
build-race:
    go build -race -ldflags '{{ldflags}}' -o {{project}} {{main}}

# Run the application (pass arguments after --)
run *args:
    go run -ldflags '{{ldflags}}' {{main}} {{args}}

# Install to $GOPATH/bin
install:
    go install -ldflags '{{ldflags}}' {{main}}

# Clean build artifacts
clean:
    rm -f {{project}}
    rm -rf dist/
    go clean -cache -testcache

# ─────────────────────────────────────────────────────────────────────────────
# Testing
# ─────────────────────────────────────────────────────────────────────────────

# Run all tests
test:
    go test -v ./...

# Run the dependency-free adaptive loading policy tests with Node.js.
test-adaptive-loading:
    node --test scripts/adaptive-loading.test.cjs

# Test the rendered Helm nginx configuration and HTTP behavior.
test-helm:
    ./helm-chart/tests/nginx-cors.sh

# Verify language-neutral rendering contract projections are current.
rendering-contract-check:
    go run ./scripts/rendering-contract --check

# Run tests with race detector
test-race:
    go test -v -race ./...

# Run tests with coverage
test-coverage:
    go test -v -coverprofile=coverage.out -covermode=atomic ./...
    go tool cover -func=coverage.out

# Generate HTML coverage report
coverage-html: test-coverage
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report: coverage.html"

# Run a specific test (e.g., just test-one TestConfig)
test-one name:
    go test -v -run {{name}} ./...

# Run tests for a specific package (e.g., just test-pkg ./pkg/config)
test-pkg pkg:
    go test -v {{pkg}}/...

# ─────────────────────────────────────────────────────────────────────────────
# Code Quality
# ─────────────────────────────────────────────────────────────────────────────

# Format all code
fmt:
    go fmt ./...
    @echo "Code formatted"

# Run go vet
vet:
    go vet ./...

# Run golangci-lint via `go run` so CI/local checks don't depend on a preinstalled binary
lint:
    go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run --timeout=5m

# Run fast lint with only essential linters (good for development iteration)
lint-fast:
    go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run --timeout=2m --fast

# Run lint on only changed files (compared to main branch)
lint-new:
    go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run --timeout=5m --new-from-rev=origin/main

# Run lint with reduced parallelism (less CPU/memory pressure, good for laptops)
lint-gentle:
    GOLANGCI_LINT_CONCURRENCY=4 go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run --timeout=5m

# Run all quality checks
check: fmt vet lint test

# Tidy dependencies
tidy:
    go mod tidy
    go mod verify

# ─────────────────────────────────────────────────────────────────────────────
# Release
# ─────────────────────────────────────────────────────────────────────────────

# Create a snapshot release (no publish)
snapshot:
    goreleaser release --snapshot --clean --skip=publish

# Check goreleaser config
release-check:
    goreleaser check

# Build for all platforms (dry run)
release-dry:
    goreleaser release --snapshot --clean

# Create and push a new version tag
# Usage: just tag v0.1.0
tag version:
    @echo "Creating tag {{version}}..."
    git tag -a {{version}} -m "Release {{version}}"
    @echo "Tag created. Push with: git push origin {{version}}"

# Push a tag to trigger release
push-tag version:
    git push origin {{version}}

# Full release: create tag and push (triggers GitHub Actions)
release version: (tag version) (push-tag version)
    @echo "Release {{version}} triggered!"

# ─────────────────────────────────────────────────────────────────────────────
# Development Helpers
# ─────────────────────────────────────────────────────────────────────────────

# Show version info that would be embedded
version-info:
    @echo "Version: {{version}}"
    @echo "Commit:  {{commit}}"
    @echo "Date:    {{date}}"

# Watch for changes and rebuild (requires watchexec)
watch:
    watchexec -e go -r -- just build

# Run the built binary with version command
version: build
    ./{{project}} version

# Generate and view docs (if godoc is installed)
docs:
    @echo "Starting godoc server at http://localhost:6060/pkg/github.com/WaylonWalker/markata/"
    godoc -http=:6060

# ─────────────────────────────────────────────────────────────────────────────
# Site Development (dogfooding)
# ─────────────────────────────────────────────────────────────────────────────

# Build the site
site-build: build
    ./{{project}} build

# Serve the site with live reload
site-serve: build
    ./{{project}} serve

# Create a new post
site-new title: build
    ./{{project}} new "{{title}}"

# ─────────────────────────────────────────────────────────────────────────────
# Setup & CI Helpers
# ─────────────────────────────────────────────────────────────────────────────

# Set up the development environment (install all required tools)
setup:
    @echo "Setting up development environment..."
    go mod download
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install golang.org/x/tools/cmd/goimports@latest
    go install github.com/goreleaser/goreleaser/v2@latest
    @echo ""
    @echo "Setup complete! You can now run:"
    @echo "  just build    - Build the binary"
    @echo "  just test     - Run tests"
    @echo "  just ci       - Run full CI checks"
    @echo "  just hooks    - Install pre-commit hooks"

# Set up the full development environment with pre-commit hooks
setup-full:
    ./scripts/setup-dev.sh

# Install pre-commit hooks
hooks:
    pre-commit install
    pre-commit install --hook-type commit-msg

# Run pre-commit on all files
pre-commit:
    pre-commit run --all-files

# Run what CI runs (lint + test + build)
ci: tidy vet lint test build
    @echo "CI checks passed!"

# Install development tools (alias for setup)
tools: setup

# ─────────────────────────────────────────────────────────────────────────────
# Performance Benchmarks
# ─────────────────────────────────────────────────────────────────────────────

# Run end-to-end build benchmarks (5 iterations for stable results)
perf:
    @echo "Running end-to-end build benchmarks..."
    @echo ""
    go test -bench=BenchmarkBuild -run='^$$' -benchmem -count=5 ./benchmarks/... | tee bench.txt
    @echo ""
    @echo "Benchmark results saved to bench.txt"
    @echo "For profiling, run: just perf-profile"

# Run benchmarks with CPU and memory profiling
perf-profile:
    @echo "Running benchmarks with profiling..."
    go test -bench=BenchmarkBuild_EndToEnd -run='^$$' -cpuprofile=cpu.prof -memprofile=mem.prof ./benchmarks/...
    @echo ""
    @echo "Profiles generated:"
    @echo "  CPU profile: cpu.prof"
    @echo "  Memory profile: mem.prof"
    @echo ""
    @echo "View with: go tool pprof cpu.prof"
    @echo "Or web UI: go tool pprof -http=:8080 cpu.prof"

# Run all stage-specific benchmarks
perf-stages:
    @echo "Running stage-specific benchmarks..."
    go test -bench='BenchmarkStage' -run='^$$' -benchmem -count=3 ./benchmarks/...

# Run concurrency benchmarks
perf-concurrency:
    @echo "Running concurrency benchmarks..."
    go test -bench='BenchmarkBuild_Concurrency' -run='^$$' -benchmem -count=3 ./benchmarks/...

# Compare benchmarks (requires benchstat: go install golang.org/x/perf/cmd/benchstat@latest)
perf-compare old new:
    benchstat {{old}} {{new}}

# Generate benchmark fixture (100 posts)
perf-generate:
    @echo "Generating benchmark fixture..."
    go run benchmarks/generate_posts.go
    @echo "Generated benchmark posts in benchmarks/site/posts/"
