---
title: "Fast PR gate rollout"
description: "Staged rollout of faster Markata pull-request checks and release-safe validation."
date: 2026-10-31
published: true
tags:
  - documentation
  - ci
---


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
- All other changes check gofmt on changed Go files, run a CLI build and smoke check, and `go test -short ./...`
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

## Observed baseline: 2026-10-29 through 2026-10-31

Queried the repository's completed Actions runs through October 31, identified
91 completed pull-request runs of the `CI` workflow (47 cancelled, 32
successful, 11 failed, and 1 other conclusion), then selected the latest **25 non-cancelled runs** for latency
statistics. The sample is 19 successes and 6 failures; cancelled runs are
excluded from latency percentiles, **not** silently treated as successes.

Wall time uses workflow `updated_at - created_at` (an approximation which can
include rerun elapsed time); job runtime uses `completed_at - started_at`.
Nearest-rank percentiles are reported in seconds, rounded to the nearest
second.

| Measure (25 sampled PR CI runs) | p50 | p90 |
| --- | ---: | ---: |
| CI workflow wall time | 323s (5m23s) | 445s (7m25s) |
| Windows short tests | 257s | 409s |
| macOS short tests | 118s | 140s |
| Linux core race/coverage | 264s | 270s |
| Linux CLI/search race/coverage | 273s | 287s |
| Linux site integration race/coverage | 280s | 312s |
| Lint | 186s | 195s |

Windows has the largest p90; Linux site/CLI also contribute substantially to
the critical path. **Do not assume Windows is always the bottleneck.** The
sample's longest observed CI wall time was 1307s, associated with a rerun.
One rerun exposed inconsistent job scheduling timestamps; do not subtract
job creation times to claim accurate runner queue time without verifying the
run attempt and dependency waits.

Representative runs:
[recent success](https://github.com/WaylonWalker/markata/actions/runs/41243717555),
[Linux + Windows comparison](https://github.com/WaylonWalker/markata/actions/runs/41220799184),
[slow rerun](https://github.com/WaylonWalker/markata/actions/runs/41133499891).

This baseline measures the **existing** CI workflow, not the newly introduced
`PR Gate`. Repeat the same sample method after its deployment; add separate
runner queue and dependency-wait metrics from attempt-specific timestamps.

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

## Editor parity split — issue #1927

Wheel construction and editor parity are separately named jobs. Parity consumes the wheel artifact from the packaging job. An editor failure can no longer obscure a successfully built wheel; screenshots and Playwright traces accompany parity failures. Explicit release publishing depends on both jobs, while the migration of required merge checks remains a separate, measured change.

### Observed shadow-gate snapshot (November 2, 2026)

A read-only sample of the latest 400 PR-triggered workflow runs contained **45 Fast PR validation** and **45 CI** runs. Among completed runs, the Fast PR gate had 13 successes, 4 failures and 27 cancellations (44 total); CI had 12 successes, 1 failure and 30 cancellations (43 total). Excluding cancelled and nonfinal runs and requiring attempt 1, fast gate latency was **p50 218s, p90 272s** (17 runs), versus full CI **p50 372s, p90 457s** (13 runs). Times are approximated by workflow `updated_at - created_at` and may include queue/dependency waits. They are not comparable matched-pair test durations.

Across 43 same-SHA workflow pairs, there were **no observed Fast PR success / full CI failure disagreements** in this narrow sample. However, fewer than 20 completed noncancelled runs exist for either workflow in the snapshot, and this is not yet sufficient to migrate branch protection. Continue collecting path-varied runs, critical-path job latencies and false-negative evidence before changing required checks. A successful browser parity job and wheel packaging job must also be evaluated independently.

For future snapshots, enumerate PR-triggered Actions runs with `gh api 'repos/WaylonWalker/markata/actions/runs?event=pull_request&per_page=100&page=1'` and paginate. Group by workflow name and SHA, exclude cancelled and incomplete attempts from successful latency percentiles, and record cancellation counts separately. Confirm every candidate required-check change against the actual latest SHA and merge-group checks; do not infer safety from one p50 comparison.
