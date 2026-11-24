---
title: "Web Awesome Components"
description: "Use Web Awesome components in Markdown, including image comparisons with vendored or CDN assets."
date: 2026-04-28
published: true
slug: /docs/guides/webawesome/
tags:
  - documentation
  - markdown
  - webawesome
---


markata-go can render [Web Awesome](https://webawesome.com/) components from Markdown. The built-in shortcuts focus on content patterns that make sense for blogs and docs: comparisons, expandable details, tabs, copy buttons, QR codes, badges, tags, tooltips, carousels, and animated images.

## Demo

::: wa-comparison {position=42 caption="Compare the same generated graphic before and after a color treatment."}
![Before: a grayscale abstract generated graphic](https://images.unsplash.com/photo-1557683316-973673baf926?auto=format&fit=crop&w=1200&q=80&sat=-100)
![After: a colorful abstract generated graphic](https://images.unsplash.com/photo-1557683316-973673baf926?auto=format&fit=crop&w=1200&q=80)
:::

## Theme Integration

Web Awesome components inherit markata-go theme variables for text, surfaces, borders, radii, shadows, and semantic colors. This keeps component chrome aligned with the active palette instead of falling back to Web Awesome defaults.

## Asset Loading

Web Awesome uses vendored assets by default (`source = "vendor"`). In that mode, markata-go downloads the Web Awesome npm archive, extracts the browser-ready `dist-cdn` subtree, and serves it from `/assets/vendor/webawesome/`.

`markata-go assets download` includes this archive in the default asset cache, so builder images can pre-seed Web Awesome and run with `MARKATA_GO_OFFLINE=true` without reaching the npm registry during a site build.

## Markdown

### Comparison

Write a `wa-comparison` container with two images. Author them in natural left-to-right order: the first image starts on the left side of the divider and the second image starts on the right.

```markdown
::: wa-comparison {position=42 caption="Compare the same graphic before and after a color treatment."}
![Before: grayscale graphic](/images/graphic-before.webp)
![After: colorful graphic](/images/graphic-after.webp)
:::
```

### Details

Use details for FAQs, spoilers, and advanced notes.

::: wa-details {summary="When should I self-host assets?"}
Self-host Web Awesome when your site must build without external CDN dependencies or when your privacy policy avoids third-party asset requests.
:::

```markdown
::: wa-details {summary="When should I self-host assets?"}
Self-host Web Awesome when your site must build without external CDN dependencies.
:::
```

### Tabs

Use tabs for platform-specific install instructions or examples in multiple languages.

<wa-tab-group>
  <wa-tab slot="nav" panel="macos">macOS</wa-tab>
  <wa-tab slot="nav" panel="linux">Linux</wa-tab>
  <wa-tab-panel name="macos">

```bash
brew install markata-go
```

  </wa-tab-panel>
  <wa-tab-panel name="linux">

```bash
curl -fsSL https://example.com/install.sh | sh
```

  </wa-tab-panel>
</wa-tab-group>

````markdown
::: wa-tabs

:::: wa-tab {label="macOS"}
```bash
brew install markata-go
```
::::

:::: wa-tab {label="Linux"}
```bash
curl -fsSL https://example.com/install.sh | sh
```
::::
:::
````

### Copy Button

Use copy buttons next to commands, URLs, or config snippets.

::: wa-copy
`go test ./...`
:::

```markdown
::: wa-copy
`go test ./...`
:::
```

### QR Code

Use QR codes for shareable article URLs, event links, or contact links.

::: wa-qr {label="Web Awesome guide QR code"}
https://example.com/docs/guides/webawesome/
:::

```markdown
::: wa-qr
https://example.com/docs/guides/webawesome/
:::
```

### Badge And Tag

Use badges and tags for compact labels like new, stable, deprecated, or experimental.

::: wa-badge {variant="brand"}
New
:::

::: wa-tag {variant="success"}
Stable
:::

```markdown
::: wa-badge {variant="brand"}
New
:::

::: wa-tag {variant="success"}
Stable
:::
```

### Tooltip

Use tooltips for short inline explanations.

::: wa-tooltip {content="Static Site Generator"}
SSG
:::

```markdown
::: wa-tooltip {content="Static Site Generator"}
SSG
:::
```

### Carousel

Use carousels for screenshot galleries or visual changelogs.

::: wa-carousel {navigation="true" pagination="true"}
![First screenshot](https://images.unsplash.com/photo-1498050108023-c5249f4df085?auto=format&fit=crop&w=1200&q=80)
![Second screenshot](https://images.unsplash.com/photo-1515879218367-8466d910aaa4?auto=format&fit=crop&w=1200&q=80)
:::

```markdown
::: wa-carousel {navigation="true" pagination="true"}
![First screenshot](/images/one.webp)
First screenshot caption
![Second screenshot](/images/two.webp)
Second screenshot caption
:::
```

### Animated Image

Use animated images to give readers play and pause controls for GIF or animated WebP content.

::: wa-animated-image
![Animated demo](https://media.giphy.com/media/26tn33aiTi1jkl6H6/giphy.gif)
:::

```markdown
::: wa-animated-image
![Animated demo](/images/demo.gif)
:::
```

The longer `::: webawesome comparison` form remains supported as an alias, but `wa-*` is the recommended authoring style. You can combine Web Awesome container classes with layout or utility classes; for example, a rendered container with `class="not-prose wa-carousel"` is still detected and converted.

```markdown
::: webawesome comparison {position=65}
![Before](/images/before.webp)
![After](/images/after.webp)
:::
```

## Configuration

Enable the plugin in your hooks list. The default `source` is `"vendor"`, which downloads and self-hosts Web Awesome through the shared `[markata-go.assets]` pipeline.

```toml
[markata-go]
hooks = ["default", "webawesome"]

[markata-go.webawesome]
enabled = true
source = "vendor"
version = "3.5.0"
output_dir = "assets/vendor/webawesome"
theme = "default"
palette = "default"
brand = "blue"
```

When `source = "vendor"`, Web Awesome is downloaded into the shared assets cache (`[markata-go.assets].cache_dir`) and published under the shared vendor root (`[markata-go.assets].output_dir`). The plugin's `output_dir` remains `assets/vendor/webawesome` by default so pages load `/assets/vendor/webawesome/...` URLs.

If you prefer to load assets from a CDN instead of self-hosting, set `source = "cdn"`:

```toml
[markata-go.webawesome]
source = "cdn"
version = "3.5.0"
# cdn_base_url defaults to https://cdn.jsdelivr.net/npm/@awesome.me/webawesome@<version>/dist-cdn
```

Self-hosting (the default) is recommended when your site must build without external CDN dependencies, when your privacy policy avoids third-party asset requests, or when you want stable, immutable component versions tied to your repo. You can prefetch the shared asset cache with `markata-go assets download`.

## Raw Components

Raw Web Awesome HTML works too. Any page containing a `<wa-*>` element automatically gets the configured Web Awesome assets.

```html
<wa-callout variant="brand">Raw Web Awesome also works.</wa-callout>
```

<wa-callout variant="brand">Raw Web Awesome also works.</wa-callout>

## Component Loading

Pages that use any `wa-*` element automatically get the Web Awesome stylesheet and initialized Web Awesome autoloader. The built-in template imports `webawesome.loader.js`, sets the configured base path, and calls `startLoader()` so the autoloader watches the document for `<wa-*>` tags and lazy-imports the matching component modules on demand. Each page only downloads what it actually uses. Pages without `wa-*` elements do not load Web Awesome CSS or JavaScript.

When `source = "vendor"`, the theme stylesheet and autoloader are served from your site under `output_dir` (for example `/assets/vendor/webawesome/styles/themes/default.css` and `/assets/vendor/webawesome/webawesome.loader.js`). The focused theme stylesheet avoids Web Awesome's global native-element styles so normal site typography and link styling stay unchanged. These files come from the npm tarball downloaded through the shared `[markata-go.assets]` cache/output pipeline. When `source = "cdn"`, they are served from the browser-ready `https://cdn.jsdelivr.net/npm/@awesome.me/webawesome@VERSION/dist-cdn/` subtree.
