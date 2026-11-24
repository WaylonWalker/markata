---
title: "Image Library"
description: "Build a searchable image inventory and authoring page with markata-go"
date: 2026-09-13
published: true
slug: /docs/guides/image-library/
tags:
  - images
  - configuration
  - authoring
---


markata-go can generate an image inventory and a page for finding images used
by your site. The default output is:

- `/images/` — a responsive, searchable image library
- `/images/index.json` — a versioned machine-readable inventory
- `/images.json` — the same inventory at the site root for other tools

The plugin runs by default. It reads public posts and, when
`include_unreferenced` is enabled, the configured static asset directory. It
does not fetch remote images or videos.

## Configure the output

Add an `images` table under `[markata-go]` when you need to change the
defaults:

```toml
[markata-go.images]
enabled = true
path = "images"
template = "images.html"
export_json = true
include_unreferenced = false
```

| Option | Default | Description |
|---|---|---|
| `enabled` | `true` | Generate the image library page and inventory. |
| `path` | `"images"` | Relative output directory below `output_dir`. |
| `template` | `"images.html"` | Template used for the HTML page. |
| `export_json` | `true` | Write `path/index.json` and `/images.json`. |
| `include_unreferenced` | `false` | Include image and video files found in `assets_dir` even when no public post uses them. Set `true` explicitly for a local authoring inventory because it makes all supported asset files enumerable. |

The output path must stay inside `output_dir`. The same options work in YAML
and JSON configuration files. Environment overrides use these names:

```text
MARKATA_GO_IMAGES_ENABLED
MARKATA_GO_IMAGES_PATH
MARKATA_GO_IMAGES_TEMPLATE
MARKATA_GO_IMAGES_EXPORT_JSON
MARKATA_GO_IMAGES_INCLUDE_UNREFERENCED
```

## Find and copy an image

The library searches image filenames, paths, alt text, figure captions, and
public post titles. Use the `All`, `Used`, `Unused`, `Cover`, and `Recently
added` views, or choose a sort order from the toolbar.

Each card includes **Copy Markdown** and **Copy URL** actions. Both actions use
the canonical `src` value from the inventory. For example:

```markdown
![Workshop](https://dropper.wayl.one/file/workshop.webp)
```

Trusted Dropper URLs may use resized URLs for previews. Those presentation URLs
never replace the canonical URL in the JSON artifact or copied Markdown.

## Image discovery

The inventory includes:

- Markdown image syntax, supported raw HTML `<img>` elements, and referenced
  `<video>`/`<source>` elements
- normalized visible text from associated Markdown or HTML `<figcaption>`
  elements, including captions shared by multiple media elements in one figure
- `image`, `video`, `cover`, `cover_image`, `og_image`, `social_image`, `thumbnail`,
  `featured_image`, `hero_image`, `avatar`, and `author_image` frontmatter
- Image and video files below `assets_dir`

The `image` frontmatter field is a cover fallback. If `cover` or `cover_image`
has a non-empty value, that field remains the first cover choice. Video URLs
are kept as video media with their video MIME type. The library renders them
with a `<video>` element and uses an explicit poster field or the same derived
poster convention used by feed and video cards.

Images supplied by an external `![embed](...)` card are still inventoried, but
they are marked with an **Embed** label. This distinguishes a remote OG image
from an image authored directly in a post.

Skipped, draft, private, and unpublished posts do not create public usage
relationships. Local files are mapped from the asset root to site-root URLs;
for example, `static/images/logo.png` becomes `/images/logo.png`.

Private body isolation is strict: the image library never scans private
Markdown, rendered HTML, captions, alt text, or remote URLs. Changing only a
private body produces identical public image-library output. A private post may
opt its public-safe `cover` frontmatter into the inventory, with `cover_alt` as
its only allowed alt field. This creates metadata only: it never creates a
`uses[]` relationship or contributes `added_at`/`last_used_at`. Other private
image fields are ignored. Private frontmatter is not presumed public; consumers
may use only explicitly documented public-safe fields. Keep private media
outside `assets_dir` when it must not be published by the static asset writer.

Local canonical URLs escape path components while preserving `/`. For example,
`my photo#1%.png` becomes `/my%20photo%231%25.png`; authored percent-encoded
references deduplicate to that same record. Remote query strings are preserved
for signed URLs, but URLs with `user:password@host` credentials are rejected.
Referenced local media is resolved and hashed when it is below `assets_dir`,
including an absolute `assets_dir` outside `content_dir`. Other local files are
eligible when they are below `content_dir`; files outside both roots are ignored.
Platform-native filesystem paths are supported, including Windows drive-letter
paths.

The public `uses[]` entries contain the public `href`, title, caption, and
flags. They do not contain repository-relative source paths.

The versioned v1 reader requires the documented top-level fields and the stable
per-image fields (`src`, `width`, `height`, `alt`, `mime_type`, `cover`, and
`uses`). Each usage requires `href` and `cover`. Unknown fields are ignored so
newer producers can add metadata without breaking readers.

The JSON artifact is deterministic. Images are sorted by source URL and usage
links are sorted by public href. Repository-relative source paths are not
included in public usage records. A missing local dimension is
reported as `0`. Local PNG, JPEG, GIF, WebP, BMP, and TIFF files report their
dimensions when the file metadata is valid. SVG files use positive `width` and
`height` values or a positive `viewBox`; ICO files use their largest directory
entry. AVIF and HEIC/HEIF records are supported, but their dimensions remain
`0` until a decoder is available. Malformed or unsupported media does not fail
the build.

Each used image includes `added_at` when a public post has a valid publication
date; it is the earliest such date in UTC. `last_used_at` is the latest public
publication date in UTC. These dates come from post frontmatter, not local file
timestamps, so remote media and local media use the same rules. Missing values
sort last. The `Recently added` view shows sources with a positive `added_at`
timestamp from the last 30 days according to the browser clock.
Each usage relationship can also include its figure caption. Captions stay with
the usage relationship because one media source can have different captions in
different posts.

The image library keeps a small media fingerprint cache for incremental builds.
When a local file's path, size, and modification time are unchanged, markata-go
reuses its content fingerprint without reopening the file. On filesystems that
expose an operating-system change time, that signal also detects replacements
that preserve size and modification time. Filesystems without that signal must
keep modification times reliable for this optimization; a changed, added,
removed, or renamed media file with changed cache state invalidates the
generated inventory.

The writer checks for collisions with posts, feeds, and static files before it
writes. It does not overwrite an existing site file. When you disable the
library or change its path, cleanup removes only files previously recorded as
image-library outputs whose contents are unchanged. A modified or unowned
stale JSON file is preserved and reported as a collision. Output paths that
traverse symlinks below the configured output root are rejected for safety.
Platform-managed ancestors above that root are allowed.

## Custom templates

Copy `images.html` into your site's `templates/` directory to customize the
page. The template receives:

- `image_library` — page counts and presentation-ready image cards
- `image_index` — the canonical `imageindex.Index`
- `image_library_config` — the resolved image configuration

Image cards expose `IsVideo`, `PosterSrc`, `Embed`, `AddedAt`, `AddedAtUnix`,
`LastUsedAt`, and `LastUsedAtUnix` in addition to the regular image fields.
`AddedAt` is the earliest valid public post date for the source. Each item in
`Uses` can also expose `Caption` when the media appears inside a figure. Custom
templates should render `IsVideo` cards with a `<video>` and `<source>` rather
than an `<img>`.

The bundled page keeps video URLs out of active `src` and `<source src>`
attributes in the initial HTML while keeping the poster active as the card's
thumbnail. Its JavaScript loads each video when its card enters the viewport,
so an image-library page with many videos does not fetch every video at once.
Without JavaScript, the card still links to the canonical media URL.

The server-rendered page remains useful when JavaScript is unavailable. The
bundled JavaScript only enhances filtering, sorting, usage expansion, copy
feedback, and video loading.

See the [configuration reference](/docs/guides/configuration/) and [built-in
plugin reference](/docs/reference/plugins/) for related settings.


## Optional server-side pagination

For large libraries, avoid sending every image card on the initial request:

```toml
[markata-go.images]
page_size = 48
```

With this opt-in setting, `/images/` renders up to 48 cards and emits native
`/images/page/2/`, `/images/page/3/`, etc. links. The JSON inventory still
contains all public media; no published media relationships are lost. The
page links work without JavaScript. Global search, filters, and sort are still
available with JavaScript: on the first interaction the browser retrieves the
other *static HTML* pages and searches across the complete collection, instead
of loading all their DOM nodes during first paint. Until that interaction, the
initial page stays light and usable. If another page fails to load, the links
remain available for ordinary browsing.

A setting of `page_size = 0` (default) preserves the existing single-page
library. Sites with a custom `images.html` should update their template
to honor `image_library.PageLinks` before enabling pagination.
