---
title: Editor settings and quick navigation
description: Personalize the editor, find commands quickly, and inspect performance.
date: 2026-11-13
published: true
tags: [documentation, editor]
---

Open **Editor settings** from the header, command palette, or Ctrl/Cmd+, with the VS Code keymap.
The window fills the screen on phones. Search settings or choose Appearance, Writing,
Preview, Keyboard Shortcuts, or Advanced. Preferences are stored in this browser;
these controls do not change your published site's configuration.

Appearance includes source text size and the existing site palette controls when available.
Writing includes line wrapping and the editor keymap. Preview controls
layout and synchronized scrolling. Keyboard Shortcuts contains the application
preset, search, conflict-aware shortcut recorder, and reset controls.

## Find files and commands

Ctrl/Cmd+P finds files; Ctrl/Cmd+Shift+P finds commands with the default application
keymap. Both support fuzzy matching. The file picker searches filename, path and
post title. Recent choices appear first on an empty query. With a query, history
breaks ties between similarly strong matches; exact matches always take priority.
Command history is local to this browser. File history is scoped to the workspace.
Neither stores your search terms or draft contents.

Advanced settings lets you disable remembering picks or clear command history and
recent files for the current workspace. Blocked browser storage does not prevent
opening files or executing commands; history then lasts only for the current tab.

## Diagnostics without sidebar clutter

Use **Diagnostics** in the footer or **Toggle Performance Diagnostics** in the
palette. Startup indexing details remain available in the drawer. **Copy full
report** includes startup status and up to 100 recent samples per measurement.
Agents can read `window.MarkataDiagnostics.snapshot()` directly. Missing samples
are not evidence of instant completion.

`file_open` ends when the editor receives the file; `save_acknowledgement`
ends when the server acknowledges a write, not when a deployment completes.
`edit_to_preview` measures the accepted preview update, including debounce and
render time. `file_request` and `publish_request` measure resource requests only:
they are not proof of a durable save, successful publication or deployment.
Search timing excludes result DOM rendering. Validate end-to-end latency on a real
site before making performance claims.

Publishing and diff changes are developed separately in PR #1982. This settings
and navigation pass does not alter those Git transactions.
