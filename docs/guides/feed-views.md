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

## Adaptive HTML layouts

The main HTML feed can use `classic`, `editorial`, `gallery`, `reading`,
`timeline`, or `resources`. Select `auto` to choose from the feed's post
types. The default is **classic** to preserve existing sites.

```toml
[[markata-go.feeds]]
slug = "archive"
presentation = "editorial"

[[markata-go.feeds]]
slug = "ping"
presentation = "timeline"

[[markata-go.feeds]]
slug = "stars"
presentation = "resources"
```

Site-wide defaults may be set in `[markata-go.feed_defaults]` with
`presentation = "classic"` (or a different layout). A feed's own
`presentation` overrides that default.

`auto` chooses a specialized layout when at least 80% of eight or more
posts are photos, notes, links, or articles. Feeds with at least twelve
posts and three content types represented by two posts each use editorial
grouping. Smaller or unclassified feeds retain the classic view.
Custom HTML templates, such as `feed-photo-grid.html`, take precedence
over automatic layout detection.

Editorial grouping uses Monday-start publishing weeks. Weekly headings show relative
calendar time (This week, Last week, 2 weeks ago, 2 months ago, 2 years ago)
with the exact week range as a smaller caption. Relative labels are rendered
at build time and refreshed in the browser so long-cached static pages remain
accurate without a rebuild. Dates and weekly grouping remain unchanged.
Articles are
emphasized, substantial links and video follow, and all photos from a
publishing week appear together in a single responsive gallery. One-photo
weeks also use the gallery layout rather than stretching a photo into a
full-width card. Each photograph retains its own publication date, caption,
and permalink, even when photos come from different days. Brief notes are
visually quieter. `featured: true` on a post explicitly elevates it
within its week. Each post retains its original HTML card and metadata.

Presentation never modifies RSS, Atom, JSON, Simple, Calendar, source
ordering, or the existing pagination boundaries. Weeks that span two
HTML pagination pages continue on the next page. JS-pagination feeds
keep classic markup because the JS paginator requires direct child cards.
