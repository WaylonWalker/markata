# Static output reproducibility comparison

Status: first measurement interface for #2097, agent-first epic #2087.

## Contract

`scripts/ci/repro_compare.py <left> <right>` compares two **existing**
generated site directories, without performing any builds. It returns a
machine-readable, deterministic JSON report on stdout:

- `schema_version: 1`
- `identical` as a boolean
- `left_file_count`, `right_file_count`
- sorted `only_left` and `only_right` relative file paths
- sorted `changed` records including each path and the previous/current
  byte size and SHA-256 digest
- `excluded_patterns` recording all explicitly requested glob exclusions

The command uses relative POSIX-normalized paths, streams file content, sorts
results for stable output, and does not compare file modification timestamps.
There are **no default exclusions**. Files inside output trees must be regular
files; symlinks are rejected rather than followed. Reads do not modify output.

Exit codes: `0` for an identical set of relative paths and bytes; `1` for
different artifact sets or bytes; `2` for invalid input directories, symlinks,
or I/O errors. Errors are printed to stderr; JSON is emitted for successful
comparisons (including unequal outputs).

## Limits and guarantees

Identical artifact bytes do **not** by themselves prove deterministic builds:
the tool does not pin source content, plugins, toolchain versions, environment,
external data, or build flags. Report these inputs separately before claiming
full reproducibility. Development fast mode is not equivalent to production
build processing. Never silently normalize away timestamps, random IDs, or
external data fields in release outputs to make results appear equal.

## Test coverage

Unit tests cover equal nested trees, edited/additional/missing files, stable
hashes and JSON, explicit exclusions, missing roots and unsafe symlinks.
CI runs `python3 -m unittest discover -s scripts/ci -p 'test_*.py'` for PRs.
Future work under #2097 must add pinned-input isolation, repeated production
builds, external/plugin policies and performance baselines.
