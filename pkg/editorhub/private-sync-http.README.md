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

Conflict-copy retrieval from draft PR #1975 remains separate, and will need
explicit routing review before being added to this composition.

Verify with `go test -race -count=1 ./pkg/editorhub ./pkg/editorsync`.
