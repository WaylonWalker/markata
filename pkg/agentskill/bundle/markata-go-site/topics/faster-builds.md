# Faster Builds

Use this topic when the task is build speed, local iteration speed, or profiling slow plugins.

## First Steps

- use `markata-go build --fast` for faster development loops
- use `markata-go serve --fast` for the normal live-edit loop
- use `markata-go serve --incremental` for a production-style live-edit loop
- use `markata-go reader update` when you only need fresh `/reader/` feed data for the next build
- prefer `[markata-go.blogroll] refresh_on_build = false` when you want to keep blogroll pages but move remote refresh work out of the normal build
- use `markata-go reader update --concurrency <n>` when reader refresh latency is dominated by many remote feeds
- use `-m fast.toml` or `--merge-config fast.toml` when you want a slimmer dev config without editing the main site config
- compare warm builds, not just cold builds
- use `markata-go build --benchmark-summary-json=/tmp/benchmark.json` for compact structured timing
- keep reusable page behavior in shared, deferred cacheable JavaScript assets
  instead of duplicating large inline scripts; preserve inline code needed for
  pre-paint theme and layout decisions
- keep benchmark reports outside the site source tree so Content Index does not fingerprint them as untracked source state
- inspect `.content.template_cache` in benchmark JSON or `.template_cache` in
  `<output_dir>/.markata/diagnostics.json` to explain template work. This optional
  version-1 field contains bounded aggregate counts, not paths, template names,
  hashes, or per-page labels. The first failing gate masks later reasons;
  `nav_preview_reset` is contextual, not a per-page cause. Reconcile
  `classified = restored + render_required` and
  `render_required = render_succeeded + render_failed + serve_deferred`.
  Compare equivalent warm samples with caches intact and account for host noise;
  recurring misses are not automatically bugs or corruption. This partial
  #1339 phase observes template cache decisions, not publish I/O, and preserves
  HTML and invalidation behavior.
- use `markata-go build -v --benchmark-detailed` when you need stage detail
- use `markata-go buildlab run --fixture <path>` to check clean, incremental, and deterministic build behavior
- read the `Slowest requests` footer section before assuming a slow plugin is CPU-bound

## Client-Side Loading

- use the site's **Loading** control to compare Auto, Save Data, and Full Quality
- keep Markdown video defaults at `autoplay = false` and `preload = "none"`; opt in only when autoplay is part of the page's purpose
- authored autoplay becomes a data attribute; Full Quality or already measured fast Auto at page entry can start unrequested videos, while Save Data and the no-JS baseline require interaction
- benchmark Auto separately from manual Save Data: verify real timing evidence without seeding session confidence, and use several runs with medians
- use responsive images, intrinsic dimensions, native lazy loading, and posters before adding client-side loading code
- keep third-party players behind a lightweight facade and reader interaction; do not add unconditional preconnect or preload hints for linked media
- benchmark generated pages on a mobile viewport with a constrained connection, and check that scrolling does not fetch every offscreen video card

Benchmark JSON retains every content entry and feed disposition/reason, along
with summaries, timings, warnings, and blogroll status. Its content serialization
is streamed with reusable per-entry buffers and buffered writes, not a slimmed
report. Large reports still require proportional disk space; temporary content
encoding memory scales with the largest entry, while non-content timing metadata
is encoded together. Use either a file path or `--benchmark-json=-` for the same
complete, two-space-indented JSON payload.

Use `--benchmark-summary-json=<path>` when you only need timings, hotspots, and
aggregate content/cache counts. Run one warm build to prime plugin caches, then
use later warm builds for steady-state comparisons.

## What `--fast` Skips

Fast mode is for local iteration. It keeps the normal content pipeline but skips expensive non-essential work.

From current code, fast mode skips or reduces:

- JS minification
- CSS minification
- CSS purging
- Tailwind rebuild work
- Pagefind indexing
- some fast-mode-aware write work such as redirects generation

## What Still Runs In `--fast`

Fast mode still does the main site build work:

- config loading and validation
- file discovery
- markdown loading and frontmatter parsing
- transforms
- markdown rendering
- template rendering
- feed and collection generation
- normal output writing for the site itself

Unchanged pages restore cached full-page HTML before output is written. This
includes published pages with frontmatter but no Markdown body, so removing
the output directory does not require deleting `.markata/` before a normal or
`--fast` build recreates the page.

Full-page restoration is eager and bounded by build concurrency; it finishes
before any fresh template can inspect cache-hit peers through Core. Do not
clear `Post.HTML`, defer hydration, or evict cache-map ownership as a supposed
memory fix: hooks, Fontpack, Tailwind, and publishing still need complete pages,
and cache/post strings share storage. Compare the separate template phase logs
(`Phase 1a classify`, `Phase 1b batch restore`, `Phase 2 render`) even for
zero-render builds. Restore latency is Phase 1b work, not fresh rendering or a
count of actual disk reads; the aggregate diagnostics schema is unchanged.

For `markata-go build --fast`, file discovery still rescans the content tree on each run. Added,
removed, and moved files should be detected without clearing `.markata/`. Only `serve --fast`
reuses in-memory and on-disk state for incremental rebuilds between change events.

So `--fast` is good for content, template, and most styling iteration, but it is not a full partial build mode. `serve --incremental` keeps the partial-build cache while retaining normal output processing.

## Guidance

- Normal JS/CSS minification persists exact results under
  `<content_dir>/.markata/asset-minify/v1/`, or the nonempty top-level `cache_dir`
  override used as supplied. It is independent of build-cache enablement and
  must stay outside published output, including derived storage paths, private
  parent directories, and symlink relationships. Unsafe locations are rejected
  before directory creation or permission changes. Never
  deploy this private cache: historical snapshots may contain stripped comments,
  though new transforms write only results and exact-input records.
  `transformed/restored/excluded/failed` log counts distinguish real
  engine work from verified exact-input restores. Only exact current-input
  digest plus matching recipe authorizes reuse; the source hash means transform
  input, not original author source. Hits publish atomically with mode `0644`,
  even for identical bytes, never chmod on a shared release inode; empty JS
  retains its no-write exception. Raw static recopies should avoid engine work
  after priming. Do not demand zero transforms for retained output or no-write
  hits: important CSS comments can make transformations non-idempotent.
  Same-size/mtime edits still miss when their contents differ from recorded input.
  CSS comment-option or engine/wrapper changes always transform current
  stage-input bytes. Never select older snapshots using processed-output or
  retry hashes, or use a prior output digest for same-recipe reuse.
  To preserve a
  stripped comment, regenerate from authoritative source; retained output or
  snapshots cannot establish that source's current intent. Missing/corrupt
  results needed for exact-input reuse warn and repair from current bytes.
  Persistence plus invalidation failure leaves
  that target unchanged under the warning-only policy. Keep fast/disabled,
  exclusions, `.min.js`, and `_pagefind` behavior intact. Legacy output sidecars
  remain untouched and untrusted: do one source-regenerating rebuild after
  upgrading, including assets not normally recopied. Do not recommend clearing
  caches alone to recover original source. Prime records before alternating
  equivalent warm timings; avoided engine calls alone do not prove a wall-time gain.
- When automatic feeds are hot, inspect the bounded `auto_feeds` / `collect`
  debug record: `generation_ns`, `filtering_sorting_ns`,
  `selection_recording_ns`, and `pagination_preparation_ns` are exclusive elapsed
  phase sums. Counts include considered occurrences across feeds, not unique
  posts. Complete exclusion recording can dominate filtering; do not drop
  diagnostics or retain membership/flags across builds to improve timings.
  Per-Collect scratch and atomic per-feed ledger batches preserve the complete
  report. Lower temporary allocation is not evidence of a whole-site RSS win.
  New ledger relationships use lazy owned slabs capped at 10 objects or 32
  initial-reason strings, with observation-bounded tails. Raw reason lists over
  32 strings use ordinary standalone append/deduplication, without counting or
  allocating by raw length; duplicate/empty-heavy input retains only unique
  nonempty values. Published storage is never cleared or reused.
  Existing updates allocate no slabs; single recording
  remains individual. There is no cross-build pool or feed registry. Measure
  full-site timings separately rather than equating allocation reduction with
  faster builds.
- The second warm build is the best steady-state comparison point.
- Avoid deleting caches unless you specifically need a cold-build measurement.
- If output is network-bound, inspect plugins that fetch remote content.
- If output is read-heavy, inspect globbing, cache loads, and broad content scans.
- If output is write-heavy, inspect cache saves, feed publishing, Pagefind output, and static output.
- CSS purge caches discovered selectors and purged stylesheets. On Linux, stable verified
  per-file identities can reuse selectors without reading unchanged HTML, even when another
  page changes. `mode=identity` is a complete hit; `mode=identity-partial` hashes only files
  with changed or unverified identities and parses only new content. Added pages are scanned,
  removed pages are dropped, and unsupported or unidentifiable HTML uses exact-content
  fallback. Compare warm builds before changing purge behavior, because disabling purge can
  increase stylesheet transfer.
- If warm builds still spend time in `configure/build_cache`, check whether template or config files actually changed before assuming the cache is stale; the build cache now fingerprints the template tree before it does a full rehash.
- The image-library cache reuses local media content fingerprints when each file's path, size, and modification time are unchanged. On filesystems that expose change time, it also detects same-size replacements that preserve mtime; output directories are excluded from source-media hashing, and the canonical image-library hash uses verified content fingerprints rather than mtime. Avoid clearing `.markata/` when measuring this warm-build path.
- If `/tags` or `/garden` writes are hot, prefer cached per-post semantic hashes so the listing hashes don't need to re-derive the same per-post summaries every build.
- Prefer targeted fixes over broad cache-busting changes.
- HTML cache restoration avoids a second page-sized byte-to-string copy in
  both executors. This reduces warm-build allocation, not cache invalidation
  or work selection. Keep existing caches when measuring it; formats and keys
  are unchanged, and unreadable cache files still cause re-rendering.
- Cached Markdown already includes heading-highlight wrappers; restore it
  unchanged rather than decorating it again. The single-pass correction
  refreshes affected derived pages/feeds using per-post render revisions,
  without a global cache reset or external downloads. A revision is certified
  only after successful fresh full-page caching. Prime warm samples after
  this refresh, and compare the first build after adding content with its
  next warm build.
- Full builds keep the complete content diagnostics artifact; serialization
  streams entries rather than buffering the whole JSON document. Do not disable
  diagnostics to make a benchmark look faster. Owned scratch grows with the
  largest entry, and repeated encoded feed fragments are reused only within
  publication, capped at 1024 entries and 1 MiB including keys; bypassing this
  retention never drops feed observations. Inspect the `diagnostics_artifact`
  cleanup debug record for snapshot, optional source metadata, serialization,
  actual file writes, flush, sync, close and replacement. Add only exclusive
  serialization/flush durations to file-write time; separately labeled inclusive
  wall totals overlap file writes. Directory/temp setup and removal are outside
  these phases; these are elapsed timings, not tracing spans or artifact fields.
  Publication uses a fixed 64 KiB buffer; large writes may bypass it, so measure
  raw calls instead of dividing artifact size by buffer size. Do not translate
  lower temporary serializer allocation into a whole-build RSS claim: live
  templates and rendered HTML can still dominate memory.
  Snapshot copies use one owned feed-reason arena per entry with isolated,
  capacity-clamped slices, including after deduplication. This reduces small
  allocations, not snapshot freshness or ownership; do not introduce snapshot
  reuse or an immutable-cache contract as a follow-on shortcut.
  Sidebar link projections are
  reused only within a build, with current-page highlights kept separate.
- Glossary matching uses longest-first keys with lexical ties, so equal-length
  aliases sharing a link limit do not randomly change cold-build output.
  Nested protected links/code are restored outermost first without internal
  marker leakage. Pre-fix results refresh once through the cache identity.
- The experimental `--dag` executor is serial and does not automatically skip
  more work or add parallelism. Compare it with explicit `--dag=false` on the
  same site/config/cache inputs, alternating at least three warm samples.
  The verbose plan digest is structural, not a content-cache fingerprint.
- On real-site benchmarks, keep encryption settings identical, disable
  publication/index-ingestion side effects, and keep unencrypted measurements
  in isolated local copies. Network-backed embeds and favicons need a shared
  external-cache snapshot before byte-level executor comparisons are meaningful.
- Warm fontpack annotation should not recopy page bodies whose root element
  already selects the resolved pack. If fontpack remains a hotspot, distinguish
  catalog/coverage resolution from HTML annotation before changing typography.
- Cold fontpack resolution shares one visible-rune analysis across packs and
  stops checking a family's coverage when full is required. It selects prebuilt
  font tiers rather than running a subsetter. Built-in cache identity includes
  canonical pack choices, default/picker settings, coverage, and bundled metadata;
  a new cache schema causes one fresh resolution. Missing generated font files
  are regenerated without disabling the picker or weakening custom validation.
- For sites that use `[markata-go.mermaid] mode = "chromium"` or `"cli"`, unchanged
  diagrams should reuse cached SVG output on warm builds; if Mermaid remains a hotspot,
  compare the diagram source and rendering inputs before recommending client mode.

## Slim Config Overrides With `-m fast.toml`

Use merged config overrides when `--fast` alone is not enough.

This is the recommended pattern for a lighter development loop because it lets you keep your main config intact while narrowing the site shape for local work.

Examples:

```bash
markata-go serve --fast -m fast.toml
markata-go build --fast -m fast.toml
markata-go build -m markata-go.local.toml -m fast.toml
```

Typical uses for `fast.toml`:

- narrower content globs
- fewer feeds
- disabling expensive optional features for local work
- lower concurrency or alternate local URLs when needed

Example `fast.toml`:

```toml
[markata-go.glob]
patterns = ["posts/current/**/*.md", "pages/*.md"]

[markata-go]
concurrency = 2
```

A starter version is included at `../examples/fast.toml`.

Use merge configs for scope changes. Use `--fast` for expensive output-step skips. Use both together for the shortest loop.

## Fast Path For Everyday Work

1. `markata-go serve --fast -m fast.toml`
2. if the output still looks wrong, run `markata-go build`
3. if the build is slow, capture `--benchmark-json` or `--benchmark-detailed`
4. compare warm builds before changing caching behavior

## Concrete Profiling Flow

1. run the build that matches the question you are asking
2. read `Resource profile` to see whether the slowdown is mostly CPU, network, disk read, or disk write
3. read `Hotspots` to find the slow plugin hook
4. read `Slowest requests` to find the exact remote waits inside that plugin
5. if the build is still CPU-heavy after network issues are understood, switch to `--cpuprofile`

## Common Culprits

Source-encrypted posts intentionally reparse without decrypted parsed/article
caches. An inferred title should compare canonical hashes to the prior build,
not Load's intermediate untitled values; the handoff is transient and hash-only.
Do not enable plaintext caching to remove template misses. Conservative Load
feed/tag/garden signals still remain.

Distinguish new entries from existing pages with missing derived semantic
metadata after a cold navigation reset. The latter use the original comparison
while repairing canonical hashes; missing metadata alone must not add unrelated
slug misses. Preserve cold -> edit scope assertions instead of adding warm
priming to hide the issue.

Password or key-name rotation, public hint changes, source-path/accessibility
changes, wrapper/browser-crypto revisions, and encrypted-cache loss can
legitimately regenerate wrappers. Those pages and their dependent closure must
be rendered **and published**, including empty-slug roots; diagnostics use
`affected_path`. Valid hits should stabilize byte-for-byte on the next warm
build. Random source-envelope ciphertext alone is not a canonical wrapper
change. Do not treat missing-cache/key errors as hits or clear unrelated
producer invalidations. Deletion drift (#1465) is a separate issue.

Wrapper regeneration clears the old full-page reference before cache writes.
A failed full-page write must leave that reference unavailable across metadata
save/reload; an unchanged wrapper hit then re-renders with
`full_html_unavailable`. Verify recovery with persisted cache and missing output,
not just fresh `ArticleHTML`; never accept republished old-key ciphertext.

For key-rotation acceptance, establish at least two consecutive warm page
restorations with unchanged config/template/nav/static-asset identities before
rotating. Decrypt actual disk-published HTML with the current key, including
reverse rotation and subsequent steady warm builds. Global cold-to-warm resets
must not serve as accidental publication invalidation.

- remote metadata or embed fetching
- blogroll and reader feed refreshes during normal builds when cache-only mode is not configured
- pagefind, minification, and purge steps
- feed-heavy sites with many aggregate pages
- custom templates or plugins doing repeated expensive work
