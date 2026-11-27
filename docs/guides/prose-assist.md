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
```

Reports use stable rule IDs and source locations. JSON results use zero-based
lines and UTF-8 byte columns. Human-readable output uses one-based positions.
The command does not edit source text and exits successfully when it finds
editorial suggestions. File I/O failures still return errors.

## Neovim / LSP

The existing `markata lsp` server emits the same findings as *hint* severity
with `markata-prose` as source. Standard Neovim diagnostics can display them.
No extra plugin is necessary. This first increment does not yet implement
prose-specific code actions or persistent ignore preferences for the LSP.

## Editor

The browser workbench already consumes Markata LSP diagnostics. Prose hints are
now available in that diagnostic stream; Google Docs-like squiggles and
Accept / Ignore / Always ignore actions are **not implemented yet**. Those
require editor UI work, stable edit actions, and a shared persisted vocabulary.

## Next increments

1. Define shared configuration and a persisted user/project ignore dictionary.
2. Add spelling and grammar via an evaluated Harper adapter.
3. Add LSP code actions for safe suggestions and ignore operations.
4. Add unobtrusive inline decorations and accessible suggestion popovers in
   Editor, with a single undoable edit per correction.
5. Evaluate Vale style packs and document-level repetition metrics, with
   false-positive tests and browser/runtime parity checks.

Editorial checks must never block builds or publishing by default. Prose
analysis must also preserve Markdown and code source locations. The initial
lightweight masks intentionally prefer omitting possible findings in protected
regions to reporting inside syntax; richer Markdown parsing will follow.
