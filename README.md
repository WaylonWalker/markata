# Markata

**Write Markdown. Publish more than pages.**

Markata builds static sites from Markdown. Define feeds to filter, sort, and paginate content into HTML, RSS, Atom, or JSON; use templates and plugins to shape the result.

[Documentation](docs/quickstart.md) · [Installation](docs/installation.md) · [Releases](https://github.com/WaylonWalker/markata/releases)

## Install

```sh
# curl (Linux / macOS)
curl -fsSL https://markata.dev/install.sh | bash

# mise
mise use -g github:WaylonWalker/markata

# uv
uv tool install markata

# uvx (run without installing)
uvx markata --help

# pip / pipx
pip install markata
pipx install markata

# Go (1.27+; current source)
go install github.com/WaylonWalker/markata/cmd/markata@main
```

Choose **one** method above, or download a binary from [GitHub Releases](https://github.com/WaylonWalker/markata/releases). Check your installation with `markata version`.

PyPI now distributes the Go implementation after a version break from [Python Markata](https://github.com/WaylonWalker/markata-py). See [installation details](docs/installation.md).

## Quick start

```sh
mkdir my-site && cd my-site
markata config init
markata new "Hello World"
markata serve
```

Build the output when you're ready:

```sh
markata build
```

## Let your agent set up the site

The installer gives you the Markata engine. To choose your own templates, plugins, and deployment environment, hand your coding agent the release-matched first-run instructions:

```sh
markata --skill
```

This prints a standalone Markdown guide without modifying the site. For a full project skill with topic references, use `markata agent install` and verify with `markata agent doctor`. Agents should run a production build and ask before deploying.

## Set up your coding agent

```sh
markata agent install -g
```

See the [documentation](docs/) for agent configuration and other commands.

## Learn more

- [Quick start](docs/quickstart.md)
- [Installation and troubleshooting](docs/installation.md)
- [Configuration](docs/)
- [Feeds, templates, and plugins](docs/)
- [Naming migration](docs/guides/name-migration.md)

Markata is the Go implementation (formerly `markata-go`). The compatibility command and existing `markata-go.toml` configuration remain supported. See the [migration guide](docs/guides/name-migration.md).
