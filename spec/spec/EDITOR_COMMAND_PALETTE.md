# Renderer Editor Command Palette

The Python, Go, and WASM workbenches MUST offer the same Markdown snippet palette.
A Commands button and Ctrl+Shift+P (Cmd+Shift+P on macOS) MUST open a modal dialog.
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

Commands MUST include note, tip, warning, aside, collapsible note, highlight,
keyboard keys, definition list, Mermaid, Chart.js, table, task list, code block,
bold, italic, link, and image. Only standalone-renderable syntax belongs in this
initial catalog. Site-dependent post embeds and wikilinks are excluded because
these editors do not load a site content index.

One canonical JavaScript and CSS implementation MUST ship to every workbench.
Packaging checks MUST prevent the Python bundled copy from drifting.
Tests MUST exercise real CodeMirror selections, history, modal interactions,
search, rendering, and local asset availability.
