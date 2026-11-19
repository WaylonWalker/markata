# Content Diagnostics and Disposition Specification

This specification defines the canonical ledger for explaining what happens to
content during a Markata Go build.

## Goals

The ledger MUST make the path from discovery to output observable without
requiring raw build logs. It MUST:

- count discovered files and content candidates;
- record load, frontmatter, post, eligibility, render, feed, and output state;
- assign stable reason codes to content that is not eligible or is not emitted;
- retain one deterministic final disposition for every discovered content
  candidate; and
- be safe to consume from human and machine-readable reporting surfaces; and
- publish a safe, versioned diagnostics artifact after a successful full build.

The ledger is the in-memory build result. The normal full build also persists a
sanitized production artifact at `.markata/diagnostics.json` below the resolved
output directory.

## Ownership and Data Flow

`lifecycle.Manager` owns one `diagnostics.ContentLedger` for each build. The
ledger is reset with the manager and is populated at these boundaries:

```text
glob       -> discovered files and content candidates
load       -> bytes loaded, frontmatter state, diagnostics, posts
collect    -> feed selection and filter reasons
render     -> rendered posts and render errors
write      -> emitted posts and output errors
cleanup    -> successful-build diagnostics artifact
```

Plugins MUST update this ledger rather than create another content inventory.
Consumers MUST use the immutable, sorted snapshot returned by the ledger.

Synthetic posts used for feed and page lookup have no source path and MUST NOT
be counted as content files or source posts.

## Candidate and Count Definitions

The snapshot contains one entry for every file returned by glob. Entries whose
extension is not a supported content extension are marked `not_candidate` and
are excluded from content-stage counts. Markdown and supported text content
extensions are candidates.

For candidate entries, the summary fields mean:

| Field | Definition |
| --- | --- |
| `discovered` | Every file returned by glob, including non-candidates. |
| `candidates` | Discovered files with a supported content extension. |
| `loaded` | Candidates whose source bytes were read successfully, including cache-backed posts whose source was inspected. |
| `frontmatter_valid` | Candidates with a structurally valid frontmatter block whose YAML parsed successfully. Empty frontmatter counts as valid. No frontmatter does not count. |
| `posts` | Candidates converted into a source-backed `Post`. |
| `eligible` | Posts with `published: true` that are not drafts or skipped. This is public content eligibility; private feed opt-in is recorded separately. |
| `rendered` | Posts for which Markdown HTML was produced or restored. |
| `emitted` | Posts for which at least one post or feed output was produced or confirmed by the write stage. |
| `excluded` | Candidates that are not eligible or have no emitted output. An unpublished shadow page can therefore be both emitted and excluded from public eligibility. |
| `warnings` | Warning-severity content diagnostics. |
| `errors` | Error-severity content diagnostics. Recoverable authoring errors do not hide valid sibling files; operational source errors fail the build. |

Counts are derived from ledger entries. They MUST NOT be assembled by adding
plugin-local counters.

## Bounded publisher I/O diagnostics

Successful full-build artifacts MAY include `publish_io` before `entries`
in schema version 1. Its optional `publish_html` and `publish_feeds` objects
contain only fixed aggregate counts and durations. Absence is valid for older
artifacts and builds without these publishers.

Each publisher MUST count logical outputs considered, skipped by the existing
cache gates, and successfully processed by its publish path. Logical units are
non-skipped, non-draft posts for HTML and configured feeds for feed publication;
a unit can emit multiple files. Actual successful `os.WriteFile` calls MUST
be counted separately with total payload bytes and summed syscall latency.
Feed byte-equality skips MAY be reported as `files_unchanged`; unsuccessful
write calls MUST increment `write_errors` but not successful bytes/files.
The bounded recorder MUST report exact max latency and an approximate p95 from
a fixed histogram, without collecting a per-file sample list.

`write_duration_ns` is the sum of successful write calls and MUST NOT be
presented as stage wall time; concurrent calls overlap. The derived
`write_throughput_bytes_per_sec` uses bytes divided by summed write durations.
Directory creation, source hashing, file comparisons, removed outputs, and
other plugins' publication work are not covered by these syscall counters.
Both publishers MUST use the same schema under both executors. Counters MUST
remain bounded, concurrent-safe, and build-local, with no paths, filenames,
payloads, secrets, or per-file traces. Reset and Discover MUST clear aggregate
state; snapshots and artifacts MUST own copies. Existing cache and output
semantics MUST NOT change merely to produce telemetry.

## Per-File Entry

### Bounded template-cache observations

Canonical title hashing MUST compare freshly parsed posts against prior-build
semantic hashes, not Load's intermediate pre-auto-title values. This prevents
unchanged source-encrypted posts (which deliberately bypass parsed caches)
from generating false InlineTitles slug changes. Load's conservative dirty
signals remain intact. Encryption-owned wrapper regeneration marks the source
path and dependent closure; template diagnostics report those misses using the
existing `affected_path` reason. The eleven reasons and artifact schema remain
unchanged; no plaintext, passwords, or private debug paths are added.

An existing processed post with incomplete prior semantic metadata MUST use
the unavailable-baseline fallback, not new-post comparison. A cold navigation
reset may leave valid page/input metadata but no semantic hashes. A subsequent
linked-target edit MUST retain unrelated cache hits without warm priming or an
extra `slug_changed` miss merely because those hashes were absent.

If encrypted wrapper regeneration is followed by a failed full-page cache
write, the invalidated full-page reference MUST remain unavailable after
metadata save/reload. A subsequent unchanged wrapper hit MUST use the existing
`full_html_unavailable` reason and render the current ciphertext, without adding
encryption affected-path marks merely for being private.

Snapshots and version-1 artifacts MAY contain `template_cache` before the final
`entries` field. Older version-1 documents without it remain valid. Benchmark
JSON exposes the same value at `.content.template_cache`; the artifact exposes
it at `.template_cache`. This section observes template cache decisions only; publish-I/O diagnostics
are separately aggregated, without changing cache invalidation or HTML.

The optional object contains scalar integer counts `classified`, `skipped`,
`cacheable`, `restored`, `render_required`, `render_succeeded`, `render_failed`,
and `serve_deferred`, a boolean `nav_preview_reset`, and a fixed `miss_reasons`
object. Its eleven integer fields report only the first failing gate:

1. `affected_path`: the existing outer affected-path gate (also on full builds);
2. `cache_unavailable`: no build cache;
3. `input_hash_missing`: the current post has no input hash;
4. `entry_missing`: metadata has no entry;
5. `input_hash_mismatch`: metadata input differs (including a partial cached
   entry with an empty input hash);
6. `template_mismatch`: metadata template differs;
7. `dependency_changed`: a post dependency is in the changed-slug set;
8. `slug_changed`: the post's own slug is in that set;
9. `feed_membership_changed`: a nonempty current membership hash differs;
10. `local_preview_changed`: the existing outer local-preview gate; and
11. `full_html_unavailable`: a phase-1a hit whose phase-1b full-page getter
    returns empty. This is not a corruption claim.

Skipped posts are neither classified nor misses. Empty bodies are classified.
Reasons MUST preserve this short-circuit precedence without reevaluating masked
gates. `nav_preview_reset` records the existing shared-navigation reset result,
not proof that any particular miss was caused by navigation.

For the posts observed by one template invocation, counts MUST reconcile:

```text
post count = skipped + classified
classified - cacheable = sum(first ten miss_reasons)
cacheable = restored + full_html_unavailable
render_required = sum(all eleven miss_reasons)
classified = restored + render_required
render_required = render_succeeded + render_failed + serve_deferred
```

Incremental canonical filtering records a deferred render without a second
miss; noncanonical full-page hits are still restored. Successful rendering
remains successful if best-effort cache storage fails. Heading revision and
encryption certification rules remain unchanged.

The manager-owned ledger MUST copy stats on setting and snapshotting under its
lock. Reset and Discover clear them. Each template invocation clears previous
stats before setup checks and publishes one replacement aggregate after worker
join, including render failures. Setup failure leaves the field absent.
Advancing one build through multiple `RunTo` calls MUST NOT clear observations.
There is no active-profiler or global-manager dependency.

Stats MUST remain bounded: no per-page telemetry records, source paths, template
names, hashes, secret payload labels, disposition reasons, or telemetry I/O.
Persistent misses are observations, not automatically bugs.

Each entry records these booleans:

- `candidate`
- `loaded`
- `frontmatter_present`
- `frontmatter_valid`
- `post_created`
- `eligible`
- `rendered`
- `emitted`
- `excluded`

It also records:

- `path`, relative to the content root when discovered by the standard glob;
- `disposition`, one of `emitted`, `shadow`, `excluded`, or `not_candidate`;
- sorted, unique `reasons`;
- frontmatter and lifecycle diagnostics; and
- feed-level dispositions with the feed name, inclusion state, and reasons.

Disposition rules are deterministic:

1. non-candidates are `not_candidate`;
2. emitted eligible posts are `emitted`;
3. emitted ineligible posts are `shadow`; and
4. all other candidates are `excluded`.

Entries, feed dispositions, reasons, and diagnostics MUST be sorted before a
snapshot is exposed. Concurrent plugin execution MUST NOT change snapshot
ordering or duplicate a diagnostic.

### Snapshot memory and concurrency

Snapshots MUST own their flat entry, diagnostic, feed-disposition, and reason
slices independently of the ledger and of other snapshots. Later ledger updates
or caller mutation of a snapshot MUST NOT change another snapshot or the ledger.
Copying observations MUST occur under the ledger's read lock; sorting and
deriving dispositions and summary counts SHOULD occur after releasing that lock.
Snapshot construction SHOULD copy feed values directly into flat slices, without
cloning internal feed maps or diagnostic deduplication indexes.
Feed reasons SHOULD be copied into one owned contiguous arena per entry rather
than allocated separately for every feed. Each feed's reason slice MUST have a
disjoint capacity-clamped segment, including after sorting and deduplication:
appending or mutating one feed MUST NOT alter siblings, the ledger, or another
snapshot, even when live input reason slices overlap. Empty-feed and empty-reason
nil semantics, ordering and summary counts MUST remain unchanged. Arenas MUST be
newly owned on each snapshot; no snapshot reuse or immutable-cache contract is
introduced. This reduces snapshot allocation, not the live HTML/template RSS.

## Ordered Feed Observation Batches

`ContentFeedObservation` contains a raw `Path`, `Included`, and `Reasons`.
`RecordFeedBatch(feed, observations)` MUST have the same result as ordered
`RecordFeed` calls, including normalization, unknown noncandidate discovery,
empty feed names, last-observation inclusion, and union/deduplication of reasons.
Entry reasons MUST include reasons from EACH excluded observation, even if a
later observation includes the source. Nil ledgers and empty batches are no-ops.
The ledger MUST own all retained slices, with no callbacks under its lock.
Each batch MUST be atomic relative to snapshots, discovery, resets, and writes.
Feed storage MUST remain sparse; recording MUST NOT allocate a sources-by-feeds
matrix. Reset and Discover MUST discard any recording-only feed identity state.

New batch relationships MAY use batch-local owned slabs instead of individual
heap objects. Object slabs MUST have stable backing arrays (never appended after
publishing pointers), at most 10 dispositions, and a tail bounded by the remaining
observation count. Initial reason slabs MUST have at most 32 string slots and
tails bounded by remaining observations and the current unique reason count.
Raw reason input longer than 32 strings MUST bypass slab reservation and use the
ordinary standalone append/deduplication path, without counting or allocating by
raw input length. Retained reason storage grows with unique nonempty values, not
duplicate or empty input. Reasons MUST be copied and deduplicated in original
order. Slab segments MUST be capacity-clamped; later appends MUST NOT alter
neighbors. Published object arrays and reason segments MUST never be reused or
cleared. Only the current batch owns allocator headers; none survive in the ledger
after recording, Reset, or Discover.
Slabs MUST be allocated lazily only for missing candidate relationships, never for
existing updates. No global/cross-build pool, dense matrix, caller storage, or
feed registry is introduced. Single recording MUST retain its ordinary individual
allocation path. Live slabs belong to the recorded ledger values, not reusable
scratch; snapshots MUST still own independent copies after Reset or Discover.

Automatic-feed collection MAY emit bounded debug phase records for generation,
filtering/sorting, selection recording, and pagination/preparation. Records MUST
contain only fixed phase names, numeric counts, and durations, never source paths,
feed names, filters, or content. Phase durations MUST cover actual work without
double-counting nested totals. These records do not add fields to schema v1.

## Reason Codes

Reason codes are stable identifiers. Human-readable messages may change without
changing the code.

### Frontmatter

| Code | Meaning |
| --- | --- |
| `frontmatter.suspicious_delimiter` | A frontmatter-like block starts with a delimiter such as `----` or `-----` instead of `---`. |
| `frontmatter.leading_whitespace` | A frontmatter-like opening delimiter is preceded by whitespace. |
| `frontmatter.malformed_closing_delimiter` | A frontmatter block ends with a delimiter other than `---`. |
| `frontmatter.missing_closing_delimiter` | An exact `---` opening has no closing delimiter. |
| `frontmatter.parse_error` | YAML or frontmatter metadata could not be parsed. |
| `frontmatter.duplicate_key` | YAML contains a duplicate key. |
| `frontmatter.invalid_type` | A common metadata field has an unsupported scalar or collection type. |
| `frontmatter.invalid_date` | A date field cannot be parsed using Markata's supported date formats. |

### Content and output

| Code | Meaning |
| --- | --- |
| `content.published_false` | The post is not eligible for public published feeds. It may still produce a shadow page. |
| `content.draft` | The post is a draft and is not emitted. |
| `content.skip` | The post is explicitly skipped. |
| `content.private` | The post is outside a public feed unless that feed opts into private content. |
| `content.filtered` | A configured feed filter did not select the post. |
| `content.duplicate_slug` | A slug conflict prevents a unique output destination. |
| `content.no_output` | No post or feed output was available for the candidate. |
| `content.load_error` | Source bytes or source-backed post state could not be loaded. |
| `content.render_error` | Markdown rendering failed for the post. |
| `content.write_error` | An output write failed. |
| `feed.offset` | The post matched the feed's privacy and filter selection, but was removed by the configured offset. |
| `feed.limit` | The post matched the feed's privacy and filter selection and offset window, but was removed by the configured limit. |

The implementation MAY add codes in a namespace, but existing codes MUST retain
their meaning.

## Frontmatter Detection

The parser recognizes frontmatter only when the first line, after an optional
UTF-8 BOM, is exactly `---`. The closing delimiter is also a line containing
exactly `---`. Line endings are normalized for parsing.

The following rules apply:

1. A BOM is accepted and does not produce a diagnostic.
2. A first line of `----`, `-----`, or another longer hyphen run produces a
   warning only when a following frontmatter-like block has an exact closing
   delimiter and YAML-shaped fields. Ordinary Markdown horizontal rules do not
   produce a warning.
3. Leading spaces or tabs before an otherwise valid opening delimiter produce a
   warning only when the surrounding block is frontmatter-like.
4. A malformed or missing closing delimiter produces a diagnostic with the
   source line and expected delimiter.
5. Invalid YAML, duplicate keys, invalid dates, and common metadata type errors
   are recorded with source positions when available.
6. Markdown that has no frontmatter remains valid and produces no frontmatter
   warning merely because YAML-shaped text appears later in the body.

Malformed frontmatter, invalid YAML, and invalid frontmatter metadata are
source-local authoring errors. They do not cause a valid sibling file to
disappear. The loader records the affected candidate and continues the content
lifecycle. Direct parser APIs continue to return parsing errors to their
callers.

Operational source errors are different. Stat failures, read failures,
permission errors, unexpected missing sources, and other filesystem or source
I/O errors are recorded as `content.load_error` and propagated from the load
stage. They MUST fail the build instead of producing a successful partial
build.

## Feed and Output Semantics

Feed selection is recorded per feed. A post can be emitted as a direct shadow
page while being excluded from a public feed because `published` is false,
because it is private, because its feed filter does not match, or because it
falls outside the feed's offset or limit window. `content.filtered` is used
only for filter rejection; `feed.offset` and `feed.limit` identify posts that
matched the feed before windowing. Feed-level reasons do not change the
source's single global final disposition.

The absence of an automatically generated tag feed is explained by the ledger
through the metadata and feed state that led to it. A tag that was lost because
frontmatter was not recognized does not create a tag feed; the source entry
still contains the frontmatter diagnostic and the resulting publication
eligibility reason.

## Console Reporting

The normal `build` summary MUST include a `Content:` section with the summary
counts. Content warnings identify the relative source path, line when known,
reason code, and concise message.

The existing global `-v, --verbose` flag is the explain mode for builds. It MUST
show every excluded or shadow entry and its sorted reason codes. It MAY also
show feed-level exclusions. The mode MUST NOT print raw configuration,
environment values, secrets, absolute source paths, or raw build logs.

The existing `--benchmark-json` machine-readable build output includes the
sanitized content snapshot. No separate diagnostics command or public HTML page
is added.

`--benchmark-summary-json` MUST write the same build identity, workload counts,
warnings, timing profile, blogroll status, content summary, and template-cache
aggregate while omitting `content.entries`. It MUST accept a file path or `-`
for stdout. The compact report is intended for timing comparisons where the
per-source, per-feed diagnostic matrix would add disproportionate disk use.
When either JSON flag uses stdout, the command MUST suppress the human build
footer. Supplying both JSON flags MUST fail before the build starts.

### Benchmark JSON serialization

Benchmark JSON MUST retain the complete raw build-result snapshot, including all
entries, feed dispositions and reasons, summaries, timings, warnings, and
blogroll status. Serialization MUST NOT sort, sanitize, or mutate its input.
The wire format MUST remain byte-for-byte compatible with the legacy indented
JSON encoder: executor normalization and validation, field order, two-space
indentation, HTML and string escaping, nil versus empty entry arrays, and the
final newline are preserved.

Content entries MUST be streamed with reusable encoding and indentation buffers
whose retained capacity scales with the largest single entry, not the complete
entry array. Non-content metadata may be encoded together. Actual writes MUST
be buffered to avoid tiny per-entry file writes. Encoding, writing, short writes,
flushing, and file closing failures MUST be returned; a close failure MUST NOT
replace an earlier serialization or flush failure.

## Production Diagnostics Artifact

The normal full build MUST write the diagnostics artifact to:

```text
<output_dir>/.markata/diagnostics.json
```

The artifact is enabled by default and does not require a configuration entry.
It is part of the generated output, so Builder Admin carries it into each
successful release automatically.

### Versioned wire format

The artifact uses this top-level shape:

```json
{
  "$schema": "markata://schemas/content-diagnostics/v1",
  "schema": "markata.content-diagnostics",
  "schema_version": 1,
  "generator": {
    "name": "markata-go",
    "version": "0.0.0",
    "commit": "..."
  },
  "source": {
    "commit": "..."
  },
  "built_at": "2026-09-10T12:00:00Z",
  "summary": {},
  "entries": []
}
```

`schema`, `schema_version`, `generator`, `built_at`, `summary`, and `entries`
are required. `generator.commit` is omitted when the binary does not provide a
reliable build commit. `source` is omitted when the source directory is not a
Git checkout or its `HEAD` cannot be read. `built_at` is an RFC 3339 timestamp
and is recorded in UTC.

`summary` and `entries` are the manager's `ContentLedgerSnapshot` copied
without recomputing counts or making a second content inventory. Entries,
reasons, feed names, and diagnostics retain the deterministic ordering defined
above. Consumers MUST dispatch on `schema` and `schema_version`; future
versions MUST NOT be interpreted as version 1.

### Complete, bounded serialization

Artifact publication MUST preserve every recorded feed disposition, including
nonincluded feeds and their reasons. Performance improvements MUST NOT change
schema, completeness, indentation, escaping, redaction, or empty `entries: []`
representation, and MUST NOT add a trailing newline or a configuration switch.
For identical build metadata (including a fixed `built_at`), streaming output
MUST be byte-for-byte equivalent to `MarshalArtifact`.

The streaming writer MUST leave the caller's snapshot unchanged, even when
reason slices share backing storage. Entries MUST be sorted by normalized path;
normalized-path collisions and duplicate feed names MUST preserve their original
relative order. Normalization and diagnostic sanitation MUST use the same rules
as the existing artifact API.

Publication SHOULD serialize one sanitized entry at a time, bounding temporary
JSON buffers by the largest entry rather than the entire artifact. An ordering
index proportional to the number of entries is permitted. The immutable
snapshot and the complete output file still scale with all recorded observations;
this contract does not truncate diagnostics or bound total artifact size.
`MarshalArtifact` MUST remain available with its existing interface and behavior.
Streaming MUST propagate serialization and writer failures, including short
writes reported without an error.

Per-call sanitized scratch MAY be reused, but MUST own all mutable slices,
including overlapping caller feed/reason storage, and grow only with the largest
entry. `NewArtifact` MUST still return independently owned mutable entries.
Serialized feed fragments MAY be reused within one write call only. Such a cache
MUST retain at most 1024 entries and 1 MiB of combined key and fragment bytes;
overlarge values and values arriving after either cap MUST bypass retention.
Keys MUST be unambiguous for arbitrary strings (including nulls, separators,
Unicode and escapes), with exact equality rather than hash-only identity.
JSON escaping MUST use the standard encoder. Fragment assembly MUST guard the
required struct layout and fall back to whole-entry encoding if it changes.
Neither scratch nor fragment caches may persist across calls or builds.

File publication SHOULD use a fixed 64 KiB output buffer, independent of
artifact size, to reduce underlying write calls. The buffer MUST still be
checked on flush before sync, close and replacement. It does not change the
complete artifact or benchmark content (including template-cache observations).
Serializer allocation improvements do not imply lower whole-build RSS: live
templates and rendered HTML remain outside this temporary-storage bound.

### Publication phase observations

Publication MAY emit one bounded structured debug record with numeric counts
and durations only, without source paths, feed names, or diagnostic messages.
Snapshot and optional source-revision lookup durations MUST include the actual
operations. Serialization/buffering and flush durations MUST exclude time spent
inside underlying file `Write` calls; file-write duration, returned bytes and
call count MUST observe those actual calls, including partial/error returns.
Sync, close and replacement durations MUST cover their actual operations.
The serialization and flush wall totals MAY also be reported, but MUST be
explicitly labeled inclusive (overlapping file-write time), never added to the
exclusive phase totals. These are elapsed measurements, not tracing spans.
No timing fields are added to the version-1 artifact or its own snapshot.

### Privacy and failure behavior

The artifact MUST contain only the sanitized ledger fields. It MUST NOT contain
raw configuration, environment values, secrets, absolute paths, raw logs,
frontmatter, Markdown bodies, rendered HTML, decrypted private content, or
encryption key names. Persisted diagnostic messages MUST be concise canonical
messages for approved diagnostic codes; arbitrary producer-provided message
text MUST NOT be copied into the artifact.

The artifact is written after all normal write and cleanup hooks complete. The
writer MUST use a temporary file in the artifact directory and replace the
destination only after the complete JSON document has been written. A failed
serialization, write, or buffered flush MUST NOT replace an existing artifact.
Buffered output MUST be flushed before syncing and closing the temporary file,
and replacement MUST occur only after those operations succeed. A failed
build MUST leave an existing artifact unchanged; a first failed build MUST NOT
leave a partial artifact. A run that completes all stages with non-critical
lifecycle warnings may publish the artifact; an artifact write failure is a
build failure.

Dry runs, fast or incremental builds, partial lifecycle calls that stop before
cleanup, and development server requests that do not complete a full build MUST
NOT publish a new artifact.
