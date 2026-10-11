---
title: "Editorial authoring specification"
description: "Deterministic attributed commits for selected human and agent edits."
date: 2026-10-11
published: true
tags:
  - specification
  - editorial
---

# Editorial authoring

`markata-go editorial` provides deterministic Git authoring for Markdown posts.
It is designed for humans and agents that need their work attributed without
including unrelated working tree changes.

## Commands

`editorial start posts/post.md --author id --role copy-editor` requires selected
tracked Markdown files to be clean and stores a session under
`.git/markata-go-editorial/<uuid>.json`.
The caller edits the selected files, then runs `editorial finish --session UUID
--summary "..."`. A session produces an editorial content commit followed by a
metadata commit containing `editorial_history` records. Unchanged files receive
`action: reviewed` and use the baseline commit.

`editorial commit posts/post.md --author id --role internal-link-editor
--summary "..."` is the short path for explicitly owned existing changes.
Selected staged files are rejected. An isolated Git index means unrelated
staged and unstaged files remain outside the commits.

## Data contract

Records use schema version 1 and contain `event_id`, `author`, `role`, `action`,
`recorded_by`, `recorded_at`, `base_commit`, `commit`, `path`, `content_digest`,
`digest_version`, `file_url`, `diff_url`, and `summary`. The Go implementation
uses the explicitly named `markata-go-editorial-v1` digest because it does not
claim byte-for-byte compatibility with another implementation.

The command appends an editor to `authors` while preserving existing entries.
Posts without a primary author require a configured default author or an
explicit `--primary-author` value, so an editor is never silently promoted to
the writer.
Frontmatter is parsed as YAML and the Markdown body is retained. The
`Editorial-Author` Git trailer carries the stable contributor ID.

`--recorded-by` MAY identify a distinct person or agent recording work on the
contributor's behalf. Store this value in the session and event. The content
commit's trailer identifies the contributor; the metadata commit's trailer
identifies the recorder. The recorder defaults to the contributor.

The `markata-go-editorial-v1` digest is SHA-256 of Go JSON serialization of the
parsed YAML frontmatter mapping with `author`, `authors`, and
`editorial_history` removed, followed by one newline and the exact Markdown
body bytes after the closing frontmatter delimiter. JSON map keys are sorted
by Go's encoder. YAML timestamps use the parser's typed representation and Go
JSON encoding. Frontmatter formatting is excluded; body whitespace is included.
This digest version is distinct from the Python `editorial-v1` normalization.
