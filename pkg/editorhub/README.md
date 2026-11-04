# Hosted editor site-discovery foundation

Tracking: #1949, #1951, #1954. Architecture review: PR #1953.

`pkg/editorhub` contains an **unmounted, opt-in** HTTP handler for a hosted
Markata site picker. This is the first safe boundary of the future editor hub,
not a production host or a replacement for the existing `markata editor`
workbench.

## Contract

- A trusted application (not this package) authenticates the request, validates
  that it traversed a trusted ingress/ForwardAuth path, rejects browser-supplied
  identity headers, and supplies a stable subject through an
  `editorhub.Authenticator`. **Never authenticate merely by reading
  `X-Hlab-*`/Forwarded headers.** A nil authenticator is rejected.
- `NewRegistry` accepts trusted server-side site configuration. It rejects
  duplicate/invalid IDs, wildcard subjects and impossible permission sets.
  It copies grants so mutating the original configuration cannot grant access.
- `Registry.Allows(subject, siteID, permission)` supplies a fail-closed site capability gate for future handlers. This is **not** a substitute for per-document private-draft/media authorization.
- `GET /api/v1/sites` returns only explicitly readable sites, sorted by
  opaque ID. Each entry contains only `id`, `name` and `permissions`.
  It returns no repository, internal builder, filesystem or draft metadata.
- Unauthenticated requests return 401. Unsupported methods return 405.
  Non-discovery paths return 404; in particular **no workbench/LSP/editor file
  routes or Builder Admin operator API are mounted**.
- Responses use `private, no-store`; CORS credentials/wildcards are not
  enabled.

The authenticator is a **trust boundary to be implemented and reviewed in a
separate PR**, with a concrete proxy/network policy and integration tests.
Never mount this handler on an externally reachable server with a test
authenticator or blindly trusted request headers.

## Test and rollout notes

Run `go test ./pkg/editorhub`. Tests exercise per-user visibility within a
shared site, missing identity, header spoofing, invalid grants, deterministic
ordering, immutability, and refusal to serve privileged routes.

This package is deliberately **not connected to an ingress, CLI, the Go
workbench, or Builder Admin**. The next POC slices are secure hosted identity
verification and then a separate durable, offline draft/sync protocol. No
production family site has been modified or exposed.
