---
title: "Maintainer CI Specification"
description: "Validation and publishing policy for the markata-go product repository."
date: 2026-10-01
published: true
tags:
  - specification
  - ci
---

# Maintainer CI

Tracking issue: #1487. This policy applies to the markata-go product repository,
not to site-author build or deployment workflows.

## Required validation

- CI MUST run for pull requests and pushes targeting `main` and `master`, and
  on the existing daily schedule and manual dispatch.
- Path selection MUST happen at job level, not by skipping the PR workflow;
  unrelated changes must still produce the ordinary lint and test checks.
- All three Linux race/coverage shards, their `test-linux` aggregate check, and
  macOS/Windows tests MUST remain. One existing Linux shard MUST smoke-test its
  already-built CLI with `version` and `--help`; no extra build job is needed.
- PR CLI tests MUST retain focused Build Lab regressions but skip the expensive
  observational `TestBuildLab_LinkedAndFixtureMutationsCharacterizeProduct`.
  Main pushes and daily/manual CI MUST continue to include that observation.
- The separate Build Lab characterization workflow MUST be manual-only: daily
  CI already covers the same observation. It MUST retain `-race` for the harness
  without generating unused coverage data. Logged product failures remain
  observational, not evidence that focused regressions passed.

## Specialized checks

Font completeness/determinism, rendering-contract generation, Helm validation,
and GoReleaser configuration validation MUST run when relevant inputs change on
PRs or pushes, and unconditionally on scheduled/manual CI.

Change detection MUST use the PR base diff (including fork PRs) and the previous
push commit for pushes, not compare a pushed default branch with itself.
Detection failure MUST prevent a successful validation run. Read access to
contents and PR files is sufficient; no write permission is needed.

Filters MUST cover:

- fonts: the manifest generator, source catalog/lock, bundled font assets,
  manifests and licenses;
- rendering: generator source, contract specification/fixtures, generated Go
  targets, imported rendering/palette Go packages and embedded palettes, and
  `go.mod`/`go.sum`;
- Helm: the entire chart, dependencies and chained chart test scripts;
- GoReleaser: `.goreleaser.yml` and Go module metadata.

Any workflow change MUST select all specialized checks. Keep filter outputs
explicit and preserve the existing check names and commands.

## Containers and benchmarks

- Docker MUST publish minimal and builder images on `main` pushes and `v*`
  tags, retaining Trivy and tagged amd64/arm64 builds.
- Feature branch pushes MUST NOT publish automatically. Manual dispatch MAY
  publish from the selected workflow ref; branch/SHA tags allow explicit
  previews. PR validation MUST remain build-only for both images.
- Non-PR publishing jobs MUST NOT request unreachable PR metadata tags.
  Existing action versions MUST remain unchanged.
- The initial performance sample MUST select exactly `BenchmarkBuild_EndToEnd`
  and `BenchmarkBuild_Incremental`; concurrency keeps its separate sample set.
  The comparison job MUST generate only the comparison fixture, retaining the
  current checkout/module metadata and downloaded current benchmark artifact.
  Comparison fixture/benchmark failures MUST NOT be silently treated as
  successful measurements.
- Release signing, SBOM generation, release tests, and unrelated workflows
  MUST remain unchanged.

See [Maintainer CI](../../docs/guides/maintainer-ci.md) for operational guidance.

## Agent preflight and failure triage

- Run a fail-closed local merge preflight before agent PR readiness, including Go formatting, short tests, new lint and selected generated docs/contract/editor checks.
- Fast PR validation MUST share executable Go formatting and short-test logic with the local preflight.
- Read-only diagnostics on completed failed workflows MUST NOT execute PR code, alter conclusions, or retry jobs.
- Required-check migration remains separate and MUST be supported by shadow fast-gate evidence.

## Pinned lint installation

Use Go 1.27 and golangci-lint v2.15.0. New findings are gated using `--new-from-rev` on PRs/pushes. A scheduled nonblocking full-tree lint audit tracks legacy debt (issue #1902), with configured linters retained.
