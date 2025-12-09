---
date: 2025-12-09
description: "The plugin copies static assets (images, CSS, JavaScript, etc.) from
  your assets directory to the output directory during the build process. Installation
  This\u2026"
published: false
slug: markata/plugins/copy-assets
title: copy_assets.py


---

---

The `markata.plugins.copy_assets` plugin copies static assets (images, CSS, JavaScript, etc.)
from your assets directory to the output directory during the build process.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.copy_assets",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.copy_assets",
]
```

## Configuration

Configure asset directories in your `markata.toml`:

```toml
[markata]
# Directory containing your static assets
assets_dir = "assets"

# Directory where assets will be copied
output_dir = "markout"
```

## Functionality

## Asset Copying

The plugin:
1. Checks if the configured assets directory exists
2. Recursively copies all files and directories from assets_dir to output_dir
3. Preserves directory structure
4. Updates existing files if they've changed
5. Maintains any existing files in the output directory

## Usage Example

Place static assets in your assets directory:
```
assets/
  ├── css/
  │   └── style.css
  ├── js/
  │   └── main.js
  └── images/
      └── logo.png
```

These will be copied to:
```
markout/
  ├── css/
  │   └── style.css
  ├── js/
  │   └── main.js
  └── images/
      └── logo.png
```

---