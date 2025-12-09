---
date: 2025-12-09
description: "The plugin handles file discovery using glob patterns. It determines
  which files should be processed by Markata based on configured patterns. A glob
  is a\u2026"
published: false
slug: markata/plugins/glob
title: glob.py


---

---

The `markata.plugins.glob` plugin handles file discovery using glob patterns.
It determines which files should be processed by Markata based on configured patterns.

> A glob is a string that can be used to match files and directories.
> https://docs.python.org/3/library/pathlib.html#pathlib.Path.glob

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.glob",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.glob",
]
```

## Configuration

Configure glob patterns and behavior in your `markata.toml`:

```toml
[markata.glob]
# Single pattern
glob_patterns = "**/*.md"

# Or multiple patterns
glob_patterns = [
    "content/**/*.md",
    "pages/**/*.md",
    "posts/**/*.md"
]

# Control .gitignore integration
use_gitignore = true  # Set to false to process gitignored files
```

## Configuration Options

- `glob_patterns`: String or list of strings specifying which files to process
  - Default: `["**/*.md"]` (all markdown files in any subdirectory)
  - Supports any valid glob pattern
- `use_gitignore`: Boolean controlling whether to respect .gitignore rules
  - Default: `true`
  - When true, files matching .gitignore patterns are skipped

## Functionality

The glob plugin runs early in the Markata pipeline to discover files for processing.
It supports:
- Single or multiple glob patterns
- .gitignore integration
- Recursive directory searching
- Common glob patterns (*, **, ?, [...])

## Dependencies

This plugin has no dependencies on other Markata plugins.

---