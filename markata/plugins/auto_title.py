"""
The `markata.plugins.auto_title` plugin automatically generates titles for posts that
don't have a title specified in their frontmatter. It uses the filename to create a
human-readable title.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.auto_title",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.auto_title",
]
```

# Configuration

This plugin requires no configuration. It automatically processes any post without a title.

# Functionality

## Title Generation

The plugin generates titles by:
1. Using the filename (without extension)
2. Replacing hyphens and underscores with spaces
3. Converting to title case

Examples:
- `my-first-post.md` becomes "My First Post"
- `python_tips_and_tricks.md` becomes "Python Tips And Tricks"

## Frontmatter Override

You can always override the automatic title by specifying one in frontmatter:

```markdown
---
title: My Custom Title
---
```

## Registered Attributes

The plugin modifies:
- `title`: Set to the generated title if none exists in frontmatter
"""

from pathlib import Path

from markata.hookspec import hook_impl, register_attr


@hook_impl
@register_attr()
def pre_render(markata) -> None:
    for article in markata.filter('title==""'):
        article["title"] = (
            Path(article["path"]).stem.replace("-", " ").replace("_", " ").title()
        )
