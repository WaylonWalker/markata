# Container Images

markata-go provides official container images for minimal runtime execution, browser-free build/publish workflows, and optional Chromium-based compatibility builds.

## Image Names

- Minimal image: `ghcr.io/waylonwalker/markata-go:<version>`
- Builder image: `ghcr.io/waylonwalker/markata-go-builder:<version>`
- Browser-enabled builder: `ghcr.io/waylonwalker/markata-go-browser-builder:<version>`

Both images MUST publish a matching `:latest` tag for the most recent stable release.

## Minimal Image

The minimal image is intended for running a single `markata-go` command with the smallest footprint.

### Requirements

- Base image MUST be `scratch` or a distroless equivalent.
- Include the `markata-go` binary at `/usr/local/bin/markata-go`.
- Include TLS roots at `/etc/ssl/certs/ca-certificates.crt`.
- Default entrypoint MUST be `markata-go`.
- Default command SHOULD be `--help`.
- Working directory MUST be `/site`.
- Default user SHOULD be a non-root UID/GID (e.g., `65532:65532`).

## Builder Image

The builder image is intended for CI pipelines, sidecar builders, and scripts that need shell tooling.

### Required Tools

The builder image MUST include:

- `markata-go` binary
- `aws` CLI
- POSIX shell at `/bin/sh`
- Filesystem utilities: `find`, `sort`, `xargs`, `sha256sum`, `awk`
- Core utilities: `date`, `ln`, `rm`, `mkdir`, `cp`, `mv`, `sleep`
- `openssl` for source archive decryption and related build workflows
- `rsync` for publish/sync workflows
- `ca-certificates` for HTTPS access

### Optional Tools

The builder image MAY include:
- `git`
- `libavif-apps` (provides `avifenc`)
- `libwebp-tools` (provides `cwebp`)
- `pagefind` (standalone binary for search index generation)
- `tailwindcss` (standalone Tailwind CLI for CSS builds)
- `openssh-client` (for rsync over SSH)
- `tzdata`

Native Mermaid rendering uses an in-process Go SVG renderer and does not need Chromium, Node.js, npm, or mmdc. The standard builder MUST NOT install Chromium. The separate optional browser-enabled builder MUST install Chromium for explicit compatibility-mode rendering, and MAY otherwise inherit all builder-runtime tools. Pagefind is installed as a standalone musl-static binary.

### Behavior

- The builder image MUST allow running shell scripts directly.
- `docker run --rm <builder-image> sh -c 'markata-go --help'` MUST succeed.
- `docker run --rm <builder-image> sh -c 'aws --version'` MUST succeed.
- `docker run --rm <builder-image> sh -c 'openssl version'` MUST succeed.
- `docker run --rm <builder-image> sh -c 'rsync --version'` MUST succeed.

## Tagging and Versioning

- All container tags MUST follow semantic versioning.
- All three images MUST use the same tag for a given release.
- `:latest` MUST point to the most recent stable release.

## Helm Site Server

The Helm chart serves published site files through nginx.

- Public site responses MUST include `Access-Control-Allow-Origin: *` by default.
- Operators MUST be able to disable this header with `site.cors.enabled: false`.
- Protected release previews MUST NOT include the public CORS header.
- The chart MUST preserve cache headers when it adds the CORS header.
- Public HTML responses MUST remain `no-cache` by default. Operators MAY opt into
  `site.cacheHtmlAtEdge: true`, which sets browser `max-age=0` and a 60-second
  shared-cache TTL. Protected release previews MUST NOT receive this policy.
- The configured Content Index path MUST return `404` when its artifact is missing.
- The site PVC MAY specify `storage.site.volumeName` to bind an existing PV.
  The chart MUST omit this field unless configured. A bound PVC's storage class
  and volume name are immutable; operators MUST migrate data before changing
  either field.
