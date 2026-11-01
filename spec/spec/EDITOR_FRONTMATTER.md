# Editor frontmatter properties contract

Stacked on PR #1906 (publish review), itself based on #1905.

## Interaction
- Closed by default, header folded in CodeMirror; Properties button, palette,
  and Ctrl/Cmd+Alt+Shift+M expose the form.
- Edit YAML reveals the original source header; Hide YAML folds it back.
- Common metadata is quickly discoverable; all other valid root-level YAML
  properties are available through an add/search picker or custom key, including
  plugin-specific fields. Nested YAML values are editable as typed maps/lists.
- Missing keys are not materialized as defaults. Explicit empty strings are
  distinct from removal.
- Draft and published values remain independent and conflicts are surfaced.

## Data integrity
- Source is canonical. Form POSTs receive a snapshot of unsaved Markdown and
  return a transformed source string; never write to disk in these handlers.
- Client applies patches only if source still equals the submitted snapshot.
- Only the edited root field is serialized via yaml.v3; all other keys,
  comments, body bytes and line endings are preserved.
- Invalid delimiter/YAML syntax, duplicate keys and invalid property keys
  refuse form operations, prompting hand-edit mode.
- Existing CodeMirror on-change path owns autosave, preview and undo/redo.
- Local same-origin checks and 200KB source limit remain in force.

## Verification
- Unit tests for preserving unrelated YAML/body and CRLF, add/remove,
  malformed YAML and duplicate-key rejection.
- HTTP tests cover local same-origin authorization and inspection/update.
- Browser smoke: editor loads with collapsed frontmatter; expand form, add
  custom nested values, edit YAML, switch site files, preview, undo, autosave,
  and open publish review without losing values.
