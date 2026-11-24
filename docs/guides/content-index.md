---
title: "Content Index"
description: "Use Markata's stable metadata index for search and tooling."
date: 2026-08-15
published: true
tags:
  - documentation
  - content-index
---


Markata can write one compact JSON file with the metadata it resolved during a
build. It does not include article bodies.

```toml
[markata-go.content_index]
enabled = true
output = "content-index.json"
```

The file is written below `output_dir` unless `output` is absolute. Its
identity is `schema = "markata.content-index"` and its current generation is
`schema_version = 2`. V2 emits `scope = "public-metadata"`. This means the file
can contain safe metadata records for private posts; it never contains their
article bodies or encryption keys.

## Reading the file

Read `schema` and `schema_version` before decoding. Ignore fields you do not
know. If the generation is newer than the reader supports, stop with a clear
upgrade message. Do not guess how a future generation works.

Use each document's repository-relative `path` to find its source file. The
`feeds` array already contains resolved Markata feed membership. Consumers
should not implement Markata's feed filter language again.

Documents may also include optional resolved frontmatter and author metadata:
`image`, `video`, `avatar`, `bio`, `thumbnail`, `cover`, `og_image`, `author`,
`authors`, `category`, and `categories`. Media values are source references;
Markata does not copy or download those assets into the index. Author order is
preserved. `cover` uses `cover` first and `cover_image` as a compatibility
fallback. `og_image` uses `og_image` first and `social_image` as a fallback.

`source.commit` identifies the Git `HEAD` observed for the build.
`source.dirty` is false when the working tree was observed clean at the Content
Index source-state boundaries. It is true when Git reports tracked, staged,
deleted, untracked, or ignored Markdown source files. Compare both values with
the source repository revision before using metadata as a bootstrap cache.
Both fields are absent when Git cannot provide source state.
Markata captures this state before content discovery and rechecks it before
writing. A changing source state fails an enabled index write instead of
publishing a mismatched identity.

The pair `source.commit` plus `dirty: false` is build provenance. It does not
prove that every document came byte-for-byte from a blob in that Git revision.
A source file can change and be restored between the two observations. Use Git
commit/blob identity, or another exact source identity mechanism, when exact
per-file synchronization is required.

The index includes non-draft, non-skipped posts, including private posts. A
`published: false` direct page may still appear. A feed named `draft` is feed
membership only and does not mean that the document has `draft: true`.

Private records contain only safe metadata. Titles and descriptions are present
only when they were explicitly provided in frontmatter; content-derived values,
article bodies, rendered HTML, key names, private media, and resolved author
biographies are omitted. Tags, dates, slugs, aliases, author identifiers,
categories, and an explicit avatar may remain. The `feeds` list for a private
record contains only feeds configured with `include_private = true` (or
`private = true`).

Markata retains readers for every released generation and normalizes them into
one current internal model. The index is not a source archive, rendered-body
cache, access-control system, or guarantee that source files remain available.
If you disable the output, use a clean output directory or remove the old
artifact before deployment.

The v1 and v2 schemas and immutable interoperability fixtures are in
`pkg/contentindex/v1_schema.json`, `pkg/contentindex/v2_schema.json`, and
`pkg/contentindex/fixtures/`. Set `schema_version = 1` when a consumer still
requires the public-only v1 artifact.
