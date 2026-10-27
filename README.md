# Markata

**Write Markdown. Publish more than pages.**

Markata builds static sites from Markdown. Define feeds to filter, sort, and paginate content into HTML, RSS, Atom, or JSON; use templates and plugins to shape the result.

[Documentation](docs/quickstart.md) · [Installation](docs/installation.md) · [Releases](https://github.com/WaylonWalker/markata/releases)

## Install

Choose one of the existing installation methods:

**curl (Linux / macOS)**

```sh
curl -fsSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash
```

**mise**

```sh
mise use -g github:WaylonWalker/markata
```

**Go (1.26+)**

```sh
go install github.com/WaylonWalker/markata/cmd/markata@main
```

**Manual download:** Pick a binary for your platform from [GitHub Releases](https://github.com/WaylonWalker/markata/releases).

```sh
markata version
```

**uv / uvx (PyPI)**

```sh
uv tool install markata
# Or run without a persistent installation:
uvx markata --help
```

**pip / pipx (PyPI)**

```sh
pip install markata
# Or install as an isolated CLI:
pipx install markata
```

The PyPI package now distributes the Go implementation following a version break from the [original Python Markata](https://github.com/WaylonWalker/markata-py). See [installation details](docs/installation.md).

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
