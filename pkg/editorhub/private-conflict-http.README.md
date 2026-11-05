# Private conflict recovery HTTP slice (#1970)

This is an **unmounted**, opt-in extension to the Go private mutation handler.
It does not deploy a public route or expose the trusted-local workbench.

`GET /api/v1/sites/{site}/conflicts/{mutation}` returns one author-scoped
rejected mutation (including orphaned source for documents now absent). The
handler requires the verified server-side subject, a server-owned site write
grant, and an author match inside `Engine.ReadPrivateConflict`. A writer on
the same site **cannot** read another writer's private rejected source, even
though both hold site write capabilities. Missing and unauthorized-owner
conflict IDs return indistinguishable 404 responses.

The response has `mutation_id`, `document_id`, `stale_base_revision`,
`current_revision`, `digest`, `source_base64` (verbatim bytes), and
`deleted`; it is `Cache-Control: private, no-store`. The handler does not
return internal site-wide change sequence values.

Cross-site Fetch Metadata and a supplied incorrect/duplicated Origin are
rejected; a GET need not include an Origin header. This is not a replacement
for the production ForwardAuth/session check, network peer ACLs, or client
conflict UI. No automatic conflict resolution or last-writer-wins behavior
is added.

Tests verify same-site writer isolation, site readers without write access,
missing login, forged HLab headers, exact source recovery, orphaned edits,
and cross-origin failures. Run:

```sh
go test -race -count=1 ./pkg/editorhub ./pkg/editorsync
```

**Deployment remains blocked** until #1956 verifies actual hlab-auth ingress,
headers, trusted peer and session revocation; #1967 finishes safe browser
outbox/reconnect and private device storage; and #1970's full read/write
routing is reviewed on a disposable origin.
