# Template Creation And Management

Use this topic when the task changes layout, partials, cards, feed pages, or page-specific rendering.

## Template Engine

Markata-go uses `pongo2`, which is Django/Jinja2-like.

For a fuller offline field list, read `../reference/template-context.md`.

Core syntax:

```django
{{ post.title }}
{% if post.published %}...{% endif %}
{% for tag in post.tags %}{{ tag }}{% endfor %}
{% extends "base.html" %}
{% include "partials/header.html" %}
```

## Variable Casing

Template variables are **lowercase/snake_case**. This is critical: pongo2 does exact-match key lookups, so `post.Title` silently produces empty output.

Correct: `{{ post.title }}`, `{{ post.article_html }}`, `{{ config.url }}`, `{{ page.has_next }}`

Wrong: `{{ post.Title }}`, `{{ post.ArticleHTML }}`, `{{ page.HasNext }}`

Known PascalCase compatibility aliases are limited. `post.Extra` and `config.Extra` are intentional for the dynamic extras namespace: `{{ post.Extra.image }}`, `{{ config.Extra.custom_key }}`.

`post.templateKey` may also appear as a compatibility alias for older content and migrations. Prefer `post.template` for new template work, but do not assume `post.templateKey` is invalid when reading or debugging an existing site.

## Search Order

Templates are resolved in this order:

1. project `templates/`
2. `themes/<theme>/templates/` if present
3. `themes/default/templates/` if present
4. embedded default theme templates bundled with markata-go

That means a file in the site's `templates/` directory wins.

The default theme includes `layouts/docs.html` for the docs layout. Use
`layout.name = "docs"` before creating a project override; add
`templates/layouts/docs.html` only when the page needs different markup.

## Default Files To Know

- `templates/base.html`
- `templates/post.html`
- `templates/feed.html`
- `templates/layouts/`
- `templates/components/`
- `templates/partials/`

Some built-in feed variants can be selected directly in config without creating a site template override. Example:

```toml
[[markata-go.feeds]]
slug = "shots"

[markata-go.feeds.templates]
html = "feed-photo-grid.html"
```

Use that first when the request is “change one feed’s archive layout” and the built-in variant already matches the need.

For a gallery page, prefer `template: gallery.html` in post frontmatter with a
`gallery` list. Each item needs `src` and `alt`; `caption`, `width`, and `height`
are optional. The built-in template provides responsive columns and an image
viewer. Keep original image links so the gallery works without JavaScript.

For first sites, assume:

- `post.html` is the single-content template
- `feed.html` is the list/archive template
- `base.html` provides the shared page shell

## High-Value Context Variables

Usually available:

- `post`
- `body`
- `config`
- `feed`
- `page`
- `posts`
- `core`

Convenience aliases are also injected:

- `title`
- `date`
- `tags`
- `slug`
- `href`
- `published`
- `draft`
- `description`
- `article_html`
- `site_title`
- `site_url`
- `site_description`
- `site_author`

Image-library templates receive these extra values when the built-in `images`
plugin renders its page:

- `image_library` — page counts and presentation-ready image cards
- `image_index` — the canonical image inventory
- `image_library_config` — the resolved `[markata-go.images]` settings

Image-library cards include `IsVideo`, `PosterSrc`, `Embed`, `AddedAt`,
`AddedAtUnix`, `LastUsedAt`, and `LastUsedAtUnix` alongside the regular image
fields. `AddedAt` is the earliest valid public post date for the source, not a
filesystem timestamp. Each item in `Uses` can include `Caption` when the media
has an associated figure caption. Render video cards with a `<video>` element
and a nested `<source>` element; use `poster_src` when it is non-empty.
The bundled image-library template keeps available posters visible but defers
video sources until a card enters the viewport. Keep that lazy-loading behavior
in custom templates when the page can contain many videos.

Keep `image_index` sources canonical in copy actions. Trusted Dropper derivatives
are for previews only. Override `images.html` in the site's `templates/`
directory when the default authoring page needs a different layout.

## Common Post Fields In Templates

Typical `post` keys used in templates:

- `post.title`
- `post.slug`
- `post.href`
- `post.date`
- `post.tags`
- `post.description`
- `post.article_html`
- `post.structured_data`
- `post.Extra.*` in some templates

The rendered article body is normally output with:

```django
{{ body | safe }}
```

## Feed And Pagination Variables

Useful feed/page fields:

- `feed.title`
- `feed.description`
- `feed.slug`
- `feed.posts`
- `page.number`
- `page.posts`
- `page.has_prev`
- `page.has_next`
- `page.prev_url`
- `page.next_url`
- `page.total_pages`
- `page.pagination_type`

## Built-In Filters Worth Knowing

- `date`
- `date_format`
- `human_date`
- `rss_date`
- `atom_date`
- `slugify`
- `truncate`
- `truncatewords`
- `default_if_none`

Use `human_date` for visible HTML dates in cards, post bylines, archive views, and reader metadata. Keep `<time datetime="...">` values machine-readable with `atom_date` or RFC 3339 strings.
- `length`
- `first`
- `last`
- `join`
- `sort`
- `striptags`
- `plaintext`
- `summary` (single-line, escaped card excerpt; use this instead of `plaintext|truncatechars` in cards)
- `safe`
- `absolute_url`
- `theme_asset`
- `theme_asset_hashed`
- `asset_url`
- `slides_reveal`

For media cards, remember that `media_url` normalizes trusted media URLs to `https`, while `with_size` and `poster_url` only decorate relative or trusted CDN URLs. When the card intentionally uses a cinematic crop, use width-only sizing such as `with_size:"1200"` so the CDN preserves the source proportions, then use CSS `object-fit: cover` to crop the displayed frame. When the source aspect ratio must be preserved, use width-only sizing and let CSS use `height: auto`.

## Common Patterns

### Base Template

```django
<!DOCTYPE html>
<html>
<head>
  <title>{% block title %}{{ site_title }}{% endblock %}</title>
</head>
<body>
  {% block content %}{% endblock %}
</body>
</html>
```

### Post Template

```django
{% extends "base.html" %}

{% block title %}{{ post.title }} | {{ site_title }}{% endblock %}

{% block content %}
<article>
  <h1>{{ post.title }}</h1>
  {% if post.date %}<time>{{ post.date | date:"January 2, 2006" }}</time>{% endif %}
  <div class="post-content">{{ body | safe }}</div>
</article>
{% endblock %}
```

### Feed Template

```django
{% extends "base.html" %}

{% block content %}
<section>
  <h1>{{ feed.title }}</h1>
  {% for post in page.posts %}
    <article>
      <h2><a href="{{ post.href }}">{{ post.title }}</a></h2>
    </article>
  {% endfor %}
</section>
{% endblock %}
```

## Web Awesome In Templates And Content

For sites using the `webawesome` hook, prefer authoring content with markdown containers such as `::: wa-details`, `::: wa-tabs`, `::: wa-tooltip`, or `::: wa-comparison` before dropping into raw custom-element HTML.

Useful guidance:

- raw `<wa-*>` elements are also supported and trigger asset loading automatically
- do not hardcode Web Awesome `<script>` or `<link>` tags if the site already relies on the built-in plugin wiring; if you must customize the template wiring, import `webawesome.loader.js`, set the configured base path, and call `startLoader()`
- if a template needs to guard extra markup or classes around Web Awesome usage, inspect existing checks first; per-page detection is usually driven by the presence of rendered `wa-*` elements in the page body
- if vendor assets are enabled, the resolved URLs come from config-driven asset mappings rather than fixed `/assets/vendor/webawesome/` assumptions

Typical authoring examples:

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

```markdown
::: wa-tooltip {content="Static Site Generator"}
SSG
:::
```

## Hover Cards In Content

The default theme's `js/tooltips.js` is the shared hover preview component. Wikilinks and glossary terms get cards automatically. For other hover previews in Markdown, use its attributes instead of hand-writing `wa-tooltip` plus inline scripts:

```html
<span data-hover-title="Static Site Generator" data-hover-description="Builds HTML ahead of time.">SSG</span>

<span data-hover-card="#term-card">rich card</span>
<template id="term-card"><strong><a href="/term/">Term</a></strong><p>Details.</p></template>
```

- Cards follow the active palette and open on hover and keyboard focus. On a wrapped link they open on the side the pointer is on (above for the first line, below for the last) and flip when there is no room.
- Custom popups can reuse the placement with `window.MarkataHoverCards.placement(target, width, height)`, which returns `{ top, left, placement }` in viewport pixels.
- Style them with `.hover-card` (plus `.wikilink-tooltip`, `.hover-card--glossary`, `.hover-card--external`, or `.hover-card--rich`).
- External links get cards only when `[markata-go.external_link_hover] enabled = true`. It reuses the embeds metadata cache, blogroll feed entries, and link titles, and makes no network requests unless `fetch = true`. Add `data-no-preview` to an inline `<a>` to skip one link. Don't hand-write `data-link-*` attributes; the plugin generates them.
- Theme scripts can add sources with `window.MarkataHoverCards.register(selector, el => node)`.
- If a template override drops `js/tooltips.js`, no cards appear.

## Per-Post Template Selection

Common frontmatter options:

```yaml
template: custom/special-page.html
layout: landing
```

For presentation decks, `template: slides.html` is available. The bundled default slides template uses reveal.js and splits rendered content with common markdown deck conventions while remaining compatible with the normal H1 content lint rule:

- `##` / rendered `h2`: new horizontal slide
- `###` / rendered `h3`: new vertical slide
- `---` / rendered `hr`: new horizontal slide

When a site self-hosts third-party assets with `[markata-go.assets]`, `slides.html` uses the shared `asset_urls` mappings for Reveal.js automatically.

Use `layout` when the site has a base-driven layout system. Use `template` when you need an explicit file.

## Layout Resolution

When a site uses the layout system, the effective layout is resolved roughly in this order:

1. frontmatter `layout`
2. feed-based layout config
3. global `[markata-go.layout].name`
4. fallback blog-like layout

Useful config shape:

```toml
[markata-go.layout]
name = "blog"

[markata-go.layout.paths]
"/docs/" = "docs"
"/about/" = "landing"

[markata-go.layout.feeds]
"docs" = "docs"
"blog" = "blog"
```

## Template Presets And Default Templates

Some sites may use config-defined template presets or default per-format templates.

Useful concepts:

- `template_presets`: named bundles for html/text/ansi/markdown/og templates
- `default_templates`: global fallback templates per format

If a post has `template: blog` and the site defines a `blog` template preset, that preset can expand to multiple output-format templates.

If the repo has no local examples, start from:

- `../examples/templates/base.html`
- `../examples/templates/post.html`
- `../examples/templates/feed.html`

## When To Choose Which Mechanism

- small structural change across many pages: layout config
- one reusable view pattern: partial or layout
- single page needs a unique shell: `template`
- many formats need coordinated template mapping: template preset
- small markup tweak: edit the smallest partial involved

## Per-Format Template Rules

Markata-go adapts template names by output format:

- `post.html` -> `post.txt`
- `post.html` -> `post.ansi`
- `post.html` -> `post.md`
- `post.html` -> `post-og.html`

Hardcoded fallbacks include:

- HTML: `post.html`
- text: `default.txt`
- ANSI: `default.ansi`
- markdown/text-ish fallback: `raw.txt`
- OG: `og-card.html`

## Guidance

- Prefer overriding the smallest template that solves the task.
- Use `layout` for base-driven page variants and `template` only when a full standalone template is needed.
- Reuse existing partials and components before creating new ones.
- Keep naming aligned with the site's current template structure.
- Confirm which template is already in use through nearby frontmatter and config before editing.
- For a first site, start by copying and editing `post.html`, `feed.html`, or a single partial instead of rebuilding the entire template tree.
- If a task sounds like sidebar placement, TOC behavior, or page shell selection, check `[markata-go.layout]` before creating new template files.

## Embedding Feeds In Content Pages

Two template helper functions are available for embedding feed content inside other pages:

### `feed_posts(slug, [limit])`

Returns a list of post maps for a given feed slug. Useful when you need raw data for custom rendering:

```django
{% with feed_posts("blog", 5) as recent %}
  {% for p in recent %}
    <a href="{{ p.href }}">{{ p.title }}</a>
  {% endfor %}
{% endwith %}
```

### `render_feed(slug, [limit], [variant], [{options}])`

Returns rendered HTML for a feed preview. Accepts a variant name (default `"card"`) and optional options map with `template`, `variant`, and `limit` keys:

```django
{{ render_feed("blog", 3) }}
{{ render_feed("blog", 5, "card") }}
```

The default template is `partials/feed_preview.html`. If that template is missing, a basic HTML fallback is used.

Feed helpers only surface posts that are not skipped. Drafts and skipped posts are omitted from feed pages and syndicated outputs; a published post with an empty body still gets a page, so its card never links to a 404.

`render_feed` output is a single HTML block with no blank lines, so goldmark never re-parses card text as Markdown. Keep card partials to inline markup and use the `summary` filter for excerpts. When a slug does not match a feed the helper renders nothing and logs a `jinja_md` warning; check the feed's `slug`, not its name.

## Text And Alternate Format Templates

Markata-go can output posts in multiple formats beyond HTML. Alternate format templates follow the naming pattern:

- `post.txt` for plain text
- `post.ansi` for terminal-colored output
- `post.md` for markdown pass-through
- `post-og.html` for Open Graph card images

Use the `plaintext` filter to strip HTML tags when producing text output. Hardcoded fallbacks: `default.txt` (text), `default.ansi` (ANSI), `raw.txt` (markdown), `og-card.html` (OG).

Feed sidebar variant links should point at the canonical short URLs for enabled formats (`/slug.md`, `/slug.txt`, `/slug.ansi`) instead of nested `/slug/index.<ext>` paths.

## Common Tasks

- customize post or feed cards
- add or adjust partials
- create a new layout under `templates/layouts/`
- apply a page-specific template via frontmatter
- embed a feed listing inside a content page with `render_feed`
