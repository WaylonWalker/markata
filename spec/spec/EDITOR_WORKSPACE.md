# Editor workspace experience

Refs #1977 and #1903. The workspace includes the publishing and Changes flow in #1982.

- Settings MUST provide one searchable, responsive window with Appearance, Writing, Preview,
  Keyboard Shortcuts and Advanced categories. Existing preferences and shortcut
  overrides MUST retain their storage keys. Settings affect this browser, not site configuration.
- Commands and files MUST use fuzzy matching with deterministic ties. Exact label
  matching MUST ignore punctuation and case in human-readable labels and MUST
  outrank history. History may reorder only similarly scored matches.
- Empty queries SHOULD show recently chosen items first. Command history is local
  to this browser; file history is scoped to the workspace. Store identifiers,
  last-use time and bounded counts, never draft contents or search queries.
- History MUST survive reload, tolerate malformed/blocked storage, support clearing
  and disabling, and never mix files across workspaces.
- Startup diagnostics MUST live outside the content sidebar in an on-demand drawer.
  Measurements MUST remain available as structured JSON for agents. Missing samples
  are unknown, not zero. Retain bounded timing samples and label synthetic benchmarks.
- Files, Feeds, Tags, New Post, layouts, settings sections and diagnostics MUST be
  reachable through the command registry. Disabled views MUST not claim success.
- Focus MUST return sensibly after closing settings/pickers. Mobile settings MUST
  fit the viewport. Shortcut recording MUST work within the unified window.
- Workspace controls MUST preserve guarded publishing, read-only conflict comparison, and Changes review.
- The first local preview MUST NOT wait for workspace metadata indexing. It MAY show unresolved
  workspace links and embeds while the index loads, then MUST refresh from the current source
  when indexing completes. A metadata refresh or older render response MUST NOT replace a
  newer editor revision.
- The preview frame MUST accept the newest pending fragment once its message handler is ready.
  It MUST NOT wait for images, fonts, or other unrelated resource loads before applying an edit.
  A late frame load event MUST NOT replay an older fragment or finish the same preview twice.
- Workbench-generated frontmatter previews MUST NOT consume Markdown source-block positions.
  Cursor markers, click-to-place behavior, and synchronized scrolling MUST map authored Markdown
  blocks to their rendered blocks even when a generated link card, heading, warning, or image
  appears before them.
- Editing text outside a Mermaid diagram MUST preserve that diagram's rendered DOM and MUST NOT
  queue another Mermaid render. Editing the diagram source MUST replace and render that diagram.
- The editor preview SHOULD render Mermaid diagrams near the viewport first. Offscreen diagrams
  MAY remain as source placeholders until scrolled near view. This MUST NOT change normal site output.
- Mermaid and Chart.js fences MUST retain their plain source markup before their feature plugins
  transform them. Repeated previews MUST avoid full syntax-lexer searches for these feature fences.
- Replacing a Chart.js block MUST dispose of the old chart instance before mounting the new canvas.
- Replacing a table MUST restore its sorting controls without reloading the preview frame.
- An edit during site metadata loading MUST preserve already parsed metadata for files whose
  source and configuration have not changed. The edited file MUST be parsed again before the
  completed index is published.
- Continuous typing MUST NOT indefinitely postpone live preview. When the renderer is idle,
  an unrendered edit burst MUST start a render within 250 ms, using the current source.

The Go workbench MUST support self-contained Excalidraw scene blocks as defined
in [EXCALIDRAW.md](EXCALIDRAW.md). Editing a scene MUST update only its fenced
JSON body in one normal editor transaction. Applying a new scene MUST insert
its fence only after the author applies it. Canceling MUST leave the source
unchanged.

## Hosted editor access

The editor MUST bind to loopback by default. A non-loopback listener MUST be
rejected unless hosted authentication is explicitly configured with all of:

- one canonical HTTPS public origin with no credentials, path, query, fragment,
  or non-default port;
- one or more explicit trusted proxy CIDRs, none of which is a universal
  address range;
- one required app-specific writer scope.

Hosted authentication MUST guard every request before route dispatch, including
the editor HTML, assets, workspace APIs, render API, language-server HTTP and
WebSocket routes. It MUST fail closed unless the remote peer is in a configured
trusted CIDR, the request Host matches the public origin, exactly one non-empty
`X-Hlab-User-Id` is present, and the comma-separated `X-Hlab-Scopes` includes
the configured writer scope as an exact member. The app MUST NOT use forwarded
IP headers or accept alternate identity header names. It MUST NOT require
`X-Hlab-Auth-Method`, because Hlab browser sessions do not emit that assertion.

The hosted app MUST use the stable Hlab user ID as the request identity. It MUST
not widen permissions based on username, display name, group, role, or a
client-controlled header. Traefik MUST strip inbound `X-Hlab-*` headers before
ForwardAuth and inject only the trusted response headers after an allow result.
The editor service MUST be reachable only through that Traefik route.

Hosted workspace mutations MUST accept only an `Origin` exactly matching the
configured public origin. A present `Sec-Fetch-Site` MUST be `same-origin`.
Hosted LSP WebSocket upgrades MUST also require exactly one `Origin` matching
the configured public origin before the upgrade. Hosted responses MUST use
private, no-store caching. Local mode MUST retain its loopback and same-origin
behavior.

The editor's route scope MUST be an app-specific administrative scope, such as
`markata-editor:admin`. It MUST NOT reuse the broader Build Admin operator
scope. A viewer-only scope MUST not authorize workspace writes, uploads, or
publishing.

ForwardAuth browser session handoff MUST include the host-only
`__Host-hlab-app-session` response cookie. The editor MUST not implement a
second login protocol or trust the browser's identity headers.
