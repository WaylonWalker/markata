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

## Trusted ForwardAuth adapter (stacked #1956)

`NewForwardAuthAuthenticator` supplies the required identity hook **only when
the actual reverse proxy has authenticated and sanitized the request**.

Configure the exact trusted proxy source CIDRs, the stable subject-ID header,
the authentication-method assertion header, and one or more explicitly allowed
methods. The test fixture uses method `session`, but deployment must verify the
actual value emitted by the installed hlab-auth/ForwardAuth release. There are
no silent defaults or permissive fallback methods.

The adapter verifies the direct TCP peer from `RemoteAddr`, refuses broad
all-network/loopback/link-local trusted CIDRs, checks exact single-valued
assertion headers, and rejects missing, duplicated, malformed, or disallowed
authentication methods. It does **not** cryptographically validate a browser
token and cannot distinguish a forged assertion from a legitimate one if a
trusted proxy fails to strip a client-supplied header.

**Required ingress preconditions before mount:**
1. Traefik strips *every* configured identity/auth-method header supplied by
   browsers before ForwardAuth and forwards only successful auth response
   assertions; verify this with a live forged-header test.
2. The writer route requires a valid hlab-auth session or other specifically
   allowed method, with writer scope distinct from the Builder Admin operator
   scope. Do not reuse a builder operator bearer token in the editor.
3. A NetworkPolicy restricts direct requests to the designated Traefik pods;
   simply trusting a shared cluster Pod CIDR is **not** a secure boundary.
4. Hosted service exposure and Helm wiring remain disabled until full
   integration tests pass on a disposable site.

This stacked slice introduces no ingress, server listener, browser auth cookie,
write endpoint, or production deployment.
