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

## Workspace appearance and rendering

When a workspace config exists, both servers MUST use it for standalone Markdown
highlighting and initial palette, light/dark mode, aesthetic, fontpack, text size,
and palette variable overrides. Workspace palette and aesthetic catalogs and an
explicit fontpack catalog MUST resolve relative to the workspace. Local custom CSS
MUST be served with appearance assets. Browser appearance preferences may override
the site's defaults. Without a config, existing standalone defaults MUST remain.
The `render --workspace` protocol MUST make these settings available to Python
without a second configuration parser. Standalone preview does not run a full
site build or apply site page templates.

Browser draft records MUST distinguish intentional empty drafts from legacy empty
placeholders. Existing nonempty Markdown drafts MUST survive format migration.

Editor and generated appearance assets MUST revalidate and use content-versioned
page references, so restarting a server at the same origin loads current code and
site defaults.

## Progressive startup

`GET /workspace` MUST return discovered paths without loading document metadata.
`GET /workspace?metadata=true` MUST provide titles and tags separately. The browser
MUST open its selected source before requesting metadata, and MUST construct the
Tags tree only when selected. Metadata failure MUST NOT block editing or preview.
CLI `workspace --paths` and Python `workspace_files(metadata=False)` expose the
same fast discovery path. Existing metadata API defaults remain compatible.

Scratch and workspace sources MUST have separate storage identities. A real
`draft.md` MUST open its source; scratch navigation MUST use a collision-free label
(such as `Unsaved draft.md`) and retain old browser scratch content during migration.
