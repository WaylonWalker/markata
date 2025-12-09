---
date: 2025-12-09
description: "The plugin generates analytics and contribution visualizations for your
  Markata site. It creates a contributions heatmap similar to GitHub's contribution\u2026"
published: false
slug: markata/plugins/analytics
title: analytics.py


---

---

The `markata.plugins.analytics` plugin generates analytics and contribution visualizations
for your Markata site. It creates a contributions heatmap similar to GitHub's contribution
graph and provides post statistics.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.analytics",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.analytics",
]
```

## Configuration

Configure analytics behavior in your `markata.toml`:

```toml
[markata.analytics]
# Maximum scale for contribution heatmap
contributions_max_post_scale = 5

# Color map for the heatmap (any matplotlib colormap)
contributions_cmap = "rocket"

# Filter to apply when generating analytics
filter = ""  # e.g. "draft != True" to exclude drafts
```

## Functionality

## Contribution Heatmap

Generates a heatmap visualization showing:
- Post frequency over time
- Relative activity levels using configurable color scales
- Year-over-year contribution patterns

## Post Statistics

Provides analytics including:
- Total number of posts
- Posts per month/year
- Word count statistics
- Reading time estimates

## Output Files

The plugin generates these files in your output directory:
- `contributions.svg`: The contribution heatmap
- `analytics.json`: Raw analytics data

## Dependencies

This plugin depends on:
- matplotlib (for heatmap generation)
- The `datetime` plugin for post date information

---