# Themes and Customization Specification

Themes control the visual appearance of the generated site. The system supports:
- **Built-in themes** - Ship with usable defaults
- **Theme packages** - Complete themes with templates, CSS, and assets  
- **User customization** - Override any theme file locally

## Design Principles

1. **Zero-config beautiful** - Sites look good with no customization
2. **Progressive customization** - Pick colors → override CSS → replace templates
3. **Theme packages** - Installable, shareable, complete visual identities
4. **Local overrides** - Any theme file can be overridden by placing it in your project
5. **Readable by default** - Typography, contrast, and spacing optimized for reading
6. **Accessible first** - WCAG 2.1 AA compliant colors, focus states, and compact touch targets

## Compact Header Preferences

When the default theme's full palette switcher is enabled, the header MUST show
a single compact preferences button beside the search input, instead of
separate mode, palette, and loading buttons. The popover MUST continue to
provide light/dark color mode, palettes, aesthetics, fonts, reading size, and
adaptive media-loading preferences (Auto, Save Data, Full Quality).

The loading control MUST retain its local persistence, automatic recommendations,
and mismatch/action announcements. The existing `t` theme-picker and `\\`
mode-toggle shortcuts MUST remain available. When the palette manifest is
missing or the full switcher is disabled, a standalone loading control MUST
remain accessible and the dark/light fallback SHOULD remain visible where
configured. Keyboard focus and Escape MUST work within the popover, including
its nested loading inputs, without triggering theme-grid shortcuts.

The navbar search MUST remain accessible and keyboard operable with either
Pagefind or Bleve, and remain lazy-loaded where it was before. CSS changes
SHOULD use existing palette tokens and MUST NOT add extra network resources.
Core styles MUST reserve the navbar search input height and mobile spacing
before deferred search styles or either search renderer hydrate the placeholder.
Loading search styles MUST NOT move the header or following content. Editorial
feed header spacing MUST remain stable when deferred calendar styles load.

## Reader Mode

On post pages, `s` toggles reader mode; on dedicated feed pages it still switches
between simple and rich feeds. The `?` shortcut help MUST describe both contexts.
Reader mode is also restored from `?reader=1` on a post URL and can be exited
with `s`. It MUST hide navigation chrome and sidebars, including pinned drawer
handles and their dimming scrims, while keeping the article and its text at full
contrast. Its exit hint MUST NOT inherit the theme's decorative `body::after`
overlay, background image, or mask; the hint text MUST remain readable when a
motif is enabled. The post layout MUST center the article without
reserving sidebar columns or reducing its available width below the normal
reading measure. It MUST fit narrow screens without horizontal overflow. Other
pages and the normal post layout are unchanged.

## Reading Size Preferences

The default theme provides a visitor-controlled reading-size preference for
long-form content. The default setting is `large`, which uses an 18px site
base, 22px article text, and a 64ch article measure. The site base remains smaller than article text so navigation and
controls do not become unnecessarily large.

Sites configure the default and whether the control is visible:

```toml
[markata-go.theme]
text_size = "large"              # small, medium, large, or x-large
show_text_size_control = true    # show the visitor control (default: true)
```

The presets are:

| Preset | Site base | Article text | Article measure |
|--------|-----------|--------------|-----------------|
| `small` | 16px | 18px | 68ch |
| `medium` | 17px | 20px | 66ch |
| `large` | 18px | 22px | 64ch |
| `x-large` | 19px | 24px | 62ch |

Article text also scales with the viewport so wide desktop displays do not
render a narrow strip of small type: `--reading-scale` multiplies the article
font size by 1.08 from 1800px, 1.16 from 2200px (1440p), and 1.25 from
3000px (4K). The `ch`-based measure follows the scaled font, so line length
stays comfortable. Site chrome (`--text-base`) is not scaled.

`--content-width` is expressed in `ch` and resolves against the article font
size, so every preset lands near 65-70 characters per line. The article card
(`article.post`) uses `box-sizing: content-box` so its padding frames the
measure instead of eating into it.

The configured `text_size` is the fallback used when a visitor has no saved
preference. When the control is enabled, the default theme renders an
accessible Small/Medium/Large/X-Large select in the header. A visitor's choice is
stored in browser-local storage under `text-size` and takes precedence over
the site default on later visits. Browser zoom remains independent and is
always supported.

Custom themes may reuse the `data-text-size` attribute and the typography
custom properties, or omit the control while retaining the configured
default. Invalid values fall back to `large` during rendering.

Preset values MUST be emitted below the `overrides` cascade layer so that
explicit `[theme.variables]` entries (which are emitted in `overrides`) take
precedence over any preset for the same token. Presets MUST still take
precedence over the base `tokens` layer.

### Reading Rhythm

The default theme sets prose rhythm in `em` units relative to the article
text size so spacing scales with the preset:

- Body copy uses `--leading-prose` (1.7) and `margin-block: 0 1.25em`.
- Headings inside `.post-content` carry `2em` above and `0.6em` below, use
  `text-wrap: balance`, and are followed by a short accent bar
  (`h1::after`/`h2::after`, colored by `--heading-rule`) instead of full-width
  rules, so section breaks read as gentle markers rather than hard lines.
- Blockquotes sit on a soft surface with a 3px accent left border; `hr` renders
  as a centered gradient hairline.

### Surfaces and Borders

Borders MUST be soft by default. The palette CSS generator emits:

- `--color-border`: `color-mix(in srgb, <ink> 16%, <background>)` — the
  everyday hairline used by cards, headers, tables and form controls.
- `--color-border-strong`: the palette's full-contrast border ink, reserved for
  focus rings and places that need a 3:1 non-text contrast guarantee.
- `--color-border-soft`: `--color-border` at 60% alpha, used for section
  dividers and post header/footer rules.

Radii follow a four-step scale — `--radius-sm` (0.375rem), `--radius`
(0.5rem), `--radius-lg` (0.75rem), `--radius-xl` (1rem) — and aesthetics MAY
remap the scale (for example the `minimal` aesthetic tightens it). Tags render
as pills; cards and the article container use `--radius-lg`/`--radius-xl`.

### Contrast Guarantees

- Default theme text in compact UI surfaces such as home metadata, card metadata, and admonition titles MUST use text tokens that maintain WCAG 2.1 AA contrast at their rendered size.
- Generic container classes such as `warning`, `note`, `info`, `success`, `error`, and `danger` MUST NOT override the palette-specific surface or foreground styles of `.admonition` elements that share those class names. Admonition body text and titles MUST maintain WCAG 2.1 AA contrast against their rendered palette surface.
- Built-in palette regression tests MUST validate the default theme's concrete foreground/background combinations, not only abstract palette token pairs, so regressions in selector usage are caught across every shipped palette.

---

## Customization Philosophy

Markata-go is designed to make customization **easy and intuitive**. Users should be able to:

### Level 1: Zero Effort (Instant Results)
- Choose a built-in color palette in one line of config
- Switch between light/dark modes automatically
- Get beautiful defaults that just work

```toml
# One line to change your entire site's look
[markata-go.theme]
palette = "catppuccin-mocha"
```

### Level 2: Quick Tweaks (Minutes)
- Override specific colors without touching CSS
- Change fonts, spacing, or accent colors
- Toggle features on/off

```toml
[markata-go.theme]
palette = "nord-dark"

[markata-go.theme.palette]
accent = "#88c0d0"  # Override just the accent color

[markata-go.theme.features]
toc = false         # Disable table of contents
```

### Level 3: Template Overrides (Hours)
- Override a single template file (e.g., footer)
- Create custom layouts for specific content types
- Maintain upgradability by only changing what you need

```
my-site/
├── templates/
│   └── partials/
│       └── footer.html    # Your custom footer only
└── markata-go.toml
```

### Level 4: Full Custom Theme (Days)
- Create a complete custom theme
- Share themes as packages
- Extend existing themes

### Design Goals

| Goal | How We Achieve It |
|------|-------------------|
| **Look good immediately** | Carefully designed default theme with professional typography |
| **Easy color changes** | Palette system with 20+ built-in options |
| **No CSS knowledge needed** | Config-driven customization for common changes |
| **CSS power when needed** | Full CSS override capability for advanced users |
| **Don't break on upgrade** | Override system means your customizations survive updates |
| **Consistent across formats** | Same visual identity in HTML, RSS readers, and more |

---

## Layout System

Markata-go uses a hierarchical layout system with different templates for different content types.

### Layout Hierarchy

```
┌─────────────────────────────────────────────────────────────────────┐
│                          base.html                                   │
│  (HTML skeleton, <head>, header, footer, common CSS/JS)             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐     │
│  │   post.html     │  │   feed.html     │  │   page.html     │     │
│  │  Single article │  │  List of posts  │  │  Static pages   │     │
│  │  (blog posts)   │  │  (index, tags)  │  │  (about, etc)   │     │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘     │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    Custom Layouts                            │   │
│  │  landing.html, gallery.html, docs.html, etc.                │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### Built-in Layouts

| Layout | Purpose | Used For |
|--------|---------|----------|
| `base.html` | HTML skeleton with blocks | All pages inherit from this |
| `post.html` | Single content item | Blog posts, articles |
| `feed.html` | List of content items | Index pages, tag pages, archives |
| `page.html` | Static pages | About, contact, custom pages |
| `card.html` | Post preview in lists | Used inside feed.html |

### Post vs Feed Layouts

**Post Layout (`post.html`)** - Single content item with full content:
- Article header (title, date, author)
- Full rendered content
- Tags/categories
- Previous/next navigation
- Comments section (optional)
- Related posts (optional)

**Feed Layout (`feed.html`)** - Collection of content items:
- Feed header (title, description)
- List of post cards/previews
- Pagination controls
- Filter/sort options (optional)

### Per-Post Layout Override

Individual posts can specify a custom layout in frontmatter:

```yaml
---
title: "My Product Landing Page"
layout: landing
---
```

This uses `templates/landing.html` instead of the default `post.html`.

**Use cases:**
- Landing pages with custom hero sections
- Gallery posts with image grids
- Documentation pages with sidebars
- Minimal pages without header/footer

### Per-Post Template Override

For complete control, posts can specify a full template:

```yaml
---
title: "Custom Page"
template: custom/special-page.html
---
```

**Layout vs Template:**
- `layout` - Uses a layout that extends `base.html` (recommended)
- `template` - Uses a completely standalone template (full control)

### Per-Post Style Override

Posts can include custom CSS or override palette colors:

```yaml
---
title: "Dark Mode Article"
palette: dracula
custom_css: |
  .post-content {
    font-size: 1.1rem;
  }
---
```

Or reference an external CSS file:

```yaml
---
title: "Styled Article"  
css: /css/special-article.css
---
```

### Per-Post Feature Flags

Toggle theme features for specific posts:

```yaml
---
title: "Minimal Post"
features:
  toc: false
  reading_time: false
  comments: false
  header: false
  footer: false
---
```

---

## RSS and Atom Stylesheets

RSS and Atom feeds SHOULD include XSL stylesheets for beautiful rendering in browsers.

### Why XSL Stylesheets?

When users click an RSS/Atom link, browsers display raw XML by default - ugly and confusing. XSL stylesheets transform this into a readable, styled page that:

1. **Explains what RSS is** - Many users don't know what they're looking at
2. **Shows feed content beautifully** - Styled preview of recent posts
3. **Provides subscription instructions** - How to use the feed
4. **Maintains brand consistency** - Uses your site's colors and typography

### Feed Stylesheet Structure

```
themes/default/
├── templates/
│   ├── rss.xml              # RSS 2.0 feed template
│   ├── atom.xml             # Atom feed template
│   └── feed-styles.xsl      # Shared XSL stylesheet
└── static/
    └── css/
        └── feed.css         # CSS for styled feed pages
```

### XSL Stylesheet Template

```xml
<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:output method="html" encoding="UTF-8" />

  <xsl:template match="/">
    <html>
      <head>
        <title><xsl:value-of select="/rss/channel/title" /> - RSS Feed</title>
        <style>
          /* Inline styles using palette colors */
          :root {
            --color-bg: {{ palette.bg_primary }};
            --color-text: {{ palette.text_primary }};
            --color-accent: {{ palette.accent }};
          }
          body {
            font-family: system-ui, sans-serif;
            max-width: 50rem;
            margin: 2rem auto;
            padding: 0 1rem;
            background: var(--color-bg);
            color: var(--color-text);
          }
          /* ... more styles ... */
        </style>
      </head>
      <body>
        <header>
          <h1>📡 <xsl:value-of select="/rss/channel/title" /></h1>
          <p class="subtitle">This is an RSS feed. Subscribe by copying the URL into your feed reader.</p>
          <p class="feed-url"><xsl:value-of select="/rss/channel/link" /></p>
        </header>

        <section class="about-feeds">
          <h2>What is an RSS feed?</h2>
          <p>RSS feeds allow you to subscribe to websites and receive updates automatically in a feed reader app.</p>
          <details>
            <summary>Popular feed readers</summary>
            <ul>
              <li><a href="https://feedly.com">Feedly</a> (Web, iOS, Android)</li>
              <li><a href="https://netnewswire.com">NetNewsWire</a> (Mac, iOS)</li>
              <li><a href="https://newsblur.com">NewsBlur</a> (Web, iOS, Android)</li>
            </ul>
          </details>
        </section>

        <section class="recent-posts">
          <h2>Recent Posts</h2>
          <xsl:for-each select="/rss/channel/item">
            <article>
              <h3><a href="{link}"><xsl:value-of select="title" /></a></h3>
              <time><xsl:value-of select="pubDate" /></time>
              <p><xsl:value-of select="description" /></p>
            </article>
          </xsl:for-each>
        </section>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>
```

### RSS Template with Stylesheet Reference

```xml
<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="/feed-styles.xsl"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ config.title }}</title>
    <link>{{ config.url }}</link>
    <description>{{ config.description }}</description>
    <atom:link href="{{ feed.url }}/rss.xml" rel="self" type="application/rss+xml"/>
    {% for post in feed.posts %}
    <item>
      <title>{{ post.title }}</title>
      <link>{{ post.absolute_url }}</link>
      <guid isPermaLink="true">{{ post.absolute_url }}</guid>
      <pubDate>{{ post.date | rfc822 }}</pubDate>
      <description>{{ post.description | escape }}</description>
    </item>
    {% endfor %}
  </channel>
</rss>
```

### Atom Template with Stylesheet

```xml
<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="/feed-styles.xsl"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>{{ config.title }}</title>
  <link href="{{ config.url }}"/>
  <link href="{{ feed.url }}/atom.xml" rel="self"/>
  <updated>{{ feed.updated | iso8601 }}</updated>
  <id>{{ config.url }}/</id>
  <author>
    <name>{{ config.author }}</name>
  </author>
  {% for post in feed.posts %}
  <entry>
    <title>{{ post.title }}</title>
    <link href="{{ post.absolute_url }}"/>
    <id>{{ post.absolute_url }}</id>
    <updated>{{ post.date | iso8601 }}</updated>
    <summary>{{ post.description | escape }}</summary>
  </entry>
  {% endfor %}
</feed>
```

### Feed Stylesheet Configuration

```toml
[markata-go.feeds]
# Enable styled RSS/Atom feeds (default: true)
xsl_stylesheet = true

# Custom stylesheet path (optional)
xsl_template = "templates/my-feed-styles.xsl"
```

### Feed Stylesheet Theming

The XSL stylesheet SHOULD use the same color palette as the main site:

```toml
[markata-go.theme]
palette = "catppuccin-mocha"  # RSS/Atom pages use the same palette, fonts and aesthetic
```

The default `rss.xsl` and `atom.xsl` MUST render with the same theme as the site pages:

- Each stylesheet contains a marker region:

  ```xml
  <!-- markata:theme-head -->
  <link rel="stylesheet" href="/css/variables.css" />
  <link rel="stylesheet" href="/css/main.css" />
  <!-- /markata:theme-head -->
  ```

- When copying stylesheets to the output, `publish_feeds` renders
  `partials/feed-xsl-head.html` with the site config and replaces the whole
  region with it. The partial sets the `<html>` theme attributes
  (`data-palette` resolved for `fallback_mode`, `data-aesthetic`,
  `data-fontpack`, `data-text-size`, `data-theme-texture-scope`), includes
  `partials/theme-bootstrap.html` (the same color-mode, visitor theme-picker,
  seasonal-palette and text-size bootstrap `base.html` uses), and links the same
  hashed core CSS (variables, palette, main, components, text-size, fonts,
  aesthetic).
- The rendered markup is made XML-safe before substitution: `<script>` bodies are
  wrapped in CDATA and `<link>`/`<meta>` are self-closed.
- Stylesheets without the marker (for example older custom files in the
  templates directory) are copied unchanged. If rendering fails, the region is
  left as-is so the static fallback links still apply.

The default browser reading view MUST:

- identify the document as a web feed and explain how to subscribe;
- show the feed title, description, entry count, and update date when present;
- render entry titles as links, with publication date and summary when present;
- preserve the RSS/Atom document as machine-readable XML; the XSL only controls
  the browser's HTML presentation;
- use responsive layout, visible keyboard focus, and reduced-motion-aware
  animation while inheriting the site theme through the marker region above.

### Theme Migration Warnings

Canonical nested theme values take precedence over legacy flat values.
The loader MUST report each conflict between explicitly supplied values once.
Built-in defaults MUST NOT cause migration warnings.
This rule applies to file loading, includes, overrides, preview overlays, string loading, and configuration-free loading.
Default merging MUST preserve the resolved theme values.

### Themed 404 Page

The generated `404.html` MUST use the site's post page shell and active theme,
identify the missing page clearly, and retain the requested path, search form,
suggestions, and a direct route back to the site. Its layout MUST remain usable
on narrow screens and provide visible keyboard focus. Both the root template set
and the packaged default theme include the built-in 404 template; site templates
may override it with `custom_404_template`.

---

## Integration with Head/Style System

Themes provide the base visual identity. For additional customization, see [HEAD_STYLE.md](./HEAD_STYLE.md):

| System | Purpose | Example |
|--------|---------|---------|
| Theme (`[markata-go.theme]`) | Base templates, CSS, layouts | Selecting "blog" theme |
| Theme Variables (`[markata-go.theme.variables]`) | Quick CSS property overrides | `--color-primary: #8b5cf6` |
| Head Config (`[markata-go.head]`) | Meta tags, scripts, links | Analytics, fonts, OG tags |
| Style Config (`[markata-go.style]`) | Legacy color overrides | `color_bg`, `color_text` |
| Post Overrides (`config_overrides`) | Per-post customization | Custom scripts for one page |

**Resolution order:** Theme defaults → Theme variables → Head/Style config → Post overrides

Theme variables are applied at build time to the generated `variables.css` and override palette-generated values. This means overrides work consistently even when a palette is enabled.

Generated theme-variable declarations MUST use lexicographic property-name
ordering. The same configuration MUST produce identical CSS bytes and asset
hashes across repeated generation, regardless of map iteration order.

---

## Theme Structure

A theme is a directory containing templates, CSS, and optionally assets:

```
themes/
└── default/
    ├── theme.toml           # Theme metadata
    ├── templates/
    │   ├── base.html        # Base layout
    │   ├── post.html        # Single post
    │   ├── feed.html        # Feed/index page
    │   ├── card.html        # Post card for feeds
    │   └── partials/
    │       ├── head.html
    │       ├── header.html
    │       ├── footer.html
    │       └── pagination.html
    ├── static/
    │   ├── css/
    │   │   ├── main.css     # Core styles
    │   │   ├── admonitions.css
    │   │   ├── code.css     # Syntax highlighting
    │   │   └── variables.css # CSS custom properties
    │   └── js/
    │       └── main.js      # Optional JavaScript
    └── assets/              # Theme-specific images, fonts
```

### theme.toml

```toml
[theme]
name = "Default"
version = "1.0.0"
description = "Clean, minimal theme with dark mode support"
author = "SSG Team"
license = "MIT"
homepage = "https://github.com/example/ssg-theme-default"

# Minimum SSG version required
min_version = "1.0.0"

# Theme supports these features
[theme.features]
dark_mode = true
responsive = true
syntax_highlighting = true
admonitions = true

# Configurable options this theme exposes
[theme.options]
primary_color = { type = "color", default = "#3b82f6", description = "Primary accent color" }
font_family = { type = "string", default = "system-ui", description = "Body font family" }
max_width = { type = "string", default = "65ch", description = "Content max width" }
show_toc = { type = "boolean", default = true, description = "Show table of contents" }
show_reading_time = { type = "boolean", default = true, description = "Show reading time" }

# Parent theme to inherit from (optional)
[theme.extends]
theme = "base"  # Inherit from another theme
```

---

## Theme Resolution

When looking for a template or static file, the system searches in order:

```
1. Project local:     ./templates/post.html
2. Project theme:     ./themes/[theme]/templates/post.html  
3. Installed theme:   ~/.config/markata-go/themes/[theme]/templates/post.html
4. Built-in theme:    [internal]/themes/[theme]/templates/post.html
5. Default theme:     [internal]/themes/default/templates/post.html
```

This allows users to override any theme file by creating it locally.

### Override Examples

**Override just the footer:**
```
my-site/
├── templates/
│   └── partials/
│       └── footer.html    # Your custom footer
├── posts/
└── markata-go.toml
```

**Override a CSS file:**
```
my-site/
├── static/
│   └── css/
│       └── admonitions.css  # Your custom admonition styles
├── posts/
└── markata-go.toml
```

---

## Configuration

### Selecting a Theme

```toml
[markata-go.theme]
name = "default"           # Theme name
```

### Customizing Theme Options

Themes expose configurable options:

```toml
[markata-go.theme]
name = "default"

# Theme-specific options (defined in theme.toml)
[markata-go.theme.options]
primary_color = "#8b5cf6"   # Purple instead of blue
font_family = "Inter, system-ui"
max_width = "70ch"
show_toc = false
```

### CSS Variable Overrides

For quick color/font changes without modifying CSS files:

```toml
[markata-go.theme]
name = "default"

# Override CSS custom properties
[markata-go.theme.variables]
"--color-primary" = "#8b5cf6"
"--color-primary-dark" = "#7c3aed"
"--font-body" = "Inter, system-ui"
"--font-heading" = "Inter, system-ui"
"--font-mono" = "JetBrains Mono, monospace"
"--content-width" = "70ch"
"--radius" = "0.5rem"
```

These are injected as a `<style>` block in the `<head>`:

```html
<style>
:root {
  --color-primary: #8b5cf6;
  --color-primary-dark: #7c3aed;
  --font-body: Inter, system-ui;
  /* ... */
}
</style>
```

### Custom CSS File

For more extensive customization:

```toml
[markata-go.theme]
name = "default"
custom_css = "my-styles.css"  # Loaded after theme CSS
```

---

## Built-in Themes

Implementations MUST provide at least the `default` theme.

### Default Theme

A clean, minimal theme with:
- Responsive layout
- Dark mode default with explicit user override
- Styled admonitions
- Syntax highlighting
- Typography optimized for reading

### Minimal Theme (Optional)

Bare-bones HTML with minimal styling for:
- Maximum customization flexibility
- Fast loading
- Print-friendly

---

## Color Palette System

Color palettes are **separate from layout themes**. This separation allows users to:
- Apply any color palette to any layout theme
- Create custom palettes without touching theme structure
- Switch between light/dark variants easily

### Architecture: Three-Layer Color System

```
┌─────────────────────────────────────────────────────────────┐
│  Layer 1: Raw Colors (Palette Definition)                   │
│  Pure color values with no semantic meaning                 │
│  e.g., red = "#f38ba8", blue = "#89b4fa"                   │
├─────────────────────────────────────────────────────────────┤
│  Layer 2: Semantic Colors (Role-Based)                      │
│  Colors mapped to meaning/purpose                           │
│  e.g., text = red, accent = blue, surface = base           │
├─────────────────────────────────────────────────────────────┤
│  Layer 3: Component Colors (Usage-Specific)                 │
│  Fine-grained component styling                             │
│  e.g., button-bg = accent, link-color = accent             │
└─────────────────────────────────────────────────────────────┘
```

### Palette File Format

Palettes are defined in TOML files within `palettes/` directory:

```toml
# palettes/catppuccin-mocha.toml
[palette]
name = "Catppuccin Mocha"
variant = "dark"              # "dark" | "light"
author = "Catppuccin Team"
license = "MIT"
homepage = "https://catppuccin.com"

# Layer 1: Raw Colors
# These are the source colors - pure values with no semantic meaning
[palette.colors]
rosewater = "#f5e0dc"
flamingo  = "#f2cdcd"
pink      = "#f5c2e7"
mauve     = "#cba6f7"
red       = "#f38ba8"
maroon    = "#eba0ac"
peach     = "#fab387"
yellow    = "#f9e2af"
green     = "#a6e3a1"
teal      = "#94e2d5"
sky       = "#89dceb"
sapphire  = "#74c7ec"
blue      = "#89b4fa"
lavender  = "#b4befe"
text      = "#cdd6f4"
subtext1  = "#bac2de"
subtext0  = "#a6adc8"
overlay2  = "#9399b2"
overlay1  = "#7f849c"
overlay0  = "#6c7086"
surface2  = "#585b70"
surface1  = "#45475a"
surface0  = "#313244"
base      = "#1e1e2e"
mantle    = "#181825"
crust     = "#11111b"

# Layer 2: Semantic Mapping
# Maps raw colors to semantic roles
[palette.semantic]
# Text colors
text-primary   = "text"       # Reference to raw color name
text-secondary = "subtext1"
text-muted     = "overlay1"

# Background colors
bg-primary   = "base"
bg-secondary = "mantle"
bg-surface   = "surface0"
bg-elevated  = "surface1"

# Accent colors
accent         = "mauve"
accent-hover   = "lavender"
link           = "blue"
link-hover     = "sapphire"
link-visited   = "mauve"

# Status colors
success = "green"
warning = "yellow"
error   = "red"
info    = "blue"

# Border colors
border         = "surface1"
border-focus   = "mauve"

# Button colors
button-primary-bg     = "mauve"
button-primary-text   = "base"
button-secondary-bg   = "surface1"
button-secondary-text = "text"

# Layer 3: Component Colors (optional overrides)
[palette.components]
# Code blocks
code-bg         = "surface0"
code-text       = "text"
code-comment    = "overlay1"
code-keyword    = "mauve"
code-string     = "green"
code-number     = "peach"
code-function   = "blue"

# Admonitions
admonition-note-bg     = "surface0"
admonition-note-border = "blue"
admonition-tip-bg      = "surface0"
admonition-tip-border  = "green"
admonition-warn-bg     = "surface0"
admonition-warn-border = "yellow"
admonition-error-bg    = "surface0"
admonition-error-border = "red"
```

### Configuration

#### Selecting a Palette

```toml
[markata-go.theme]
name = "default"              # Layout theme
palette = "catppuccin-mocha"  # Color palette

# Or with explicit light-mode override
[markata-go.theme]
name = "default"
palette = "catppuccin-latte"        # Light mode palette
palette_dark = "catppuccin-mocha"   # Dark mode palette
```

#### Inline Palette Customization

Override specific colors without creating a new palette:

```toml
[markata-go.theme]
name = "default"
palette = "catppuccin-mocha"

# Override specific semantic colors
[markata-go.theme.palette]
accent = "#ff79c6"            # Use Dracula's pink as accent
link = "#8be9fd"              # Use Dracula's cyan for links
```

#### Custom Palette Definition

Define a complete custom palette inline:

```toml
[markata-go.theme]
name = "default"

# Full custom palette
[markata-go.theme.palette]
name = "My Brand"
variant = "light"

[markata-go.theme.palette.colors]
brand-primary = "#4f46e5"
brand-secondary = "#7c3aed"
# ... more colors

[markata-go.theme.palette.semantic]
accent = "brand-primary"
link = "brand-secondary"
# ... more mappings
```

### Built-in Palettes

Implementations MUST include at least `default-light` and `default-dark` palettes.
Implementations SHOULD include these popular community palettes (all MIT licensed):

| Palette | Variants | Description |
|---------|----------|-------------|
| `catppuccin-latte` | light | Soothing pastel theme |
| `catppuccin-frappe` | dark | Medium contrast dark |
| `catppuccin-macchiato` | dark | Higher contrast dark |
| `catppuccin-mocha` | dark | Highest contrast dark |
| `nord-light` | light | Arctic, north-bluish |
| `nord-dark` | dark | Arctic, north-bluish |
| `gruvbox-light` | light | Retro groove, warm |
| `gruvbox-dark` | dark | Retro groove, warm |
| `tokyo-night` | dark | Tokyo city lights |
| `tokyo-night-storm` | dark | Stormy variant |
| `tokyo-night-day` | light | Day variant |
| `rose-pine` | dark | Soho vibes, natural |
| `rose-pine-moon` | dark | Softer variant |
| `rose-pine-dawn` | light | Light variant |
| `everforest-light` | light | Nature-inspired green |
| `everforest-dark` | dark | Nature-inspired green |
| `dracula` | dark | Vibrant on dark |
| `solarized-light` | light | Scientifically designed |
| `solarized-dark` | dark | Scientifically designed |
| `kanagawa-wave` | dark | Japanese art inspired |
| `kanagawa-dragon` | dark | Warmer variant |
| `kanagawa-lotus` | light | Light variant |

### Palette Directory Structure

```
palettes/
├── catppuccin-latte.toml
├── catppuccin-frappe.toml
├── catppuccin-macchiato.toml
├── catppuccin-mocha.toml
├── nord-light.toml
├── nord-dark.toml
├── gruvbox-light.toml
├── gruvbox-dark.toml
├── tokyo-night.toml
├── tokyo-night-storm.toml
├── tokyo-night-day.toml
├── rose-pine.toml
├── rose-pine-moon.toml
├── rose-pine-dawn.toml
├── everforest-light.toml
├── everforest-dark.toml
├── dracula.toml
├── solarized-light.toml
├── solarized-dark.toml
├── kanagawa-wave.toml
├── kanagawa-dragon.toml
└── kanagawa-lotus.toml
```

### Palette Resolution Order

1. User's project: `./palettes/{name}.toml`
2. User config: `~/.config/markata-go/palettes/{name}.toml`
3. Built-in: `[internal]/palettes/{name}.toml`

### Default Palette Behavior

If no palette is specified in configuration:
- **Default palette:** `default-light` is used
- **Mode fallback:** If no explicit user preference is saved, markata-go uses `theme.fallback_mode` (default: `"dark"`)
- **User preference wins:** Light mode is only used when the visitor explicitly selects light mode

```toml
[markata-go.theme]
palette = "catppuccin-latte"
palette_dark = "catppuccin-mocha"
fallback_mode = "dark"  # "dark" (default) or "light"
```

```toml
# Explicit light/dark configuration (recommended)
[markata-go.theme]
palette = "catppuccin-latte"        # Light mode
palette_dark = "catppuccin-mocha"   # Dark mode

# Single-variant palette: the other mode is derived automatically
[markata-go.theme]
palette = "dracula"                  # dark: dracula, light: dracula-light (derived)
```

### Derived Counterparts

Every palette MUST resolve to both a light and a dark palette. Variant
resolution (`palettes.GetEffectivePalettes`) proceeds in order:

1. Explicit `palette_light` / `palette_dark` configuration.
2. Known family mappings (`catppuccin-*`, `rose-pine*`, `tokyo-night*`,
   `kanagawa-*`).
3. Standard naming: `<base>-light` / `<base>-dark` files that exist.
4. A derived counterpart named `<palette>-light` or `<palette>-dark`,
   produced by `palettes.DeriveCounterpart` and served by the loader as if
   it were a real palette.

The derivation classifies raw colors by their semantic/component usage:

- **Surfaces** (`bg-*`, `code-bg`, `nav-bg`, `card-*`, `border`,
  admonition backgrounds) keep hue, scale saturation by 0.7, and are
  compressed into a 0.12-wide lightness band anchored at L=0.985 (light) or
  L=0.11 (dark) while preserving their relative order.
- **Text** (`text-*`, `code-text`, `code-comment`, `nav-text`) reflects
  lightness (L → 1−L) and is then adjusted until it meets 7:1 (primary and
  secondary) or 4.5:1 (muted) against every derived surface.
- **Accents** (everything else) reflect lightness, clamp into a readable
  mid band (0.24–0.46 light, 0.56–0.80 dark), and are adjusted to 4.5:1
  against the surfaces.

Semantic and component references are copied unchanged, so a derived palette
has the same structure as its source. `Loader.Discover()` includes derived
counterparts (flagged `Derived: true`) for families that lack an explicit
partner, so the multi-palette switcher offers both modes for every family.
The rendering contract's derived entries are kept in sync with
`go run ./scripts/rendering-contract --derive-palettes`.

Templates pin `data-palette` on `<html>` to the palette matching the resolved
color mode (`config.Extra.palette_light_effective` /
`palette_dark_effective`), and the mode toggle updates it, so per-palette CSS
blocks never disagree with `data-theme`.

### Color Reference Rules

Semantic and component colors follow strict reference rules:

| Layer | Can Reference | Cannot Reference |
|-------|---------------|------------------|
| Raw colors | (literal hex values) | - |
| Semantic colors | Raw colors only | Other semantic colors |
| Component colors | Raw colors, Semantic colors | Other component colors |

**Valid:**
```toml
[palette.semantic]
accent = "mauve"              # References raw color ✓

[palette.components]
button-bg = "accent"          # References semantic color ✓
button-bg = "mauve"           # References raw color ✓
```

**Invalid:**
```toml
[palette.semantic]
accent-hover = "accent"       # Cannot reference semantic ✗

[palette.components]
button-hover = "button-bg"    # Cannot reference component ✗
```

This prevents circular dependencies and ensures predictable color resolution.

### CSS Generation

Palettes are compiled to CSS custom properties:

```css
/* Generated from catppuccin-mocha palette */
:root {
  /* Raw colors */
  --palette-rosewater: #f5e0dc;
  --palette-flamingo: #f2cdcd;
  --palette-pink: #f5c2e7;
  /* ... all 26 Catppuccin colors */

  /* Semantic colors */
  --color-text-primary: var(--palette-text);
  --color-text-secondary: var(--palette-subtext1);
  --color-text-muted: var(--palette-overlay1);
  --color-bg-primary: var(--palette-base);
  --color-bg-secondary: var(--palette-mantle);
  --color-bg-surface: var(--palette-surface0);
  --color-accent: var(--palette-mauve);
  --color-link: var(--palette-blue);
  --color-success: var(--palette-green);
  --color-warning: var(--palette-yellow);
  --color-error: var(--palette-red);
  --color-info: var(--palette-blue);
  --color-border: var(--palette-surface1);

  /* Component colors */
  --code-bg: var(--palette-surface0);
  --code-text: var(--palette-text);
  --code-keyword: var(--palette-mauve);
  /* ... */
}

/* Light mode override (when visitor explicitly chooses light mode) */
[data-theme="light"] {
  /* Light palette colors override here */
}
```

### Contrast Validation

Palettes SHOULD pass WCAG 2.1 contrast requirements:

| Combination | Minimum Ratio | Level |
|-------------|---------------|-------|
| text-primary on bg-primary | 4.5:1 | AA |
| text-secondary on bg-primary | 4.5:1 | AA |
| text-muted on bg-primary | 3:1 | AA Large |
| accent on bg-primary | 3:1 | AA Large |
| link on bg-primary | 4.5:1 | AA |

Implementations SHOULD provide a contrast validation command:

```bash
$ markata-go palette check catppuccin-mocha

Checking palette: Catppuccin Mocha

Contrast Ratios:
  text-primary on bg-primary:     11.8:1 ✓ (AA, AAA)
  text-secondary on bg-primary:    8.2:1 ✓ (AA, AAA)
  text-muted on bg-primary:        4.1:1 ✓ (AA Large)
  accent on bg-primary:            5.2:1 ✓ (AA)
  link on bg-primary:              6.4:1 ✓ (AA, AAA)
  error on bg-primary:             5.8:1 ✓ (AA)
  warning on bg-primary:           9.1:1 ✓ (AA, AAA)
  success on bg-primary:           7.3:1 ✓ (AA, AAA)

All contrast checks passed!
```

### Contrast Ratio Calculation

Implementations MUST use the WCAG 2.1 relative luminance formula for contrast calculations.

**Algorithm:**

1. **Calculate contrast ratio between foreground and background colors:**
   - Compute relative luminance for each color
   - Ensure lighter color is used as L1
   - Return: (L1 + 0.05) / (L2 + 0.05)
   - Result ranges from 1:1 (same color) to 21:1 (black/white)

2. **Calculate relative luminance (per WCAG 2.1 sRGB):**
   - Convert each RGB channel (0-255) to 0-1 range
   - Linearize each channel:
     - If value ≤ 0.04045: divide by 12.92
     - Otherwise: ((value + 0.055) / 1.055) ^ 2.4
   - Apply ITU-R BT.709 coefficients:
     - L = 0.2126 × R_linear + 0.7152 × G_linear + 0.0722 × B_linear

### WCAG Compliance Levels

| Level | Normal Text | Large Text | UI Components |
|-------|-------------|------------|---------------|
| A     | 3:1         | 3:1        | 3:1           |
| AA    | 4.5:1       | 3:1        | 3:1           |
| AAA   | 7:1         | 4.5:1      | 4.5:1         |

**Large text** is defined as:
- 18pt (24px) or larger for normal weight
- 14pt (18.5px) or larger for bold weight

### Required Contrast Checks

Implementations MUST validate these combinations:

| Foreground | Background | Min Ratio | Level |
|------------|------------|-----------|-------|
| text-primary | bg-primary | 4.5:1 | AA |
| text-primary | bg-surface | 4.5:1 | AA |
| text-primary | bg-elevated | 4.5:1 | AA |
| text-secondary | bg-primary | 4.5:1 | AA |
| text-muted | bg-primary | 3.0:1 | AA Large |
| link | bg-primary | 4.5:1 | AA |
| accent | bg-primary | 3.0:1 | AA Large |
| success | bg-primary | 3.0:1 | UI |
| warning | bg-primary | 3.0:1 | UI |
| error | bg-primary | 3.0:1 | UI |
| info | bg-primary | 3.0:1 | UI |
| code-text | code-bg | 4.5:1 | AA |
| code-comment | code-bg | 3.0:1 | AA Large |
| code-keyword | code-bg | 4.5:1 | AA |
| button-primary-text | button-primary-bg | 4.5:1 | AA |
| button-secondary-text | button-secondary-bg | 4.5:1 | AA |

### Testing Integration

Contrast validation SHOULD be available in multiple contexts:

**1. CLI Command:**
```bash
markata-go palette check <palette-name> [--strict]
```

**2. Automated Tests:**

Tests SHOULD verify all required contrast checks pass:
- Load the palette
- Run contrast checks for all required combinations
- Report any failures with actual vs required ratios

**3. CI Integration:**
```yaml
# .github/workflows/test.yml
- name: Validate palette contrast
  run: markata-go palette check --all --strict
```

### Accessibility Report

The `palette check` command SHOULD generate an accessibility report:

```bash
$ markata-go palette check catppuccin-mocha --report

=== Accessibility Report: Catppuccin Mocha ===

WCAG 2.1 Level AA Compliance: PASS

Text Readability:
  ✓ Primary text on backgrounds:     11.8:1 - 8.2:1 (excellent)
  ✓ Secondary text on backgrounds:    8.2:1 - 5.4:1 (good)
  ✓ Muted text on backgrounds:        4.1:1 - 3.2:1 (acceptable)

Interactive Elements:
  ✓ Links are distinguishable:        6.4:1 (AA compliant)
  ✓ Focus indicators visible:         5.2:1 (AA compliant)

Color Blindness Simulation:
  ✓ Protanopia:   Status colors distinguishable
  ✓ Deuteranopia: Status colors distinguishable  
  ✓ Tritanopia:   Status colors distinguishable

Recommendations:
  - Consider increasing muted text contrast for AAA compliance
  - Warning color could be slightly darker for better contrast

Full report saved to: .markata/palette-report-catppuccin-mocha.html
```

---

## Live Theme Picker

Every site ships a visitor-facing theme picker by default
(`theme.switcher.enabled` defaults to `true`). Setting it to `false` restores
single-palette output and hides the picker; the light/dark toggle remains
governed by `theme.switcher.mode_toggle` and `header.show_theme_toggle`.

```toml
[markata-go.theme.switcher]
enabled = true          # default
mode_toggle = true      # light/dark button
include_all = true      # every discovered palette
include = []            # or an explicit allow-list
exclude = []
```

### Initial palette stylesheet

When the switcher is enabled, `css/palette-core.css` MUST contain the configured
light/dark pair, picker manifest, fallback rules, and site overrides. The full
`css/palette.css` MUST remain available for custom templates and picker previews.
Both artifacts MUST have content hashes registered before HTML rendering.

The default template MUST load the core stylesheet initially. Its theme
bootstrap MUST request the full stylesheet with render blocking when restoring
a saved or seasonal palette outside the configured pair. Opening the picker
MUST load the full stylesheet before rendering palette previews. Direct palette
selection and keyboard shortcuts MUST also request it when necessary. Failed
requests MUST be retryable, with the configured pair remaining usable.

Sites without a switcher MUST retain their existing single stylesheet behavior.

Palette blocks MUST style both preview cards and the selected document root.
Root palette selectors MUST outrank generic `:root` defaults so a deferred
library keeps its selected colors when page bundles change during navigation.

Image zoom assets MUST load on the first zoom interaction rather than page
startup. The first interaction MUST open the selected image after assets load.
Subsequent zooms and view-transition navigation MUST reuse the loaded library.

### Page stylesheet bundles

`css_bundle.pages` MAY list output-relative HTML files, such as `index.html`.
With CSS bundling enabled, production builds MUST consolidate eligible local
head stylesheets and inline styles for these pages into content-hashed CSS.
The cascade order MUST remain unchanged. Synchronous scripts and ineligible
stylesheets MUST form bundle boundaries. External, conditional, integrity-
protected, missing, imported, or relative-URL stylesheets MUST remain intact.
Bundling MUST run after CSS generation, purging, and minification. Empty page
lists and fast builds MUST preserve existing behavior. All paths MUST remain
inside the output directory. A second pass MUST be idempotent.

### Header UI

The header renders one compact button (a painter's-palette icon whose dots use
the live `--color-primary`, `--color-code-string`, and `--color-warning`) plus
the mode toggle. No other nav items are added. The button opens a
fixed-height popover (a bottom sheet at `max-width: 640px`) with an ARIA
tablist of three tabs, each a `role="listbox"` grid of preview cards:

| Tab | Cards | Scoping |
|-----|-------|---------|
| Colors | Seasonal card, then one card per palette for the browsed mode; search, Light/Dark, Shuffle, seasonal schedule toggle | card sets `data-palette` |
| Style | one card per aesthetic | card sets `data-aesthetic` |
| Font | text-size radio buttons, then one card per font pack | inline font stacks from the manifest, applied when the card intersects the grid (IntersectionObserver) |

Hover preview: for mouse pointers (`pointerType === "mouse"`), entering a
card, a year-strip segment, or a Coming-up row applies its palette (with the
palette's variant as `data-theme`), aesthetic, or fontpack to `<html>` without
storing anything, at most once per animation frame. The real attributes are
saved on the first preview and restored 90ms after the pointer leaves all
cards (gaps inside a grid, the year strip, or the list keep the last preview),
and before any choice is committed, the tab or schedule view changes, the mode
changes, Reset, or the picker closes.

Tabs with fewer than two choices are hidden (Font stays while the size
control is shown). The top bar has ‹ › step buttons and a live label with
the current choice. The buttons move through the active tab's visible cards,
wrap at the ends, and apply each choice immediately. The last tab is stored
as `theme-picker-tab`. The bottom bar holds a visibly labeled Reset button
(clears the stored palette, aesthetic, and font choices)
and, under `markata-go serve` only, Bake. The popover uses a fixed system UI font, so it does not
reflow while fonts are previewed. On small screens, every control is at least
40px tall and the sheet pads for `env(safe-area-inset-bottom)`.

### Fonts

When the picker is enabled, the fontpack plugin resolves every catalog pack
(not only the configured one) into `css/fonts.css`, one
`[data-fontpack="<name>"]` block per pack, and appends:

```css
:root{--fontpack-default:"brush";--fontpack-manifest:'[{"name":"brush","displayName":"Brush","heading":"Knewave","body":"Space Grotesk","code":"DM Mono"},...]'}
```

Each manifest entry also carries `headingFont`, `headingWeight`, `bodyFont`,
and `codeFont`, the CSS values of the pack's `--font-*` roles. Quotes are
normalized to single quotes so the JSON survives the custom-property round
trip. `@font-face` files are only fetched when used. A visitor choice is stored as
`theme-fontpack` and wins over per-post `fontpack` overrides. Without a
stored choice, SPA navigation keeps each page's own `data-fontpack`.

### Styles

The built-in preset TOML files define the shape, spacing, border, shadow and
typography tokens; `aesthetic.css` supplies their CSS consumption rules.
Generated aesthetic CSS includes these rules after the preset variables so
the picker and a built site use the same styling. Tokens use bare
`[data-aesthetic="<id>"]` selectors for Style card previews. Component rules
are scoped to `html[data-aesthetic]` in the `overrides` layer, leaving
unlayered site custom CSS and theme variable overrides authoritative.
Balanced retains the theme's native spacing, surfaces and hover treatment;
other presets adjust cards, code, quotes, admonitions, tables, media, forms,
tags, header, picker and prose rhythm. Drawer shapes/depth only change when
the fixed doc/feed sidebar is visible. No aesthetic changes palette colors.

Brutal uses square corners, thick borders, hard offset shadows and mono UI
accents; precision uses tight spacing, small radii, hairlines and crisp
shadows; elevated uses generous spacing, rounded corners and layered soft
shadows; minimal is airy and flat, with borders only where needed.

Because low-alpha shadows and small radius changes are nearly invisible on
dark palettes, each non-balanced preset MUST also carry one surface signature
that reads on light and dark palettes:

| Preset | Signature |
|--------|-----------|
| brutal | thick strong borders, hard offset primary shadows |
| precision | visible hairline borders, gridded table cells, uppercase mono table headers, 2px primary rule on code blocks, mono tags, no soft shadow |
| elevated | borderless tinted surfaces with an inset highlight and a deep drop shadow, tinted quotes, pill tags, cards lift 2px on hover |
| minimal | no card/table borders or backgrounds, faintly tinted code, 1px quote rule, underlined text tags |

Signature rules use `html[data-aesthetic][data-aesthetic="<id>"]` so they
outrank the shared non-balanced rules in the same layer.

### Seasonal

`seasonal` is a pseudo-palette: the picker stores it as
`theme-palette-light` and `theme-palette-dark`. `theme.seasonal = true` makes
it the default when no palette is stored. With a seasonal default, picking
any concrete palette, including the site pair, stores it explicitly.

The palette CSS plugin resolves a calendar against the manifest and publishes
it as `config.Extra.palette_seasonal_json`. `base.html` inlines it as
`window.__markataSeasonal`, together with `__markataSeasonalDefault` and
`__markataSeasonalPick(mode, date?)`, which returns `{label, name}`:

```json
{"lead":3,
 "s":[{"l":"Spring","f":320,"lt":"pollen8","dk":"pollen8-dark"}, ...],
 "h":[{"l":"Christmas","d":["12-25"],"lt":"christmas-light","dk":"christmas"},
      {"l":"Hanukkah","d":["2026-12-05", ...],"n":8,"lt":"...","dk":"..."}, ...]}
```

Resolution uses the visitor's local date:

1. **Holidays.** A holiday matches from `lead` days before its date through
   the end of the celebration (`n` days, default 1). `MM-DD` dates are tried
   in the previous, current, and next year. When several match, the one with
   the soonest date wins.
2. **Seasons.** With no holiday match, the season is the last one whose start
   (`f`, as MMDD) is on or before today. Winter wraps the year.

Seasons follow the northern hemisphere, astronomical (spring 03-20, summer
06-21, autumn 09-22, winter 12-21).

Holidays:

- **Fixed dates:** New Year, Valentine's Day, St. Patrick's Day, Earth Day,
  Halloween, Christmas.
- **Computed:** Western Easter, using the Gregorian computus for the build
  year −1 through +10.
- **Tabulated:** Lunar New Year (through 2036), Diwali/Lakshmi Puja (through
  2033), and Hanukkah (first full day, 8 days, through 2035).

For each mode, the first listed palette with that variant is used, then the
first palette's counterpart. Entries without both variants are dropped. When
no season resolves, the schedule is omitted and the Seasonal card is hidden.

#### Schedule View

When a schedule exists, the Colors tab shows a calendar toggle
(`[data-picker-schedule]`, `aria-pressed`) that swaps the grid for
`#theme-picker-schedule`:

- Periods are computed by calling `__markataSeasonalPick(mode, date)` for each
  day and merging consecutive days with the same pick, so the view always
  matches what the head script applies. `mode` is the mode being browsed.
- A year strip (Jan 1 to Dec 31 of the visitor's current year) with one
  segment per period, width proportional to its days, scoped with
  `data-palette`, and a marker for today. It has an `aria-label` listing every
  period.
- A **Coming up** list for the next 366 days: the current period ("Now · until
  <date>"), then each holiday with its date range and each season the first
  time it starts ("From <date>"). A season resuming after a holiday is not
  listed again.
- Durations: holidays show their day count (lead days included); seasons show
  the full season length (start to the next season's start, holidays
  included); the current period shows the time left. Under 14 days is shown in
  days, otherwise in weeks (`~` when not a whole number of weeks).
- A readout under the strip shows the hovered (or tapped, on touch) period's
  label, dates, length, and palette. The strip segment and list row of the same
  period share a `data-period-key` and get `.is-linked` together. While any
  hover preview is active, the top label reads `Preview: <name>`.
- **Use Seasonal** selects seasonal; it is disabled when seasonal is active.
- Searching, Shuffle, the ‹ › step buttons, closing the picker, or pressing the
  toggle again returns to the grid. Changing the browsed mode re-renders the
  view.
Bake writes `seasonal = true` and leaves the configured palettes as the fallback.

### Bake

Bake writes the picker's current choices into the site config. It exists only
under `markata-go serve`: the dev server injects
`window.__markataThemeBakeEndpoint = "/__markata/theme/bake"` into served
HTML, and the `[data-picker-bake]` button stays `hidden` unless that global is
set. Static builds never contain the endpoint.

The first click sends `GET /__markata/theme/bake`, which returns
`{target, path, exists, target_kind, warnings}`, and relabels the button
`Bake into <target>?` (with `(override)` for a `--merge-config` file). When
`target_kind` is `global` the button reads "Bake unavailable" instead. A second
click within 4s sends `POST` with a JSON body:

```json
{"palette": "<current>", "palette_light": "<light pick>", "palette_dark": "<dark pick>",
 "seasonal": false, "fallback_mode": "light|dark", "aesthetic": "<id>",
 "fontpack": "<id>", "text_size": "small|medium|large|x-large"}
```

The server writes these `[markata-go.theme]` keys, in this order:

- When `seasonal` is true, it writes `seasonal = true` and does not touch the palettes.
- Otherwise it writes `palette`, `palette_light`, and `palette_dark`, plus
  `seasonal = false` when the target file or effective config has seasonal on.
- Then `fallback_mode`, `aesthetic`, `fontpack`, and `text_size`. Empty fields are skipped.

`fallback_mode` and `text_size` are visitor preferences, so the picker sends
them only when the visitor explicitly chose them in this browser (the
`color-mode` or `text-size` key is stored); otherwise the site's settings are
left alone.

The picks are converted to `theme.*` settings changes and applied through the
settings sidebar's apply path (see "Serve Settings Sidebar" in
[CONFIG.md](CONFIG.md)): the same target selection, `override`/`global`
target kinds (a global target is refused with 409), reload verification,
atomic rollback, and serve config lock. Baked keys are dropped from any
active settings preview, so the rebuilt page shows the baked value rather than
a stale preview. The response lists the baked `keys`.

On success the picker removes the stored palette, aesthetic, and fontpack
picks (and text size when it was baked), notifies which keys were baked, the
server queues a full rebuild, and live reload shows the new defaults.

**Target file.** Candidates are the root config, its `include` files
(recursively, in load order), then `--merge-config` files. This is the same
lowest-to-highest precedence order the loader uses.

- Each key goes to the candidate that defines the deepest prefix of
  `theme.<key>` (so normally the last file that defines `markata-go.theme`).
- If none defines it, the root config is used.
- With no config file, `markata-go.toml` is created in the site directory.

**Editing rules.** Edits MUST preserve everything except the baked keys:

- Existing keys are replaced in place, keeping their trailing comments.
- Missing keys are added to the existing group. TOML appends them after the
  table's last key; YAML inserts them first in the block mapping, using the
  file's indentation.
- A missing group is appended as `[markata-go.theme]` (TOML) or inserted under
  `markata-go:` (YAML). JSON keeps key order and indentation.
- TOML dotted keys (`theme.palette = ...` under `[markata-go]`) are edited in place.

**Refusals.** Layouts that cannot be edited without reformatting fail with
HTTP 409 and leave the file unchanged: TOML inline tables, YAML flow
mappings, block scalars, and multi-line strings. (Multi-line TOML arrays and
YAML block sequences are replaced whole.) The edited document is
re-parsed, and it MUST equal the original plus the baked keys, or the write
is aborted. Writes are atomic (temp file plus rename) and keep file permissions.

**Request safety.** The endpoint accepts only GET and POST.

- POST requires `Content-Type: application/json`, and requests from another
  origin (`Origin` host mismatch or `Sec-Fetch-Site: cross-site`) get 403.
- Bodies over 8 KiB, unknown fields, names outside
  `[A-Za-z0-9][A-Za-z0-9 ._'&+()-]{0,79}`, and unknown modes or sizes get 400.
- `MARKATA_GO_THEME_*` environment overrides are reported as warnings.

Included config files are part of the build-cache config hash, so editing
any of them, by bake or by hand, invalidates cached pages.

Bake is the theme-only front end of the general settings sidebar, which edits
any config setting with the same editor; see "Serve Settings Sidebar" in
[CONFIG.md](CONFIG.md). The picker shows an `[data-picker-settings]`
**All settings** button, hidden unless `window.__markataSettingsEndpoint` is
set, that closes the picker and calls `window.markataDevSettings.open('theme')`.

### Manifest

`palette.css` is emitted with one scoped block per palette
(`[data-palette="<name>"] { ... }`) and a JSON manifest stored in the
`--palette-manifest` custom property (aesthetics use `--aesthetic-manifest`).
Each palette entry has:

| Field | Meaning |
|-------|---------|
| `name` | Palette id used in `data-palette` |
| `displayName` | Human label |
| `variant` | `light` or `dark` |
| `baseName` | Family label shared by light/dark partners |
| `counterpart` | Name of the opposite-mode partner, if any |
| `derived` | `true` when generated by `Loader.Discover()` |

The manifest is deterministic: `Discover()` snapshots explicit palettes
before deriving counterparts and never derives from a derived palette.
Aesthetics use normalized ids (`aestheticID()`, e.g. `minimal`) in both the
manifest and `[data-aesthetic]` selectors.

### No-Flash Restore

A blocking inline script in `<head>`, before any stylesheet, reads
`localStorage` (`color-mode`, `theme-palette-<mode>`, `theme-aesthetic`,
`theme-fontpack`), resolves `seasonal` to today's palette, and sets `data-theme`, `data-palette`, `data-aesthetic`, and
`data-fontpack` on `<html>` (recording the page's own pack in
`window.__markataPageFontpack` first). Without a stored mode, the site's
`fallback_mode` (default `dark`) is used; the OS preference is ignored so the
site is dark by default. `color-scheme` follows `data-theme`. It falls
back to the configured light/dark palettes (`window.__markataThemeDefaults`)
when nothing is stored; the runtime discards stored palettes missing from the
manifest and ignores legacy keys (`selected-palette`, `selected-family`,
`selected-aesthetic`). Contract fallbacks for `:root` dark/light use
`:where(...)` so their specificity never beats a `[data-palette]` block.
View transitions preserve `data-palette` and `data-aesthetic` across SPA
navigation.

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `t` | Open/close picker |
| `.` / `,` | Next / previous palette family (applies immediately) |
| `>` / `<` | Next / previous aesthetic |
| `f` / `F` | Next / previous font pack |
| `1` `2` `3` (in popover) | Colors / Style / Font tab |
| `\` | Toggle light/dark within the current family |
| Arrows, Home/End, PageUp/PageDown | Move and live-preview inside the popover; `Enter` selects, `Esc` reverts |

### Motif Recoloring

The motif is off by default (`theme.motif.kind = "off"` in the contract
defaults). The recipe compiler still compiles texture and heading passes when
the motif is off and drops the unused motif asset. Aesthetic motif images (SVG) bake the configured palette's colors. When
`html[data-palette]` differs from the configured palette, the motif is
reused as a CSS mask painted with
`color-mix(in srgb, var(--color-text) N%, var(--color-background))` so it
matches the live palette.

### Config-less Builds

When no config file is found (`LoadWithDefaults`), the theme is normalized
exactly as an empty config file would be, so single-file builds get the
contract defaults (`ayu-dark` palette, `minimal` aesthetic) and a working
picker. If no manifest is available at runtime, the picker container is
hidden.

---

## Feature Flags

Themes can expose optional features that users can enable/disable:

```toml
[markata-go.theme]
name = "default"

[markata-go.theme.features]
dark_mode = true              # Enable dark mode toggle
toc = true                    # Show table of contents
reading_time = true           # Show estimated reading time
copy_code = true              # Add copy button to code blocks
back_to_top = true            # Show back-to-top button
search = false                # Disable search (if theme supports it)
comments = false              # Disable comments integration
```

### Feature Implementation

Themes declare supported features in `theme.toml`:

```toml
[theme.features]
dark_mode = { default = true, description = "Dark mode toggle in header" }
toc = { default = true, description = "Table of contents sidebar" }
reading_time = { default = true, description = "Reading time estimate" }
copy_code = { default = true, description = "Copy button on code blocks" }
back_to_top = { default = false, description = "Floating back-to-top button" }
search = { default = false, description = "Site search functionality" }
comments = { default = false, description = "Comments via Giscus/Utterances" }
```

Features are available in templates:

```jinja2
{% if features.toc and post.toc %}
<aside class="toc">
  {{ post.toc | safe }}
</aside>
{% endif %}

{% if features.reading_time %}
<span class="reading-time">{{ post.reading_time }} min read</span>
{% endif %}

{% if features.copy_code %}
<script src="{{ 'js/copy-code.js' | theme_asset }}"></script>
{% endif %}
```

---

## Conditional CSS Loading

To optimize page size and load times, markata-go conditionally loads CSS files based on page content. This is implemented in two places:

### Content-Based Detection (RenderMarkdownPlugin)

The `render_markdown` plugin scans rendered HTML content and sets flags in `post.Extra`:

| Flag | Detected When |
|------|---------------|
| `needs_admonitions_css` | HTML contains `class="admonition` |
| `needs_code_css` | HTML contains `class="chroma"`, `class="highlight"`, `<pre><code`, or `<code class="language-` |
| `needs_image_zoom` | HTML contains `data-zoomable` or `class="glightbox"` (set by `image_zoom` plugin) |
| `needs_roads` | HTML has a `road` class token or a `data-road` attribute (see [ROADS.md](ROADS.md)) |

### Other Plugin Detection

Other plugins set their own flags in `post.Extra`:

| Flag | Set By | Detected When |
|------|--------|---------------|
| `has_encrypted_content` | `encryption` plugin | Post has `private: true` frontmatter |
| `glightbox_enabled` | `image_zoom` plugin | GLightbox is configured and post has zoomable images |
| `needs_cards_css` | `publish_feeds` plugin | Page is a feed/index page that renders post cards |
| `needs_webmentions_css` | `webmentions` plugin | Post has webmentions or webmentions are enabled globally |

### Template Implementation

The base template (`base.html`) uses these flags to conditionally load CSS:

```jinja2
<!-- Core CSS (always loaded) -->
<link rel="stylesheet" href="{{ 'css/variables.css' | theme_asset }}">
<link rel="stylesheet" href="{{ 'css/main.css' | theme_asset }}">
<link rel="stylesheet" href="{{ 'css/components.css' | theme_asset }}">

<!-- Cards CSS (only on feed/index pages) -->
{% if needs_cards_css %}
<link rel="stylesheet" href="{{ 'css/cards.css' | theme_asset }}">
{% endif %}

<!-- Admonitions CSS (only when content uses admonitions) -->
{% if needs_admonitions_css %}
<link rel="stylesheet" href="{{ 'css/admonitions.css' | theme_asset }}">
{% endif %}

<!-- Code CSS (only when content has code blocks) -->
{% if needs_code_css %}
<link rel="stylesheet" href="{{ 'css/code.css' | theme_asset }}">
<link rel="stylesheet" href="/css/chroma.css">
{% endif %}

<!-- Webmentions CSS (only when post has webmentions) -->
{% if post.webmentions or needs_webmentions_css %}
<link rel="stylesheet" href="{{ 'css/webmentions.css' | theme_asset }}">
{% endif %}

<!-- Encryption CSS (only for encrypted/private posts) -->
{% if has_encrypted_content %}
<link rel="stylesheet" href="{{ 'css/encryption.css' | theme_asset }}">
{% endif %}

<!-- GLightbox CSS (only when images need zoom) -->
{% if glightbox_enabled and needs_image_zoom %}
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/glightbox/dist/css/glightbox.min.css">
{% endif %}

<!-- Feature-based CSS -->
{% if config.theme.switcher.enabled or (config.theme.switcher.mode_toggle and config.header.show_theme_toggle) %}
<link rel="stylesheet" href="{{ 'css/palette-switcher.css' | theme_asset }}">
{% endif %}
```

### View-Transition Safety Net

Because view transitions (SPA-like navigation) preserve `<head>` across navigations, CSS `<link>` tags persist from the original page. When navigating from a page without conditional CSS (e.g., a plain text post) to a page that needs it (e.g., a post with code blocks), the CSS would be missing unless the runtime adds it.

The `conditional-css.js` script solves this by:
1. Listening for `view-transition-complete` events dispatched by `view-transitions.js`
2. Scanning the new page's body for content markers (admonition classes, chroma classes, etc.)
3. Injecting the needed `<link>` tags into `<head>` if not already present
4. Removing conditional CSS that is no longer needed

Any JS-driven theme controls rendered in `<body>` (palette family selector, dark/light toggle, and aesthetic selector) MUST also re-initialize on `view-transition-complete` so controls and option lists remain interactive after navigation.

View-transition navigations SHOULD fetch and parse the destination document before calling `document.startViewTransition()` so the current page is not visually frozen during network I/O. Implementations SHOULD allow normal browser caching for HTML navigations and MAY prefetch likely next destinations (for example hovered links or pagination targets) to reduce latency for keyboard-driven page jumps. Prefetching MUST be bounded so dense navigation regions do not create unbounded fetch and parse work. Expired prefetch entries MUST NOT permanently consume bounded cache capacity. An in-memory prefetched document MUST have a finite freshness window. A document ready for less than five seconds MAY be consumed without another request; an older document MUST be discarded and fetched with HTTP cache revalidation semantics (`cache: "no-cache"`). The HTTP cache MAY use retained ETag or Last-Modified validators to return an updated representation or 304; responses marked `no-store` may require a full response. Failed revalidation MUST fall back to ordinary document navigation.

When possible, implementations SHOULD replace only the primary view-transition regions (for example the progress area and page wrapper) instead of replacing the entire `<body>`. The runtime MUST still synchronize `<body>` attributes such as classes so page-specific layout state remains correct after navigation.

The `document.startViewTransition()` update callback SHOULD stay limited to critical state updates and DOM replacement. Script reinitialization and inline module re-execution SHOULD run after the transition finishes so expensive page hydration does not block the transition snapshot. When synchronizing `<html>` attributes, implementations MUST preserve runtime-only theme and transition state such as `data-theme`, the `dark` class, and active transition data attributes.

### CSS File Categories

| Category | Files | Loading Condition |
|----------|-------|-------------------|
| Core | `variables.css`, `main.css`, `components.css` | Always |
| Layout | `cards.css` | `needs_cards_css` (feed/index pages) |
| Content | `admonitions.css` | `needs_admonitions_css` (content detection) |
| Content | `code.css`, `chroma.css` | `needs_code_css` (content detection) |
| Features | `webmentions.css` | `post.webmentions` or `needs_webmentions_css` |
| Features | `encryption.css` | `has_encrypted_content` |
| Content | `roads.css`, `roads.js` | `needs_roads` (content detection) |
| Features | `glightbox.min.css` | `glightbox_enabled` and `needs_image_zoom` |
| Features | `palette-switcher.css` | `config.theme.switcher.enabled` or (`config.theme.switcher.mode_toggle` and `config.header.show_theme_toggle`) |

### Cascade Layers

The default theme's stylesheets are organised with CSS cascade layers,
declared once in `variables.css`:

```css
@layer reset, tokens, base, components, utilities, overrides;
```

Rules:

- Every top-level block in a theme stylesheet that uses layers MUST be an
  `@layer` block. Unlayered rules outrank every layered rule regardless of
  specificity, so a single stray brace that closes a layer early silently
  changes the cascade for the whole site. `pkg/themes/css_layers_test.go`
  enforces balanced braces and layered top-level blocks for shipped CSS.
- Palette tokens and plugin-generated theme tokens live in `tokens`; component
  styling lives in `components`; text-size presets live in `utilities`;
  author `[theme.variables]` overrides are emitted in `overrides`.
- Every page — posts, feeds, tag pages, archives — MUST receive the same
  `data-palette`, `data-aesthetic`, `data-fontpack` and texture attributes on
  `<html>`; `config.theme` exposes those resolved values to templates so
  typography is never scoped to post pages only.

### Estimated Savings

| Page Type | CSS Skipped | Approx. Savings |
|-----------|-------------|-----------------|
| Plain text post (no code, no admonitions) | code.css, chroma.css, admonitions.css, cards.css | ~92KB |
| Post with code blocks only | admonitions.css, cards.css | ~84KB |
| Feed/index page | code.css, chroma.css, admonitions.css | ~75KB |

### Benefits

1. **Reduced page size**: Pages without code blocks skip ~15KB of syntax highlighting CSS
2. **Faster initial render**: Less CSS to parse before first paint
3. **Better caching**: Shared core CSS files are cached across all pages
4. **View-transition safe**: Runtime script ensures CSS correctness during SPA navigation

---

## CSS Custom Properties

Built-in themes SHOULD use CSS custom properties for consistency:

### Colors

```css
:root {
  /* Primary brand color */
  --color-primary: #3b82f6;
  --color-primary-light: #60a5fa;
  --color-primary-dark: #2563eb;

  /* Semantic colors */
  --color-text: #1f2937;
  --color-text-muted: #6b7280;
  --color-background: #ffffff;
  --color-surface: #f9fafb;
  --color-border: color-mix(in srgb, #1f2937 16%, #ffffff);
  --color-border-strong: #1f2937;
  --color-border-soft: color-mix(in srgb, var(--color-border) 60%, transparent);

  /* Status colors */
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --color-info: #3b82f6;
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  :root {
    --color-text: #f9fafb;
    --color-text-muted: #9ca3af;
    --color-background: #111827;
    --color-surface: #1f2937;
    --color-border: #374151;
  }
}
```

### Typography

```css
:root {
  /* Font families */
  --font-body: system-ui, -apple-system, sans-serif;
  --font-heading: var(--font-body);
  --font-mono: ui-monospace, 'Cascadia Code', 'Fira Code', monospace;

  /* Font sizes (modular scale) */
  --text-xs: 0.75rem;
  --text-sm: 0.875rem;
  --text-base: 1rem;
  --text-lg: 1.125rem;
  --text-xl: 1.25rem;
  --text-2xl: 1.5rem;
  --text-3xl: 1.875rem;
  --text-4xl: 2.25rem;

  /* Line heights */
  --leading-tight: 1.25;
  --leading-normal: 1.5;
  --leading-relaxed: 1.75;
}
```

### Spacing

```css
:root {
  /* Spacing scale */
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;
  --space-12: 3rem;
  --space-16: 4rem;

  /* Layout */
  --content-width: 66ch;
  --page-width: 1200px;
  --radius-sm: 0.375rem;
  --radius: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-xl: 1rem;
}
```

---

## Media Styles

Built-in themes MUST include styles for images and videos.

### Images

```css
/* Base image styles */
img {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius);
}

/* Images in article content */
.post-content img {
  display: block;
  margin: var(--space-6) auto;
  border-radius: var(--radius-lg);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* Dark mode: use border instead of shadow */
@media (prefers-color-scheme: dark) {
  .post-content img {
    box-shadow: none;
    border: 1px solid var(--color-border);
  }
}

figure {
  margin: var(--space-6) 0;
}

figcaption {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
  text-align: center;
  margin-top: var(--space-2);
}
```

### Videos

The `md_video` plugin converts markdown image syntax to video tags. Themes MUST style these appropriately:

```css
/* Base video styles */
video,
.md-video {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius);
  display: block;
  margin: var(--space-4) 0;
  background: var(--color-surface);
}

/* Videos in article content */
.post-content video,
.post-content .md-video {
  display: block;
  margin: var(--space-6) auto;
  border-radius: var(--radius-lg);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  .post-content video,
  .post-content .md-video {
    box-shadow: none;
    border: 1px solid var(--color-border);
  }
}

/* Loading state placeholder */
video {
  background: var(--color-surface) url("data:image/svg+xml,...") center center no-repeat;
  min-height: 200px;
}
```

### Fullscreen Media Performance

Themes MUST keep fullscreen video playback smooth when decorative backgrounds are enabled.

- Content containers MUST NOT apply `backdrop-filter: blur(0px)` or an equivalent no-op filter by default.
- Backdrop blur for article, feed, or main content wrappers MUST be opt-in via configuration.
- Theme-level `::backdrop` styling MUST NOT apply blur globally because it also affects fullscreen media backdrops.
- When blur is disabled, the computed backdrop filter for those wrappers SHOULD be `none` so browsers can promote fullscreen video without repeated page recompositing.

Recommended background configuration fields:

```toml
[markata-go.theme.background]
enabled = true
article_bg = "color-mix(in srgb, var(--color-background) 88%, transparent)"
article_blur_enabled = false  # default; safer for fullscreen video
article_blur = "12px"        # only used when article_blur_enabled = true
```

### Full-Width Media

Support for edge-to-edge media display:

```css
.post-content img.full-width,
.post-content video.full-width {
  width: 100%;
  max-width: none;
  margin-left: calc(-1 * var(--space-4));
  margin-right: calc(-1 * var(--space-4));
  width: calc(100% + var(--space-8));
  border-radius: 0;
}
```

### Media Borders

Themes SHOULD provide configurable media borders via CSS custom properties:

```css
:root {
  /* Media border settings */
  --media-border-width: 3px;
  --media-border-style: solid;
  --media-border-color: var(--color-border);
  --media-border-radius: var(--radius-lg);

  /* Generic gradient presets */
  --gradient-accent: linear-gradient(135deg, var(--color-primary), var(--color-primary-dark));
  --gradient-vibrant: linear-gradient(135deg, #667eea, #764ba2, #f093fb);
  --gradient-warm: linear-gradient(135deg, #f093fb, #f5576c, #f8b500);
  --gradient-cool: linear-gradient(135deg, #4facfe, #00f2fe);
  --gradient-sunset: linear-gradient(135deg, #fa709a, #fee140);
  --gradient-ocean: linear-gradient(135deg, #2193b0, #6dd5ed);

  /* Palette-specific gradients - use colors from popular palettes */
  --gradient-catppuccin: linear-gradient(135deg, #cba6f7, #f5c2e7, #f38ba8);
  --gradient-nord: linear-gradient(135deg, #88c0d0, #81a1c1, #5e81ac);
  --gradient-dracula: linear-gradient(135deg, #bd93f9, #ff79c6, #8be9fd);
  --gradient-gruvbox: linear-gradient(135deg, #fabd2f, #fe8019, #fb4934);
  --gradient-rose-pine: linear-gradient(135deg, #c4a7e7, #ebbcba, #f6c177);
  --gradient-solarized: linear-gradient(135deg, #268bd2, #2aa198, #859900);
  --gradient-tokyo-night: linear-gradient(135deg, #7aa2f7, #bb9af7, #f7768e);
}
```

### Palette-Specific Gradient Borders

Themes SHOULD include gradient presets that match popular color palettes. This allows users to maintain visual consistency when using a specific palette:

| Gradient Class | Palette | Colors |
|---------------|---------|--------|
| `gradient-catppuccin` | Catppuccin | Mauve → Pink → Red |
| `gradient-nord` | Nord | Frost colors (cyan → blue) |
| `gradient-dracula` | Dracula | Purple → Pink → Cyan |
| `gradient-gruvbox` | Gruvbox | Yellow → Orange → Red |
| `gradient-rose-pine` | Rosé Pine | Iris → Rose → Gold |
| `gradient-solarized` | Solarized | Blue → Cyan → Green |
| `gradient-tokyo-night` | Tokyo Night | Blue → Purple → Pink |

Usage in frontmatter:

```yaml
---
title: "My Catppuccin-Themed Post"
css_class: gradient-catppuccin
---
```

CSS implementation:

```css
/* Palette-specific gradient classes */
.post-content.gradient-catppuccin img,
.post-content.gradient-catppuccin video {
  border: none;
  padding: var(--media-border-width);
  background: var(--gradient-catppuccin);
  background-origin: border-box;
}

.post-content.gradient-nord img,
.post-content.gradient-nord video {
  background: var(--gradient-nord);
}

/* ... similar for dracula, gruvbox, rose-pine, solarized, tokyo-night */
```

This allows users who select a palette like `catppuccin-mocha` to also use matching gradient borders that complement their color scheme.

### Gradient Borders

Themes MAY support gradient borders for media elements. The technique uses padding and background:

```css
/* Apply gradient borders to all media in post content */
.post-content.gradient-borders img,
.post-content.gradient-borders video {
  border: none;
  padding: var(--media-border-width);
  background: var(--gradient-accent);
  background-origin: border-box;
}

/* Specific gradient variants */
.post-content.gradient-vibrant img { background: var(--gradient-vibrant); }
.post-content.gradient-warm img { background: var(--gradient-warm); }
.post-content.gradient-cool img { background: var(--gradient-cool); }
.post-content.gradient-sunset img { background: var(--gradient-sunset); }
.post-content.gradient-ocean img { background: var(--gradient-ocean); }
```

### Animated Gradient Borders

For eye-catching media display:

```css
@keyframes gradient-rotate {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.post-content.gradient-animated img,
.post-content.gradient-animated video {
  background: linear-gradient(270deg, #667eea, #764ba2, #f093fb, #667eea);
  background-size: 300% 300%;
  animation: gradient-rotate 6s ease infinite;
}
```

### Glow Effects

For dramatic media presentation:

```css
.post-content.glow img,
.post-content.glow video {
  box-shadow:
    0 0 20px rgba(102, 126, 234, 0.3),
    0 0 40px rgba(118, 75, 162, 0.2),
    0 0 60px rgba(240, 147, 251, 0.1);
}

/* Enhanced glow in dark mode */
@media (prefers-color-scheme: dark) {
  .post-content.glow img,
  .post-content.glow video {
    box-shadow:
      0 0 30px rgba(102, 126, 234, 0.4),
      0 0 60px rgba(118, 75, 162, 0.3),
      0 0 90px rgba(240, 147, 251, 0.2);
  }
}
```

### Usage in Templates

Apply gradient borders via class on the content wrapper:

```html
<!-- Default borders -->
<div class="post-content">
  {{ post.ArticleHTML | safe }}
</div>

<!-- Gradient borders -->
<div class="post-content gradient-borders">
  {{ post.ArticleHTML | safe }}
</div>

<!-- Vibrant gradient with glow -->
<div class="post-content gradient-vibrant glow">
  {{ post.ArticleHTML | safe }}
</div>

<!-- Animated gradient -->
<div class="post-content gradient-animated">
  {{ post.ArticleHTML | safe }}
</div>
```

### Configuration via Frontmatter

Posts can specify their preferred media style:

```yaml
---
title: "My Cool Demo"
media_style: gradient-vibrant glow
---
```

Template usage:

```html
<div class="post-content {{ post.Frontmatter.media_style | default('') }}">
  {{ post.ArticleHTML | safe }}
</div>
```

---

## Admonition Styles

Built-in themes MUST include styles for all admonition types.

Built-in themes MUST keep interactive text and compact controls readable by default. This includes using link and accent treatments that pass WCAG AA contrast against the active surface and giving compact interactive elements at least a 24px effective touch target when they appear in dense UI such as cards, sidebars, or feed navigation.

### Admonition CSS

```css
/* Base admonition styles */
.admonition {
  margin: var(--space-4) 0;
  padding: var(--space-4);
  border-left: 4px solid var(--admonition-color, var(--color-primary));
  background: var(--admonition-bg, var(--color-surface));
  border-radius: 0 var(--radius) var(--radius) 0;
}

.admonition-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-weight: 600;
  margin-bottom: var(--space-2);
  color: var(--color-text);
  color: color-mix(in srgb, var(--admonition-color, var(--color-primary)) 45%, var(--color-text) 55%);
}

.admonition-title::before {
  content: var(--admonition-icon, "");
  font-size: 1.25em;
}

/* Type-specific styles */
.admonition.note {
  --admonition-color: #3b82f6;
  --admonition-bg: #eff6ff;
  --admonition-icon: "ℹ️";
}

.admonition.info {
  --admonition-color: #3b82f6;
  --admonition-bg: #eff6ff;
  --admonition-icon: "ℹ️";
}

.admonition.tip {
  --admonition-color: #10b981;
  --admonition-bg: #ecfdf5;
  --admonition-icon: "💡";
}

.admonition.hint {
  --admonition-color: #10b981;
  --admonition-bg: #ecfdf5;
  --admonition-icon: "💡";
}

.admonition.success {
  --admonition-color: #10b981;
  --admonition-bg: #ecfdf5;
  --admonition-icon: "✅";
}

.admonition.warning {
  --admonition-color: #f59e0b;
  --admonition-bg: #fffbeb;
  --admonition-icon: "⚠️";
}

.admonition.caution {
  --admonition-color: #f59e0b;
  --admonition-bg: #fffbeb;
  --admonition-icon: "⚠️";
}

.admonition.danger {
  --admonition-color: #ef4444;
  --admonition-bg: #fef2f2;
  --admonition-icon: "🚨";
}

.admonition.error {
  --admonition-color: #ef4444;
  --admonition-bg: #fef2f2;
  --admonition-icon: "❌";
}

.admonition.bug {
  --admonition-color: #ef4444;
  --admonition-bg: #fef2f2;
  --admonition-icon: "🐛";
}

.admonition.example {
  --admonition-color: #8b5cf6;
  --admonition-bg: #f5f3ff;
  --admonition-icon: "📝";
}

.admonition.quote {
  --admonition-color: #6b7280;
  --admonition-bg: #f9fafb;
  --admonition-icon: "💬";
}

.admonition.abstract {
  --admonition-color: #06b6d4;
  --admonition-bg: #ecfeff;
  --admonition-icon: "📋";
}

/* Dark mode adjustments */
@media (prefers-color-scheme: dark) {
  .admonition.note { --admonition-bg: #1e3a5f; }
  .admonition.info { --admonition-bg: #1e3a5f; }
  .admonition.tip { --admonition-bg: #064e3b; }
  .admonition.hint { --admonition-bg: #064e3b; }
  .admonition.success { --admonition-bg: #064e3b; }
  .admonition.warning { --admonition-bg: #451a03; }
  .admonition.caution { --admonition-bg: #451a03; }
  .admonition.danger { --admonition-bg: #450a0a; }
  .admonition.error { --admonition-bg: #450a0a; }
  .admonition.bug { --admonition-bg: #450a0a; }
  .admonition.example { --admonition-bg: #2e1065; }
  .admonition.quote { --admonition-bg: #1f2937; }
  .admonition.abstract { --admonition-bg: #083344; }
}
```

### Aside (Sidebar) Styles

```css
/* Aside - Sidebar/Marginal Note */
aside.admonition.aside {
  --admonition-color: #6b7280;
  --admonition-bg: #f9fafb;
  border-left: none;
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  font-size: var(--text-sm);
  max-width: 280px;
}

/* Default: float right */
aside.admonition.aside {
  float: right;
  margin: 0 0 var(--space-4) var(--space-4);
  clear: right;
}

/* Inline left modifier */
aside.admonition.aside.aside-inline {
  float: left;
  margin: 0 var(--space-4) var(--space-4) 0;
  clear: left;
}

/* Inline end (right) modifier - explicit */
aside.admonition.aside.aside-inline-end {
  float: right;
  margin: 0 0 var(--space-4) var(--space-4);
  clear: right;
}

/* Aside title styling */
aside.admonition.aside .admonition-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

aside.admonition.aside .admonition-title::before {
  content: none; /* No icon for asides */
}

/* Responsive: full width on small screens */
@media (max-width: 768px) {
  aside.admonition.aside,
  aside.admonition.aside.aside-inline,
  aside.admonition.aside.aside-inline-end {
    float: none;
    max-width: 100%;
    margin: var(--space-4) 0;
  }
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  aside.admonition.aside {
    --admonition-bg: #1f2937;
  }
}
```

### Chat/Conversation Styles

```css
/* Chat container */
.chat-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin: var(--space-4) 0;
  padding: var(--space-4);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
}

/* Individual message */
.chat-message {
  display: flex;
  align-items: flex-end;
  gap: var(--space-2);
  max-width: 85%;
}

/* Left-aligned messages */
.chat-message.chat-left {
  align-self: flex-start;
  flex-direction: row;
}

/* Right-aligned messages */
.chat-message.chat-right {
  align-self: flex-end;
  flex-direction: row-reverse;
}

/* Avatar */
.chat-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: var(--color-primary);
  background-size: cover;
  background-position: center;
  flex-shrink: 0;
}

/* Message bubble */
.chat-bubble {
  padding: var(--space-3);
  border-radius: var(--radius-lg);
  position: relative;
}

/* Left bubble styling */
.chat-left .chat-bubble {
  background: #e5e7eb;
  color: #1f2937;
  border-bottom-left-radius: var(--radius);
}

/* Right bubble styling */
.chat-right .chat-bubble {
  background: var(--color-primary);
  color: white;
  border-bottom-right-radius: var(--radius);
}

/* Author name */
.chat-author {
  font-size: var(--text-xs);
  font-weight: 600;
  margin-bottom: var(--space-1);
  opacity: 0.8;
}

.chat-left .chat-author {
  color: var(--color-text-muted);
}

.chat-right .chat-author {
  color: rgba(255, 255, 255, 0.8);
}

/* Message content */
.chat-content {
  font-size: var(--text-sm);
  line-height: var(--leading-normal);
}

.chat-content p {
  margin: 0;
}

.chat-content p + p {
  margin-top: var(--space-2);
}

/* Timestamp */
.chat-timestamp {
  font-size: var(--text-xs);
  margin-top: var(--space-1);
  opacity: 0.6;
}

.chat-left .chat-timestamp {
  color: var(--color-text-muted);
}

.chat-right .chat-timestamp {
  color: rgba(255, 255, 255, 0.7);
}

/* System messages */
.chat-message.chat-system {
  align-self: center;
  max-width: 100%;
}

.chat-system .chat-content {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  background: transparent;
  padding: var(--space-2) var(--space-4);
  text-align: center;
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  .chat-container {
    background: var(--color-surface);
  }

  .chat-left .chat-bubble {
    background: #374151;
    color: #f9fafb;
  }

  .chat-left .chat-author {
    color: #9ca3af;
  }
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .chat-message {
    max-width: 95%;
  }

  .chat-avatar {
    width: 28px;
    height: 28px;
  }
}
```

### Collapsible Admonitions

```css
/* Collapsible admonitions use <details>/<summary> */
details.admonition {
  border-left: 4px solid var(--admonition-color, var(--color-primary));
}

details.admonition > summary {
  cursor: pointer;
  list-style: none;
}

details.admonition > summary::before {
  content: "▶";
  display: inline-block;
  margin-right: var(--space-2);
  transition: transform 0.2s;
}

details.admonition[open] > summary::before {
  transform: rotate(90deg);
}
```

---

## Code Block Styles

### Base Code Styles

```css
/* Inline code */
code {
  font-family: var(--font-mono);
  font-size: 0.9em;
  padding: 0.125em 0.375em;
  background: var(--color-surface);
  border-radius: var(--radius);
}

/* Code blocks */
pre {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  line-height: normal;
  padding: var(--space-4);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  overflow-x: auto;
}

pre code {
  padding: 0;
  background: none;
}

/* Line numbers (optional) */
pre.line-numbers {
  padding-left: 3.5em;
  position: relative;
}

pre.line-numbers::before {
  content: attr(data-line-numbers);
  position: absolute;
  left: 0;
  top: var(--space-4);
  width: 3em;
  text-align: right;
  color: var(--color-text-muted);
  border-right: 1px solid var(--color-border);
  padding-right: var(--space-2);
  user-select: none;
}

/* Language label */
pre[data-language]::after {
  content: attr(data-language);
  position: absolute;
  top: var(--space-2);
  right: var(--space-2);
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  text-transform: uppercase;
}
```

### Syntax Highlighting Themes

Implementations SHOULD support multiple syntax highlighting themes. Syntax highlighting can be:

1. **Explicitly configured** via `markdown.highlight.theme`
2. **Automatically styled** from the site's palette component colors (`theme.palette`)
3. **Defaulted** to a stable internal Chroma theme when no explicit override is set

#### Palette-native Styling

When no explicit theme is set, syntax highlighting CSS MUST be generated from the active palette's code component roles instead of relying on nearest-match Chroma themes. The generated CSS SHOULD use these palette-driven variables when available:

- `code-bg`
- `code-text`
- `code-comment`
- `code-keyword`
- `code-string`
- `code-number`
- `code-function`
- `code-type`
- `code-operator`

If a palette omits some of these roles, implementations SHOULD fall back to sensible site variables such as `--color-text`, `--color-text-muted`, `--color-link`, `--color-primary`, `--color-success`, and `--color-error`.

Explicit `markdown.highlight.theme` overrides MAY still use any supported Chroma theme.

#### Available Themes

| Theme | Description |
|-------|-------------|
| `github` | GitHub's light theme |
| `github-dark` | GitHub's dark theme |
| `monokai` | Classic dark theme |
| `onedark` | Atom's One Dark |
| `dracula` | Dracula theme |
| `nord` | Nord color palette |
| `solarized-light` | Solarized light |
| `solarized-dark` | Solarized dark |
| `catppuccin-latte` | Catppuccin light |
| `catppuccin-frappe` | Catppuccin medium-dark |
| `catppuccin-macchiato` | Catppuccin dark |
| `catppuccin-mocha` | Catppuccin darkest |
| `rose-pine` | Rosé Pine dark |
| `rose-pine-moon` | Rosé Pine medium |
| `rose-pine-dawn` | Rosé Pine light |
| `gruvbox` | Gruvbox dark |
| `gruvbox-light` | Gruvbox light |
| `tokyonight-night` | Tokyo Night dark |
| `tokyonight-storm` | Tokyo Night storm |
| `tokyonight-day` | Tokyo Night light |
| `vim` | Classic Vim |
| `evergarden` | Nature-inspired |

See https://xyproto.github.io/splash/docs/ for the complete list of Chroma themes.

Configuration:

```toml
# Automatic (derived from theme.palette)
[markata-go.theme]
palette = "catppuccin-mocha"  # Will use catppuccin-mocha for code highlighting

# Explicit override
[markata-go.markdown.highlight]
theme = "dracula"  # Overrides palette-derived theme
line_numbers = false
```

---

## Plugin Output Styles

Built-in themes MUST include styles for plugin-generated HTML elements.

### Heading Anchors (`heading_anchors` plugin)

```css
.heading-anchor {
  opacity: 0;
  margin-left: var(--space-2);
  color: var(--color-text-muted);
  text-decoration: none;
  font-weight: 400;
  transition: opacity 0.2s, color 0.2s;
}

h1:hover .heading-anchor,
h2:hover .heading-anchor,
h3:hover .heading-anchor,
h4:hover .heading-anchor,
h5:hover .heading-anchor,
h6:hover .heading-anchor,
.heading-anchor:focus {
  opacity: 1;
}

.heading-anchor:hover {
  color: var(--color-primary);
  text-decoration: none;
}
```

### Table of Contents (`toc` plugin)

```css
.toc {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  margin: var(--space-6) 0;
}

.toc-title {
  font-size: var(--text-sm);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  margin-bottom: var(--space-3);
}

.toc ul {
  list-style: none;
  padding-left: 0;
  margin: 0;
}

.toc ul ul {
  padding-left: var(--space-4);
  margin-top: var(--space-1);
}

.toc li {
  margin-bottom: var(--space-1);
}

.toc a {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
  text-decoration: none;
  transition: color 0.2s;
}

.toc a:hover {
  color: var(--color-primary);
}
```

### CSV Tables (`csv_fence` plugin)

```css
.csv-table {
  width: 100%;
  border-collapse: collapse;
  margin: var(--space-6) 0;
  font-size: var(--text-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
}

.csv-table th,
.csv-table td {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.csv-table th {
  background: var(--color-surface);
  font-weight: 600;
}

.csv-table tbody tr:nth-child(even) {
  background: var(--color-surface);
}

.csv-table tbody tr:hover {
  background: var(--color-border);
}
```

### Glossary Terms (`glossary` plugin)

```css
.glossary-term {
  color: var(--color-primary);
  text-decoration: underline;
  text-decoration-style: dotted;
  text-underline-offset: 3px;
  cursor: help;
  transition: text-decoration-color 0.2s;
}

.glossary-term:hover {
  text-decoration-style: solid;
  text-decoration-color: var(--color-primary);
}
```

### Mermaid Diagrams (`mermaid` plugin)

Client-rendered Mermaid diagrams follow the active site palette when the theme
or palette changes. Palette values may use modern CSS expressions such as
`color-mix()`, but the Mermaid theme API must receive concrete colors supported
by the Mermaid version in use. Theme changes must rerender diagrams without
discarding their original source, and server-rendered diagram markup must remain
readable when JavaScript is disabled.

```css
.mermaid {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  margin: var(--space-6) 0;
  text-align: center;
  overflow-x: auto;
}
```

### Wikilinks (`wikilinks` plugin)

```css
a.wikilink {
  color: var(--color-primary);
  text-decoration: underline dotted;
  text-underline-offset: 0.18em;
}

a.wikilink.wikilink-missing {
  color: var(--color-error);
  text-decoration: underline;
  text-decoration-style: wavy;
}
```

### Reading Time (`reading_time` plugin)

```css
.reading-time {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
```

---

## Template Requirements

### Base Template (`base.html`)

Every theme MUST provide a base template with these blocks:

```jinja2
<!DOCTYPE html>
<html lang="{{ config.lang | default('en') }}">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">

  {% block meta %}
  <title>{% block title %}{{ config.title }}{% endblock %}</title>
  <meta name="description" content="{% block description %}{{ config.description }}{% endblock %}">
  {% endblock %}

  {% block head %}
  {# Theme CSS #}
  <link rel="stylesheet" href="{{ 'css/main.css' | theme_asset }}">

  {# CSS variable overrides from config #}
  {% if config.theme.variables %}
  <style>:root { {% for k, v in config.theme.variables.items() %}{{ k }}: {{ v }}; {% endfor %} }</style>
  {% endif %}

  {# Custom CSS #}
  {% if config.theme.custom_css %}
  <link rel="stylesheet" href="{{ config.theme.custom_css | asset_url }}">
  {% endif %}
  {% endblock %}
</head>
<body>
  {% block header %}
  {% include "partials/header.html" %}
  {% endblock %}

  <main>
    {% block content %}{% endblock %}
  </main>

  {% block footer %}
  {% include "partials/footer.html" %}
  {% endblock %}

  {% block scripts %}{% endblock %}
</body>
</html>
```

### Post Template (`post.html`)

```jinja2
{% extends "base.html" %}

{% block title %}{{ post.title }} | {{ config.title }}{% endblock %}
{% block description %}{{ post.description | default(post.content | truncate(160)) }}{% endblock %}

{% block meta %}
{{ super() }}
<meta property="og:title" content="{{ post.title }}">
<meta property="og:type" content="article">
<meta property="og:url" content="{{ post.absolute_url }}">
{% if post.date %}
<meta property="article:published_time" content="{{ post.date | isoformat }}">
{% endif %}
{% endblock %}

{% block content %}
<article class="post">
  <header class="post-header">
    <h1>{{ post.title }}</h1>
    {% if post.date or post.reading_time %}
    <div class="post-meta">
      {% if post.date %}<time datetime="{{ post.date | isoformat }}">{{ post.date | date }}</time>{% endif %}
      {% if post.reading_time %}<span>{{ post.reading_time }} min read</span>{% endif %}
    </div>
    {% endif %}
  </header>

  <div class="post-content">
    {{ body | safe }}
  </div>

  {% if post.tags %}
  <footer class="post-footer">
    <div class="tags">
      {% for tag in post.tags %}
      <a href="/tags/{{ tag | slugify }}/" class="tag">{{ tag }}</a>
      {% endfor %}
    </div>
  </footer>
  {% endif %}
</article>

{% if post.prev or post.next %}
<nav class="post-nav">
  {% if post.prev %}<a href="{{ post.prev.href }}" class="prev">← {{ post.prev.title }}</a>{% endif %}
  {% if post.next %}<a href="{{ post.next.href }}" class="next">{{ post.next.title }} →</a>{% endif %}
</nav>
{% endif %}
{% endblock %}
```

### Feed Template (`feed.html`)

Feed pages render in the feed-only page wrapper state. When `feed` is present in
template context, the base layout does not render content or doc sidebars, so
feed templates must opt into the centered feed wrapper explicitly.

```jinja2
{% extends "base.html" %}

{% block title %}{{ feed.title | default(config.title) }}{% endblock %}

{% block content %}
<div class="feed">
  {% if feed.title %}
  <header class="feed-header">
    <h1>{{ feed.title }}</h1>
    {% if feed.description %}<p>{{ feed.description }}</p>{% endif %}
  </header>
  {% endif %}

  <div class="posts">
    {% for post in feed.posts %}
    {% include "card.html" %}
    {% endfor %}
  </div>

  {% if feed.pagination.total_pages > 1 %}
  {% include "partials/pagination.html" %}
  {% endif %}
</div>
{% endblock %}
```

### Card Template (`card.html`)

```jinja2
<article class="card">
  <h2><a href="{{ post.href }}">{{ post.title }}</a></h2>
  {% if post.description %}
  <p class="card-description">{{ post.description }}</p>
  {% endif %}
  <div class="card-meta">
    {% if post.date %}<time datetime="{{ post.date | isoformat }}">{{ post.date | date }}</time>{% endif %}
    {% if post.reading_time %}<span>{{ post.reading_time }} min</span>{% endif %}
  </div>
</article>
```

---

## Template Filters

Themes should have access to these filters:

| Filter | Description | Example |
|--------|-------------|---------|
| `theme_asset` | URL to theme static file | `{{ 'css/main.css' \| theme_asset }}` |
| `asset_url` | URL to project static file | `{{ 'images/logo.png' \| asset_url }}` |
| `date` | Format date | `{{ post.date \| date('%B %d, %Y') }}` |
| `isoformat` | ISO 8601 date | `{{ post.date \| isoformat }}` |
| `truncate` | Truncate text | `{{ text \| truncate(160) }}` |
| `slugify` | Generate slug | `{{ title \| slugify }}` |
| `safe` | Mark as safe HTML | `{{ body \| safe }}` |

---

## Installing Themes

### From Package Manager

```bash
# Python
pip install [name]-theme-blog

# npm
npm install @[name]/theme-blog

# Direct download
[name] theme install https://github.com/user/theme-blog
```

### Manual Installation

Download theme to `~/.config/[name]/themes/` or `./themes/`:

```bash
mkdir -p themes
git clone https://github.com/user/theme-blog themes/blog
```

### Configuration

```toml
[markata-go.theme]
name = "blog"  # Matches directory name
```

---

## Creating a Theme

### Scaffold a New Theme

```bash
[name] theme new my-theme
```

Creates:
```
themes/my-theme/
├── theme.toml
├── templates/
│   ├── base.html
│   ├── post.html
│   ├── feed.html
│   └── card.html
└── static/
    └── css/
        ├── main.css
        └── variables.css
```

### Extending a Theme

Inherit from an existing theme and override only what you need:

```toml
# themes/my-theme/theme.toml
[theme]
name = "My Theme"

[theme.extends]
theme = "default"
```

Then only include the files you want to override:

```
themes/my-theme/
├── theme.toml
├── templates/
│   └── partials/
│       └── footer.html    # Custom footer only
└── static/
    └── css/
        └── variables.css  # Custom colors only
```

---

## CLI Commands

### `theme list`

List available themes:

```bash
$ markata-go theme list

Installed themes:
  default     Clean, minimal theme with dark mode (built-in)
  minimal     Bare-bones HTML for maximum customization (built-in)
  blog        Feature-rich blog theme (~/.config/markata-go/themes/blog)

Current theme: default
```

### `theme info`

Show theme details:

```bash
$ markata-go theme info blog

Name: Blog Theme
Version: 2.1.0
Author: Theme Author
License: MIT
Homepage: https://github.com/user/theme-blog

Features:
  ✓ Dark mode
  ✓ Responsive
  ✓ Syntax highlighting
  ✓ Admonitions

Options:
  primary_color   Color   #3b82f6   Primary accent color
  font_family     String  system-ui Body font family
  show_toc        Boolean true      Show table of contents
```

### `theme install`

Install a theme:

```bash
$ markata-go theme install https://github.com/user/theme-blog
Installing theme from https://github.com/user/theme-blog...
Theme 'blog' installed to ~/.config/markata-go/themes/blog
```

### `theme new`

Create a new theme:

```bash
$ markata-go theme new my-theme
Created theme scaffold in themes/my-theme/
```

### `palette list`

List available color palettes from all sources.

**Usage:**
```bash
markata-go palette list [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--variant` | `-v` | Filter by variant (`light` or `dark`) |
| `--json` | | Output as JSON |

**Output:**
```bash
$ markata-go palette list

Built-in palettes:
  default-light         Light   Clean, minimal light theme
  default-dark          Dark    Clean, minimal dark theme
  catppuccin-latte      Light   Soothing pastel theme
  catppuccin-frappe     Dark    Medium contrast dark
  catppuccin-macchiato  Dark    Higher contrast dark
  catppuccin-mocha      Dark    Highest contrast dark
  nord-light            Light   Arctic, north-bluish
  nord-dark             Dark    Arctic, north-bluish
  gruvbox-light         Light   Retro groove, warm
  gruvbox-dark          Dark    Retro groove, warm
  dracula               Dark    Vibrant on dark
  rose-pine             Dark    Soho vibes, natural
  tokyo-night           Dark    Tokyo city lights
  solarized-dark        Dark    Scientifically designed

User palettes:
  my-brand              Light   ./palettes/my-brand.toml

Current: catppuccin-mocha (from markata-go.toml)
```

**Filtered output:**
```bash
$ markata-go palette list --variant light

Light palettes:
  default-light         Clean, minimal light theme
  catppuccin-latte      Soothing pastel theme
  nord-light            Arctic, north-bluish
  gruvbox-light         Retro groove, warm
  my-brand              ./palettes/my-brand.toml
```

**JSON output:**
```bash
$ markata-go palette list --json
```
```json
{
  "current": "catppuccin-mocha",
  "palettes": [
    {
      "name": "catppuccin-mocha",
      "variant": "dark",
      "description": "Highest contrast dark",
      "source": "built-in",
      "author": "Catppuccin Team"
    }
  ]
}
```

**Implementation notes:**
- Scan directories in order: `./palettes/`, `~/.config/markata-go/palettes/`, built-in
- Read `[palette]` section from each TOML file for metadata
- Current palette determined from `markata-go.toml` config

---

### `palette pick`

Interactively browse and select a palette using a full-screen TUI fuzzy picker. Sets the chosen palette in the project config by default.

**Usage:**
```bash
markata-go palette pick [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--no-set` | Only print the palette name without updating the config file |

**Behavior:**

Opens a full-screen alternate-screen TUI with two panels:

- **Left panel:** Scrollable, fuzzy-filterable list of all discovered palettes (built-in, user, project). Each entry shows the palette name and a variant badge (`[dark]`/`[light]`). A search prompt at the top accepts typed input for fuzzy matching. Scroll indicators appear when the list overflows.

- **Right panel:** Live preview of the currently highlighted palette showing:
  - Palette metadata (name, author, description, source)
  - Raw color swatch grid (alphabetically sorted `"██"` blocks with labels)
  - Semantic color swatches for key roles (text, bg, accent, link, success, warning, error, border)
  - Contrast preview block rendering sample text with `bg-primary` background and `text-primary` foreground, plus underlined link text

**Key bindings:**
| Key | Action |
|-----|--------|
| Type characters | Fuzzy-filter the palette list |
| Backspace | Delete last filter character |
| Up / Down | Move cursor |
| Page Up / Page Down | Scroll by page |
| Enter | Select the highlighted palette |
| Esc / Ctrl+C | Cancel without selecting |

**Default behavior (sets config):**
```bash
$ markata-go palette pick
# (interactive TUI opens)
# On Enter, prints the palette name and updates config:
catppuccin-mocha
Set palette to "catppuccin-mocha" in markata-go.toml
```

On selection, the command discovers the project config file (or uses `--config`), parses it, sets `theme.palette` to the chosen name, and writes it back. The palette name is printed to stdout and the status message to stderr.

**With `--no-set` (echo only):**
```bash
$ markata-go palette pick --no-set
# (interactive TUI opens)
# On Enter, only prints the selected palette name to stdout:
catppuccin-mocha
```

**Composability:**
```bash
# Use pick output as input to other commands (skip config write)
markata-go palette info "$(markata-go palette pick --no-set)"

# Pick, set, and rebuild
markata-go palette pick && markata-go build
```

**Implementation notes:**
- Uses Bubble Tea (`tea.NewProgram` with `tea.WithAltScreen()`) for the TUI
- Uses Lipgloss for styled color swatch rendering
- Pre-loads all palettes at startup for instant preview switching
- Fuzzy matching checks that all characters in the query appear in the palette name in order
- Config writing reuses the `config set` infrastructure (`parseConfigToMap`)
- Palette discovery via `palettes.NewLoader().Discover()`, loading via `loader.Load(name)`

---

### `palette info`

Show detailed information about a specific palette.

**Usage:**
```bash
markata-go palette info <name> [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Output as JSON |
| `--colors` | `-c` | Show only raw colors |
| `--semantic` | `-s` | Show only semantic mappings |

**Output:**
```bash
$ markata-go palette info catppuccin-mocha

Name:     Catppuccin Mocha
Variant:  dark
Author:   Catppuccin Team
License:  MIT
Homepage: https://catppuccin.com
Source:   [built-in]

Raw Colors (26):
  rosewater  #f5e0dc    flamingo   #f2cdcd    pink       #f5c2e7
  mauve      #cba6f7    red        #f38ba8    maroon     #eba0ac
  peach      #fab387    yellow     #f9e2af    green      #a6e3a1
  teal       #94e2d5    sky        #89dceb    sapphire   #74c7ec
  blue       #89b4fa    lavender   #b4befe    text       #cdd6f4
  subtext1   #bac2de    subtext0   #a6adc8    overlay2   #9399b2
  overlay1   #7f849c    overlay0   #6c7086    surface2   #585b70
  surface1   #45475a    surface0   #313244    base       #1e1e2e
  mantle     #181825    crust      #11111b

Semantic Mappings:
  text-primary    -> text (#cdd6f4)
  text-secondary  -> subtext1 (#bac2de)
  text-muted      -> overlay1 (#7f849c)
  bg-primary      -> base (#1e1e2e)
  bg-secondary    -> mantle (#181825)
  bg-surface      -> surface0 (#313244)
  bg-elevated     -> surface1 (#45475a)
  accent          -> mauve (#cba6f7)
  accent-hover    -> lavender (#b4befe)
  link            -> blue (#89b4fa)
  link-hover      -> sapphire (#74c7ec)
  link-visited    -> mauve (#cba6f7)
  success         -> green (#a6e3a1)
  warning         -> yellow (#f9e2af)
  error           -> red (#f38ba8)
  info            -> blue (#89b4fa)
  border          -> surface1 (#45475a)
  border-focus    -> mauve (#cba6f7)

Component Colors:
  code-bg         -> surface0 (#313244)
  code-text       -> text (#cdd6f4)
  code-keyword    -> mauve (#cba6f7)
  code-string     -> green (#a6e3a1)
  code-function   -> blue (#89b4fa)
  ...
```

**Implementation notes:**
- Load palette TOML file
- Resolve all semantic references to actual hex values
- Display in three columns for raw colors (fits 80-char terminal)

---

### `palette check`

Validate palette contrast ratios against WCAG 2.1 accessibility guidelines.

**Usage:**
```bash
markata-go palette check [name] [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--strict` | | Require AAA compliance (7:1 for text) |
| `--all` | `-a` | Check all installed palettes |
| `--json` | | Output as JSON |
| `--report` | `-r` | Generate detailed HTML report |
| `--fix` | | Suggest color adjustments to pass |

**Basic check:**
```bash
$ markata-go palette check catppuccin-mocha

Checking palette: Catppuccin Mocha

Text Contrast (AA requires 4.5:1):
  text-primary on bg-primary:      11.79:1  [PASS] AA, AAA
  text-secondary on bg-primary:     8.21:1  [PASS] AA, AAA
  text-muted on bg-primary:         4.08:1  [PASS] AA Large (3:1)
  text-primary on bg-surface:       9.42:1  [PASS] AA, AAA
  text-primary on bg-elevated:      7.51:1  [PASS] AA, AAA

Interactive Elements (AA requires 4.5:1, UI 3:1):
  link on bg-primary:               6.41:1  [PASS] AA, AAA
  link-hover on bg-primary:         5.87:1  [PASS] AA
  accent on bg-primary:             5.18:1  [PASS] AA
  border-focus on bg-primary:       5.18:1  [PASS] UI

Status Colors (UI requires 3:1):
  success on bg-primary:            7.28:1  [PASS]
  warning on bg-primary:            9.14:1  [PASS]
  error on bg-primary:              5.84:1  [PASS]
  info on bg-primary:               6.41:1  [PASS]

Code Blocks (AA requires 4.5:1):
  code-text on code-bg:             9.42:1  [PASS] AA, AAA
  code-comment on code-bg:          3.25:1  [PASS] AA Large
  code-keyword on code-bg:          4.12:1  [PASS] AA Large
  code-string on code-bg:           5.79:1  [PASS] AA
  code-function on code-bg:         5.10:1  [PASS] AA

Buttons:
  button-primary-text on button-primary-bg:    11.79:1  [PASS] AA, AAA
  button-secondary-text on button-secondary-bg: 9.42:1  [PASS] AA, AAA

Summary: 19/19 checks passed
WCAG 2.1 Level AA: COMPLIANT
```

**Strict mode (AAA):**
```bash
$ markata-go palette check catppuccin-mocha --strict

Checking palette: Catppuccin Mocha (AAA strict mode)

Text Contrast (AAA requires 7:1):
  text-primary on bg-primary:      11.79:1  [PASS] AAA
  text-secondary on bg-primary:     8.21:1  [PASS] AAA
  text-muted on bg-primary:         4.08:1  [FAIL] AAA requires 4.5:1
  ...

Summary: 15/19 checks passed
WCAG 2.1 Level AAA: NOT COMPLIANT

Failed checks:
  - text-muted on bg-primary: 4.08:1 (need 4.5:1, adjust to #8a8faa)
  - code-comment on code-bg: 3.25:1 (need 4.5:1, adjust to #9399b2)
```

**Check all palettes:**
```bash
$ markata-go palette check --all

Checking all installed palettes...

  default-light         [PASS] AA compliant (19/19)
  default-dark          [PASS] AA compliant (19/19)
  catppuccin-latte      [PASS] AA compliant (19/19)
  catppuccin-mocha      [PASS] AA compliant (19/19)
  nord-dark             [PASS] AA compliant (19/19)
  gruvbox-dark          [PASS] AA compliant (19/19)
  dracula               [WARN] 1 marginal (18/19, code-comment: 3.01:1)
  rose-pine             [PASS] AA compliant (19/19)
  tokyo-night           [PASS] AA compliant (19/19)
  solarized-dark        [PASS] AA compliant (19/19)
  my-brand              [FAIL] 3 failures (16/19)

Summary: 10/11 palettes fully AA compliant
         1 palette has failures (run with palette name for details)
```

**Generate report:**
```bash
$ markata-go palette check catppuccin-mocha --report

=== Accessibility Report: Catppuccin Mocha ===

WCAG 2.1 Level AA Compliance: PASS

Text Readability:
  Excellent (7:1+):  text-primary, text-secondary
  Good (4.5:1+):     link, accent, status colors
  Acceptable (3:1+): text-muted, code-comment

Interactive Elements:
  [PASS] Links are distinguishable from text
  [PASS] Focus indicators meet 3:1 contrast
  [PASS] Buttons have sufficient contrast

Color Blindness Simulation:
  [PASS] Protanopia:   Status colors distinguishable
  [PASS] Deuteranopia: Status colors distinguishable
  [PASS] Tritanopia:   Status colors distinguishable

Recommendations:
  - Consider increasing text-muted contrast for AAA compliance
  - code-comment is marginal; consider lightening to #9399b2

Report saved: .markata/palette-report-catppuccin-mocha.html
```

**JSON output:**
```bash
$ markata-go palette check catppuccin-mocha --json
```
```json
{
  "palette": "catppuccin-mocha",
  "compliant": true,
  "level": "AA",
  "checks": [
    {
      "foreground": "text-primary",
      "background": "bg-primary",
      "fg_color": "#cdd6f4",
      "bg_color": "#1e1e2e",
      "ratio": 11.79,
      "required": 4.5,
      "level": "AA",
      "passed": true,
      "grades": ["AA", "AAA"]
    }
  ],
  "summary": {
    "total": 19,
    "passed": 19,
    "failed": 0
  }
}
```

**Implementation notes:**
- Use WCAG 2.1 relative luminance formula (see Contrast Ratio Calculation section)
- Check all combinations defined in RequiredChecks
- `--fix` uses binary search to find minimum luminance adjustment while preserving hue

**Exit codes:**

| Scenario | Exit Code |
|----------|-----------|
| All checks pass | 0 |
| Any AA failure (default mode) | 1 |
| AA passes, AAA fails (no `--strict`) | 0 |
| AA passes, AAA fails (with `--strict`) | 1 |
| `--all`: All palettes pass | 0 |
| `--all`: Any palette fails | 1 |

---

### `palette preview`

Generate an HTML preview page showing the palette applied to sample content.

**Usage:**
```bash
markata-go palette preview <name> [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file path (default: `.markata/palette-preview-{name}.html`) |
| `--open` | | Open in default browser after generating |
| `--compare` | `-c` | Compare with another palette side-by-side |

**Basic preview:**
```bash
$ markata-go palette preview catppuccin-mocha

Generated: .markata/palette-preview-catppuccin-mocha.html

Preview includes:
  - Color swatches with names and hex values
  - Typography samples (headings, body, code)
  - Code block with syntax highlighting
  - All admonition types
  - Button and form components
  - Card and navigation examples

Open in browser to view.
```

**Open automatically:**
```bash
$ markata-go palette preview catppuccin-mocha --open

Generated: .markata/palette-preview-catppuccin-mocha.html
Opening in default browser...
```

**Compare two palettes:**
```bash
$ markata-go palette preview catppuccin-mocha --compare catppuccin-latte

Generated: .markata/palette-compare-catppuccin-mocha-vs-catppuccin-latte.html

Side-by-side comparison of:
  - catppuccin-mocha (dark)
  - catppuccin-latte (light)
```

**Preview HTML structure:**

The generated HTML includes:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <title>Palette Preview: Catppuccin Mocha</title>
  <style>
    :root {
      /* All palette CSS variables injected here */
      --palette-base: #1e1e2e;
      --color-text-primary: var(--palette-text);
      /* ... */
    }
  </style>
</head>
<body>
  <header>
    <h1>Catppuccin Mocha</h1>
    <p>Dark variant by Catppuccin Team</p>
  </header>

  <section id="colors">
    <h2>Raw Colors</h2>
    <!-- Color swatches grid -->
  </section>

  <section id="typography">
    <h2>Typography</h2>
    <!-- Heading samples, body text, links -->
  </section>

  <section id="code">
    <h2>Code Blocks</h2>
    <!-- Sample code with syntax highlighting -->
  </section>

  <section id="admonitions">
    <h2>Admonitions</h2>
    <!-- All admonition types -->
  </section>

  <section id="components">
    <h2>Components</h2>
    <!-- Buttons, cards, forms -->
  </section>

  <section id="contrast">
    <h2>Contrast Ratios</h2>
    <!-- Visual contrast check results -->
  </section>
</body>
</html>
```

**Implementation notes:**
- Use embedded HTML template
- Inject palette as CSS custom properties
- Include sample content demonstrating all components
- For `--compare`, use CSS grid with two columns

---

### `palette new`

Create a new custom palette from scratch or based on an existing one.

**Usage:**
```bash
markata-go palette new <name> [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--variant` | `-v` | Variant type: `light` or `dark` (default: `dark`) |
| `--from` | `-f` | Base palette to copy from |
| `--output` | `-o` | Output directory (default: `./palettes/`) |
| `--minimal` | `-m` | Create minimal palette (semantic only, no components) |

**Create from scratch:**
```bash
$ markata-go palette new my-brand --variant light

Created: palettes/my-brand.toml

The palette has been initialized with placeholder colors.
Edit the file to customize your brand colors.

Next steps:
  1. Edit palettes/my-brand.toml
  2. Validate: markata-go palette check my-brand
  3. Preview:  markata-go palette preview my-brand --open
  4. Use:      Add 'palette = "my-brand"' to markata-go.toml
```

**Create from existing palette:**
```bash
$ markata-go palette new my-brand --from catppuccin-mocha

Created: palettes/my-brand.toml (based on catppuccin-mocha)

The palette has been copied from catppuccin-mocha.
Customize the colors to match your brand.

Changes from base:
  - Updated [palette] metadata (name, author)
  - All colors preserved - edit as needed
```

**Create minimal palette:**
```bash
$ markata-go palette new my-brand --minimal --variant dark

Created: palettes/my-brand.toml (minimal)

Minimal palette with only semantic colors.
Component colors will inherit from semantic defaults.
```

**Generated palette structure:**

Full palette (default):
```toml
# palettes/my-brand.toml
# Generated by markata-go palette new

[palette]
name = "My Brand"
variant = "dark"
author = "Your Name"
license = "MIT"
homepage = ""

# Layer 1: Raw Colors
# Define your brand colors here
[palette.colors]
# Primary brand colors
primary     = "#3b82f6"   # TODO: Replace with your primary color
secondary   = "#8b5cf6"   # TODO: Replace with your secondary color

# Neutral colors (dark variant)
gray-900    = "#111827"
gray-800    = "#1f2937"
gray-700    = "#374151"
gray-600    = "#4b5563"
gray-500    = "#6b7280"
gray-400    = "#9ca3af"
gray-300    = "#d1d5db"
gray-200    = "#e5e7eb"
gray-100    = "#f3f4f6"

# Status colors
green       = "#10b981"
yellow      = "#f59e0b"
red         = "#ef4444"
blue        = "#3b82f6"

# Layer 2: Semantic Mapping
[palette.semantic]
text-primary   = "gray-100"
text-secondary = "gray-300"
text-muted     = "gray-500"

bg-primary   = "gray-900"
bg-secondary = "gray-800"
bg-surface   = "gray-800"
bg-elevated  = "gray-700"

accent       = "primary"
accent-hover = "secondary"
link         = "primary"
link-hover   = "secondary"
link-visited = "secondary"

success = "green"
warning = "yellow"
error   = "red"
info    = "blue"

border       = "gray-700"
border-focus = "primary"

# Layer 3: Component Colors (optional)
[palette.components]
code-bg       = "gray-800"
code-text     = "gray-100"
code-comment  = "gray-500"
code-keyword  = "secondary"
code-string   = "green"
code-number   = "primary"
code-function = "primary"

# Add more component colors as needed
```

Minimal palette (`--minimal`):
```toml
# palettes/my-brand.toml
# Minimal palette - component colors inherit from semantic

[palette]
name = "My Brand"
variant = "dark"
author = "Your Name"

[palette.colors]
primary   = "#3b82f6"
secondary = "#8b5cf6"
bg        = "#111827"
surface   = "#1f2937"
text      = "#f3f4f6"
muted     = "#6b7280"
green     = "#10b981"
yellow    = "#f59e0b"
red       = "#ef4444"

[palette.semantic]
text-primary   = "text"
text-secondary = "text"
text-muted     = "muted"
bg-primary     = "bg"
bg-surface     = "surface"
accent         = "primary"
link           = "primary"
success        = "green"
warning        = "yellow"
error          = "red"
info           = "primary"
border         = "muted"
```

**Implementation notes:**
- Create `./palettes/` directory if it doesn't exist
- Use template with TODO comments for guidance
- `--from` copies and modifies metadata
- Validate name doesn't conflict with built-in palettes

---

### `palette export`

Export a palette to different formats.

**Usage:**
```bash
markata-go palette export <name> [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--format` | `-f` | Output format: `css`, `scss`, `json`, `tailwind` |
| `--output` | `-o` | Output file (default: stdout) |

**Export as CSS:**
```bash
$ markata-go palette export catppuccin-mocha --format css

:root {
  /* Raw palette colors */
  --palette-rosewater: #f5e0dc;
  --palette-flamingo: #f2cdcd;
  --palette-pink: #f5c2e7;
  --palette-mauve: #cba6f7;
  /* ... */

  /* Semantic colors */
  --color-text-primary: var(--palette-text);
  --color-text-secondary: var(--palette-subtext1);
  --color-bg-primary: var(--palette-base);
  --color-accent: var(--palette-mauve);
  /* ... */

  /* Component colors */
  --code-bg: var(--palette-surface0);
  --code-text: var(--palette-text);
  /* ... */
}
```

**Export as Tailwind config:**
```bash
$ markata-go palette export catppuccin-mocha --format tailwind

module.exports = {
  theme: {
    extend: {
      colors: {
        'palette': {
          'rosewater': '#f5e0dc',
          'flamingo': '#f2cdcd',
          'pink': '#f5c2e7',
          'mauve': '#cba6f7',
          // ...
        },
        'text': {
          'primary': '#cdd6f4',
          'secondary': '#bac2de',
          'muted': '#7f849c',
        },
        'bg': {
          'primary': '#1e1e2e',
          'secondary': '#181825',
          'surface': '#313244',
        },
        // ...
      }
    }
  }
}
```

**Export as JSON:**
```bash
$ markata-go palette export catppuccin-mocha --format json
```
```json
{
  "name": "Catppuccin Mocha",
  "variant": "dark",
  "colors": {
    "rosewater": "#f5e0dc",
    "flamingo": "#f2cdcd",
    "pink": "#f5c2e7",
    "mauve": "#cba6f7",
    "red": "#f38ba8",
    "text": "#cdd6f4",
    "base": "#1e1e2e"
  },
  "semantic": {
    "text-primary": "#cdd6f4",
    "text-secondary": "#bac2de",
    "text-muted": "#7f849c",
    "bg-primary": "#1e1e2e",
    "bg-secondary": "#181825",
    "accent": "#cba6f7",
    "link": "#89b4fa",
    "link-hover": "#74c7ec",
    "success": "#a6e3a1",
    "warning": "#f9e2af",
    "error": "#f38ba8"
  },
  "components": {
    "code-bg": "#313244",
    "code-text": "#cdd6f4",
    "code-keyword": "#cba6f7",
    "button-primary-bg": "#cba6f7",
    "button-primary-text": "#1e1e2e"
  }
}
```

**Export as SCSS:**
```bash
$ markata-go palette export catppuccin-mocha --format scss -o _palette.scss

// _palette.scss - Generated from catppuccin-mocha

// Raw colors
$palette-rosewater: #f5e0dc;
$palette-flamingo: #f2cdcd;
// ...

// Semantic colors
$color-text-primary: $palette-text;
$color-bg-primary: $palette-base;
// ...

// As a map for programmatic access
$palette: (
  'rosewater': #f5e0dc,
  'flamingo': #f2cdcd,
  // ...
);
```

**Implementation notes:**
- Read palette TOML
- Transform to target format
- Resolve references for formats that don't support variables

---

## See Also

- [TEMPLATES.md](./TEMPLATES.md) - Template system details
- [CONFIG.md](./CONFIG.md) - Theme configuration
- [CONTENT.md](./CONTENT.md) - Admonition syntax
- [SPEC.md](./SPEC.md) - Core specification

### Deferred palette catalog

The default theme's initial palette-core stylesheet MUST include only configured
light/dark palette metadata and colors. The full palette stylesheet MUST retain
the complete catalog for custom templates and lazy picker loading. If the
configured pair cannot be resolved, retain the full manifest as a fallback.
The picker MUST refresh a previously cached core manifest after full palette
styles become ready, so every color remains available on the first opening.

Generated palette hashes participate in HTML cache invalidation after all
Configure hooks finish. Consecutive full builds must update cached archive
stylesheet URLs when generated colors or metadata change, even if source
content, templates, and static files are unchanged. Unchanged generated
assets preserve warm page and feed caches.

Font preloads include each distinct body, display, and heading source once.
A separate display role must not suppress the heading face used by ordinary
page titles. Code and extended coverage remain demand-loaded. The configured
tier selection and returning visitor font-pack choice are preserved.

Initial title fitting and heading contour measurements use ResizeObserver
delivery when available, allowing the browser's first layout to complete before
measurement. Font-ready subscriptions begin after the initial rendering
opportunity. Font loading, later resizing, heading text changes, and article
navigation retain the same fitted sizes and contour geometry. Browsers without
ResizeObserver retain the animation-frame measurement fallback.

The article progress indicator also uses ResizeObserver for its initial
measurement and to refresh after article height changes. Reinitializing after
navigation disconnects the prior observer. Scroll and viewport resize updates
remain animation-frame throttled, with the existing fallback when observation
is unavailable.

Closed overlay drawers must not measure header or active navigation geometry
during startup. Opening a drawer centers the active feed item. Full-height overlay drawers
need no header offset; desktop offsets are scoped to drawer elements, preserving
desktop persistence and resize/scroll behavior without restyling article trees.
An overlay drawer focuses its active navigation link when available, otherwise
its first interactive content, and Escape returns focus to the drawer handle.

Where content visibility and discrete transitions are supported, closed feed
and TOC drawers skip rendering their contents. Opening restores rendering
immediately before focus and active-item measurements. Closing keeps contents
rendered through the slide. Saved-open desktop drawers render immediately.
Reduced motion and startup restoration disable the content transition with the
drawer transition. Other browsers retain the existing drawer rendering.
