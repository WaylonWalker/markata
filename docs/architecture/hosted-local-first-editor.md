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

### Privacy model and scoped change feeds

There are three visibility classes: `private` (owner only), `shared_draft` (explicitly authorized site collaborators), and `published_source` (available to authorized site readers; actual public visibility is determined by the site's publishing rules). The server—not the browser—enforces document/media ACLs for discovery, read, write, change feed, conflict copies, and publication. Private drafts and their attachments **must never** be returned merely because a user has site-level read access. The same ACL applies to historical revisions, search/autocomplete, conflict metadata, and document existence.

A single per-site sequence may order events internally; user-visible change pages filter entries *before* serialization using current subject, site role, and document visibility. A page may contain zero visible changes and still advance its `next_checkpoint` across hidden events; clients never infer hidden document IDs, counts, or existence from sequence gaps. Permission grants and revocations cause a scoped re-bootstrap/visibility refresh, not silent reliance on old checkpoints. `403` responses must not reveal whether an inaccessible document exists. The POC must include two collaborators on the same site, with private drafts and media from each writer remaining invisible to the other.

SQLite rows are retained for server-durable drafts and source revisions; an eviction policy only applies to superseded change-history entries and tombstones. Never trim history based on all clients acknowledging, since offline devices may never return. Define an operator-configurable retention floor (initial proposed minimum **30 days**, subject to backup/volume tests), and report `checkpoint_expired` with an authorized snapshot cursor. The client must keep its unsent outbox and any conflict copies while rebuilding its local index.

### Proposed wire contract (version 1; subject to spike)

* `GET /api/v1/sites` — authorized site summaries.
* `POST /api/v1/sites/{site}/snapshots` — return an opaque, time-bounded snapshot token and captured per-site high-watermark `as_of_seq`. Server pins an internally consistent logical view (or materializes a bounded snapshot manifest) for the token lifetime; it must not mix live index pages from different revisions.
* `GET /api/v1/sites/{site}/documents?snapshot=<token>&cursor=<opaque>&limit=<n>` — paginated **metadata-only** authorized index, stable throughout the snapshot. Return document IDs, titles, paths, content hashes, revisions and `next_cursor`, not entire large-site post bodies.
* `GET /api/v1/sites/{site}/documents/{document_id}?revision=<revision>` — on-demand content hydration of an authorized revision. Snapshot pages identify the exact revision to fetch; an expired revision returns a typed resnapshot error, not a silently newer body.
* `GET /api/v1/sites/{site}/changes?after=<seq>&limit=<n>` — bounded ordered **authorized** changes, next checkpoint and pagination. After applying the full snapshot, start at the snapshot's `as_of_seq` (not the current head) to avoid missing changes between index and feed. A stale or permission-incompatible checkpoint returns `checkpoint_expired` / `visibility_changed` and triggers a fresh snapshot while retaining the local unsent outbox. Empty visible pages may still advance checkpoints across filtered events.
* `POST /api/v1/sites/{site}/mutations` — mutation ID, document ID, base revision, operation, content or blob references. Returns durable receipt and assigned revision, or explicit conflict.
* `POST /api/v1/sites/{site}/media` — resumable upload session; blob not attached until checksum/size verified.
* `POST /api/v1/sites/{site}/publish` — expected synced revision(s), idempotency key and scope; returns durable operation ID.
* `GET /api/v1/sites/{site}/publications/{operation}` — queued, committed, building, failed or verified-live, including Git SHA and deployed release ID where available.

All requests are authenticated, authorized, size-limited and origin/CSRF-protected. No endpoint accepts an arbitrary repository path. Publish cannot claim success until the deployed revision is verified.

## Conflict handling and external Git changes

**Source-of-truth separation:** Git's configured branch is authoritative for the **publishable site source** and build history; SQLite is authoritative for **unpublished, server-accepted drafts, private versions, sync receipts and document visibility**. Device storage remains authoritative for changes that the server has not acknowledged. A server-side draft must never be treated as automatically Git-committed or publicly visible.

For an incoming mutation, compare the submitted base revision to the current server revision. If stale, preserve the original and return both revisions; attempt a deterministic three-way merge only when a common ancestor is available and the result is unambiguous. Otherwise create a reviewable conflict copy, not a last-writer-wins overwrite. Never discard local content on a 409.

External Git pushes still reach Builder Admin through signed webhooks. Builder Admin currently tracks one branch/checkout and performs a fast-forward-only pull; webhook delivery can be delayed or missed. The sync service must ingest only **confirmed** Git revisions (or content-tree deltas), preserving Git commit SHA, document path and content hash, before publishing new `published_source` change events. Import is idempotent by site + branch + SHA. Periodic bounded reconciliation against the authoritative remote branch or builder source SHA is a required missed-webhook recovery path; new external changes become visible in the editor without the client downloading the whole repository. The service must not overwrite unpublished drafts when importing Git changes.

A build is deduplicated only when site, source commit SHA, build configuration and requested target are equivalent. Existing Builder Admin may enqueue an additional watcher-triggered build after a webhook pull; deduplication across both trigger paths is an **implementation follow-up**, not a presently guaranteed behavior. Distinct publication operations retain their own durable receipts.

### Scoped publish and builder trust

**The editor/sync service authenticates the writer, but never passes that writer an operator token.** Today's Builder Admin `/api/v1/` is an operator API, requires bearer ForwardAuth and appropriate operator scope, and is **not** directly usable as a writer publish API. The POC needs a new internal, independently authorized publish capability (a narrow endpoint or service identity with exact site/paths/revisions and permission checks); it cannot proxy arbitrary operator endpoints. Service-to-service identity is short-lived/scoped (or mutually authenticated), bound to site and operation, checked against replay and recorded in audit events. No builder credential is stored in the browser.

The **site builder** owns an isolated checkout/worktree, selective file staging, commit, push, build enqueue, and deployed-revision verification. The sync service owns the durable publish request, permission evaluation, snapshot of exact server-synced draft revision(s), and receipt exposed to the browser. Publishing requires an explicit `expected_draft_revision`, `expected_source_commit` and idempotency key; builder rechecks these just before Git changes. The operation must not commit unrelated work or private draft material. On upstream push races, fetch and reconcile against the new source or return a conflict preserving the draft. Existing content publishing flags and Markata's review UI remain authoritative for public/private decisions.

Publish is a durable state machine (`accepted → preparing → git_committed → push_confirmed → building → verified_live`, plus retryable and terminal failures). After a timeout, the caller polls by idempotency key/operation ID; neither end retries a Git commit blindly without checking remote SHA and durable receipt. If Git push succeeds but build fails, retain both Git SHA and draft revision and report `build_failed` instead of `published`. Confirmation of the **live source revision** uses Builder Admin's deployment record, not the push acknowledgment. Atomic promotion and build restart recovery remain Builder Admin concerns (#1446).

## Rendering and offline behavior

Preview is a separate interface from persistence and synchronization. The proposed `Render(markdown, context, revision)` contract can target authenticated server Go or cached WASM using the same Go Markdown core. Reject stale responses by revision. **This is not yet equivalent to a full Markata site build:** standalone WASM cannot execute site discovery, full frontmatter processing, feeds, or site-dependent plugins; cross-document embeds need precomputed authorized site metadata. Define capabilities per mode: pure Markdown, highlighting, basic containers and cached assets can be previewed locally where supported; site-dependent lookups require the server or a bounded, permission-filtered cached index; native-only rendering and uncacheable resources must report `unsupported` or an explicit approximate-preview indicator, never silently present a different result as exact parity. Compare HTML/asset manifests via #1866, including Mermaid, frontmatter and themes. Rendered HTML from untrusted authors stays sandboxed with scripts/unsafe HTML disabled by default.

POC selection: prefer Go while online; fall back to WASM on request failure/offline after assets are cached. Expose a diagnostic override for developers. Latency-based adaptive selection requires later benchmarks on mobile and desktop, accounting for network RTT, bundle load and memory use. If WASM is not cached, preserve editable Markdown and display the last cached preview rather than block writing.

Service worker caches the application shell and versioned renderer assets, **not** authenticated API responses by default. IndexedDB stores the outbox and local document snapshots with atomic writes. A background sync event is an optimization only: resume on app open, network reconnect and explicit retry.

## Minimal acceptance suite

1. Sign in as a single-site writer; no vault, GitHub token, or repo URL. A different site's content and builder operator routes return 403. Two people on the **same** site cannot enumerate, search, read, or sync each other's private drafts or media.
2. Edit offline, refresh/force-close, reopen offline: exact bytes and pending mutation IDs survive.
3. Restore connection; retry the same mutation twice after simulated lost acknowledgment: exactly one server revision.
4. Edit on device A while B changes the same post: no silent overwrite; both versions recoverable.
5. Push an external Git commit: it appears incrementally in a clean client without re-downloading a whole large repository; an unsynced local edit is preserved. Simulate lost webhook and prove periodic Git reconciliation catches up.
6. Interrupt server during SQLite commit/ack and restart: either a durable receipt exists or client safely retries.
7. Upload media with interrupted chunks and verify checksum; no broken published attachment references.
8. Preview while connected through Go, then offline through cached WASM; edits never block on rendering.
9. Publish a specific synced revision; verify Git commit and actual live release; a failed build reports failure without losing drafts.
10. Repeat with a disposable large site (thousands of posts): paginate an authorized, stable snapshot index; fetch only selected document bodies; mutate during snapshot and prove the following change feed has no gap. Measure index transfer, incremental bytes, UI responsiveness and storage consumption.
11. Expire a stale checkpoint and revoke a collaborator's access: client re-snapshots without losing unsynced edits, and the revoked subject cannot access cached server drafts via API.

## First implementation slices

1. **Site discovery and authorization**: hosted-only API on a disposable site; leave local workbench unchanged.
2. **Durable sync spike**: Go + SQLite schema, one Markdown document, idempotent mutations, checkpoints and offline browser outbox.
3. **External Git ingestion**: incremental source refresh and conflict preservation.
4. **Publish integration**: writer-scoped action through existing Builder Admin and verified release receipt.
5. **Hybrid rendering**: shared preview adapter and parity fixtures.
6. **Family pilot**: phone photo posting, large-site indexing, accessibility and recovery drills.

Open decisions for the spike: RxDB versus a minimal IndexedDB outbox (choose by failure tests and bundle size); shared central SQLite versus site-scoped databases (evaluate backups and isolation); user-visible conflict resolution UI; exact service-to-service authentication mechanism for the new narrow publish capability (no reuse of the operator credential); revision retention and backup restoration tests. The protocol above is a **proposal**, not implemented endpoints.

## Existing entry points

* `docs/guides/go-workbench.md`: editor, LSP, local file saves and Go preview.
* `docs/guides/deployment/builder-admin.md`: authenticated operator service, signed webhooks, build queue, verified releases.
* #1950: Plaindown/Markata feature audit; #1854, #1856, #1866: renderer parity and offline assets; #1446: builder recovery.

No production family content should be modified in this POC. Keep Plaindown running until all recovery and publication tests pass.
