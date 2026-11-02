---
title: "Build DAG plugin ownership inventory"
description: "This document is the conservative ownership audit for the feature-flagged DAG executor. It follows the resource vocabulary in pkg/builddag: site, files,..."
date: 2026-09-30
published: true
tags:
  - documentation
---


This document is the conservative ownership audit for the feature-flagged DAG executor. It follows the resource vocabulary in `pkg/builddag`: `site`, `files`, `post`, `feed`, `cache`, `output`, `template`, `external`, and `plugin`.

The inventory is intentionally stricter than an implementation-level code review. A plugin is **not** parallel-safe merely because its loop body looks independent. Unknown ownership, manager mutation, package/global caches, filesystem walks, external processes, and network access all keep work serialized until an explicit task boundary owns those resources.

## Classification

- **exclusive** — site/global state or unclear ownership; keep serialized.
- **aggregate** — reads or mutates a collection and therefore forms a barrier.
- **per-post candidate** — plausible one-task-per-post boundary after an API refactor.
- **per-feed candidate** — plausible one-task-per-feed boundary after an API refactor.
- **per-output candidate** — plausible one-task-per-owned-output-path boundary after an API refactor.
- **external/network** — external state, subprocesses, or network access requires a separate isolation/rate-limit policy even if item-local work exists.

All current `LegacyTasks` remain `ScopeSite`, `Exclusive=true`, and `ParallelSafe=false`. The classifications below describe possible future decomposition only.

## Configure / discovery / load

| Plugin | Classification | Conservative resources / side effects | Refactor required before fan-out |
| --- | --- | --- | --- |
| `build_cache` | exclusive | write `cache:*`, plugin state, lifecycle persistence | Keep barrier; expose namespace-safe cache operations before any overlap. |
| `theme_calendar` | exclusive | read/write `site:config`, plugin state | Configuration mutation must finish before downstream tasks. |
| `glob` | exclusive | write `files:*`, read `site:config`, filesystem walk | Discovery is a barrier that establishes the file set. |
| `background` | exclusive | read/write `site:config` | Configuration mutation must remain ordered. |
| `cdn_assets` | external/network + per-output candidate | network, cache, write `output:static/*`, config/plugin state | Separate fetch/cache planning from deterministic per-asset writes; add rate limits. |
| `icon_vendor` | external/network + per-output candidate | network/cache, vendor asset outputs, plugin index | Split vendor fetch/index from per-pack output ownership. |
| `tailwind` | external/network/exclusive | subprocess/toolchain, templates/includes, CSS output, config | Isolate subprocess inputs/outputs and eliminate shared include mutation. |
| `load` | aggregate/barrier | read `files:*`, write `post:*`, post index, cache | Load establishes the post collection; per-file parse tasks require deterministic merge/index phase. |
| `python_docs` | external + aggregate | source filesystem/import inspection, creates posts | Split source discovery/parse from deterministic post insertion. |
| `tag_aggregator` | aggregate | reads/writes all posts/tags and tag-derived state | Requires item-local normalization plus deterministic aggregate merge. |

## Transform

| Plugin | Classification | Conservative resources / side effects | Refactor required before fan-out |
| --- | --- | --- | --- |
| `auto_title` | per-post candidate | write one `post:<id>` | Extract a post-local transform that does not iterate/mutate Manager collections. |
| `inline_titles` | per-post candidate | read/write one post title representation | Expose post-local transform. |
| `authors` | per-post candidate | read author config, write one post | Freeze author registry/config then expose post-local resolution. |
| `description` | per-post candidate | read/write one post | Expose post-local transform. |
| `structured_data` | per-post candidate | read/write one post metadata | Expose post-local transform with immutable site metadata input. |
| `reading_time` | per-post candidate | read content, write one post | Expose post-local calculation. |
| `stats` | aggregate | scans posts and writes shared statistics | Keep barrier or split item facts from aggregate reduction. |
| `breadcrumbs` | per-post candidate | read site/post routing state, write one post | Freeze route/index view then expose post-local generation. |
| `embeds` | per-post candidate + external/network | post mutation, cache/network for remote embeds | Split cached/local expansion from network resolution; explicit cache/external claims and rate limits. |
| `wikilinks` | aggregate/barrier | reads post index, writes post content/dependency edges | Freeze lookup index and collect dependency-edge deltas for deterministic merge. |
| `mentions` | per-post candidate | read blogroll/config, write one post | Freeze mention lookup state then expose post-local transform. |
| `hashtag_tags` | per-post candidate | read/write one post tags/content | Expose post-local transform after tag normalization barrier. |
| `webmentions_fetch` | external/network + per-post candidate | cache/network, attaches webmentions to posts | Separate fetch/cache refresh from post-local attachment. |
| `webmentions_leaderboard` | aggregate | scans webmention counts across posts | Keep barrier/reduction. |
| `toc` | per-post candidate | read/write one post | Expose post-local transform. |
| `jinja_md` | per-post candidate | read `template:*`, write one post | Make template environment read-only/thread-safe and expose post-local render. |
| `icons` | per-post candidate | read icon index/plugin state, write one post | Freeze icon index then expose post-local expansion. |

## Render

| Plugin | Classification | Conservative resources / side effects | Refactor required before fan-out |
| --- | --- | --- | --- |
| `render_markdown` | per-post candidate | write one post rendered HTML, read renderer/plugin config | Expose a pure/post-local render operation with immutable renderer state. |
| `heading_anchors` | per-post candidate | read/write one post HTML | Expose post-local transform. |
| `image_zoom` | per-post candidate | read/write one post HTML | Expose post-local transform. |
| `webawesome` | per-post candidate | read/write one post HTML | Expose post-local transform. |
| `contribution_graph` | per-post candidate | read/write one post HTML | Expose post-local transform; isolate any shared renderer cache. |
| `md_video` | per-post candidate | read/write one post HTML | Expose post-local transform. |
| `youtube` | per-post candidate | read/write one post HTML | Expose post-local transform; keep network discovery out of task if introduced later. |
| `chartjs` | per-post candidate | read/write one post HTML | Expose post-local transform. |
| `csv_fence` | per-post candidate | read/write one post HTML | Expose post-local transform. |
| `mermaid` | per-post candidate + external | post HTML plus renderer/browser/cache depending on mode | Separate deterministic cache lookup from external rendering; explicit cache/external claims. |
| `glossary` | aggregate + per-post candidate | shared glossary/index plus post mutation; also writes glossary output | Freeze glossary dictionary before post-local linking; keep output generation separate. |
| `wikilink_hover` | aggregate + per-post candidate | reads post index/metadata, mutates one post | Freeze hover lookup snapshot then fan out post-local application. |
| `external_link_hover` | external/network + per-post candidate | network/cache and post HTML | Separate fetch/cache from deterministic application and rate-limit external tasks. |
| `link_collector` | per-post candidate | reads one rendered post, writes one post link metadata | Expose post-local collection; aggregate consumers wait on barrier. |
| `encryption` | per-post candidate | post content plus key/config state | Freeze immutable key/config inputs; ensure secrets never appear in graph diagnostics. |
| `link_avatars` | external/network + per-post candidate | favicon network/cache and post HTML | Separate fetch/cache refresh from deterministic post-local insertion. |
| `templates` | per-post candidate | read `template:*`, post/site metadata, write rendered page state | Make template environment immutable/thread-safe and expose one-post rendering API. |

## Collect

| Plugin | Classification | Conservative resources / side effects | Refactor required before fan-out |
| --- | --- | --- | --- |
| `slug_conflicts` | aggregate/barrier | reads all post output routes | Keep validation barrier; could reduce per-post route facts first. |
| `series` | aggregate | reads posts, creates/mutates feeds/series metadata | Keep aggregate or split series-key grouping from deterministic feed creation. |
| `subscription_feeds` | aggregate | mutates shared feed collection | Keep feed-collection barrier. |
| `feeds` | aggregate | scans posts and constructs feeds | Keep aggregate until filtering/sorting can produce immutable feed plans. |
| `auto_feeds` | aggregate | scans tags/posts and creates feeds | Keep aggregate; deterministic feed-plan generation could precede per-feed work. |
| `blogroll` | external/network + aggregate | network/cache and shared blogroll/feed state | Separate external fetch plan/cache from deterministic aggregate merge. |
| `prevnext` | per-feed candidate | reads ordered feed posts, writes post navigation metadata | Build immutable feed ordering first; resolve conflicts for posts present in multiple feeds before fan-out. |
| `overwrite_check` | aggregate/barrier | reads all planned output paths | Keep validation barrier after output planning. |
| `static_file_conflicts` | aggregate/barrier | filesystem/static inventory + planned outputs | Keep validation barrier. |

## Write

| Plugin | Classification | Conservative resources / side effects | Refactor required before fan-out |
| --- | --- | --- | --- |
| `static_assets` | per-output candidate | filesystem walk, writes `output:static/*`, asset hashes | Produce immutable copy plan then one task per destination namespace/path; serialize shared hash ledger merge. |
| `fontpack` | per-output candidate | reads config/assets, writes font/CSS namespace | Build immutable manifest then own disjoint output paths. |
| `palette_css` | per-output candidate | config -> known CSS output | Make generation pure and claim exact output path. |
| `aesthetic_css` | per-output candidate | config -> known CSS output | Make generation pure and claim exact output path. |
| `chroma_css` | per-output candidate | config -> known CSS output | Make generation pure and claim exact output path. |
| `css_bundle` | aggregate/per-output | reads multiple CSS outputs and writes bundle | Barrier on CSS producers; then own bundle path. |
| `publish_feeds` | per-feed candidate | reads one feed, writes feed-format paths | Generate immutable feed plan and one task per feed/format output. |
| `well_known` | per-output candidate | site metadata -> `.well-known/*` | Claim exact output namespace and remove shared writer state. |
| `publish_html` | per-post candidate + per-output candidate | reads one rendered post/template state, writes post output paths | **Top candidate:** expose one-post publish function and exact output claims; template state must be immutable. |
| `images` | aggregate | scans posts/assets and writes image inventory/pages | Keep aggregate until image records are planned; individual thumbnails/assets may fan out later. |
| `random_post` | aggregate/per-output | scans eligible posts, writes random endpoint/data | Keep eligibility reduction; output write itself can own fixed path. |
| `redirects` | per-output candidate | reads redirect declarations, writes redirect paths | Build/validate unique redirect plan, then one task per destination path. |
| `error_pages` | per-output candidate | templates/site config -> fixed error-page outputs | Freeze template state and claim exact paths. |
| `tags_listing` | aggregate/per-output | scans tag/feed state, writes tags page | Aggregate barrier, then fixed output task. |
| `feeds_listing` | aggregate/per-output | scans feed state, writes feeds page | Aggregate barrier, then fixed output task. |
| `garden_view` | aggregate | scans post graph, writes graph/page assets | Keep aggregate; later split graph reduction from fixed output writes. |
| `sitemap` | aggregate/per-output | scans published routes, writes sitemap outputs | Keep route reduction; fixed output task after plan. |
| `content_index` | aggregate/per-output | scans posts/feeds/source state, writes content index | Keep aggregate due cross-content contract. |

## Cleanup

| Plugin | Classification | Conservative resources / side effects | Refactor required before fan-out |
| --- | --- | --- | --- |
| `css_minify` | per-output candidate | filesystem walk/read/write CSS outputs | Build immutable file list first, then one task per CSS path; no shared mutable minifier state. |
| `js_minify` | per-output candidate | filesystem walk/read/write JS outputs | Build immutable file list first, then one task per JS path. |
| `css_purge` | aggregate/external | scans generated HTML/CSS and rewrites CSS | Keep barrier until usage analysis can be separated from per-file rewrite. |
| `pagefind` | external/exclusive | subprocess scans final output and writes search index | Keep serialized external barrier with explicit output namespace. |
| `diagnostics_artifact` | exclusive | reads final lifecycle diagnostics/executor state, writes `.markata/diagnostics.json` | Must remain final publication barrier; never overlap unfinished build work. |

## First fan-out candidates

The first fan-out PRs should stay **serial** even after the graph contains item-level nodes. Resource claims should be added before any task is marked parallel-safe.

1. **`publish_html` — one task per post/output path.** High payoff and naturally disjoint output ownership. Refactor the plugin to produce an immutable publish plan, then expose a one-post writer. Claims: read `post:<id>` and `template:*`; write `output:<post-path>`.
2. **`publish_feeds` — one task per feed/format.** Feed objects already form natural units and outputs are usually disjoint. Claims: read `feed:<slug>`; write exact `output:<feed-path>`.
3. **`render_markdown` — one task per post.** CPU-heavy and item-local in intent. First make renderer/config state immutable and extract a post-local API. Claims: write `post:<id>` plus read-only plugin/config resources.
4. **`reading_time` — one task per post.** Small, low-risk proving ground for the per-post task API and graph fan-out mechanics. Claims: write `post:<id>`.
5. **`css_minify` / `js_minify` — one task per planned output file.** Straightforward filesystem ownership once directory walking is separated into a deterministic planning barrier. Claims: write exact `output:<path>`.

`auto_title`, `description`, `toc`, `heading_anchors`, `image_zoom`, and several HTML transforms are also strong per-post candidates, but they should follow the first shared per-post adapter rather than each inventing a scheduler integration.

## Required safety rules before concurrency

1. A fan-out task must name every mutable/shared resource it reads or writes.
2. Unknown ownership remains exclusive.
3. A write claim owns reads performed as part of mutating the same resource; a task must not declare duplicate read+write claims for one resource.
4. Collection planning/validation remains a barrier before item-level tasks consume the plan.
5. External/network work needs explicit cache ownership, cancellation, and a bounded rate/concurrency policy separate from CPU/file task concurrency.
6. Output tasks claim normalized destination paths/namespaces; overlap checks happen before scheduling.
7. Package/global template/browser/plugin caches must be proven read-only or made task-safe before post-level concurrency.
8. Build Lab differential/incremental/determinism checks and race tests remain mandatory gates for each fan-out slice.

## Next implementation sequence

1. Land deterministic resource claims (#1452).
2. Add compiled graph/resource diagnostics (#1454) so ownership is inspectable.
3. Introduce a serial item-task adapter and use `reading_time` as the smallest per-post proof.
4. Move `render_markdown` to serial per-post tasks and verify Build Lab equivalence.
5. Add immutable publish planning, then split `publish_html` and `publish_feeds` into serial per-output tasks.
6. Only after the graph and diagnostics prove disjoint ownership should any task become `ParallelSafe=true` or the executor accept `MaxParallel > 1`.

## Compiled graph core contract

The [core specification](../../spec/spec/BUILD_DAG.md) defines the snapshot and
diagnostic guarantees supporting this inventory (#1451, #1454):

- `Compile` snapshots artifact/resource slices and external inputs.
  After compilation, callers may reuse or mutate declaration arrays without
  changing the compiled graph. `Graph.Task` returns independent metadata slices,
  retaining the distinction between nil and empty slices. Task functions still
  own the mutable state captured by their closures; snapshotting declarations
  does not make plugin state immutable or parallel-safe.
- Artifact identity and canonical ordering use `(Kind, Key)`, not `String()`.
  For example, `post` / `rendered:example` and `post:rendered` / `example` are
  distinct artifacts even though their display strings match. Resource reports
  sort by `(Resource.Kind, Resource.Key, Access)`. Keep keys opaque; do not parse
  display strings to recover ownership.
- The compiler chooses the lexicographically smallest currently ready task,
  using reverse adjacency and a min-heap rather than rescanning every task.
  Fan-out/fan-in does not change this serial order or stable cycle diagnostics.
- Empty ownership remains valid for conservative compatibility tasks. Reports
  explicitly show empty claims as `[]`; serialization omits empty optional
  fields. Scope is descriptive and does not infer dependencies or claims.
- Structural digests identify declarations only. They do not prove that content
  or external state is unchanged and must not be used to skip execution.

These guarantees harden graph inspection and compilation, not concurrency:
`MaxParallel` remains one, existing scheduling defaults remain unchanged, and
task failures and cancellation continue to surface as errors.
