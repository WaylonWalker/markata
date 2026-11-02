---
title: "Build Lab"
description: "Compare clean and incremental Markata-Go builds in isolated workspaces"
date: 2026-08-31
published: true
tags:
  - documentation
  - build-lab
  - reproducibility
---


Build Lab compares a baseline binary with a candidate binary. It runs each
build in an isolated temporary workspace and records output manifests and
correctness checks. It runs ordinary `markata-go build` commands. It does not
use a task graph or scheduler flags.

Build Lab copies the fixture, redirects cache and temporary paths, disables
fixture `.env` loading, and pins child processes to UTC, `C.UTF-8`, and
`SOURCE_DATE_EPOCH=0`. Custom plugins, arbitrary binaries, and fixture content
remain trusted inputs. Process groups or Windows Job Objects provide best-effort
descendant cleanup; they are not an operating-system sandbox.

## Run a smoke comparison

Build the binary first:

```bash
go build -o /tmp/markata-go ./cmd/markata-go
markata-go buildlab run \
  --fixture /path/to/site \
  --baseline /tmp/markata-go \
  --candidate /tmp/markata-go \
  --result buildlab-result.json
```

Each checkpoint compares:

- clean baseline output with clean candidate output;
- candidate incremental output with clean candidate output; and
- two clean candidate builds when determinism checking is enabled.

The default scenario performs a cache/output cleanup, a priming build, and a
no-op incremental build. Use `--fast` for a quicker smoke check. Fast mode is
not a substitute for a full production-output comparison.

## Use a scenario file

Scenario files contain an identifier, version, seed, and ordered operations:

```json
{
  "id": "title-change",
  "version": "1",
  "seed": 42,
  "operations": [
    {"type": "build"},
    {
      "type": "replace-exact",
      "path": "content/foo.md",
      "old": "title: Old",
      "new": "title: New"
    },
    {"type": "build"}
  ]
}
```

Run it with:

```bash
markata-go buildlab run \
  --fixture /path/to/site \
  --baseline /tmp/markata-go \
  --candidate /tmp/markata-go \
  --scenario title-change.json \
  --result buildlab-title-change.json
```

Build Lab supports `build`, `write-file`, `replace-exact`, `delete-file`,
`rename-file`, `copy-file`, `set-config`, `touch-file`, `clean-cache`, and
`clear-output`. Paths stay inside the fixture, and `replace-exact` must match
exactly one occurrence. Semantic mutations must be followed by a build.

The foundation covers clean, no-op, edit, add, delete, rename, linked-target,
and configuration mutations. Add fixture-specific scenarios for other product
surfaces.

A linked-target checkpoint compares the dependent page and its feed outputs.
The incremental build must rerender known dependents before feeds and pages are
written. Unrelated posts must remain eligible for cache hits.

## Read failures correctly

Build Lab is a measuring instrument. A result with `failure_class: "product"`
means that a child build or an output comparison found a Markata-Go behavior
problem. It does not mean that the harness should be changed to hide the
problem. A `failure_class: "harness"` result means that the comparison itself
was not trustworthy, for example because of a timeout, output truncation, path
violation, or manifest error.

Product failures include concise diagnostics in the JSON result and on stderr.
Known failures can therefore be recorded while the foundation remains focused
on measurement.

Setup errors also emit a versioned JSON result with `failure_class: "harness"`
when the command reaches the Build Lab runner. Invalid command-line syntax is
reported by the normal CLI parser instead.

## Output classes

Manifest records are deterministic by default. Known minification and fontpack
metadata files are excluded because they are build state, not publishable
output. Use `--volatile` for known changing publishable paths, such as the
default `.well-known/time` file:

```bash
markata-go buildlab run --fixture /path/to/site --volatile .well-known/time
```

Use `--check-determinism=false` only when a fixture cannot support the clean
rebuild check. Keep that exception visible in CI configuration.

## Isolation and environment

Use the global `--merge-config` flag to apply a fixture-local profile to both
child binaries:

```bash
markata-go buildlab run --fixture /path/to/site --merge-config fast.toml
```

The `--env` option accepts only `PATH`,
`MARKATA_GO_ENCRYPTION_ENABLED`, and `MARKATA_GO_OFFLINE`. Build Lab does not
inherit arbitrary environment variables, including proxies, cloud settings,
Git settings, agent settings, or credentials.

```bash
markata-go buildlab run \
  --fixture /path/to/site \
  --env MARKATA_GO_ENCRYPTION_ENABLED=false
```

Use `--timeout` and `--gomaxprocs` to control execution. Record external tool
versions with repeated comma-separated entries:

```bash
markata-go buildlab run \
  --fixture /path/to/site \
  --tool-version pagefind=1.5.2,tailwind=4.1.0
```
