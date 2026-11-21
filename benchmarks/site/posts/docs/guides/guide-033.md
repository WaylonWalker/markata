---
title: "Guide 33: Templates"
slug: docs/guides/guide-033
date: 2024-09-05
published: true
description: "Complete guide to Templates in markata-go."
tags:
  - api
  - database
  - documentation
---

# Templates Guide

Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.

## Getting Started

Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt.

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

### Rust Example

```rust
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

#[derive(Debug, Clone)]
pub struct Cache<T> {
    data: Arc<RwLock<HashMap<String, T>>>,
}

impl<T: Clone> Cache<T> {
    pub fn new() -> Self {
        Self {
            data: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn get(&self, key: &str) -> Option<T> {
        let data = self.data.read().await;
        data.get(key).cloned()
    }

    pub async fn set(&self, key: String, value: T) {
        let mut data = self.data.write().await;
        data.insert(key, value);
    }
}
```

## Advanced Usage

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.

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
