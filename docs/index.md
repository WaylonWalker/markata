---
title: "Build with Markata"
description: "Write Markdown, build a fast static site, and publish it with Markata."
slug: ""
date: 2026-10-18
published: true
tags: [documentation, getting-started]
---


Markata turns a folder of posts into a site you own. Start with the short path below; add feeds, themes, and plugins when you need them.

::: road steps
- [Install Markata](/docs/installation/)
- [Create a site](/docs/quickstart/)
- [Write your first post](/docs/guides/authoring/)
- [Publish your site](/docs/guides/deployment/)

:::

## Keep learning

Explore the [authoring guide](/docs/guides/authoring/) to write posts, then adjust [configuration](/docs/guides/configuration/) and [themes](/docs/guides/themes/) as your site grows. For more, see [feeds and RSS](/docs/guides/feeds/) or the [CLI reference](/docs/reference/cli/).

For agent-assisted authoring, install the bundled [Markata site skill](/docs/guides/agent-skills/). Keep drafts private, review generated changes, and publish when they're ready.

## Releases

- [Release archive](/docs/releases/) — themes, highlights, and upgrade notes for every published version
- [v0.15.0 release notes](/docs/releases/v0-15-0/) — rendering beyond the static-site build
- [Release and deployment runbook](/docs/guides/releasing/) — tag, artifacts, GitOps rollout, build verification and rollback

Markata is still **0.x**. The older `markata-go` command remains compatible. New Go sites should continue to use `markata-go.toml` and `[markata-go]` until the configuration migration is complete.

**PyPI version break:** Current `markata` releases install the Go CLI; older pinned releases remain the original Python implementation. See [installation](/docs/installation/) for details.

If something fails, run `markata triage` or see the [triage guide](/docs/guides/triage/).

Use Markata in an application with the [Python and FastAPI renderer](/docs/guides/python-renderer/), [Go workbench](/docs/guides/go-workbench/), or [browser workbench](/docs/guides/wasm-workbench/).
