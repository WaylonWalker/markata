---
title: "Guide 35: Themes"
slug: docs/guides/guide-035
date: 2024-11-07
published: true
description: "Complete guide to Themes in markata-go."
tags:
  - cache
  - go
  - documentation
---

# Themes Guide

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.

## Getting Started

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.

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

### Go Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```

## Advanced Usage

Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis.

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
