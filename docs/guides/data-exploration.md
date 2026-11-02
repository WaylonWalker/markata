---
title: "Data Exploration"
description: "Explore posts, tags, and feeds from the CLI"
date: 2026-02-16
published: true
tags:
  - documentation
  - guides
  - cli
---


Use `list` for scriptable output and `tui` for interactive browsing.

## List Posts

```bash
# Default table output
markata-go list posts

# Path-only output for piping
markata-go list posts --format path

# Filter posts with an expression
markata-go list posts --filter "published == true and 'go' in tags"

# Sort by title, ascending
markata-go list posts --sort title --order asc

# List posts for a feed
markata-go list posts --feed blog
```

## List Tags

```bash
# Table output
markata-go list tags

# JSON output for scripts
markata-go list tags --format json
```

## List Feeds

```bash
# Table output
markata-go list feeds

# CSV output for spreadsheets
markata-go list feeds --format csv

# List posts for a feed
markata-go list feeds posts blog

# List post paths for a feed
markata-go list feeds posts blog --format path
```

## Path-Only Output

The `path` format prints one value per line for easy piping:

```bash
# Open each post in your editor
markata-go list posts --format path | xargs -n 1 $EDITOR

# Create a tag list file
markata-go list tags --format path > tags.txt
```

When `--site-dir` or `MARKATA_GO_SITE_DIR` selects a site, post paths are
absolute. This lets agents and scripts edit results from any directory:

```bash
markata-go --site-dir ~/sites/blog search "configuration" --format path
# Edit the returned Markdown path directly.
```

## Interactive Browsing

For a full-screen view with filtering and sorting:

```bash
markata-go tui
```

## Cache

`list` and `tui` use a persistent cache at `.markata/cache/list.json`. It updates automatically when files change. Delete the file to force a full refresh.
