---
title: "Guide 24: Configuration"
slug: docs/guides/guide-024
date: 2024-12-24
published: true
description: "Complete guide to Configuration in markata-go."
tags:
  - tutorial
  - guide
  - documentation
---

# Configuration Guide

Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt.

## Getting Started

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.

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

### Sql Example

```sql
-- Create tables for blog application
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_published ON posts(published);

-- Insert sample data
INSERT INTO posts (title, slug, content, published) VALUES
    ('Hello World', 'hello-world', 'Welcome to my blog!', TRUE),
    ('Second Post', 'second-post', 'More content here.', TRUE);
```

## Advanced Usage

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.

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
