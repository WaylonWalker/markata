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

The Python standard-library `sqlite3` test exercises the actual SQLite
schema, constraints, immutable revisions, and transaction rollback without
bringing a second Go SQLite dependency into every Markata build. The Go tests
check mutation validation and fingerprints.

## Next steps before enabling a sync endpoint

- Choose and pin a maintained Go SQLite driver after cross-platform build,
  dependency and performance evaluation; the current repository has no SQLite
  driver. Run Go integration tests against a **real temporary SQLite database**
  on restart/crash, concurrent mutations, and WAL replay.
- Configure `PRAGMA foreign_keys = ON` **per connection**; a file database
  must also enable WAL, a bounded busy timeout and a backup/restore policy.
- Implement transactional ApplyMutation that compares the exact base revision,
  persists document + immutable version + change event + mutation receipt in
  one commit, and acknowledges **only after that commit succeeds**.
- Reject mutation UUID reuse with a different digest, preserve conflicts, and
  preserve unsigned offline changes if a client checkpoint expires.
- Implement permission-filtered snapshots and feeds and subject-bound cursors;
  nothing in this package may bypass the authorization boundaries of #1951.
- Add a true offline browser outbox and attach it to this API only after
  server restart, interrupted acknowledgment and conflict tests pass.

No public ingress, Git operations, browser sync, or production family site is
modified by this slice.
