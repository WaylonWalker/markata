---
title: "Pins"
description: "Build a theme-aware link board from ordinary post metadata."
date: 2026-09-27
published: true
slug: /docs/guides/pins/
tags:
  - documentation
  - feeds
  - pins
---


Markata includes a built-in `/pins/` page for saved links. Published posts appear
there when they have a non-empty `link` or `url` field. Items appear newest first.

```yaml
---
title: "A field guide to useful maps"
link: "https://example.com/field-guide"
description: "A thoughtful guide to reading a landscape."
image: "/images/maps.webp"
date: 2026-09-27
published: true
---

A note about why this is worth saving.
```

Link posts can also use `url` for their destination.

The pinboard keeps the title linked to your local post and provides a separate source link to the external URL. Images, descriptions, dates, tags, and authored notes remain ordinary post data.

The default board uses the available viewport width, with responsive masonry
columns and 100 pins per page. The header has extra breathing room; card footers
share the card background. Hovering subtly emphasizes cover art without moving
cards or opening floating notes.

## Cover art and notes

Pins uses `image`, `cover`, `cover_image`, or `og_image`, in that order. When those
fields are empty, it reuses image previews from already-rendered embeds, YouTube
embed thumbnails, or images in the post body. The pins page itself does
not fetch metadata for arbitrary posts. Link-style posts automatically generate
a preview from their `link` or `url` frontmatter destination. For other posts,
use the normal embed syntax to supply an external preview:

```markdown
![[https://example.com/field-guide]]

I saved this for its clear examples and detailed maps.
```

Your commentary appears below the title as a short note. Click or tap **Note**,
or activate it with the keyboard, to expand the text inside the card. Embedded
destination descriptions and standalone source URLs are excluded from the note.
If there are no authored paragraphs or list items, a non-URL `description` is used
instead. Both cover art and notes work without JavaScript.

## Customize pagination

Define a Pins feed to use a different page size or pagination mode:

```toml
[[markata-go.feeds]]
slug = "pins"
title = "Pins"
filter = "published == true and (link or url)"
sort = "date"
reverse = true
items_per_page = 50
pagination_type = "manual"

[markata-go.feeds.templates]
html = "pins.html"

[markata-go.feeds.formats]
html = true
```

See [feeds](./feeds.md) for pagination and format options, and
[templates](./templates.md) for overriding the board design.

The `link` field is intentionally render-neutral outside the pinboard. Adding `link:` does not inject an external preview into the individual post and does not cause an external metadata request by itself. Existing sites can therefore keep using `link` as ordinary frontmatter without changing post rendering.

## Route ownership

The built-in Pins feed is a convenience route, not a reservation. If your site already has a publishable post with the slug `pins`, that post keeps `/pins/` and Markata does not inject the implicit feed. A draft, skipped, or unpublished `pins` post does not reserve the route.

An explicitly configured feed with the slug `pins` also takes precedence over the built-in definition.

## Default behavior

The implicit Pins feed uses:

```text
filter: published == true and (link or url)
sort: date
reverse: true
items_per_page: 100
pagination_type: manual
format: HTML
path: /pins/
template: pins.html
```

Disabling built-in subscription feeds disables the implicit Pins feed along with the implicit root and archive feeds. You can still define your own `pins` feed when you want different selection, sorting, pagination, templates, or output formats.
