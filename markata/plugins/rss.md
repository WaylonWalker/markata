---
date: 2025-12-09
description: "[DEPRECATED] The plugin is deprecated and will be removed in a future
  version. Please use instead, which provides more comprehensive feed generation\u2026"
published: false
slug: markata/plugins/rss
title: rss.py


---

---

[DEPRECATED] The `markata.plugins.rss` plugin is deprecated and will be removed in a
future version. Please use `markata.plugins.feeds` instead, which provides more
comprehensive feed generation capabilities.

## Installation

This plugin is deprecated. Use `markata.plugins.feeds` instead:

```toml
hooks = [
    "markata.plugins.feeds",  # Use this instead
    # "markata.plugins.rss",  # Deprecated
]
```

# Migration Guide

To migrate to the new feeds plugin:

1. Remove rss plugin from hooks:
```toml
# Remove or comment out
# "markata.plugins.rss"
```

2. Add feeds plugin:
```toml
hooks = [
    "markata.plugins.feeds"
]
```

3. Update configuration:
```toml
[markata.feeds]
# RSS feed configuration
rss = { output = "rss.xml" }

# Optional: Add other feed formats
atom = { output = "atom.xml" }
json = { output = "feed.json" }
```

See the feeds plugin documentation for more configuration options.

# Legacy Configuration

If you must continue using this plugin temporarily, configure in `markata.toml`:

```toml
[markata]
url = "https://example.com"
title = "Site Title"
author_name = "Author Name"
author_email = "author@example.com"
icon = "favicon.ico"
lang = "en"
```

# Dependencies

This plugin depends on:
- feedgen for RSS generation
- pytz for timezone handling

WARNING: This plugin is deprecated and will be removed in a future version.
Please migrate to `markata.plugins.feeds` as soon as possible.

---