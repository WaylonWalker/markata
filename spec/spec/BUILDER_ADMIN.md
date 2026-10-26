# Builder Admin Specification

This document specifies the long-lived builder admin service for markata-go Kubernetes deployments.

## Goals

- Provide a warm, long-lived build worker so authoring builds avoid per-Job startup cost.
- Expose an operator-facing HTTP UI and API for builds, releases, logs, refresh tasks, and rollback.
- Preserve the existing release model based on `releases/<id>/` plus a `current` symlink.
- Keep remote-content refresh work out of normal content builds unless explicitly configured otherwise.

## Scope

The builder admin service is intended for self-hosted and Kubernetes workflows, especially hostPath-backed authoring deployments.

The first required capabilities are:

- serialized build queue
- manual HTTP-triggered builds
- file-watch triggered builds that enqueue through the same queue
- build history with full raw logs
- release history with current/live indicator
- promote-previous-release rollback
- scheduled refresh tasks for reader/blogroll/other external data commands
- operator UI that shows running, queued, successful, and failed work

## Runtime Model

The builder admin service MUST run as a long-lived HTTP process.

It MUST mount the same site-authoring paths as the existing build workflow:

- source tree
- rendered site root
- optional dedicated cache volume

The service MUST process queued work one item at a time for a given site.

Triggers MUST enqueue work rather than executing builds directly.

Kubernetes rollouts MUST support clean cutover without dropping the currently active builder-admin pod before the replacement pod is ready.

To support that requirement, the service MUST tolerate an active/standby deployment shape where multiple pods may be live briefly but only one pod is allowed to perform mutating work.

Required trigger sources:

- manual UI action
- manual HTTP API call
- file watch
- scheduled refresh completion when configured to enqueue a build
- rollback action
- signed GitHub or Forgejo push webhook after a successful source update
- Helm startup/leadership acquisition when the running Builder commit differs from the latest successful build; this MUST go through the serialized queue rather than a direct build

## Leadership And Handoff

When more than one builder-admin pod is running for the same site, exactly one pod MUST hold leadership for mutating work.

The leader is responsible for:

- draining the serialized work queue
- running file watching
- running scheduled refresh tasks
- executing builds and rollbacks
- persisting queue/running/history state

Standby pods MUST:

- serve the read-only HTTP UI and API state
- remain ready so rolling updates can keep the old leader serving while the new pod starts
- refuse or forward mutating requests unless they become leader

If a standby pod receives a mutating HTTP request while another pod is leader, it SHOULD forward that request to the active leader so operator actions do not fail during rollout handoff.

On leadership acquisition after a restart or rollout, persisted queued work MUST be replayed. A previously running in-flight operation MAY be marked interrupted instead of resumed.

## Build Workflow

Successful builds MUST preserve the existing atomic release publication model:

1. prepare cache symlinks when a dedicated cache mount is configured
2. seed a stable work directory from the current release when one exists
3. run `markata-go build` into the work directory
4. stage an independent release under `releases/`, then atomically rename it to `releases/<release-id>/`
5. atomically repoint `current` to the new release
6. publish build success and queue release pruning according to retention policy

Release pruning MUST NOT block successful publication or the next queued build.
Pruning MUST remain serialized with promotion and rollback, and MUST re-check
the current release before deleting each expired release.

Every retained successful release MUST be available through the protected preview path
`/__preview/<release-id>/`. Preview routing MUST use the same ForwardAuth policy as builder-admin
and MUST rewrite root-absolute HTML asset URLs to the selected release prefix so previews do not
mix current-release assets with historical content.

The service MUST record phase timings for at least:

- queue wait
- prepare
- build
- promote
- prune (recorded when asynchronous cleanup completes)
- total

The service MUST store the full raw build log and a parsed performance summary that includes any `Duration:` and `Hotspots:` lines emitted by markata-go.

## File Watching

When file watching is enabled, the service MUST watch the configured source roots recursively.

Watch events MUST be debounced and coalesced into a single queued build request.

The recorded build trigger MUST include:

- trigger type `file-watch`
- the set of changed paths captured during the debounce window

The watcher SHOULD ignore internal cache and admin-state paths.

## Git Push Webhooks

Builder-admin MAY expose `POST /webhook` when webhook support is enabled. The endpoint MUST
accept GitHub and Forgejo push deliveries, validate their HMAC-SHA256 signature in constant time,
and reject deliveries with missing or invalid signatures. It MUST not require the operator
ForwardAuth identity because Git providers cannot supply it.

Each builder-admin instance MUST have one configured branch, defaulting to `main`. It MUST ignore
push deliveries for every other branch. Deployments for production, development, QA, or previews
MUST use separate builder-admin instances with independent source checkouts, site roots, webhook
secrets, and branch configuration.

For an accepted matching push, the active leader MUST run `git -C <source-dir> pull --ff-only origin <branch>`.
It MUST enqueue a build only when the checked-out commit changes. Pull failures and non-fast-forward
updates MUST not enqueue a build. The corresponding build record MUST use trigger type `webhook`
and retain provider, repository, branch, commit, and delivery metadata when supplied.

Filesystem changes remain visible to the source watcher while a webhook pull runs. This prioritizes
local author edits over deduplicating watcher events: a concurrent or delayed filesystem change MAY
enqueue a follow-up file-watch build.

Webhook configuration MUST be available in `[markata-go.builder_admin.webhook]`, through
`MARKATA_GO_BUILDER_ADMIN_WEBHOOK_*` environment variables, and through builder-admin flags.
Environment values MUST override file configuration; explicitly supplied flags MUST override both.
An enabled webhook MUST have a non-empty secret.

## Build History

Each build record MUST include:

- unique build id
- operation kind: `build`, `refresh`, or `rollback`
- status: `queued`, `running`, `success`, `failed`, `cancelled`
- trigger type
- trigger detail text
- changed paths when available
- enqueue, start, and finish timestamps
- per-phase timings
- total duration
- raw log path
- parsed performance summary
- produced release id, when applicable
- whether the result became live

The UI MUST show current queue state, running build state, and the current live release.

### Operator UI Information Hierarchy

The default operator view MUST prioritize the work that needs attention over implementation
metadata. It MUST present, in this order:

1. the currently running operation and its elapsed time, or a clear idle state;
2. the number of queued operations and the oldest queue wait;
3. the latest build result and current live release; and
4. the build history list.

The default view MUST use one primary **Jobs** workspace. It MUST combine queued, running, and
completed builds, refreshes, and rollbacks in operational priority order: running work first,
then queued work, then completed work ordered by completion time. A separate releases view MAY
remain for promotion actions. Queue, running-operation, and release information MAY be presented
as compact supporting regions, but MUST NOT require an operator to scan separate queue and
history tables to understand all work.

The build-history list MUST make status, trigger, queue wait, build duration, completion time,
and resulting release scannable. Per-phase timings, changed paths, performance-summary lines,
raw filesystem paths, and raw logs MUST be available through a build details disclosure or
detail view, rather than occupying the default table columns.

Raw filesystem paths MUST NOT be a primary dashboard metric. If shown, they MUST be secondary
detail information and use intentional wrapping or copyable overflow treatment.

Raw operation IDs, release IDs, user IDs, email addresses, groups, roles, and scopes MUST NOT
appear in the default operator view. The header MAY show the signed-in operator's display name
and avatar. Diagnostic identity assertions and raw IDs MAY appear only in an explicitly opened
detail or diagnostic disclosure.

Decorative backgrounds MUST NOT reduce text contrast or run through operational data. The UI
MUST retain readable, useful layouts at narrow viewport widths; repeated operational records
MUST reflow into labeled rows or otherwise provide an explicit, usable narrow-screen treatment.

The workspace view MUST display one primary panel at a time when switching between builds, refresh runs, and releases. Tabs MUST NOT leave multiple primary panels visually stacked on top of each other.

Build, refresh, running, and release state labels SHOULD use distinct visual status treatment so success, failure, queued, running, and live states are scannable without reading raw text.

Build and release timestamps shown in the UI SHOULD pair an absolute timestamp with a relative age label such as `(5m ago)`.

The UI MUST expose an at-a-glance browser-tab indicator by updating the favicon to reflect the current admin state:

- idle when nothing is running and the queue is empty
- queued when work is waiting but not yet running
- build when a build or rollback is running
- refresh when a refresh task is running
- error when live polling fails

## Releases And Rollback

The service MUST discover releases from the site root `releases/` directory and the `current` symlink.

Rollback in the first version is defined as:

- selecting a previously successful rendered release directory
- atomically repointing `current` to that release
- recording a rollback operation in history

The UI MUST clearly indicate that rollback promotes a prior rendered release rather than restoring the historical source tree.

## Refresh Tasks

The builder admin service MUST support configured scheduled refresh commands.

Each refresh task MUST define:

- stable task name
- command argv
- interval duration
- whether a successful run enqueues a build

The first version MAY use fixed interval durations instead of cron expressions.

Refresh runs MUST have their own history with:

- task name
- status
- duration
- raw log path
- optional follow-up build id when a build was enqueued

## Persistence

The service MUST persist operator state on disk so restarts do not lose build history.

Required persisted data:

- build records
- refresh records
- release metadata derived from disk and linked build ids when known
- full raw logs

The first version MAY use a JSON state file plus log files instead of a relational database.

## HTTP UI And API

The service MUST expose an HTTP admin interface.

### Operator Authentication

Every builder-admin UI, API, and log route except the Kubernetes `/health` probe MUST require a
trusted ForwardAuth identity. The first deployment uses Traefik `ForwardAuth`, but the service
MUST NOT require a specific SSO product. The chain MUST first remove client-supplied configured
identity headers, then forward the authentication decision and configured identity assertions.
Response cookies required by the selected SSO MAY be copied with `addAuthCookiesToResponse`.
The service MUST reject a request unless its remote address belongs to an explicitly configured
trusted-proxy CIDR and it includes exactly one non-empty configured durable-user-ID header.
Universal, loopback, and link-local trusted-proxy CIDRs MUST be rejected. The leader-forwarding
marker MUST never substitute for source provenance: it is valid only from a configured trusted CIDR.

The operator UI MUST use only the configured durable-user-ID header as the identity key. Username,
display name, email, groups, roles, and scopes are display-only assertions; they MUST NOT grant or
widen privileges. The defaults MUST preserve hlab-auth's `X-Hlab-*` contract. CLI and Helm
configuration MUST permit replacement header names, including empty optional headers. It MAY show
the operator's own profile picture only when an explicit HTTPS public auth-origin configuration is
set; it MUST construct the escaped URL from that origin and the stable user ID as
`/users/{user_id}/picture`, with a useful no-image fallback.

The identity-header mapping MUST be configurable in `[markata-go.builder_admin.auth.headers]` and
through `MARKATA_GO_BUILDER_ADMIN_AUTH_HEADERS_*` environment variables. Environment values MUST
override file configuration, and explicitly supplied CLI header flags MUST take precedence over
both. The durable `user_id` mapping MUST remain non-empty after resolution; empty optional mappings
disable their display-only assertions.

Direct Service access and `kubectl port-forward` cannot establish proxy provenance and MUST fail
closed. They are not supported operator-access paths for authenticated deployments.

### Bearer Machine API

The service MUST expose a versioned JSON API under `/api/v1/` for non-interactive operator clients.
Every machine API request MUST pass through the configured trusted ForwardAuth proxy, include the
durable operator ID, and include exactly one configured auth-method header whose value is `bearer`.
The auth-method header MUST be removed from client requests before ForwardAuth and copied to the
service only from a successful ForwardAuth response. The service MUST NOT infer bearer validation
from a client-supplied `Authorization` header or from display-only scope headers.

The API MUST expose state, build and refresh detail, and operation logs; it MUST enqueue builds,
refresh tasks, and release rollbacks through the existing serialized queue and leader handoff. The
ForwardAuth route for `/api/v1/` MUST require the app's configured operator-access scope. Mutations
MUST record the authenticated durable user ID on the queued operation and resulting history record.
Successful enqueue responses MUST use HTTP 202 and return the operation ID and queued status. API
errors MUST be JSON with a stable error code and message. Queue-full responses MUST include a
retryable status and `Retry-After` header.

### Operator UI Theme

The operator UI MUST derive its default color scheme from the site theme configured for the
builder's source directory. It MUST select the configured fallback-mode palette, including
explicit light or dark palette overrides, and map palette semantic and component colors to its
background, text, border, action, code, and status UI colors. A missing or invalid site palette
MUST fall back to the built-in builder-admin colors without preventing the service from starting.

### Build Workspace

Job and release rows MUST disclose their full operational details when activated anywhere on the
row and MUST remain keyboard accessible. Release summaries MUST use a human-readable creation
time and keep full release and build identifiers in the expanded detail area. Build records MUST
provide a prominent link to their full detail page.

For completed successful builds, the workspace MUST compare build duration with up to 60 builds
finished in the preceding 30 days. It MUST not classify a duration when fewer than eight samples
are available. The visual comparison MUST show mean and recorded maximum markers and use a robust
median/MAD band to apply restrained fast, slow, and extreme-slow emphasis.

### Browser Mutation Protection

An authenticated `GET /` MUST mint a cryptographically random double-submit CSRF token. It MUST
set the token in a host-only `__Host-` cookie with `Secure`, `HttpOnly`, `SameSite=Strict`, and
`Path=/`, and embed the same token in every browser mutation form. Every browser mutation MUST execute
only after the active leader validates a cookie token against either the submitted form token or
`X-CSRF-Token` using a constant-time comparison. It MUST also require an exact `Origin` match to
the configured HTTPS public origin and reject a present `Sec-Fetch-Site` other than `same-origin`.
Standbys MUST forward mutations before CSRF validation so that the active leader performs the
validation and active/standby handoff remains correct. The forwarding marker MUST be stripped at
the public ingress and accepted only from a configured trusted CIDR for builder-admin peer traffic
protected by the required ingress NetworkPolicy.

Required UI views:

- dashboard summary
- build history list
- build detail/log view
- release list with current/live indicator
- refresh task list and refresh history

Required actions:

- enqueue manual build
- trigger refresh task immediately
- promote a prior release to live

The service SHOULD also expose JSON endpoints for the same core operations.

The browser UI SHOULD derive favicon state from the same `/api/state` polling payload used to update the dashboard so the tab indicator stays consistent with the visible live state.

## Helm Integration

The Helm chart MUST support enabling the builder admin service independently of the scheduled build CronJob.

Required chart configuration includes:

- service enable/disable
- host/port
- file-watch enable/disable and debounce
- release retention
- build history retention
- refresh task definitions
- rollout strategy settings for clean active/standby cutover
- protected builder-admin ingress using a configurable ForwardAuth path, response cookies,
  identity-header mapping, optional public auth origin for self profile pictures, and trusted-proxy CIDRs
- a ForwardAuth internal URL that uses HTTPS or a cluster-local Service endpoint
  (`http://<service>.<namespace>.svc.cluster.local:<port>`); the public auth origin remains HTTPS
  when used for browser login or profile pictures
- an exact public origin derived as `https://<builder-admin ingress host>` and passed to the
  service for CSRF origin validation, never inferred from request headers
- an enabled ingress NetworkPolicy that allows the configured builder-admin port only from the
  configured Traefik namespace/pod selectors and builder-admin peers

The first protected deployment target MUST use the dedicated authenticated ingress; `kubectl
port-forward` is intentionally not an operator-access path because it bypasses proxy provenance.

When webhook support is enabled, the chart MUST render a separate exact `/webhook` ingress path
without the operator ForwardAuth middleware. Its HMAC secret MUST be supplied from a Kubernetes
Secret rather than a command-line argument.

The default rendered-release retention MUST keep at least 25 releases, including the current live release, so operators have more than ten rollback targets by default.

## Retained build workspace reuse

A workspace retained after successful staged publication MUST be marked with the exact successfully promoted release ID. The marker MUST live beside the workspace, outside generated output. Preparation MAY reuse the retained workspace only when its marker matches the current live release, and MUST consume the marker before mutation. Missing or mismatched proof, a failed or interrupted build, and rollback MUST force a fresh seed from current. Same-filesystem publication MUST also retain an independent workspace. A marker write failure after successful publication MUST be logged without changing the successful build result. Warm preparation MUST log when it reuses output instead of copying current.

## Prune and promotion concurrency

Pruning MUST recheck the live release while holding the release lock, atomically rename an obsolete release to an internal `.pruning-` directory on the release filesystem, then release the lock before recursively deleting that directory. Promotion and rollback MUST NOT wait for recursive deletion. Release discovery and rollback MUST exclude internal dot-prefixed directories. Interrupted deletions MUST be retried by later prune runs. Prune failures MUST preserve current and remain observable without changing a successfully published build to failure.

## Published release root permissions

Completed release roots MUST use permission mode `0755` before publication so an unprivileged static server can traverse them. This applies to same-filesystem rename and cross-filesystem staged copy, including workspaces seeded from releases created with private root permissions. Staging directories MUST remain private while copying and become `0755` only after their copy completes. Child file permissions MUST remain unchanged. Permission failures MUST prevent publication and preserve current.

## Span presentation foundation

The span view model MUST sort completed spans by start offset, breaking ties by
ID, and expose ID, parent ID, name, stage, plugin, status, attributes, depth,
and start/duration in integer milliseconds. Negative times MUST display as
zero. Attribute maps MUST be copied so presentation changes cannot mutate
recorded input. Missing parents MUST produce root depth; cycles MUST terminate.

The highlighted completion chain MUST select the latest-finishing leaf, prefer
the longer leaf on equal end times, and mark its ancestors until the root, a
missing parent, or a repeated ID. A graph with no leaves MUST have no highlighted
chain. This heuristic MUST NOT be described as a dependency-aware critical path.

The initial model is internal groundwork. It does not require benchmark artifact
retention, API exposure, or an operator waterfall; those remain follow-up work.

## Retained Workspace And Incremental Publication

Builder Admin MUST retain the mutable workspace after successful publication on
both same-filesystem and cross-filesystem layouts. A successful-release marker
MUST match current before reuse. Preparation MUST consume the marker before
mutation. Failure, interruption, missing markers, and rollback MUST force reseeding.

Publication MUST enumerate the completed workspace into a fresh hidden staging
tree. It MUST omit paths absent from the workspace. An unchanged regular file
MAY share an inode with the previous immutable release only after content-digest and mode
comparison. A workspace file MUST NOT share an inode with a published release.
Changed files and unsupported hard-link operations MUST use independent copies.
Symlinks MUST be copied as links, never traversed during workspace enumeration.
Baseline lookup and hard linking MUST remain confined to the release root.

The staging root MUST use mode 0755. Copied files MUST retain their permissions
and modification times. Linked files retain the previous release's modification
time. Publication MUST finish the staged tree before exposing its final name or
switching current. A failure MUST preserve current and historical releases.

Publication logs MUST report linked and copied file counts, copied bytes, and
compared bytes. Timings alone MUST NOT substitute for these work counts.
Publication MAY retain a private SHA-256 manifest beside the workspace. The
manifest MUST bind its checksum, workspace path, and baseline release. On Linux,
a cached source digest requires unchanged device, inode, size, mode, modification
time, and change time. Both timestamps MUST precede the manifest commit by more
than one second. A backwards clock MUST disable this shortcut. Other platforms
MUST read source content. Immutable baseline digests require matching device,
inode, size, mode, and modification time; hard-link change times are excluded.
Missing, corrupt, or mismatched records MUST fall back to content hashing.
Manifest write failure MUST NOT fail publication. The manifest MUST be private,
atomically replaced, bounded to 32 MiB when loading, and omitted from releases.
Logs MUST report source and baseline digest cache hits and manifest persistence.
Retained releases MUST remain immutable; privileged out-of-band modifications
that preserve all metadata are outside this cache contract.
This feature MUST NOT claim power-loss durability beyond the existing staged
publication contract.

Incremental publication MUST use a bounded pool with at most eight concurrent
file operations. Each worker MUST own its comparison buffers and counters.
Parent directories MUST exist before file processing. Publication MUST wait for
all workers before applying directory metadata or exposing the release. A worker
failure MUST prevent publication, and all workers MUST finish before staging
cleanup. Independent copies MUST retain the existing per-file synchronization.
