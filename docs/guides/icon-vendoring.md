---
title: "Vendoring Icon Packs"
description: "Download pinned SVG icon packs for local icon shortcodes"
published: true
slug: /docs/guides/icon-vendoring/
tags:
  - documentation
  - markdown
  - icons
  - assets
---


markata-go has two related icon-vendoring paths:

1. **Lucide is automatic on first use.** A shortcode such as `:lucide-smile:` can be used with no configuration. If no local Lucide asset exists, markata-go vendors its pinned default `lucide-static` pack, then resolves the shortcode from the resulting local SVG tree.
2. **Explicit vendoring remains available** for custom Lucide pins, additional packs, custom archive sources, or builds that want missing packs to be a hard error.

Ordinary builds that never reference a missing Lucide icon do not download the default pack.

## Explicit Lucide pin from npm

```toml
[icons.vendor]
enabled = true

[[icons.vendor.packs]]
name = "lucide"
version = "1.48.0"
source = "npm"
package = "lucide-static"
icons_path = "icons"
license_path = "LICENSE"
```

On the first build, markata-go downloads and caches the pinned npm tarball, copies its SVGs to `static/.icons/lucide/`, and writes the pack license to `static/.icons/licenses/lucide.txt`.

You can then use the local icons normally:

```markdown
:lucide-smile: Hello!

:lucide/github: Source code
```

The first form follows Zensical's hyphenated shortcode spelling. The slash form is a markata-go alias for the same local asset.

Explicit vendoring runs during Configure, so it downloads/materializes configured packs whether or not a particular icon appears in the current content. The zero-config Lucide fallback is different: it waits until a missing Lucide shortcode is encountered.

## Disable the automatic Lucide fallback

If you want icon rendering to remain strictly local unless `[icons.vendor]` is explicitly configured:

```toml
[icons]
auto_vendor = false
```

This does not disable explicitly configured `[icons.vendor]` packs.

## Pin versions

Every explicitly vendored pack requires a version. This keeps the cache key and generated icon tree deterministic rather than silently following `latest`.

The built-in Lucide fallback is also pinned by markata-go. When markata-go updates that pin, already cached/materialized packs continue to follow the normal vendor fingerprint and replacement behavior.

When you explicitly upgrade a pack, change its configured version and rebuild. markata-go replaces that pack's materialized icon directory so icons removed upstream do not remain as stale files.

## Custom cache and target

```toml
[icons.vendor]
enabled = true
cache_dir = ".cache/markata-icons"
target = "static/.icons"
```

Defaults for explicit vendoring are:

- cache: `.markata/cache/icon-packs`
- target: `static/.icons`

The zero-config Lucide fallback reuses those same vendor/cache mechanics and selects an icon root that the shortcode renderer already scans.

## Pack-specific archive layout

Different packages arrange SVGs and licenses differently. Configure those archive-relative paths per pack:

```toml
[[icons.vendor.packs]]
name = "my-pack"
version = "2.4.0"
source = "npm"
package = "my-icon-package"
archive_path = "package"
icons_path = "dist/icons"
license_path = "LICENSE.md"
```

`archive_path` defaults to `package`, which matches standard npm tarballs. `icons_path` defaults to `icons`, and `license_path` defaults to `LICENSE`.

## Direct `.tgz` archives

You can also use a versioned HTTP or HTTPS tarball directly:

```toml
[[icons.vendor.packs]]
name = "house-icons"
version = "2026.09.1"
source = "url"
url = "https://example.com/house-icons-2026.09.1.tgz"
archive_path = "package"
icons_path = "svg"
license_path = "LICENSE.txt"
```

Direct sources use the same archive cache and extraction safety checks as other markata-go vendored assets.

## Offline builds

The vendor plugin reuses markata-go's asset cache. After a successful online build, deleting a materialized icon directory and rebuilding can recreate it from the cache without downloading the archive again.

To prohibit runtime network access entirely:

```bash
MARKATA_GO_OFFLINE=1 markata-go build
```

For **explicit** `[icons.vendor]` packs, offline mode requires the archive cache or an existing materialized pack; otherwise the build fails.

For the **automatic Lucide fallback**, a cache miss does not fail the site build. The unresolved shortcode remains literal. This lets normal offline builds keep working even when a previously unused Lucide icon appears.

## Licenses

Each vendored pack must provide the file named by `license_path`. markata-go copies it to:

```text
static/.icons/licenses/<pack-name>.txt
```

That preserves the upstream license alongside the files being redistributed. You are still responsible for confirming that a pack's license is appropriate for your site.

## SVG safety

Downloaded SVGs are checked before publication. markata-go rejects a pack if an SVG contains common active-content constructs such as scripts, `foreignObject`, JavaScript URLs, or inline event handlers.

Only `.svg` files under the configured icon directory are copied; symlinks are ignored.

## Multiple packs

Add more `[[icons.vendor.packs]]` blocks as needed. Pack names must be unique because the name becomes the first directory component and shortcode namespace.

```toml
[icons.vendor]
enabled = true

[[icons.vendor.packs]]
name = "lucide"
version = "1.48.0"
source = "npm"
package = "lucide-static"

[[icons.vendor.packs]]
name = "brand"
version = "1.3.0"
source = "url"
url = "https://example.com/brand-icons-1.3.0.tgz"
icons_path = "svg"
license_path = "LICENSE.txt"
```

Because vendoring only populates local files, rendered pages do not perform browser-time requests to an icon service or CDN.
