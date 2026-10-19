# Markata Notes Helm Chart

This chart deploys a reusable markata-go notes workload that:

- pulls a source archive from object storage
- optionally decrypts it with `MARKATA_GO_ENCRYPTION_KEY_DEFAULT`
- renders the site in-cluster on a schedule or manual trigger
- serves the rendered site with nginx
- serves bleve search from the same host under `/api/search`

## Notes

- `work_notes/justfile` now uploads `source.tar.gz` from `git archive HEAD`.
- Set `MARKATA_GO_SOURCE_ARCHIVE_ENCRYPT=true` when publishing if you want the uploaded archive encrypted with `MARKATA_GO_ENCRYPTION_KEY_DEFAULT`.
- The search pod runs `markata-go search-server --mode watch-content --host 0.0.0.0` so bleve stays in sync when the source PVC changes.
- Source-backed search modes can opt into `search.localCache.enabled: true` to mount pod-local `emptyDir` storage at `/data/cache`. Use this when the source tree contains `.markata` (or another cache path) symlinked to `/data/cache/build`; leave it disabled by default. It does not change the source or index volumes and is not mounted for `read-only-index` mode.
- When `search.waitForSource.enabled` is true, source-backed search and builder-admin pods wait for the source-ready marker before starting. This also works when `build.enabled` is false and an external bootstrap job populates the source volume.
- Site and search probes now check every 5s instead of every 10s, which trims rollout and restart latency without making the health checks brittle.
- Runtime pods use a dedicated ServiceAccount with `automountServiceAccountToken: false`, disable service links, and apply `RuntimeDefault` seccomp with stricter container security settings where they are low risk.
- A NetworkPolicy now limits runtime pod egress to cluster DNS, so the site and search pods cannot freely call other cluster services by default.
- The build CronJob uses a source PVC lock by default so overlapping manual jobs cannot update the shared source and site PVCs at the same time.
- Default lock polling and search debounce values are intentionally short (`2s`) so the feedback loop feels tighter when a developer triggers back-to-back builds or restarts the search pod.
- The build CronJob now overrides `[markata-go.mermaid].mode` to `client` by default, which avoids Chromium hangs seen in-cluster; set `build.mermaid.mode` back to `chromium` or `""` if you want to rely on the source repo config instead.
- If your new namespace does not already have `aws-default`, enable `aws.sealedSecret` and provide sealed credentials in values.
- By default the chart reuses a shared `markata-go-encryption` Secret name for `MARKATA_GO_ENCRYPTION_KEY_DEFAULT`; override `markataEncryption.secretName` only if your established secret uses a different name.
- If your notes contain private content, provide that Secret externally in the namespace or enable `markataEncryption.sealedSecret`.
- The nginx site adds `Access-Control-Allow-Origin: *` to all public file responses by default. Set `site.cors.enabled` to `false` to disable this header. Protected release previews never receive it.
- Public HTML is `no-cache` by default. Set `site.cacheHtmlAtEdge: true` for a 60-second shared-cache TTL while browsers continue to revalidate each visit. Protected release previews never receive this policy.
- `site.contentIndexPath` keeps the Content Index on an exact nginx location. This location returns `404` instead of the HTML fallback when the artifact is missing.
- Set `site.contentIndexPath` to an empty string to disable the exact Content Index location. This setting does not disable site-wide CORS.
- Set `site.healthPath` when the published site does not serve a successful response at `/`. The path must be public and stable.
- Set `storage.site.volumeName` only to bind an existing PV with copied site data. Keep the site's StorageClass and requested size consistent with that PV. A bound PVC cannot change its volume name or StorageClass in place.

## Setting up the encryption secret

1. Reuse the same plaintext value you already use locally for `MARKATA_GO_ENCRYPTION_KEY_DEFAULT`.
2. Seal it for the target namespace:

```bash
printf '%s' "$MARKATA_GO_ENCRYPTION_KEY_DEFAULT" | kubeseal \
  --controller-name sealed-secrets \
  --namespace <project>-notes \
  --name markata-go-encryption \
  --raw
```

3. Put the returned encrypted blob into your ArgoCD values:

```yaml
markataEncryption:
  secretName: "markata-go-encryption"
  sealedSecret:
    enabled: true
    encryptedData:
      MARKATA_GO_ENCRYPTION_KEY_DEFAULT: "Ag..."
```

4. If you also want the uploaded `source.tar.gz` encrypted, set:

```yaml
sourceArchive:
  encryption:
    enabled: true
```

Then publish with:

```bash
export MARKATA_GO_SOURCE_ARCHIVE_ENCRYPT=true
just push
```

## Generic defaults

- Set `project_identifier`, `project_name`, `sourceArchive.bucket`, and `sourceArchive.location` for your site before installing.
- Set `ingress.host` explicitly if you do not want the default `<project>.example.com` hostname.
- Enable `ingress.tls.enabled` and set `ingress.tls.secretName` when your ingress should terminate HTTPS with a named secret.
- Ingress auth is disabled by default. For builder-admin, set
  `builderAdmin.ingress.auth.internalUrl` to either an HTTPS auth origin or the
  cluster-local hlab-auth Service URL
  (`http://hlab-auth.hlab-auth.svc.cluster.local:<port 1-65535>`). The chart appends
  `/api/v1/forward-auth`; keep `publicAuthOrigin` on the browser-facing HTTPS
  auth origin.
- Use `build.extraEnv` for host-specific `MARKATA_GO_*` overrides such as `MARKATA_GO_URL` or `MARKATA_GO_SEARCH_ENDPOINT`.
- `storage.source.mode` and `storage.site.mode` default to `pvc` and can be set to `hostPath` for node-local content.
- Use `nodeSelector` when hostPath-backed content exists only on a specific node.
- `build.cacheDir` defaults to `/data/site/.cache/xdg` so tool caches survive between builds on the site volume.
- Enable `storage.cache` when you want a dedicated PVC for build caches. The build CronJob recreates `.markata` and `.markata-cache` in the source tree as symlinks into that cache volume, which keeps incremental state out of source sync paths.

## Offline builds

The official `ghcr.io/waylonwalker/markata-go-builder` image ships with:

- preloaded bundled CDN asset cache
- bundled Mermaid JS source
- `aws` and `openssl` for the source fetch/decrypt step

That means the chart can run without internet egress as long as your workload can still reach the configured source archive location and your site only needs assets that are already bundled or already cached.

Simple Helm values:

```yaml
offline:
  enabled: true
```

## Host-specific build overrides

Use `build.extraEnv` when the same content repo needs a different deploy hostname or search endpoint in one environment.

```yaml
build:
  extraEnv:
    - name: MARKATA_GO_URL
      value: "https://go.waylonwalker.com"
    - name: MARKATA_GO_SEARCH_ENDPOINT
      value: "https://go.waylonwalker.com/api/search"
```

If the search server itself also needs matching config overrides, set them under `search.extraEnv`.

## HostPath source and site mode

Use hostPath mode when your source repo or rendered site already lives on the node filesystem.

When `storage.source.mode` is `hostPath`, the build job skips the source archive download step and reads directly from the mounted directory.

The site volume still uses the chart's release layout under the mounted path:

- `releases/<timestamp-hostname>/`
- `current` symlink

That means a site hostPath should point at the root directory that serves the site, not directly at a single rendered release directory.

```yaml
storage:
  source:
    mode: hostPath
    hostPath:
      path: /mnt/main/walkershare/waylon/vaults/go.waylonwalker.com
      type: Directory
  site:
    mode: hostPath
    hostPath:
      path: /mnt/main/walkershare/waylon/sites/go.waylonwalker.com
      type: DirectoryOrCreate
```

This mode is node-local. The chart does not add any scheduling guard, so make sure the workloads land on a node where those paths exist.

For search-heavy dev loops, you can also point the search PVC at a node-local storage class and give it a new claim suffix so it rebuilds on local disk instead of Longhorn:

```yaml
storage:
  search:
    storageClassName: hostpath
    pvcNameSuffix: "-local"
```

The search deployment defaults to `Recreate` strategy because the bundled search index PVC is typically `ReadWriteOnce`, which prevents safe rolling updates across old and new pods.

In `watch-content` mode the search server runs with `/data/search` as its working directory and uses `search.configPath` to load the site config from the mounted source tree. This keeps `.markata/cache` and the Bleve index on writable storage instead of trying to write cache files into the source hostPath.

To preserve build-time tool caches between runs, keep `build.cacheDir` on persistent storage. The default points at the site volume so cached helper binaries and assets survive CronJob restarts instead of being lost with the pod `emptyDir`.

For deployments that sync the source tree with `rsync --delete` or otherwise treat `/data/source` as disposable, prefer a dedicated cache volume:

```yaml
storage:
  cache:
    enabled: true
    storageClassName: longhorn
    size: 10Gi
    accessModes:
      - ReadWriteOnce

build:
  cacheDir: /data/cache/xdg
```

That keeps XDG caches plus markata's `.markata` and `.markata-cache` incremental state outside the synced source tree while still letting the build recreate source-root symlinks on each run.

You can also enable a best-effort build notification hook:

```yaml
build:
  notify:
    enabled: true
    url: https://ntfy.example.com/topic
```

When enabled, the build CronJob posts a small result summary on success or failure. Notification failures do not fail the build.

For large sites using hardlinked releases, the build job now seeds a stable work directory at `/data/site/.build-work` from `current` before running `markata-go build`. That gives incremental write stages a real on-disk baseline for unchanged files while still publishing via an atomic `current` symlink swap at the end.

```yaml
nodeSelector:
  kubernetes.io/hostname: falcon3
```

## Example test rollout for `go.waylonwalker.com`

```yaml
project_identifier: go-waylonwalker-com
project_name: waylonwalker-com

sourceArchive:
  bucket: YOUR_BUCKET
  location: /YOUR_SHARE/outputs/waylonwalker-com

ingress:
  host: go.waylonwalker.com
  ingressClassName: traefik
  tls:
    enabled: true
    secretName: go-waylonwalker-com-tls

build:
  extraEnv:
    - name: MARKATA_GO_URL
      value: "https://go.waylonwalker.com"
    - name: MARKATA_GO_SEARCH_ENDPOINT
      value: "https://go.waylonwalker.com/api/search"
```

That lets you validate the chart on the existing `go.waylonwalker.com` hostname before switching `waylonwalker.com` itself.

When `offline.enabled` is true, the build and search workloads run with:

- `MARKATA_GO_OFFLINE=true`
- `MARKATA_GO_BUNDLED_ASSETS_CACHE_DIR=/usr/local/share/markata-go/assets-cache`
- `MARKATA_GO_BUNDLED_MERMAID_DIR=/usr/local/share/markata-go/mermaid`

If your site depends on extra self-hosted CDN assets beyond the bundled cache, pre-populate `.markata/assets-cache` in the source archive before publishing it:

```bash
markata-go assets download
just push
```

Or override `offline.bundledAssetsCacheDir` / `offline.bundledMermaidDir` to point at a custom preloaded location in your image.

## Manual rebuild

```bash
kubectl create job --from=cronjob/<project>-notes-build <project>-notes-build-manual-$(date +%s) -n <project>-notes
```

If another build is already running, the manual job waits on the build lock until the active build finishes or the lock becomes stale.
