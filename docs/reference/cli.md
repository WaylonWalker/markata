---
title: "CLI Reference"
description: "Complete reference for all markata-go commands, flags, and options"
date: 2024-01-15
published: true
slug: /docs/reference/cli/
tags:
  - documentation
  - reference
  - cli
---

# CLI Reference

markata-go provides a command-line interface for building, serving, and managing your static site.

## Overview

The markata-go CLI follows a subcommand pattern similar to tools like `git` and `docker`. Each command has its own set of flags, with some global flags available across all commands.

```bash
markata-go [global flags] <command> [command flags] [arguments]
```

The common short forms are `b` for `build`, `s` and `serv` for `serve`, `a` for
`admin`, and `ls` for `list`. For example, `markata-go s --bind 0.0.0.0` starts the development
server on every network interface. `--bind` and `--host` select the same
listening address; if both are set, their values must agree.
The semantic names `dev` and `preview` suggest `serve` but do not execute it.

If a command, flag, or config key is misspelled, markata-go suggests close
matches and shows the exact help command to run. A clear alias executes; a
semantic synonym or typo only suggests. Suggestions are advisory: the
CLI never executes or writes a guessed name. For example:

```text
Error: unknown flag: --hst
Did you mean:
  --host
Next: Run 'markata-go serve --help' to see valid flags.
```

Unknown command names exit with status `2`. The direct file shortcut accepts
Markdown paths such as `markata-go post.md`; other unknown names produce an
error. Errors use color when `stderr` is a terminal, and plain text when it is
redirected or `--no-color` is set.

For a one-file preview, pass a Markdown file directly. The file is rendered
with the default theme and published at the site root (`<output>/index.html`),
along with the theme's CSS, JS, and fonts. Site-wide output is skipped: feeds,
the archive, tag/feed listings, sitemaps, `.well-known` files, the image
library, 404 page, and the Pagefind search index (search UI is hidden). The
build cache is not used, so a single-file build never touches the full-site
cache:

```bash
markata-go pages/sample.md
# equivalent:
markata-go build pages/sample.md
# live preview of just that file, rebuilt on every change:
markata-go serve pages/sample.md
```

## Agent-assisted troubleshooting: `triage`

`markata triage` collects a small, privacy-conscious diagnostic summary and passes a troubleshooting prompt to an interactive coding agent, clipboard, or stdout. The default is a fuzzy picker of **detected installed** coding agents plus clipboard and stdout. When either terminal stream is redirected or `--no-input` is used, it prints **only the Markdown prompt** to stdout.

```sh
markata triage
markata triage --to opencode
markata triage --agent claude
markata triage --issue "My posts aren't appearing" --to clipboard
markata triage --to stdout | codex exec -
markata triage --to stdout > triage-prompt.md
```

| Flag | Description | Default |
| --- | --- | --- |
| `--to` | `auto`, `stdout`, `clipboard`, `opencode`, `claude`, `codex`, `gemini`, or `cursor` | `auto` |
| `--agent` | Convenience alias for an installed agent name; cannot be combined with `--to` | none |
| `--issue` | Optional brief user-provided problem description (up to 4000 bytes) | none |

Task-specific subcommands and flags:

```sh
markata triage lint --scope site
markata triage lint --scope post --post posts/story.md
markata triage post --post posts/story.md --action debug --feature mermaid --issue "Parse error"
markata triage post --post posts/story.md --action enhance --feature road --goal "Add project milestones"
```

| Intake flag | Description |
| --- | --- |
| `--category` | `general`, `lint`, or `post`; usually implied by a subcommand |
| `--scope` | `site` or `post` for `triage lint` |
| `--post` | Relative Markdown post path; metadata only (512-byte limit) |
| `--action` | `debug`, `enhance`, or `review` for `triage post` |
| `--feature` | `chart`, `road`, `mermaid`, `markdown`, or `other` |
| `--goal` | Desired enhancement (4000-byte limit) |

Interactive missing values are asked **in a deterministic order** before the
destination picker. Category, scope, action, and feature use a searchable fuzzy
picker. Missing `--post` values show a fuzzy list of local Markdown paths from
common content folders, with manual entry as a fallback. Only filenames are
scanned; Markata does not parse the post contents. Escape cancels safely.
Noninteractive use never reads stdin for answers: missing
questions become an ordered section of the emitted prompt, with clean
prompt-only stdout preserved. Supplied answers are not asked again.

Only supported agent executables found on PATH are offered. An explicit agent launch needs a terminal so the agent can ask follow-up questions. For pipelines use stdout and the selected agent's own noninteractive mode. Clipboard uses local platform helpers and fails clearly if unavailable. No agent is launched automatically in a noninteractive session.

The prompt contains the **running CLI** version, not necessarily the version where the bug occurred, plus OS, architecture, runtime, executable/install hints, project path and **config filenames only**. It never reads config values, logs, full environment values or secrets. Review local paths before sharing a prompt. The agent should first clarify the user's issue, confirm the actual affected version/deployment and ask before applying changes.

For examples and privacy details, see [Triage a Markata problem](../guides/triage.md).

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to configuration file | Auto-discovered |
| `--site-dir` | | Site directory for all site operations | `MARKATA_GO_SITE_DIR` or current directory |
| `--merge-config` | `-m` | Additional config file(s) to merge (can be used multiple times) | None |
| `--output` | `-o` | Output directory (overrides config) | From config (`output` built-in default) |
| `--quiet` | `-q` | Suppress non-essential progress and status output | `false` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--color` | | Force ANSI color output | `false` |
| `--no-color` | | Disable ANSI color on all streams | `false` |
| `--log-format` | | Log formatting: `auto`, `plain`, or `rich` | `auto` |
| `--no-input` | | Disable prompts and interactive UI | `false` |

### Config File Discovery

When `--config` is not specified, markata-go searches for configuration files in the following order:

1. `markata-go.toml`
2. `markata-go.yaml`
3. `markata-go.yml`
4. `markata-go.json`
5. `~/.config/markata-go/config.toml`

If none is found, commands that load the full site configuration continue with
built-in defaults plus `MARKATA_GO_*` environment overrides. Commands that edit
or read a value directly from a config file, such as `config get` and
`config set`, require a config file.

See [[configuration-guide|Configuration]] for details on configuration options.

### Site Directory

Use `--site-dir` to run a command against a site from any directory. Set
`MARKATA_GO_SITE_DIR` to use the same site in a shell alias or wrapper. The flag
takes precedence over the environment variable, which takes precedence over the
current directory.

```bash
alias blog='markata-go --site-dir "$HOME/sites/blog"'
alias work-notes='markata-go --site-dir "$HOME/notes/work"'

blog search "release process" --format path
work-notes new "Weekly review" --template note --no-input
```

The selected directory must exist. Configuration discovery, relative
`--config` and `--merge-config` paths, content templates, linting, and caches
all resolve from that directory. `list` and `search` print absolute Markdown
source paths when a site directory is explicit, so their results can be edited
directly.

### CLI UX Conventions

markata-go follows a few shared CLI rules:

- primary command results are written to `stdout`
- warnings, progress, prompts, and errors are written to `stderr`
- interactive terminal output uses color by default when the target stream is a TTY
- operational logs use centralized formatting; `--log-format auto` chooses rich logs for TTYs and plain logs otherwise
- `--color` forces ANSI color output, while `--no-color` disables it everywhere
- rich logs can use structured metadata such as lifecycle phase, and builds/serve runs can tint log colors from the configured site palette
- color is disabled when output is not a terminal, when `NO_COLOR` is set, when
  `TERM=dumb`, or when `--no-color` is passed
- interactive commands such as `new` and `init` honor `--no-input`
- commands that expose structured data should prefer explicit machine-readable
  modes such as `--json` or `--format json`

Docs: https://github.com/WaylonWalker/markata-go/tree/main/docs

Issues: https://github.com/WaylonWalker/markata-go/issues

---

## Commands

### version

Print the version, commit, build date, Go runtime, and platform for this CLI
binary. Release and container builds inject the version, commit, and build date
at compile time. Use `--short` when a script needs only the version string.

```bash
markata-go version
markata-go version --short
```

### lsp

Start the Markdown language server, check its local prerequisites, or print
editor setup guidance.

```bash
markata-go lsp
markata-go lsp doctor
markata-go lsp doctor --no-verify-editor
markata-go lsp setup
markata-go lsp setup --editor neovim
```

`markata-go lsp` communicates using LSP over standard input and output. It
provides wikilink completion, broken-wikilink diagnostics, hover information,
and go-to-definition for Markdown files.

| Subcommand | Description |
| --- | --- |
| `doctor` | Read-only check of the LSP binary, active project configuration, and installed editors |
| `setup --editor <name>` | Print a copy-paste editor configuration snippet without editing files |

`setup` supports `generic`, `neovim`, `helix`, `emacs`, `zed`, and `vscode`.
Use `generic` for another client. Editor recipes are guidance for documented
client configuration surfaces, not a guarantee that every editor version can be
configured automatically. See [[editor-integration|Editor Integration]] for
full instructions.

By default, `doctor` performs headless validation for installed Neovim, Helix,
and Emacs configurations. This can load user startup/plugin code. Pass
`--no-verify-editor` to only detect installed editors. VS Code, Cursor, and
Zed detection is advisory because their LSP clients do not expose stable
headless verification APIs.

When verification is skipped and a deeper check is available, `doctor` prints
the default follow-up command to stderr.

On an interactive terminal, the doctor report derives its status colors from
the configured site palette while retaining text labels for each status.

Interactive help and human-facing setup headings also use the configured site
palette. If a site palette cannot be resolved, Markata falls back to its default
CLI theme; structured command output remains uncolored for scripting.

When `--editor` is omitted, `setup` prints snippets for each supported editor
installed on `PATH`; it prints the generic integration contract if none are
found. Multiple snippets have editor headings and are syntax highlighted on an
interactive terminal; redirected output remains plain text.

The LSP server reads TOML and `.yaml` configuration for mention indexing. If
your valid site configuration is `.yml` or `.json`, `doctor` warns that those
LSP-specific mentions are unavailable.

### agent

Install bundled agent integrations for markata-go site repositories.

#### Usage

```bash
markata-go agent install [site-path] [flags]
```

#### Subcommands

##### list-agents

List the supported agent identifiers and their project/global skill directories.

##### install

Install the bundled `markata-go-site` skill into a repository.

| Flag | Description | Default |
|------|-------------|---------|
| `--agent` | Install for a specific agent such as `opencode`, `claude-code`, or `cursor` | auto-detected, else `universal` |
| `-g`, `--global` | Install into the selected agent's user-level skill directory | `false` |
| `--name` | Installed skill directory name | `markata-go-site` |
| `--force` | Overwrite bundled skill files if they already exist | `false` |
| `--dry-run` | Show what would be installed without writing files | `false` |

##### update

Update the installed `markata-go-site` skill in place.

This is the user-friendly equivalent of reinstalling with `--force`.

| Flag | Description | Default |
|------|-------------|---------|
| `--agent` | Update a specific agent install such as `opencode` or `claude-code` | auto-detected, else `universal` |
| `-g`, `--global` | Update the selected agent's user-level skill directory | `false` |
| `--name` | Installed skill directory name | `markata-go-site` |
| `--dry-run` | Show what would be updated without writing files | `false` |

##### doctor

Check the installed skill for drift against the versions bundled in the current binary. Reports files that are modified, missing, or newly added since the skill was installed.

| Flag | Description | Default |
|------|-------------|---------|
| `--agent` | Check a specific agent install such as `opencode` or `claude-code` | auto-detected, else `universal` |
| `-g`, `--global` | Check the selected agent's user-level skill directory | `false` |
| `--name` | Installed skill directory name | `markata-go-site` |

Exit codes:

| Code | Meaning |
|------|---------|
| `0` | Skill is up to date |
| `1` | Drift detected |
| `2` | Error (skill not installed, manifest unreadable) |

##### remove

Remove the installed `markata-go-site` skill directory.

`uninstall` is an alias for `remove`.

| Flag | Description | Default |
|------|-------------|---------|
| `--agent` | Remove a specific agent install such as `opencode` or `claude-code` | auto-detected, else `universal` |
| `-g`, `--global` | Remove from the selected agent's user-level skill directory | `false` |
| `--name` | Installed skill directory name | `markata-go-site` |

#### Examples

```bash
# Install into the detected agent's project layout
markata-go agent install

# List supported agent ids and paths
markata-go agent list-agents

# Install into Claude Code's .claude layout
markata-go agent install --agent claude-code

# Install into OpenCode's global skill directory
markata-go agent install --agent opencode -g

# Preview files without writing
markata-go agent install --dry-run

# Install into a different repository
markata-go agent install ../my-site

# Replace an existing installed skill
markata-go agent install --force

# Update an installed skill in place
markata-go agent update

# Preview an update without writing
markata-go agent update --dry-run

# Check for drift after upgrading the binary
markata-go agent doctor

# Check a Claude Code layout
markata-go agent doctor --agent claude-code

# Check a global OpenCode install
markata-go agent doctor --agent opencode -g

# Check a different repository
markata-go agent doctor ../my-site

# Remove an installed skill
markata-go agent remove

# Alias for remove
markata-go agent uninstall

# Remove a Claude Code install
markata-go agent remove --agent claude-code

# Remove a global OpenCode install
markata-go agent remove --agent opencode -g
```

#### Installed Layouts

- detected agent or `universal`: project install into that agent's project skill path
- `--agent claude-code`: `.claude/skills/markata-go-site/`
- `--agent opencode`: `.agents/skills/markata-go-site/`
- `--agent opencode -g`: `~/.config/opencode/skills/markata-go-site/`

The installed skill is split into `SKILL.md` plus focused topic files under `topics/`, reference material under `reference/`, starter files under `examples/`, and regression prompts under `evals/` so agents can read only the sections relevant to the current task while maintainers still have a starter eval set for bundled-skill reviews. A `.manifest.json` file is written alongside the skill for drift detection via `agent doctor`.

When `-g` is used, `--agent` is required and `site-path` is not accepted.

The `agent` command group is intentionally generic so future subcommands can add export or MCP-oriented integrations without changing the bundled skill format.

### assets

Manage external CDN assets used by self-hosted/vendor asset mode.

#### Usage

```bash
markata-go assets <subcommand>
```

#### Subcommands

##### download

Download all registered external assets into the configured asset cache, defaulting to `.markata/assets-cache`. This includes the Web Awesome npm archive so builder images can pre-seed the cache for offline site builds.

```bash
markata-go assets download
```

##### list

List registered assets and whether each one is available in the asset cache.

```bash
markata-go assets list
```

##### clean

Remove the asset cache so assets are fetched again on the next download/build.

```bash
markata-go assets clean
```

### build

Build the static site by processing all content files through the plugin lifecycle.

#### Usage

```bash
markata-go build [flags]
```

#### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--clean` | | Remove output directory and build cache before building | `false` |
| `--clean-all` | | Also remove external plugin caches before building | `false` |
| `--dry-run` | | Show what would be built without writing files | `false` |
| `--fast` | | Skip minification, CSS purge, Tailwind rebuilds, and Pagefind indexing | `false` |
| `--dag` | | Use the experimental serial DAG executor | `false` |
| `--benchmark-json` | | Write benchmark details as JSON; use `-` for stdout | `""` |
| `--benchmark-summary-json` | | Write compact benchmark JSON without content entries; use `-` for stdout | `""` |
| `--benchmark-detailed` | | Print per-stage benchmark resource summaries | `false` |
| `--verbose` | `-v` | Enable verbose logging | `false` |
| `--output` | `-o` | Override output directory | from config |

#### Examples

```bash
# Standard build
markata-go build

# Clean build (replaces output after the new build succeeds)
markata-go build --clean

# Also clear external plugin caches
markata-go build --clean-all

# Preview what would be built
markata-go build --dry-run

# Fast dev build
markata-go build --fast

# Experimental serial task graph
markata-go build --dag

# Build with verbose output
markata-go build -v

# Build to a custom output directory
markata-go build -o dist

# Build with a specific config file
markata-go build -c production.toml

# Combine flags
markata-go build --clean -v -o dist
```

#### Content diagnostics

Every build prints a `Content:` summary. It counts discovered files, content
candidates, loaded sources, valid frontmatter, posts, eligible posts, rendered
posts, emitted outputs, exclusions, warnings, and errors.

Content diagnostics identify the source path, line, stable reason code, and a
short message. For example, a frontmatter opening of `----` reports
`frontmatter.suspicious_delimiter`; it does not silently discard the file.
Malformed author content is isolated to that source, but source stat, read, and
other filesystem errors still fail the build.

Use `-v` to show the final disposition and reason codes for every excluded or
shadow source:

```text
Content dispositions:
  posts/draft.md [excluded]: content.draft
  posts/review.md [shadow]: content.published_false
```

The JSON written by `--benchmark-json` includes the same sanitized content
snapshot under `content`. It does not include raw frontmatter, environment
values, secrets, or build logs.

A successful full build writes the sanitized
content diagnostics artifact to `<output_dir>/.markata/diagnostics.json`. It includes a versioned schema,
summary counts, per-source dispositions, stable reason codes, and feed-level
selection details. Inspect it with a JSON tool such as:

```bash
jq . output/.markata/diagnostics.json
```

Use the configured output directory instead of `output` when it differs. The
artifact is not written by `build --dry-run`, `build --fast`, fast or incremental
`serve`, or an incomplete or failed build. A normal `serve` build can write the
artifact after a successful full lifecycle. See [[content-diagnostics|Content Diagnostics]]
for the field and privacy contract.

### buildlab

Compare baseline and candidate builds in isolated workspaces.

#### Usage

```bash
markata-go buildlab run --fixture PATH [flags]
```

Build Lab runs ordinary `build` commands. It does not use a task graph or an
experimental scheduler. Each checkpoint compares clean baseline output with
clean candidate output, candidate incremental output when available, and a
second clean candidate build when determinism checks are enabled.

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--fixture` | Site fixture directory | required |
| `--baseline` | Baseline markata-go binary | current executable |
| `--candidate` | Candidate markata-go binary | current executable |
| `--build-config` | Config path relative to the fixture | `markata-go.toml` |
| `--scenario` | Versioned scenario JSON file | clean/no-op scenario |
| `--result` | Write the JSON result to this path | none |
| `--seed` | Scenario seed | `1` |
| `--timeout` | Per-build timeout | `10m` |
| `--gomaxprocs` | GOMAXPROCS for child builds | `1` |
| `--check-determinism` | Compare two clean candidate builds | `true` |
| `--volatile` | Comma-separated volatile output paths | `.well-known/time` |
| `--fast` | Add `--fast` to both child builds | `false` |
| `--tool-version` | External tool metadata, such as `pagefind=1.5.2` | none |
| `--env` | Allowed entries: `PATH`, `MARKATA_GO_ENCRYPTION_ENABLED`, or `MARKATA_GO_OFFLINE` | none |

Build Lab copies the fixture, redirects cache and temporary paths, disables
fixture `.env` loading, and pins children to UTC, `C.UTF-8`, and
`SOURCE_DATE_EPOCH=0`. It rejects paths and environment entries that escape the
workspace. Setup failures emit a versioned harness result to stdout and a
concise diagnostic to stderr. See [[build-lab|Build Lab]] for scenario and
output-class details.

### builder-admin

Run the long-lived builder admin HTTP service for Kubernetes and authoring workflows.

#### Usage

```bash
markata-go builder-admin [flags]
```

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--host` | Host to bind to | `127.0.0.1` |
| `--port` | Port to listen on | `8080` |
| `--source-dir` | Mounted source directory to watch and build from | `.` |
| `--release-dir` | Mounted release root that contains `releases/` and `current` | `public` |
| `--watch` | Enable recursive file watching and queued rebuilds | `true` |
| `--watch-debounce` | Debounce window for coalescing file changes | `2s` |
| `--fast` | Use `markata-go build --fast` for queued builds | `false` |
| `--mermaid-mode` | Override `[markata-go.mermaid].mode` for queued builds | `""` |
| `--cache-mount` | Optional dedicated cache mount used for `.markata` and `.markata-cache` symlinks | `""` |
| `--history-dir` | Directory for persisted admin state and logs | `<release-dir>/.builder-admin` |
| `--releases-keep` | Number of rendered releases to keep on disk | `25` |
| `--refresh-task` | Repeatable task spec in the form `name|every|enqueue|arg1|arg2...` | none |
| `--trusted-proxy-cidr` | Repeatable CIDR allowed to supply trusted identity headers | none |
| `--auth-user-id-header` | Trusted header containing the durable operator ID | `X-Hlab-User-Id` |
| `--auth-username-header` | Trusted header containing the operator username | `X-Hlab-Username` |
| `--auth-display-name-header` | Trusted header containing the operator display name | `X-Hlab-Display-Name` |
| `--auth-email-header` | Trusted header containing the operator email | `X-Hlab-Email` |
| `--auth-groups-header` | Trusted header containing operator groups | `X-Hlab-Groups` |
| `--auth-roles-header` | Trusted header containing operator roles | `X-Hlab-Roles` |
| `--auth-scopes-header` | Trusted header containing operator scopes | `X-Hlab-Scopes` |
| `--auth-method-header` | Trusted ForwardAuth assertion for bearer-authenticated API requests | `X-Hlab-Auth-Method` |

Header mappings may also be set in `[markata-go.builder_admin.auth.headers]` or with
`MARKATA_GO_BUILDER_ADMIN_AUTH_HEADERS_*` environment variables. Explicit CLI flags override
environment and file configuration. See [Builder Admin](../guides/deployment/builder-admin.md)
for a complete ForwardAuth and Authentik example.

#### Examples

```bash
# Run locally against mounted source and site paths
markata-go builder-admin \
  --config /data/source/markata-go.toml \
  --source-dir /data/source \
  --release-dir /data/site \
  --cache-mount /data/cache \
  --fast \
  --watch

# Add a scheduled reader refresh that enqueues a build when complete
markata-go builder-admin \
  --refresh-task 'reader-update|30m|true|markata-go|--config|/data/source/markata-go.toml|reader|update'
```

#### Exit Codes

| Code | Description |
|------|-------------|
| `0` | Build completed successfully |
| `1` | Build failed (configuration error, plugin error, etc.) |
| `2` | No content files found |

#### Build Process

The build command executes the full 9-stage lifecycle:

1. **Configure** - Load and merge configuration
2. **Validate** - Validate configuration settings
3. **Glob** - Discover content files
4. **Load** - Parse markdown files and frontmatter
5. **Transform** - Pre-render processing (wikilinks, jinja-md, etc.)
6. **Render** - Convert markdown to HTML
7. **Collect** - Build feeds and collections
8. **Write** - Output files to disk
9. **Cleanup** - Release resources

When `--verbose` is enabled, build output includes per-stage timing to highlight slow stages.
`--fast` keeps the same HTML output path but skips minification, CSS purge, Tailwind rebuilds,
and Pagefind indexing for a tighter dev loop.

Successful builds also print a compact benchmark summary with:

- estimated wall-time spent on CPU work, network wait, disk read wait, disk write wait, and idle time
- the slowest lifecycle hotspots so you can spot expensive plugins quickly
- the slowest outbound HTTP requests so network-bound plugins stop hiding in aggregate wait time

For deeper analysis:

- `markata-go build --benchmark-json=benchmark.json` writes machine-readable benchmark data to a file
- `markata-go build --benchmark-json=-` writes only the benchmark JSON to stdout
- `markata-go build --benchmark-summary-json=/tmp/benchmark.json` writes timing and aggregate content statistics without the large per-content entry matrix
- `markata-go build -v --benchmark-detailed` adds per-stage resource summaries to the build footer

When requests are present, the footer prints the 10 slowest network waits with stage, plugin, method, sanitized URL, duration, and either HTTP status or the request error.

The resource profile is approximate. It is intended for local hotspot hunting, not precise system profiling.

The experimental `--dag` flag is available on `build`, `serve`, and
`builder-admin`. It retains serial plugin ordering and is not a parallel-build
switch. `MARKATA_GO_DAG=true` also opts in; an explicit `--dag=false` overrides
the environment, including queued Builder Admin child builds. Invalid
environment boolean values report an error unless an explicit flag overrides
them. `build --dry-run` still uses the legacy partial lifecycle and does not
claim a DAG build ran. Benchmark JSON records the executor that actually ran.
See [[dag-builds|Experimental DAG builds]] for comparison guidance.

---

### reader

Refresh the cached external feed data used by the `/reader/` page without running a build.

#### Usage

```bash
markata-go reader update
```

##### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--concurrency` | Override feed refresh concurrency for this run | config value, else `5` |

#### Subcommands

##### update

Fetch all configured blogroll feeds and rewrite the on-disk reader cache so the next build can reuse fresh reader data.

This command does not render pages or write site output. It only updates the external feed cache used by the blogroll and reader features.

#### Examples

```bash
# Refresh reader cache using the auto-discovered config
markata-go reader update

# Refresh reader cache with higher parallelism for this run
markata-go reader update --concurrency 12

# Refresh reader cache using a specific config file
markata-go reader update -c production.toml
```

#### Notes

- requires `[markata-go.blogroll] enabled = true`
- requires at least one `[[markata-go.blogroll.feeds]]` entry
- bypasses normal cache age checks so it can fetch fresh remote feed data immediately
- reuses stale cached feed data on fetch failures when available

---

### serve

Start a development server with live reload support.

#### Usage

```bash
markata-go serve [flags] [markdown-file]
```

Pass a Markdown file to serve only that page at `/`, using the same single-file mode as `build`. Every rebuild stays in single-file mode.

#### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--port` | `-p` | Port to listen on | `8000` |
| `--host` | | Host address to bind to | `localhost` |
| `--bind` | | Address to listen on (alias for `--host`) | `localhost` |
| `--admin` | | Explicitly enable the local Control Center on loopback | Loopback default |
| `--no-admin` | | Disable the local Control Center | `false` |
| `--no-tui` | | Use line-oriented output | `false` |
| `--watch` | | Enable file watching and auto-rebuild | `true` |
| `--no-watch` | | Disable file watching (legacy, overrides --watch) | `false` |
| `--fast` | | Skip minification/CSS purge, Tailwind rebuilds, Pagefind indexing, and disable blogroll/mentions | `false` |
| `--incremental` | | Reuse unchanged posts and dependencies while keeping normal output processing | `false` |
| `--verbose` | `-v` | Enable verbose logging | `false` |

#### Examples

```bash
# Serve on default port (localhost:8000)
markata-go serve

# Serve on a custom port
markata-go serve -p 3000
markata-go serve --port 3000

# Bind to all network interfaces (accessible from other devices)
markata-go serve --host 0.0.0.0

# Disable the local Control Center on a loopback listener
markata-go serve --no-admin

# Explicitly enable file watching (default behavior)
markata-go serve --watch

# Serve without file watching
markata-go serve --watch=false
markata-go serve --no-watch

# Serve with production-style incremental rebuilds
markata-go serve --incremental

# Serve with verbose logging
markata-go serve -v

# Combine options
markata-go serve -p 8080 --host 0.0.0.0 -v
```

#### Features

- **Automatic rebuild**: File changes trigger a rebuild automatically
- **Live reload**: Connected browsers refresh when content changes
- **Static file serving**: Serves the output directory
- **MIME type detection**: Correct content types for all file types
- **Immediate serve**: Server starts before the initial build completes
- **Local Control Center**: `http://localhost:8000/_markata/` is available by default on loopback. A non-loopback bind leaves it disabled. `--admin` is rejected there, and `--admin --no-admin` is invalid. The TUI shows `⚙ Admin`; press `a` to open it in a browser.
- **Build status banner**: Shows build progress and errors during serve
- **Early 404**: Minimal 404 is served until the generated 404.html exists
- **Incremental mode**: `serve --incremental` rebuilds only changed posts and dependents while retaining minification, CSS purging, Tailwind, redirects, and Pagefind processing
- **Fast mode**: `serve --fast` rebuilds only changed posts and dependents, but skips production output steps for iteration speed
- **Bake theme**: the theme picker gains a **Bake** button that writes your current palette, style, and font (and text size or color mode when you picked them) into `[markata-go.theme]`. It edits the config file (including `include`d files) that already holds the theme, keeping comments, then rebuilds. It never writes your global `~/.config/markata-go/config.toml`. Only `serve` offers it. See [Themes Guide](/docs/guides/themes/#baking-a-look-into-your-config).
- **Settings sidebar**: a gear button at the bottom left (or `Alt+,`) opens a sidebar with every config setting, its current value, and its description. Changes preview live on the site right away without touching your files; click **Bake** to review a diff and write them into the config file that owns each one, or **Reset** to discard them. **Reset to default** removes a setting from your config. Settings with fixed values are dropdowns, and invalid values are rejected. Previews build into `.markata/serve-preview/` so your main build cache is kept. Only `serve` shows it. See [Configuration Guide](/docs/guides/configuration/#edit-settings-while-serving).
- **Local only**: the settings and Bake endpoints accept requests only from the machine running `serve` (a loopback client address). With `--host 0.0.0.0`, other devices can browse the site but cannot change your config.

#### Development Workflow

```bash
# Terminal 1: Start the server
markata-go serve -v

# The server will:
# - Start serving at http://localhost:8000 immediately
# - Build your site in the background
# - Watch for file changes
# - Rebuild and reload on changes
```

---

### tui

Interactive terminal UI for browsing and managing your markata site. Inspired by [k9s](https://k9scli.io/).

#### Usage

```bash
markata-go tui [flags]
```

#### Description

The TUI provides a keyboard-driven interface for exploring your site's posts, tags, and feeds. It features a k9s-style table layout with filtering, sorting, and quick actions.

#### Views

The TUI has three main views, accessible via number keys:

| Key | View | Description |
|-----|------|-------------|
| `1` | Posts | Browse all posts in a table |
| `2` | Tags | View and filter by tags |
| `3` | Feeds | View configured feeds |

#### Global Keybindings

These keys work in all views:

| Key | Action |
|-----|--------|
| `1`, `2`, `3` | Switch between views |
| `/` | Enter filter mode |
| `Esc` | Clear filter / exit current mode |
| `?` | Show help |
| `q`, `Ctrl+C` | Quit |

#### Posts View

The posts view displays all posts in a table with columns for title, date, word count, tags, and path.

##### Posts Keybindings

| Key | Action |
|-----|--------|
| `↑`/`↓`, `j`/`k` | Navigate up/down |
| `Enter` | View post details |
| `e` | Edit post in $EDITOR |
| `s` | Open sort menu |
| `/` | Filter posts |

##### Post Detail View

Press `Enter` on any post to see its full details including:
- Title, path, and date
- Published status and tags
- Word count and description
- Content preview

Press `Esc` to return to the post list, or `e` to edit the post.

#### Sorting

Press `s` in the posts view to open the sort menu:

| Sort Option | Description |
|-------------|-------------|
| Date | Sort by post date (default) |
| Title | Sort alphabetically by title |
| Word Count | Sort by word count |
| Path | Sort by file path |

Use `a` for ascending order, `d` for descending order.

#### Filtering

Press `/` to enter filter mode. Type a filter expression and press `Enter` to apply.

##### Filter Expression Examples

```
published==true           # Only published posts
tags contains "go"        # Posts tagged with "go"  
date > "2024-01-01"       # Posts after a specific date
title startswith "How"    # Titles starting with "How"
```

See [Filter Expressions](/docs/guides/filters) for the complete filter syntax.

#### Feeds View

The feeds view displays all configured feeds with:
- Feed name
- Number of posts in the feed
- Filter expression used
- Output path

#### Environment Variables

| Variable | Description |
|----------|-------------|
| `EDITOR` | Primary editor for `e` key (falls back to `vim`) |
| `VISUAL` | Secondary editor if `EDITOR` is not set |

#### Examples

```bash
# Launch the TUI
markata-go tui

# The TUI will display your posts in a table format
# Use j/k to navigate, Enter to view details, e to edit
```

---

### list

List posts, tags, or feeds in table, JSON, CSV, or path-only formats.

#### Usage

```bash
markata-go list posts [flags]
markata-go list tags [flags]
markata-go list feeds [flags]
markata-go list feeds posts <feed>
```

#### Shared Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format: `table`, `json`, `csv`, `path` | `table` |
| `--sort` | Sort field (varies by subcommand) | posts: `date`, tags: `count`, feeds: `name` |
| `--order` | Sort order: `asc` or `desc` | posts/tags: `desc`, feeds: `asc` |

#### Posts Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--filter` | Filter expression for posts | none |
| `--feed` | Limit posts to a feed by name | none |

#### Examples

```bash
# List posts as a table
markata-go list posts

# Path-only output for piping
markata-go list posts --format path

# List tags as JSON
markata-go list tags --format json

# List feeds as CSV
markata-go list feeds --format csv

# List posts for a feed
markata-go list posts --feed blog

# List posts for a feed
markata-go list feeds posts blog

# List post paths for a feed
markata-go list feeds posts blog --format path
```

#### Cache

`list` and `tui` use a persistent cache at `.markata/cache/list.json`. Delete this file to force a full refresh.

---

### search

Full-text search across post content, titles, descriptions, and tags. Uses a bleve full-text index for BM25-ranked results with optional fuzzy matching.

#### Usage

```bash
markata-go search <query> [flags]
```

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format: `table`, `json`, `csv`, `path` | `table` |
| `--sort` | Sort field: `score`, `date`, `title`, `words`, `path`, `reading_time`, `tags` | `score` |
| `--order` | Sort order: `asc` or `desc` | `desc` |
| `--filter` | Additional filter expression to narrow results | none |
| `--fields` | Fields to search: `title,content,description,tags` (default: all) | all |
| `--fuzzy` | Enable fuzzy matching (tolerates typos) | `false` |
| `--limit` | Maximum number of results (0 = no limit) | `0` |

#### Examples

```bash
# Search for posts about golang (BM25-ranked)
markata-go search golang

# Search with JSON output for scripting
markata-go search "error handling" --format json

# Search with a filter and limit
markata-go search docker --filter "published == True" --limit 10

# Search only in titles and tags
markata-go search golang --fields title,tags

# Fuzzy search (tolerates typos like "tutoral" → "tutorial")
markata-go search tutoral --fuzzy

# Search and sort by date instead of relevance
markata-go search cli --sort date --order desc

# Get matching file paths for piping
markata-go search kubernetes --format path

# Combine filter expressions with search
markata-go search tutorial --filter '"go" in tags and date >= "2024-01-01"'
```

#### Search Behavior

- Uses bleve full-text index with BM25 ranking (results sorted by relevance by default)
- Falls back to substring matching if the bleve index cannot be built
- Index is cached at `.markata/cache/search.bleve` and rebuilt when content changes
- Combines with `--filter` to narrow results before searching (filter applied first)
- Table output shows relevance scores when using ranked search
- `--fuzzy` enables edit-distance matching for typo tolerance
- **Synonym expansion** is enabled by default — searching "land" also finds posts about "shore" (powered by a bundled WordNet thesaurus with 10,000+ synonym groups)
- **Privacy-safe** — private posts are searchable by public metadata, while their body content and media remain excluded from search

---

### search build-index

Build a reusable bleve index artifact without starting the search server.

#### Usage

```bash
markata-go search build-index [flags]
```

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--index-dir` | Directory to write the bleve index | `.markata/cache/search.bleve` |
| `--hash-path` | Path for the content hash file | `.markata/cache/search.hash` |
| `--index-name` | Named index suffix inside the default cache directory | `""` |
| `--force` | Rebuild even if content hash is unchanged | `false` |

#### Examples

```bash
# Build the default cached index
markata-go search build-index

# Build a reusable artifact at a custom path
markata-go search build-index --index-dir /data/search.bleve --hash-path /data/search.hash

# Build a named index alongside other local indexes
markata-go search build-index --index-name server --force
```

#### Behavior

- Loads the site content the same way as `markata-go search`
- Writes or refreshes a bleve index on disk
- Intended for CI, builder jobs, and future read-only search server deployments

---

### search-server

Start a standalone bleve-backed search API server.

#### Usage

```bash
markata-go search-server [flags]
```

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--port` | Port to listen on | `3001` |
| `--host` | Host to bind to | `localhost` |
| `--mode` | Server mode: `runtime-index`, `watch-content`, or `read-only-index` | `runtime-index` |
| `--index-dir` | Directory of the bleve index | runtime default |
| `--hash-path` | Path for the content hash file in runtime-index mode | runtime default |
| `--index-name` | Named index suffix inside the default cache directory | `server` |
| `--rebuild-index` | Force a rebuild in runtime-index mode | `false` |
| `--watch-debounce` | Debounce duration for watch-content reloads | `750ms` |

#### Examples

```bash
# Start the search API server
markata-go search-server

# Start on a custom port
markata-go search-server --port 8081

# Start in watch mode
markata-go search-server --mode watch-content

# Serve a prebuilt index artifact without loading content
markata-go search-server --mode read-only-index --index-dir /data/search.bleve

# Query the API
curl "http://localhost:3001/api/search?q=golang&fuzzy=true&limit=10"
curl "http://localhost:3001/api/search?q=docker&tags=devops&from=2024-01-01"
```

#### API Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `q` | Search query (required) | — |
| `fuzzy` | Enable fuzzy matching (`true`/`false`) | `false` |
| `limit` | Max results | `20` |
| `tags` | Comma-separated tag filter | none |
| `from` | Date range start (`YYYY-MM-DD`) | none |
| `to` | Date range end (`YYYY-MM-DD`) | none |

#### Search API Behavior

- Returns JSON with `query`, `total`, `results` fields
- Private post content is never searchable — only metadata (title, description, tags)
- CORS follows `search.bleve.cors_origins`
- Also available at the configured endpoint during `markata-go serve`
- Health check at `/health`

#### Modes

- `runtime-index` loads site content and builds or refreshes a local bleve index
- `watch-content` reloads content and refreshes the local index when source files change
- `read-only-index` opens an existing prebuilt index without loading or rebuilding site content

#### Configuration

```toml
[search]
backend = "bleve"           # "pagefind" (default) or "bleve"
endpoint = "/api/search"    # API endpoint path

[search.bleve]
endpoint = "https://search.example.com/api/search"  # Optional remote endpoint for navbar client
fuzzy = false               # Default fuzzy matching
limit = 20                  # Default result limit
max_limit = 100             # Maximum allowed limit
cors_origins = ["*"]        # Allowed CORS origins
```

---

### init

Initialize a new markata-go project with interactive setup.

#### Usage

```bash
markata-go init [flags]
```

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Overwrite existing files | `false` |
| `--plain` | Use plain text prompts instead of the TUI wizard | `false` (auto-detected) |

#### Examples

```bash
# Start interactive project setup (uses TUI wizard when available)
markata-go init

# Overwrite existing configuration
markata-go init --force

# Use plain text prompts (for scripts or non-TTY environments)
markata-go init --plain
```

#### Interactive Flow

The `init` command provides a rich TUI (Text User Interface) wizard powered by [charmbracelet/huh](https://github.com/charmbracelet/huh) for an enhanced user experience. The wizard automatically falls back to plain text prompts when:

- The `--plain` flag is specified
- stdin is not a TTY (e.g., when piping input or running in CI)

**TUI Mode Features:**
- Searchable palette selection
- Multi-select for features
- Form validation
- Styled using your site's configured palette (for existing projects)
- Post format prompts include optional ANSI terminal output

**Plain mode example:**

```
$ markata-go init --plain

Welcome to markata-go!

Site title [My Site]: My Awesome Blog
Description [A site built with markata-go]: A blog about things
Author []: Your Name
URL [https://example.com]: https://myblog.com

Creating project structure...
  ✓ Created posts/
  ✓ Created static/
  ✓ Created markata-go.toml

Create your first post? (Y/n): y
Post title [Hello World]: My First Post

  ✓ Created pages/post/my-first-post.md

Done! Run 'markata-go serve' to start.
```

#### What Gets Created

The command creates:

1. **markata-go.toml** - Configuration file with your site settings
2. **posts/** - Directory for your blog posts
3. **static/** - Directory for static assets (images, CSS, etc.)
4. **(Optional) First post** - A starter markdown file

---

### new

Create a new content file with frontmatter template.

#### Usage

```bash
markata-go new [title] [flags]
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `title` | The title of the new content | No (prompted if not provided) |

#### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--template` | `-t` | Content template to use | `post` |
| `--list` | `-l` | List available templates | `false` |
| `--dir` | | Directory to create the content in (overrides template) | Template's default |
| `--draft` | | Create as a draft | `false` |
| `--tags` | | Comma-separated list of tags | `""` |
| `--plain` | | Use plain text prompts instead of TUI wizard | `false` |

#### Built-in Templates

| Template | Directory | Description |
|----------|-----------|-------------|
| `post` | `pages/post/` | Blog posts with standard frontmatter |
| `page` | `pages/` | Static pages |
| `docs` | `docs/` | Documentation pages |
| `article` (`blog-post`, `essay`, `tutorial`) | `pages/article/` | Long-form article content |
| `note` (`ping`, `thought`, `status`, `tweet`) | `pages/note/` | Short note and micro-post content |
| `photo` (`shot`, `shots`, `image`, `gallery`) | `pages/photo/` | Image-focused posts |
| `video` (`clip`, `cast`, `stream`) | `pages/video/` | Video posts |
| `link` (`bookmark`, `til`, `stars`) | `pages/link/` | Link and bookmark posts |
| `quote` (`quotation`) | `pages/quote/` | Quote posts with attribution |
| `guide` (`series`, `step`, `chapter`) | `pages/guide/` | Guides and multi-part content |
| `inline` (`gratitude`, `micro`) | `pages/inline/` | Inline feed-first content |
| `contact` (`character`, `person`) | `pages/contact/` | Profile/contact pages |
| `author` | `pages/author/` | Author profile pages |

#### Examples

```bash
# Create a new post (creates pages/post/my-first-post.md)
markata-go new "My First Post"

# Use a different template
markata-go new "About" --template page
# Creates: pages/about.md

markata-go new "Getting Started" -t docs
# Creates: docs/getting-started.md

# List available templates
markata-go new --list

# Override the template's directory
markata-go new "Hello World" --template post --dir blog
# Creates: blog/hello-world.md

# Create as a draft (opt-in, default is published)
markata-go new "Work in Progress" --draft
# Creates: pages/post/work-in-progress.md with draft: true

# Create as published (default behavior)
markata-go new "Ready to Publish"
# Creates: pages/post/ready-to-publish.md with draft: false, published: true

# Create with tags
markata-go new "Go Tutorial" --tags "go,tutorial,programming"
# Creates: pages/post/go-tutorial.md with tags: ["go", "tutorial", "programming"]

# Create content in another site from any directory
markata-go --site-dir ~/sites/blog new "Release notes" --no-input

# Interactive mode (no arguments) - launches TUI wizard
markata-go new
# TUI wizard prompts for: template, title, directory, tags, privacy, authors

# Interactive mode with plain text prompts (no TUI)
markata-go new --plain
```

#### Template System

Templates control the default frontmatter and output directory for new content.

**Template Discovery:**

1. **Built-in templates** - `post`, `page`, `docs`, `article`, `note`, `photo`, `video`, `link`, `quote`, `guide`, `inline`, `contact`, and `author` are always available
2. **Config templates** - Defined in `markata-go.toml` under `[content_templates]`
3. **File templates** - Markdown files in `content-templates/` directory

Built-in content types with aliases are shown in `markata-go new --list` as `type (aka: alias1, alias2)`. You can also pass those aliases directly to `--template`; for example, `--template ping` resolves to the built-in `note` template while preserving `template: ping` in frontmatter.

**Configuration Example:**

```toml
[content_templates]
directory = "content-templates"

[content_templates.placement]
post = "blog"       # Override: posts go to blog/
page = "pages"
docs = "documentation"

[[content_templates.templates]]
name = "tutorial"
directory = "tutorials"
body = "## Prerequisites\n\n## Steps\n\n## Summary"

[content_templates.templates.frontmatter]
templateKey = "tutorial"
series = ""
```

**File Template Example:**

Create `content-templates/recipe.md`:

```markdown
---
templateKey: recipe
_directory: recipes
prep_time: ""
cook_time: ""
servings: 4
---

## Ingredients

-

## Instructions

1.
```

The `_directory` field in frontmatter sets the output directory (removed from generated content).

#### Interactive Mode

When called without a title argument (or no arguments at all), the command launches an
interactive TUI wizard powered by [charmbracelet/huh](https://github.com/charmbracelet/huh).

The wizard guides you through:

1. **Template selection** - Pick from built-in, config, and file-based templates
2. **Title** - Enter the post title
3. **Directory** - Choose from existing directories or enter a custom one, with config-driven defaults per template type
4. **Tags** - Select from existing site tags and/or type custom tags
5. **Privacy** - Mark the post as private
6. **Authors** - If multiple authors are configured, choose whether to use the default or select specific authors
7. **Summary** - Review all choices before creating the file

```
$ markata-go new

? Select a template
  > post - Blog post (pages/post/)
    page - Static page (pages/)
    docs - Documentation (docs/)

? Title: My New Post
? Select a directory: pages/post (default)
? Select tags: [go, tutorial]
? Make this post private? No
? Create my-new-post.md in pages/post/? Yes

Created: pages/post/my-new-post.md
```

**Plain text mode:** Use `--plain` or pipe input to get simple text prompts instead of the TUI.
Falls back automatically when stdin is not a TTY.

#### Generated File

The command creates a markdown file with template-specific frontmatter:

```markdown
---
title: My First Post
slug: my-first-post
date: "2024-01-15"
published: true
draft: false
templateKey: post
tags:
  - go
  - tutorial
description: ""
---

Write your content here...
```

The generated body does not include a starter H1. markata-go renders the page title from frontmatter, so content should start with regular body text or an `##` heading.

The slug is automatically generated from the title by:
- Converting to lowercase
- Replacing spaces with hyphens
- Removing special characters

---

### lint

Lint markdown files for common issues that can cause build failures.

#### Usage

```bash
markata-go lint [files...] [flags]
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `files` | Glob patterns or file paths to lint | Yes |

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--fix` | Automatically fix detected issues | `false` |

#### Detected Issues

The linter checks for these common problems:

| Issue | Severity | Auto-fixable | Description |
|-------|----------|--------------|-------------|
| `duplicate-key` | Error | Yes | Duplicate YAML keys in frontmatter |
| `invalid-date` | Warning | Yes | Non-ISO 8601 date formats |
| `missing-alt-text` | Warning | Yes | Image links without alt text `![]()` |
| `protocol-less-url` | Warning | Yes | URLs starting with `//` instead of `https://` |
| `encryption-key-policy` | Error | No | Missing or weak encryption keys based on `encryption.*` policy |

#### Examples

```bash
# Lint all markdown files in posts directory
markata-go lint posts/*.md

# Lint with glob pattern (recursive)
markata-go lint 'posts/**/*.md'

# Lint multiple patterns
markata-go lint posts/*.md pages/*.md docs/**/*.md

# Lint and auto-fix issues
markata-go lint posts/*.md --fix

# Lint a specific file
markata-go lint posts/hello-world.md

# Lint with verbose output
markata-go lint posts/*.md -v
```

#### Output Format

```
posts/example.md:
  error [line 5]: duplicate key 'title' (first occurrence at line 2)
  warning [line 3]: invalid date format for 'date': 2020-1-15 (single-digit month/day)
  warning [line 12, col 1]: image link missing alt text

✗ 5 file(s) linted, 3 issue(s) in 1 file(s)
```

#### Auto-fix Behavior

When `--fix` is enabled:

| Issue | Fix Applied |
|-------|-------------|
| `duplicate-key` | Keeps last occurrence, removes earlier duplicates |
| `invalid-date` | Pads single-digit months/days (e.g., `2020-1-5` → `2020-01-05`) |
| `missing-alt-text` | Adds placeholder alt text: `![]()` → `![image]()` |
| `protocol-less-url` | Adds HTTPS protocol: `//example.com` → `https://example.com` |

#### Exit Codes

| Code | Description |
|------|-------------|
| `0` | No errors found (warnings are allowed) |
| `1` | Errors found (or file read errors) |

#### Use Cases

**Migration from Python markata:**
```bash
# Find and fix common issues when migrating content
markata-go lint content/**/*.md --fix
```

**CI/CD integration:**
```bash
# Fail CI if markdown has errors
markata-go lint posts/**/*.md || exit 1
```

**Pre-commit hook:**
```bash
#!/bin/bash
markata-go lint $(git diff --cached --name-only -- '*.md')
```

---

### config

Configuration management commands for viewing, validating, and initializing configuration.

#### Usage

```bash
markata-go config <subcommand> [flags]
```

Running `markata-go config` with no subcommand behaves like `markata-go config show`.

#### Subcommands

##### show

Display the resolved configuration after merging defaults, config file, and environment variables.

`config show` honors the root `--config` and `--merge-config` flags, so it reports the same effective configuration used by `build` and `serve`.

```bash
markata-go config show [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format: `yaml`, `json`, or `toml` | `yaml` |
| `--json` | Output as JSON | `false` |
| `--toml` | Output as TOML | `false` |
| (none) | Output as YAML | default |

**Examples:**

```bash
# Show config as YAML (default)
markata-go config show

# Show config as JSON
markata-go config show --json

# Show config as TOML
markata-go config show --toml

# Show config from specific file
markata-go config show -c production.toml

# Show config with merged overrides
markata-go config show -m fast.toml
```

Notes:

- `markata-go config` is equivalent to `markata-go config show`
- `--json` and `--toml` are shorthands for `--format json` and `--format toml`
- conflicting combinations such as `--json --toml` fail with usage error exit code `2`

##### get

Get a specific configuration value using dot notation for nested keys.

```bash
markata-go config get <key>
```

**Examples:**

```bash
# Get top-level value
markata-go config get output_dir
# Output: output

# Get nested value (list values print one item per line)
markata-go config get glob.patterns
# Output:
# pages/**/*.md
# posts/**/*.md

# Get deeply nested value
markata-go config get feed_defaults.items_per_page
# Output: 10

# Get from specific config file
markata-go config get url -c production.toml
# Output: https://example.com
```

Notes:

- Values are read directly from the config file; `config get` does not synthesize absent defaults.
- TOML/YAML are parsed with tree-sitter in CGO-enabled builds to preserve formatting.
- CGO-disabled builds parse TOML/YAML via full decode.
- JSON values may be re-emitted for structured output.

##### set

Set a configuration value using dot notation for nested keys.

```bash
markata-go config set <key> <value>
```

**Examples:**

```bash
# Set a string value
markata-go config set url "https://example.com"

# Set an integer value
markata-go config set concurrency 4

# Set an array value (JSON syntax)
markata-go config set glob.patterns '["pages/**/*.md", "posts/**/*.md"]'
```

Notes:

- TOML/YAML are updated with byte-range edits to preserve formatting and comments in CGO-enabled builds.
- CGO-disabled builds re-emit TOML/YAML with standard encoders.
- JSON is re-emitted with stable indentation.
- File permissions are preserved.

##### validate

Validate the configuration file and report any errors or warnings.

`config validate` honors the root `--config` and `--merge-config` flags.

```bash
markata-go config validate [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Config file to validate |

**Examples:**

```bash
# Validate auto-discovered config, or built-in defaults if no config exists
markata-go config validate

# Validate specific config file
markata-go config validate -c production.toml
markata-go config validate -c custom.yaml

# Validate with merged overrides
markata-go config validate -m fast.toml
```

**Output:**

```
Configuration is valid: markata-go.toml

# Or with errors:
Errors:
  - output_dir: directory does not exist
  - glob.patterns: at least one pattern required
  - feeds[0].filter: invalid filter expression
```

##### init

Create a new configuration file with sensible defaults.

```bash
markata-go config init [filename] [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Overwrite existing config file | `false` |

**Examples:**

```bash
# Create markata-go.toml with defaults
markata-go config init

# Create YAML config
markata-go config init markata-go.yaml

# Create with specific name
markata-go config init site.toml

# Overwrite existing config
markata-go config init --force
```

---

### aesthetic

Manage and inspect aesthetic presets for your site's visual styling.

#### Usage

```bash
markata-go aesthetic <subcommand>
```

#### Subcommands

##### list

List all available aesthetic presets.

```bash
markata-go aesthetic list
```

**Example Output:**

```
Available aesthetics:

  balanced    Default. Comfortable rounding, subtle shadows, normal spacing
  brutal      Sharp corners, thick borders, tight spacing, no shadows
  elevated    Generous rounding, layered shadows, generous spacing
  minimal     No rounding, maximum whitespace, no shadows, hairline borders
  precision   Subtle corners, compact spacing, hairline borders, minimal shadows
```

##### show

Display the details of a specific aesthetic preset.

```bash
markata-go aesthetic show <name>
```

**Arguments:**

| Argument | Description | Required |
|----------|-------------|----------|
| `name` | Name of the aesthetic to show | Yes |

**Examples:**

```bash
# Show details of the elevated aesthetic
markata-go aesthetic show elevated

# Show details of the brutal aesthetic
markata-go aesthetic show brutal
```

**Example Output:**

```
Aesthetic: elevated

Description: Generous rounding, layered shadows, generous spacing

Values:
  border_radius:    16px
  spacing_scale:    1.5
  border_width:     0px
  border_style:     none
  shadow_intensity: 0.6
  shadow_size:      lg
```

#### Available Aesthetics

| Aesthetic | Description |
|-----------|-------------|
| `balanced` | Default. Comfortable rounding, subtle shadows, normal spacing |
| `brutal` | Sharp corners, thick borders, tight spacing, no shadows |
| `elevated` | Generous rounding, layered shadows, generous spacing |
| `minimal` | No rounding, maximum whitespace, no shadows, hairline borders |
| `precision` | Subtle corners, compact spacing, hairline borders, minimal shadows |

See [[themes-and-styling|Themes Guide]] for detailed aesthetic customization.

---

### palette

Manage, browse, and select color palettes for your site's theme.

#### Usage

```bash
markata-go palette <subcommand> [flags]
```

#### Subcommands

##### list

List all available palettes from built-in, user, and project sources.

```bash
markata-go palette list
markata-go palette list --variant dark
markata-go palette list --json
```

##### info

Show detailed information about a specific palette including colors, semantic mappings, and component colors.

```bash
markata-go palette info catppuccin-mocha
markata-go palette info catppuccin-mocha --json
```

##### pick

Open an interactive full-screen TUI to browse and select a palette with live color previews. Sets the chosen palette in your config by default.

```bash
# Pick and set in your config (default)
markata-go palette pick

# Only print the name without updating config
markata-go palette pick --no-set
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--no-set` | Only print the palette name without updating the config file |

The picker shows a two-panel layout with a fuzzy-filterable palette list on the left and a live color swatch preview on the right. Type to filter, arrow keys to navigate, Enter to select, Esc to cancel.

**Composability:**

```bash
# View info for the palette you pick (without setting it)
markata-go palette info "$(markata-go palette pick --no-set)"
```

##### check

Check a palette for accessibility (WCAG contrast ratios).

```bash
markata-go palette check catppuccin-mocha
```

##### preview

Preview a palette's colors in the terminal.

```bash
markata-go palette preview catppuccin-mocha
```

##### new

Create a new custom palette interactively.

```bash
markata-go palette new my-brand
```

##### clone

Clone an existing palette as a starting point for customization.

```bash
markata-go palette clone catppuccin-mocha my-custom
```

##### export

Export a palette to different formats (CSS, SCSS, JSON).

```bash
markata-go palette export catppuccin-mocha --format css
```

##### fetch

Import a palette from [Lospec.com](https://lospec.com/palette-list).

```bash
markata-go palette fetch https://lospec.com/palette-list/sweetie-16.txt
```

See [[themes-and-styling|Themes Guide]] for detailed palette customization.

---

## Environment Variables

markata-go configuration can be overridden using environment variables with the `MARKATA_GO_` prefix.

### Variable Naming Convention

Environment variables follow this pattern:
- Prefix: `MARKATA_GO_`
- Key names: UPPERCASE with underscores
- Nested keys: Use single underscore separation

### Common Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `MARKATA_GO_OUTPUT_DIR` | Output directory | `dist` |
| `MARKATA_GO_URL` | Site base URL | `https://staging.example.com` |
| `MARKATA_GO_TITLE` | Site title | `My Staging Site` |
| `MARKATA_GO_CONCURRENCY` | Worker count (0=auto) | `4` |
| `MARKATA_GO_TEMPLATES_DIR` | Templates directory | `themes/custom/templates` |
| `MARKATA_GO_ASSETS_DIR` | Static assets directory | `assets` |

### Examples

```bash
# Build to a different directory
MARKATA_GO_OUTPUT_DIR=dist markata-go build

# Build for staging environment
MARKATA_GO_URL=https://staging.example.com markata-go build

# Limit concurrency for debugging
MARKATA_GO_CONCURRENCY=1 markata-go build -v

# Multiple overrides
MARKATA_GO_OUTPUT_DIR=dist \
MARKATA_GO_URL=https://prod.example.com \
markata-go build --clean

# In CI/CD pipelines
export MARKATA_GO_URL="${DEPLOY_URL}"
export MARKATA_GO_OUTPUT_DIR="build"
markata-go build --clean
```

### Override Precedence

Configuration values are resolved in this order (later sources override earlier):

1. Built-in defaults
2. Configuration file
3. Environment variables
4. Command-line flags

---

## Exit Codes

markata-go uses these common exit code patterns:

| Code | Name | Description |
|------|------|-------------|
| `0` | Success | Command completed successfully |
| `1` | Error | General error (configuration, plugin, I/O, etc.) |
| `2` | Usage / Command-specific | Invalid CLI usage such as unknown flags or missing args, or command-specific non-success states documented by that command |
| `130` | Interrupted | Command interrupted by user (Ctrl+C) |

### Using Exit Codes in Scripts

```bash
#!/bin/bash

# Build and check for success
if markata-go build --clean; then
    echo "Build successful!"
    rsync -av output/ user@server:/var/www/
else
    echo "Build failed with exit code $?"
    exit 1
fi
```

```bash
# Validate before deploying
markata-go config validate
if [ $? -eq 0 ]; then
    markata-go build --clean
fi
```

---

## Common Workflows

### Development

```bash
# Start development server
markata-go serve -v

# In another terminal, create new content
markata-go new "My New Post"
```

### Production Build

```bash
# Clean build for production
markata-go build --clean -v

# Or with environment-specific URL
MARKATA_GO_URL=https://example.com markata-go build --clean
```

### CI/CD Pipeline

```bash
#!/bin/bash
set -e

# Validate configuration
markata-go config validate

# Build with production settings
MARKATA_GO_URL="${DEPLOY_URL}" markata-go build --clean

# Deploy (example)
aws s3 sync output/ s3://my-bucket/ --delete
```

See [[deployment-guide|Deployment]] for detailed deployment guides.

---

### migrate

Migrate from Python markata to markata-go. Analyzes configuration files, filter expressions, and templates for compatibility.

#### Usage

```bash
markata-go migrate [flags]
markata-go migrate config [flags]
markata-go migrate filter [expression]
markata-go migrate templates [path]
```

#### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--input` | `-i` | Input config file path | Auto-detect |
| `--output` | `-o` | Output config file path | None |
| `--dry-run` | `-n` | Show changes without writing | `false` |
| `--format` | `-f` | Output format (toml, yaml) | `toml` |
| `--json` | | Output results as JSON | `false` |

#### Subcommands

##### config

Migrate configuration file only.

```bash
markata-go migrate config -i markata.toml -o markata-go.toml
```

##### filter

Check and migrate filter expressions.

```bash
# Check a specific expression
markata-go migrate filter "published == 'True'"

# Check all filters in config
markata-go migrate filter
```

##### templates

Check template compatibility with pongo2.

```bash
markata-go migrate templates
markata-go migrate templates ./my-templates
```

#### Examples

```bash
# Full migration analysis
markata-go migrate

# Dry run to see what would change
markata-go migrate --dry-run

# Migrate and write output
markata-go migrate -o markata-go.toml

# Migrate from pyproject.toml
markata-go migrate -i pyproject.toml -o markata-go.toml

# Check a filter expression
markata-go migrate filter "templateKey in ['blog', 'til']"

# Check templates for compatibility
markata-go migrate templates
```

#### Exit Codes

| Code | Description |
|------|-------------|
| `0` | Migration successful, no issues |
| `1` | Migration completed with warnings |
| `2` | Migration has incompatibilities |
| `3` | Migration failed (invalid input) |

See [[migration-guide|Migration Guide]] for detailed migration instructions.

---

### explain

Show detailed information about markata-go for AI agents and developers.

#### Usage

```bash
markata-go explain [topic]
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `topic` | The topic to explain | No (shows overview if omitted) |

#### Available Topics

| Topic | Description |
|-------|-------------|
| (none) | General overview of markata-go |
| `build` | The build command and build process |
| `serve` | The development server |
| `new` | Creating new content |
| `content` | Finding, creating, and editing site content |
| `init` | Initializing projects |
| `config` | Configuration system |
| `agents` | Agent skill installation and usage |
| `plugins` | Plugin system and development |
| `lifecycle` | Build lifecycle stages |
| `templates` | Template system |
| `feeds` | Feed generation system |

#### Examples

```bash
# Show general overview
markata-go explain

# Get detailed info about the build process
markata-go explain build

# Learn about the plugin system
markata-go explain plugins

# Learn about the bundled agent skill
markata-go explain agents

# Follow the agent-oriented content workflow
markata-go explain content

# Understand the build lifecycle
markata-go explain lifecycle

# Learn about feeds and filtering
markata-go explain feeds
```

#### Use Cases

The `explain` command is particularly useful for:

- **AI coding agents**: Provides comprehensive context about markata-go's architecture and commands
- **Content work from any directory**: `explain content` documents how to select a site, search, edit returned source paths, create content, and validate changes
- **Site maintainers installing the bundled skill**: Shows how `markata-go agent install` fits into agent workflows
- **New developers**: Quick reference for understanding how different parts work together
- **Debugging**: Understanding the build process and configuration options
- **Writing plugins**: Learning the plugin interfaces and lifecycle stages

---

### encryption

Utilities for managing encryption passwords.

#### Subcommands

##### generate-password

Generate a password that satisfies the default encryption policy (≥ 14 characters, strong entropy).

```
markata-go encryption generate-password [--length <n>]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--length` | | Desired password length (must be at least the configured minimum) | `14` |

The command prints only the generated password to stdout, making it easy to pipe into environment files or password managers. The optional `--length` flag requests a longer password; the minimum is always at least the configured `min_password_length`.

#### Examples

```
# Generate a default length password
markata-go encryption generate-password

# Generate a longer password for a stricter key
markata-go encryption generate-password --length 20
```

##### check

Check configured encryption keys against the active policy.

```
markata-go encryption check [--key <name>]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--key` | | Check a specific key name instead of all configured keys | `""` |

By default this command checks all keys referenced by `encryption.default_key` and `encryption.private_tags`.
It fails with a non-zero exit code if any key is missing from environment variables or too weak for the configured policy.

##### encrypt-posts

Encrypt all private Markdown source bodies matched by the active content glob configuration.

```
markata-go encryption encrypt-posts [--dry-run]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dry-run` | | Report files that would be encrypted without modifying them | `false` |

The command rewrites matching private posts in place by default. It skips draft, skipped, public, and already source-encrypted posts. Missing or weak keys fail before any files are written.

##### decrypt-posts

Decrypt source-encrypted Markdown bodies back to plaintext. This is the inverse of `encrypt-posts`.

```
markata-go encryption decrypt-posts [path...] [--dry-run]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dry-run` | | Report files that would be decrypted without modifying them | `false` |

With no path arguments the command scans the active content glob configuration. Explicit paths may be files or
directories (directories are scanned recursively for `.md` files).

The key name is read from the `<!-- markata-encrypted-source:v1 key=... -->` marker, falling back to
`encryption.default_key`. The password is read from `MARKATA_GO_ENCRYPTION_KEY_<KEY>`. Files that are not
source-encrypted are skipped, and frontmatter is preserved exactly as written.

#### Examples

```
# Preview source encryption changes
markata-go encryption encrypt-posts --dry-run

# Encrypt private post source bodies in place
markata-go encryption encrypt-posts

# Decrypt a single file so you can edit it
export MARKATA_GO_ENCRYPTION_KEY_DEFAULT='your-password'
markata-go encryption decrypt-posts pages/ar-calibration.md

# Preview decrypting every encrypted post
markata-go encryption decrypt-posts --dry-run
```

---

## See Also

- [[getting-started|Getting Started]] - Quick introduction to markata-go
- [[configuration-guide|Configuration]] - Detailed configuration reference
- [[deployment-guide|Deployment]] - Deployment guides for various platforms
---

## Editorial authoring

Commit selected post edits with transparent contributor credit and exact
revision links. These commands preserve unrelated staged and unstaged work.

```bash
markata-go editorial start posts/article.md --author editor-id --role copy-editor
markata-go editorial finish --session UUID --summary "Corrected grammar" --json
markata-go editorial commit posts/article.md --author human-id \
  --role internal-link-editor --summary "Added related links" --dry-run
```

`start` requires clean tracked posts; `commit` explicitly attributes existing
owned changes. `finish` appends revision records after the content SHA exists.
Use `--json` to obtain the full list of resulting commits. See
[Editorial authoring](../../guides/editorial-authoring/) and
[Editorial history](../../guides/editorial-history/) for ownership rules,
primary author configuration, and the built-site comparison view.
