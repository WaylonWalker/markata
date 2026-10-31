# Renderer Editor Command Palette

The Python, Go, and WASM workbenches MUST offer the same Markdown snippet palette.
A Commands button, Ctrl+K (Cmd+K on macOS), and Ctrl+Shift+P (Cmd+Shift+P on macOS) MUST open a modal dialog.
The shortcut MUST work with every supported editor keymap, outside IME composition.
The palette MUST use local assets and MUST work in the offline WASM workbench.

Search MUST match all query words against command labels, descriptions, syntax,
and aliases without case sensitivity. Results MUST show the label, description,
and syntax hint. Empty results MUST explain how to recover.
Arrow keys MUST change the active result. Enter MUST insert that result.
Escape, a Close button, and a backdrop click MUST dismiss without changing text.
The dialog MUST contain keyboard focus and expose its active result to assistive
technology. Dismissal MUST restore editor focus and the original selection.

Insertion MUST replace the selection saved when the palette opened. Commands
with a text placeholder MUST preserve selected text there. Block commands MUST
add blank lines around their output when surrounding text requires separation.
The inserted placeholder MUST be selected for replacement. Insertion MUST be one
undoable editor operation and MUST trigger existing draft saving and rendering.
The initial version operates on the primary selection.

The catalog MUST cover every documented custom Markdown family: all admonition
types, collapsible variants and split layouts; all eight road styles and
editorial modifiers; article layout/card/icon/feature containers; status, stat,
delta, meter, progress, sparkline, checklist, and source containers; Mermaid and Chart.js;
all ten Web Awesome Markdown shortcuts; Lucide icons; and inline extensions.
All bundled Lucide icon names MUST be searchable by name and insertable offline.
A named icon command MUST replace selected text with the chosen icon; it MUST NOT
interpret selected prose as an icon name.
A generic local icon-pack shortcode MUST also be offered for site authors.

Site-dependent syntax (wikilinks, post embeds, external embed cards, TOC,
glossary frontmatter, custom icon packs, and Web Awesome components) MUST be
insertable. Descriptions MUST state the site/plugin/runtime requirement rather
than promise a functioning standalone preview. Web Awesome shortcuts MUST use
the exact forms in WEBAWESOME.md. A raw Web Awesome HTML template MUST support
components without a Markdown shorthand. This palette does not change renderer
capabilities or fetch third-party runtimes. Interactive checklist enhancement
MUST accept the classless list items emitted by the standalone Goldmark renderer
as well as site markup with task-list classes. Ordinary task lists MUST stay disabled.

Search MUST remain responsive with the full icon catalog. At most 60 results
are displayed at once. The count MUST show the total number of matches and
instruct authors to refine the search when results are truncated.

The repository MUST include a self-contained example document that demonstrates
the custom elements that work in the standalone preview. Site-dependent examples
MUST be clearly identified. Tests MUST cover catalog completeness and syntax
conformance in addition to interaction behavior.

One canonical JavaScript and CSS implementation MUST ship to every workbench.
Packaging checks MUST prevent the Python bundled copy from drifting.
Tests MUST exercise real CodeMirror selections, history, modal interactions,
search, rendering, and local asset availability.

Each workbench MUST support `?example`: load the shared showcase Markdown into
the editor, render its preview, and open the command palette. Example assets MUST
be included in the Python bundle and offline WASM cache.
