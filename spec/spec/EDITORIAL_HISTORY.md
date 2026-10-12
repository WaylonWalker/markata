---
title: "Editorial history specification"
description: "Built-site editorial history, safe comparisons, and contributor layout."
date: 2026-10-11
published: true
tags:
  - specification
  - authors
---

# Editorial history

Posts MAY supply append-only `editorial_history` frontmatter using schema version
1. Each event identifies `author`, `role`, `action` (`edited` or `reviewed`),
`recorded_at`, `recorded_by`, `base_commit`, `commit`, repository-relative `path`,
`file_url`, `diff_url`, and `summary`. Digest fields remain authoring data; the
renderer MUST NOT claim a review is current without checking its digest.

The default post template MUST expose a collapsed history disclosure only when
valid events exist. Events display newest first without mutating source history.
Each event MUST identify the contributor (configured display name when known,
otherwise author ID), role, action, recorded timestamp, summary, and edit SHA.
Reviewed events MUST NOT imply that an edit occurred. Importer attribution MUST
be displayed when `recorded_by` differs from `author`.

GitHub file links MUST match the event's full 40-hex commit SHA and normalized
repository-relative path. Diff links MUST be reconstructed for the same file,
repository, and commit. Unsafe URLs, path traversal, unsupported versions, and
malformed records MUST NOT create executable links or break the page. Invalid
revision coordinates MAY leave a metadata-only record.

Inline comparison MUST load its implementation and revision sources only after
an explicit reader action. It MUST compare Markdown bodies at `base_commit` and
`commit`, excluding YAML frontmatter. The comparison MUST show text insertions
and deletions with semantic `ins`/`del`, escaped DOM text, and a readable legend.
It MUST describe this as a historical source comparison, not reconstructed
ownership of the current rendered article. Metadata-only edits MUST report no
body changes. GitHub links and history MUST work without JavaScript.

Only public GitHub raw-file URLs derived from validated revision coordinates
MAY be fetched. No build-time network requests, embedded historical private
sources, credentials, or third-party JavaScript dependencies are required.
Failures, private repositories, oversized files, missing revisions, and CSP
restrictions MUST show a useful fallback and allow retry. Comparison work MUST
be bounded, and revision responses MUST have size and time limits.

The default byline MUST keep non-editor contributors prominent, preserve their
order, and group recognized editing/reviewing roles beneath them. The first
credited editor is displayed by name; additional editors appear as a compact
avatar stack and count in a native disclosure exposing every name and role.
Editor credit MUST remain usable with keyboard input, without hover, and without
JavaScript. Existing author-only bylines remain compatible. Both active template
trees MUST implement the same behavior.

Client initialization MUST support view-transition navigation and avoid duplicate
handlers. History and inline viewer styles use theme tokens, responsive layout,
visible focus, and text/decorations alongside color. Asset URLs use theme hashing.
