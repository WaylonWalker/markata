---
title: "Find heavy posts during builds"
description: "Use markata-go page-weight warnings and JSON reports to prioritize post performance work."
published: true
tags:
  - performance
  - build
---


Run a normal build:

```sh
markata-go build --clean
```

After a successful build, markata-go reports unusually large HTML and writes
`output/.markata/page-weight.json` (or under your configured `output_dir`).
The report is intended for agents, scripts and humans. It sorts the heaviest
pages first and includes each page's path, byte breakdown, stable warning codes,
budgets, and suggested improvements.

```sh
# Disable diagnostics for a build.
markata-go build --page-weight=false

# Warn on posts larger than 96 KiB of uncompressed HTML source.
markata-go build --page-weight-html-kb=96
```

The default post warning is 160 KiB of uncompressed HTML; other page types get
320 KiB. Embedded JSON over 32 KiB, inline JS over 24 KiB, and sidebar JSON
over 16 KiB each produce additional warnings.

## Working through a warning

Give your coding agent the affected `path`, `findings[].code`,
`findings[].bytes`, and `findings[].advice` from the JSON. Compare clean
builds before and after changing templates; do not hide warnings by trimming
useful article or sidebar content. Prefer loading supplementary features on
interaction, using static files for shared logic, and removing duplicated data
from HTML.

**Important:** These are source-size checks, not measured compressed transfer
or Lighthouse metrics. Validate consequential optimizations with cold-browser
Lighthouse and browser network traces before claiming a performance win.
The report is local and does not upload visitor or site data.
