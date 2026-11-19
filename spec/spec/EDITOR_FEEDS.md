# Feed Studio — editor contract

Tracking: [#2032](https://github.com/WaylonWalker/markata/issues/2032)

## Purpose
A configured feed is an editable site object. The editor MUST expose feed
settings and effective membership without asking authors to edit TOML, YAML or
JSON, and MUST distinguish the live membership model from the site's actual
published HTML.

## User workflow
1. Select **Feeds** in the workspace sidebar. The Feed Studio opens in the
   main editor area; the existing Files and Changes views remain available.
2. Pick a configured feed to edit title, description, filter, sort/reverse,
   presentation, and HTML template selection.
3. A filter change triggers a debounced, cancellable Go filter evaluation.
   Invalid syntax MUST be shown inline and MUST disable saving.
4. **Matching posts** shows the membership result and an **Edit post** action
   for each visible entry. This action opens the existing post editor, not a
   new editing implementation. **Back to feed** restores the feed workspace.
5. **Built site** loads the latest generated HTML from this site's output
   folder, including the feed's actual custom template and site assets.
   It is explicitly labeled a *last-built* preview. It MUST NOT claim to
   include unsaved draft changes or an unrun build.
6. Saving feed settings writes the originating site configuration and makes
   it eligible for the existing Git review/publish workflow.

## Data and API
- `GET /workspace/feed-studio`: full configured feed list, per-feed
  membership counts, config version, and whether each entry is directly
  editable in the root site configuration.
- `GET /workspace/feed-studio?slug=S&filter=EXPR&sort=FIELD&reverse=true`:
  same plus membership preview, filter errors, total count and first 60
  visible posts. Request ordering is authoritative: clients must cancel or
  ignore stale responses.
- `POST /workspace/feed-studio`: JSON with slug, version and the editable
  configuration fields. Input is validated. The server uses optimistic
  concurrency, rejecting stale versions with HTTP 409.
- `GET /workspace/feed-snapshot/... `: sandboxed resource access to the
  selected workspace's *generated* output, not arbitrary source files.

All workspace endpoints retain the editor's loopback-host and peer restriction.
Mutations additionally require same-origin headers. The output route validates
path ancestry, regular-file status, and resolved symlinks.

## Compatibility and constraints
- The existing Go `filter.Parse` / `Match` logic determines preview
  membership, with the same metadata-to-post projection as Feeds navigation.
- Public membership previews exclude private and draft posts even if a
  malformed authored filter would otherwise match them.
- The author's template path is treated as configuration data, never as an
  arbitrary file read by the save endpoint.
- Config fields not edited by Feed Studio MUST be preserved. Multiple-tab
  edits MUST result in a conflict rather than silently overwriting.
- Feeds defined by recursively included config files inside the trusted
  workspace MUST be edited in their actual source file. Each save uses the
  *owning file* SHA-256 version; the root file and unrelated includes MUST
  remain byte-for-byte unchanged. Unsafe or outside-workspace targets MUST
  remain read-only.
- The built-site route never generates HTML on its own; no build or cached
  preview should block basic filter editing or editor startup.

## Acceptance checks
- Valid filters update membership; invalid filters surface errors.
- Post edit and return preserve the selected feed.
- All writes honor expected-version compare-and-swap and preserve unrelated
  config values.
- Last-built previews show real HTML and are marked stale with respect to
  unsaved filters/templates.
- No nonlocal fetch or traversal is possible through the snapshot endpoint.
- Keyboard, narrow-screen and empty-feed states remain usable.

The initial implementation is an iterative foundation. Follow-up work may
add exact rebuilt preview scheduling, and
template-specific post card rendering in the live membership view.
