---
title: "Rendering migration merge-readiness evidence"
description: "Validation record for the versioned rendering contract migration."
date: 2026-08-10
published: true
tags:
  - documentation
  - rendering
  - migration
---


This record closes the rendering-contract acceptance checklist. It does not
claim full application WCAG compliance or exact cross-browser pixel parity.

## Commands and outcomes

### Markata-Go

All commands passed:

```text
go run ./scripts/rendering-contract --check
go test ./...
go vet ./...
go build ./cmd/markata-go
go test ./pkg/renderingcontract -run TestContract_FinalRenderSweepCoversEveryCanonicalPalette -count=1
git diff --check
```

The final CSS projection sweep audited 120 canonical variants: 120 audited,
120 passing, 0 failing. The contract contains 60 light/dark families.

### Theme Lab

All commands passed:

```text
bun test                         # 10 pass, 0 fail
just browser-regression          # 3 pass, 0 fail, 45 assertions
git diff --check
```

The browser matrix uses `brush-poster`, waits for the Knewave font load,
checks H1-H6 fixture geometry and bounds, and captures six generated
screenshots: none, splatter 0, splatter half, splatter full, dry-brush half,
and screenprint half. The zero/none geometry is equal, intermediate/full and
recipe masks differ, and the headings remain inside the preview bounds.

### Plaindown

Focused and browser checks passed:

```text
bun test tests/markata-theme.test.js tests/editor-theme-accessibility.test.js tests/render-plan-fixtures.test.js
just browser-regression          # 22 pass, 0 fail, 130 assertions across 2 browser files
git diff --check
```

The browser files include two Plaindown rendering-regression tests plus the
existing browser smoke tests. The suite asserts both Site and Editor mode
states, including Site-mode canonical variables and Editor-mode ownership.

The full suite is a baseline, not a rendering gate:

```text
current worktree: 415 pass, 18 fail, 2 errors
origin/main:     389 pass, 20 fail, 4 errors
```

The 16 named assertion failures are identical in both runs: selectable-hunk
review, draft/discard persistence, file-tree actions, navigation aliases,
frontmatter media completion, semantic preview routing, remote update actions,
wikilink resolution, stale-route recovery, draft deep links, live-preview
debouncing, and repository-content sweep checks. The two unhandled errors are
environment/dependency failures (`@codemirror/*`/vendor `fflate`), not
rendering changes.

## Evidence boundaries

- Go, Theme Lab, and Plaindown intentionally use different rendering recipes:
  Go emits generated site CSS, Theme Lab is the visual reference, and
  Plaindown owns editor chrome while Site mode projects the contract. Their
  shared evidence is normalized state, complete render plans, semantic roles,
  and final CSS variables—not identical implementation details.
- Plaindown browser endpoint checks exercise texture, heading-mask, and motif
  values at 0, 0.5, and 1. They verify readable heading contrast and that
  motif color mix changes the image, not layer opacity.
- Plaindown editor source-to-final-variable accessibility uses the separate
  editor projection audit for every editor-owned theme plus representative
  browser checks of final body, muted, and boundary variables.
- Chromium is available and passes. Firefox/WebKit smoke execution is not
  claimed: installed browser binaries are older than the Playwright package's
  requested revisions, and no large download was performed.

## Repository hygiene and payload

The intended files are the source, tests, generated contract projections,
docs, checklist, and Theme Lab `package.json`/`bun.lock`. Theme Lab
`node_modules/` and `test-results/` are ignored and are not intended for the
PR. Browser screenshots are transient review artifacts.

The contract JSON and generated consumer contract are about 34.7 KB and 39.6
KB respectively. Theme Lab's HTML is about 160 KB; Plaindown's rendering
contract and theme module are about 39.7 KB and 38.8 KB. Theme Lab declares all
font families in one CSS import for its catalog, but font-face files are
requested by browser use; no 120-palette asset requests were observed.

The public authoring names audited and frozen are `palette`, `aesthetic`,
`fontpack`, the nested `texture`, `heading_texture`, and `motif` fields, and
`variables`. No accepted public field is silently ignored.
