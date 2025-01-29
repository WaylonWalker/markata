"""
DidYouMean Plugin for Markata

Automatically generates redirect pages for URLs that may be mistyped by users.

# Installation

```toml
hooks = [
    "markata.plugins.didyoumean",
]
```

# Configuration

```toml
[markata.didyoumean]
output_dir = "markout"  # Directory where HTML files will be saved
didyoumean_filter = "True"  # A filter expression to determine which pages should be included in suggestions
search_hotkey = "/"  # Hotkey to focus the search input. Set to None to disable. Default is "/"
```

# Usage

This plugin will generate HTML redirect pages for missing URLs that forward users
to the most relevant existing page, or present a list of suggested pages when ambiguity exists.
"""

from markata.hookspec import hook_impl
from pathlib import Path
import pydantic
from pydantic import Field
from typing import TYPE_CHECKING, Optional

if TYPE_CHECKING:
    from markata import Markata


def render_template(markata: "Markata", template_name: str, **context) -> str:
    """Render a template with the given context."""
    template = markata.config.jinja_env.get_template(template_name)
    return template.render(markata=markata, body="", config=markata.config, **context)


class Config(pydantic.BaseModel):
    didyoumean_filter: str = Field(
        default="True",
        description="A string to filter slugs by. If set, only slugs containing this string will have a 404 page generated.",
    )
    search_hotkey: Optional[str] = Field(
        default="/",
        description="Hotkey to focus the search input. Set to None to disable.",
    )


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
def save(markata: "Markata") -> None:
    output_dir = Path(markata.config.output_dir)
    output_dir.mkdir(exist_ok=True, parents=True)

    # Get all existing pages
    existing_pages = set(
        p.strip("/") for p in markata.map("slug", filter="slug is not None") if p
    )

    # Filter slugs based on configuration
    didyoumean_filter = getattr(markata.config, "didyoumean_filter", "")
    filtered_slugs = (
        [slug for slug in existing_pages if didyoumean_filter in slug]
        if didyoumean_filter
        else list(existing_pages)
    )

    # Generate 404 page with all slugs
    html = render_template(markata, "404.html", slugs=filtered_slugs)

    # Write 404 page
    not_found_file = output_dir / "404.html"
    not_found_file.write_text(html, encoding="utf-8")

    # Maps for redirects and suggestions
    redirect_map = {}  # path -> target
    suggestions_map = {}  # path -> set of possible targets

    # Process each page
    for full_slug in existing_pages:
        if not full_slug:
            continue

        # Clean the slug and split into parts
        parts = full_slug.strip("/").split("/")

        # Generate all possible paths
        possible_paths = []

        # Build all possible combinations of parts
        for i in range(len(parts)):  # Start with each part
            current = [parts[i]]
            possible_paths.append("/".join(current))  # Add single part

            # Add combinations with this part and later parts
            for j in range(i + 1, len(parts)):
                current.append(parts[j])
                path = "/".join(current)
                if path != full_slug:  # Don't include the full path
                    possible_paths.append(path)

            # Also try combinations that skip parts
            if i > 0:  # Only if we're not at the first part
                for j in range(i + 1, len(parts)):
                    path = f"{parts[0]}/{parts[j]}"  # First part + current part
                    if path != full_slug:
                        possible_paths.append(path)

        # Check each possible path
        for path in possible_paths:
            if path not in existing_pages:
                if path in redirect_map:
                    if redirect_map[path] != full_slug:
                        suggestions_map.setdefault(path, set()).add(redirect_map[path])
                        suggestions_map[path].add(full_slug)
                        redirect_map.pop(path)
                elif path in suggestions_map:
                    suggestions_map[path].add(full_slug)
                else:
                    redirect_map[path] = full_slug

    # Generate redirect pages
    for redirect_path, target in redirect_map.items():
        redirect_file = output_dir / redirect_path / "index.html"
        if redirect_file.exists():
            continue

        target_url = f"/{target}"
        redirect_file.parent.mkdir(parents=True, exist_ok=True)
        html = render_template(markata, "redirect.html", target_url=target_url)
        redirect_file.write_text(html, encoding="utf-8")

    # Generate suggestion pages
    for suggestion_path, suggestions in suggestions_map.items():
        suggestion_file = output_dir / suggestion_path / "index.html"
        if suggestion_file.exists():
            continue

        full_suggestions = {f"/{s}" for s in suggestions}
        suggestion_file.parent.mkdir(parents=True, exist_ok=True)
        html = render_template(
            markata,
            "suggestions.html",
            path=suggestion_path,
            suggestions=sorted(full_suggestions),
        )
        suggestion_file.write_text(html, encoding="utf-8")
