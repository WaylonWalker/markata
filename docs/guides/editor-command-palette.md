---
title: "Editor Command Palette"
description: "Insert special Markdown syntax in the Python, Go, and WASM workbenches."
date: 2026-10-29
published: true
tags:
  - documentation
  - markdown
  - editor
---

Select **Commands** above the Markdown editor, or press **Ctrl+K**.
On macOS, press **Cmd+K**. The **Ctrl/Cmd+Shift+P** shortcut also works.
The palette works with the Default, Vim, Emacs, and Sublime keymaps.

Search for a command, such as `callout`, `diagram`, or `highlight`.
Use the arrow keys to choose a result, then press **Enter**.
You can also select a result with a mouse or touch.
Each result shows the syntax that it inserts.

For example, the **Note box** command inserts:

```markdown
!!! note "Note"
    Write your note here.
```

The editor selects the placeholder so you can replace it immediately.

In the Python and Go browser editors, the connected language service also completes
container classes and modifiers from these commands. For example, type `::::: md-i`
to suggest `md-icon`. Inside its attribute block, type `{data-s` to suggest
`data-size`, then type `data-size="` to choose `sm` or `lg`.

Inside a Mermaid fence, the language service suggests diagram types such as
`flowchart`, `sequenceDiagram`, and `classDiagram`; flowchart directions; a
flowchart starter; and Mermaid class style properties. Press **Ctrl+Space** to
show completions at the cursor. The WASM editor keeps the full command palette
and does not currently connect to the LSP.
If you select text before opening Commands, the command uses that text instead.
For example, **Highlighted text** changes selected `important` to `==important==`.
Named Lucide icon commands replace selected text with the chosen icon.
Block commands add blank lines around their output to keep surrounding paragraphs separate.
Multiline selections receive the indentation required by admonitions.
Commands use the primary selection when the editor has multiple selections.

Use the editor's normal undo command to remove an insertion in one step.
The existing draft saving and preview update run after insertion.
Press **Escape**, select **Close**, or click outside the dialog to cancel.
The editor keeps your text and selection and receives focus again.

The catalog includes 110 authoring snippets and all 2,118 icons in the bundled Lucide pack.
Search matches command names, syntax hints, and aliases.
The palette shows up to 60 results at once. Refine the search to find more specific commands.
For example, search `Lucide: zodiac-aries` to find that icon directly.

| Family | Commands |
| --- | --- |
| Admonitions | All documented types, closed and expanded disclosures, and split columns |
| Roads | Winding, timeline, steps, snake, cycle, constellation, git, and editorial; story and reveal variants |
| Cards and layouts | Card, box, layout, columns, grid, rail, pending cards, stack, divider, connector, and story section |
| Icons and feature rows | Icon tile, icon feature row, generic Lucide shortcode, every bundled Lucide icon, and custom icon-pack shortcode |
| Metrics | Status, statistic, delta, native meter, and progress gauge |
| Interactive content | Checklist and Markdown source disclosure |
| Mermaid | Flow, sequence, class, state, entity relationship, Gantt, pie, mind map, timeline, and git diagrams |
| Chart.js | Bar, line, area, stacked bar, pie, doughnut, scatter, bubble, radar, and polar area |
| Web Awesome | Comparison, details, tabs, copy button, QR code, badge, tag, tooltip, carousel, animated image, and raw component template |
| Inline extensions | Highlight, keyboard keys, emoji, superscript, subscript, attributes, footnote, and definition list |
| Images and citations | Captioned figure, image URL, and attributed quotation |
| Site features | Wikilink, post embed, external embed card, YouTube embed, hover cards, table of contents, glossary frontmatter, and contribution calendar |
| Standard Markdown | Table, task list, code fence, inline code, bold, italic, strikethrough, and link |

## Preview and site requirements

The palette uses local assets in all three workbenches.
It remains available offline after the WASM workbench caches its assets.
Bundled Lucide icons, roads, article containers, metrics, Mermaid, and Chart.js render in the standalone preview.

Commands for site features state their requirements in the result description.
The standalone renderers do not load a site content index or the Web Awesome plugin and custom-element runtime.
Web Awesome commands insert the documented Markdown shortcuts for a configured site build.
They do not load third-party scripts in the standalone preview.
The raw component template supports components that have no Markdown shorthand.
Replace its tag with the required Web Awesome tag.

Wikilinks and post embeds need site content. External embeds need the site embed pipeline.
Hover cards and contribution calendars need their site plugins or runtime assets.
Custom icon packs need site configuration and local SVGs.
The TOC command inserts `[[toc]]` for the site table-of-contents plugin.
Use glossary frontmatter at the start of a new glossary post.
Replace sample image paths and IDs before publishing.

## Try a complete example

Open the workbench with `?example` to load the sample in the editor and open the command palette.
For example, visit `http://localhost:8766/?example` for the Go workbench.
You can also open `examples/command-palette-showcase.md` with **Open .md**.
The example combines icon cards, a steps road, a feature row, metrics, Mermaid, and a Chart.js chart.
It also shows Web Awesome source with its site requirement clearly identified.
The chart data is illustrative.

See the [Markdown guide](markdown.md), [container guide](markdown-containers.md),
[roads guide](roads.md), [icon guide](icons.md), and [Web Awesome guide](webawesome.md)
for syntax details. The [editor comparison](renderer-integration-benchmarks.md#shared-editor-features-and-runtime-differences)
explains runtime differences.

## Open workspace files

Press **Ctrl+P** (macOS: **Cmd+P**) to open the file picker in the workspace editor.
You can also run **Open file** from the command palette. Ctrl/Cmd+Shift+P still
opens commands.

Type a filename abbreviation, path, or document title. Filename matches rank
above path and title matches. Titles become searchable after background metadata
loads. Matching characters are highlighted, and full paths distinguish duplicate
filenames. The picker shows up to 60 results; type more to narrow the list.
Use Up/Down and Enter to open a file, or Escape to cancel. Opening another file
preserves the current browser draft. Source files stay unchanged on disk.

The picker uses locally bundled [fuzzysort](https://github.com/farzher/fuzzysort).
Search runs in the browser over discovered files without fetching source content.
