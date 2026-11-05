# Stateless verified identity for offline sync

`fetchVerifiedSyncIdentity(siteId)` is a small, dependency-free browser helper
for #1979. It calls the same-origin `GET /api/v1/session` endpoint using
`credentials: same-origin`, `mode: same-origin`, `redirect: error` and
`cache: no-store`. It accepts only a current, valid opaque subject ID and
exactly one authorized site with both read and write grants.

Inject it as `verifiedIdentity: () => fetchVerifiedSyncIdentity(siteId)`
into `OfflineMutationSync`. The drain already rechecks identity before
sending each frozen mutation and after asynchronous boundaries. Logout,
account switching, lost connectivity, read-only grant, malformed JSON and
session expiry return `null`, causing sync to pause without deleting the
IndexedDB outbox. The helper stores no credentials, tokens, session state or
positive identity verdict in localStorage/IndexedDB.

Node tests cover fresh repeated calls, read-only site exclusion, account
switches, session/network failures, malformed data, and a fake end-to-end
`OfflineMutationSync` acceptance. Run:

```sh
node --test pkg/workbench/web/offline-session.test.mjs
```

**Security boundary:** the server endpoint is intentionally unmounted until
#1956 proves ForwardAuth header stripping, trusted proxy IP/NetworkPolicy,
session method, revocation and exact origin/CSRF handling. The browser fetch
helper does not turn plaintext IndexedDB into secure multi-user local storage.
Do not integrate this with a public editor origin or share a browser profile
among different family accounts until cache encryption/locking and logout
recovery UX are independently verified.
