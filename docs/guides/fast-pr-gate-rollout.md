---
title: "Fast PR gate rollout"
description: "Staged rollout of faster Markata pull-request checks and release-safe validation."
date: 2026-10-31
published: true
tags:
  - documentation
  - ci
---

# Fast PR gate rollout

Tracking issue: [#1747](https://github.com/WaylonWalker/markata/issues/1747).

## Current stage: shadow validation

The new [Fast PR validation](../../.github/workflows/pr-fast.yml) workflow provides
an always-present `PR Gate` job on pull requests, pushes to main, merge groups,
and manual dispatch. It runs alongside existing `CI`; **do not replace any
required check or weaken release gates until the new check has been observed**.

PR classification is intentionally conservative:

- Only changes confined to `docs/**`, `docs-site/**`, `README.md`, or
  `CONTRIBUTING.md` qualify for the docs-only Go-test exemption. The existing
  documentation-site workflow still validates relevant documentation changes.
- All other changes run gofmt, a CLI build and smoke check, and `go test -short ./...`
  with the explicitly observational Build Lab characterization excluded.
- Python-related source/packaging changes also run `python/tests` unit tests.
- Main pushes, merge-group checks, and manual runs receive the broad path.
- Failure to classify changed files is an error; unknown paths are **not**
  considered documentation-only. Patch whitespace is always verified.
- Use the stable job name `PR Gate` (not the workflow name) for rulesets.

The full Linux race/coverage shards, macOS/Windows matrix, and specialized
browser/packaging checks continue to run under existing workflows.
The Docker PR job no longer repeats a separate `go test ./...`; Docker image
build and smoke tests remain in place, and publishing still tests source before
advancing images.

## Baseline and progression

Before switching branch protection, record at least 20 completed, non-cancelled
PR runs with each workflow's wall time (created to completion), queue delay,
job-level duration, pass/fail and critical path. Report p50/p90, not only the
fastest run. Report cancelled runs separately because frequent commits can bias
the successful sample. Include a mix of docs, Go, frontend and packaging PRs.

1. **Shadow:** Compare `PR Gate` against the existing `CI` check. Investigate
   any disagreement (especially tests skipped by `-short`).
2. **Gate migration:** After evidence is accepted, make only `PR Gate` and
   relevant focused checks required. Keep the broad matrix on main and nightly.
   Coordinate the change with the existing [CI specification](../../spec/spec/CI.md).
3. **Release hardening:** Validate immutable SHAs and artifacts before any
   stable publishing or promotion. Preserve the existing release/test/sign/SBOM
   checks throughout the transition.

### GitHub configuration (maintainer action)

The workflow itself cannot configure repository branch protection. When ready:
protect `main` with `PR Gate` as a required check, require a current branch or
merge-queue validation, and forbid bypasses. Do **not** drop old required checks
until both are green on real PRs and the policy change has been reviewed.
For merge queue, `merge_group` is already subscribed.

### Rollback

Restore the existing required check set and disable/remove `pr-fast.yml` in a
small revert PR. Docker source tests on actual publication are unchanged.
Never roll back by bypassing a failed release test or forcing a publication.
