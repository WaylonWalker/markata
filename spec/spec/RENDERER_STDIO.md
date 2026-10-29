# Persistent renderer stdio protocol (first implementation)

Tracking: [#1858](https://github.com/WaylonWalker/markata/issues/1858). Baseline renderer contract: [#1855](https://github.com/WaylonWalker/markata/issues/1855).

Use one long-running process to render multiple Markdown documents without
launching a new Go executable for each document:

```sh
markata render --stdio
```

The process reads one JSON object per line from stdin and writes one JSON object
per line to stdout. Embedded newlines in source must be escaped as `\\n` JSON
characters. **Stdout is reserved for protocol output**; stderr is available for
CLI errors. No site build, site config loading, or network request is performed
by `RenderStandalone` itself. Authored HTML and URLs still have the usual
security considerations.

## Request examples

```json
{"id":"doc-1","op":"render","source":"# Hello\n\nWorld","options":{"inline_assets":false}}
```

```json
{"id":"caps","op":"capabilities"}
```

```json
{"id":"bye","op":"shutdown"}
```

Every nonempty input line receives one JSON response. `id` is a caller-supplied,
nonempty string of at most 128 UTF-8 bytes, echoed back for valid envelopes.
A successful render response contains `{"id":"doc-1","ok":true,"result":{...}}`;
a capability response contains `{"id":"caps","ok":true,"capabilities":{...}}`.
Errors return `ok:false` and `error:{code,message}`. Malformed JSON or
oversized framing might not have a recoverable request ID, so `id` is omitted.

Supported operations:

| Operation | Behavior |
| --- | --- |
| `render` | Render one Markdown source via `RenderStandalone` |
| `capabilities` | Return the shared capability descriptor without rendering |
| `shutdown` | Acknowledge, then exit normally without reading more input |

Supported per-request `options`: `unsafe` (default false),
`inline_assets` (default false), `asset_base` (default
`/markata-assets`). All options are isolated to their request.

## Failure behavior and safety

- A malformed JSON line or unknown top-level/option field receives
  `invalid_json`; the process continues.
- An unknown operation returns `unsupported_operation`, with the original ID.
- A missing or excessively long ID returns `invalid_id`.
- A `source` greater than 200,000 bytes returns `source_too_large`.
- An input frame greater than 1,048,576 bytes returns `frame_too_large`,
  and the frame is drained to its newline before accepting another request.
- Renderer errors return `render_failed`; later requests remain usable.
- Errors writing to stdout or reading from stdin terminate with a CLI error.

Responses are **sequential, not concurrent**; a slow diagram render blocks
later requests. There is no per-request cancellation yet. No shared site config
is loaded and no session-wide options or cache protocol is exposed. Feature
assets may be repeated in each result; optimized asset deduplication and a
Python persistent session wrapper are separate follow-ups.

The shared render result stays at `schema_version:1`. The stdio envelope
does not introduce a new Markdown renderer or change the legacy
`markata render --json` command.

## Verification

```sh
go test ./pkg/renderprotocol -count=1
go test ./cmd/markata-go/cmd -run '^TestRenderCommand_' -count=1
```

Tests cover multiple render requests, capability discovery, malformed frames,
unknown options, size limits, shutdown, and one-shot output parity.
