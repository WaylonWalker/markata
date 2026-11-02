---
title: "Content Diagnostics"
description: "Inspect the content diagnostics artifact produced by a markata-go build"
date: 2026-09-10
published: true
tags:
  - documentation
  - content
  - diagnostics
---


After a successful full build completes, markata-go writes a
sanitized diagnostics artifact below the configured output directory:

```text
<output_dir>/.markata/diagnostics.json
```

For example, a build with `-o dist` writes `dist/.markata/diagnostics.json`.
The artifact is enabled by default. It is part of the generated output, so a
Builder Admin release from such a build carries it with the rest of the site.

## Inspect the artifact

Use a JSON viewer such as `jq`:

```bash
jq . public/.markata/diagnostics.json
jq '.executor' public/.markata/diagnostics.json
jq '.summary' public/.markata/diagnostics.json
jq '.entries[] | select(.disposition == "excluded")' public/.markata/diagnostics.json
```

Replace `public` with the output directory used by the build. There is no
separate diagnostics command or generated public HTML page.

## Stable fields and reason codes

The top-level `schema`, `schema_version`, `generator`, `built_at`, `executor`,
`summary`, and `entries` fields are stable parts of the versioned format.
Consumers must check both `schema` and `schema_version` before interpreting the
document.

`executor` records the build engine that produced the artifact. Normal builds
report `legacy`; builds explicitly using the feature-flagged DAG executor report
`dag`. The same identity is included in `markata-go build --benchmark-json`
output, which makes local benchmark reports and Builder Admin release artifacts
unambiguous without parsing human-readable logs.

`summary` contains counts for discovered files, candidates, loaded sources,
valid frontmatter, posts, eligible posts, rendered posts, emitted outputs,
excluded candidates, warnings, and errors. Each entry includes a relative
source `path`, lifecycle state booleans, a final `disposition`, sorted `reasons`,
diagnostics, and feed-level dispositions.

Feed dispositions are complete: feeds that did not include a source remain in
the artifact with their exclusion reasons. On large sites this can make the
JSON file substantial. Full builds stream one sanitized entry at a time instead
of constructing full-artifact serializer buffers, reducing temporary memory
without dropping observations or changing the indented JSON format. The
immutable ledger snapshot and final file still contain the complete diagnostics;
there is no additional configuration to enable this behavior.
Snapshot construction copies feed reasons into one owned arena per entry, with
each feed's slice capacity clamped even after deduplication. Appending or mutating
one feed cannot change sibling feeds, the ledger, or another snapshot. Every
snapshot remains a fresh copy; counts, ordering and empty-value semantics are
unchanged. Fewer small allocations do not imply a reduction in whole-build RSS.

Reason codes are stable identifiers. Messages can change without changing a
code. The current codes cover:

- frontmatter: `frontmatter.suspicious_delimiter`,
  `frontmatter.leading_whitespace`, `frontmatter.malformed_closing_delimiter`,
  `frontmatter.missing_closing_delimiter`, `frontmatter.parse_error`,
  `frontmatter.duplicate_key`, `frontmatter.invalid_type`, and
  `frontmatter.invalid_date`;
- content and output: `content.published_false`, `content.draft`,
  `content.skip`, `content.private`, `content.filtered`,
  `content.duplicate_slug`, `content.no_output`, `content.load_error`,
  `content.render_error`, and `content.write_error`;
- feed windowing: `feed.offset` and `feed.limit`.

The artifact uses schema `markata.content-diagnostics` and schema version `1`.
Older version-1 artifacts created before executor reporting may omit the
`executor` field; readers should treat that as unknown rather than infer it from
other fields.

## Template cache decisions

Completed template rendering adds an optional `template_cache` object to the
version-1 artifact, and to `content` in benchmark JSON:

```bash
jq '.template_cache' public/.markata/diagnostics.json
jq '.content.template_cache' benchmark.json
```

Older version-1 artifacts without this object remain valid. It is independent
of the profiler and contains bounded counts, not per-page cache diagnostics.
`classified` counts non-skipped posts (including empty bodies); `skipped` is
separate. `cacheable` counts metadata hits, `restored` counts usable full-page
hits, and `render_required` includes misses later deferred by incremental serve.
`render_succeeded` and `render_failed` count selected render results;
`serve_deferred` counts pages excluded by the canonical incremental selection.
A best-effort cache write failure does not turn a successful render into failure.

Bounded parallel restoration does not change these counts or the artifact
schema. All reads finish before fresh templates run; restored peer pages remain
visible through Core. Unavailable pages retain their deterministic input order
after classification misses, and an existing in-memory full-page hit remains
usable even if its disk file is gone. Template phase logs measure classification,
restoration, and fresh rendering separately, including zero-render builds;
aggregate counts alone are not disk-read measurements.

`miss_reasons` reports only the **first failing gate**, in this order:
`affected_path`, `cache_unavailable`, `input_hash_missing`, `entry_missing`,
`input_hash_mismatch`, `template_mismatch`, `dependency_changed`, `slug_changed`,
`feed_membership_changed`, `local_preview_changed`. A metadata hit can instead
fall back to rendering with `full_html_unavailable` when the full-page getter
returns empty. Empty cached input hashes are partial entries compared normally,
not evidence of corruption. Empty or unavailable full HTML is not a corruption
claim either.

Earlier reasons mask later ones; the histogram does not enumerate every
possible reason. `nav_preview_reset` reports whether shared navigation metadata
reset the page cache in this invocation, not proof that each miss was caused
by navigation. Repeated invocations replace the observations rather than
accumulating them.

To reconcile counts:

```text
classified - cacheable = sum(first ten reasons)
cacheable = restored + full_html_unavailable
render_required = sum(all eleven reasons)
classified = restored + render_required
render_required = render_succeeded + render_failed + serve_deferred
```

Expect cold-build misses and some recurring misses. Persistent misses are not
automatically bugs: compare equivalent warm builds and their aggregate reasons
before investigating invalidation. Timing varies with host activity and I/O.
These counts contain no paths, template names, hashes, secrets, or per-page
labels, and add no telemetry I/O. This is the template-cache portion of #1339;
publish-I/O diagnostics are not implemented by this phase.

## Unchanged private posts and encryption changes

Source-encrypted posts are parsed on every build to avoid caching decrypted
Markdown or article HTML. Their inferred titles are compared using canonical
post-transform hashes from the previous build, not the temporary untitled
values produced during Load. Losing a plain post's parsed cache follows the
same rule. The handoff is transient and hash-only; diagnostics gain no titles,
decrypted bodies, passwords, or per-private-post debug paths.

Changing a password, key name, hint, or encrypted-wrapper format, or losing the
encrypted wrapper cache, legitimately regenerates a wrapper. Encryption marks
the source path and dependent pages for fresh templates and publication,
including an empty-slug homepage. These misses use `affected_path`, not a new
diagnostic reason. An unchanged valid wrapper hit does not add affected paths.
An existing page can have missing semantic metadata immediately after a cold
navigation reset. Those missing hashes are treated as an unavailable baseline,
not evidence that the post is new. A subsequent edit therefore keeps unrelated
pages cacheable while repairing canonical hashes, without requiring an extra
warm-up build.
Regeneration also clears the old full-page cache reference. If fresh page-cache
storage fails, the next persisted reload re-renders using the current wrapper
and reports `full_html_unavailable`; it cannot restore the old-password page.
Load may still conservatively dirty feeds, tag indexes, or garden metadata;
this fix does not clear those signals or address deletion drift (#1465).

## Privacy and publication behavior

Feed collection retains every considered source occurrence and exclusion,
including feeds whose publishing work is skipped. Ordered per-feed recording is
atomic with respect to ledger snapshots: the last observation controls inclusion,
while reasons are unioned and excluded observations also contribute entry reasons.
Duplicate raw paths and normalized aliases are not discarded during selection.
Collection scratch is invocation-local, and later builds refresh flags and
membership. The bounded `auto_feeds` / `collect` debug record separates generation,
filtering/sorting, selection recording, and pagination/preparation durations and
counts without paths, feed names, or new artifact fields.

New batch feed relationships use owned slabs capped at 10 dispositions or 32
initial-reason string slots, with smaller tails bounded by remaining observations.
Raw lists above 32 strings use ordinary standalone append/deduplication, not
reservation by raw length or a counting pass. Retained storage grows with unique
nonempty values. Reasons are copied and deduplicated in occurrence order; slab
segments are disjoint and capacity-clamped, so later updates cannot append into a
neighbor. Published storage is never cleared or reused. Only missing candidate
relationships allocate slabs. Existing updates and ordinary single-record
producers do not use this allocator; there is no registry or cross-build pool.
Snapshots remain independently owned even after ledger reset or rediscovery.

Streaming reuses owned sanitation scratch sized by the largest entry. Within
one publication it also reuses standard-encoded feed fragments, capped at 1024
fragments and 1 MiB of combined keys and JSON bytes. Overlarge or high-cardinality
values are still encoded completely, without retention. Nothing is cached across
builds; the snapshot and final file still grow with all recorded observations.
The schema, ordering, escaping and indented bytes are unchanged.
Publication uses a fixed 64 KiB file buffer, checked on flush before sync, close
and replacement. Large individual writes can bypass buffering, so raw-write
counts must be measured rather than inferred from artifact size divided by the
buffer size. These improvements concern temporary serialization storage, not
whole-build RSS: live templates and rendered HTML still retain their own memory.

The `diagnostics_artifact` cleanup debug record separates `snapshot`,
`source_metadata`, `serialization_buffering_exclusive`, `file_write`,
`flush_exclusive`, `sync`, `close`, and `replace`. `file_write_bytes` and
`file_write_calls` count actual underlying file writes, not buffered calls.
Serialization/buffering and flush **exclusive** durations exclude those file
writes. Their separately labeled **inclusive** wall totals overlap file-write
time: do not add them to the exclusive components. Source metadata includes the
optional Git lookup; missing Git metadata is still harmless. These elapsed
measurements contain no paths or messages and are not tracing spans or fields in
the artifact. Directory/temp setup and temp removal are outside these phases;
`published=0` means replacement did not complete.

The artifact contains sanitized ledger data only. It does not contain raw
configuration, environment values, secrets, absolute paths, raw logs,
frontmatter, Markdown bodies, rendered HTML, decrypted private content, or
encryption key names.

The file is published only after the normal full lifecycle completes. A failed
build leaves an existing artifact unchanged and does not leave a partial new
file. Streaming output is flushed, synced, and closed in a temporary sibling
file before replacing the artifact; serialization or write failures leave the
previous artifact intact and fail the build. `markata-go build --dry-run`,
`markata-go build --fast`, partial lifecycle calls, and fast or incremental
development server requests do not publish a new artifact. A normal `serve`
build may publish one when its full build completes successfully.