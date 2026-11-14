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
- An edit during site metadata loading MUST preserve already parsed metadata for files whose
  source and configuration have not changed. The edited file MUST be parsed again before the
  completed index is published.
- Continuous typing MUST NOT indefinitely postpone live preview. When the renderer is idle,
  an unrendered edit burst MUST start a render within 500 ms, using the current source.
