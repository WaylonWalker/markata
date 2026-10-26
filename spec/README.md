---
title: "Static Site Generator Specification"
description: "Specification index for markata-go behavior and maintainer contracts."
date: 2026-10-01
published: true
tags:
  - specification
---

# Static Site Generator Specification

The [galleries specification](spec/GALLERIES.md) defines built-in gallery pages, photo-grid feeds, and initial post title sizing.

A language-agnostic specification for building plugin-driven static site generators.

## Philosophy

This project follows the **"spec-as-product"** approach inspired by [whenwords](https://github.com/dbreunig/whenwords): the specification and test suite ARE the product. Implementations in any language should produce identical results.

**Two choices to make:**
1. `[name]` - What to call your SSG
2. `[language]` - Python, TypeScript, Go, or Rust

Everything else is specified.

## Quick Start

1. Read [INSTALL.md](./spec/INSTALL.md) - Fill in `[name]` and `[language]`
2. Hand the spec files to an AI agent
3. Run the test suite to verify your implementation

## Documentation Requirements

**Every feature in the spec MUST have corresponding user documentation.**

The spec defines *what* and *how* - technical behavior, data models, edge cases.
The docs explain *why* and *usage* - user-friendly guides, examples, tutorials.

| Spec File | Required User Doc |
|-----------|-------------------|
| `EDITORIAL_HISTORY.md` | `docs/guides/editorial-history.md` |
| `EDITORIAL_AUTHORING.md` | `docs/guides/editorial-authoring.md` |
| `CONFIG.md` | `docs/guides/configuration.md` |
| `AGENTS.md` | `docs/guides/agent-skills.md`, `docs/reference/cli.md` |
| `THEMES.md` | `docs/guides/themes.md` |
| `ADAPTIVE_LOADING.md` | `docs/guides/adaptive-loading.md` |
| `FEEDS.md` | `docs/guides/feeds.md` |
| `TEMPLATES.md` | `docs/guides/templates.md` |
| `CONTENT.md` | `docs/guides/markdown.md`, `docs/guides/frontmatter.md` |
| `REDIRECTS.md` | `docs/guides/nginx-redirects.md`, `docs/reference/redirects.md` |
| `PLUGINS.md` | `docs/guides/plugin-development.md`, `docs/reference/plugins.md` |
| `SPEC.md` (CLI) | `docs/reference/cli.md` |
| `CONTAINERS.md` | `docs/guides/deployment/docker.md` |
| `CI.md` | `docs/guides/maintainer-ci.md` |
| `ROADS.md` | `docs/guides/roads.md` |

**Documentation lives in `docs/` and is built as part of the site itself.**

## Specification Files

| File | Description |
|------|-------------|
| [INSTALL.md](./spec/INSTALL.md) | Entry point with language/library choices |
| [SPEC.md](./spec/SPEC.md) | Core architecture, CLI, concurrency |
| [EDITORIAL_HISTORY.md](./spec/EDITORIAL_HISTORY.md) | Editorial history, lazy source comparisons, and editor credits |
| [EDITORIAL_AUTHORING.md](./spec/EDITORIAL_AUTHORING.md) | Isolated attributed commits for humans and agents |
| [CLI_LIST.md](./spec/CLI_LIST.md) | CLI list command and output formats |
| [CLI_SITE_DIR.md](./spec/CLI_SITE_DIR.md) | Cross-directory site selection and agent workflows |
| [CLI_UX.md](./spec/CLI_UX.md) | Shared CLI UX rules for streams, color, and prompts |
| [CLI_TRIAGE.md](./spec/CLI_TRIAGE.md) | Agent-first troubleshooting, safe diagnostics, and pipeable prompt handoff |
| [SERVE_CONTROL_CENTER.md](./spec/SERVE_CONTROL_CENTER.md) | Shared local serve session runtime, resources, and terminal/web clients |
| [LSP.md](./spec/LSP.md) | Language Server Protocol integration and setup guidance |
| [CONTAINERS.md](./spec/CONTAINERS.md) | Container images and runtime environments |
| [CI.md](./spec/CI.md) | Product repository validation, container publishing, and benchmark policy |
| [CONFIG.md](./spec/CONFIG.md) | Configuration system, file discovery, env vars, CLI |
| [AGENTS.md](./spec/AGENTS.md) | Agent skill packaging, install targets, and CLI |
| [THEMES.md](./spec/THEMES.md) | **Theming system, customization, built-in styles** |
| [ROADS.md](./spec/ROADS.md) | Scroll-drawn roads built from `::: road` containers |
| [LIFECYCLE.md](./spec/LIFECYCLE.md) | 13 build stages, incremental builds |
| [BUILD_DAG.md](./spec/BUILD_DAG.md) | Immutable graph declarations, canonical diagnostics, deterministic serial ordering |
| [FEEDS.md](./spec/FEEDS.md) | **Feed system - the core differentiator** |
| [REDIRECTS.md](./spec/REDIRECTS.md) | Native nginx redirects and portable HTML fallbacks |
| [DEFAULT_PLUGINS.md](./spec/DEFAULT_PLUGINS.md) | All 15 built-in plugins |
| [PLUGINS.md](./spec/PLUGINS.md) | Plugin development guide |
| [DATA_MODEL.md](./spec/DATA_MODEL.md) | Post/Config schemas, querying, error types |
| [CONTENT_INDEX.md](./spec/CONTENT_INDEX.md) | Versioned derived metadata artifact and parser contract |
| [IMAGE_INDEX.md](./spec/IMAGE_INDEX.md) | Canonical image inventory, JSON artifact, and authoring page |
| [CONTENT.md](./spec/CONTENT.md) | Markdown processing, frontmatter, admonitions, wikilinks |
| [CONTENT_DIAGNOSTICS.md](./spec/CONTENT_DIAGNOSTICS.md) | Deterministic content dispositions and frontmatter diagnostics |
| [TEMPLATES.md](./spec/TEMPLATES.md) | Template system, engine differences |
| [OPTIONAL_PLUGINS.md](./spec/OPTIONAL_PLUGINS.md) | Optional enhancement plugins |
| [VERSIONING.md](./spec/VERSIONING.md) | Versioning scheme and release process |
| [tests.yaml](./spec/tests.yaml) | 350+ test cases for verification |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI                                  │
├─────────────────────────────────────────────────────────────┤
│                    Core Orchestrator                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Lifecycle  │  │   Plugin    │  │    Data     │         │
│  │  Manager    │  │   Manager   │  │   Access    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│                       Plugins                                │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐    │
│  │ glob │ │ load │ │render│ │feeds │ │ save │ │ ...  │    │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Key Design Principles

1. **Feeds are the core feature** - One definition, many outputs (HTML, RSS, Atom, JSON, Markdown, Text)
2. **Core provides only what plugins cannot** - Lifecycle orchestration, plugin loading, caching, data access
3. **Plugins are first-class** - Single file can hook any stage, extend schemas, add CLI commands
4. **Content is queryable data** - `filter("published == True and 'python' in tags")`
5. **Configuration is hierarchical** - CLI > local config > global config > defaults
6. **Markdown is enhanced, not replaced** - Frontmatter, admonitions, wikilinks, Jinja-in-Markdown
7. **Builds are deterministic and cacheable** - Same inputs = same outputs, incremental rebuilds

## Lifecycle Stages

```
CONFIGURATION PHASE
  config_model → post_model → create_models → load_config → configure → validate_config

CONTENT PHASE  
  glob → load → pre_render → render → post_render

OUTPUT PHASE
  save → teardown
```

## Language Support

| Language | Templating | Markdown | Plugin System |
|----------|------------|----------|---------------|
| Python | Jinja2 | markdown-it-py | pluggy |
| TypeScript | Nunjucks | unified/remark | Custom hooks |
| Go | pongo2 | goldmark | Custom hooks |
| Rust | Tera | pulldown-cmark | Custom hooks |

## Running Tests

The `tests.yaml` file contains 350+ test cases. Implementations should:

1. Parse each test case
2. Run the described operation
3. Compare output to expected result
4. Report pass/fail

```bash
# Example test runner (implement for your language)
[name] test --spec tests.yaml
```

## Contributing

To improve the specification:

1. Identify a gap or inconsistency
2. Propose changes with rationale
3. Add test cases that exercise the new behavior
4. Update relevant spec files

## Inspiration

- [markata](https://github.com/waylonwalker/markata) - The Python SSG this spec is based on
- [whenwords](https://github.com/dbreunig/whenwords) - The "spec-as-product" philosophy
- [Hugo](https://gohugo.io/) - Fast builds, good CLI
- [Eleventy](https://www.11ty.dev/) - Plugin flexibility
- [Zola](https://www.getzola.org/) - Single binary simplicity

## License

This specification is open source. Implementations may use any license.
