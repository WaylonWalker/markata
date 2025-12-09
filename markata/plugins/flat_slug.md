---
date: 2025-12-09
description: "The plugin handles URL slug generation for your articles. It sanitizes
  special characters that don't work in browsers and provides a consistent URL structure.\u2026"
published: false
slug: markata/plugins/flat-slug
title: flat_slug.py


---

---

The `markata.plugins.flat_slug` plugin handles URL slug generation for your articles.
It sanitizes special characters that don't work in browsers and provides a consistent
URL structure.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.flat_slug",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.flat_slug",
]
```

Note: Disabling this plugin will prevent automatic URL slug generation and sanitization.

## Configuration

Configure the plugin behavior in your `markata.toml`:

```toml
[markata.flat_slug]
slugify = true  # Set to false to disable automatic slug sanitization
```

## Functionality

## Explicit Slug in Frontmatter

If you explicitly set the slug in the frontmatter of a post, markata will not
overwrite it:

```markdown
---
title: My First Post
slug: /my-post
---

This is my first post it will be at `<markata.config.url>/my-post/`
regardless of filename.
```

## Automatic Slug Based on Filename

By default the flat_slug plugin will use the `stem` of your filename (filename without
extension) unless you explicitly set your slug in frontmatter.

Examples:
* `/pages/my-post.md` becomes `<markata.config.url>/my-post/`
* `/pages/blog/a-blog-post.md` becomes `<markata.config.url>/a-blog-post/`

## Registered Attributes

The plugin registers the following attributes on Post objects:
- `should_slugify`: Boolean indicating if the post's slug should be sanitized

---

!!! function
    <h2 id="pre_render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">pre_render <em class="small">function</em></h2>

    Sets the article slug if one is not already set in the frontmatter.

???+ source "pre_render <em class='small'>source</em>"
    ```python
    def pre_render(markata: "Markata") -> None:
        """
        Sets the article slug if one is not already set in the frontmatter.
        """
        from slugify import slugify

        for article in markata.iter_articles(description="creating slugs"):
            stem = article.get(
                "slug",
                Path(article.get("path", article.get("title", ""))).stem,
            )
            if article.should_slugify:
                article.slug = "/".join([slugify(s) for s in stem.split("/")])
            else:
                article.slug = stem
    ```