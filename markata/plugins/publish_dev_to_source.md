---
date: 2025-12-09
description: "The plugin enables synchronization between your Markata posts and dev.to
  articles. It handles frontmatter translation and content formatting for dev.to\u2026"
published: false
slug: markata/plugins/publish-dev-to-source
title: publish_dev_to_source.py


---

---

The `markata.plugins.publish_dev_to_source` plugin enables synchronization between
your Markata posts and dev.to articles. It handles frontmatter translation and content
formatting for dev.to compatibility.

## Installation

This plugin is built-in but not enabled by default. Add it to your plugins list:

```toml
hooks = [
    "markata.plugins.publish_dev_to_source",
]
```

## Uninstallation

Remove the plugin from your hooks list in `markata.toml`:

```toml
hooks = [
    # Remove or comment out the line below
    # "markata.plugins.publish_dev_to_source",
]
```

## Configuration

No explicit configuration is required. The plugin automatically processes posts with
dev.to-specific frontmatter.

## Functionality

## Frontmatter Mapping

Automatically maps these frontmatter fields:
- title
- published
- description
- tags
- canonical_url
- cover_image
- series

## Content Processing

The plugin:
1. Joins consecutive text lines for dev.to formatting
2. Preserves code blocks and lists
3. Maintains markdown compatibility
4. Updates canonical URLs

## Post Model

Extends the base Post model with:
- dev.to-specific fields
- Validation rules
- Performance optimizations

## Publishing Flow

1. Processes posts after rendering
2. Updates frontmatter
3. Formats content
4. Saves changes to source files

## Dependencies

This plugin depends on:
- python-frontmatter for YAML handling
- pydantic for model validation

---