# Composed private sync HTTP proof of concept

`NewPrivateSyncHandler` combines the existing site discovery, private
mutation write, snapshot, paginated manifest, lazy source hydration and
incremental private change feed. It is **not mounted** in any production,
hosted or local editor server.

The constructor requires a validated site registry, trusted Authenticator,
SQLite mutation engine, private reader with encrypted cursors, and an exact
HTTPS editor origin. It fails closed when dependencies are missing. Each
inner handler evaluates site permissions again, and private source bytes are
accessible only to their original author.

In-process integration tests with disposable SQLite exercise initial writes,
idempotent lost-response replay, other-author isolation, snapshot checkpoint,
historical hydration, subsequent change feed, forged browser identity headers
and unknown operator paths.

This handler has no Git, publish, file system, Builder Admin or LSP route.
It does not wire `hlab-auth`, open an ingress or write family content.
Before a disposable hosted pilot, #1956 must prove real ForwardAuth header
stripping, trusted proxy peer and NetworkPolicy, session assertion,
revocation, CSRF and same-origin behavior. Browser-local private drafts
remain plaintext until the separate device-account isolation work is done.

Owner-only conflict-copy retrieval from merged PR #1975 is included here.
Existing and orphaned rejected drafts can be recovered only by their original
authenticated author; a site write grant alone cannot expose another writer's
private conflict source. This is recovery, not automatic conflict resolution.

Verify with `go test -race -count=1 ./pkg/editorhub ./pkg/editorsync`.

## Verified browser identity bootstrap (issue #1979)

The opt-in `NewPrivateSessionHandler` adds a read-only
`GET /api/v1/session` endpoint to this composed handler. It returns the
current verified stable `subject_id` and the caller's currently authorized
site summaries. It does **not** return or store a session cookie, JWT,
ForwardAuth header, GitHub credential, storage path or operator token. The
response is `private, no-store`, `nosniff`, and same-origin-only; the
handler derives the subject freshly via the injected trusted Authenticator
for every request and never accepts an identity from browser-supplied JSON.

A future browser adapter must fetch this endpoint with same-origin credentials
and `cache: no-store`; checking a previous subject stored in IndexedDB is
not proof of an active session. On logout, session expiry, account switch or
revoked site grant, network sync must pause without discarding the outbox.
No browser adapter, public ingress or production route is added in this PR.

**Note:** verified identity for HTTP synchronization does not protect
plaintext IndexedDB drafts from other people using the same browser profile.
Device-local encryption/locking and logout/export remain separate required
work before any family/shared-device deployment.
