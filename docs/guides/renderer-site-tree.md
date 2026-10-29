---
title: "Explore an Example Site in the Renderer Editors"
description: "Navigate site files and feeds in the Python and Go workbenches."
date: 2026-10-29
published: true
tags:
  - documentation
  - editor
  - feeds
---

The [Python editor](fastapi-editor/) and [Go workbench](go-workbench/) include an example site sidebar.
Pierre `@pierre/trees` provides the file tree.

Select a Markdown file under **Files** to open its source and live preview.
Edit the source, then select another file.
When you return, the editor restores your edits.
Each file has a separate browser draft that survives a reload when browser storage is available.
The previous standalone draft remains available as **draft.md**.

Select **Feeds** in the sidebar to browse the `notes` and `garden` feeds.
Expand a feed, then select a member file.
The editor opens the same source and draft as the Files view.
Feed membership comes from the bundled example tags.
Changes to tags do not regroup the demo tree.

Select `markata-go.toml` to read the example configuration.
The editor displays configuration files as read-only text and shows a code preview.
Markdown previews omit the leading YAML frontmatter block.

The source files live in `examples/editor-site/` in the repository.
Edits stay in memory and browser storage. They do not change files on the server.
Use **Download .md** to save the current file to disk.
Use **Open .md** to import a file into **draft.md**.
If browser storage fails, the editor shows a warning. Download your edits before closing the page.

The tree supports keyboard navigation and built-in search.
On narrow screens, the sidebar appears before the editor and preview.
The WASM workbench keeps its existing single-document editor.

Maintainers regenerate the tree bundle with `npm ci --prefix tools/site-tree` and `npm run build --prefix tools/site-tree`.
The build also updates the Python editor archive.
Browser checks run with `python tests/editors/site_tree_checks.py --url http://127.0.0.1:8766` for a running Go workbench.
Without `--url`, the browser checks start the Python demo using the installed Markata package.
Both servers serve the tree runtime locally, without a CDN.
