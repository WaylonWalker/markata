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

## CMS navigation and mobile workflow

On desktop viewports, the top bar keeps the active page title, save status, a
sidebar toggle, a **Write / Split / Preview** layout switcher, **+ New**, and
**Publish post** within reach without opening settings.

On narrow screens (`<= 768px`), the editor switches to a full-height mobile CMS
shell with a bottom navigation bar (**Content**, **Write**, **Preview**, and
**New**):

- **Content** shows the full-screen browser for Recent, Files, Feeds, Tags, and
  Changes. Opening a file transitions directly to **Write**.
- **Write** gives the full viewport to Markdown editing, Page Properties, and
  Drawings, plus a mobile formatting bar (**Undo**, **Redo**, **Head**, **List**,
  **Code**, **Link**, **Cmd**, **Done**) that stays above the virtual keyboard.
- **Preview** renders the live site preview full-screen. Swipe left/right to
  move between **Write** and **Preview**, pull down at the top to open the
  Command Palette, or tap the **⌘** button in the top bar from any view.
- **New** opens the content template picker bottom sheet directly.

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

See the [Go workbench guide](go-workbench.md) for publishing, Changes review, and
comparing unsynced copies.


### Share photos from Android

Install the editor PWA from an HTTPS origin and choose **Markata** from Android's Share menu. The share target accepts text, URLs, and up to eight PNG/JPEG/GIF/WebP photos (15 MiB per image, 48 MiB total). Shared photos are staged only in this browser until you review and save them in **Quick Note**. They are written to the workspace via the existing new-post image upload route, using local storage for multi-photo captures, or the site's configured Dropper provider for a single photo (multi-photo Dropper is rejected until rollback/cleanup is supported); the first photo becomes the cover and all photos are added to the Markdown body. Canceling or reloading before saving preserves the staged share for up to 24 hours as long as the share URL is retained. Daily-note append does not support photo attachments. Actual publishing remains an explicit action.
