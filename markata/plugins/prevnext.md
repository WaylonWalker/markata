---
date: 2025-12-09
description: "The plugin adds previous and next navigation links to each post, allowing
  readers to easily navigate between related content. Installation This plugin is\u2026"
published: false
slug: markata/plugins/prevnext
title: prevnext.py


---

---

The `markata.plugins.prevnext` plugin adds previous and next navigation links to each
post, allowing readers to easily navigate between related content.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.prevnext",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.prevnext",
]
```

## Configuration

Configure navigation behavior in `markata.toml`:

```toml
[markata.prevnext]
# Strategy for finding prev/next posts
# 'first': Use first map where post is found
# 'all': Use all maps
strategy = 'first'

# Custom colors (optional)
prevnext_color_text = "white"
prevnext_color_text_light = "black"
prevnext_color_angle = "white"
prevnext_color_angle_light = "black"

# Navigation maps
[[markata.prevnext.maps]]
# Map posts by category
category = "tutorials"
filter = "post.category == 'tutorials'"
sort = "post.date"
reverse = true

[[markata.prevnext.maps]]
# Map posts by series
category = "python-series"
filter = "post.series == 'python'"
sort = "post.part"
reverse = false
```

## Functionality

## Navigation Maps

The plugin allows you to:
1. Define multiple navigation maps
2. Filter posts by attributes
3. Sort posts by any field
4. Control sort direction
5. Group related content

## Navigation Strategies

Two navigation modes:
- `first`: Use first map containing the post
- `all`: Use all maps containing the post

## Template Integration

Adds to each post:
- Previous post link
- Next post link
- Navigation styling
- Responsive design

## Styling

Customizable colors:
- Text colors
- Arrow colors
- Light/dark mode support
- Hover effects

## Dependencies

This plugin depends on:
- jinja2 for templating
- pydantic for configuration

---

!!! class
    <h2 id="UnsupportedPrevNextStrategy" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">UnsupportedPrevNextStrategy <em class="small">class</em></h2>

    A custom error class to raise when an unsupporte prevnext strategy is
    defined.

???+ source "UnsupportedPrevNextStrategy <em class='small'>source</em>"
    ```python
    class UnsupportedPrevNextStrategy(NotImplementedError):
        """
        A custom error class to raise when an unsupporte prevnext strategy is
        defined.
        """
    ```