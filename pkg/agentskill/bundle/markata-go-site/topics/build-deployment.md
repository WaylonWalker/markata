# Build And Deployment

Use this topic when the task involves local preview, CI, publishing, or hosting strategy.

## Local Iteration

- `markata-go serve --fast` for active editing
- `markata-go build` for full output validation
- `markata-go build --clean` only when you need to rule out stale output
- a clean build publishes its output after the new build succeeds; a failed build leaves the previous output available
- `markata-go build -o dist` when CI or previews need an isolated artifact directory
- `markata-go buildlab run --fixture /path/to/site` when you need clean, incremental, and deterministic build evidence

## Production Build Basics

- set the correct `url` before production builds
- treat `output_dir` as the deploy artifact root
- prefer clean builds for deployment validation
- validate config before deploy if the workflow can afford it
- if the deploy target runs with limited or no internet egress, prefetch self-hosted CDN assets with `markata-go assets download` before shipping the repo or build input

Examples:

```bash
markata-go config validate
markata-go build --clean
MARKATA_GO_URL=https://example.com markata-go build
```

Use Build Lab to compare a baseline and candidate binary without changing the
ordinary build command:

```bash
markata-go buildlab run \
  --fixture /path/to/site \
  --baseline /tmp/markata-go-baseline \
  --candidate /tmp/markata-go-candidate
```

Build Lab runs child builds in isolated workspaces and checks clean output,
incremental output, and deterministic rebuilds. See the Build Lab guide for
scenario mutations and output classes.

## Recommended CI Shape

For most hosts, the safe build flow is:

1. checkout repo
2. install Go
3. install or build `markata-go`
4. run `markata-go config validate`
5. run `markata-go build --clean`
6. publish the build artifact from `output/` or the configured output dir

A successful full build also includes `.markata/diagnostics.json`
below the chosen output directory. Preserve this file when publishing a release
if build disposition and content-quality data will be inspected after deployment.
It is sanitized and does not contain raw content, secrets, or absolute paths. Dry
runs, fast or incremental development-server builds, failed builds, or incomplete
builds do not publish a new diagnostics artifact. A normal `serve` build can
publish it after a successful full lifecycle.

Minimal GitHub Actions shape:

```yaml
steps:
  - uses: actions/checkout@v4
  - uses: actions/setup-go@v5
    with:
      go-version: '1.26'
  - run: go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
  - run: markata-go config validate
  - run: markata-go build --clean
    env:
      MARKATA_GO_URL: https://example.com
```

## Hosting Targets That Fit Well

- GitHub Pages
- Netlify
- Vercel
- Cloudflare Pages
- AWS S3 or another static bucket
- self-hosted nginx, Caddy, or Docker-based static hosting

## Provider Patterns

### GitHub Pages

- build in GitHub Actions
- upload `output/` (or the configured output dir) as the Pages artifact
- set `MARKATA_GO_URL` to the final GitHub Pages or custom domain URL

### Netlify

- build command usually installs markata-go then runs `markata-go build --clean`
- publish directory is `output/` unless `output_dir` is configured
- set production and preview `MARKATA_GO_URL` separately if needed

### Vercel

- use a custom build command because this is a static Go-built site, not a framework preset
- set `outputDirectory` to `output` unless `output_dir` is configured
- verify preview and production URLs separately

### Cloudflare Pages

- framework preset is usually none
- build output directory is `output/` unless `output_dir` is configured
- `_headers` and `_redirects` can live under `static/` so they copy into output

### Self-hosted

- deploy the built `output/` directory (or the configured output dir) behind nginx, Caddy, S3, or another static file server
- ensure the server preserves nested `index.html` routing and static asset paths
- for nginx-native redirects, include the generated `redirects.conf` from the configured output directory and reload nginx when that file changes; see `docs/guides/nginx-redirects.md`

## Guidance

- Prefer simple static hosting first.
- Keep deployment changes focused on reproducible builds and correct output paths.
- If the repo already has CI, extend the existing workflow instead of replacing it.
- If a deploy bug is path-related, inspect `output_dir`, `url`, asset paths, and feed URLs before changing templates.
- If previews and production use different domains, inject `MARKATA_GO_URL` per environment instead of hardcoding one value.
- if the runtime build environment is offline, make sure `.markata/assets-cache` or another configured asset cache is already populated before relying on self-hosted CDN assets.
- for Helm or ArgoCD source-archive deployments, prefer environment-specific `MARKATA_GO_*` overrides such as `MARKATA_GO_URL` instead of editing the repo just to change hostnames
- for Helm or ArgoCD feedback loops, inspect the chart's lock and debounce knobs (`build.lock.pollIntervalSeconds`, `search.waitForSource.pollIntervalSeconds`, `search.watchDebounce`) before chasing deeper build bugs; shorter values make manual rebuilds and search restarts feel noticeably snappier
- for Kubernetes hostPath authoring deployments, prefer the long-lived `builder-admin` service over one-shot build Jobs when the goal is fast interactive rebuilds, release history, rollback, and scheduled remote refreshes
- for agent-triggered builder-admin operations, use its `/api/v1/` JSON API with a bearer credential only after ForwardAuth documents and validates bearer access for that route; a browser-session scope alone does not make UI forms callable from a headless agent
- for builder-admin operational orientation, read authenticated `/llms.txt` from a signed-in browser or `/api/v1/llms.txt` using an authorized bearer token. The embedded guide explains real SSO login, read-only state/logs, and approval before builds/refreshes/rollbacks; it is distinct from site-generated `llms.txt`
- for Kubernetes hostPath deployments, confirm the mounted source path and served site root are the real node paths, and remember the served site root may contain release directories plus a `current` symlink rather than a flat output tree
- builder-admin is an operator surface: expose it only through its dedicated protected Traefik/hlab-auth ingress, configure `builderAdmin.auth.trustedProxyCIDRs` for the actual Traefik sources and required builder-admin peers, and do not use Service access or `kubectl port-forward` as an authentication bypass
- when enabling builder-admin, configure its TLS host, HTTPS ForwardAuth URL, and explicit ingress NetworkPolicy selectors for the live Traefik pods; a shared Pod CIDR is acceptable only with those selectors for peer forwarding, while universal, loopback, and link-local CIDRs are rejected; the chart derives the exact `https://<host>` CSRF origin, so do not derive an origin from forwarded request headers or change hlab-auth's primary RP/origin
- Validate that feed URLs, social URLs, and asset URLs use the expected domain after build.

## Markata-Go-Specific Checks Before Shipping

- `url` matches the actual deployment domain
- `output_dir` matches the host publish directory
- static assets copied from `static/` are present in output
- feeds like `rss.xml` or `atom.xml` contain the correct absolute URLs
- homepage and archive slugs resolve to the intended directories
- any `_headers`, `_redirects`, or `CNAME` files under `static/` are present in output when the host expects them

## Preview Deployments

Good preview strategy:

- keep the same build command
- inject a preview URL only when absolute links must be correct
- if the provider gives automatic preview URLs, prefer provider env configuration over editing config files

## Inspect

- workflow files under `.github/workflows/` or other CI directories
- deployment docs in the repo
- output path settings in config
- asset or CDN self-hosting settings when debugging missing static files
- any `static/CNAME`, `static/_headers`, or `static/_redirects` files
# Git Push Rebuilds

For deployments using builder-admin, GitHub and Forgejo push webhooks can pull and rebuild a
single configured branch. Read the site's `[markata-go.builder_admin.webhook]` configuration before
changing deployment behavior. Each production, development, QA, or preview environment should use
an independent builder-admin deployment, source checkout, release root, ingress host, and webhook
secret. The webhook endpoint is `/webhook`; configure its secret through
`MARKATA_GO_BUILDER_ADMIN_WEBHOOK_SECRET` or a Kubernetes Secret, never commit it to the site
repository. See the Builder Admin deployment guide for Git-provider and Helm setup.

## Builder workspace performance

When builder-admin preparation dominates warm builds, inspect the workspace mount and prepare/promote timings. A retained node-local workspace skips reseeding only when its successful-release marker matches current. Failed builds and rollbacks reseed safely. Keep published releases on durable storage; a fast workspace does not remove the cross-filesystem promotion copy. Use a persistent hostPath workspace to retain warm output across pod replacement, and verify the physical disk behind its path.

When promotion latency overlaps release pruning, check the deployed engine version. Builder-admin detaches obsolete releases under the publication lock and deletes their trees afterward; internal `.pruning-` and `.staging-` directories are not rollback targets. Interrupted pruning is retried on the next cleanup.

A successful build record alone does not prove the static server can read the release. Check readiness and public endpoints. For nginx `current/index.html` permission errors after staged publication, inspect the release root: builder-admin normalizes it to `0755`; older versions could retain `0700` from a temporary directory. Correct only affected root permissions and upgrade the engine.

## Builder Admin Workspace Reuse

Builder Admin retains an independent workspace after successful publication.
Release staging links unchanged files from the previous immutable release and
copies changed files. Never hard-link a mutable workspace to a published release.
Warm logs report `reusing build work from current release`. Publication logs show
linked files, copied files, copied bytes, and compared bytes. Failed builds and
rollbacks force independent reseeding. Compare complete build timings, including
preparation and publication. Direct comparison still reads unchanged output.

Publication uses up to eight concurrent file operations to overlap storage
latency. It still compares exact bytes and synchronizes independent copies.

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
