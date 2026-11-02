---
title: "Performance Benchmarking"
description: "How to run, interpret, and optimize markata-go build performance"
date: 2024-01-15
published: true
tags:
  - performance
  - benchmarking
  - profiling
  - optimization
---


markata-go includes a comprehensive benchmarking suite for measuring and optimizing build performance. This guide covers how to run benchmarks locally, interpret results, and use profiling tools to identify bottlenecks.

### Explain template work before optimizing it

Inspect `jq '.content.template_cache' benchmark.json`, or
`jq '.template_cache' public/.markata/diagnostics.json` after a successful build.
The optional, additive version-1 field counts classification, full-page
restoration, render success/failure, and incremental serve deferral. Its bounded
histogram records the first failing cache gate only; earlier gates mask later
ones. `nav_preview_reset` is context, not proof of why individual pages missed.
No paths, template names, hashes, or per-page telemetry are recorded.

Compare equivalent warm builds without clearing their caches, and alternate
samples to account for shared-host noise. A persistent miss count is not
automatically a bug or evidence of corruption. Reconcile `classified =
restored + render_required` and `render_required = render_succeeded +
render_failed + serve_deferred` before attributing time to cache selection.
See [Content Diagnostics](content-diagnostics.md#template-cache-decisions) for
the complete precedence and reconciliation rules. This phase measures template
cache decisions only, not publish I/O, and does not change rendered HTML.

Use `--benchmark-summary-json=/tmp/benchmark.json` for timing comparisons. The
summary retains hotspots and aggregate cache statistics without the much larger
per-source feed matrix. Keep benchmark artifacts outside the site source tree;
Content Index fingerprints untracked source-tree files at both build boundaries.
Run a clean build, one cache-priming warm build, and then repeated steady-state
warm builds before comparing executors or code changes.

Warm template builds still restore every usable full page into public
`Post.HTML` before rendering any misses, so custom templates can inspect
cache-hit peers and a deleted output tree can be repaired. Restoration uses
up to four workers, never exceeding the configured build concurrency, without
changing rendering concurrency or introducing another configuration flag.
The template phase logs separate `Phase 1a classify`, `Phase 1b batch restore`,
and `Phase 2 render`; restore time belongs to Phase 1b even when no pages need
fresh rendering. Compare both restore latency and total template time.
This is a latency optimization, not removal of the live full-page memory floor:
hooks, Fontpack, Tailwind, and publication still receive complete pages.

## Quick Start

Run the end-to-end build benchmark:

```bash
just perf
```

This runs the benchmark 5 times and outputs results to `bench.txt`.

## Persistent CSS and JavaScript Minification

Normal builds reuse minified results for exact current inputs and matching
options, including static assets copied into output again as raw bytes.
No extra setting is required:

```bash
markata-go build
markata-go build
```

The per-plugin log separates work into `transformed` (engine calls), `restored`
(verified exact-input result reuse), `excluded`, and `failed`. Restored assets
receive verified result bytes without invoking the minifier, but still publish
by atomic replacement with mode `0644`, just like fresh transformations. This
also applies when result bytes already match, without changing hardlinked
releases. Empty JavaScript retains its existing no-write exception.
Contents, not file size or mtime, determine reuse. There is no general
processed-output shortcut: CSS important comments matching preservation
patterns can make repeated minification non-idempotent, so retained output may
need another engine call. Static-asset hashed aliases still
receive the final canonical
minified bytes during Cleanup.

Private result blobs and per-asset exact-input records live under
`<content_dir>/.markata/asset-minify/v1/`, independently of the build-cache enabled
setting. A nonempty `[markata-go] cache_dir` override is used as supplied (a
relative path is relative to the build's working directory), like the build
cache. Sites and output roots are isolated even when they share cache storage.
Keep this directory outside published output; unsafe locations, including
symlink relationships, produce a warning and disable reuse. Do not deploy the
private cache: historical snapshots from earlier implementations may contain
comments removed from public assets. New transforms do not write source snapshots.
The derived storage paths and private parent directories must not overlap
published output either; unsafe configuration is rejected before directory
creation or permission changes, then minification proceeds uncached.

Changing CSS `preserve_comments`, the minifier/parser versions, or the maintained
transform revision invalidates affected recipes. Changed recipes always minify
the current stage-input bytes, even when they equal an earlier minified result.
Only exact input plus the same recipe can restore a verified result; matching a
historical output or retry digest never authorizes reuse. The recorded source
hash means exact transform input, not original author source. Historical
snapshots are never read or substituted for current input.
To preserve a previously stripped comment, regenerate the asset from
authoritative source before building; changing options alone cannot recover it
from retained minified output. Missing or corrupt cached results needed for
exact-input reuse are warned about and repaired from current bytes.
Minification errors remain warnings rather than failing the whole build. Cache
write errors also warn; if obsolete provenance cannot safely be invalidated,
that target is left unchanged rather than risking restoration of an older input.

`--fast`, disabled plugins, exclusions, `.min.js`, and `_pagefind` retain their
existing selection behavior. They do not change transform recipes. Legacy
`.markata-{js,css}_minify-cache` sidecars are ignored and left untouched, not
promoted into trusted provenance. After upgrading, perform one rebuild that
regenerates raw source assets, especially assets not normally recopied. Clearing
only the new cache cannot recover original source from an already-minified file.

For measurements, keep caches intact, prime the new records, then compare
equivalent alternating warm builds. Repeated raw copies should avoid engine
work after priming, but do not expect zero transformations for retained
non-idempotent output, or no target writes on hits. Restored counts measure
avoided engine work, not by themselves an end-to-end speedup.

## Running Benchmarks Locally

### Prerequisites

Install `benchstat` for analyzing benchmark results:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

### Available Commands

| Command | Description |
|---------|-------------|
| `just perf` | Run end-to-end benchmarks (5 iterations) |
| `just perf-profile` | Generate CPU and memory profiles |
| `just perf-stages` | Benchmark individual lifecycle stages |
| `just perf-concurrency` | Test performance at different concurrency levels |
| `just perf-compare old.txt new.txt` | Compare two benchmark runs |
| `just perf-generate` | Regenerate the benchmark fixture |

### Running Specific Benchmarks

```bash
# All benchmarks
go test -bench=. -run='^$' -benchmem ./benchmarks/...

# Only end-to-end
go test -bench=BenchmarkBuild_EndToEnd -run='^$' -benchmem ./benchmarks/...

# Only stage-specific
go test -bench='BenchmarkStage' -run='^$' -benchmem ./benchmarks/...

# With more iterations for stability
go test -bench=BenchmarkBuild -run='^$' -benchmem -count=10 ./benchmarks/...
```

## Understanding Benchmark Output

### Build Summary Hotspots

The default `markata-go build` summary now includes two fast feedback signals:

- **Resource profile** - estimated wall-time spent on CPU work, network wait, disk read wait, disk write wait, and idle time
- **Hotspots** - the slowest lifecycle plugin hooks from that build
- **Slowest requests** - the longest outbound HTTP requests with plugin attribution

Example:

```text
Build completed successfully!
  Resource profile (estimated wall time):
    CPU             18.2s (20.1%)
    Network wait    42.7s (47.1%)
    Disk read       10.1s (11.1%)
    Disk write      14.3s (15.8%)
    Idle             5.4s ( 5.9%)
  Hotspots:
    collect/blogroll 31.77s
    cleanup/pagefind 26.32s
    write/publish_feeds 8.54s
  Slowest requests:
    collect/blogroll GET https://example.com/feed.xml 12.40s (HTTP 200)
    transform/mentions GET https://slow.example.com/ 6.80s (context deadline exceeded)
```

Use this summary to decide what tool to reach for next:

- mostly `CPU` -> capture a CPU profile with `just perf-profile`
- mostly `Network wait` -> inspect `Slowest requests` first, then the owning plugins that fetch remote content or external metadata
- mostly `Disk read` -> inspect globbing, cache loads, index reads, and wide content scans
- mostly `Disk write` -> inspect publishing, cache saves, index generation, and emitted static output
- if `configure/build_cache` is hot on warm builds, check whether the template or config tree actually changed; unchanged template trees now reuse a cheap fingerprint before falling back to a full content hash
- if `/tags` or `/garden` hashing is hot, prefer cached per-post semantic hashes over re-deriving the same per-post summaries inside the listing hashers
- mostly `Idle` -> look for subprocess waits, scheduler gaps, or work that is happening outside the Go process

### Concrete Flow For Real Site Builds

Use this flow when you want to answer "what feature is slowing this build down?"

1. Run a build that matches the feature set you care about.
   - for the everyday dev loop, use `markata-go build --fast`
   - for production-style live iteration, use `markata-go serve --incremental`
   - for real feature attribution, run the full build without `--fast`
2. Read the footer in this order:
   - `Resource profile` tells you whether the build is mostly CPU, network, disk read, or disk write bound
   - `Hotspots` tells you which plugin hooks are slow overall
   - `Slowest requests` tells you which exact network calls dominated wall time
3. If `Slowest requests` points to one plugin repeatedly, fix that plugin first.
   - add or verify cache reuse
   - reduce duplicate fetches
   - batch requests or lower request count
   - make timeouts and concurrency explicit
4. Export structured data when you need a diffable artifact:

```bash
markata-go build --benchmark-json benchmark.json
```

5. If the build is still mostly CPU after network fixes, capture `--cpuprofile` and inspect the hottest functions with `go tool pprof`.

### Automatic-Feed Collection

Automatic-feed collection emits one bounded `auto_feeds` / `collect` debug
record. `generation_ns`, `filtering_sorting_ns`, `selection_recording_ns`, and
`pagination_preparation_ns` are exclusive elapsed phase sums in nanoseconds.
`feeds`, `matched`, `selected`, and `observations` are numeric counts (post
counts include occurrences across feeds). There are no per-feed names, filters,
or paths in this record, and these timings are not diagnostics v1 fields.
Use the selection phase when profiling complete diagnostics recording: every
considered occurrence is recorded, including exclusions. Configured and automatic
feeds reuse transient source flags, observation rows, membership maps, and a
reason arena only within one synchronous Collect invocation, and write one
atomic ledger batch per feed. This does not cache selection across builds, change
filters or publication, reduce artifact completeness, or establish an RSS win.

For new batch relationships, ledger-owned object and initial-reason slabs hold
at most 10 dispositions or 32 string slots each, with smaller observation-bounded
tails. Raw reason lists over 32 strings use ordinary standalone append/deduplication
without a counting pass or allocation by raw length; duplicate/empty-heavy input
retains storage only for unique nonempty values.
Allocation is lazy; existing relationships allocate no slabs, and single-record
producers keep their individual allocation path. Published object addresses stay
stable, and copied, ordered-deduplicated reason segments are capacity-clamped so
later appends cannot touch neighbors. This reduces per-relation heap allocations
without changing sparse maps or introducing a registry or cross-build pool.
Allocation results alone do not demonstrate faster full builds.

### Diagnostics Publication

Full-build diagnostics publication has a separate bounded `diagnostics_artifact`
cleanup debug record: snapshot, optional source metadata, serialization/buffering,
actual file-write time/bytes/calls, flush, sync, close and replacement. Use
`serialization_buffering_exclusive` and `flush_exclusive` with `file_write`
to avoid counting file I/O twice; the corresponding `*_inclusive` wall totals
overlap it. Directory/temp setup and removal are not included. No publication
timings are inserted into the diagnostics v1 document.

The diagnostics artifact writer reuses owned per-entry scratch and already
encoded/indented feed fragments within a single call (at most 1024 retained
fragments and 1 MiB including keys). Values outside these caps bypass retention,
not output. This reduces repeated sanitation allocations and JSON indentation
when posts share feed dispositions; it does not remove observations, persist a
cache, or alter output bytes. Snapshot and artifact size still scale with the
complete ledger. Compare equivalent warm builds and fixed-metadata serializer
benchmarks rather than disabling diagnostics to improve a timing.
Publication uses a fixed 64 KiB buffer. Large entry writes can bypass it, so use
the measured raw-write counts rather than estimating calls from artifact size.
Reduced serializer allocation does not establish a whole-build RSS reduction;
live templates and rendered HTML may still dominate memory.
Snapshot copying also replaces per-feed reason allocations with one owned arena
per entry. Feed slices remain isolated and capacity-clamped after deduplication,
and feed sorting uses a typed comparator. This lowers allocation count and
snapshot work without reusing snapshots or changing ownership, counts or output.

### JSON Benchmarks

Use machine-readable output when you want to compare builds over time or ingest
results into another tool:

```bash
markata-go build --benchmark-json benchmark.json
markata-go build --benchmark-json - > benchmark.json
```

The JSON output includes:

- whole-build resource totals
- per-stage timings and per-stage estimated resources
- plugin timing entries used for hotspot ranking
- request timing entries used for the slowest-request list
- build counts and warnings
- the complete content summary and per-source entries, including every feed
  disposition and selection reason

Reports stream content entries with reusable per-entry buffers and buffered
writes, reducing serialization memory on large sites without omitting details.
The two-space-indented JSON format and trailing newline are unchanged. The
report itself can still be large: its size grows with content and feed
observations, while temporary content-encoding buffers grow with the largest
single entry rather than the complete report. Non-content timing metadata is
still encoded together. File and stdout modes retain the same complete payload;
write, flush, and file-close failures are reported as errors.

### Per-Stage Detail

Keep the default footer small for everyday use, and opt into stage detail when
debugging:

```bash
markata-go build -v --benchmark-detailed
```

This adds a per-stage estimated wall-time breakdown so you can see whether a
slow build is CPU-heavy in `render`, read-heavy in `glob`, write-heavy in `write`, or mostly idle in a
subprocess-oriented cleanup stage.

### Raw Output

```
BenchmarkBuild_EndToEnd-8    	       5	 234567890 ns/op	123456789 B/op	 1234567 allocs/op
```

| Field | Meaning |
|-------|---------|
| `BenchmarkBuild_EndToEnd-8` | Test name with GOMAXPROCS |
| `5` | Number of iterations |
| `234567890 ns/op` | Nanoseconds per operation |
| `123456789 B/op` | Bytes allocated per operation |
| `1234567 allocs/op` | Number of allocations per operation |

### Using benchstat

`benchstat` provides statistical analysis of benchmark results:

```bash
# Single run analysis
benchstat bench.txt

# Compare two runs
benchstat old.txt new.txt
```

Example output:

```
name                 time/op
Build_EndToEnd-8     235ms ± 2%

name                 alloc/op
Build_EndToEnd-8     124MB ± 0%

name                 allocs/op
Build_EndToEnd-8     1.23M ± 0%
```

The `±` value shows the variation between runs. Lower is better for reproducibility.

### Comparing Runs

When comparing two benchmark files:

```
name              old time/op    new time/op    delta
Build_EndToEnd-8    250ms ± 3%     235ms ± 2%   -6.00%  (p=0.008 n=5+5)
```

| Column | Meaning |
|--------|---------|
| `old time/op` | Time from first file |
| `new time/op` | Time from second file |
| `delta` | Percentage change (negative = faster) |
| `p=0.008` | Statistical significance (p < 0.05 is significant) |
| `n=5+5` | Number of samples in each file |

## Profiling

### Generating Profiles

```bash
just perf-profile
```

This creates:
- `cpu.prof` - CPU profile
- `mem.prof` - Memory allocation profile

### Analyzing CPU Profiles

#### Interactive CLI

```bash
go tool pprof cpu.prof
```

Common commands:
- `top` - Show top functions by CPU time
- `top -cum` - Show by cumulative time
- `list FunctionName` - Show annotated source
- `web` - Open in browser (requires graphviz)

#### Web Interface

```bash
go tool pprof -http=:8080 cpu.prof
```

This opens an interactive web UI with:
- Flame graphs
- Call graphs
- Source code annotation
- Top functions

### Analyzing Memory Profiles

```bash
go tool pprof mem.prof
```

Useful options:
- `go tool pprof -alloc_space mem.prof` - Total allocations
- `go tool pprof -alloc_objects mem.prof` - Number of allocations
- `go tool pprof -inuse_space mem.prof` - Live memory

### Profile Types

| Profile | What it Measures | When to Use |
|---------|-----------------|-------------|
| CPU | Time spent in functions | Slow builds |
| Memory | Allocations | High memory usage |
| Block | Blocking on sync primitives | Deadlocks/contention |
| Mutex | Mutex contention | Lock performance |

## Benchmark Fixture

The benchmark suite uses a deterministic fixture at `benchmarks/site/`:

```
benchmarks/
├── site/
│   ├── markata-go.toml      # Benchmark config
│   └── posts/
│       ├── blog/2024/01/    # 60 blog posts
│       └── docs/guides/     # 40 documentation guides
└── benchmark_test.go        # Benchmark tests
```

### Fixture Characteristics

- **100 posts total** - Representative of a medium-sized site
- **Code-heavy content** - Syntax highlighting stress test
- **Nested paths** - Tests path handling
- **Multiple languages** - Go, Python, JavaScript, Rust, SQL
- **Deterministic** - Same content every generation for reproducible results

### Regenerating the Fixture

```bash
just perf-generate
```

Or directly:

```bash
go run benchmarks/generate_posts.go
```

## CI Performance Tracking

Performance benchmarks run automatically:
- **Weekly** - Sunday at 2 AM UTC
- **Manual** - Via workflow dispatch

### Viewing Results

1. Go to Actions tab in GitHub
2. Select "Performance Benchmarks" workflow
3. View the job summary for benchstat output
4. Download artifacts for profiles

### Comparing Branches

Trigger a manual workflow with a comparison branch:

1. Go to Actions > Performance Benchmarks
2. Click "Run workflow"
3. Enter the branch name to compare against
4. View the comparison in the job summary

## Optimization Tips

### Common Bottlenecks

1. **Markdown rendering** - goldmark processing
2. **Template execution** - Pongo2 templates
3. **File I/O** - Reading/writing files
4. **Syntax highlighting** - Chroma processing
5. **Memory allocations** - String operations

### Improving Performance

#### Concurrency

Adjust the concurrency level in `markata-go.toml`:

```toml
[markata-go]
concurrency = 8  # 0 = auto (NumCPU)
```

#### Profile-Guided Optimization

1. Run profiling: `just perf-profile`
2. Analyze: `go tool pprof -http=:8080 cpu.prof`
3. Identify hot paths
4. Optimize targeted functions
5. Re-benchmark to verify improvement

#### Memory Optimization

If memory is the bottleneck:

```bash
go tool pprof -alloc_space mem.prof
```

Look for:
- Large allocations (`alloc_space`)
- Many small allocations (`alloc_objects`)
- Functions with high cumulative allocations

### Writing Efficient Plugins

1. **Reuse allocations** - Use `sync.Pool` for buffers
2. **Minimize copies** - Use pointers where appropriate
3. **Batch operations** - Group file writes
4. **Cache results** - Use the lifecycle cache

## Troubleshooting

### Benchmarks Skip

```
--- SKIP: BenchmarkBuild_EndToEnd
    benchmark_test.go:35: Benchmark fixture not found
```

Fix: Run `just perf-generate` to create the fixture.

### High Variance

If `±` values are high (>10%), try:
- More iterations: `-count=10`
- Close other applications
- Use a consistent environment

### Profile is Empty

Ensure the benchmark runs long enough:

```bash
go test -bench=BenchmarkBuild_EndToEnd -run='^$' \
  -benchtime=30s \
  -cpuprofile=cpu.prof \
  ./benchmarks/...
```

### Source-encrypted warm builds

Source-encrypted Markdown deliberately bypasses parsed-post and plaintext
article caches. Reparsing alone should not force fresh page templates: the
canonical title/feed/tag/garden hashes are compared against the previous build
through a build-local, hash-only handoff. No decrypted snapshot is retained by
that handoff, and no cache schema migration is required.

For already processed posts, incomplete semantic hashes (for example after a
cold navigation reset) mean the previous canonical baseline is unavailable.
The build retains the original comparison for that pass and repairs the hashes,
rather than marking every reparsed page changed. New entries still compare
against an explicitly empty/partial baseline. Cold-to-edit cache scope does not
require warm priming.

Encrypted wrappers are reused only when article HTML, key name, resolved
password, hint, source path, and wrapper/browser-crypto revision all match.
Real wrapper regeneration invalidates the page and dependent closure before
templates and publication; subsequent equivalent warm builds stabilize.
It also clears the old full-page reference before cache writes. After a failed
fresh full-page cache write, a persisted reload reports `full_html_unavailable`
and re-renders with the current wrapper rather than restoring old ciphertext.
Randomized source ciphertext for the same plaintext/key contract is not itself
a wrapper miss. Older wrapper identities miss once. Conservative Load dirty
signals for feeds, tags, and garden output remain unchanged.

## See Also

- [Configuration Guide](/docs/guides/configuration/) - Concurrency settings
- [Plugin Development](/docs/guides/plugin-development/) - Writing efficient plugins
- [Go Profiling](https://go.dev/blog/pprof) - Official pprof documentation

## Nested operation spans

Save build measurements with `markata-go build --benchmark-json=benchmark.json`,
then inspect completed spans with `jq '.benchmark.Spans' benchmark.json`.
Durations and start offsets use nanoseconds; a zero duration is valid on clocks
with coarse resolution. IDs and parent IDs describe nesting, and stage/plugin
fields identify the active build work when the span started.

Plugin authors can instrument an operation using `buildstats.StartSpan(ctx,
"template.execute")` and finish the returned handle with `End()` or
`EndError(err)` before the build profile stops. Pass the returned context to
child operations to preserve nesting. With no active profile, these calls are
safe no-ops. This release provides the recording API; builds without explicit
span instrumentation may report an empty span list.

Use fixed operation names and safe attributes. Credential-bearing attribute
keys are dropped, and error text is discarded, but values under other keys are
not automatically redacted. Avoid query strings, sensitive paths, and secrets.
See the [lifecycle specification](../../spec/spec/LIFECYCLE.md) for the API
contract and limits.

## OTLP configuration foundation

The Go API `buildstats.ResolveOTLPConfig(os.Getenv)` resolves optional trace
settings. For example, a generic endpoint `http://collector:4318` resolves to
`http://collector:4318/v1/traces` with the default `http/protobuf` protocol.
The resolver alone does not send traces, and setting these variables does not
yet enable export from `markata-go build`.

| Variable | Meaning |
| --- | --- |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` | Complete trace URL; takes precedence over the generic endpoint |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Generic URL; HTTP protocols append `/v1/traces`, while gRPC preserves the URL |
| `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` | Trace protocol; takes precedence over the generic protocol |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | Generic protocol; defaults to `http/protobuf` |

Supported protocols are `grpc`, `http/protobuf`, and `http/json`. Endpoints must
use HTTP or HTTPS and include a host. With no endpoint, the returned config is
disabled. This foundation follows [OpenTelemetry endpoint conventions](https://opentelemetry.io/docs/specs/otel/protocol/exporter/);
exporter integration, headers, sampling, and collector deployment remain
follow-up work.
