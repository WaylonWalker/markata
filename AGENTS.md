# AGENTS.md - AI Coding Agent Guidelines for markata-go

This document provides guidelines for AI coding agents working in the markata-go codebase.

## Merge readiness and CI triage

Before marking a PR ready, fetch origin/main and run `just pr-check`. For editor changes also run `just pr-check --browser` with Playwright/Chromium. Local checks are **not** substitutes for platform or release gates.

- Batch related changes before pushing, rather than repeatedly pushing untested commits.
- Check the latest head SHA with `gh pr checks`, not old or cancelled runs.
- Use `just ci-diagnose RUN_ID` or `gh run view RUN_ID --log-failed` to find the first underlying failure.
- Retry identical SHAs only for evidence-backed flakiness or infrastructure failures; never bypass required checks.
- CI uses Go 1.27 and golangci-lint v2.15.0. `just lint-new` gates new issues; `just lint-full` audits legacy debt.

## GitHub Workflow

**This project uses a formal GitHub workflow.** All work should be tracked and managed through GitHub:

### Issues

- **All work starts with an issue** - Before implementing features, fixes, or changes, ensure there's a GitHub issue tracking the work
- **Reference issues in commits** - Use `Fixes #123` or `Refs #123` in commit messages to link work to issues
- **Check existing issues first** - Before creating new issues, search for existing ones that may cover the same topic

### Pull Requests

- **Feature branches** - Create branches from `main` for all changes (e.g., `feat/chartjs-plugin`, `fix/template-error`)
- **PR descriptions** - Include a summary, link to the issue, and list of changes
- **CI must pass** - All PRs require passing CI (lint, test, build) before merge
- **Keep PRs focused** - One feature or fix per PR when possible

### Releases

- **Semantic Versioning** - This project follows [semver](https://semver.org/): `MAJOR.MINOR.PATCH`
- **CRITICAL: Do NOT bump MAJOR version** - Only the project owner can authorize major version bumps. This includes any breaking changes.
- **Minor versions** - New features that are backward compatible
- **Patch versions** - Bug fixes and minor improvements
- **Releases are created via GitHub Releases** - Use `gh release create` or the GitHub UI
- **GoReleaser handles binaries** - Release automation is configured in `.goreleaser.yml`

### Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feat/<description>` | `feat/chartjs-plugin` |
| Bug fix | `fix/<description>` | `fix/feed-pagination` |
| Docs | `docs/<description>` | `docs/plugin-guide` |
| Refactor | `refactor/<description>` | `refactor/lifecycle` |
| Chore | `chore/<description>` | `chore/update-deps` |

### Useful Commands

```bash
# Issues
gh issue list                        # List open issues
gh issue create                      # Create new issue
gh issue view 123                    # View issue details

# Pull Requests
gh pr create                         # Create PR from current branch
gh pr list                           # List open PRs
gh pr checks                         # View CI status
gh pr merge                          # Merge PR (after approval)

# Releases
gh release list                      # List releases
gh release create v0.2.0 --generate-notes  # Create release with auto-generated notes
```

---

## Project Overview

**markata-go** is a static site generator written in Go 1.27+. It processes Markdown files with frontmatter through a plugin-based lifecycle system to generate HTML sites.

Key dependencies: cobra (CLI), pongo2 (templating), goldmark (Markdown), fsnotify (live reload)

## Spec-Driven Development

**This is a spec-driven project.** All changes must follow this workflow:

1. **Spec First** - Before implementing any feature or change, update or create the relevant specification in `spec/`. The spec defines the expected behavior, data models, and interfaces.

2. **Review Spec** - Ensure the spec is complete and covers edge cases, error handling, and integration with existing components.

3. **Implement from Spec** - Write code that conforms to the specification. The spec is the source of truth.

4. **Test Against Spec** - Tests should verify behavior defined in the spec, not implementation details.

5. **Document for Users** - Every feature MUST have corresponding user documentation in `docs/`. See Documentation Requirements below.

**Spec directory structure:**
```
spec/
├── README.md         # Spec overview and index
├── spec/
│   ├── SPEC.md           # Core architecture overview
│   ├── CONFIG.md         # Configuration system
│   ├── LIFECYCLE.md      # Build stages and hooks
│   ├── PLUGINS.md        # Plugin development guide
│   ├── DATA_MODEL.md     # Post, Config, Feed schemas
│   ├── TEMPLATES.md      # Template system
│   ├── THEMES.md         # Theme and palette system
│   ├── FEEDS.md          # Feed generation
│   └── ...
```

**When making changes:**
- For new features: Create or update the relevant spec file first
- For bug fixes: Check if the spec needs clarification
- For refactoring: Ensure changes still conform to the spec
- If spec and code disagree: The spec should be updated first, then code follows

## Documentation Requirements

**CRITICAL: Every feature must be documented.** The documentation lives in `docs/` and is built as part of this site.

### Documentation Workflow

1. **Spec defines behavior** - Technical specification in `spec/spec/*.md`
2. **Docs explain usage** - User-friendly guides in `docs/guides/*.md`
3. **Docs are part of the site** - All docs are Markdown files processed by markata-go itself
4. **Bundled agent skill stays current** - The bundled site skill in `pkg/agentskill/bundle/markata-go-site/` is a user-facing artifact for coding agents working in markata-go site repositories

### Documentation Structure

```
docs/
├── index.md              # Documentation home
├── getting-started.md    # Quick start guide
├── quickstart.md         # 5-minute tutorial
├── guides/
│   ├── configuration.md  # Config reference
│   ├── themes.md         # Themes and palettes
│   ├── templates.md      # Template customization
│   ├── feeds.md          # Feed system
│   ├── frontmatter.md    # Frontmatter fields
│   ├── markdown.md       # Markdown features
│   └── ...
├── reference/
│   ├── cli.md            # CLI commands
│   └── plugins.md        # Plugin reference
└── troubleshooting.md    # Common issues
```

### Documentation Standards

**Every doc file must have frontmatter:**
```yaml
---
title: "Human-Readable Title"
description: "One-line description for SEO and feeds"
date: 2024-01-15
published: true
tags:
  - documentation
  - relevant-topic
---
```

**Documentation principles:**
- **User-focused** - Explain how to USE features, not how they're implemented
- **Examples first** - Show working examples before explaining details
- **Progressive disclosure** - Simple cases first, advanced options later
- **Complete** - Cover all options, edge cases, and common mistakes
- **Tested** - All code examples should actually work

### When to Create Documentation

| Change Type | Documentation Required |
|-------------|----------------------|
| New feature | New guide or section in existing guide |
| New CLI command | Add to `docs/reference/cli.md` |
| New config option | Add to `docs/guides/configuration.md` |
| New plugin | Add to `docs/reference/plugins.md` |
| Bug fix | Update docs if behavior was unclear |
| Breaking change | Update all affected docs + migration guide |

### Bundled Skill Maintenance

**The bundled site skill is part of the product surface.** Any change that affects how an AI agent should work in a markata-go site repository MUST review the bundled skill and update it when needed.

Common triggers include changes to:

- configuration behavior or recommended config patterns
- frontmatter semantics or content creation workflows
- CLI usage, flags, or command output relied on by agents
- build and deployment workflows
- performance guidance such as `--fast`, caching, and benchmark workflows
- theme, palette, or aesthetic workflows
- template, layout, or template-context behavior
- plugin authoring constraints or extension workflows

When a change affects site-author workflows, the work is not complete until the bundled skill is updated or the PR explicitly documents why no skill update is needed.

### Documentation Checklist

Before completing any feature, verify:
- [ ] Spec is updated in `spec/spec/` (behavior, config options, data model changes)
- [ ] User documentation exists in `docs/`
- [ ] Config options are documented in `docs/guides/configuration.md`
- [ ] New plugins are documented in `docs/reference/plugins.md`
- [ ] CLI commands are documented in `docs/reference/cli.md`
- [ ] Examples are provided and tested
- [ ] Cross-references link to related docs
- [ ] Bundled agent skill reviewed and updated if the change affects site-building agents

## Build/Lint/Test Commands

## Agent Gotchas

### Template trees

- This repo has two active template trees:
  - `pkg/themes/default/templates/`
  - `templates/`
- When a template change appears to have no effect, check both trees before assuming the build or cache is wrong.

### Shared template rebuilds

- Incremental builds may not refresh pages affected only by shared templates or global CSS.
- For changes to `base.html`, shared components, or layout CSS, prefer a clean rebuild of the target site:
  - remove `output/`, `cache/`, `.markata-cache/`, and `.markata/`
  - then rebuild from the site directory

### Feed page layout

- Feed pages should be treated differently from post pages when debugging layout.
- Before changing feed centering rules, inspect the rendered HTML for sidebars like `content-sidebar` or `feed-sidebar`; those can shift `main` even when the inner feed content is centered.

### Verify rendered output first

- When source changes and browser output disagree, inspect generated HTML and the exact hashed CSS/JS assets linked by the page.
- This codebase can produce stale-looking results from incremental rebuilds even when source files are correct.

### Building

```bash
go build ./cmd/markata-go        # Build the binary
go run ./cmd/markata-go build    # Run without building
```

### Testing

```bash
go test ./...                           # Run all tests
go test -v ./...                        # Verbose output
go test -race ./...                     # Race detection
go test -coverprofile=coverage.out ./...  # Coverage
```

### Cross-platform tests

When changing timing, filesystem, or process tests, check the Windows contract
before pushing:

- Validate measured work with operation counts and outputs; durations may be
  zero within one clock tick. Assert non-negative durations, not positive ones.
- Use `filepath` for filesystem paths and native-path map keys; convert to
  slash-separated paths only at URL or explicitly normalized key boundaries.
- Capture file identity with an opened file's `Stat`, then close the handle
  before replacement. Windows `os.Stat` may defer identity lookup until
  `os.SameFile`, after the original path has changed.
- Check missing-path errors before `syscall.ENOTDIR` special cases: Windows
  aliases `ENOTDIR` to `ERROR_PATH_NOT_FOUND`.
- Test portable permission behavior on Windows and exact POSIX modes on Unix.
  Gate only inherently Unix-specific assertions; retain byte, publication,
  containment, and hard-link safety checks on every supported platform.
- Run affected tests locally, then verify the PR's actual Windows check.
  Cross-compilation verifies compilation, not filesystem or runtime behavior.

### Running a Single Test

```bash
go test -v -run TestParseFrontmatter ./pkg/plugins/    # Specific test
go test -v -run "TestConfig.*" ./pkg/config/           # Pattern match
go test -v ./pkg/filter/...                            # Entire package
go test -v ./tests/...                                 # Integration tests
```

### Linting and Formatting

```bash
go fmt ./...          # Format all code (ALWAYS run before committing)
go vet ./...          # Check for common issues
```

**Lint commands (choose based on context):**

| Command | Time | Use When |
|---------|------|----------|
| `just lint-fast` | ~1s | Quick iteration while coding |
| `just lint-new` | ~2-5s | Before committing (only changed files vs main) |
| `just lint` | ~15-20s | Full validation before PR, uses all CPU cores |
| `just lint-gentle` | ~20-25s | Full lint with reduced CPU (4 cores), good when multitasking |

**Recommended workflow for agents:**
1. Run `just lint-fast` during iteration.
2. Run `just pr-check` before pushing.
3. For editor changes also run `just pr-check --browser`, then inspect latest PR checks.

**Note:** Full lint runs 36 linters across ~120K lines of code. The `--fast` flag runs only linters that don't require type-checking, which is sufficient for catching most issues during development.

## Project Structure

```
cmd/markata-go/       # CLI entry point and commands
pkg/
├── config/           # Configuration loading, parsing, validation
├── models/           # Data models (Post, Config, Feed, errors)
├── filter/           # Filter expression lexer/parser/evaluator
├── lifecycle/        # Build lifecycle manager (9 stages)
├── plugins/          # Built-in plugins (15+)
└── templates/        # Pongo2 template engine wrapper
templates/            # Default HTML templates
tests/                # Integration tests
spec/                 # Specification documents
```

## Code Style Guidelines

### Import Organization

**Three groups separated by blank lines, alphabetically sorted within each:**
```go
import (
    // 1. Standard library
    "errors"
    "fmt"
    "sync"

    // 2. Project-internal
    "github.com/example/markata-go/pkg/models"

    // 3. Third-party
    "github.com/yuin/goldmark"
    "gopkg.in/yaml.v3"
)
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Files | snake_case | `render_markdown.go` |
| Exported | PascalCase | `NewManager()`, `Post` |
| Unexported | camelCase | `runStage()`, `memoryCache` |
| Constructors | `New` prefix | `NewConfig()`, `NewFrontmatterParseError()` |
| Sentinel errors | `var Err` prefix | `var ErrConfigNotFound` |
| Error types | `Error` suffix | `FrontmatterParseError` |
| Interfaces | Capability name | `Plugin`, `Cache`, `ConfigurePlugin` |

### Error Handling

**Sentinel errors** - for expected, checkable conditions:
```go
var ErrConfigNotFound = errors.New("no configuration file found")
```

**Custom error types** - with context and `Unwrap()` for wrapping:
```go
type FrontmatterParseError struct {
    Path    string
    Message string
    Err     error  // wrapped error
}

func (e *FrontmatterParseError) Error() string {
    return fmt.Sprintf("frontmatter parse error in %s: %s", e.Path, e.Message)
}

func (e *FrontmatterParseError) Unwrap() error { return e.Err }

func NewFrontmatterParseError(path, message string, err error) *FrontmatterParseError {
    return &FrontmatterParseError{Path: path, Message: message, Err: err}
}
```

**Error wrapping** - always use `%w` verb:
```go
return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
```

**Error checking** - use `errors.Is()` and `errors.As()`:
```go
if errors.Is(err, ErrConfigNotFound) { ... }

var parseErr *FrontmatterParseError
if errors.As(err, &parseErr) { ... }
```

### Struct Patterns

**Struct tags** - support all three serialization formats:
```go
type Post struct {
    Path  string  `json:"path" yaml:"path" toml:"path"`
    Title *string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
}
```

**Optional fields** - use pointers for optional scalars:
```go
Title       *string    `json:"title,omitempty"`
PublishDate *time.Time `json:"publish_date,omitempty"`
```

### Interface Patterns

**Small, composable interfaces:**
```go
type Plugin interface {
    Name() string
}

type ConfigurePlugin interface {
    Plugin
    Configure(m *Manager) error
}
```

**Compile-time interface verification:**
```go
var (
    _ lifecycle.Plugin          = (*RenderMarkdownPlugin)(nil)
    _ lifecycle.ConfigurePlugin = (*RenderMarkdownPlugin)(nil)
)
```

### Concurrency Patterns

**Mutex for shared state** - use `sync.RWMutex`:
```go
type Manager struct {
    mu    sync.RWMutex
    posts []*models.Post
}

func (m *Manager) Posts() []*models.Post {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return append([]*models.Post{}, m.posts...)  // Return copy
}
```

**Semaphore for limiting concurrency:**
```go
semaphore := make(chan struct{}, concurrency)
for _, item := range items {
    go func(i Item) {
        semaphore <- struct{}{}        // Acquire
        defer func() { <-semaphore }() // Release
        process(i)
    }(item)
}
```

### Documentation

**Package docs** in `doc.go` with `# Heading` sections:
```go
// Package templates provides a Jinja2-like template engine.
//
// # Template Engine
//
// The Engine type manages template loading:
//
//  engine, err := templates.NewEngine("templates/")
```

**Function docs** - describe behavior, parameters, return values:
```go
// ExtractFrontmatter splits content into frontmatter and body.
// Returns ErrInvalidFrontmatter if the frontmatter is malformed.
```

### Testing Patterns

**Test naming:** `Test{Type}_{Scenario}`
```go
func TestConfig_DefaultOutputDir(t *testing.T)
func TestParseFrontmatter_MissingDelimiter(t *testing.T)
```

**Table-driven tests:**
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"basic", "input", "expected", false},
    {"empty", "", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Process(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
        }
        if got != tt.want {
            t.Errorf("got %q, want %q", got, tt.want)
        }
    })
}
```

**Test helpers** - always call `t.Helper()`:
```go
func newTestSite(t *testing.T) *testSite {
    t.Helper()
    dir := t.TempDir()
    // setup...
    return &testSite{dir: dir, t: t}
}
```

### File Organization

- `doc.go` - Package documentation
- `errors.go` - Error types and sentinels
- One primary type per file for large types
- Tests adjacent to source: `foo.go` -> `foo_test.go`

## Lifecycle Stages

Build runs through 9 stages: Configure -> Validate -> Glob -> Load -> Transform -> Render -> Collect -> Write -> Cleanup

## Plugin Development

Implement `lifecycle.Plugin` plus stage interfaces (`ConfigurePlugin`, `RenderPlugin`, etc.). Use `PriorityPlugin` for execution order control.

## Configuration

Primary: TOML (`markata-go.toml`). Also: YAML, JSON. Env override: `MARKATA_GO_` prefix.

## Performance Optimization

### Build Performance Context

The site at `waylonwalker.com-markata-go-migration` (3341 posts, 396 feeds, 100 blogroll feeds) is the primary performance benchmark target. Build types:

- **Cold build** (~105s): No `.markata-cache/` directory. I/O bound (9% CPU utilization). Clear caches with `rm -rf .markata.cache .markata-cache output cache .markata`.
- **Warm build** (~4s steady-state): Cache populated. CPU bound. The first warm build after a cold build may show cache misses for plugin-level caches due to hash changes from config detection; the **second warm build** is the true steady-state number.
- **`--fast` build**: Skips JS/CSS minification and CSS purging for faster dev iteration.

### Profiling Methodology

```bash
# Benchmark warm builds (run build 3 times, use 2nd warm as steady-state):
rm -rf .markata-cache output
time go run ./cmd/markata-go build -c /path/to/config   # Cold
time go run ./cmd/markata-go build -c /path/to/config   # Warm 1 (cache priming)
time go run ./cmd/markata-go build -c /path/to/config   # Warm 2 (steady-state)

# CPU/memory profiling:
go tool pprof -http=:8080 /tmp/cpu.prof
go tool pprof -top /tmp/mem.prof
```

### Benchmark Prerequisites For Real Sites

- If the target site uses encrypted/private posts, benchmark commands must either:
  - provide the required `MARKATA_GO_ENCRYPTION_KEY_*` environment variables, or
  - explicitly disable encryption for both baseline and candidate runs with `MARKATA_GO_ENCRYPTION_ENABLED=false`
- Do not assume the benchmark shell has access to a developer's local secrets or `.env` loader.
- When encryption is disabled to unblock local measurements, call that out in the benchmark report so comparisons are interpreted correctly.

The lifecycle timing instrumentation in `hooks.go` logs any plugin taking >50ms. This is a permanent feature -- use these logs to identify bottlenecks.

### Plugin Caching Pattern

The build cache (`pkg/buildcache/cache.go`) stores per-post results keyed by content hashes. When optimizing a plugin for warm builds, follow this two-phase pattern:

1. **Add fields to `PostCache` struct** in `cache.go`:
   - `XxxHash string` -- the content hash used to detect changes
   - `XxxHTML string` or `XxxContent string` -- the cached output

2. **Add cache methods** to `buildcache.Cache`:
   - `GetCachedXxx(slug string) (hash, html string, ok bool)` -- reads from in-memory `sync.Map` first, falls back to disk cache
   - `CacheXxx(slug, hash, html string)` -- writes through to both in-memory `sync.Map` and disk cache

3. **Rewrite the plugin's `Render()`/`Transform()`** with two phases:
   - **Phase 1 (restore)**: Iterate all posts, compute current hash (e.g., `buildcache.ContentHash(post.ArticleHTML)`), compare to cached hash. If match, restore cached result directly. Track cache hits/misses.
   - **Phase 2 (process)**: Process only changed posts concurrently using `ProcessPostsSliceConcurrently`. After processing, call `CacheXxx()` to store results.

**Hash inputs** should capture everything that affects the plugin's output. For example:
- `render_markdown`: raw markdown body
- `glossary`: `ArticleHTML` + hash of glossary terms (via `computeTermsHash()`)
- `link_avatars`: `ArticleHTML` (contains links that get avatar icons)
- `embeds`: raw markdown `Content` (embeds are transformed pre-render)

### Regex Optimization

All regex patterns should be compiled once at package level, not inside functions:

```go
// Good: compiled once at init
var headingPattern = regexp.MustCompile(`^#+\s+(.+)$`)

// Bad: recompiled on every call
func process(s string) {
    re := regexp.MustCompile(`^#+\s+(.+)$`)
}
```

This project has ~30+ hoisted regex patterns across 6 files (`blogroll.go`, `oembed.go`, `critical_css.go`, `mentions.go`, `updater.go`, `avatar.go`, `lint.go`).

### Concurrency in Plugins

- Use `ProcessPostsSliceConcurrently` (worker pool with configurable concurrency) for per-post work
- Use `sync.Pool` for buffer reuse in minification plugins
- Use `sync.Map` for concurrent caches (icon lookups, XML tag caches) that don't need eviction
- Use semaphore pattern (`chan struct{}`) for bounded concurrency in I/O-heavy work

### Performance Anti-Patterns to Avoid

- **O(N^2) feed membership**: Never scan all feeds for every post. Pre-compute a `map[slug][]feedName` index (see `templates.go` `computeFeedMembershipHash`).
- **Double cache checks**: Don't call `ShouldRebuild()` then separately read from cache. Use `ShouldRebuildBatch()` or combine check+restore in one pass.
- **Per-item regex scanning**: Don't re-extract data with regex when you already have it. Use single-pass extraction (see `link_collector.go` `extractHrefsAndText()`).
- **Allocations in hot loops**: Reuse buffers, pre-allocate slices to known capacity, use `sync.Pool` for temporary objects.

### Current Warm Build Budget (~4s total)

| Plugin | Time | Notes |
|--------|------|-------|
| glob | ~600ms | Filesystem scanning, likely irreducible |
| build_cache cleanup | ~400ms | Writing cache to disk |
| js_minify | ~380ms | 66 JS files, already parallelized |
| css_minify | ~330ms | 68 CSS files, already parallelized |
| configure/build_cache | ~280ms | Loading cache from disk |
| collect/auto_feeds | ~240ms | |
| transform/mentions | ~185ms | 32s cold -- potential optimization target |
| collect/feeds | ~150ms | |
| collect/blogroll | ~140ms | |
| render_markdown | ~120ms | Cached, only processes changed posts |
| Everything else | <100ms each | Cached or fast |

### Dependencies Note

The following dependencies may appear unused but are actively used elsewhere in the codebase. Do NOT remove or flag them:
- `chromedp` -- used for critical CSS extraction
- `charmbracelet/bubbletea` -- used for TUI components
- `steam` plugin -- generates steam game pages
- Theme switcher, keyboard shortcuts -- active UI features
