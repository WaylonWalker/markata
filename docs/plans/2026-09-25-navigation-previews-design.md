---
title: "Navigation Preview Design"
description: "Build-time feed and post previews for the default navigation component"
date: 2026-09-25
published: false
tags:
  - design
  - navigation
---

## Navigation previews

The default header and footer navigation components add a compact editorial card to local links that resolve to a generated feed, a published post, or an enabled Reader, Blogroll, or Random route. The card appears on hover and keyboard focus. It uses the current theme colors, a subtle underline cue, and a short entrance animation. On narrow and touch layouts, the card is hidden so the scrollable nav remains usable. Reduced-motion settings remove the animation.

The render stage resolves configured nav URLs against feed configs and public posts once per build. It computes the preview data from the same feed membership used to publish feeds. Feed cards include available description, public post count, total words, reading time, and a monthly publication sparkline. Post cards include available description, tags, word count, and reading time. Private post metadata is excluded. External and unmatched URLs remain plain links.

Both template trees use one link partial. A hash of the preview data invalidates cached pages and feeds when shared metadata changes. Focused tests cover URL resolution, privacy, statistics, sparkline presence, generated routes, footer links, and rendered output.
