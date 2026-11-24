---
title: "Your words. On your web."
description: "A fast static site generator for people and their agents. Write Markdown and publish a website you own."
slug: ""
date: 2026-11-02
published: true
tags: [documentation, getting-started]
hero:
  title: "Your words. On your web."
  title_lines:
    - "Your words."
    - "On your web."
  eyebrow: "Markdown to website"
  subtitle: "A fast static site generator for people and their agents."
  cta_primary:
    text: "Install Markata"
    href: "#install"
  cta_secondary:
    text: "Read the quickstart"
    href: "/docs/quickstart/"
---

## One command to start {#install}

Choose an install method. Markata is a Go CLI; Python is not required to run it.

::: md-tabs {data-label="Installation methods" style="boxed"}

:::: md-tab {label="curl"}

For Linux and macOS, read [the installer source](https://github.com/WaylonWalker/markata/blob/main/install.sh) first if you want to inspect it.

```bash title="Terminal"
curl -fsSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash
markata version
```

::::

:::: md-tab {label="Bash"}

Download and inspect the installer before you run it:

```bash title="Terminal"
curl -fsSLo install-markata.sh https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh
less install-markata.sh
bash install-markata.sh
markata version
```

::::

:::: md-tab {label="mise"}

Install a release globally, or omit `-g` to pin Markata to this project:

```bash title="Terminal"
mise use -g github:WaylonWalker/markata
markata version
```

::::

:::: md-tab {label="Go"}

Install from source with Go 1.27 or newer:

```bash title="Terminal"
go install github.com/WaylonWalker/markata/cmd/markata@main
markata version
```

::::

:::: md-tab {label="GitHub"}

Choose the archive for your operating system and architecture from [GitHub Releases](https://github.com/WaylonWalker/markata/releases). Verify its checksum, extract it, and put the executable on your `PATH`.

::::

:::: md-tab {label="PyPI"}

The current `markata` package installs the Go CLI. Older pinned package versions are the original Python project. Check the selected version and platform support before installing.

```bash title="Terminal"
uv tool install markata
uvx markata --help
```

See [installation details](/docs/installation/) for pip, pipx, and the version transition.

::::

:::

The installer provides both `markata` and the compatible `markata-go` command. See [all installation options](/docs/installation/).

## Make something small

::: md-grid {data-columns="2" data-gap="lg"}

:::: md-card

### Give your agent the manual

Markata ships a site-authoring skill for Codex, Claude Code, OpenCode, Cursor, and other coding agents.

```bash title="Terminal"
markata agent install
markata agent doctor
```

See the [agent setup guide](/docs/guides/agent-skills/) for project and global installs.

::::

:::: md-card

### Write Markdown. Publish HTML.

Start a site and preview your first page:

```bash title="Terminal"
mkdir my-site && cd my-site
markata config init
markata new "Hello World"
markata serve
```

Open [localhost:8000](http://localhost:8000). Run `markata build` when you are ready to publish.

::::

:::

Continue with the [quickstart](/docs/quickstart/), learn [how to write posts](/docs/guides/authoring/), or [explore all documentation](/docs/).
