---
title: "Guide 36: Feeds"
slug: docs/guides/guide-036
date: 2024-12-08
published: true
description: "Complete guide to Feeds in markata-go."
tags:
  - go
  - python
  - documentation
---

# Feeds Guide

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.

## Getting Started

Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis.

### Installation

```bash
go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
```

### Configuration

```toml
[markata-go]
title = "My Site"
output_dir = "public"
```

### Python Example

```python
import asyncio
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Post:
    title: str
    content: str
    tags: List[str]
    published: bool = False

async def fetch_posts(limit: Optional[int] = None) -> List[Post]:
    """Fetch posts from the database asynchronously."""
    await asyncio.sleep(0.1)  # Simulate I/O
    posts = [
        Post("Hello", "World", ["intro"], True),
        Post("Second", "Content", ["update"], True),
    ]
    return posts[:limit] if limit else posts
```

## Advanced Usage

Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.

## Checklist

- [x] Initial setup
- [x] Basic configuration
- [ ] Advanced features
- [ ] Production deployment

## Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| enabled | bool | true | Enable the feature |
| timeout | int | 30 | Timeout in seconds |
| verbose | bool | false | Enable verbose output |

## See Also

- [Configuration Guide](/docs/guides/configuration/)
- [Plugin Development](/docs/guides/plugins/)
