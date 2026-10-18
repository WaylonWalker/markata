# Markata

A fast, plugin-driven static site generator with a powerful feed system.

[![CI](https://github.com/WaylonWalker/markata/actions/workflows/ci.yml/badge.svg)](https://github.com/WaylonWalker/markata/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

**Markata** (formerly **markata-go**) is a static site generator written in Go that processes Markdown files with YAML frontmatter and generates a complete static website. It features a flexible plugin architecture, a powerful feed system for creating archives and syndication feeds, and Jinja2-like templating.

> **Status:** Markata remains a **0.x** Go project; the rename is not a 1.0 release decision. The `markata` executable is preferred and `markata-go` remains supported. The original Python project is at [WaylonWalker/markata-py](https://github.com/WaylonWalker/markata-py); this repository is the Go implementation.

### Key Features

- **Plugin-driven architecture** - Extensible system with 30+ built-in plugins
- **Powerful feed system** - Define feeds with filtering, sorting, and pagination; generate multiple output formats (HTML, RSS, Atom, JSON) from a single definition
- **Jinja2-like templates** - Familiar template syntax via pongo2 with custom filters
- **9-stage lifecycle** - Configure, Validate, Glob, Load, Transform, Render, Collect, Write, Cleanup
- **Concurrent processing** - Parallel processing with configurable worker count
- **Markdown with extensions** - GFM tables, strikethrough, task lists, admonitions, syntax highlighting, wikilinks, and table of contents generation
- **Live reload development server** - Built-in server with file watching and automatic rebuilds

## Installation

### Quick Install (Recommended)

```bash
# One-liner install script (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash

# Using jpillora/installer (Linux/macOS)
curl -sL https://i.jpillora.com/WaylonWalker/markata | bash

# Using eget
eget WaylonWalker/markata

# Using mise (installs from GitHub releases)
mise use -g github:WaylonWalker/markata
```

### Go Install

Until a tagged release contains the new entry point, install the canonical binary from `main`:

```bash
go install github.com/WaylonWalker/markata/cmd/markata@main
```

Existing historical tagged versions remain available under the former `github.com/WaylonWalker/markata-go` module path when pinned explicitly. For current source, use `go install github.com/WaylonWalker/markata/cmd/markata-go@main` to install the legacy-named executable.

### Manual Download

Download pre-built binaries from [GitHub Releases](https://github.com/WaylonWalker/markata/releases).

Available platforms: Linux (amd64, arm64, armv7), macOS (Intel, Apple Silicon), Windows (amd64), FreeBSD (amd64), Android (arm64/Termux).

See [docs/installation.md](docs/installation.md) for detailed installation instructions.

## Quick Start

```bash
# Create a new post
markata new "Hello World"

# Build the site
markata build

# Serve with live reload
markata serve
```

## Configuration

The existing `markata-go.toml`, `[markata-go]`, and `MARKATA_GO_*` configuration contracts remain supported. Do **not** rename these to `markata` yet: that namespace is also used by the historical Python product. See [the name migration guide](docs/guides/name-migration.md).

markata-go uses TOML configuration (also supports YAML and JSON). Create a `markata-go.toml` in your project root:

```toml
[markata-go]
# Site metadata
title = "My Site"
description = "A site built with markata-go"
url = "https://example.com"
author = "Your Name"

# Build settings
output_dir = "output"
templates_dir = "templates"
assets_dir = "static"

# Reading time calculation
words_per_minute = 200

# Concurrency (0 = auto based on CPU cores)
concurrency = 0

# Plugin configuration
hooks = ["default"]
disabled_hooks = []

# Content discovery
[markata-go.glob]
patterns = ["posts/**/*.md", "pages/**/*.md"]
use_gitignore = true

# Markdown extensions
[markata-go.markdown]
extensions = ["tables", "strikethrough", "autolinks", "tasklist"]

# Auto-generated descriptions
[markata-go.description]
max_length = 160

# Table of contents
[markata-go.toc]
min_level = 2
max_level = 4

# Wikilinks
[markata-go.wikilinks]
warn_broken = true
```

### Environment Variable Overrides

All configuration options can be overridden via environment variables using the `MARKATA_GO_` prefix:

```bash
MARKATA_GO_OUTPUT_DIR=dist markata-go build
MARKATA_GO_URL=https://staging.example.com markata-go build
MARKATA_GO_CONCURRENCY=4 markata-go build
```

## Content

### Frontmatter Format

Posts use YAML frontmatter to define metadata:

```yaml
---
title: "My First Post"
slug: "my-first-post"
date: 2024-01-15
published: true
draft: false
tags: ["go", "static-site"]
description: "An introduction to markata-go"
template: "post.html"
---
```

### Supported Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | - | Post title |
| `slug` | string | auto | URL-safe identifier (auto-generated from title if not set) |
| `date` | date | - | Publication date |
| `published` | bool | `false` | Whether the post is published |
| `draft` | bool | `false` | Whether the post is a draft |
| `tags` | []string | `[]` | List of tags |
| `description` | string | auto | Post description (auto-generated if not set) |
| `template` | string | `post.html` | Template file to use |

Any additional fields in frontmatter are stored in `Extra` and accessible in templates.

## Feed System

The feed system is markata-go's most powerful feature. Define feeds to create filtered, sorted, paginated collections of posts with multiple output formats.

### Feed Configuration

```toml
# Feed defaults - inherited by all feeds
[markata-go.feed_defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feed_defaults.formats]
html = true
rss = true
atom = false
json = false

[markata-go.feed_defaults.templates]
html = "feed.html"
card = "card.html"

[markata-go.feed_defaults.syndication]
max_items = 20
include_content = true
```

### Defining Feeds

```toml
# Main blog feed
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 10

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = true

# Home page (empty slug = root index.html)
[[markata-go.feeds]]
slug = ""
title = "Latest Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5

[markata-go.feeds.formats]
html = true

# Featured posts
[[markata-go.feeds]]
slug = "featured"
title = "Featured"
description = "Featured posts"
filter = "published == True and featured == True"
sort = "date"
reverse = true
```

### Multiple Output Formats

Each feed can generate multiple formats from a single definition:

| Format | Output | Description |
|--------|--------|-------------|
| `html` | `/{slug}/index.html` | Paginated HTML pages |
| `rss` | `/{slug}/rss.xml` | RSS 2.0 feed |
| `atom` | `/{slug}/atom.xml` | Atom 1.0 feed |
| `json` | `/{slug}/feed.json` | JSON Feed |
| `markdown` | `/{slug}/index.md` | Markdown output |
| `text` | `/{slug}/index.txt` | Plain text output |

## CLI Commands

### `markata-go build`

Build the static site.

```bash
markata-go build              # Standard build
markata-go build --clean      # Clean output directory and build cache first
markata-go build --dry-run    # Show what would be built
markata-go build -v           # Verbose output
markata-go build -o dist      # Custom output directory
```

### `markata-go serve`

Development server with live reload.

```bash
markata-go serve              # Serve on localhost:8000
markata-go serve -p 3000      # Custom port
markata-go serve --host 0.0.0.0  # Bind to all interfaces
markata-go serve --no-watch   # Disable file watching
markata-go serve -v           # Verbose logging
```

### `markata-go new <title>`

Create a new post.

```bash
markata-go new "My First Post"         # Creates pages/post/my-first-post.md
markata-go new "Hello World" --dir blog  # Creates blog/hello-world.md
markata-go new "Draft" --draft          # Create as draft
markata-go new "Published" --draft=false # Create as published (default)
```

### `markata-go config`

Configuration management commands.

```bash
markata-go config show           # Show resolved config (YAML)
markata-go config show --json    # Show as JSON
markata-go config show --toml    # Show as TOML

markata-go config get output_dir          # Get specific value
markata-go config get glob.patterns       # Nested values with dot notation

markata-go config validate       # Validate configuration

markata-go config init           # Create markata-go.toml
markata-go config init site.yaml # Create YAML config
markata-go config init --force   # Overwrite existing
```

### `markata-go version`

Show version information.

```bash
markata-go version           # Full version info
markata-go version --short   # Just the version number
markata-go --version         # Short version flag
```

### Global Flags

```bash
-c, --config string   Config file path (default: auto-discover)
-o, --output string   Output directory (overrides config)
-v, --verbose         Verbose output
```

## Documentation

- [Installation Guide](docs/installation.md) - Detailed installation instructions
- [Configuration Guide](docs/guides/configuration.md) - Full configuration reference
- [CLI Reference](docs/reference/cli.md) - All CLI commands and options
- [Plugin Reference](docs/reference/plugins.md) - Built-in plugins documentation

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/WaylonWalker/markata.git
cd markata-go

# Build
just build

# Run tests
just test

# Run all quality checks
just check
```

### Using just

The project uses [just](https://github.com/casey/just) for development commands:

```bash
just build          # Build with version info
just test           # Run tests
just test-race      # Run tests with race detector
just lint           # Run linters
just check          # Run all quality checks (fmt, vet, lint, test)
just snapshot       # Test goreleaser locally
just ci             # Run what CI runs
```

See the [justfile](justfile) for all available commands.

## Contributing

This project is **open source but not open contribution**. See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

If you want to build your own static site generator, feel free to fork this project or use the [spec](spec/) as a starting point for your own implementation.

## Inspiration

- [markata](https://github.com/waylonwalker/markata) - The Python SSG this project is based on
- [Hugo](https://gohugo.io/) - Fast builds, good CLI
- [Eleventy](https://www.11ty.dev/) - Plugin flexibility
- [Zola](https://www.getzola.org/) - Single binary simplicity

## License

MIT License - see [LICENSE](LICENSE) for details.

Copyright (c) 2024 Waylon Walker
