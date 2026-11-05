# Offline draft outbox — persistence foundation (PR #1967)

**Not enabled in the public workbench or family sites.**

`offline-outbox.mjs` is a small native IndexedDB adapter and testable
local-draft state machine. No framework, GitHub credentials or network
connection is required. `saveDraft` resolves **after the database transaction
commits**; only then may the editor claim “Saved on this device.” A quota
failure or aborted transaction rejects the promise and must show a warning,
not a false saved indicator.

`enqueue` snapshots the exact locally committed source, base revision,
deletion intent and a stable mutation UUID. It blocks multiple unresolved
mutations on one document. A retry reuses the frozen payload, even if the user
has typed newer content. `acknowledge` updates the server revision only for
a matching committed receipt, without overwriting newer local edits.
`recordConflict` retains the original proposal for explicit reconciliation.

All reads and writes are keyed by *site ID, authenticated stable subject ID,
and document or mutation ID*. The IDs must come from a trusted authenticated
bootstrap, **never from an arbitrary user-editable browser setting**.

## Verification

```sh
node --test pkg/workbench/web/offline-outbox.test.mjs
```

The tests use a fault-injectable in-memory transactional driver to check 100
offline edits, hard-reload semantics, idempotency, failed commits, stale
replies, conflicts and logical per-user/site partitioning.

**This is not a cryptographic account isolation boundary.** IndexedDB is
origin-scoped and contains plaintext drafts. Any script running on the same
origin can inspect its records. Before shared-device/family deployment, add a
session-bound local lock/encryption design, real browser IndexedDB upgrade and
transaction tests, storage eviction/quota UX, per-account logout and “Forget
this device” flows, Web Locks single-coordinator sync, and a secure hlab-auth
bridge. Do not mount this module in a publicly reachable editor until those
requirements are verified. “Saved on this device” is not a server backup or a
published post; users can still lose it if the browser clears local data.

This slice neither adds any HTTP route nor sends source to the server.
