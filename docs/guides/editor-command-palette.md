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

Select **Commands** above the Markdown editor, or press **Ctrl+Shift+P**.
On macOS, press **Cmd+Shift+P**.
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
If you select text before opening Commands, the command uses that text instead.
For example, **Highlighted text** changes selected `important` to `==important==`.
Block commands add blank lines around their output to keep surrounding paragraphs separate.
Multiline selections receive the indentation required by admonitions.
Commands use the primary selection when the editor has multiple selections.

Use the editor's normal undo command to remove an insertion in one step.
The existing draft saving and preview update run after insertion.
Press **Escape**, select **Close**, or click outside the dialog to cancel.
The editor keeps your text and selection and receives focus again.

Available commands include:

| Commands | Inserted syntax |
| --- | --- |
| Note, tip, warning, aside | `!!!` admonitions |
| Collapsible note | `??? note` |
| Highlighted text | `==text==` |
| Keyboard keys | `++Ctrl+Alt+Del++` |
| Definition list | Term followed by `: Definition` |
| Mermaid diagram | A `mermaid` code fence |
| Chart.js chart | A `chartjs` code fence with sample JSON |
| Table | Two columns with a header row |
| Task list | `- [ ] Task` |
| Code block | A `text` code fence |
| Bold, italic, link, image URL | Standard Markdown formatting |

The palette uses local assets in all three workbenches.
It remains available offline after the WASM workbench caches its assets.
Site-dependent commands, such as post embeds, require a site index and are outside this initial catalog.

See the [Markdown guide](markdown.md) for syntax details and the
[editor comparison](renderer-integration-benchmarks.md#shared-editor-features-and-runtime-differences) for runtime differences.
