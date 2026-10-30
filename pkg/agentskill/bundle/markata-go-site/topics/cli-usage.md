# CLI Usage

Use this topic for everyday site work and safe project inspection.

## Terminal Colors

CLI output, logs, forms, and TUIs share semantic colors from the configured site palette.
The theme's `fallback_mode` selects the light or dark palette variant.
Without a usable site palette, these surfaces share one fallback theme.
Palette picker swatches still show each previewed palette.
Use `--no-color` or `--log-format plain` to suppress terminal colors.

## High-Value Commands

### Build And Serve

- `markata-go build`
- `markata-go build pages/sample.md` (render one Markdown file to `<output>/index.html` with the default theme; no feeds, listings, sitemap, search index, or build cache)
- `markata-go pages/sample.md` (equivalent one-file invocation)
- `markata-go build --clean` (remove output + build cache)
- `markata-go build --clean-all` (also remove external caches: blogroll, embeds, mentions)
- `markata-go build --fast` (skip minification, CSS purging, Pagefind indexing)
- `markata-go build --dry-run` (run through Collect, show counts, skip Write)
- `markata-go reader update` (refresh external feed cache for `/reader/` without building)
- `markata-go reader update --concurrency 12` (override reader refresh parallelism for one run)
- `MARKATA_GO_BLOGROLL_REFRESH_ON_BUILD=false markata-go build` (keep blogroll pages but skip remote refresh during the build)
- `markata-go build --benchmark-json=benchmark.json`
- `markata-go build --dag` (experimental serial task graph, not extra parallelism)
- `markata-go build --dag=false` (override `MARKATA_GO_DAG=true` for one build)
- `markata-go build -v --benchmark-detailed`
- `markata-go buildlab run --fixture /path/to/site` (compare clean and incremental builds)
- `markata-go serve` (dev server with live reload; its theme picker adds a **Bake** button that writes the chosen look into the config file holding `[markata-go.theme]`, and a gear button opens a settings sidebar that live-previews any setting in memory (Reset discards) and bakes it into the config file that owns it after showing a diff, and can reset a setting to its default by removing it from the file; neither exists in `build` output, and neither writes the global `~/.config/markata-go/config.toml`. Previews build into `.markata/serve-preview/`, so the main cache is untouched. What the served site shows may be an unbaked preview, so trust config files, not the browser, and re-read them before editing)
- `markata-go serve --fast`
- `markata-go serve pages/post.md` (serve only that file at `/`, same single-file mode as `build pages/post.md`)
- `markata-go serve --incremental` (reuse unchanged posts while retaining normal output processing)
- `markata-go serve --no-tui` (line-oriented logs for scripts, CI, or terminal capture)
- `markata-go serve` on loopback also serves the local Control Center at `http://localhost:8000/_markata/`; press `a` in the TUI to open it
- `markata-go serve --no-admin` (disable the local Control Center)
- `markata-go serve --admin` (explicitly enable it on loopback; rejected on non-loopback binds)
- `markata-go admin` (same local serve session and web dashboard, with line-oriented terminal output)

In a real terminal, `serve` opens the Serve Control Center. Inspect the warnings
and errors inboxes or select a page to find a diagnostic's source location and
suggested fix. Build failures remain in job history after later rebuilds. The
local web dashboard and terminal use one session state. Use `builder-admin` for
the separate deployed release workflow behind proxy authentication.

In the local browser dashboard, the URL fragment preserves the section, selected
item, and filter for sharing or browser Back/Forward. Use `j`/`k` or arrows to
move through results, Enter to inspect, `/` to focus the filter, `w`/`e` for
warnings/errors, and `p`/`f`/`l` for pages/feeds/logs. Press `?` for the full key
list. Press `r` in Jobs to rerun the selected completed job. These shortcuts
pause while editing text. Review each source file's full before/after diff
before applying suggested fixes; the local source actions are not part of deployed
Builder Admin.

`serve` and `builder-admin` also accept `--dag`; an explicit `--dag=false`
overrides the environment for their builds, including Builder Admin child
processes. An invalid `MARKATA_GO_DAG` boolean reports an error unless an
explicit flag overrides it. `build --dry-run` uses the legacy partial lifecycle,
not the DAG scheduler. Trust the executor recorded in benchmark JSON and
full-build diagnostics rather than inferring it from an environment variable.

### Config And Inspection

- `markata-go config show`
- `markata-go config show` (annotated YAML by default; use `--no-annotate` for plain YAML)
- `markata-go config show --diff`
- `markata-go config get <key>`
- `markata-go config validate`

Bare `markata-go config` behaves like `markata-go config show`.
- `markata-go list posts`
- `markata-go list feeds`
- `markata-go list tags`

### Search

- `markata-go search <query>` (BM25-ranked full-text search with synonym expansion)
- `markata-go search <query> --format json` (machine-readable output)
- `markata-go search <query> --filter "published == True"` (combine with filter)
- `markata-go search <query> --fields title,tags` (restrict fields)
- `markata-go search <query> --fuzzy` (typo-tolerant matching)
- `markata-go search <query> --limit 10` (cap results)
- `markata-go search <query> --format path` (file paths only, for piping)
- `markata-go search <query> --sort date` (sort by date instead of relevance)
- `markata-go search-server` (standalone search API server)
- `markata-go search-server --port 8081` (custom port)
- During `markata-go serve`, the search API is auto-mounted at `/api/search`
- `curl "http://localhost:3001/api/search?q=golang&fuzzy=true&limit=10"` (query API)

### Content Creation

- `markata-go new` (create content from built-in templates)
- `markata-go new --list` (list available content templates)
- `markata-go init` (initialize a new project with TUI wizard)
- `markata-go init --plain` (plain text prompts for non-TTY environments)

### Content Workflow For Agents

Use `markata-go explain content` for the complete inspect, find, edit, create,
and validate workflow. From another directory, select the site explicitly:

```bash
markata-go --site-dir ~/sites/blog list posts --format json
markata-go --site-dir ~/sites/blog search "release process" --format path
markata-go --site-dir ~/sites/blog new "Release notes" --no-input
markata-go --site-dir ~/sites/blog lint
markata-go --site-dir ~/sites/blog build --fast
```

`list` and `search` return absolute Markdown source paths with `--site-dir`.
Edit a returned path directly; there is no separate non-interactive edit command.
Set `MARKATA_GO_SITE_DIR` in an alias or wrapper when repeatedly using one site.

### Content Quality

- `markata-go lint` (lint markdown files for common issues)
- `markata-go lint --fix` (auto-fix fixable issues)
- `markata-go lint --dry-run` (show files without linting)
- `markata editor --workspace <site> --port 8766` (start the Go browser editor on loopback)
- The Go and FastAPI browser editor examples use the existing LSP through a same-origin `/lsp` WebSocket.
  Go uses `markata editor --workspace`; FastAPI uses `MARKATA_LSP_WORKSPACE`.
  FastAPI uses `markata-go lsp --browser-workspace` behind its stdio bridge.
  Drafts stay in memory. The Go and Python editors open files discovered by Markata glob rules; edits stay in browser drafts and do not write source files.
  The editors complete command-catalog containers and attributes, plus Mermaid
  diagram types and directions. Press Ctrl+Space to show completions.
  See `docs/guides/go-workbench.md` and `docs/guides/fastapi-editor.md` in the Markata repository.
- `markata-go lsp doctor` (check local LSP command and project configuration)
- `markata-go lsp doctor --no-verify-editor` (detect installed editors without loading their configuration)
- `markata-go lsp setup` (print setup guidance for each installed supported editor; use `--editor <editor>` for `generic`, `neovim`, `helix`, `emacs`, `zed`, or `vscode`)

Every build prints a `Content:` summary. Use it to check the path from
discovery to output: candidates, loaded sources, valid frontmatter, posts,
eligible posts, rendered posts, emitted outputs, exclusions, warnings, and
errors. Content diagnostics include the source path, line, stable reason code,
and message. `-v` also lists excluded and shadow sources with their reasons.

If a frontmatter block starts with `----` or has a malformed closing delimiter,
fix the source delimiter. Do not assume the file was ignored. Inspect the
content diagnostic before changing feed filters or templates. The
`--benchmark-json` output includes the same sanitized snapshot under `content`.
Malformed author content can be isolated to one source. Source stat, read, and
other filesystem errors are operational failures and make the build fail.

After a successful full build, inspect the persisted snapshot at
`<output_dir>/.markata/diagnostics.json`:

```bash
jq . public/.markata/diagnostics.json
```

The artifact has versioned summary and entry fields, stable reason codes, and
feed-level dispositions. It is sanitized: it does not include raw configuration,
secrets, frontmatter, Markdown bodies, rendered HTML, or absolute paths. Replace
`public` when the site uses another output directory. `build --dry-run`,
`build --fast`, fast or incremental `serve`, failed, or incomplete builds do not
publish a new artifact. A normal `serve` build can publish it after a successful
full lifecycle. See the Content Diagnostics guide when the site
repository includes the markata-go documentation.

### Encryption

- `markata-go encryption generate-password` (generate a policy-compliant password)
- `markata-go encryption check` (verify configured encryption keys)
- `markata-go encryption encrypt-posts --dry-run` (preview source encryption for private posts)
- `markata-go encryption encrypt-posts` (encrypt private Markdown bodies in place)
- `markata-go encryption encrypt-posts --workers 4` (bound concurrent preparation; default `0` uses `GOMAXPROCS`)
- `markata-go encryption decrypt-posts --dry-run` (preview decrypting source-encrypted posts)
- `markata-go encryption decrypt-posts [path...]` (decrypt source-encrypted Markdown bodies in place)
- Bulk encryption and decryption prepare every file before writing, so a missing key or invalid ciphertext does not leave the repository partially transformed.

### Theme And Palette

- `markata-go palette list`
- `markata-go palette info <name>`
- `markata-go palette check <name>` (WCAG contrast validation)
- `markata-go palette check <name> --strict` (WCAG AAA instead of AA)
- `markata-go palette check --all` (check all palettes)
- `markata-go palette preview <name>`
- `markata-go palette new <name>`
- `markata-go palette clone <source>`
- `markata-go theme render-all`
- `markata-go theme gallery`
- `markata-go theme check-all` (check 16 contrast combos per palette)
- `markata-go theme check-all --colorblindness` (simulate color vision deficiencies)
- `markata-go aesthetic list`
- `markata-go aesthetic show <name>`

### Explain

- `markata-go explain` (list topics)
- `markata-go explain config`
- `markata-go explain templates`
- `markata-go explain plugins`
- `markata-go explain agents`
- `markata-go explain feeds`
- `markata-go explain lifecycle`

### Migration And Import

- `markata-go migrate config` (convert Python markata config)
- `markata-go migrate filter [expression]` (check filter expression compatibility)
- `markata-go migrate templates [path]` (validate template compatibility)
- `markata-go migrate compare --old <dir> --new <dir>` (compare site outputs)
- `markata-go import rss <url>` (import from RSS/Atom feed)
- `markata-go import jsonfeed <url>` (import from JSON Feed)
- Shared import flags: `--output`, `--since`, `--dry-run`, `--tags`

### Maintenance

- `markata-go update` (self-update from GitHub releases)
- `markata-go update --check` (check for updates without installing)
- `markata-go benchmark --scenario small|medium|large` (performance benchmarks)
- `markata-go agent list-agents` (list supported agent ids and their install paths)
- `markata-go agent install` (install bundled agent skill into the detected project agent or the universal layout)
- `markata-go agent install --agent <name> [-g]` (choose a specific agent and optional global install scope)
- `markata-go agent install -g` (open an fzf-style picker with a bottom prompt and colored matches; scripts must pass `--agent`, such as `github-copilot`; `--no-color` disables colors)
- `markata-go agent doctor` (check for drift after binary upgrades)
- `markata-go version`

## Lint Checks

`markata-go lint` detects:

- duplicate YAML keys in frontmatter
- invalid date formats (non-ISO 8601)
- malformed image links (missing alt text)
- protocol-less URLs (should use `https://`)
- encryption policy issues (when encryption is configured)

Use `--fix` to auto-fix fixable issues. Only error-severity issues cause a non-zero exit code; warnings alone pass.

## Global Flags Agents Should Know

- `-c`, `--config`: use a specific config file
- `--site-dir`: select a site directory for commands run elsewhere; overrides `MARKATA_GO_SITE_DIR`
- `-m`, `--merge-config`: merge override configs such as `fast.toml`
- `-o`, `--output`: override the output directory without editing config
- `-v`, `--verbose`: show detailed logs and plugin-stage hints
- `-q`, `--quiet`: suppress non-essential progress output
- `--no-input`: disable prompts for scripted or non-interactive runs

Examples:

```bash
markata-go build -c markata-go.toml
markata-go --site-dir ~/sites/blog search golang --format path
markata-go serve -m fast.toml
markata-go build -o dist
markata-go new "My Post" --no-input
markata-go build -v
markata-go lint --fix
markata-go encryption encrypt-posts --dry-run
```

## What To Use When

- inspect content inventory: `markata-go list posts`
- search for content by keyword: `markata-go search <query>`
- inspect feed definitions and sizes: `markata-go list feeds`
- inspect resolved configuration: `markata-go config show`
- create new content: `markata-go new`
- lint content before committing: `markata-go lint`
- set up Markdown wikilink IDE support: `markata-go lsp doctor`, then `markata-go lsp setup --editor <editor>`; inspect existing editor configuration before applying the printed values
- source-encrypt private Markdown before committing: `markata-go encryption encrypt-posts --dry-run`, then `markata-go encryption encrypt-posts`
- edit an already source-encrypted post: `markata-go encryption decrypt-posts <path>`, edit, then `markata-go encryption encrypt-posts`
- validate config before deploy: `markata-go config validate`
- validate palette contrast: `markata-go palette check <name>`
- interactive local editing: `markata-go serve --fast`
- refresh reader data without a build: `markata-go reader update`
- final verification: `markata-go build`
- bootstrap a new project: `markata-go init`
- migrate from Python markata: `markata-go migrate config`
- import content from external feeds: `markata-go import rss <url>`

## Operator Patterns

- use `-m fast.toml` when you want a lighter dev build without rewriting main config
- use the same `-m` overrides with `config show` and `config validate` when you need to inspect or verify the exact config that `build` or `serve` will use
- use `-c` when a repo has multiple configs or examples and you need the exact active one
- use `-o dist` in CI or preview contexts when you want a temporary artifact path
- use `--no-input` for automation or when the agent must avoid prompts
- use `-v` when debugging plugin order, missing outputs, or config resolution issues
- use `-v` when debugging why a source became excluded or a shadow page
- inspect the `Content:` build summary before changing feed filters for a missing post
- use `--dry-run` on build or lint to preview behavior without side effects
- use `encryption encrypt-posts --dry-run` before rewriting private source files; make sure required `MARKATA_GO_ENCRYPTION_KEY_*` environment variables are set first
- use `encryption decrypt-posts --dry-run` before decrypting; it is the inverse of `encrypt-posts` and needs the same `MARKATA_GO_ENCRYPTION_KEY_*` variables

## Guidance

- Prefer `list` commands when you need structured inspection.
- Prefer `search` when you need to find posts by content or keyword.
- Prefer `explain content` when an agent needs the complete content workflow, including direct editing of returned paths.
- Prefer `explain` when you need command-specific or subsystem context.
- Prefer `serve` for interactive local work and `build` for validation or CI-like runs.
- Run `lint` before committing content changes.
- Run `palette check` after creating or modifying palettes.
- Use `--verbose` only when normal output is not enough.
- Keep primary results script-friendly by using built-in machine-readable output when available.
- Prefer merged config overrides over editing the main config for temporary local or CI changes.
# CLI discovery and mistakes

`markata-go b`, `markata-go s`/`markata-go serv`, and `markata-go ls` are
aliases for `build`, `serve`, and `list`. The words `dev` and `preview` suggest
`serve`, but do not run it. Use the canonical names in scripts and examples. To listen
on all network interfaces during local preview, run
`markata-go serve --bind 0.0.0.0`; `--bind` and `--host` mean the same thing.

Unknown commands, flags, and close config key typos report possible matches.
Treat each suggestion as a hint and use the fully spelled command or key.
Run the specific command's `--help` when no suggestion fits. A bare Markdown
path such as `markata-go post.md` still runs a one-file build.


### Mention attribution and lint warnings

In Markdown, `@handle` refers to a person and `~handle` attributes a
quote. Both resolve from active, configured blogroll feeds and aliases.
When blogroll is enabled, `markata-go lint` reports
`unknown-mention` warnings for unresolved handles. Add a verified
entry or alias to the blogroll before linking: CLI lint/build never
searches social networks or guesses a person's identity.

Use `markata workspace --workspace <site>` to inspect discovered document metadata as JSON.
Add `--file <relative-path>` to read a discovered source. Python exposes the same API as `markata.workspace_files`.
