---
title: "Verify reproducible Markata builds"
description: "Compare generated static sites by SHA-256 and understand what deterministic output proves."
date: 2026-11-22
published: true
tags: [documentation, building, reproducibility]
---

# Reproducible build verification

Markata generates static site artifacts that you can deploy anywhere. To **prove
that two production builds produced byte-identical files**, compare their
outputs with the repository's offline, standard-library-only comparator:

```sh
python3 scripts/ci/repro_compare.py /tmp/build-one /tmp/build-two
```

The command prints JSON with a stable schema and sorted paths, including any
added/removed files and both SHA-256 digests for changed files. Exit codes:

- `0`: identical file names and byte contents
- `1`: changed, added, or removed output
- `2`: invalid or empty output directory (including a missing build), or unsafe/unreadable output entry

By default, an empty output directory fails verification rather than silently reporting reproducibility. Use `--allow-empty` only if you intentionally want to compare empty directory structures.\n\nNo dynamic files are silently excluded. When needed, pass an explicit
`--exclude 'some/relative/path/*'` and preserve the exclusion list in CI
results. The comparator refuses symlinks rather than reading outside an output
tree. File modification times are deliberately not compared.

Run the offline unit tests with:

```sh
python3 -m unittest discover -s scripts/ci -p 'test_repro_compare.py'
```

## What this proves (and what it does not)

**An identical result only proves the compared output bytes are identical.**
It does not prove that input content, environment, toolchain, plugin versions,
remote fetched assets, or build flags were the same. Record all of those
alongside the report before claiming a deterministic build pipeline.

For meaningful verification, use two isolated, clean workspaces at the same
commit, pin Markata and plugin versions, fix remote data or enable an explicit
offline policy, then perform a complete production build in each workspace.
`markata build --fast` skips some production processing and is not a substitute.
Use the existing build benchmarking options for runtime measurements; this
comparator reports byte-level correctness, not speed.

This is the first measurement slice under
[issue #2097](https://github.com/WaylonWalker/markata/issues/2097).
Next work: input/dependency manifests, repeated clean builds in CI, controlled
plugin/network policies, artifact-level nondeterminism triage, and baseline
cold/warm build budgets.
