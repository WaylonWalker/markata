---
title: "Real-Site Warm Build Optimization"
description: "Measure and reduce steady-state build costs on a large Markata site"
date: 2026-10-02
published: false
tags:
  - performance
  - build
---


## Findings

The supplied legacy and DAG comparisons each contain a clean build followed by
one warm build. That first warm build primes plugin caches and is not the
steady-state measurement described by the performance guide. A third build on
the same output and cache skips all 421 feeds, reducing `publish_feeds` from
9.03 seconds to 58 milliseconds and total time from 24.05 seconds to 10.36
seconds.

The steady-state costs are full-page template cache restoration, automatic feed
collection and diagnostics, Content Index source snapshots, and asset-cache
checks. The four benchmark JSON files also add nearly 1 GB of untracked data.
Content Index deliberately fingerprints untracked source-tree files at both
build boundaries, so storing benchmark artifacts in the site root adds about
1.47 CPU-seconds to every measured build.

## Options

1. Add a compact benchmark report. This preserves the existing complete report
   while avoiding the per-post feed matrix when timing is the goal.
2. Increase the bounded full-page cache read pool. An experiment with sixteen
   workers produced overlapping baseline and candidate timings after the OS
   page cache was warm, so retain the existing cap of four.
3. Relax Content Index source fingerprinting. This could remove more time in a
   dirty benchmark checkout, but it would weaken the source-stability contract.
4. Change feed diagnostics. This could reduce the 1.4 million observations made
   by automatic feeds, but it would alter the published diagnostics contract.

Use option 1. Keep the existing complete report for content-level diagnosis.
Benchmark files should also be written outside the source tree when measuring
the build itself.

## Implementation and validation

Add a `--benchmark-summary-json` output that retains timing and aggregate
content statistics but omits content entries. Compare its size with the full
report and verify both output modes. Compare at least three consecutive builds
after cache priming.

This adds an optional CLI output and changes no site output or cache keys. Update
the CLI reference, performance guide, and bundled site skill.

## Measured result

On `waylonwalker.com`, the compact report is 20,361 bytes. Each existing full
report is about 248 MB, so the compact form is more than 12,000 times smaller.
The steady-state build completed in 9.73 seconds and retained the hotspot list,
stage/resource profile, workload counts, content summary, and template-cache
statistics. The full-page restoration concurrency experiment showed no stable
benefit and was reverted.
