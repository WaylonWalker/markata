# Editor workspace navigation

The Python and Go workbenches MUST use Pierre `@pierre/trees` for a workspace
file tree and a tag grouping view. Tag entries MUST open the same draft as
their canonical file entries.

Selecting a Markdown file MUST update the editor and preview. Switching files
MUST preserve edits in memory and attempt to save drafts per source path in
browser storage. Storage failures MUST leave navigation usable and show a warning.
The existing single-document draft MUST remain available as `draft.md`.
Markdown previews MUST omit the leading YAML frontmatter block.

The sidebar MUST support keyboard navigation, have an accessible label, and adapt
to narrow viewports. Runtime assets MUST be served locally without CDN requests.
WASM remains a single-document editor in this change.

## Server workspaces

Go and Python servers MUST expose the selected local workspace through the same
Go discovery and loading API. Discovery uses configured Markata glob patterns,
gitignore behavior and the load plugin. Only discovered, root-contained regular
files may be read. Symlink escapes and arbitrary path requests MUST be rejected.
The browser MUST fetch the file list and load source on selection. Source files
remain unchanged; drafts and downloads retain the example's existing behavior.
Python MUST delegate discovery and reads to the bundled Go executable.
