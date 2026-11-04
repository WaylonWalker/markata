# Editor sync storage foundation

Issue: #1958; umbrella: #1949. **This is not a running sync service.**

This package introduces the SQLite schema for private and shared draft
revisions, idempotent mutation receipts, a per-site internal change log, and
reviewable conflicts. The Go contract validates mutation metadata and
computes a length-framed fingerprint covering the complete mutation intent.

The schema separates immutable document history (`editor_document_versions`)
from the current materialized document (`editor_documents`). This allows a
future paginated snapshot to hydrate the exact revisions visible when the
snapshot cursor was issued, even when changes arrive concurrently. Client-facing
resume cursors must remain opaque and bound to the authenticated subject.
Queries MUST enforce document and media visibility; SQL tables alone do not
provide authorization.

## Verification

```sh
go test ./pkg/editorsync
python3 -m unittest discover -s pkg/editorsync -p 'test_schema.py' -v
```

The Python standard-library `sqlite3` test exercises schema constraints,
immutable revisions, and rollback. The Go tests additionally exercise the real
pure-Go SQLite driver, schema migration, file restart and receipt persistence.
`OpenSQLite` is the opt-in file-backed opener: it requires an existing trusted
parent directory, creates new database files with owner-only permissions,
enables WAL/foreign keys/bounded busy timeout, and serializes initial writes
through one connection. Existing database-file permissions and backups remain
an operator responsibility. This may increase test/binary build costs, so
cross-platform release and CI impact must be reviewed.

## Private-draft transaction POC (stacked PR #1961)

`NewEngine` requires both an authenticated, caller-supplied site write
authorizer and a SQLite handle. `Engine.Apply` currently accepts **private
drafts only**, enforces owner identity and base revision, and atomically stores
the current document, immutable version, per-site change record and mutation
receipt. A repeated UUID and unchanged payload returns the original durable
receipt, even after database restart; a changed payload is rejected. Server
acknowledgment means *SQLite committed*, not Git publication.

The Go tests simulate interrupted acknowledgment and retry, concurrent
creation, stale revisions, identity spoofing, and a forced receipt-write failure.
For stale edits to an **existing private document**, the server now commits
the exact rejected source, full mutation fingerprint, deletion intent and
observed revision to `editor_conflicts` **before** returning a
`ConflictError{Preserved: true}`. The same UUID/intent replays the original
conflict, while a changed intent is rejected. `ReadPrivateConflict` retrieves
a copy only for the original authenticated subject; no HTTP route is mounted.

Orphaned mutations whose document no longer exists are stored in the
`editor_orphan_conflicts` journal instead of weakening the original foreign
key. Their exact contents, deletion intent and mutation receipts remain
recoverable after restart, scoped to the authenticated author. Conflicts do
not create live document versions or change-feed entries.

This is still not family-grade multi-device sync: a browser offline outbox,
authorized change feed, conflict resolution UI, backup/restore and
token-revocation handling remain unimplemented.

## Next steps before enabling a sync endpoint

- Finish assessing the pinned `modernc.org/sqlite v1.61.1` driver's
  cross-platform build and binary-size impact. Add restore/WAL replay and
  multi-process busy-contention tests before production use.
- Configure and validate backup/restore, media storage and site isolation.
- Extend transactions to shared drafts and Git-imported source revisions,
  and implement an explicit acknowledged conflict-resolution workflow.
- Preserve unsigned offline changes if client checkpoints expire, then
  reconcile conflicts explicitly without silently overwriting content.
- Implement permission-filtered snapshots and feeds and subject-bound cursors;
  nothing in this package may bypass the authorization boundaries of #1951.
- Add a true offline browser outbox and attach it to this API only after
  server restart, interrupted acknowledgment and conflict tests pass.

No public ingress, Git operations, browser sync, or production family site is
modified by this slice.
