"""
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
"""

import filecmp
import os
import shutil
from pathlib import Path
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def copy_if_changed(src_dir, dst_dir):
    src_dir = Path(src_dir)
    dst_dir = Path(dst_dir)
    for root, _, files in os.walk(src_dir):
        rel_root = Path(root).relative_to(src_dir)
        dst_root = dst_dir / rel_root
        dst_root.mkdir(parents=True, exist_ok=True)
        for fname in files:
            src_file = Path(root) / fname
            dst_file = dst_root / fname
            if not dst_file.exists() or not filecmp.cmp(
                src_file, dst_file, shallow=False
            ):
                shutil.copy2(src_file, dst_file)  # preserves metadata


@hook_impl
def save(markata: "Markata") -> None:
    with markata.console.status("copying assets", spinner="aesthetic", speed=0.2):
        if markata.config.assets_dir.exists():
            copy_if_changed(markata.config.assets_dir, markata.config.output_dir)
            # shutil.copytree(
            #     markata.config.assets_dir,
            #     markata.config.output_dir,
            #     dirs_exist_ok=True,
            # )
