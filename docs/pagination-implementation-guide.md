---
title: Pagination Implementation Guide
description: Guide for implementing pagination in Markata templates with manual, HTMX, and JavaScript options

---

# Pagination Implementation Guide

This guide provides clear instructions for implementing pagination in your Markata templates. The pagination system is integrated into the feeds plugin and supports three types: manual, HTMX, and JavaScript infinite scroll.

## Core Components

### Feeds Plugin with Pagination
**Location:** `markata/plugins/feeds.py`

The feeds plugin includes built-in pagination support with these features:
- Automatic pagination for any feed
- Three pagination types: manual, HTMX, JavaScript
- Configurable items per page
- SEO-friendly URL generation
- Template context variables

### Pagination Plugin (Core Logic)
**Location:** `markata/plugins/pagination.py`

Core pagination functionality that provides:
- Pagination configuration models
- URL generation helpers
- Template rendering context
- Support for all pagination types

## Quick Setup

### 1. Basic Configuration

Add to your `markata.yaml`:

```yaml
# Enable pagination globally (optional, per-feed config also available)
pagination:
  default:
    enabled: true
    items_per_page: 20
    pagination_type: 'manual'  # 'manual', 'htmx', 'js'

# Or configure per feed
feeds:
  blog:
    template: "feed.html"
    partial_template: "feed_partial.html"
    enabled: true
    items_per_page: 10
    pagination_type: 'htmx'
```

### 2. Template Variables

All pagination templates receive these variables:

```jinja2
{{ markata }}           # Markata instance
{{ feed }}              # Current feed object
{{ posts }}             # Posts for current page
{{ page }}              # Current page number (1-based)
{{ total_pages }}       # Total number of pages
{{ has_prev }}          # Boolean: has previous page?
{{ has_next }}          # Boolean: has next page?
{{ prev_page }}         # Previous page number or null
{{ next_page }}         # Next page number or null
{{ pagination_enabled }} # Boolean: is pagination enabled?
{{ pagination_config }} # PaginationConfig object
{{ feed_name }}         # Feed slug/name
{{ pagination_js_url }}     # URL to pagination config JS (when using JS pagination)
```

## Pagination Types

### 1. Manual Pagination

**Best for:** SEO, accessibility, traditional blogs
**Features:**
- Traditional click navigation
- Page numbers
- Previous/Next buttons
- Clean permanent URLs
- Works without JavaScript

**Configuration:**
```yaml
feeds:
  blog:
    pagination_type: 'manual'
    items_per_page: 10
    show_page_numbers: true
    max_page_links: 7
```

**Template Implementation:**

The feeds plugin automatically handles manual pagination when `pagination_type: 'manual'`. Your template just needs to include pagination controls:

```jinja2
{% if pagination_enabled %}
{% set config = pagination_config %}

<div class="pagination-controls mt-8 flex justify-center space-x-4">
    {% if prev_page %}
    {% if page > 2 %}
    <a href="/{{ feed_name }}/{{ prev_page }}/" 
    {% else %}
    <a href="/{{ feed_name }}/" 
    {% endif %}
       class="px-4 py-2 bg-neutral-700 text-white rounded hover:bg-neutral-600 transition">
        ← Previous
    </a>
    {% endif %}
    
    <span class="px-4 py-2 text-neutral-400">
        {{ page }} / {{ total_pages }}
    </span>
    
    {% if has_next %}
    <a href="/{{ feed_name }}/{{ next_page }}/" 
       class="px-4 py-2 bg-neutral-700 text-white rounded hover:bg-neutral-600 transition">
        Next →
    </a>
    {% endif %}
</div>
{% endif %}
```

### 2. HTMX Pagination

**Best for:** Progressive enhancement, modern UX with fallback
**Features:**
- Infinite scroll with 14KB HTMX library
- Progressive enhancement built-in
- SEO-friendly URLs
- Graceful JavaScript fallback

**Configuration:**
```yaml
feeds:
  blog:
    pagination_type: 'htmx'
    items_per_page: 15
    show_loading_skeleton: true
    auto_load_threshold: 200
```

**Template Implementation:**

```jinja2
{% if pagination_enabled and pagination_context.pagination_type == 'htmx' %}
{% if has_next %}
<div id="load-more-trigger" 
     hx-get="/{{ feed_name }}/{{ next_page }}/partial/"
     hx-target="#feed"
     hx-swap="beforeend"
     hx-trigger="revealed"
     hx-indicator=".loading-indicator"
     hx-push-url="/{{ feed_name }}/{{ next_page }}/">
</div>
{% endif %}

<div class="loading-indicator htmx-indicator" style="display: none;">
  <div class="spinner"></div>
  <span>Loading more...</span>
</div>

<script src="/static/js/htmx.org@1.9.10.min.js"></script>
{% endif %}
```

### 3. JavaScript Pagination

**Best for:** Custom infinite scroll, zero external dependencies
**Features:**
- Custom infinite scroll using Intersection Observer
- Zero external dependencies
- AJAX content loading
- URL history management
- Loading indicators

**Configuration:**
```yaml
feeds:
  blog:
    pagination_type: 'js'
    items_per_page: 12
    show_loading_skeleton: true
    auto_load_threshold: 300
    show_end_message: true
```

**Template Implementation:**

```jinja2
{% if pagination_enabled and pagination_context.pagination_type == 'js' %}
<script src="{{ pagination_js_url|default('/static/js/pagination-config.js') }}"></script>
<script src="/static/js/pagination.js"></script>

<div id="scroll-trigger" style="height: 1px;"></div>

<div class="loading-indicator" style="display: none;">
  <div class="spinner"></div>
  <span>Loading more...</span>
</div>
{% endif %}
```

## URL Structure

The pagination system generates clean, SEO-friendly URLs:

- **First page:** `/feed-name/`
- **Subsequent pages:** `/feed-name/page/2/`, `/feed-name/page/3/`, etc.
- **Partial files:** `/feed-name/partial/`, `/feed-name/page/2/partial/`

## Static Assets

### Required Files

The pagination system includes these static files (automatically created):

```bash
markata/static/js/
├── pagination.js              # JavaScript infinite scroll module
├── pagination-config.js       # Generated config (JS pagination only)
└── htmx.org@1.9.10.min.js  # Secure HTMX download

markata/static/css/
└── pagination.css              # Pagination styles and animations
```

## Advanced Configuration

### Per-Feed Customization

```yaml
feeds:
  blog:
    enabled: true
    pagination_type: 'manual'
    items_per_page: 8
    show_page_numbers: true
    max_page_links: 5
    show_loading_skeleton: false
    auto_load_threshold: 100
    show_end_message: false
    
  news:
    enabled: true
    pagination_type: 'htmx'
    items_per_page: 20
    show_loading_skeleton: true
    auto_load_threshold: 200
    show_end_message: true
    
  portfolio:
    enabled: true
    pagination_type: 'js'
    items_per_page: 12
    show_loading_skeleton: true
    auto_load_threshold: 300
    show_end_message: true
```

### Template Includes

Create reusable template components:

**`includes/pagination_info.html`:**
```jinja2
{% if pagination_enabled %}
<div class="pagination-info text-gray-500 text-sm mb-4">
    <span>Page {{ page }} of {{ total_pages }}</span>
    <span class="ml-4">Showing {{ posts|length }} items</span>
</div>
{% endif %}
```

**`includes/post_card.html`:**
```jinja2
<article class="post mb-6">
    <h2><a href="/{{ post.slug }}/">{{ post.title }}</a></h2>
    {% if post.date %}
    <time datetime="{{ post.date }}">{{ post.date.strftime('%B %d, %Y') }}</time>
    {% endif %}
    {% if post.description %}
    <p>{{ post.description }}</p>
    {% endif %}
</article>
```

Then in your main template:
```jinja2
{% include "includes/pagination_info.html" %}

<div class="posts">
    {% for post in posts %}
    {% include "includes/post_card.html" %}
    {% endfor %}
</div>
```

## CLI Helper

Get example configuration:

```bash
markata pagination config-example
```

This outputs a ready-to-use configuration block for your `markata.yaml`.

## Testing Your Implementation

1. **Manual Pagination:** Verify page numbers and Previous/Next links work
2. **HTMX Pagination:** Test infinite scroll and JavaScript fallback
3. **JavaScript Pagination:** Verify auto-loading and URL updates
4. **SEO:** Check that each page has unique titles and meta tags
5. **Accessibility:** Test keyboard navigation and screen readers

## Troubleshooting

### Common Issues

**Pagination not showing:**
- Ensure `enabled: true` is set for your feed
- Check that pagination plugin loads before feeds plugin

**HTMX not working:**
- Verify HTMX script is loaded
- Check that partial template exists and is accessible

**JavaScript errors:**
- Ensure browser supports Intersection Observer
- Check console for fetch API errors

**URL issues:**
- Verify your web server supports clean URLs
- Check that page 1 redirects work correctly

### Debug Mode

Add this to templates to debug pagination data:

```jinja2
{% if markata.config.debug %}
<pre>{{ pagination_config | pprint }}</pre>
<pre>Page: {{ page }}, Total: {{ total_pages }}</pre>
<pre>Has Prev: {{ has_prev }}, Has Next: {{ has_next }}</pre>
{% endif %}
```

This comprehensive guide should help you implement any pagination type in your Markata templates. Choose the pagination type that best fits your use case and customize templates to match your site's design.