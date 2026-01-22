# Static Site Generator Specification

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

## Specification Files

| File | Description |
|------|-------------|
| [INSTALL.md](./spec/INSTALL.md) | Entry point with language/library choices |
| [SPEC.md](./spec/SPEC.md) | Core architecture, CLI, concurrency |
| [CONFIG.md](./spec/CONFIG.md) | Configuration system, file discovery, env vars, CLI |
| [THEMES.md](./spec/THEMES.md) | **Theming system, customization, built-in styles** |
| [LIFECYCLE.md](./spec/LIFECYCLE.md) | 13 build stages, incremental builds |
| [FEEDS.md](./spec/FEEDS.md) | **Feed system - the core differentiator** |
| [DEFAULT_PLUGINS.md](./spec/DEFAULT_PLUGINS.md) | All 15 built-in plugins |
| [PLUGINS.md](./spec/PLUGINS.md) | Plugin development guide |
| [DATA_MODEL.md](./spec/DATA_MODEL.md) | Post/Config schemas, querying, error types |
| [CONTENT.md](./spec/CONTENT.md) | Markdown processing, frontmatter, admonitions |
| [TEMPLATES.md](./spec/TEMPLATES.md) | Template system, engine differences |
| [OPTIONAL_PLUGINS.md](./spec/OPTIONAL_PLUGINS.md) | Optional enhancement plugins |
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
