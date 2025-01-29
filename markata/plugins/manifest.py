"""
The `markata.plugins.manifest` plugin generates a Web App Manifest (manifest.json) file
for your site. This enables Progressive Web App (PWA) features and provides metadata
for mobile devices.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.manifest",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.manifest",
]
```

# Configuration

Configure manifest settings in your `markata.toml`:

```toml
[markata]
# Site information
site_name = "My Blog"              # Full site name
short_name = "Blog"                # Short name for app icon
start_url = "/"                    # Starting URL when launched
display = "standalone"             # Display mode (standalone/fullscreen/etc)
background_color = "#ffffff"       # App background color
theme_color = "#4a9eff"           # Theme color for browser UI
description = "My awesome blog"    # Site description

# Output directory for manifest.json
output_dir = "markout"
```

# Functionality

## Manifest Generation

The plugin:
1. Collects site configuration
2. Integrates icon information from icon_resize plugin
3. Generates a standards-compliant manifest.json
4. Places it in the output directory

## Generated Output

Creates a manifest.json file with:
```json
{
    "name": "My Blog",
    "short_name": "Blog",
    "start_url": "/",
    "display": "standalone",
    "background_color": "#ffffff",
    "theme_color": "#4a9eff",
    "description": "My awesome blog",
    "icons": [
        {
            "src": "/favicon-192x192.png",
            "sizes": "192x192",
            "type": "image/png"
        },
        // ... other icon sizes
    ]
}
```

## Integration

Works with:
- icon_resize plugin for PWA icons
- service_worker plugin for offline support
- HTML templates for manifest linking

## Dependencies

This plugin works best with:
- The `icon_resize` plugin for PWA icons
- The `service_worker` plugin for full PWA support
"""

import json
from pathlib import Path
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata.plugins.icon_resize import MarkataIcons


@hook_impl
def render(markata: "MarkataIcons") -> None:
    icons = markata.icons if "icons" in markata.__dict__ else []
    manifest = {
        "name": markata.config.site_name,
        "short_name": markata.config.short_name,
        "start_url": markata.config.start_url,
        "display": markata.config.display,
        "background_color": str(markata.config.background_color),
        "theme_color": str(markata.config.theme_color),
        "description": markata.config.description,
        "icons": icons,
    }
    filepath = Path(markata.config["output_dir"]) / "manifest.json"
    filepath.touch(exist_ok=True)
    with open(filepath, "w+") as f:
        json.dump(manifest, f, ensure_ascii=True, indent=4)
