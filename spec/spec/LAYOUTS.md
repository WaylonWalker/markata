# Layout System Specification

The layout system provides configurable page structures for different content types. It enables documentation sites, blogs, landing pages, and custom layouts with zero-config defaults and progressive customization.

## Design Principles

1. **Zero-config beautiful** - `layout: docs` gives a professional 3-panel documentation layout immediately
2. **Progressive customization** - Enable components, configure them, override per-page
3. **Responsive by default** - All layouts adapt to mobile with collapsible panels
4. **Component-based** - Sidebar, TOC, header, footer are independent, composable components
5. **Theme integration** - Layouts use the existing palette and CSS variable system
6. **Template inheritance** - Layouts extend `base.html` using pongo2 blocks

---

## Layout Types

Markata-go provides four built-in layout presets. Each preset configures components appropriately for its use case.

### Layout Preset Overview

| Layout | Sidebar | TOC | Header | Footer | Best For |
|--------|---------|-----|--------|--------|----------|
| `docs` | Left (nav) | Right | Minimal | Minimal | Documentation sites |
| `blog` | None | Optional | Full | Full | Blog posts, articles |
| `landing` | None | None | Minimal | Full | Home pages, marketing |
| `bare` | None | None | None | None | Embeds, custom pages |

### `docs` Layout

The 3-panel documentation layout optimized for technical content.

```
+------------------------------------------------------------------+
|                         HEADER (minimal)                          |
|  [Logo]              [Nav Links]              [Search] [Theme]    |
+----------+------------------------------------------+-------------+
|          |                                          |             |
|  SIDEBAR |              CONTENT                     |     TOC     |
|  (nav)   |                                          |  (on-page)  |
|          |  # Page Title                            |             |
|  [Home]  |                                          |  - Section  |
|  [Guide] |  Content here...                         |  - Section  |
|  [API]   |                                          |    - Sub    |
|  [FAQ]   |                                          |  - Section  |
|          |                                          |             |
+----------+------------------------------------------+-------------+
|                         FOOTER (minimal)                          |
+------------------------------------------------------------------+
```

**Default configuration:**
```toml
[markata-go.layout]
name = "docs"

[markata-go.layout.docs]
sidebar_position = "left"
sidebar_width = "280px"
toc_position = "right"
toc_width = "220px"
content_max_width = "800px"
header_style = "minimal"
footer_style = "minimal"

# TOC Configuration (new in v0.x.x)
[markata-go.layout.docs.toc]
enabled = true                              # Enable TOC display
auto_enable = true                          # Auto-show based on thresholds
min_toc_links = 3                           # Min headings for auto-show
min_word_count = 500                        # Min words for auto-show
```

### `blog` Layout

Traditional blog layout with optional TOC for long-form content.

```
+------------------------------------------------------------------+
|                           HEADER                                  |
|  [Logo]    [Home] [Blog] [About] [Contact]    [Search] [Theme]   |
+------------------------------------------------------------------+
|                                                                   |
|                          CONTENT                                  |
|                                                                   |
|  # Article Title                                                  |
|  Posted on January 15, 2024 by Author                            |
|  Tags: [go] [static-site]                                        |
|                                                                   |
|  Article content here with full width...                         |
|                                                                   |
|  [Previous Post]                      [Next Post]                |
|                                                                   |
+------------------------------------------------------------------+
|                           FOOTER                                  |
|  Copyright 2024  |  RSS  |  Twitter  |  GitHub                   |
+------------------------------------------------------------------+
```

**Default configuration:**
```toml
[markata-go.layout]
name = "blog"

[markata-go.layout.blog]
content_max_width = "720px"
show_toc = false          # Show the TOC sidebar on blog-layout pages
toc_position = "right"
header_style = "full"
footer_style = "full"
show_author = true
show_date = true
show_tags = true
show_reading_time = true
show_prev_next = true
```

**Table of contents (`show_toc`):** The `blog` layout renders with `post.html`.
When `show_toc = true`, pages rendered with `post.html` MUST show the table of
contents sidebar (`components/doc_sidebar.html`) for posts that have TOC
entries. The sidebar uses `toc_position` and `toc_width` from
`[markata-go.layout.blog]`, and `min_depth`/`max_depth` from
`[markata-go.components.doc_sidebar]`.

- If `[markata-go.components.doc_sidebar] enabled = true`, that component's
  settings take precedence on all pages, and `show_toc` has no extra effect.
- `show_toc` MUST NOT enable the sidebar for pages rendered with other
  templates (for example `docs`, `landing`, `bare`, or custom templates).
- Templates read the effective settings from the `resolved_doc_sidebar`
  context variable (`enabled`, `position`, `width`, `min_depth`, `max_depth`).

### `landing` Layout

Full-width layout for marketing pages and home pages.

```
+------------------------------------------------------------------+
|  HEADER (minimal, transparent)        [Nav] [CTA Button]         |
+------------------------------------------------------------------+
|                                                                   |
|                     HERO SECTION                                  |
|                                                                   |
|              Your Amazing Product                                 |
|         Build beautiful sites with ease                          |
|                                                                   |
|              [Get Started]  [Learn More]                         |
|                                                                   |
+------------------------------------------------------------------+
|                                                                   |
|                   CONTENT (full width)                            |
|                                                                   |
|  Features, testimonials, pricing, etc.                           |
|                                                                   |
+------------------------------------------------------------------+
|                           FOOTER                                  |
+------------------------------------------------------------------+
```

**Default configuration:**
```toml
[markata-go.layout]
name = "landing"

[markata-go.layout.landing]
content_max_width = "100%"
header_style = "transparent"
footer_style = "full"
hero_enabled = true
```

### `bare` Layout

Minimal layout with no chrome - just the content.

```
+------------------------------------------------------------------+
|                                                                   |
|                         CONTENT                                   |
|                                                                   |
|  Pure content, no header, no footer, no navigation               |
|                                                                   |
+------------------------------------------------------------------+
```

**Default configuration:**
```toml
[markata-go.layout]
name = "bare"

[markata-go.layout.bare]
content_max_width = "100%"
```

---

## Configuration

### Global Layout Configuration

Set site-wide layout defaults in `markata-go.toml`:

```toml
[markata-go.layout]
# Default layout for all pages
name = "docs"

# Global component settings (apply to all layouts unless overridden)
[markata-go.layout.defaults]
content_max_width = "800px"
header_sticky = true
footer_sticky = false
```

### Layout-Specific Configuration

Configure each layout preset:

```toml
# Documentation layout
[markata-go.layout.docs]
sidebar_position = "left"     # "left" | "right"
sidebar_width = "280px"
sidebar_collapsible = true
sidebar_default_open = true   # Desktop default
toc_position = "right"        # "left" | "right" (opposite of sidebar)
toc_width = "220px"
toc_collapsible = true
toc_default_open = true
content_max_width = "800px"
header_style = "minimal"      # "full" | "minimal" | "transparent" | "none"
footer_style = "minimal"      # "full" | "minimal" | "none"

# Blog layout
[markata-go.layout.blog]
content_max_width = "720px"
show_toc = false
toc_position = "right"
toc_width = "200px"
header_style = "full"
footer_style = "full"
show_author = true
show_date = true
show_tags = true
show_reading_time = true
show_prev_next = true

# Landing page layout
[markata-go.layout.landing]
content_max_width = "100%"
header_style = "transparent"
header_sticky = true
footer_style = "full"

# Bare layout
[markata-go.layout.bare]
content_max_width = "100%"
```

### Component Configuration

Individual components can be configured globally:

```toml
# Sidebar component
[markata-go.components.sidebar]
enabled = true
position = "left"
width = "280px"
collapsible = true
default_open = true
# Navigation structure (auto-generated from feeds if not specified)
# Or manually defined:
[[markata-go.components.sidebar.nav]]
title = "Getting Started"
href = "/docs/getting-started/"

[[markata-go.components.sidebar.nav]]
title = "Guides"
children = [
    { title = "Configuration", href = "/docs/guides/configuration/" },
    { title = "Themes", href = "/docs/guides/themes/" },
    { title = "Layouts", href = "/docs/guides/layouts/" },
]

[[markata-go.components.sidebar.nav]]
title = "Reference"
children = [
    { title = "CLI", href = "/docs/reference/cli/" },
    { title = "Plugins", href = "/docs/reference/plugins/" },
]

# Table of Contents component
[markata-go.components.toc]
enabled = true
position = "right"
width = "220px"
min_depth = 2          # Start at h2
max_depth = 4          # End at h4
title = "On this page"
collapsible = true
default_open = true
scroll_spy = true      # Highlight current section

# Header component
[markata-go.components.header]
style = "full"         # "full" | "minimal" | "transparent" | "none"
sticky = true
show_logo = true
show_title = true
show_nav = true
show_search = true
show_theme_toggle = true

# Footer component
[markata-go.components.footer]
style = "full"         # "full" | "minimal" | "none"
sticky = false
show_copyright = true
copyright_text = "Copyright 2024 My Company"
show_social_links = true
show_nav_links = true

# Feed sidebar component
[markata-go.components.feed_sidebar]
enabled = true
position = "left"
width = "250px"
feeds = ["docs", "guides"]
max_posts = 51
```

The feed sidebar component renders navigation for posts that share a feed with the current post. Feed selection order is:

1. series metadata
2. explicit post frontmatter such as `prev_next_feed` or `sidebar_feed`
3. configured `feed_sidebar.feeds`
4. automatic discovery from feeds that contain the current post

Automatic discovery MUST choose the smallest matching feed as the most specific navigation set. When two matching feeds have the same size, a primary feed wins the tie.

`max_posts` limits the **initial server-rendered sidebar window**. With JavaScript, public HTML feeds progressively fetch additional static navigation chunks as the reader scrolls; the setting does not limit how many posts can eventually be reached. Without JavaScript, the initial links remain available and the full feed page remains linked.

`max_posts` limits the rendered sidebar items. A value of `0` or less means no limit. When `max_posts` is positive and the matching feed has more posts than the limit, implementations MUST render a contiguous window that includes the current post and SHOULD center that window around the current post when possible. If the current post cannot be found in the selected feed, implementations SHOULD render the first `max_posts` entries.

---

## Frontmatter Overrides

Any layout setting can be overridden per-page using frontmatter.

### Basic Override

```yaml
---
title: "API Reference"
layout: docs
---
```

### Component Overrides

```yaml
---
title: "Quick Start Guide"
layout: docs

# Override specific components for this page
layout_config:
  sidebar:
    enabled: true
    position: left
  toc:
    enabled: true
    position: right
    min_depth: 2
    max_depth: 3
---
```

### Disable Components

```yaml
---
title: "Full-Width Demo"
layout: docs

# Disable sidebar and TOC for this page
layout_config:
  sidebar:
    enabled: false
  toc:
    enabled: false
  content_max_width: "100%"
---
```

### Switch Layout

```yaml
---
title: "Welcome to Our Docs"
layout: landing   # Use landing layout instead of default

layout_config:
  hero_enabled: true
  header_style: transparent
---
```

### Feature Flags (Shorthand)

Common overrides have shorthand frontmatter keys:

```yaml
---
title: "My Post"
toc: false              # Shorthand for layout_config.toc.enabled = false
sidebar: false          # Shorthand for layout_config.sidebar.enabled = false
full_width: true        # Shorthand for layout_config.content_max_width = "100%"
---
```

---

## Components

### Sidebar Component

The sidebar provides navigation for documentation and multi-page content.

#### Structure

```html
<aside class="layout-sidebar" data-position="left" data-collapsible="true">
  <button class="sidebar-toggle" aria-label="Toggle sidebar">
    <span class="icon-menu"></span>
  </button>
  <nav class="sidebar-nav" aria-label="Documentation">
    <ul class="sidebar-nav-list">
      <li class="sidebar-nav-item">
        <a href="/docs/" class="sidebar-nav-link">Home</a>
      </li>
      <li class="sidebar-nav-item sidebar-nav-item--has-children">
        <button class="sidebar-nav-toggle" aria-expanded="true">
          Guides
          <span class="icon-chevron"></span>
        </button>
        <ul class="sidebar-nav-children">
          <li><a href="/docs/guides/config/">Configuration</a></li>
          <li><a href="/docs/guides/themes/">Themes</a></li>
        </ul>
      </li>
    </ul>
  </nav>
</aside>
```

#### Auto-Generation

If no manual navigation is specified, the sidebar is auto-generated from:

1. **Feed structure** - Posts in the same feed become nav items
2. **Directory structure** - Nested directories become nested nav items
3. **Frontmatter order** - Use `nav_order: 1` to control ordering
4. **Frontmatter groups** - Use `nav_group: "Guides"` to group items

```yaml
---
title: "Configuration Guide"
nav_order: 2
nav_group: "Guides"
---
```

### TOC Component

The table of contents component shows an outline of the current page.

#### Structure

```html
<aside class="layout-toc" data-position="right">
  <nav class="toc" aria-label="Table of Contents">
    <h2 class="toc-title">On this page</h2>
    <ul class="toc-list">
      <li class="toc-item toc-item--level-2">
        <a href="#overview" class="toc-link">Overview</a>
      </li>
      <li class="toc-item toc-item--level-2">
        <a href="#installation" class="toc-link toc-link--active">Installation</a>
        <ul class="toc-list toc-list--nested">
          <li class="toc-item toc-item--level-3">
            <a href="#prerequisites" class="toc-link">Prerequisites</a>
          </li>
        </ul>
      </li>
    </ul>
  </nav>
</aside>
```

#### Scroll Spy

When `scroll_spy: true`, the TOC highlights the currently visible section.

- The active entry is the last heading whose top has passed the reading line
  (the sticky header offset plus `min(120px, 20vh)`). At the very bottom of the
  page, the last heading still on screen becomes active.
- The active link gets `.toc-link--active` and `aria-current="location"`; its
  ancestors get `.toc-link--parent-active`. A `.toc-marker` element inside the
  top-level `.toc-list` slides to the active link, and the TOC scroller keeps
  that link in view.
- Clicking a TOC link MUST land its heading at the reading offset. The click
  marks the entry active at once and the page moves: first a cut to within
  0.6 viewport heights of the target, then an ease-out glide of 150–250ms that
  re-measures the target every frame. Late layout shifts are corrected at
  +120ms and +450ms. Any wheel, touch or navigation key input cancels the glide
  and the corrections. `prefers-reduced-motion: reduce` jumps instantly.
- The URL hash is updated with `history.pushState`, the heading receives focus
  (`tabindex="-1"`, no scroll), and `.toc-target` plays a brief highlight.
- The clicked entry stays active until the reader scrolls, even when the page
  cannot scroll far enough to bring the heading to the reading line.
- Listeners are bound once and re-initialised after SPA navigation.

#### TOC ids

TOC hrefs MUST equal the ids of the rendered headings:

- `#` lines inside fenced code blocks are not headings.
- Explicit `{#id}` ids win and do not consume an auto-id collision slot.
- Ids are slugged from the visible text: inline code is kept literally
  (`update_meta` keeps its underscore), paired emphasis delimiters are
  removed, and typographer entities (smart quotes) are unescaped first.
- Duplicate auto ids use `-1`, `-2`, … suffixes (GitHub style), counting
  every heading, including levels outside the TOC range.

### Sidebar Drawers (default theme)

The feed sidebar (`.feed-sidebar`) and document sidebar (`.doc-sidebar`) are
fixed drawers at every width. Each has a `.sidebar-toggle` handle with
`aria-expanded`, `aria-keyshortcuts` and `data-sidebar-hotkey` (`b` for the
left drawer, `Shift+B` for the right) and a `<kbd class="sidebar-toggle__kbd">`
hint shown on hover or keyboard focus (hidden on coarse, hover-less pointers).
The handle's `title`/`aria-label` also names the hotkey. Open drawers carry
`.sidebar--pinned`; closed drawers are `visibility: hidden` so their links are
out of the tab order. Each drawer is one scroller on a single surface color
with a themed scrollbar.

- **Docked (≥ 1201px):** opening a drawer pushes `.main-content` aside via
  `.page-wrapper:has(> .X--side.sidebar--pinned)`. Content width is
  `min(page width, 100% − open drawer widths)`, so the article never sits
  under a drawer. Widths come from `--feed-sidebar-width` and
  `--doc-sidebar-width` on `.page-wrapper`. The open state persists per side
  in `localStorage`.
- **Overlay (≤ 1200px):** drawers open over the page (z-index 1100) above a
  scrim (`.page-wrapper::before`). Only one drawer is open at a time. Handles
  are bottom-corner pills, offset right when the dev settings button is
  present. A scrim click, `Escape`, a link click inside the drawer, or the
  hotkey closes it, and focus returns to the handle. Opening moves focus into
  the drawer. Overlay state is never persisted, and entering overlay mode
  closes any open drawers.
- Print hides drawers and handles; reduced motion removes the slide.

### Hover Cards (default theme)

`tooltips.js` is the single hover/focus preview component for the default
theme. It exposes `window.MarkataHoverCards` with `register(selector, build,
options)`, `placement(target, width, height)`, `show(el)`, `hide()`, and
`refresh()`. `build(el)` returns a DOM node,
or `null` to show no card. Built-in sources are tried first, in this order:

| Source | Card content | Extra class |
|--------|--------------|-------------|
| `a.wikilink[data-title]` | title, description (falls back to `data-preview`), date, target path | `wikilink-tooltip` |
| `a.glossary-term` | `data-hover-title` (or link text), definition from `data-hover-description` or `title`, "Glossary" label, path | `hover-card--glossary` |
| `a[data-link-preview]` | optional image, `data-link-title` (or host), description, icon, site name (or host), path | `hover-card--external` |
| `[data-hover-card]` | children of the element or `<template>` matched by the attribute's CSS selector, with `id`s removed | `hover-card--rich` |
| `[data-hover-title]`, `[data-hover-description]` | title and description (falls back to `title`) | none |

Rules:

- Every card is a `div.hover-card[role=tooltip]` appended to `<body>` with
  `position: fixed`. The target gets `aria-describedby` while it's open. Only
  one card is open at a time.
- Built-in text sources set content with `textContent`, never HTML. Rich cards
  clone author-controlled DOM from the page.
- When a card opens for an element with a `title`, the title moves to
  `data-hover-native-title` so the browser tooltip doesn't also appear.
- Pointer hover (not touch) opens after 180ms, or immediately when moving
  from one open card's target to another. Keyboard focus opens only for
  `:focus-visible`. Leaving the target and card closes after 150ms. Moving
  onto the card or focusing inside it keeps it open.
- It closes on `Escape` (returning focus to the target if focus was inside the
  card), window scroll, or a pointer down outside the card.
- Placement follows the pointer. For a target that wraps across lines,
  hovering the first line opens the card above the first line, hovering the
  last line opens it below the last line, and a middle line picks the nearer
  end. Single-line targets and keyboard focus open below. When the preferred
  side is too short for the card and the other side has more room, it flips.
  The card aligns with the left edge of the line it attaches to, clamped to
  8px from the viewport edges, and never overlaps any line of the target when
  either side fits.
- `MarkataHoverCards.placement(target, width, height)` returns
  `{ top, left, placement }` in viewport pixels using these rules. Other theme
  popups (`mention-cards.js`, `tag-cards.js`) use it when it's available, so
  every popup opens on the same side.
- `[data-hover-card]`, `[data-hover-title]`, and `[data-hover-description]`
  elements that aren't focusable get `tabindex="0"` on load and after view
  transitions.
- Events are delegated from `document`, so cards work after SPA navigation
  without re-binding. `window.initTooltips` and `view-transition-complete`
  remove stale cards.

### Header Component

#### Styles

**Full Header:**
```html
<header class="layout-header layout-header--full">
  <div class="header-container">
    <a href="/" class="header-logo">
      <img src="/logo.svg" alt="Site Logo">
      <span class="header-title">Site Title</span>
    </a>
    <nav class="header-nav" aria-label="Main navigation">
      <a href="/docs/">Docs</a>
      <a href="/blog/">Blog</a>
      <a href="/about/">About</a>
    </nav>
    <div class="header-actions">
      <button class="search-toggle" aria-label="Search">
        <span class="icon-search"></span>
      </button>
      <button class="theme-toggle" aria-label="Toggle theme">
        <span class="icon-sun"></span>
        <span class="icon-moon"></span>
      </button>
    </div>
  </div>
</header>
```

**Minimal Header:**
```html
<header class="layout-header layout-header--minimal">
  <div class="header-container">
    <a href="/" class="header-logo">Site</a>
    <div class="header-actions">
      <button class="search-toggle"></button>
      <button class="theme-toggle"></button>
    </div>
  </div>
</header>
```

**Transparent Header:**
```html
<header class="layout-header layout-header--transparent">
  <!-- Same as full, but with transparent background -->
</header>
```

### Footer Component

#### Styles

**Full Footer:**
```html
<footer class="layout-footer layout-footer--full">
  <div class="footer-container">
    <div class="footer-nav">
      <div class="footer-nav-group">
        <h3>Product</h3>
        <a href="/features/">Features</a>
        <a href="/pricing/">Pricing</a>
      </div>
      <div class="footer-nav-group">
        <h3>Resources</h3>
        <a href="/docs/">Documentation</a>
        <a href="/blog/">Blog</a>
      </div>
    </div>
    <div class="footer-social">
      <a href="https://twitter.com/..." aria-label="Twitter">
        <span class="icon-twitter"></span>
      </a>
      <a href="https://github.com/..." aria-label="GitHub">
        <span class="icon-github"></span>
      </a>
    </div>
    <div class="footer-legal">
      <p>&copy; 2024 My Company. All rights reserved.</p>
      <a href="/privacy/">Privacy</a>
      <a href="/terms/">Terms</a>
    </div>
  </div>
</footer>
```

**Minimal Footer:**
```html
<footer class="layout-footer layout-footer--minimal">
  <div class="footer-container">
    <p>&copy; 2024 My Company</p>
  </div>
</footer>
```

### Share Component

The share component provides a configurable "Share this post" experience at the end of article templates. It renders icon buttons for each platform, exposes hover labels, and keeps copy-to-clipboard functionality accessible via keyboard and screen readers.

**Placement**: Injected into `post.html` (and theme equivalents) between the article body and ancillary sections (guide navigation, webmentions) so it appears at the end of every article before comments/footers.

**Behavior**:

1. Displays icon buttons in a responsive grid. Mobile stacks vertically, desktop shows multi-column layout.
2. Hover or focus reveals the platform name via visible text and `aria-label` attributes.
3. Clicks open the platform share dialog in a new tab (or copy the link when `copy` is enabled).
4. Copy button uses the Clipboard API with a DOM fallback and provides live feedback (including a desktop-only toast).
5. Uses CSS variables for colors, spacing, and border radius so palettes can theme the component.

#### Configuration

```toml
[markata-go.components.share]
enabled = true
position = "bottom"            # Controls the style hook (CSS adds `share-panel--bottom`).
title = "Share this post"     # Heading above the buttons.
platforms = ["twitter", "bluesky", "linkedin", "whatsapp", "facebook", "telegram", "pinterest", "reddit", "hacker_news", "email", "copy"]

[markata-go.components.share.custom]
mastodon = { name = "Mastodon", icon = "mastodon.svg", url = "https://mastodon.social/share?text={{title}}&url={{url}}" }
```

| Option | Type | Description |
| --- | --- | --- |
| `enabled` | `bool` | Flip the entire component. Defaults to `true`. |
| `position` | `string` | Adds a modifier class `share-panel--<position>` (default `bottom`). |
| `title` | `string` | Heading text shown above the buttons.
| `platforms` | `[]string` | Ordered list of platform keys to render. Missing list falls back to the built-in order. |
| `custom` | `table` | Keyed definitions (`name`, `icon`, `url`). Icon paths are resolved via `theme_asset_hashed` when they do not start with `/`, `http`, or `data:`.

Valid placeholders in share URLs:

| Placeholder | Description |
| --- | --- |
| `{{title}}` | URL-encoded post title (falls back to site title). |
| `{{url}}` | URL-encoded absolute post URL (`config.url` + `post.href`). |
| `{{excerpt}}` | URL-encoded post description/excerpt when provided.

The component ships with these built-in platforms:

| Key | Template | Notes |
| --- | --- | --- |
| `twitter` | `https://twitter.com/intent/tweet?text={{title}}&url={{url}}` | Icon: `icons/share/twitter.svg`. |
| `bluesky` | `https://bsky.app/intent/compose?text={{url}}` | Icon: `icons/share/bluesky.svg`. |
| `linkedin` | `https://www.linkedin.com/sharing/share-offsite/?url={{url}}` | Icon: `icons/share/linkedin.svg`. |
| `whatsapp` | `https://wa.me/?text={{url}}` | Icon: `icons/share/whatsapp.svg`. |
| `facebook` | `https://www.facebook.com/sharer/sharer.php?u={{url}}` | Icon: `icons/share/facebook.svg`. |
| `telegram` | `https://t.me/share/url?url={{url}}` | Icon: `icons/share/telegram.svg`. |
| `pinterest` | `https://pinterest.com/pin/create/button/?url={{url}}` | Icon: `icons/share/pinterest.svg`. |
| `reddit` | `https://reddit.com/submit?url={{url}}&title={{title}}` | Icon: `icons/share/reddit.svg`. |
| `hacker_news` | `https://news.ycombinator.com/submitlink?u={{url}}&t={{title}}` | Icon: `icons/share/hacker_news.svg`. |
| `email` | `mailto:?subject={{title}}&body={{url}}` | Icon: `icons/share/email.svg`. |
| `copy` | Copy clipboard | Icon: `icons/share/copy.svg`. |

Custom entries can reuse the built-in keys (e.g., override the icon for `copy`) or introduce new platforms (`mastodon`, `bluesky`, etc.). If the `platforms` list omits a built-in key, that button is hidden.

#### Accessibility & Interaction

- Buttons are focusable, provide `aria-label`s like "Share on Twitter" (copy button reads "Copy link to clipboard"), and update their label to "Link copied" after copying.
- The component loads a tiny Clipboard helper script only once.
- Responsive states stack items on narrow viewports and transition with `transform`/`opacity` to feel deliberate.


### Post Copy Control

The post copy control provides a compact "Copy this post" menu near the top of article templates so readers can copy a canonical URL, rich HTML, markdown, or plain text without selecting the whole document manually.

**Placement**: Injected into `post.html` (and theme equivalents) inside the article header after description/metadata so it is visible before the article body.

**Configuration**: Set `enabled = false` under `[markata-go.components.post_copy]` to omit the control from rendered post pages. The default is `true`; omitting the key preserves the existing menu. This setting does not disable the separate share component or post-format outputs. An explicitly disabled control must not expose its clipboard payload in the page HTML, including when a custom post template includes the component.

**Behavior**:

1. Renders a disclosure button labeled "Copy this post" with four actions: `URL`, `Rich`, `Markdown`, and `Text`.
2. `URL` copies the canonical post URL (`config.url` + `post.href`).
3. `Rich` writes both `text/html` and `text/plain` clipboard payloads when the browser supports `ClipboardItem`, using the article title plus rendered article content so pasting into chat tools such as Slack or Teams preserves structure.
4. `Markdown` copies a markdown-friendly payload containing the post title, source URL, and raw markdown body (`post.content`).
5. `Text` copies a plain-text payload containing the title, URL, and text derived from the rendered article body.
6. When rich clipboard APIs are unavailable, the control falls back to plain-text clipboard writes and still reports success/failure through an `aria-live` status message.

**Implementation notes**:

- The control is automatic on post pages unless disabled in site configuration.
- Rich copy must rewrite relative links/media URLs in the copied fragment to absolute URLs before writing HTML to the clipboard.
- The menu closes after a successful copy and exposes short-lived visual feedback for the chosen action.


---

## CSS Structure

### CSS Variables

Layout system uses CSS custom properties for easy theming:

```css
:root {
  /* Layout dimensions */
  --layout-sidebar-width: 280px;
  --layout-toc-width: 220px;
  --layout-content-max-width: 800px;
  --layout-header-height: 64px;
  --layout-footer-min-height: 200px;

  /* Responsive breakpoints */
  --breakpoint-sm: 640px;
  --breakpoint-md: 768px;
  --breakpoint-lg: 1024px;
  --breakpoint-xl: 1280px;

  /* Spacing */
  --layout-gutter: var(--space-6);
  --layout-padding: var(--space-4);

  /* Z-index layers */
  --z-sidebar: 100;
  --z-header: 200;
  --z-overlay: 300;
  --z-modal: 400;

  /* Transitions */
  --transition-sidebar: transform 0.3s ease-in-out;
  --transition-toc: opacity 0.2s ease;
}
```

### Layout Classes

```css
/* Base layout container */
.layout {
  display: grid;
  min-height: 100vh;
  grid-template-rows: auto 1fr auto;
  grid-template-columns: 1fr;
}

/* Docs layout: 3-column */
.layout--docs {
  grid-template-columns: var(--layout-sidebar-width) 1fr var(--layout-toc-width);
  grid-template-areas:
    "header  header  header"
    "sidebar content toc"
    "footer  footer  footer";
}

/* Blog layout: single column */
.layout--blog {
  grid-template-columns: 1fr;
  grid-template-areas:
    "header"
    "content"
    "footer";
}

/* Landing layout: full width */
.layout--landing {
  grid-template-columns: 1fr;
  grid-template-areas:
    "header"
    "content"
    "footer";
}

/* Bare layout: content only */
.layout--bare {
  grid-template-columns: 1fr;
  grid-template-rows: 1fr;
  grid-template-areas: "content";
}
```

### Component Classes

```css
/* Header */
.layout-header {
  grid-area: header;
  position: sticky;
  top: 0;
  z-index: var(--z-header);
  background: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border);
}

.layout-header--sticky {
  position: sticky;
}

.layout-header--transparent {
  background: transparent;
  border-bottom: none;
  position: absolute;
  width: 100%;
}

/* Sidebar */
.layout-sidebar {
  grid-area: sidebar;
  position: sticky;
  top: var(--layout-header-height);
  height: calc(100vh - var(--layout-header-height));
  overflow-y: auto;
  padding: var(--layout-padding);
  border-right: 1px solid var(--color-border);
  background: var(--color-bg-secondary);
}

.layout-sidebar[data-position="right"] {
  grid-area: toc;
  border-right: none;
  border-left: 1px solid var(--color-border);
}

/* Content */
.layout-content {
  grid-area: content;
  max-width: var(--layout-content-max-width);
  margin: 0 auto;
  padding: var(--layout-padding) var(--layout-gutter);
  width: 100%;
}

/* TOC */
.layout-toc {
  grid-area: toc;
  position: sticky;
  top: var(--layout-header-height);
  height: calc(100vh - var(--layout-header-height));
  overflow-y: auto;
  padding: var(--layout-padding);
  border-left: 1px solid var(--color-border);
}

/* Footer */
.layout-footer {
  grid-area: footer;
  background: var(--color-bg-secondary);
  border-top: 1px solid var(--color-border);
}
```

---

## Responsive Behavior

### Breakpoints

| Breakpoint | Width | Behavior |
|------------|-------|----------|
| `xl` | >= 1280px | Full 3-panel layout |
| `lg` | >= 1024px | TOC hidden, sidebar visible |
| `md` | >= 768px | Sidebar as overlay |
| `sm` | < 768px | All panels as overlays |

### Mobile Behavior

On mobile devices (< 768px):

1. **Sidebar** - Collapses to hamburger menu
2. **TOC** - Collapses to floating button or hides
3. **Header** - Remains sticky with hamburger toggle
4. **Content** - Expands to full width

```css
/* Responsive layout */
@media (max-width: 1279px) {
  .layout--docs {
    grid-template-columns: var(--layout-sidebar-width) 1fr;
    grid-template-areas:
      "header  header"
      "sidebar content"
      "footer  footer";
  }

  .layout-toc {
    display: none;
  }
}

@media (max-width: 1023px) {
  .layout--docs {
    grid-template-columns: 1fr;
    grid-template-areas:
      "header"
      "content"
      "footer";
  }

  .layout-sidebar {
    position: fixed;
    left: 0;
    top: var(--layout-header-height);
    transform: translateX(-100%);
    transition: var(--transition-sidebar);
    z-index: var(--z-sidebar);
    height: calc(100vh - var(--layout-header-height));
    width: min(var(--layout-sidebar-width), 85vw);
  }

  .layout-sidebar[data-open="true"] {
    transform: translateX(0);
  }
}

@media (max-width: 767px) {
  :root {
    --layout-header-height: 56px;
    --layout-padding: var(--space-3);
    --layout-gutter: var(--space-3);
  }
}
```

### Mobile Navigation

The hamburger menu button toggles the sidebar:

```html
<button class="mobile-menu-toggle"
        aria-label="Toggle navigation"
        aria-expanded="false"
        aria-controls="sidebar">
  <span class="hamburger-line"></span>
  <span class="hamburger-line"></span>
  <span class="hamburger-line"></span>
</button>
```

JavaScript handles the toggle:

```javascript
// Sidebar toggle behavior
const toggle = document.querySelector('.mobile-menu-toggle');
const sidebar = document.querySelector('.layout-sidebar');
const overlay = document.querySelector('.sidebar-overlay');

toggle?.addEventListener('click', () => {
  const isOpen = sidebar.dataset.open === 'true';
  sidebar.dataset.open = !isOpen;
  toggle.setAttribute('aria-expanded', !isOpen);
  document.body.classList.toggle('sidebar-open', !isOpen);
});

// Close on overlay click
overlay?.addEventListener('click', () => {
  sidebar.dataset.open = 'false';
  toggle.setAttribute('aria-expanded', 'false');
  document.body.classList.remove('sidebar-open');
});

// Close on escape key
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && sidebar.dataset.open === 'true') {
    sidebar.dataset.open = 'false';
    toggle.setAttribute('aria-expanded', 'false');
    document.body.classList.remove('sidebar-open');
  }
});
```

### Overlay

When sidebar is open on mobile, an overlay covers the content:

```css
.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  top: var(--layout-header-height);
  background: rgba(0, 0, 0, 0.5);
  z-index: calc(var(--z-sidebar) - 1);
}

.sidebar-open .sidebar-overlay {
  display: block;
}
```

---

## Template Integration

### Base Layout Template

The layout system extends `base.html` with layout-specific structure:

```jinja2
{# templates/layouts/docs.html #}
{% extends "base.html" %}

{% block body_class %}layout layout--docs{% endblock %}

{% block header %}
{% include "components/header.html" with style=layout.header_style %}
{% endblock %}

{% block main %}
<div class="layout-container">
  {% if layout.sidebar.enabled %}
  {% include "components/sidebar.html" %}
  {% endif %}

  <main class="layout-content" id="main-content">
    {% block content %}{% endblock %}
  </main>

  {% if layout.toc.enabled and post.Extra.toc %}
  {% include "components/toc.html" %}
  {% endif %}
</div>
{% endblock %}

{% block footer %}
{% include "components/footer.html" with style=layout.footer_style %}
{% endblock %}

{% block scripts %}
{{ super() }}
<script src="{{ 'js/layout.js' | theme_asset }}" defer></script>
{% endblock %}
```

### Using Layouts in Post Templates

```jinja2
{# templates/post.html #}
{% extends layout_template %}

{% block content %}
<article class="post">
  <header class="post-header">
    <h1>{{ post.title }}</h1>
    {% if layout.show_date and post.date %}
    <time datetime="{{ post.date | atom_date }}">{{ post.date | date:"January 2, 2006" }}</time>
    {% endif %}
    {% if layout.show_tags and post.tags %}
    <div class="tags">
      {% for tag in post.tags %}
      <a href="/tags/{{ tag | slugify }}/">{{ tag }}</a>
      {% endfor %}
    </div>
    {% endif %}
  </header>

  <div class="post-content">
    {{ body | safe }}
  </div>

  {% if layout.show_prev_next %}
  <nav class="post-nav">
    {% if post.prev %}<a href="{{ post.prev.href }}" class="prev">&larr; {{ post.prev.title }}</a>{% endif %}
    {% if post.next %}<a href="{{ post.next.href }}" class="next">{{ post.next.title }} &rarr;</a>{% endif %}
  </nav>
  {% endif %}
</article>
{% endblock %}
```

### Layout Resolution

The template system resolves layouts in this order:

1. **Frontmatter `layout`** - Post specifies layout name
2. **Feed default** - Feed configuration specifies default layout
3. **Global default** - `[markata-go.layout].name` setting
4. **Fallback** - `blog` layout

```python
def get_layout_template(post, feed, config):
    layout_name = (
        post.frontmatter.get('layout') or
        feed.layout or
        config.layout.name or
        'blog'
    )
    return f"layouts/{layout_name}.html"
```

---

## Examples

### Complete Documentation Site

```toml
# markata-go.toml

[markata-go.layout]
name = "docs"

[markata-go.layout.docs]
sidebar_position = "left"
sidebar_width = "280px"
toc_position = "right"
toc_width = "220px"
content_max_width = "800px"
header_style = "minimal"
footer_style = "minimal"

[markata-go.components.sidebar]
enabled = true
collapsible = true

[[markata-go.components.sidebar.nav]]
title = "Getting Started"
href = "/docs/"

[[markata-go.components.sidebar.nav]]
title = "Guides"
children = [
    { title = "Configuration", href = "/docs/guides/configuration/" },
    { title = "Themes", href = "/docs/guides/themes/" },
    { title = "Layouts", href = "/docs/guides/layouts/" },
    { title = "Plugins", href = "/docs/guides/plugins/" },
]

[[markata-go.components.sidebar.nav]]
title = "Reference"
children = [
    { title = "CLI Commands", href = "/docs/reference/cli/" },
    { title = "Plugin API", href = "/docs/reference/plugins/" },
    { title = "Config Options", href = "/docs/reference/config/" },
]

[markata-go.components.toc]
enabled = true
min_depth = 2
max_depth = 4
scroll_spy = true

[markata-go.components.header]
style = "minimal"
sticky = true
show_search = true
show_theme_toggle = true
```

### Blog with Optional TOC

```toml
# markata-go.toml

[markata-go.layout]
name = "blog"

[markata-go.layout.blog]
content_max_width = "720px"
show_toc = false       # Disabled by default
toc_position = "right"
show_author = true
show_date = true
show_tags = true
show_reading_time = true
show_prev_next = true

[markata-go.components.header]
style = "full"
sticky = true

[markata-go.components.footer]
style = "full"
show_copyright = true
show_social_links = true
```

Individual posts can enable TOC:

```yaml
---
title: "Complete Guide to Go Interfaces"
date: 2024-01-15
toc: true    # Enable TOC for this long post
---
```

### Landing Page

```toml
# markata-go.toml for landing page sections

[markata-go.layout]
name = "landing"

[markata-go.layout.landing]
content_max_width = "100%"
header_style = "transparent"
hero_enabled = true
```

```yaml
---
title: "Welcome to Markata-Go"
layout: landing
hero:
  title: "Beautiful static sites in Go"
  subtitle: "Zero-config, blazing fast, infinitely customizable"
  cta_primary:
    text: "Get Started"
    href: "/docs/getting-started/"
  cta_secondary:
    text: "View on GitHub"
    href: "https://github.com/..."
---
```

### Mixed Layouts in One Site

```yaml
# docs/index.md - uses docs layout
---
title: "Documentation"
layout: docs
---

# docs/api/index.md - uses docs layout (inherited)
---
title: "API Reference"
---

# blog/my-post.md - uses blog layout
---
title: "Announcing v2.0"
layout: blog
---

# index.md - uses landing layout
---
title: "Home"
layout: landing
---
```

---

## Accessibility

### ARIA Landmarks

```html
<header role="banner">...</header>
<nav role="navigation" aria-label="Main">...</nav>
<main role="main" id="main-content">...</main>
<aside role="complementary" aria-label="Table of Contents">...</aside>
<footer role="contentinfo">...</footer>
```

### Skip Links

```html
<a href="#main-content" class="skip-link">Skip to main content</a>
<a href="#sidebar" class="skip-link">Skip to navigation</a>
```

### Keyboard Navigation

- **Tab** - Navigate between interactive elements
- **Escape** - Close open menus/overlays
- **Arrow keys** - Navigate within menus
- **Enter/Space** - Activate buttons and links

### Focus Management

```css
/* Visible focus indicators */
:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}

/* Skip link */
.skip-link {
  position: absolute;
  top: -40px;
  left: 0;
  padding: var(--space-2) var(--space-4);
  background: var(--color-bg-primary);
  z-index: 9999;
}

.skip-link:focus {
  top: 0;
}
```

---

## Configuration Reference

### Layout Config Model

```go
type LayoutConfig struct {
    // Layout preset name
    Name string `toml:"name" default:"blog"`

    // Layout-specific settings
    Docs    DocsLayoutConfig    `toml:"docs"`
    Blog    BlogLayoutConfig    `toml:"blog"`
    Landing LandingLayoutConfig `toml:"landing"`
    Bare    BareLayoutConfig    `toml:"bare"`

    // Global defaults
    Defaults LayoutDefaults `toml:"defaults"`
}

type LayoutDefaults struct {
    ContentMaxWidth string `toml:"content_max_width" default:"800px"`
    HeaderSticky    bool   `toml:"header_sticky" default:"true"`
    FooterSticky    bool   `toml:"footer_sticky" default:"false"`
}

type DocsLayoutConfig struct {
    SidebarPosition    string `toml:"sidebar_position" default:"left"`
    SidebarWidth       string `toml:"sidebar_width" default:"280px"`
    SidebarCollapsible bool   `toml:"sidebar_collapsible" default:"true"`
    SidebarDefaultOpen bool   `toml:"sidebar_default_open" default:"true"`
    TocPosition        string `toml:"toc_position" default:"right"`
    TocWidth           string `toml:"toc_width" default:"220px"`
    TocCollapsible     bool   `toml:"toc_collapsible" default:"true"`
    TocDefaultOpen     bool   `toml:"toc_default_open" default:"true"`
    ContentMaxWidth    string `toml:"content_max_width" default:"800px"`
    HeaderStyle        string `toml:"header_style" default:"minimal"`
    FooterStyle        string `toml:"footer_style" default:"minimal"`
}

type BlogLayoutConfig struct {
    ContentMaxWidth string `toml:"content_max_width" default:"720px"`
    ShowToc         bool   `toml:"show_toc" default:"false"`
    TocPosition     string `toml:"toc_position" default:"right"`
    TocWidth        string `toml:"toc_width" default:"200px"`
    HeaderStyle     string `toml:"header_style" default:"full"`
    FooterStyle     string `toml:"footer_style" default:"full"`
    ShowAuthor      bool   `toml:"show_author" default:"true"`
    ShowDate        bool   `toml:"show_date" default:"true"`
    ShowTags        bool   `toml:"show_tags" default:"true"`
    ShowReadingTime bool   `toml:"show_reading_time" default:"true"`
    ShowPrevNext    bool   `toml:"show_prev_next" default:"true"`
}

type LandingLayoutConfig struct {
    ContentMaxWidth string `toml:"content_max_width" default:"100%"`
    HeaderStyle     string `toml:"header_style" default:"transparent"`
    HeaderSticky    bool   `toml:"header_sticky" default:"true"`
    FooterStyle     string `toml:"footer_style" default:"full"`
    HeroEnabled     bool   `toml:"hero_enabled" default:"true"`
}

type BareLayoutConfig struct {
    ContentMaxWidth string `toml:"content_max_width" default:"100%"`
}
```

### Component Config Models

```go
type SidebarConfig struct {
    Enabled     bool           `toml:"enabled" default:"true"`
    Position    string         `toml:"position" default:"left"`
    Width       string         `toml:"width" default:"280px"`
    Collapsible bool           `toml:"collapsible" default:"true"`
    DefaultOpen bool           `toml:"default_open" default:"true"`
    Nav         []NavItem      `toml:"nav"`
}

type NavItem struct {
    Title    string    `toml:"title"`
    Href     string    `toml:"href"`
    Children []NavItem `toml:"children"`
}

type TocConfig struct {
    Enabled     bool   `toml:"enabled" default:"true"`
    Position    string `toml:"position" default:"right"`
    Width       string `toml:"width" default:"220px"`
    MinDepth    int    `toml:"min_depth" default:"2"`
    MaxDepth    int    `toml:"max_depth" default:"4"`
    Title       string `toml:"title" default:"On this page"`
    Collapsible bool   `toml:"collapsible" default:"true"`
    DefaultOpen bool   `toml:"default_open" default:"true"`
    ScrollSpy   bool   `toml:"scroll_spy" default:"true"`
}

type HeaderConfig struct {
    Style           string `toml:"style" default:"full"`
    Sticky          bool   `toml:"sticky" default:"true"`
    ShowLogo        bool   `toml:"show_logo" default:"true"`
    ShowTitle       bool   `toml:"show_title" default:"true"`
    ShowNav         bool   `toml:"show_nav" default:"true"`
    ShowSearch      bool   `toml:"show_search" default:"true"`
    ShowThemeToggle bool   `toml:"show_theme_toggle" default:"true"`
}

type FooterConfig struct {
    Style           string `toml:"style" default:"full"`
    Sticky          bool   `toml:"sticky" default:"false"`
    ShowCopyright   bool   `toml:"show_copyright" default:"true"`
    CopyrightText   string `toml:"copyright_text"`
    ShowSocialLinks bool   `toml:"show_social_links" default:"true"`
    ShowNavLinks    bool   `toml:"show_nav_links" default:"true"`
}
```

---

## See Also

- [THEMES.md](./THEMES.md) - Color palettes and theme customization
- [TEMPLATES.md](./TEMPLATES.md) - Template system and blocks
- [CONFIG.md](./CONFIG.md) - Configuration system
- [HEAD_STYLE.md](./HEAD_STYLE.md) - Head/style injection
