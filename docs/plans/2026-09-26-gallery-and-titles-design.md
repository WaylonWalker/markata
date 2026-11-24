---
title: "Gallery pages and title sizing"
description: "Issue: #1289"
date: 2026-09-26
published: true
tags:
  - documentation
---


Issue: #1289

## Design

The default theme will provide `gallery.html` for a post whose `gallery` frontmatter is a list of images. Each item has a required `src` and `alt`, with optional `caption`, `width`, and `height`. The template renders an ordinary linked image gallery. CSS creates a responsive masonry layout; a small script upgrades links to a dialog with keyboard and touch navigation. Full image links remain useful without JavaScript. The theme asset helpers make the CSS and script available to projects without copied assets.

The existing `feed-photo-grid.html` already supports shots. Its cards will gain the word count used by Rhiannon's site for longer entries and load priority will favor only the first visible row. Site-specific colors and typography stay in site CSS.

Post titles will receive a build-time length category derived from the plain title. CSS uses it as a conservative initial size and preserves wrapping. An optional browser fit pass can refine the size against the byline after fonts load and on resize. This keeps the first paint readable and avoids assuming the build knows the final font metrics or viewport.

## Verification

Render a demo gallery and a shots feed from local fixture content, inspect the generated HTML, test the gallery viewer and long title in a browser, and run relevant Go tests. Include screenshots in the pull request when capture is available.
