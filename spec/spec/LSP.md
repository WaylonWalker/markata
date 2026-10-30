# Language Server Protocol Specification

## Browser editor examples

The Go and FastAPI workbench examples expose `/lsp` as a same-origin WebSocket.
Each text message contains one JSON-RPC object without Content-Length headers.
The Go example uses the existing LSP handlers in process. FastAPI translates
WebSocket messages to the existing stdio protocol of a server-owned executable.
Each connection owns an independent LSP session and index.

The server selects the workspace at startup. Browser initialize requests cannot
change it. `/lsp/config` supplies a root URI and a synthetic `.markata/browser-draft.md` URI.
Documents remain in memory. Neither example reads, saves, or switches source
files through the browser. File trees and file operations are separate work.
Workspace file notifications from the browser MUST NOT cause filesystem reads.

Both editors synchronize full unsaved buffers and request completion from the
same LSP. They support wiki-link slugs and aliases, Markdown document paths,
container names and modifiers from the command catalog, and documented container
attributes and values. Mermaid completions include diagram types, flowchart
directions, a flowchart starter, and class style properties. Existing frontmatter, admonition, and mention completions remain available.
Completion edits use UTF-16 positions,
preserve closing delimiters, and cannot apply to a changed buffer or cursor.
Ctrl+Space requests completion explicitly. Broken-link diagnostics appear in
the editor. Connection failures leave editing and preview available.

The WebSocket endpoints reject foreign or missing browser origins, bound message
sizes, cap completion responses at 200 items, and release connections or subprocesses
on disconnect. They are trusted local examples, not a multi-user remote workspace
service. Existing stdio clients
retain their workspace selection and Content-Length framing by default.

## Overview

`markata-go lsp` starts the markata-go Language Server Protocol (LSP) server on
standard input and standard output. It provides Markdown wikilink completion,
diagnostics, hover information, and go-to-definition.

The LSP server is editor-neutral. Editor clients MUST start it with:

```text
markata-go lsp
```

The client MUST communicate using LSP over stdio and register the server for
Markdown files. The workspace root SHOULD be the site directory that contains
the active markata-go configuration.

## Setup Guidance

The CLI MUST provide editor-configuration snippets without mutating editor
configuration files:

```text
markata-go lsp setup --editor <editor>
```

Supported editor identifiers are `generic`, `neovim`, `helix`, `emacs`, `zed`,
and `vscode`. The output MUST state the command, Markdown filetype, and
workspace-root requirements. It MUST identify when an editor requires a
third-party LSP client or an extension rather than implying native support.

The command MUST reject an unsupported editor identifier with a usage error and
list the supported identifiers. It MUST write snippets to stdout so agents and
users can copy or redirect them safely.

When `--editor` is omitted, the command MUST detect supported editors on
`PATH` and print setup snippets for each installed editor. If no supported
editor is installed, it MUST print the generic integration contract.
When multiple snippets are printed, they MUST have clear editor headings. On
an interactive terminal, snippets SHOULD use syntax highlighting; redirected
output MUST remain uncolored copyable text.

## Diagnostics

The CLI MUST provide a read-only diagnostic command:

```text
markata-go lsp doctor
```

`doctor` MUST check that the current binary can serve LSP requests and validate
an explicit or auto-discovered markata-go configuration. A missing
configuration is informational because the LSP can still start; an invalid
explicitly selected or discovered configuration is an error. LSP mention
indexing reads TOML and `.yaml` configuration only; the command MUST warn when
a valid `.yml` or `.json` site configuration is selected.

The command MUST detect supported editor executables on `PATH` and, by default,
load supported editor configuration in headless or batch mode. It MUST report
configured, unconfigured, and inconclusive results without modifying
configuration. `--no-verify-editor` MUST opt out of loading editor
configuration and report only installed editor detection. Loading editor
configuration can execute user plugin startup code.

The command MUST write its report to stdout and return exit code 0 when all
required checks pass. It MUST return exit code 1 when a required check fails.
It MUST not start a long-running LSP server or modify editor configuration.
Human-oriented reports SHOULD use semantic status color when stdout is a
terminal and MUST honor the CLI color-disable controls. When verification is
skipped and a detected editor supports it, the report MUST print the default
doctor follow-up on stderr.

When a site palette is available, the report MUST derive its title, pass,
information, warning, and failure colors from that palette. It MUST retain the
same status labels so color is an enhancement rather than the only signal.

## Version Support

Editor-specific snippets are examples for current, documented client
configuration surfaces. They are not a promise to automatically configure all
editor versions. The generic integration contract is the stable compatibility
surface; agents SHOULD inspect a repository's existing editor configuration and
adapt the snippet to the installed client.
