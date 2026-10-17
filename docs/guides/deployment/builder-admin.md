---
title: "Builder Admin"
description: "Run a long-lived builder admin service for fast Kubernetes authoring loops"
date: 2026-06-27
published: true
tags:
  - documentation
  - deployment
  - kubernetes
  - performance
---

# Builder Admin

The builder admin service keeps a warm markata-go build worker running inside your cluster.

Instead of starting a new Kubernetes Job for every authoring build, it keeps one HTTP service alive with:

- a serialized build queue
- file watching that enqueues builds
- build history and full raw logs
- release history and current live release
- the sanitized content diagnostics artifact from successful full builds
- rollback by promoting an older rendered release
- scheduled refresh tasks for reader, blogroll, or other remote-content commands
- signed GitHub and Forgejo push webhooks that pull and build matching branches

## What It Is Good For

Use builder admin when:

- your site content already lives on a mounted filesystem such as a hostPath
- you care about fast authoring loops more than one-shot batch builds
- you want an operator UI for builds and releases
- you want remote-content refreshes to stay decoupled from normal content builds

## Basic Helm Values

```yaml
builderAdmin:
  enabled: true
  port: 8080
  fast: false
  watch:
    enabled: true
    debounce: 2s
  releases:
    keep: 25
  history:
    successfulBuilds: 60
    failedBuilds: 100
    refreshRuns: 100
  refreshTasks:
    - name: reader-update
      every: 30m
      enqueueBuildOnSuccess: true
      args:
        - markata-go
        - --config
        - /data/source/markata-go.toml
        - reader
        - update
```

### Put transient build work on node-local storage

By default Builder Admin preserves the historical behavior and creates `.build-work` under the
release directory. For a write-heavy site whose durable site volume is slow at small-file I/O,
move only the transient workspace onto node-local storage:

```yaml
builderAdmin:
  workspace:
    enabled: true
    mode: emptyDir
    mountPath: /data/work
    workDir: /data/work/build
    emptyDir: {}
```

`emptyDir` is node-local and disappears with the pod. The next build safely seeds
a new workspace from the current release. Retained releases and the `current` symlink remain under `release-dir` on the
durable site volume.

Builder Admin retains an independent mutable workspace after publication. It creates
a hidden staged release and compares its files with the previous immutable release.
Matching file contents and permissions reuse that release's inode through a hard
link. Changed files use independent copies. Unsupported hard links fall back to
copies. The workspace never shares writable file inodes with published releases.

Publication renames the complete staging directory into its retained release path
before switching `current`. Deleted workspace paths remain absent from the new
release. A failed stage preserves the current site and historical releases. This
guarantee covers atomic visibility, not full power-loss durability. Linked files
retain the previous release's modification time.

Logs report `linked_files`, `copied_files`, `copied_bytes`, and `compared_bytes`.
A bounded pool processes up to eight files concurrently. Each worker owns its
comparison buffers. Publication waits for all workers before exposing a release.
Independent copies retain per-file synchronization.

Direct byte comparisons still read unchanged output. They avoid copying those
bytes onto replicated storage. The first publication needs independent copies
when no baseline exists, and requires space for both the workspace and release.

Use `mode: hostPath` and set `workspace.hostPath.path` when you want a specific node-local disk
instead of the pod's ephemeral storage. Keep `workDir` underneath `mountPath`; the CLI also rejects
filesystem roots, the source tree, and the release root as explicit `--work-dir` values because the
workspace can be deleted and recreated when it needs a fresh seed.

Keep `builderAdmin.fast` at `false` when queued builds publish the live site. In this repo,
`--fast` is an authoring optimization, not a production-equivalent build mode: it skips
blogroll, mentions, and other expensive work that can affect user-facing output. Enable it
only when the admin service is being used as a preview loop and a separate full build path
still exists for public releases.

Current chart defaults also prefer clean rolling cutover instead of stop-then-start replacement.
Builder-admin keeps one active leader for queue draining, file watching, refresh scheduling, and
 release promotion while standby pods stay ready during rollout handoff.

## Match the Site Theme

Builder admin reads the site's configured palette from the source directory and uses its fallback
mode for the operator UI. For example, this site configuration gives the admin UI the Everforest
Dark colors:

```toml
[markata-go.theme]
palette = "everforest-dark"
```

When your site config uses `palette_light`, `palette_dark`, or `fallback_mode`, builder admin uses
the same fallback palette. Keep the site config and any custom palette files under the mounted
source directory. If the palette cannot be loaded, builder admin remains available with its
default colors. Builder Admin and local Serve share a semantic browser token stylesheet for these
colors and keyboard focus treatment; their page structure and production/local capabilities remain
separate while a shared shell is developed.

## Read Build History Quickly

Jobs and releases are expandable: select a row to see full IDs, timings, logs, and the build
record link without losing your place in the list. Release rows use their creation time as the
primary label; full IDs remain available in the expanded view.

For successful builds, the duration bar compares each build with up to 60 builds completed in the
last 30 days. It marks the mean and recorded maximum. Builder admin needs at least eight matching
builds before it flags unusually fast or slow work, using a median/MAD baseline so a single slow
build does not distort the result.

## Secure Access with ForwardAuth

Builder admin is an operator control plane. Expose it only through its protected Traefik ingress;
direct Service access and `kubectl port-forward` fail closed because they cannot prove that the
request came through the trusted proxy.

```yaml
builderAdmin:
  enabled: true
  auth:
    # CIDR used by the Traefik instances and builder-admin peers that forward this
    # Ingress. 0.0.0.0/0 and ::/0 are rejected.
    trustedProxyCIDRs:
      - 10.42.0.0/24
  ingress:
    enabled: true
    host: builder.example.com
    ingressClassName: traefik
    tls:
      enabled: true
      secretName: builder-example-com-tls
    auth:
      enabled: true
      # Traefik uses this direct Service URL for ForwardAuth. The chart appends
      # /api/v1/forward-auth.
      internalUrl: http://hlab-auth.hlab-auth.svc.cluster.local:8000
      # Optional: browser-reachable auth origin for the operator's own picture.
      # This does not alter hlab-auth login, session, or WebAuthn configuration.
      publicAuthOrigin: https://auth.wayl.one
  networkPolicy:
    enabled: true
    # Verify these against the live Traefik installation before applying.
    traefikNamespace: kube-system
    traefikNamespaceSelector:
      kubernetes.io/metadata.name: kube-system
    traefikPodSelector:
      app.kubernetes.io/name: traefik
      app.kubernetes.io/instance: traefik-kube-system
```

The chart strips every configured identity header before calling the ForwardAuth service. On a
successful decision, Traefik forwards only configured response headers to builder admin. The
ForwardAuth `internalUrl` must use HTTPS or a cluster-local Service URL in the form
`http://<service>.<namespace>.svc.cluster.local:<port>`; Helm rejects other HTTP URLs. Keep
`publicAuthOrigin` on the browser-facing HTTPS auth origin when you use profile pictures.
Set `trustedProxyCIDRs` to only the actual Traefik source CIDRs seen by the pod. A shared Pod CIDR
is permitted only when it is needed for builder-admin peer forwarding and the required selector
NetworkPolicy restricts that CIDR to the configured Traefik pods and builder-admin peers. Do not
use universal, loopback, or link-local CIDRs; the service rejects them.

The default values use hlab-auth's `/api/v1/forward-auth` endpoint, `X-Hlab-*` headers, and the
`__Host-hlab-app-session` cookie. Builder admin itself has no hlab-auth API dependency: any SSO
that supports Traefik ForwardAuth can be used by setting the endpoint path, response cookies, and
identity-header mapping. Builder admin makes no authorization decision from a username, display
name, group, role, or scope; the ForwardAuth service is the access decision point.

### Authentik example

Authentik's proxy outpost supplies durable `X-Authentik-Uid` values. Configure its Traefik
ForwardAuth endpoint, then map those headers in the chart:

```yaml
builderAdmin:
  auth:
    headers:
      userID: X-Authentik-Uid
      username: X-Authentik-Username
      displayName: X-Authentik-Name
      email: X-Authentik-Email
      groups: X-Authentik-Groups
      roles: ""
      scopes: ""
  ingress:
    auth:
      internalUrl: http://ak-outpost.authentik.svc.cluster.local:9000
      forwardAuthPath: /outpost.goauthentik.io/auth/traefik
      responseCookies: []
```

For a non-Helm deployment, use the same names with `builder-admin` flags such as
`--auth-user-id-header=X-Authentik-Uid`. Configure Traefik to strip every configured identity
header before ForwardAuth and forward the same headers from its successful response. Do not map a
mutable username to `userID`; use the SSO's stable subject or UID header.

You can keep this mapping in the site configuration instead of repeating CLI flags:

```toml
[markata-go.builder_admin.auth.headers]
user_id = "X-Authentik-Uid"
username = "X-Authentik-Username"
display_name = "X-Authentik-Name"
email = "X-Authentik-Email"
groups = "X-Authentik-Groups"
roles = ""
scopes = ""
```

Environment variables override this file, for example
`MARKATA_GO_BUILDER_ADMIN_AUTH_HEADERS_USER_ID=X-Authentik-Uid`. Explicit CLI flags have the
highest precedence. Empty optional header values disable those display-only assertions; the durable
`user_id` header is always required.

The operator panel uses the trusted durable user ID and any supplied profile headers only for
identity display; they never grant access. Set `publicAuthOrigin` only when the browser can reach
hlab-auth at that HTTPS origin and you want the panel to request the signed-in operator's own
picture. The URL is derived from the authenticated stable user ID; a `No image` fallback remains
visible when it is unset, unavailable, or cannot load. Do not use this setting to change hlab-auth's
primary public auth origin, cookies/sessions, or primary WebAuthn RP.

### Builder Admin llms.txt for agents

Builder Admin serves a **service-specific** `llms.txt` directly from its
binary. This is separate from the `llms.txt` files generated for published
sites by markata-go. Both entry points contain the same operator guide:

- `GET /llms.txt`: protected by the browser operator/SSO route.
- `GET /api/v1/llms.txt`: protected by the existing trusted bearer API route,
  so agent clients can fetch the guide using the same bearer token as
  `/api/v1/state`.

Both endpoints support `HEAD`, reject writes, return `text/plain`, and set
`Cache-Control: private, no-store`. They do not grant anonymous read access,
mint session cookies, or expose a way to bypass ForwardAuth. The ingress
continues to protect both paths with the same operator-access policy.

**Browser login:** open the configured HTTPS Builder Admin ingress host. Its
SSO/ForwardAuth provider handles login; complete the configured passkey or
session flow and then browse `/llms.txt` or the dashboard. There is no
separate Builder Admin username/password endpoint.

**Agent login:** obtain a personal bearer token approved for the app's operator
scope (for example `gobuilder:access`) from the configured auth provider.
For hlab-auth use the same token as the `hlab` CLI, supplied as
`HLAB_TOKEN` in the agent's secure environment:

```sh
curl -fsS -H "Authorization: Bearer $HLAB_TOKEN" \
  https://builder.example.com/api/v1/llms.txt
```

The hostname is an example; use the actual deployed ingress. Never paste
tokens into the guide, put them in source control, spoof identity headers, or
use port-forward to bypass the trusted proxy. The guide recommends inspecting
state and logs first and getting approval before queueing builds, refreshes, or
rollbacks.

### Agent API

Builder admin exposes a machine API at `/api/v1/` for reading state and logs, inspecting build and
refresh operations, enqueueing builds and refreshes, and rolling back releases. Successful
mutations return HTTP 202 with an operation ID and queued status. Poll the operation detail endpoint
or `/api/v1/state` until it completes. The API records the trusted operator ID on operations.

The API is restricted to requests that pass through ForwardAuth with a trusted `bearer` auth-method
assertion. The default header is `X-Hlab-Auth-Method`; the chart strips any client-supplied copy and
only accepts the value when ForwardAuth returns it. The ForwardAuth route for `/api/v1/` must
require the site's operator-access scope (for example, `gobuilder:access`). Hlab-auth ForwardAuth
validates personal bearer tokens on `/api/` routes against the matched route's required scopes and
returns the identity plus auth-method headers. Deploy a hlab-auth release with bearer ForwardAuth
support before enabling this API. Browser UI requests continue to use passkey sessions and CSRF
validation.

Use the API with the same credential used for `hlab`:

```sh
curl -fsS -H "Authorization: Bearer $HLAB_TOKEN" \
  https://builder.example.com/api/v1/state

curl -fsS -X POST -H "Authorization: Bearer $HLAB_TOKEN" \
  https://builder.example.com/api/v1/builds
```

The API also provides `GET /api/v1/builds/{id}`, `GET /api/v1/refreshes/{id}`,
`GET /api/v1/logs/{filename}`, `POST /api/v1/refresh/{task}`, and
`POST /api/v1/releases/{release-id}/rollback`. It returns JSON errors; queue saturation returns
503 with `Retry-After`. The browser UI remains CSRF-protected, and bearer API authentication does
not grant access to the API unless ForwardAuth validates the token and returns the trusted marker.

When `builderAdmin.enabled` is true, Helm fails unless the protected ingress, TLS secret, ingress
host/class, ForwardAuth URL, trusted proxy CIDRs, and builder-admin ingress NetworkPolicy are all
configured. The service receives its CSRF public origin as exactly `https://<ingress host>` from
the chart. Do not supply that value from request headers or alter hlab-auth's primary RP/origin.
The NetworkPolicy permits the configured `builderAdmin.port` only from the configured Traefik
namespace/pod selectors and other builder-admin pods; confirm the default labels match the
installed controller. The peer-forwarding marker is accepted only from a configured trusted CIDR;
it cannot authenticate a direct request by itself.

Authenticated `GET /` responses set a host-only, secure, HttpOnly, strict-SameSite CSRF cookie
and include its token in each mutation form. Browser mutations must present the matching token and
the exact configured Origin; programmatic clients may use `X-CSRF-Token` instead of a form field.

## What You Can See

The UI shows:

- queued and running builds
- recent successful and failed builds
- the trigger source for each build
- per-build duration and phase timings
- full raw logs
- parsed markata-go performance summary lines
- current live release
- old releases that can be promoted back to live
- refresh task history

## Reading the Operator UI

The default **Jobs** workspace is the primary operational view. It puts running work, queued
work, and completed builds, refreshes, and rollbacks in one list. Use **Enqueue Build** for a
manual build, then inspect its row for status, queue wait, build duration, trigger, and release.

Select **Details** on a build to see phase timings, changed paths, parsed performance lines, and
the raw log. Release pruning runs asynchronously after publication, so a completed build can be live
before old release files finish deleting; the recorded prune timing updates when cleanup completes.
Raw IDs and filesystem paths are detail data rather than dashboard metrics. On narrow screens,
each build and release reflows into a labeled record so operators do not need to read a clipped
desktop table.

## Protected Release Previews

Each retained successful release is available at `https://<site-host>/__preview/<release-id>/`.
The chart protects this path with the builder-admin ForwardAuth middleware, so unpromoted content
is not public. Preview links expire when release retention prunes the underlying release.

The browser tab favicon also reflects live admin state so you can spot activity without keeping the tab focused:

- idle when nothing is running
- queued when work is waiting
- build when a build or rollback is running
- refresh when a refresh task is running
- error when UI polling fails

The workspace tabs show one primary view at a time so build history, refresh runs, and releases do
not visually stack on top of each other during tab switches. Build and release timestamps also pair
the absolute RFC3339 value with a relative age label such as `(5m ago)`.

## Build Triggers

Builder admin can enqueue builds from:

- the UI
- HTTP API calls
- debounced file-watch events
- successful refresh task runs when configured to enqueue a build
- GitHub or Forgejo push webhooks

## Build on Git Push

Enable one builder-admin service for each environment. Each service owns one checked-out branch,
source directory, site directory, and webhook secret. The branch defaults to `main`.

```toml
# markata-go.toml
[markata-go.builder_admin.webhook]
enabled = true
branch = "main"
```

Set the HMAC secret outside version control. Environment variables override the config file, and
explicit flags override both:

```bash
export MARKATA_GO_BUILDER_ADMIN_WEBHOOK_SECRET='replace-with-a-random-secret'
markata-go builder-admin --webhook-enabled --webhook-branch main
```

The endpoint is `https://<builder-admin-host>/webhook`. In GitHub, create a repository webhook
with that payload URL, set the same secret, select JSON payloads, and subscribe to **Push** events.
Forgejo uses the same URL and secret; select its push event. Builder admin accepts GitHub's
`X-Hub-Signature-256` and Forgejo's `X-Gitea-Signature` headers. Invalid signatures, events other
than push, and pushes to another branch are ignored or rejected without running a build.

For an accepted matching push, builder admin runs `git -C <source-dir> pull --ff-only origin <branch>`. It queues
a build only when that changes the checkout. Ensure the deployed source directory is a clean clone
with an authenticated `origin` remote and that its checked-out branch tracks the configured branch.
The normal source watcher remains active during the pull so a local author edit is never ignored.
That safety measure can produce a follow-up file-watch build after a push.

### Helm setup

The chart exposes `/webhook` without the operator ForwardAuth middleware; HMAC validation protects
that endpoint. Create a Kubernetes Secret, then reference it in the environment's Helm values:

```yaml
builderAdmin:
  webhook:
    enabled: true
    branch: main
    existingSecretName: builder-webhook
    existingSecretKey: webhook-secret
```

Create separate releases for `main`, `dev`, and `qa` when they publish different domains. Give each
release a distinct source checkout, site storage, builder-admin ingress host, and webhook secret.
Point each repository webhook at its matching builder-admin host. Do not share a production webhook
secret with preview environments.

If GitHub or Forgejo reports a successful delivery but no build appears, verify the delivery branch,
the service's configured branch, the HMAC secret, and that `git -C <source-dir> pull --ff-only origin <branch>`
succeeds in the builder container. A delivery for an unchanged commit intentionally does not queue
a build.

## Rollback

Rollback in the first version promotes a previous rendered release by repointing the `current` symlink.

This is fast and keeps the release model simple, but it does not restore the historical source tree.

## Scheduled Refresh Tasks

Refresh tasks are external commands that run on an interval.

Examples:

- `markata-go reader update`
- `markata-go blogroll update --force`
- a custom remote-asset fetch command

Use them to keep reader/blogroll data or other remote caches fresh without slowing down every normal content build.

### Retaining a warm workspace

Successful publication retains the completed workspace on both same-filesystem and
cross-filesystem layouts. Builder Admin reuses it only when its successful-release
marker matches current. Logs report `reusing build work from current release`.
Failed or interrupted builds and rollbacks trigger a fresh independent seed. A
persistent workspace survives pod restarts. An `emptyDir` survives container
restarts but disappears with the pod. Incremental publication still compares
unchanged bytes, so storage and host load affect elapsed time.

### Background release cleanup

Builder Admin first removes obsolete releases from the published release list with an atomic rename, then deletes their files in the background. Large cleanup jobs do not hold the lock needed for promotion or rollback. Internal `.pruning-` directories are not rollback targets; interrupted cleanup is retried on the next prune run. Published build success is separate from cleanup errors, which appear in service logs.

### Successful build but nginx returns 403

If nginx reports `current/index.html is forbidden (13: Permission denied)`, inspect the completed release root permissions. Builder-admin publishes release roots as `0755` so nginx can traverse them. Older cross-filesystem publication could retain a private `0700` staging directory, and subsequent seeded builds could inherit it. Upgrade the engine and correct affected retained release roots to `0755`; child file permissions do not need a recursive change for this defect.

### Warm publication digest cache

Builder Admin stores a private SHA-256 manifest beside its mutable workspace.
On Linux, unchanged file identity and change timestamps allow cached digest reuse.
Files with recent or changed timestamps require content hashing. Missing or corrupt
cache records fall back to hashing. Other platforms always hash source files.

The manifest binds the workspace and previous release. Published releases must
stay immutable. Do not edit retained releases in place. Delete the workspace's
`.publication.json` sibling to force content verification on the next publication.
Manifest write failures leave publication usable but disable cache reuse. Logs
report `cached_source_hashes`, `cached_release_hashes`, and `manifest_saved`.
The manifest contains file identities and digests and never appears in a release.


## Revision provenance and new-build status

Builder Admin shows the running **Builder executable revision** beside the
revision used by the latest successful build, as well as the current site Git
HEAD. Its top status banner highlights **New Builder revision has not built
yet** when the running commit differs from every retained successful build
(and there is reliable previous-build provenance). A queued or running build
shows a pending state. Once a successful build using that Builder revision is
recorded the warning clears.

New build records persist `builder_version`, `builder_commit`, and
`source_commit`. These are visible in Jobs details, build detail pages and
the authenticated `/api/v1/state` response under `revision`. Retained
releases also show recorded revisions where available. Older build records may
lack these fields; the current release's diagnostics artifact provides a
best-effort fallback. Unknown revision identity is explicitly reported as
**unknown**, not mistaken for an unbuilt revision.

Source revision comparison is Git HEAD based. Uncommitted edits can require a
rebuild even if HEAD matches; this status does not claim to fingerprint dirty
source files. The displayed Builder revision is taken from the running binary's
build metadata, not the site's Git HEAD or a mutable image tag.
