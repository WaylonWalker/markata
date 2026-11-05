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

## Private mutation HTTP transport POC (#1970)

`NewPrivateMutationHandler(registry, authenticate, engine, allowedOrigin)`
returns one **unmounted**, private-draft-only route:
`POST /api/v1/sites/{site}/mutations`. No LSP, local filesystem, Git, builder
admin, publish, list or public site endpoint is added here.

The handler requires an injected **trusted** `editorhub.Authenticator` that
derives a stable subject from a proven ingress decision, and an independent
explicit site `write` grant. Browser JSON cannot supply or override the
subject/site/visibility; these are taken from the verified context, route and
server-owned private-draft policy. The underlying `editorsync.Engine` checks
the same author identity/permission and transactionally commits before
returning a receipt. Replays of the same UUID return the original receipt;
stale revisions return a recoverable conflict response.

Browser writes must have an exact configured HTTPS `Origin` and
`X-Markata-Sync-Intent: write`. Unsupported origins, missing/duplicate
headers, cross-site fetch metadata, absent login, read-only users and
unexpected JSON fields fail closed. This is an additional CSRF barrier, not a
replacement for a validated hlab-auth browser session, trusted proxy and
NetworkPolicy. The handler does not enable permissive CORS.

Successful responses expose only the mutation ID, accepted document revision
and intent hash. Internal per-site change sequences are never returned:
they may contain other users' private activity.

Run `go test -race ./pkg/editorhub ./pkg/editorsync`. The handler is not
installed in any server/CLI/ingress. It is **not safe to mount publicly** until
the #1956 live ForwardAuth/Origin/CSRF tests prove trust boundaries, and the
read endpoints, log retention, offline browser outbox and conflict UX are
reviewed independently.
## Private read HTTP POC (#1970)

`NewPrivateReadHandler` is an opt-in, **unmounted** HTTP adapter for
`editorsync.PrivateReader` (private drafts only). It requires an injected
verified subject authenticator, server-owned site registry with explicit
read grants, and the private SQLite reader with its own per-request
authorization. This does not mount the loopback workbench, Git, Builder
Admin or any site file route.

Routes are `POST /api/v1/sites/{site}/snapshots` (first metadata page),
`GET /api/v1/sites/{site}/documents?cursor=...&limit=...` (continued
metadata pages), `GET /api/v1/sites/{site}/changes?cursor=...` (incremental
owner-private events and tombstones), and
`GET /api/v1/sites/{site}/documents/{id}?revision=...` (on-demand exact
source bytes encoded as base64). Pages are capped at 100 items.
The wire DTOs use explicit snake_case field names and never emit internal
per-site change sequence numbers, other users' source, or absolute paths.

The initial snapshot POST requires an exact, configured HTTPS `Origin`.
Every request checks authenticated subject and site permissions; the
underlying read engine also enforces document ownership. An invalid or
other-user cursor is rejected, while expiry returns HTTP 410 with
`{"code":"checkpoint_expired"}` to trigger a new authorized snapshot.
The browser MUST keep its unsent outbox during that resnapshot.

**Not ready for hosting:** no actual hlab-auth ingress, distributed
session invalidation, retention-backed change history, shared/published
document transitions, browser-side reconnect, media, or performance
benchmark on thousands of posts. Before enabling public ingress, validate
#1956's peer/header stripping and session assertions and review HTTP
cache, CSP, Origin and CORS behavior in the full deployment. Never use a
test authenticator in a deployed service.
