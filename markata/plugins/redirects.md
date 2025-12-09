---
date: 2025-12-09
description: "The plugin creates static redirects for your site using a simple configuration
  file. Compatible with services like Cloudflare Pages and Netlify. Installation\u2026"
published: false
slug: markata/plugins/redirects
title: redirects.py


---

---

The `markata.plugins.redirects` plugin creates static redirects for your site using
a simple configuration file. Compatible with services like Cloudflare Pages and Netlify.

## Installation

This plugin is built-in but not enabled by default. Add it to your plugins list:

```toml
hooks = [
    "markata.plugins.redirects",
]
```

## Uninstallation

Remove the plugin from your hooks list in `markata.toml`:

```toml
hooks = [
    # Remove or comment out the line below
    # "markata.plugins.redirects",
]
```

## Configuration

Configure redirects in `markata.toml`:

```toml
[markata]
# Default redirects file location
redirects = "static/_redirects"

# Or use assets_dir to set default location
assets_dir = "static"
```

## Redirects File Format

Create a `_redirects` file with entries:

```text
# Basic redirect
/old-path    /new-path

# Force specific status code
/api/*    /v2/api/:splat    301

# Redirect with query parameters
/search    /new-search    301?q=:q

# Proxy to external URL
/external    https://api.example.com    200

# Redirect based on country
/app/*    /app/us/:splat    200    Country=us
/app/*    /app/mx/:splat    200    Country=mx

# Redirect with placeholders
/blog/:year/:month    /posts/:year/:month
```

## Functionality

## Redirect Types

Supports:
- Basic redirects
- Path patterns
- Query parameters
- Status codes
- Country rules
- Placeholders

## File Generation

The plugin:
1. Reads redirect rules
2. Creates HTML files
3. Handles directories
4. Preserves parameters

## Compatibility

Works with:
- Cloudflare Pages
- Netlify
- Static hosting
- Local development

## Performance

Features:
- Efficient file creation
- Pattern matching
- Rule validation
- Error handling

## Dependencies

This plugin depends on:
- pydantic for configuration
- pathlib for file operations

---

!!! class
    <h2 id="Redirect" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">Redirect <em class="small">class</em></h2>

    DataClass to store the original and new url

???+ source "Redirect <em class='small'>source</em>"
    ```python
    class Redirect(pydantic.BaseModel):
        "DataClass to store the original and new url"

        original: str
        new: str
        markata: Markata
        model_config = pydantic.ConfigDict(
            validate_assignment=True, arbitrary_types_allowed=True
        )
    ```
!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

    saves an index.html in the directory called out by the redirect.

???+ source "save <em class='small'>source</em>"
    ```python
    def save(markata: "Markata") -> None:
        """
        saves an index.html in the directory called out by the redirect.
        """
        redirects_file = Path(markata.config.redirects.redirects_file)
        if redirects_file.exists():
            raw_redirects = redirects_file.read_text().split("\n")
        else:
            raw_redirects = []

        redirects = [
            Redirect(original=s[0], new=s[1], markata=markata)
            for r in raw_redirects
            if "*" not in r and len(s := r.split()) == 2 and not r.strip().startswith("#")
        ]

        if "redirect_template" in markata.config:
            template_file = Path(str(markata.config.get("redirect_template")))
        else:
            template_file = DEFAULT_REDIRECT_TEMPLATE

        # Get template mtime to bust cache when template changes
        template_mtime = template_file.stat().st_mtime if template_file.exists() else 0

        key = markata.make_hash("redirects", "raw_redirects", raw_redirects, str(template_mtime))
        with markata.cache as cache:
            cache.get(key)
            if cache.get(key) == "done":
                return

            cache.set(key, "done", expire=markata.config.default_cache_expire)

        template = Template(template_file.read_text())

        for redirect in redirects:
            file = markata.config.output_dir / redirect.original.strip("/") / "index.html"
            file.parent.mkdir(parents=True, exist_ok=True)
            new_content = template.render(redirect.dict(), config=markata.config)
            current_content = file.read_text() if file.exists() else ""
            if current_content != new_content:
                file.write_text(new_content)
    ```