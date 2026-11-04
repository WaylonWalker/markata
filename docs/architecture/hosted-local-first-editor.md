# Hosted local-first editor architecture (POC contract)

Status: **proposed, not implemented**. Tracks #1949 and #1951. This document defines the first testable vertical slice, not a claim that the hosted workbench is safe to expose today.

## Outcome and deployment boundaries

One shared browser editor origin serves multiple sites. Each site retains an independent source repository and Builder Admin. A shared Go sync API brokers authenticated site-scoped reads/writes; clients never clone repositories or hold Git credentials.

```text
browser PWA -- local IndexedDB draft/outbox
    | HTTPS (authenticated site-scoped API)
    v
shared Go sync service -- transactional SQLite revision log
    | authenticated internal API, per-site capability
    v
site Builder Admin -- isolated Git checkout -- build queue -- public release
    ^                         |
    | signed webhook          | external Git pushes
    +-------------------------+
```

The current `markata editor` Go workbench is loopback-oriented, has same-origin protections, and explicitly lacks authentication/shared workspace management. Do **not** expose its file-writing or LSP endpoints through public ingress. Hosted mode needs separate route authorization and an allowlisted site registry. Builder Admin is currently an **operator control plane**; writer roles must never inherit operator capabilities.

## Identity and site authorization

* The trusted ingress authenticates with hlab-auth or another ForwardAuth provider. The sync service only trusts identity headers from a validated proxy path; reject direct connections and untrusted/duplicate identity headers.
* Use stable opaque subject IDs, not mutable usernames/emails, for ACL and draft ownership.
* Site registry binds an opaque site ID to an internal builder endpoint, allowed content roots, repository identity, and explicit `read`, `write`, `publish` grants. Evaluate permissions on **every** operation, including WebSocket/stream connections and media access.
* The site picker returns only authorized sites. A single-site user lands directly in their site. Never expose arbitrary filesystem paths or repository URLs in the browser.
* Offline mode may edit previously authorized local drafts, but server authorization is rechecked on reconnect. Logout clears sensitive cached content when requested; warn about unsynced drafts first.

## Draft and revision model

Distinct states: `local_saved` → `server_synced` → `git_committed` → `build_verified_live`. A synced draft is **not** published. Private drafts never enter a public Git branch implicitly.

Each document has `site_id`, immutable `document_id`, path, content hash, base revision, server revision, optional Git source commit, author subject and deletion marker. Rename preserves `document_id`; deleted documents remain tombstones until clients have advanced beyond them. Media are separate immutable content-addressed blobs with explicit attachment references.

Server assigns monotonically increasing **per-site** change sequence numbers in SQLite transactions. Every client mutation includes a stable UUID, base revision and payload hash. The server durably records a mutation receipt and revision before acknowledging it. Retrying the same mutation ID and payload returns the same receipt; reusing an ID with different content is rejected. Client checkpoints advance only after all earlier page entries are applied locally.

### Proposed wire contract (version 1; subject to spike)

* `GET /api/v1/sites` — authorized site summaries.
* `GET /api/v1/sites/{site}/changes?after=<seq>&limit=<n>` — bounded ordered changes, next checkpoint and pagination; a stale checkpoint returns an explicit resnapshot instruction, never an empty success.
* `POST /api/v1/sites/{site}/mutations` — mutation ID, document ID, base revision, operation, content or blob references. Returns durable receipt and assigned revision, or explicit conflict.
* `POST /api/v1/sites/{site}/media` — resumable upload session; blob not attached until checksum/size verified.
* `POST /api/v1/sites/{site}/publish` — expected synced revision(s), idempotency key and scope; returns durable operation ID.
* `GET /api/v1/sites/{site}/publications/{operation}` — queued, committed, building, failed or verified-live, including Git SHA and deployed release ID where available.

All requests are authenticated, authorized, size-limited and origin/CSRF-protected. No endpoint accepts an arbitrary repository path. Publish cannot claim success until the deployed revision is verified.

## Conflict handling and external Git changes

For an incoming mutation, compare the submitted base revision to the current server revision. If stale, preserve the original and return both revisions; attempt a deterministic three-way merge only when a common ancestor is available and the result is unambiguous. Otherwise create a reviewable conflict copy, not a last-writer-wins overwrite. Never discard local content on a 409.

External Git pushes still reach Builder Admin through signed webhooks. The sync service must import the changed Git tree as a new source revision, generate incremental document events and reconcile against pending drafts. Use stable content hashes and Git SHAs for idempotent ingestion. A build is deduplicated by site + commit + build configuration, while distinct publication operations retain their own receipts.

Publish uses an isolated checkout/worktree and explicit changed-path allowlist. It must not accidentally commit unrelated source changes or overwrite a newer upstream branch; on push race, fetch and reconcile or return a retryable conflict.

## Rendering and offline behavior

Preview is a separate interface from persistence and synchronization. `Render(markdown, siteConfig, revision)` can target authenticated server Go or cached WASM using the same Markata rendering rules. Reject stale responses by revision. Do not assume parity until tested: compare HTML/asset output for a representative corpus, including embeds, Mermaid, frontmatter and site themes.

POC selection: prefer Go while online; fall back to WASM on request failure/offline after assets are cached. Expose a diagnostic override for developers. Latency-based adaptive selection requires later benchmarks on mobile and desktop, accounting for network RTT, bundle load and memory use. If WASM is not cached, preserve editable Markdown and display the last cached preview rather than block writing.

Service worker caches the application shell and versioned renderer assets, **not** authenticated API responses by default. IndexedDB stores the outbox and local document snapshots with atomic writes. A background sync event is an optimization only: resume on app open, network reconnect and explicit retry.

## Minimal acceptance suite

1. Sign in as a single-site writer; no vault, GitHub token, or repo URL. A different site's content and builder operator routes return 403.
2. Edit offline, refresh/force-close, reopen offline: exact bytes and pending mutation IDs survive.
3. Restore connection; retry the same mutation twice after simulated lost acknowledgment: exactly one server revision.
4. Edit on device A while B changes the same post: no silent overwrite; both versions recoverable.
5. Push an external Git commit: it appears incrementally in a clean client without re-downloading a whole large repository; an unsynced local edit is preserved.
6. Interrupt server during SQLite commit/ack and restart: either a durable receipt exists or client safely retries.
7. Upload media with interrupted chunks and verify checksum; no broken published attachment references.
8. Preview while connected through Go, then offline through cached WASM; edits never block on rendering.
9. Publish a specific synced revision; verify Git commit and actual live release; a failed build reports failure without losing drafts.
10. Repeat with a disposable large site (thousands of posts), measure initial index transfer, incremental bytes, UI responsiveness and storage consumption.

## First implementation slices

1. **Site discovery and authorization**: hosted-only API on a disposable site; leave local workbench unchanged.
2. **Durable sync spike**: Go + SQLite schema, one Markdown document, idempotent mutations, checkpoints and offline browser outbox.
3. **External Git ingestion**: incremental source refresh and conflict preservation.
4. **Publish integration**: writer-scoped action through existing Builder Admin and verified release receipt.
5. **Hybrid rendering**: shared preview adapter and parity fixtures.
6. **Family pilot**: phone photo posting, large-site indexing, accessibility and recovery drills.

Open decisions for the spike: RxDB versus a minimal IndexedDB outbox (choose by failure tests and bundle size); shared central SQLite versus site-scoped databases (evaluate backups and isolation); browser WASM renderer parity; user-visible conflict resolution UI; exact internal builder authentication protocol.

## Existing entry points

* `docs/guides/go-workbench.md`: editor, LSP, local file saves and Go preview.
* `docs/guides/deployment/builder-admin.md`: authenticated operator service, signed webhooks, build queue, verified releases.
* #1950: Plaindown/Markata feature audit; #1854, #1856, #1866: renderer parity and offline assets; #1446: builder recovery.

No production family content should be modified in this POC. Keep Plaindown running until all recovery and publication tests pass.
