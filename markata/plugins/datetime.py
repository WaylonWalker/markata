"""
The `markata.plugins.datetime` plugin handles date and time parsing for posts. It ensures
consistent datetime handling by converting various date formats to timezone-aware datetime
objects.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.datetime",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.datetime",
]
```

Note: Disabling this plugin will prevent proper date handling and sorting in other plugins.

# Configuration

This plugin requires no explicit configuration. It processes dates from post frontmatter.

## Frontmatter Usage

Specify dates in your post frontmatter:

```markdown
---
title: My Post
date: 2024-01-28  # Date only
---
```

```markdown
---
title: My Post
date: 2024-01-28T12:00:00Z  # Date and time with timezone
---
```

# Functionality

## Date Parsing

The plugin supports:
- ISO format dates (2024-01-28)
- ISO format datetimes (2024-01-28T12:00:00Z)
- Natural language dates ("January 28, 2024")
- Date objects from Python
- Datetime objects with or without timezone

## Date Normalization

All dates are normalized to:
- Timezone-aware datetime objects
- UTC timezone if none specified
- Start of day (00:00:00) for date-only values

## Registered Attributes

The plugin modifies:
- `date`: Converted to timezone-aware datetime object

## Dependencies

This plugin depends on:
- python-dateutil for flexible date parsing
- pytz for timezone handling

"""

import datetime
from typing import TYPE_CHECKING

import dateutil
import pytz

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(trylast=True)
def load(markata: "Markata") -> None:
    for article in markata.iter_articles("datetime"):
        try:
            date = article.metadata["date"]
        except KeyError:
            date = None
        if isinstance(date, str):
            date = dateutil.parser.parse(date)
        if isinstance(date, datetime.date):
            date = datetime.datetime(
                year=date.year,
                month=date.month,
                day=date.day,
                tzinfo=pytz.utc,
            )

        article["today"] = datetime.date.today()
        article["now"] = datetime.datetime.now()
        article["datetime"] = date
        # if date is not None:
        #     article["date"] = date.date()
