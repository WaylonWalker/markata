---
title: "Editorial history"
description: "Credit editors and inspect their recorded changes on built posts"
date: 2026-10-11
published: true
tags:
  - documentation
  - authors
  - frontmatter
---

Keep the writer prominent while giving every editor credit. The default theme
shows editors below the author byline. The first credited editor is visible;
additional editors appear in a compact avatar stack that expands to show their
names and roles. Credit order follows the post's `authors` list. Editing and
reviewing roles, including `copy editor`, `internal-link-editor`, `editor`, and
`writing-review-editing`, are grouped as editors. Coauthors keep their original
order and prominence.

```yaml
authors:
  - waylon
  - id: copy-editor-agent
    role: copy editor
  - id: linking-agent
    role: internal link editor
```

Configure contributor names and avatars as described in [Authors](../authors/).
The native [editorial authoring commands](../editorial-authoring/) can commit
selected human or agent edits and fill in the revision history automatically.

## Record a revision

Add an event to `editorial_history` after committing an editor's changes. Keep
previous events. Each event points to one file at one exact revision, even when
a commit edits several posts.

```yaml
editorial_history:
  - schema_version: 1
    event_id: "session-id:pages/blog/example.md"
    author: copy-editor-agent
    role: copy-editor
    action: edited
    recorded_by: copy-editor-agent
    recorded_at: "2026-10-11T15:00:00Z"
    base_commit: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
    commit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
    path: pages/blog/example.md
    file_url: https://github.com/OWNER/REPO/blob/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb/pages/blog/example.md
    summary: Corrected spelling and grammar.
```

Replace the example SHAs with real full 40-character commit SHAs and the URL
with the matching repository and file. Use `action: reviewed` for an unchanged
review. Set `recorded_by` to the importer when someone else records existing
work; the page displays that distinction. The native tools also record digest
fields, diff URLs, and event IDs for later automation. The viewer does not claim
that a historical review covers today's content.

## Inspect changes

The built post has a collapsed **Editorial history** section, newest event
first. Each entry identifies the contributor, role, action, timestamp, summary,
and commit. Readers can open the exact file or its diff in GitHub. These links
and the history work without JavaScript.

For edit events with valid before and after revisions, **Show changes here**
loads a small comparison viewer and fetches those two files directly from
GitHub. Added words are underlined and removed words are struck through.
The viewer compares historical Markdown source, excluding YAML frontmatter.
Markdown syntax remains visible; it does not execute historical HTML or
reconstruct ownership of the current rendered article. A change affecting only
frontmatter reports that no Markdown body changes occurred.

Revision fetching starts only after a reader clicks the comparison control.
The build does not contact GitHub or require the site's full Git history.
The viewer has no third-party script dependency, limits each file to 512 KiB,
and caps comparison work for extensive rewrites. Large rewrites may highlight
a whole changed passage rather than each word.

Push the recorded commits before publishing their links. Private repositories,
missing revisions, offline connections, and restrictive content security
policies can prevent inline comparisons. Readers can retry or open the GitHub
diff with their own account. No credentials are sent to GitHub by the inline
viewer. For a CSP, allow the site's own scripts and styles and permit
`connect-src https://raw.githubusercontent.com` for inline comparisons.

## Custom templates

The normalized template fields are `post.editorial_history_view`,
`post.byline_authors`, `post.byline_editors`, and
`post.additional_editor_count`. The original `post.author_objects` and
`post.Extra.editorial_history` remain available. The default template uses
`components/editorial_history.html`; sites overriding `post.html` should
include that component and conditionally load `css/editorial-history.css`
and `js/editorial-history.js` when normalized history exists. Asset paths
should use `theme_asset_hashed`. The component supplies the heavy viewer's
hashed asset URL for loading only on demand.

Malformed records and unsupported schema versions are skipped. Unsafe or
mismatched revision coordinates leave metadata without executable links.
GitHub diff links are rebuilt from validated file coordinates, so they always
target the recorded file in the recorded commit.
