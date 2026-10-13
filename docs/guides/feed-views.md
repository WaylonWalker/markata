---
title: "Feed Views"
description: "Choose default, simple, and calendar presentations globally or per feed"
date: 2026-10-02
published: true
tags:
  - documentation
  - feeds
---

# Feed views

Built-in feeds expose three presentation views by default: `default`, `simple`, and `calendar`. An opt-in fourth view, `shorts`, adds a full-screen photo/video player.
`default` is the feed's normal card/grid page and is always required. `simple` is the compact
`/simple/` page, and `calendar` is the interactive `?view=calendar` peer view.

To disable an alternate view for every feed, set site-wide feed defaults:

```toml
[markata-go.feed_defaults]
views = ["default", "calendar"]
```

A feed can override that set:

```toml
[[markata-go.feeds]]
slug = "notes"
views = ["default", "simple"]
```

Omitting `views` inherits the site defaults. Unknown view names are configuration errors, and an
explicit list must include `default`. Disabling `simple` also stops Markata from generating that
feed's `/simple/` HTML output even when `formats.simple_html` is enabled. Disabling `calendar`
removes Calendar controls and avoids loading its CSS, JavaScript, and hidden calendar source data.
Output-format settings for RSS, Atom, JSON, Markdown, text, and sitemaps are independent of views.

An empty list is invalid. Config overlays replace the view list rather than
adding to it. Keep `default` in every explicit list. Rebuilding after disabling
Simple removes its stale generated `/simple/` pages.

See [feed configuration](./configuration.md) and [the feed guide](./feeds.md)
for pagination and output-format options.

## Shorts view

The Shorts presentation uses the existing feed's selected posts and sort order, without copying the filter. To enable it for a Shots gallery:

```toml
[[markata-go.feeds]]
slug = "shots"
views = ["default", "simple", "calendar", "shorts"]
shorts_path = "/shorts/"

[markata-go.feeds.templates]
html = "feed-photo-grid.html"
```

The grid remains at `/shots/`, and the full-screen player is available at
`/shorts/`. A URL fragment such as `#id=shots%2Fdog` identifies an individual
post. The reader traverses the complete feed using static JSON chunks,
independent of normal feed pagination. Only one video plays at a time;
surrounding posts are loaded as previews. The previous and next buttons stay
fixed at the right center of the screen, separate from the video controls. Use
`j`/`k` or the up/down arrows to move, Escape to return to the grid, Space to
play or pause a video, left/right arrows to seek five seconds, and `m` to mute
or unmute.

If not configured, the Shorts route defaults to `/<feed-slug>/shorts/`. Only enable Shorts on a feed that should expose a full-screen player.
