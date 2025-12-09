---
date: 2025-12-09
description: "The plugin automatically generates favicons in multiple sizes from a
  single source image. This ensures your site has appropriate icons for different
  devices\u2026"
published: false
slug: markata/plugins/icon-resize
title: icon_resize.py


---

---

The `markata.plugins.icon_resize` plugin automatically generates favicons in multiple
sizes from a single source image. This ensures your site has appropriate icons for
different devices and platforms.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.icon_resize",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.icon_resize",
]
```

## Configuration

Configure icon settings in your `markata.toml`:

```toml
[markata]
# Directory for output files
output_dir = "markout"

# Directory containing static assets
assets_dir = "static"

# Path to source icon file
icon = "static/icon.png"

# Optional: List of icon sizes to generate
# Default sizes are provided if not specified
icon_sizes = [16, 32, 48, 128, 180, 192, 512]
```

## Functionality

## Icon Generation

The plugin:
1. Loads the source icon image
2. Resizes it to multiple standard sizes
3. Saves each size as a separate PNG file
4. Generates appropriate HTML meta tags

## Default Sizes

If not configured, generates these standard sizes:
- 16x16: Classic favicon size
- 32x32: Modern favicon size
- 48x48: Windows taskbar
- 128x128: Chrome Web Store
- 180x180: Apple Touch icon
- 192x192: Android homescreen
- 512x512: PWA splash screen

## Output Files

Generated files follow this pattern:
- `favicon-{size}x{size}.png`
- `favicon.ico` (16x16 and 32x32 combined)
- `apple-touch-icon.png` (180x180)

## HTML Integration

The plugin adds appropriate meta tags to your HTML:
```html
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
```

## Dependencies

This plugin depends on:
- Pillow (PIL) for image manipulation

---