---
title: "Markata documentation"
description: "Install Markata, teach your agent the site skill, and publish your first Markdown page."
slug: ""
date: 2026-10-18
published: true
tags: [documentation, getting-started]
---


**A fast static site generator for people and their agents.**

Let's be honest. We're all running agents these days. They manage everything.

Write Markdown. Give your agent the right context. Publish a site you own.

## Install Markata

Choose an installation method. Markata is a single Go-based CLI; you don't need Python to run it.

::: md-tabs {#homepage-install data-label="Installation methods" style="boxed"}

:::: md-tab {label="curl"}

```bash
curl -fsSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash
markata version
```

Linux / macOS. [Review the installer source](https://github.com/WaylonWalker/markata/blob/main/install.sh) before running it.

::::

:::: md-tab {label="Bash"}

```bash
curl -fsSLo install-markata.sh https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh
less install-markata.sh
bash install-markata.sh
markata version
```

Download and inspect the installer before executing it.

::::

:::: md-tab {label="mise"}

```bash
mise use -g github:WaylonWalker/markata
markata version
```

Install from GitHub Releases. Older release archives may use the `markata-go` executable name.

::::

:::: md-tab {label="Go"}

```bash
go install github.com/WaylonWalker/markata/cmd/markata@main
markata version
```

Build the current source with Go 1.27 or newer.

::::

:::: md-tab {label="GitHub"}

Download the appropriate binary from [GitHub Releases](https://github.com/WaylonWalker/markata/releases), verify the checksum, and add it to your `PATH`.

::::

:::: md-tab {label="PyPI"}

```bash
uv tool install markata
# or run without installing:
uvx markata --help
# or use pip:
pip install markata
```

Current `markata` PyPI releases install the Go CLI; older pinned versions are the original Python implementation.

::::

:::

See [all installation options](/docs/installation/) for platform details and troubleshooting.

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

- [Release archive](/docs/releases/) — themes, highlights, and upgrade notes for every published version
- [v0.15.0 release notes](/docs/releases/v0-15-0/) — rendering beyond the static-site build
- [Release and deployment runbook](/docs/guides/releasing/) — tag, artifacts, GitOps rollout, build verification and rollback

Markata is still **0.x**. The older `markata-go` command remains compatible. New Go sites should continue to use `markata-go.toml` and `[markata-go]` until the configuration migration is complete.

**PyPI version break:** Current `markata` releases install the Go CLI; older pinned releases remain the original Python implementation. See [installation](/docs/installation/) for details.

If something fails, run `markata triage` or see the [triage guide](/docs/guides/triage/).

Use Markata in a web application: [Python and FastAPI rendering](/docs/guides/python-renderer/),
the [Go workbench](/docs/guides/go-workbench/), or the [browser WASM workbench](/docs/guides/wasm-workbench/).
See [renderer integration benchmarks](/docs/guides/renderer-integration-benchmarks/) for measured timings.
