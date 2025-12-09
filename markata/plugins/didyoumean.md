---
date: 2025-12-09
description: "DidYouMean Plugin for Markata Automatically generates redirect pages
  for URLs that may be mistyped by users. Installation Configuration Usage This plugin
  will\u2026"
published: false
slug: markata/plugins/didyoumean
title: didyoumean.py


---

---

DidYouMean Plugin for Markata

Automatically generates redirect pages for URLs that may be mistyped by users.

## Installation

```toml
hooks = [
    "markata.plugins.didyoumean",
]
```

## Configuration

```toml
[markata.didyoumean]
output_dir = "markout"  # Directory where HTML files will be saved
didyoumean_filter = "True"  # A filter expression to determine which pages should be included in suggestions
search_hotkey = "/"  # Hotkey to focus the search input. Set to None to disable. Default is "/"
```

# Usage

This plugin will generate HTML redirect pages for missing URLs that forward users
to the most relevant existing page, or present a list of suggested pages when ambiguity exists.

---

!!! function
    <h2 id="render_template" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">render_template <em class="small">function</em></h2>

    Render a template with the given context.

???+ source "render_template <em class='small'>source</em>"
    ```python
    def render_template(markata: "Markata", template_name: str, **context) -> str:
        """Render a template with the given context."""
        template = markata.jinja_env.get_template(template_name)
        return template.render(markata=markata, body="", config=markata.config, **context)
    ```