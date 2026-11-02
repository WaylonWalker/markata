---
title: "Maintainer CI"
description: "Run focused product-repository checks, Docker previews, and performance comparisons."
date: 2026-10-01
published: true
tags:
  - documentation
  - ci
  - maintenance
---


This guide describes CI for the **markata-go product repository** (issue #1487).
For site-author deployments, use the [[ci-cd-guide|CI/CD guide]] instead.
The policy is specified in [CI.md](../../spec/spec/CI.md).

## Routine validation

CI runs on pushes and pull requests targeting `main`/`master`, daily at 08:23 UTC,
and through **Actions → CI → Run workflow**. Lint, all three Linux race/coverage
shards, their aggregate `Test (ubuntu-latest, 1.27)` check, and macOS/Windows tests
remain active. The CLI/search shard smoke-tests its existing binary with
`version` and `--help`; there is no separate duplicate build.

Specialized checks use job-level changed-file filters on PRs and pushes:

| Checks | Inputs |
|--------|--------|
| Font completeness and determinism | `scripts/generate_font_manifests.py`, `internal/fontcatalog/**` (catalog, lock, manifests, licenses, font assets), `vendor/google-fonts/**` |
| Rendering contract | `scripts/rendering-contract/**`, `spec/rendering-contract/**`, `pkg/renderingcontract/**`, `pkg/palettes/**`, `palettes/**`, `go.mod`, `go.sum` |
| Helm | `helm-chart/**`, including chart dependencies and chained test scripts |
| GoReleaser config | `.goreleaser.yml`, `go.mod`, `go.sum` |

Changes under `.github/workflows/**` select every specialized check. Scheduled
and manual CI runs execute all specialized checks regardless of paths.
Docs-only PRs still receive ordinary checks rather than skipping the entire
workflow and potentially leaving required checks pending.

The detection job uses `dorny/paths-filter@v3`: PRs use the changed-files API,
including fork PRs, while pushes compare against the previous push commit.
Only this job gets `pull-requests: read`, alongside `contents: read`.
If detection fails, do not interpret skipped specialized jobs as validation.
When adding a generator dependency, update its filter in `ci.yml`.

## Spec conformance manifest

When adding a normative requirement, assign it a stable `DOMAIN-NNN` ID in
`spec/conformance/requirements.yaml`. Record its specification file, scope,
conformance status, and the implementation and test files that support it. Mark
incomplete or intentionally different behavior as `partial` and explain it in
`notes`.

Run `go test ./pkg/specconformance` to check IDs, allowed values, referenced
files, and the evidence required for `implemented` status. The normal CI test
suite runs the same check. See [the conformance manifest guide](../../spec/conformance/README.md)
for the field definitions and rollout plan.

## Build Lab observation

The expensive `TestBuildLab_LinkedAndFixtureMutationsCharacterizeProduct` logs
product failures as observations. PR CLI validation skips that test while
retaining focused regression tests. Main push and daily/manual CI include it.

For a focused rerun, select **Actions → Build Lab characterization → Run
workflow** and the desired ref. This separate workflow is manual-only because
daily CI already runs the observation; a second daily run would duplicate work.
It retains harness race detection, but no longer writes unused coverage data.
Its success does not mean that every logged product observation passed.

## Docker preview publishing

Minimal and builder images publish automatically only on `main` pushes and
`v*` tags. Relevant PR changes still build and validate both images without
publishing. Tagged builds retain amd64/arm64 images and Trivy scanning.

To publish a feature preview explicitly, select **Actions → Docker → Run
workflow** and choose the feature branch in **Use workflow from**. Both images
publish to GHCR with branch and `sha-` tags; feature previews do not receive
`latest`. Use the SHA tag when a stable preview reference is needed. Manual
dispatch on `main` follows the existing default-branch tagging policy, and
dispatch on a version tag follows the existing semver/multiarch policy.

## Performance comparisons

The initial benchmark sample selects only `BenchmarkBuild_EndToEnd` and
`BenchmarkBuild_Incremental`. Concurrency keeps its separate three-sample run.
Profiling and stage-specific sampling remain unchanged.

For a comparison, manually run **Performance Benchmarks** with `compare_branch`.
The current checkout remains available for Go version/module setup; its
benchmark measurements come from the downloaded current-results artifact.
Only the comparison checkout generates a fixture and executes comparison
benchmarks. Fixture or benchmark failures fail the job instead of being hidden
as a successful comparison.

## Scope

Release tests, signing, SBOM generation, action versions, and repository settings
are unchanged. Required-check or remote workflow registry changes are separate
maintainer operations, not part of these repository workflows.

The bundled `markata-go-site` skill was reviewed (entrypoint, build/deployment,
and faster-build guidance). No update is needed: these changes govern product
maintainer CI, not site-author commands, build behavior, deployment workflows,
or benchmarking guidance.

## Golangci-lint v2 on Go 1.27

Go 1.27 uses `golangci-lint` **v2.15.0**, rather than the older v1.64.8
tool which cannot decode the Go 1.27 compiler export format. The migrated
`.golangci.yml` retains all existing configured linters.

The new v2 analyzers also surfaced hundreds of historical findings in
untouched files. To keep meaningful merge protection while progressively
resolving that debt (tracked in
[issue #1902](https://github.com/WaylonWalker/markata/issues/1902)),
PR and push CI now rejects **new findings relative to the PR base or previous
main commit**. Scheduled CI runs the same new-issue gate plus a
**full-tree lint audit** that reports failures without blocking the daily
build. The full audit is deliberately visible, not suppressed.

For local work, run `just lint` after fetching `origin/main` to enforce
new findings, or `just lint-full` to inspect all existing diagnostics.
`just lint-fast` uses the supported v2 `--fast-only` mode. Once issue
#1902 is resolved, restoring blocking full-tree lint is the exit criterion.
