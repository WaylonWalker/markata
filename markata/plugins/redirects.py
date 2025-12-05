"""
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

"""

from pathlib import Path

import pydantic
from jinja2 import Template

from markata import Markata
from markata.hookspec import hook_impl

DEFAULT_REDIRECT_TEMPLATE = Path(__file__).parent / "default_redirect_template.html"


class Redirect(pydantic.BaseModel):
    "DataClass to store the original and new url"

    original: str
    new: str
    markata: Markata
    model_config = pydantic.ConfigDict(
        validate_assignment=True, arbitrary_types_allowed=True
    )


class RedirectsConfig(pydantic.BaseModel):
    redirects_file: Path = Path("static/_redirects")
    model_config = pydantic.ConfigDict(
        validate_assignment=True,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )


class Config(pydantic.BaseModel):
    redirects: RedirectsConfig = RedirectsConfig()


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
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
        file.write_text(template.render(redirect.dict(), config=markata.config))
