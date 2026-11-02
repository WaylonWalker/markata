---
title: "Experimental DAG builds"
description: "Compare the serial DAG executor with normal builds without changing your site's configuration."
date: 2026-09-29
published: true
tags:
  - documentation
  - builds
  - performance
---


```bash
markata-go build --dag --benchmark-json=dag.json
markata-go build --dag=false --benchmark-json=legacy.json
```

`--dag` selects an experimental task-graph executor. The normal lifecycle
executor remains the default. Both retain the same plugin order, lifecycle
stages, output formats, and cache behavior. The DAG executor runs one graph
task at a time; existing plugins may still use their own internal worker pools.
It is not a switch for additional parallelism or a guarantee of faster builds.

The DAG executor takes one snapshot of the plugins scheduled for each stage.
It builds each plugin's graph when that plugin is reached, after earlier plugins
have finished. A plugin added during a stage can run in a later stage, but does
not join the stage already in progress. This preserves the normal lifecycle's
plugin ordering while allowing later graphs to use the current post and feed
state.

The flag is also available on `serve` and `builder-admin`. Builder Admin passes
the choice to queued build processes. Use `MARKATA_GO_DAG=true` for an
environment-wide opt-in, and `--dag=false` to override it for one command.
An invalid environment boolean produces a clear error unless an explicit flag
overrides it.

`build --dry-run` keeps its existing partial lifecycle through Collect. It does
not execute the DAG scheduler or publish the full-build diagnostics artifact.
Full-build benchmark JSON and `<output>/.markata/diagnostics.json` identify the
executor that actually ran.

## Compare builds fairly

Run a cold build, one warm-up build, and at least three warm samples per
executor. Alternate the order of the warm samples. Use the same source,
configuration, worker count, output-processing flags, and external cache
snapshot. Report cold and warm measurements separately.

Keep encryption settings identical for both executors. If local measurements
disable encryption because keys are unavailable, do so only in isolated local
copies, state that limitation, and do not deploy those outputs. Disable external
index ingestion or other publication side effects for local benchmarks.

Live network responses can change embeds, favicons, and feed content between
builds. Reuse a populated external cache for scheduler comparisons rather than
interpreting remote variability as a DAG difference. Compare generated output
bytes as well as timings; diagnostics timestamps, executor identity, and
`.well-known/time` intentionally differ.

The verbose `plan=` digest describes the graph's declarations. It is not a
content fingerprint, proof that inputs are unchanged, or permission to skip
tasks. Graph compilation snapshots metadata and uses canonical ordering,
but task functions still own their runtime inputs and mutations.

For development iteration, use [[cli-reference|CLI Reference]] and the existing
`--fast` or `serve --incremental` options rather than assuming `--dag` skips
output work.
