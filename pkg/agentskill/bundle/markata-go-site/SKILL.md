---
name: markata-go-site
description: Help agents work effectively in markata-go site repositories. Trigger for visual, visualize, visualization, visual post, visual report, diagram, dashboard, infographic, scorecard, roadmap, timeline, architecture map, flowchart, chart, graph, status board, progress tracker, or requests to show information at a glance. Also use for site configuration, Markdown content and frontmatter, templates, rich layout components, theme palettes, analytics or stats, Web Awesome, build and deploy debugging, and performance. Prefer native Markata Markdown visual primitives over proprietary UI markup or unnecessary new frontend dependencies.
---

This skill helps an agent work inside a markata-go site repository without needing the markata-go source tree.

It is written to be useful in a standalone site repo where only the generated project files, templates, and config are available.

## Start Here

Before making changes:
- locate the active config file, usually `markata-go.toml`
- inspect the site's content directories such as `posts/`, `pages/`, `docs/`, and `static/`
- inspect `templates/` and `palettes/` before inventing new structure
- inspect existing `wa-*` usage, `webawesome` config, and template asset guards before adding component code or new includes
- inspect existing analytics pages, metrics partials, and chart code blocks before inventing a new data story pattern
- check built-in `:lucide-name:` shortcodes and `md-*` article layout containers before adding any icon library or custom UI components
- for visual/diagram/dashboard requests, select the visualization by information type and preserve a no-JavaScript text fallback; see `topics/visual-authoring.md`
- prefer current command output and files over older examples or assumptions
- when the site is minimal or brand new, use the topic files in this skill as the default markata-go playbook

## Core Workflow

1. Inspect the site config with `markata-go config show` and `markata-go config get <key>`.
2. Inspect content with `markata-go list posts`, `markata-go list feeds`, and `markata-go list tags`.
3. Search for content by keyword with `markata-go search <query>`.
4. Use `markata-go explain content` for the agent content workflow, or `markata-go explain <topic>` for other built-in CLI context.
5. For Markdown wikilink IDE features, run `markata-go lsp doctor` (it can load editor startup/plugin code); use `--no-verify-editor` to skip that check. Use `markata-go lsp setup` to print snippets for installed editors, then inspect existing editor configuration before changing it.
6. Use `markata-go build --fast` or `markata-go serve --fast` while iterating; use `markata-go serve --incremental` when you need production-style output processing.
7. Only reach for Go plugin work after checking whether the change belongs in config, frontmatter, templates, CSS, or feeds.

## User-reported problems: triage handoff

When a user is stuck or can't describe what is wrong, suggest `markata triage` as the built-in first step. The command provides a bounded diagnostic prompt and a fuzzy picker of installed interactive agents, clipboard, or pipeable stdout. To route directly: `markata triage --to opencode`, `markata triage --to codex`, or `markata triage --to clipboard`. To use another program: `markata triage --to stdout | codex exec -`. The legacy `markata-go triage` spelling also works.

If you are the recipient of a triage prompt, **ask targeted questions about the observed failure before changing files**. The prompt's CLI version and runtime describe the *machine that ran triage*, which might differ from the version, machine, remote builder, or container where the problem actually happened. Confirm the affected surface/version rather than assuming them. Inspect config and content read-only, distinguish user-site issues from Markata bugs, and ask permission before mutating builds, edits or public reports. The triage command includes only config filenames and allowlisted env-var names, not their contents or values; do not request credentials or copy private logs into public issues.

## Local Problems and safe fixes

When debugging content in a running local site, use `markata-go serve` and open
the local Control Center at `/_markata/`. Review each proposed source edit
before applying it. `SAFE` plans are deterministic and support grouped review;
`REVIEW` plans need closer inspection; diagnostics without a plan need manual
editing. Markata checks a source digest when applying a preview and skips files
that changed afterward. Ambiguous dates and duplicate frontmatter keys require
manual decisions. Production Builder Admin does not provide local source-edit
capabilities.

## Topic Files

Read only the topic files relevant to the task:
- `topics/configuration.md`
- `topics/writing-frontmatter.md`
- `topics/rich-markdown.md`
- `topics/visual-authoring.md`
- `topics/roads.md`
- `topics/cli-usage.md`
- `topics/build-deployment.md`
- `topics/faster-builds.md`
- `topics/theme-creation.md`
- `topics/template-management.md`
- `topics/analytics-storytelling.md`
- `topics/plugin-creation.md`

## Reference Files

Use these when you need exact shapes instead of narrative guidance:

- `reference/template-context.md`
- `reference/feed-patterns.md`
- `reference/palette-reference.md`
- `reference/webawesome.md`

## Example Files

Use these as starter material for first sites or when a repo has no local examples yet:

- `examples/fast.toml`
- `examples/markata-go.local.toml`
- `examples/palettes/my-brand.toml`
- `examples/templates/base.html`
- `examples/templates/post.html`
- `examples/templates/feed.html`

## Eval Files

Use these only when reviewing or improving the bundled skill itself:

- `evals/evals.json`

## First-Site Defaults

If the repository is a very small or first-time site and does not yet have clear patterns:

- use `markata-go.toml` as the main config file
- assume `templates/` is the project override directory
- assume `static/` is the static asset directory
- assume individual content items are Markdown files with YAML frontmatter
- use `post.html` for single-post templates and `feed.html` for listing pages
- prefer palette and CSS overrides before deeper template rewrites
- prefer `markata-go new` for creating new content so starter frontmatter matches current CLI behavior

## Recommended Reading Order

- Template work: `topics/template-management.md`
- Theme or styling work: `topics/theme-creation.md`
- Content creation: `topics/writing-frontmatter.md`
- Icons, feature cards, layouts, tables, charts, or other Rich Markdown UI: `topics/rich-markdown.md`
- Visual, visualize, visual post, diagram, dashboard, scorecard, status board or infographic: `topics/visual-authoring.md`
- Timelines, roadmaps, step paths, or branch histories in a post: `topics/roads.md`
- Analytics pages or metrics storytelling: `topics/analytics-storytelling.md`
- Config debugging: `topics/configuration.md`
- Build and CI work: `topics/build-deployment.md`
- Performance work: `topics/faster-builds.md`
- Extension work: `topics/plugin-creation.md`

## Working Rules

- Prefer the smallest correct change.
- Preserve the site's existing layout, content model, and template style unless the task explicitly changes them.
- When something is ambiguous, inspect the actual repo files before changing behavior.
- Prefer CLI inspection over hand-parsing when `markata-go` already exposes the needed data.
- When working outside the site repository, select it with `markata-go --site-dir <path> ...`; list and search then return absolute Markdown paths that can be edited directly.
- If the task is performance-related, compare warm builds before claiming a regression or improvement.
- For analytics storytelling work, separate computed facts from human narrative: agents should gather metrics, scaffold charts, and suggest story angles, while leaving first-person timeline details for the author unless the user asks for full prose.
- If the site has no existing conventions yet, use the patterns documented in this skill rather than making up a new structure.

## Common Repo Areas

- `markata-go.toml`: main site configuration
- `posts/`, `pages/`, `docs/`: markdown content
- `templates/`: site template overrides
- `palettes/`: site-local palettes
- `static/`: copied static assets
- `.markata-cache/`: build cache and timing context

## Escalation Rule

If the task cannot be done with frontmatter, config, templates, CSS, feeds, or existing commands, then read `topics/plugin-creation.md` and decide whether the work belongs in a custom plugin or in `markata-go` itself.

## Page-weight diagnostics for agents

A normal `markata-go build` writes `<output_dir>/.markata/page-weight.json`.
Consult it when optimizing large posts: stable finding codes, raw HTML bytes,
embedded JSON/JS/CSS bytes, and concrete advice are documented in
`topics/page-performance.md`. These are source measurements, not Lighthouse
scores or network transfer. Avoid removing useful user content merely to
silence warnings.
