# Prose Assist (initial foundation)

Prose Assist is an optional, non-blocking writing-analysis feature. This
first increment contains a small, deterministic Go rule engine for editorial
phrases. It is **not** an AI-authorship detector, grammar checker, or spell
checker.

## CLI

```sh
markata prose check posts/example.md
markata prose check posts/example.md --json
markata prose check posts/example.md --ignore-rule prose/cliche
markata prose ignore prose/cliche posts/example.md
markata prose ignore-phrase prose/intensifier "very unique" posts/example.md
```

Reports use stable rule IDs and source locations. JSON results use zero-based
lines and UTF-8 byte columns. Human-readable output uses one-based positions.
The `ignore` command writes a project-level `.markata-prose.json` file in the nearest existing preferences directory (or next to the specified Markdown file). A CLI check and the LSP both discover the nearest such file by walking upward from the document. For example:

```json
{
  "ignored_rules": ["prose/cliche"],
  "ignored_phrases": [{"rule": "prose/intensifier", "phrase": "very unique"}]
}
```

Use `ignore-phrase` to accept one pattern without suppressing all the other patterns in the same rule. `ignore` remains available for disabling the entire rule. These are project preferences, not personal dictionary entries. The command does not edit source text and exits successfully when it finds
editorial suggestions. File I/O failures still return errors.

## Neovim / LSP

Set `MARKATA_PROSE_ENABLED=1` when starting `markata lsp` to opt in.
The existing LSP server then emits the same findings as *hint* severity
with `markata-prose` as source. Standard Neovim diagnostics can display them.
No extra plugin is necessary. Rule-level project ignores are shared with the CLI. With `MARKATA_PROSE_ENABLED=1`, native LSP sessions also expose standard `textDocument/codeAction` fixes: Replace (when a deterministic suggestion exists), Ignore this occurrence (until the document closes or position changes), and Always ignore this phrase (saved to project preferences), and Disable rule (for all matches). For example, `very unique` can be replaced with `unique` as a single undoable text edit. Personal dictionaries are not implemented yet.

## Editor

When the editor server starts with `MARKATA_PROSE_ENABLED=1`, the browser
workbench shows quiet amber squiggles for style hints rather than mixing them
with structural error diagnostics. Click or tap a squiggle, or press
`Alt+Enter` with the cursor on a hint, to review it in a small popover.

- **Accept** uses a versioned LSP text edit when the rule provides an explicit
  safe replacement. The edit is a single CodeMirror operation, preserving undo.
- **Ignore once** hides the current observation in the browser session.
- **Always ignore phrase** persists the selected rule and exact phrase to
  the workspace-root `.markata-prose.json` using a same-origin route. The route
  accepts known phrases only and cannot select a filesystem path.

This is an initial editor UI, **not** Harper spelling/grammar or a personal
word dictionary. Always ignore accepts one *phrase*, not the entire rule. The browser LSP remains read-only; persistence happens
through the workbench's separately authorized route.

## Next increments

1. Add personal dictionaries and richer per-occurrence ignore semantics to existing project rule/phrase preferences.
2. Add spelling and grammar via an evaluated Harper adapter.
3. Expand LSP code actions beyond the first deterministic replacement and rule-level ignores.
4. Expand the initial Editor popover to support spelling dictionaries,
   phrase-level ignores, keyboard navigation and reliable mobile behavior.
5. Evaluate Vale style packs and document-level repetition metrics, with
   false-positive tests and browser/runtime parity checks.

Editorial checks must never block builds or publishing by default. Prose
analysis must also preserve Markdown and code source locations. The initial
lightweight masks intentionally prefer omitting possible findings in protected
regions to reporting inside syntax; richer Markdown parsing will follow.
