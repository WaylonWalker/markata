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
```

Reports use stable rule IDs and source locations. JSON results use zero-based
lines and UTF-8 byte columns. Human-readable output uses one-based positions.
The `ignore` command writes a project-level `.markata-prose.json` file in the nearest existing preferences directory (or next to the specified Markdown file). A CLI check and the LSP both discover the nearest such file by walking upward from the document. For example:

```json
{
  "ignored_rules": ["prose/cliche"]
}
```

This is rule-level suppression, not a personal dictionary, and it intentionally applies to all occurrences of that rule within the project. The command does not edit source text and exits successfully when it finds
editorial suggestions. File I/O failures still return errors.

## Neovim / LSP

Set `MARKATA_PROSE_ENABLED=1` when starting `markata lsp` to opt in.\nThe existing LSP server then emits the same findings as *hint* severity
with `markata-prose` as source. Standard Neovim diagnostics can display them.
No extra plugin is necessary. Rule-level project ignores are shared with the CLI. With `MARKATA_PROSE_ENABLED=1`, native LSP sessions also expose standard `textDocument/codeAction` fixes: Replace (when a deterministic suggestion exists), Ignore this occurrence (until the document closes or position changes), and Always ignore this rule (saved to project preferences). For example, `very unique` can be replaced with `unique` as a single undoable text edit. Personal dictionaries are not implemented yet.

## Editor

The browser workbench already consumes Markata LSP diagnostics. Prose hints are
now available in that diagnostic stream; Google Docs-like squiggles and
Accept / Ignore / Always ignore actions are **not implemented yet**. Those
require editor UI work, stable edit actions, and a shared persisted vocabulary.

## Next increments

1. Extend project rule ignores with user dictionaries and per-occurrence ignore semantics.
2. Add spelling and grammar via an evaluated Harper adapter.
3. Expand LSP code actions beyond the first deterministic replacement and rule-level ignores.
4. Add unobtrusive inline decorations and accessible suggestion popovers in
   Editor, with a single undoable edit per correction.
5. Evaluate Vale style packs and document-level repetition metrics, with
   false-positive tests and browser/runtime parity checks.

Editorial checks must never block builds or publishing by default. Prose
analysis must also preserve Markdown and code source locations. The initial
lightweight masks intentionally prefer omitting possible findings in protected
regions to reporting inside syntax; richer Markdown parsing will follow.
