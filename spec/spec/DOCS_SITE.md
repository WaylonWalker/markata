---
title: "Documentation Site"
description: "Rendering and asset requirements for the separate documentation website."
date: 2026-11-01
published: true
tags:
  - specification
  - documentation
---

# Documentation Site

The public documentation website MUST be built from `docs/**/*.md` by the
Markata Go CLI using `markata-docs.toml`. The generated site is written to
`docs-site/_site/` for local preview and CI validation. Production builds MUST
use the same configuration and renderer through Builder Admin.

## Markdown rendering

- Fenced code blocks MUST be highlighted by Markata when a block declares a
  supported language.
- Highlighting MUST remain present in the generated HTML without client-side
  JavaScript.
- The site MUST use templates bundled with Markata's default theme, without
  templates from the Markata repository's main-site override directory.
- The default theme MUST provide a documentation page template that renders
  the configured documentation sidebar and post content.
- Rendered documentation pages MUST have one page-level H1. Markdown content
  MUST follow Markata's frontmatter-title convention and use H2 headings for
  sections below the generated page title.
- The docs site MUST use the built-in template and theme behavior so its own
  public documentation exercises Markata's default rendering path.
- The docs layout MUST render a hierarchical sidebar from Markdown files under
  `docs/`, so newly added guides appear in navigation without a separate nav
  list.
- Markdown files under `docs/` MUST preserve their `docs/` URL prefix unless a
  page declares a custom slug.
- `docs/index.md` MUST be the root homepage through an explicit empty slug.

## Site assets

- Every generated page MUST link to the site's favicon.
- The build output MUST contain the favicon at the linked path.
- Static files under `static/` MUST be published, including `favicon.ico`.
- The build MUST include all stylesheets linked by generated pages.

## Production publishing

- Production MUST run Markata Builder Admin from the Markata Helm chart.
- The builder MUST source its checkout from the configured Git repository and
  branch.
- A signed GitHub push webhook MUST fast-forward the configured branch and
  enqueue a build only when source changes.
- The webhook HMAC secret MUST be supplied from a Kubernetes Secret.
- Published output MUST use retained, atomic releases and preserve the current
  live release if a build fails.
