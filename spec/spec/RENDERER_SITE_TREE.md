# Renderer example site navigation

The Python and Go workbenches MUST include a bundled example Markata site and a
sidebar rendered by Pierre `@pierre/trees`. The sidebar MUST provide Files and
Feeds views. Feed folders MUST contain references to their member source files.
Opening a reference MUST open the same draft as opening its source file.

Selecting a Markdown file MUST update the editor and preview. Switching files
MUST preserve edits in memory and attempt to save drafts per source path in
browser storage. Storage failures MUST leave navigation usable and show a warning.
The existing single-document draft MUST remain available as `draft.md`.
Configuration files MUST be read-only and appear as code in the preview.
Markdown previews MUST omit the leading YAML frontmatter block.

The example files MUST exist in the repository. A reproducible asset build MUST
bundle their contents and feed membership for both runtimes. Editing this demo
MUST NOT modify repository files or require server filesystem access.
The sidebar MUST support keyboard navigation, have an accessible label, and adapt
to narrow viewports. Runtime assets MUST be served locally without CDN requests.
WASM remains a single-document editor in this change.
