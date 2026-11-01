---
title: "Markata documentation"
description: "Install Markata, teach your agent the site skill, and publish your first Markdown page."
slug: ""
date: 2026-10-18
published: true
tags: [documentation, getting-started]
---

# Markata documentation

**A fast static site generator for people and their agents.**

Let's be honest. We're all running agents these days. They manage everything.

Write Markdown. Give your agent the right context. Publish a site you own.

## Start here

1. [Install Markata](/docs/installation/) using curl, a shell, mise, GitHub Releases, or Go.
2. [Quickstart](/docs/quickstart/) — create your first site and preview it locally.
3. [Set up your agent](/docs/guides/agent-skills/) — install the bundled `markata-go-site` skill.
4. [Write your first post](/docs/guides/authoring/) — Markdown, frontmatter, drafts, images, and links.
5. [Publish](/docs/guides/deployment/) — deploy the generated static files.

## Explore as you grow

- [Configuration](/docs/guides/configuration/)
- [Themes](/docs/guides/themes/)
- [Feeds and RSS](/docs/guides/feeds/)
- [Templates](/docs/guides/templates/)
- [Plugins](/docs/reference/plugins/)
- [CLI reference](/docs/reference/cli/)
- [Naming migration](/docs/guides/name-migration/)

## Releases

- [v0.14.0 release notes](/docs/releases/v0.14.0/) — what's new, compatibility and upgrade guidance (release preparation)
- [Release and deployment runbook](/docs/guides/releasing/) — tag, artifacts, GitOps rollout, build verification and rollback

Markata is still **0.x**. The older `markata-go` command remains compatible. New Go sites should continue to use `markata-go.toml` and `[markata-go]` until the configuration migration is complete.

**PyPI version break:** Current `markata` releases install the Go CLI; older pinned releases remain the original Python implementation. See [installation](/docs/installation/) for details.

If something fails, run `markata triage` or see the [triage guide](/docs/guides/triage/).

Use Markata in a web application: [Python and FastAPI rendering](/docs/guides/python-renderer/),
the [Go workbench](/docs/guides/go-workbench/), or the [browser WASM workbench](/docs/guides/wasm-workbench/).
See [renderer integration benchmarks](/docs/guides/renderer-integration-benchmarks/) for measured timings.
