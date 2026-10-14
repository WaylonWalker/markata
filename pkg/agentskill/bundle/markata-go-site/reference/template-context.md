# Template Context Reference

This reference is for agents working in a markata-go site repo without the markata-go source tree.

## Important: Variable Casing

All template variables use **lowercase/snake_case**. Pongo2 does exact-match key lookups, so PascalCase silently produces empty output.

Correct: `{{ post.title }}`, `{{ post.article_html }}`, `{{ page.has_next }}`

Wrong: `{{ post.Title }}`, `{{ post.ArticleHTML }}`, `{{ page.HasNext }}`

Known PascalCase compatibility aliases are limited. `post.Extra` and `config.Extra` provide access to custom frontmatter and config values.

`post.templateKey` may also be present as a compatibility alias. Prefer `post.template` for new work, but expect `post.templateKey` in older content, migrations, and some template/debugging contexts.

## Core Template Variables

Usually available in HTML templates:

- `post`: current post as a map-like object
- `body`: rendered article HTML string
- `config`: resolved site configuration
- `feed`: current feed config when rendering a feed page
- `page`: current paginated feed page
- `posts`: current post list for feed/index contexts
- `core`: lifecycle manager-like object exposed to templates in some contexts

## Common `post` Keys

- `post.path`
- `post.content`
- `post.slug`
- `post.href`
- `post.hrefs`
- `post.inlinks`
- `post.outlinks`
- `post.published`
- `post.draft`
- `post.private`
- `post.skip`
- `post.tags`
- `post.template`
- `post.templateKey` (compatibility alias for `post.template`)
- `post.html`
- `post.article_html`
- `post.title`
- `post.title_html` (safe rendered inline Markdown for visible headings)
- `post.title_text` (semantic plain text for metadata, feeds, search, and slugs)
- `post.title_text_derived` (whether semantic title derivation completed)
- `post.title_source` (authored/source title)
- `post.date`
- `post.description`
- `post.author`
- `post.authors`
- `post.author_objects`
- `post.Extra`

`post.title` uses the semantic plain title. Use `post.title_html` only in a
trusted visual heading context. A successfully derived empty `title_text` is
intentional and must not fall back to `title_source`. Use
`post.title_text_derived` when selecting a site-title fallback so an empty
derived title stays empty.

## Stats And Analytics Fields

Common per-post stats in `post.Extra`:

- `post.Extra.word_count`
- `post.Extra.char_count`
- `post.Extra.reading_time`
- `post.Extra.reading_time_text`
- `post.Extra.code_lines`
- `post.Extra.code_blocks`
- `post.Extra.stats`

Common site-wide stats in `config.Extra.site_stats`:

- `config.Extra.site_stats.total_posts`
- `config.Extra.site_stats.total_words`
- `config.Extra.site_stats.total_chars`
- `config.Extra.site_stats.total_reading_time`
- `config.Extra.site_stats.total_reading_time_text`
- `config.Extra.site_stats.average_words`
- `config.Extra.site_stats.average_reading_time`
- `config.Extra.site_stats.average_reading_time_text`
- `config.Extra.site_stats.total_code_lines`
- `config.Extra.site_stats.total_code_blocks`
- `config.Extra.site_stats.posts_by_year`
- `config.Extra.site_stats.words_by_year`
- `config.Extra.site_stats.posts_by_tag`

Advanced helper access may also be available via `config.Extra.stats` for site and feed KPIs.

Common feed-helper calls when the helper is exposed:

- `config.Extra.stats.ForFeed("blog").PostCount()`
- `config.Extra.stats.ForFeed("blog").TotalWords()`
- `config.Extra.stats.ForFeed("blog").TotalReadingTimeText()`
- `config.Extra.stats.ForFeed("blog").PostsByYear()`
- `config.Extra.stats.ForFeed("blog").PostsByTag()`

## Link Graph Fields

Common link-analysis fields on `post`:

- `post.hrefs`
- `post.inlinks`
- `post.outlinks`

These are useful for backlink sections, related-note sections, hub-note detection, and orphan-note analysis.

## Common Top-Level Aliases

Markata-go also injects convenient top-level aliases:

- `title`
- `date`
- `tags`
- `slug`
- `href`
- `published`
- `draft`
- `private`
- `description`
- `article_html`

## Common Plugin-Provided Extras

When optional hooks are enabled, templates may also receive extra top-level fields and matching `config.Extra` values.

Web Awesome commonly exposes:

- `needs_webawesome` on posts/pages that rendered `wa-*` elements
- `config.Extra.webawesome_enabled`
- `config.Extra.webawesome_css_url`
- `config.Extra.webawesome_loader_url`
- `config.Extra.webawesome_theme_class`

Use these only after inspecting the active site templates and rendered output. Some sites rely entirely on the built-in base template wiring and never reference these values directly in project templates.

## Site Aliases

- `site_title`
- `site_url`
- `site_description`
- `site_author`
- `authors`
- `default_author`
- `default_author_id`

## Feed Keys

- `feed.slug`
- `feed.base_url`
- `feed.title`
- `feed.description`
- `feed.filter`
- `feed.sort`
- `feed.reverse`
- `feed.items_per_page`
- `feed.limit`
- `feed.offset`
- `feed.posts`
- `feed.formats`

## Page Keys

- `page.number`
- `page.posts`
- `page.has_prev`
- `page.has_next`
- `page.prev_url`
- `page.next_url`
- `page.total_pages`
- `page.total_items`
- `page.items_per_page`
- `page.page_urls`
- `page.pagination_type`

## Reader Page Context

- `day_groups` on `/reader/` templates
- `reader_type` for the active reader stream (`all`, `written`, `video`, `podcast`)
- `reader_types` for switch links between the combined reader and typed reader pages
- each group provides `date_iso`, `weekday`, `title`, `year`, `count`, and `entries`
- reader entry maps include `feed_type`, `source_icon_url`, `published_datetime`, `published_label`, `preview_url`, `preview_kind`, `source_image_url`, and `source_initial`

## Sidebar/Navigation Context

- `sidebar_items`
- `sidebar_title`
- `resolved_content_sidebar`

## Discovery Feed Context

- `discovery_feed.slug`
- `discovery_feed.title`
- `discovery_feed.rss_url`
- `discovery_feed.atom_url`
- `discovery_feed.json_url`
- `discovery_feed.has_rss`
- `discovery_feed.has_atom`
- `discovery_feed.has_json`

Used in `<head>` for `<link rel="alternate">` feed discovery tags. Automatically populated per post based on the post's sidebar feed or the root subscription feed.

## Feed Helper Functions

- `feed_posts(slug, [limit])`: returns a list of post maps for the given feed slug
- `render_feed(slug, [limit], [variant], [{options}])`: returns rendered HTML for a feed preview

See `topics/template-management.md` for usage examples.

## For Loop Variables

Inside `{% for %}` blocks, pongo2 provides `forloop`:

- `forloop.Counter`: 1-indexed iteration number
- `forloop.Counter0`: 0-indexed iteration number
- `forloop.First`: true on first iteration
- `forloop.Last`: true on last iteration
- `forloop.Revcounter`: iterations remaining (1-indexed)

## Common Filters

- `date`
- `date_format`
- `rss_date`
- `atom_date`
- `slugify`
- `truncate`
- `truncatewords`
- `default_if_none`
- `length`
- `first`
- `last`
- `join`
- `sort`
- `reverse`
- `selectattr`
- `rejectattr`
- `striptags`
- `plaintext`
- `linebreaks`
- `linebreaksbr`
- `safe`
- `urlencode`
- `absolute_url`
- `theme_asset`
- `theme_asset_hashed`
- `asset_url`
- `google_fonts_url`
- `reading_time`
- `excerpt`

## Safe Usage Pattern

Use `{{ body | safe }}` for already-rendered HTML body content.

Use plain `{{ post.title }}` or `{{ config.title }}` for text values.

Use `post.Extra.<field>` when the site relies on custom frontmatter not covered by built-in fields.

### Editorial attribution

- `post.byline_authors`: resolved non-editor contributors in original order.
- `post.byline_editors`: resolved editing/reviewing contributors in original order.
- `post.additional_editor_count`: number of editors after the first.
- `post.editorial_history_view`: newest-first normalized schema-v1 events;
  validated GitHub file/diff links and optional before/after raw URLs for lazy
  inline source comparisons. Original history remains in `post.Extra`.
