# Analytics And Storytelling

Use this topic when the user wants an `analytics.md` page, a metrics dashboard, chart-driven posts, contribution graphs, or template badges for reading time and content stats.

## What Agents Should Do First

1. Inspect the active config with `markata-go config show`.
2. Confirm whether the site already enables `stats`, `chartjs`, and `contribution_graph`.
3. Inspect existing analytics pages, feed pages, and template partials before adding a new pattern.
4. Decide whether the task belongs in:
   - a markdown page with `jinja: true`
   - an existing HTML template or partial
   - config only
   - a plugin change

Prefer content, templates, and built-in stats over custom plugin work.

## Reading-Time Configuration

The `reading_time` plugin reads `words_per_minute` from the top-level
`[markata-go]` section and excludes fenced and inline code. The `stats` plugin
keeps its independent metrics in `[markata-go.stats]`; it supports
`words_per_minute`, `include_code_in_count`, and `track_code_blocks`. The
`include_code_in_count` setting controls Stats' word and reading-time metrics.
The `track_code_blocks` setting controls code block and code line metrics; when
false, those per-post and aggregate metrics are zero. The top-level
`words_per_minute` value takes precedence for Stats when both values exist.

## Built-In Data Surfaces

This topic should distinguish between three categories:

- built-in today: fields and blocks users can use directly now
- derivable today: analytics agents can compute from posts, frontmatter, and link data with Jinja/templates
- external-only: traffic and reader behavior that the build cannot know on its own

### Per-post stats

Available on posts after the `stats` plugin runs:

- `post.Extra.word_count`
- `post.Extra.char_count`
- `post.Extra.reading_time`
- `post.Extra.reading_time_text`
- `post.Extra.code_lines`
- `post.Extra.code_blocks`
- `post.Extra.stats`

`post.Extra.stats` is Stats' authoritative `PostStats` value and supplies
feed/site aggregates. In the default plugin registry, `reading_time` runs
later and owns the public `word_count`, `reading_time`, and
`reading_time_text` fields. It uses the top-level reading-time configuration
and excludes code. When `stats` runs alone, it also populates those public
fields from its own configuration.

These fields are good for post headers, feed cards, article footers, and
"long read" badges.

### Site-wide stats

Available in templates and Jinja-enabled markdown via `config.Extra.site_stats`:

- `total_posts`
- `total_words`
- `total_chars`
- `total_reading_time`
- `total_reading_time_text`
- `average_words`
- `average_reading_time`
- `average_reading_time_text`
- `total_code_lines`
- `total_code_blocks`
- `posts_by_year`
- `words_by_year`
- `posts_by_tag`

Use these for analytics pages, about pages, and site summary components.

### Advanced helper

`config.Extra.stats` exposes a helper object for site and feed KPIs. Prefer `config.Extra.site_stats` for straightforward template access. Reach for the helper only when the site already uses it or the task specifically needs helper-driven feed metrics.

Common helper-driven feed access patterns:

- `config.Extra.stats.ForFeed("blog").PostCount()`
- `config.Extra.stats.ForFeed("blog").TotalWords()`
- `config.Extra.stats.ForFeed("blog").TotalReadingTimeText()`
- `config.Extra.stats.ForFeed("blog").PostsByYear()`
- `config.Extra.stats.ForFeed("blog").PostsByTag()`

### Link graph data

Available on each post today:

- `post.hrefs`: raw href values found in the post
- `post.inlinks`: template-friendly links from other posts pointing to this post
- `post.outlinks`: template-friendly links from this post to other pages

Use these to derive:

- orphan notes or posts with no inlinks
- hub notes with many inlinks
- heavily connected posts by total links
- posts with many outbound references
- internal link lists and backlink sections

### Common built-ins available today

The easiest analytics pages should start with these built-ins:

- word count
- reading time
- code block metrics
- posts by year
- words by year
- posts by tag
- total site posts
- total site reading time
- feed-level totals through `config.Extra.stats.ForFeed(...)`
- contribution graphs from post dates
- `chartjs` blocks fed by built-in metrics

### Derivable today with Jinja/templates

These are not first-class dashboard features, but agents can build them with existing data:

- posts per month
- author contribution mix when author fields are present
- draft/private/published counts
- posts by template or layout when frontmatter is consistent
- stale-content lists using `modified` or `lastmod`
- inlink/outlink leaderboards
- orphan notes or orphan posts
- simple second-brain graph summaries

### External-only analytics

These require outside analytics tools or logs and should not be presented as build-derived facts:

- pageviews
- unique visitors
- referrers
- conversions
- search queries
- subscriber counts unless imported from elsewhere

## Chart And Graph Options

Read `topics/chart-authoring.md` for choosing a bar, line, filled-line
area, stacked bar, doughnut, scatter or other chart from the information.
A chart is the right primary visual for **measured comparisons or trends**,
not simply a decorative replacement for a table. Use real `chartjs`
fenced JSON with labels, datasets, and semantic units; use
`contribution-graph` for calendar-style publishing/activity views.

A **line** compares ordered periods, a **bar** compares categories,
a **doughnut** summarizes a known whole with few groups, and a **stacked
bar** shows comparable subgroup contributions. Use the published
Chart.js data-table fallback with aligned label/data arrays.

**Do not put total articles, total words, and total code blocks in one
ordinary bar chart**: the units and magnitudes are incompatible, so
the visual comparison is misleading. Use individual `md-stat` values
or separate same-unit charts. Reuse any site chart palette or CSS treatment;
never claim that the canvas has been verified in every theme without a
browser check.

## Recommended `analytics.md` Pattern

Use a normal markdown page with frontmatter like:

```yaml
---
title: "Site Analytics"
description: "Publishing metrics, activity charts, and yearly trends for this site."
date: 2026-04-09
published: true
slug: /analytics/
jinja: true
tags:
  - analytics
  - meta
  - writing
---
```

Then structure the page in three layers:

1. KPI summary
2. one or two charts
3. author-written interpretation

Agents should compute or scaffold the first two layers and leave space for the third.

## Good Analytics Page Structure

Use sections like:

- `## At a glance`
- `## Publishing activity`
- `## Topic mix`
- `## Link graph`
- `## What changed this year`
- `## Notes from the author`

The first sections can be data-driven. The last sections should invite the author to explain timeline changes, experiments, and tradeoffs in their own voice.

## Safe Authoring Pattern

When asked to create an analytics post, prefer this workflow:

1. Gather site metrics from built-in stats, list commands, and nearby content.
2. Draft an `analytics.md` page with working charts and KPI callouts.
3. Add short prompts or placeholders where the human should add interpretation.
4. Avoid inventing personal history, timelines, or motivations unless the user explicitly asks for prose generation.

Example placeholder style:

```markdown
## What changed this year

Write about the periods where output increased or dipped.
Possible prompts:
- Which projects or life events changed your publishing pace?
- Which topics became more important?
- Were there intentional breaks, experiments, or format changes?
```

## Template Snippets

### Post badge

```django
<div class="post-meta">
  <span>{{ post.Extra.reading_time_text }}</span>
  <span>{{ post.Extra.word_count }} words</span>
</div>
```

### Site KPI block

```django
<section class="site-kpis">
  <p>{{ config.Extra.site_stats.total_posts }} posts</p>
  <p>{{ config.Extra.site_stats.total_words }} words</p>
  <p>{{ config.Extra.site_stats.total_reading_time_text }} total reading time</p>
</section>
```

## Markdown Snippets

### KPI summary in `analytics.md`

```markdown
## At a glance

- {{ config.Extra.site_stats.total_posts }} published posts
- {{ config.Extra.site_stats.total_words }} total words
- {{ config.Extra.site_stats.total_reading_time_text }} total reading time
- {{ config.Extra.site_stats.average_words }} average words per post
```

### Injecting variables into markdown

Use `jinja: true` in frontmatter, then render values directly inside normal markdown or inside JSON fenced blocks.

Guidance:

- keep the surrounding JSON valid after Jinja renders
- prefer simple loops and explicit commas
- reuse existing repo patterns before inventing more dynamic logic
- if a page only needs site totals, prefer `config.Extra.site_stats.*`
- in HTML templates, prefer lowercase map keys like `post.inlinks`; in Jinja markdown loops over `core.filter(...)`, many sites use object-style fields like `post.Inlinks` and `post.Title`

### Contribution graph from post dates

````markdown
## Publishing activity

```contribution-graph
{
  "data": [
    {% for post in core.filter("published == true") %}
    {"date": "{{ post.Date.Format \"2006-01-02\" }}", "value": 1}{% if not loop.last %},{% endif %}
    {% endfor %}
  ],
  "options": {
    "domain": "year",
    "subDomain": "day",
    "maxPercentile": 95
  }
}
```
````

Use this when the author wants a calendar-like view of publishing cadence.
`maxPercentile` is useful when a few unusually large days would otherwise wash out the rest of the heatmap.
The built-in contribution graph runtime also fits graphs to narrow content columns, so agents should not add page-specific resize JavaScript unless the site has a special layout need.

### `chartjs` chart from comparable built-in site stats

This one series of same-unit article counts can be charted meaningfully.
When the site has `posts_by_year`, derive chronologically ordered
labels and values using the repository's existing safe Jinja pattern.
Do not silently invent counts or render raw map values as if sorted.

Illustrative static example (replace with verified annual counts):

````markdown
## Articles by year

```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["2023", "2024", "2025"],
    "datasets": [
      {"label": "Articles", "data": [7, 12, 9]}
    ]
  },
  "options": {
    "scales": {"y": {"beginAtZero": true}}
  }
}
```
````

For site-wide totals measured in *different units*, show
`total_posts`, `total_words`, and `total_code_blocks` as
separately labeled `md-stat` KPIs instead of a single shared bar axis.

### Topic mix chart

Agents may either:

- generate a static `chartjs` block from inspected repo data
- or wire the page to existing site stats if the repo already has a stable pattern for dynamic Jinja-generated chart data

Keep the chart JSON valid and prefer a small number of clearly labeled datasets.

### Feed KPI snippet

````markdown
## Blog feed totals

- {{ config.Extra.stats.ForFeed("blog").PostCount() }} posts in the blog feed
- {{ config.Extra.stats.ForFeed("blog").TotalWords() }} words in the blog feed
- {{ config.Extra.stats.ForFeed("blog").TotalReadingTimeText() }} total reading time
````

Use this when the site has meaningful feed partitions such as `blog`, `docs`, or `notes`.

### Orphan notes section

An orphan note is usually a post with no `inlinks`. Agents should confirm the site's conventions before presenting this as a strong quality signal, because some sites intentionally publish standalone pages.

````markdown
## Orphan notes

{% for post in core.filter("published == true") %}
{% if post.Inlinks|length == 0 %}
- [{{ post.Title }}]({{ post.Href }})
{% endif %}
{% endfor %}
````

If the site already uses a different field casing pattern in markdown/Jinja, mirror the local convention instead of forcing a rewrite. The important thing is the pattern: filter posts, check the inlink count, and render a list.

### Most-linked posts

````markdown
## Most linked notes

{% for post in core.filter("published == true") %}
- [{{ post.Title }}]({{ post.Href }}) has {{ post.Inlinks|length }} inbound links and {{ post.Outlinks|length }} outbound links
{% endfor %}
````

Agents should usually sort or pre-filter this in the site's preferred way if the repo already has helper patterns for ordered post collections.

## Link Analytics Ideas

Good build-side link analytics include:

- orphan notes
- most-linked notes
- posts with the most outbound references
- posts with no outbound links
- dense hub pages for docs or second brains
- backlink sections for related-post discovery
- broken-wikilink follow-up lists when the site treats warnings as actionable

For second-brain sites, these are often more valuable than pageview analytics because they describe knowledge structure rather than reader traffic.

## Storytelling Guidance

Agents should help authors create stories from their data by:

- identifying notable changes in publishing volume, topic mix, or reading-time distribution
- identifying useful link graph patterns such as isolated notes, strong hubs, and unusually reference-heavy posts
- suggesting section headings and narrative prompts grounded in the metrics
- calling out where the data is strong and where it is only directional
- keeping claims tied to visible charts or explicit metrics

Agents should not overstate causation from site metrics alone.
