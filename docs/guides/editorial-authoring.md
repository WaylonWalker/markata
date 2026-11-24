---
title: "Attributed editorial authoring"
description: "Use markata-go to make isolated, attributed post edits."
date: 2026-10-11
published: true
tags:
  - editorial
  - cli
---

## Editorial authoring

Use the editorial commands when a person or agent is editing a post and the
result needs a durable author, role, commit, and GitHub link.

Start a session from clean tracked Markdown files:

```bash
markata-go editorial start posts/my-post.md \
  --author editor-id --role copy-editor
```

Edit the files, then finish the session:

```bash
markata-go editorial finish --session SESSION-UUID \
  --summary "Corrected grammar and tightened headings"
```

An edited session creates two owned commits. The first contains the post edits
and contributor credit. The second appends an `editorial_history` event with
the exact file URL, diff URL, commit SHA, digest, role, and summary. Unchanged
posts are recorded as reviews with a metadata commit only.

For changes that already exist in the working tree, use the explicit commit
path:

```bash
markata-go editorial commit posts/my-post.md \
  --author linker-id --role internal-link-editor \
  --summary "Added links to related posts"
```

Selected staged changes are rejected so ownership is explicit. Other staged
and unstaged files stay in place. Add `--json` for automation or `--dry-run`
to validate inputs without changing the repository.

When a post has no `author` or `authors` field, the command uses the configured
default author when one is marked in the site's author configuration. You can
override it explicitly with `--primary-author writer-id`; the editor is never
silently promoted to the post's primary writer.

When an agent runs commands for a human who wrote the edits, keep the human's
ID in `--author` and use `--recorded-by` for the agent's ID:

```bash
markata-go editorial commit posts/my-post.md \
  --author alex --role copy-editor --recorded-by agent-model-id \
  --summary "Alex corrected grammar"
```

The content commit credits Alex; the attribution commit credits the recorder.
The built history shows that distinction. For a session, supply `--recorded-by`
when starting it; the recorder ID is retained for finish and retries.
