# Offline mutation sync drain POC

Tracks #1967 and #1949. This module is **not loaded by the public workbench**,
does not own login/session handling, and is not a production sync client.

`offline-sync.mjs` consumes already-committed, immutable queued mutations from
`LocalDraftOutbox`. A hosting app must inject an identity callback whose
`{siteId, subjectId, verified: true}` was derived from the server's authenticated
identity endpoint—not from localStorage, user-entered text, an email, or a
client-controlled header. It must revalidate identity when a session changes.

`OfflineMutationSync.syncOnce()` uses an exclusive `navigator.locks` lock for
one site/subject. If Web Locks is unavailable or another tab owns the lock,
no upload occurs. The drain is private, so callers cannot bypass the lock.
Every queued mutation is posted verbatim (same mutation ID, base revision,
source bytes represented by JSON, and deletion intent) to the **relative
same-origin** private write route, using `credentials: same-origin`,
`mode: same-origin`, `redirect: error`, `cache: no-store` and a custom write
intent header. A session check is repeated before each upload and after each
response, including account switching while a request is in flight.

A 200 response is accepted only when its `mutation_id`, document revision
and SHA-256 intent digest **match the frozen original proposal**, byte-for-byte
with Go's length-framed `editorsync.Mutation.Digest`. Only then is the
matching local SQLite receipt recorded into IndexedDB. A lost network reply
retains the same mutation for safe idempotent retry. 401/403 pause sync,
409 with `preserved: true` moves the exact draft to a separate conflict
recovery list, and errors leave the outbox unchanged. `synced_drafts` means
the server accepted a private draft; it does **not** mean published.

Run the dependency-free Node tests:

```sh
node --test pkg/workbench/web/offline-outbox.test.mjs pkg/workbench/web/offline-sync.test.mjs
```

They cover the Go SHA-256 digest contract, lost-response replay, concurrent
browser-tab lock behavior, 401/403/network failures, forged receipts, retained
conflicts, missing Web Locks, and logout during an in-flight request.

## Deferred before hosted rollout

- Bind to hlab-auth **only after** #1956 verifies the actual proxy
  boundary and server identity behavior. This module is deliberately not
  wired into an existing app, UI, service worker or public ingress.
- Add reconnect/visibility triggers, bounded jitter/backoff, explicit
  `Sync now` UX, and safe offline state indicators.
- Browser-local private drafts remain plaintext IndexedDB. Build a
  session-bound local lock/encryption and multi-user logout/forget/export
  story before deploying to shared family devices.
- Add real-browser Web Locks concurrency tests, multi-device conflict
  resolution, expired-cursor full resnapshots preserving the outbox, media
  hydration, backup/restore and large-site performance acceptance.
