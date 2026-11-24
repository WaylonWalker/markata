---
title: "CLI Output"
description: "Understand Markata's terminal colors, themed help, and plain output behavior."
date: 2026-07-29
published: true
tags:
  - documentation
  - cli
  - themes
---


Markata uses your configured site palette for interactive help, status output,
and section headings. Help sections use Unicode box-drawing rules sized to your
terminal, up to 80 columns.

```bash
markata-go --help
markata-go lsp --help
```

When output is redirected, or when you pass `--no-color`, Markata keeps the
same text and separators but removes ANSI color codes. Structured output such
as JSON, YAML, and TOML remains suitable for scripts.
