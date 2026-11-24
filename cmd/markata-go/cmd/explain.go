package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// topicServe and topicTemplates are constants for built-in explain topics.
const topicTemplates = "templates"
const topicServe = "serve"

// explainCmd represents the explain command.
var explainCmd = &cobra.Command{
	Use:   "explain [topic]",
	Short: "Show detailed information about markata-go for AI agents",
	Long: `Outputs comprehensive context about markata-go commands and concepts.

This command is designed to provide AI coding agents with the context they need
to work effectively with markata-go projects.

Topics:
  (none)     General overview of markata-go
  build      The build command and build process
  serve      The development server
  new        Creating new content
  content    Finding, creating, and editing site content
  init       Initializing projects
  config     Configuration system
  agents     Agent skill installation and usage
  plugins    Plugin system and development
  lifecycle  Build lifecycle stages
  templates  Template system
  feeds      Feed generation system

Example:
  markata-go explain              # General overview
  markata-go explain build        # Explain build command
  markata-go explain plugins      # Explain plugin system
  markata-go explain agents       # Explain agent integrations
  markata-go explain content      # Explain the agent content workflow`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: explainValidArgs,
	RunE:              runExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

// explainValidArgs provides tab completion for explain topics.
func explainValidArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"build\tThe build command and process",
		topicServe + "\tThe development server",
		"new\tCreating new content",
		"content\tFinding, creating, and editing site content",
		"init\tInitializing projects",
		"config\tConfiguration system",
		"agents\tAgent integrations and skills",
		"plugins\tPlugin system",
		"lifecycle\tBuild lifecycle stages",
		topicTemplates + "\tTemplate system",
		"feeds\tFeed generation",
	}, cobra.ShellCompDirectiveNoFileComp
}

func runExplain(cmd *cobra.Command, args []string) error {
	currentCmd = cmd
	topic := ""
	if len(args) > 0 {
		topic = args[0]
	}

	var content string
	switch topic {
	case "":
		content = explainGeneral
	case "build":
		content = explainBuild
	case topicServe:
		content = explainServe
	case "new":
		content = explainNew
	case "content":
		content = explainContent
	case "init":
		content = explainInit
	case "config":
		content = explainConfig
	case "agents":
		content = explainAgents
	case "plugins":
		content = explainPlugins
	case "lifecycle":
		content = explainLifecycle
	case topicTemplates:
		content = explainTemplates
	case "feeds":
		content = explainFeeds
	default:
		return fmt.Errorf("unknown topic %q; run 'markata-go explain --help' to see available topics", topic)
	}

	outln(content)
	return nil
}

const explainGeneral = `# markata-go

A fast, plugin-based static site generator written in Go.

## Overview

markata-go processes Markdown files with YAML frontmatter through a plugin-based
lifecycle system to generate static HTML sites. It features a powerful feed system,
live reload development server, and extensive customization options.

## Quick Start

    # Initialize a new project
    markata-go init

    # Create a new post
    markata-go new "My First Post"

    # Start development server with live reload
    markata-go serve

    # Build for production
    markata-go build --clean

## Key Features

- **Plugin Architecture**: 9-stage lifecycle with 15+ built-in plugins
- **Feed System**: Automatic RSS/Atom/JSON feeds with filtering
- **Live Reload**: Development server with automatic rebuilds
- **Markdown Extensions**: Tables, syntax highlighting, wikilinks, admonitions
- **Theme System**: Color palettes and customizable templates
- **Fast Builds**: Concurrent processing, incremental builds

## Project Structure

    my-site/
    ├── markata-go.toml      # Configuration file
    ├── posts/               # Markdown content
    │   └── my-post.md
    ├── pages/               # Static pages
    ├── templates/           # Custom templates (optional)
    ├── static/              # Static assets
    └── output/              # Generated output (default)

## Key Source Files

    cmd/markata-go/          # CLI entry point
    pkg/config/              # Configuration loading
    pkg/lifecycle/           # Build lifecycle manager
    pkg/plugins/             # Built-in plugins
    pkg/models/              # Data models (Post, Config, Feed)
    pkg/filter/              # Filter expression parser
    pkg/templates/           # Template engine (pongo2)

## Configuration

Primary config file: markata-go.toml (also supports YAML/JSON)

    [markata-go]
    title = "My Site"
    url = "https://example.com"
    output_dir = "output"

    [markata-go.glob]
    patterns = ["posts/**/*.md"]

## Common Commands

    markata-go build              # Build the site
    markata-go serve              # Development server
    markata-go new "Title"        # Create new post
    markata-go config show        # Show resolved config
    markata-go explain <topic>    # Get help on a topic

## Getting Help

    markata-go --help             # List all commands
    markata-go <command> --help   # Command-specific help
    markata-go explain <topic>    # Detailed explanations

## Documentation

- Docs: docs/ directory in the repo
- Specs: spec/spec/ directory for technical specifications
- Config reference: docs/guides/configuration.md
- CLI reference: docs/reference/cli.md`

const explainBuild = `# markata-go build

Build the static site by processing all content through the plugin lifecycle.

## Usage

    markata-go build [flags]

## Flags

    --clean       Remove output directory and build cache before building
    --clean-all   Also remove external plugin caches before building
    --dry-run     Show what would be built without writing files
    -v, --verbose Enable verbose logging
    -o, --output  Override output directory
    -c, --config  Specify config file path

## Examples

    # Standard build
    markata-go build

    # Clean build (removes output and build cache)
    markata-go build --clean

    # Clean all local and external plugin caches
    markata-go build --clean-all

    # Preview what would be built
    markata-go build --dry-run

    # Build with verbose output
    markata-go build -v

    # Build to custom directory
    markata-go build -o dist

## Build Process (9 Stages)

1. **Configure** - Load and merge configuration from files and environment
2. **Validate** - Validate configuration settings
3. **Glob** - Discover content files matching patterns
4. **Load** - Parse markdown files and extract frontmatter
5. **Transform** - Pre-render processing (wikilinks, jinja-md, descriptions)
6. **Render** - Convert markdown to HTML using templates
7. **Collect** - Build feeds, archives, and collections
8. **Write** - Output all files to disk
9. **Cleanup** - Release resources

## Key Source Files

    cmd/markata-go/cmd/build.go   # Build command implementation
    pkg/lifecycle/manager.go      # Lifecycle orchestration
    pkg/plugins/                  # Individual stage plugins

## Related Configuration

    [markata-go]
    output_dir = "output"         # Where to write output

    [markata-go.glob]
    patterns = ["posts/**/*.md"]  # Files to process

## Exit Codes

    0  - Build completed successfully
    1  - Build failed (configuration, plugin, or I/O error)
    2  - No content files found

## Common Issues

1. **No content files found**
   - Check glob.patterns in config matches your file structure
   - Ensure files have .md extension

2. **Template not found**
   - Verify templates_dir points to correct location
   - Check template names in frontmatter match actual files

3. **Build errors in specific posts**
   - Run with -v to see which file failed
   - Check frontmatter YAML syntax
   - Verify date formats (YYYY-MM-DD)

## Environment Variables

    MARKATA_GO_OUTPUT_DIR    Override output directory
    MARKATA_GO_URL           Override site URL (for production builds)`

const explainServe = `# markata-go serve

Start a development server with live reload support.

## Usage

    markata-go serve [flags]

## Flags

    -p, --port     Port to listen on (default: 8000)
    --host         Host address to bind (default: localhost)
    --no-watch     Disable file watching and auto-rebuild
    -v, --verbose  Enable verbose logging

## Examples

    # Serve on default port (localhost:8000)
    markata-go serve

    # Serve on custom port
    markata-go serve -p 3000

    # Serve with merged config overrides
    markata-go serve -m fast.toml

    # Bind to all interfaces (accessible from network)
    markata-go serve --host 0.0.0.0

    # Serve without auto-rebuild
    markata-go serve --no-watch

## Features

- **Automatic Rebuild**: File changes trigger instant rebuilds
- **Live Reload**: Connected browsers refresh automatically
- **Static Serving**: Serves the output directory
- **MIME Types**: Correct content types for all file formats

## Development Workflow

    # Terminal: Start dev server
    markata-go serve -v

    # Make changes to posts/*.md
    # Browser automatically refreshes

## Key Source Files

    cmd/markata-go/cmd/serve.go   # Serve command implementation
    pkg/server/                   # HTTP server and file watcher

## Related Configuration

    [markata-go]
    output_dir = "output"         # Directory to serve

## Common Issues

1. **Port already in use**
   - Use -p to specify different port
   - Kill existing process on port 8000

2. **Changes not detected**
   - Check if watching is disabled (--no-watch)
   - Verify file is in watched directory
   - Check for file system permission issues

3. **Live reload not working**
   - Check browser console for WebSocket errors
   - Ensure JavaScript is enabled
   - Try hard refresh (Ctrl+Shift+R)`

const explainNew = `# markata-go new

Create a new content file with frontmatter template.

## Usage

    markata-go new [title] [flags]

## Arguments

    title    The title of the new post (prompted if not provided)

## Flags

    --template, -t   Content template to use (default: post)
    --list, -l       List available templates
    --dir            Directory to create post in (default: template's default)
    --draft          Create as draft (default: false)
    --tags           Comma-separated list of tags
    --plain          Use plain text prompts instead of TUI wizard

## Interactive Wizard

When called without arguments, launches a TUI wizard (charmbracelet/huh) that
guides you through template selection, title, directory, tags, privacy, and
author configuration. Falls back to plain text prompts when stdin is not a TTY
or when --plain is passed.

## Examples

    # Launch interactive TUI wizard
    markata-go new

    # Create with title (skips wizard)
    markata-go new "My First Post"

    # Create in specific directory
    markata-go new "Hello World" --dir blog

    # Create as a draft
    markata-go new "Work in Progress" --draft

    # Create with tags
    markata-go new "Go Tutorial" --tags "go,tutorial,programming"

    # Use plain text prompts instead of TUI
    markata-go new --plain

## Generated File

The command creates a markdown file with this template:

    ---
    title: "My First Post"
    slug: "my-first-post"
    date: 2024-01-15
    draft: false
    published: true
    tags: []
    description: ""
    ---

    # My First Post

    Write your content here...

## Slug Generation

The slug is auto-generated from the title:
- Converts to lowercase
- Replaces spaces with hyphens
- Removes special characters

Examples:
    "My First Post"      → my-first-post
    "Hello, World!"      → hello-world
    "Post #1: The Start" → post-1-the-start

## Key Source Files

    cmd/markata-go/cmd/new.go      # New command implementation
    cmd/markata-go/cmd/new_huh.go  # TUI wizard implementation

## Related Configuration

    [markata-go.glob]
    patterns = ["pages/**/*.md", "posts/**/*.md"]  # Ensure new posts match patterns

    [content_templates.placement]
    post = "blog"     # Override default directory per template type

## Common Issues

1. **Post not appearing in build**
   - New posts default to draft: false, published: true
   - Verify file matches glob patterns
   - If created with --draft, set draft: false when ready

2. **Date format errors**
   - Use YYYY-MM-DD format
   - Don't include time unless needed`

const explainContent = `# markata-go Content Workflow

Use this workflow to inspect, find, create, and update content. It is designed
for coding agents and works from any directory when you select a site.

## Select A Site

Use --site-dir for a one-off command, or set MARKATA_GO_SITE_DIR for a shell
alias or wrapper. The flag takes precedence over the environment variable.

    markata-go --site-dir ~/sites/blog list posts
    MARKATA_GO_SITE_DIR=~/notes/work markata-go search "meeting notes"

When a site directory is explicitly selected, source paths returned by list and
search are absolute and can be edited directly.

## Inspect And Find Content

    # Inspect the resolved configuration and content inventory
    markata-go --site-dir ~/sites/blog config show
    markata-go --site-dir ~/sites/blog list posts --format json

    # Find Markdown source files by content, title, description, or tag
    markata-go --site-dir ~/sites/blog search "release process" --format path
    markata-go --site-dir ~/sites/blog list posts --filter "'go' in tags" --format path

Edit a returned Markdown path directly. Markata-go does not provide a separate
non-interactive edit command. Human users can also run markata-go tui and press
e to open the selected post in $EDITOR.

## Create Content

    # Inspect site-provided and built-in content templates
    markata-go --site-dir ~/sites/blog new --list

    # Create content without prompts, suitable for automation
    markata-go --site-dir ~/sites/blog new "Release notes" --template post --tags "release,go" --no-input

New content uses the selected site's configuration and content templates. With
an explicit site directory, the created path is printed as an absolute path.

## Visual Sequences (Roads)

For timelines, roadmaps, step-by-step paths, or release histories, wrap a list
of links in a road container. The default theme draws it as the reader scrolls.

    ::: road timeline

    - [Kickoff](/notes/kickoff/) Jan 2026
    - [Launch](/notes/launch/) Jun 2026
    :::

The text after each link is the label. Styles: winding, timeline, steps, snake,
cycle, constellation, git, and editorial. For editorial, use a heading per
entry, with Markdown under it. For git, write the branch after each link
(feat/x, fix/y from feat/x, main <- feat/x). See docs/guides/roads.md.

## Validate Changes

    markata-go --site-dir ~/sites/blog lint
    markata-go --site-dir ~/sites/blog build --fast

Run a full build before publishing or deployment.

## Install Agent Guidance Globally

If agents often begin outside a site repository, install the bundled skill in
their global skill directory so they can discover this workflow:

    markata-go agent install --agent opencode --global`

const explainInit = `# markata-go init

Initialize a new markata-go project with interactive setup.

## Usage

    markata-go init [flags]

## Flags

    --force    Overwrite existing files

## Examples

    # Interactive project setup
    markata-go init

    # Overwrite existing config
    markata-go init --force

## Interactive Flow

    $ markata-go init

    Welcome to markata-go!

    Site title [My Site]: My Awesome Blog
    Description [A site built with markata-go]: A blog about things
    Author []: Your Name
    URL [https://example.com]: https://myblog.com

    Creating project structure...
      ✓ Created posts/
      ✓ Created static/
      ✓ Created markata-go.toml

## What Gets Created

1. **markata-go.toml** - Configuration with your site settings
2. **posts/** - Directory for blog posts
3. **static/** - Directory for static assets
4. **(Optional) First post** - Starter markdown file

## Generated Configuration

    [markata-go]
    title = "My Site"
    description = "A site built with markata-go"
    url = "https://example.com"
    author = "Your Name"

    output_dir = "output"
    templates_dir = "templates"
    assets_dir = "static"

    [markata-go.glob]
    patterns = ["posts/**/*.md", "pages/**/*.md"]

## Key Source Files

    cmd/markata-go/cmd/init.go    # Init command implementation

## Next Steps After Init

    # Create your first post
    markata-go new "Hello World"

    # Start development server
    markata-go serve

    # View your site at http://localhost:8000`

const explainConfig = `# markata-go Configuration

Configuration system for markata-go sites.

## Config File Locations

Unless --config is provided, markata-go searches for config files in this order:
1. markata-go.toml
2. markata-go.yaml
3. markata-go.yml
4. markata-go.json
5. ~/.config/markata-go/config.toml

If none is found, commands that load the full site configuration use built-in
defaults plus MARKATA_GO_* environment overrides. Commands that read or edit a
config file directly, such as config get and config set, require a file.

## Basic Configuration

    [markata-go]
    title = "My Site"
    description = "A site built with markata-go"
    url = "https://example.com"
    author = "Your Name"
    output_dir = "output"

## Content Discovery

    [markata-go.glob]
    patterns = [
        "posts/**/*.md",
        "pages/**/*.md"
    ]

## Feed Configuration

    [[markata-go.feeds]]
    name = "all"
    filter = "published == true"
    sort = "-date"

    [[markata-go.feeds]]
    name = "python"
    filter = "'python' in tags"
    template = "feed.html"

## Markdown Extensions

    [markata-go.markdown]
    extensions = [
        "tables",
        "strikethrough",
        "autolinks",
        "tasklist"
    ]

## Theme Configuration

    [markata-go.theme]
    palette = "dracula"

    [markata-go.theme.colors]
    primary = "#bd93f9"
    background = "#282a36"

## Config Commands

    # Show resolved configuration
    markata-go config show

    # Show as JSON
    markata-go config show --json

    # Get specific value
    markata-go config get output_dir
    markata-go config get feeds.defaults.items_per_page

    # Validate configuration
    markata-go config validate

    # Create new config file
    markata-go config init

## Environment Variables

Override config with MARKATA_GO_ prefix:

    MARKATA_GO_OUTPUT_DIR=dist markata-go build
    MARKATA_GO_URL=https://staging.example.com markata-go build

## Override Precedence (highest first)

1. Command-line flags
2. Environment variables
3. Configuration file
4. Built-in defaults

## Key Source Files

    pkg/config/config.go         # Configuration loading
    pkg/config/defaults.go       # Default values
    cmd/markata-go/cmd/config.go # Config commands

## Common Issues

1. **Config not found**
   - Check the local names above and the global fallback
   - Use -c to specify a custom path
   - Run markata-go config init to create a local file before using config get/set

2. **Invalid TOML syntax**
   - Use online TOML validator
   - Check for missing quotes around strings
   - Verify array syntax uses brackets []

3. **Feed filter errors**
   - Check filter expression syntax
   - Verify field names match frontmatter
   - Use quotes around string values in filters`

const explainAgents = `# markata-go Agent Integrations

Markata-go can install a bundled, topic-based skill into a site repository so AI
coding agents can work with markata-go sites more effectively.

## Install the Bundled Skill

List supported agent ids and install paths:

    markata-go agent list-agents

    markata-go agent install

By default this installs the skill into the current agent's project layout when
markata-go can detect one from the environment. Otherwise it falls back to the
portable universal layout:

    .agents/skills/markata-go-site/

To install for Claude Code explicitly:

    markata-go agent install --agent claude-code

This installs to:

    .claude/skills/markata-go-site/

To install into OpenCode's global skill directory instead:

    markata-go agent install --agent opencode --global

This installs to:

    ~/.config/opencode/skills/markata-go-site/

## Skill Layout

The installed skill is intentionally split into focused topic files:

- configuration
- writing/frontmatter
- cli usage
- build and deployment
- faster builds
- theme creation
- template creation and management
- plugin creation

The entrypoint is ` + "`SKILL.md`" + ` and the detailed guidance lives under ` + "`topics/`" + `.
Reference material lives under ` + "`reference/`" + `, starter files live under ` + "`examples/`" + `, and bundled regression prompts live under ` + "`evals/`" + `.

## Why a Command Group?

The ` + "`agent`" + ` command group is intentionally generic so future subcommands can add:

- alternate install/export targets
- tool-specific packaging helpers
- MCP-oriented integrations

without changing the underlying bundled skill structure.

## Useful Related Commands

    markata-go explain config
    markata-go explain templates
    markata-go explain plugins
    markata-go config show
    markata-go list posts`

const explainPlugins = `# markata-go Plugin System

Plugins extend markata-go's functionality at specific lifecycle stages.

## Plugin Architecture

Plugins implement interfaces for lifecycle stages they participate in:

    type Plugin interface {
        Name() string
    }

    type ConfigurePlugin interface {
        Plugin
        Configure(m *Manager) error
    }

    type RenderPlugin interface {
        Plugin
        Render(m *Manager) error
    }

## Built-in Plugins

| Plugin | Stage | Purpose |
|--------|-------|---------|
| glob | Glob | Discover content files |
| frontmatter | Load | Parse YAML frontmatter |
| auto_title | Load | Generate title from H1 |
| auto_description | Transform | Generate descriptions |
| jinja_md | Transform | Jinja templates in markdown |
| wikilinks | Transform | [[wikilink]] expansion |
| render_markdown | Render | Markdown to HTML |
| toc | Render | Table of contents |
| feeds | Collect | RSS/Atom/JSON feeds |
| prev_next | Collect | Navigation links |
| copy_assets | Write | Static file copying |
| write_posts | Write | HTML file output |
| sitemap | Write | Sitemap generation |

## Plugin Priority

Plugins can specify execution priority within their stage:

    type PriorityPlugin interface {
        Priority() int  // Lower = runs first
    }

Default priority is 100. Use lower values to run early, higher to run late.

## Plugin Configuration

Plugins read from config sections:

    [markata-go.feeds]
    defaults.items_per_page = 10

    [markata-go.markdown]
    extensions = ["tables", "strikethrough"]

## Plugin Development (Go)

    package myplugin

    import "github.com/WaylonWalker/markata/pkg/lifecycle"

    type MyPlugin struct{}

    func (p *MyPlugin) Name() string { return "my-plugin" }

    func (p *MyPlugin) Transform(m *lifecycle.Manager) error {
        for _, post := range m.Posts() {
            // Process post
        }
        return nil
    }

    // Verify interface compliance
    var _ lifecycle.TransformPlugin = (*MyPlugin)(nil)

## Key Source Files

    pkg/lifecycle/lifecycle.go    # Plugin interfaces
    pkg/lifecycle/manager.go      # Plugin registration
    pkg/plugins/                  # Built-in plugins

    # Individual plugins:
    pkg/plugins/glob.go
    pkg/plugins/frontmatter.go
    pkg/plugins/render_markdown.go
    pkg/plugins/feeds.go

## Plugin Best Practices

1. **Check skip flag**: Honor posts marked for skipping
2. **Handle errors gracefully**: Return errors, don't panic
3. **Use priority sparingly**: Most plugins should use default
4. **Cache expensive operations**: Use content hashes for cache keys
5. **Log appropriately**: Use verbose flag for debug output`

const explainLifecycle = `# markata-go Build Lifecycle

The build process runs through 9 ordered stages.

## Stage Overview

    ┌─────────────────────────────────────────────────┐
    │              CONFIGURATION PHASE                 │
    ├─────────────────────────────────────────────────┤
    │  1. Configure  →  2. Validate                   │
    │     [load]          [check]                     │
    ├─────────────────────────────────────────────────┤
    │                CONTENT PHASE                     │
    ├─────────────────────────────────────────────────┤
    │  3. Glob → 4. Load → 5. Transform → 6. Render  │
    │  [find]   [parse]    [preproc]      [html]     │
    │                                                 │
    │  7. Collect                                     │
    │     [feeds, nav]                                │
    ├─────────────────────────────────────────────────┤
    │                OUTPUT PHASE                      │
    ├─────────────────────────────────────────────────┤
    │  8. Write  →  9. Cleanup                        │
    │     [files]    [resources]                      │
    └─────────────────────────────────────────────────┘

## Stage Details

### 1. Configure
- Load config files (TOML/YAML/JSON)
- Apply environment variable overrides
- Initialize plugins

### 2. Validate
- Validate configuration values
- Check required fields
- Verify paths exist

### 3. Glob
- Find files matching patterns
- Filter by gitignore (optional)
- Build file list

### 4. Load
- Read file contents
- Parse YAML frontmatter
- Create Post objects
- Handle encoding issues

### 5. Transform
- Pre-render content processing
- Expand Jinja templates in markdown
- Process wikilinks
- Generate descriptions
- Calculate reading time

### 6. Render
- Convert markdown to HTML
- Apply syntax highlighting
- Generate table of contents
- Wrap in templates

### 7. Collect
- Build feed collections
- Create prev/next navigation
- Build tag pages
- Generate archives

### 8. Write
- Write HTML files
- Generate RSS/Atom/JSON feeds
- Create sitemap
- Copy static assets

### 9. Cleanup
- Release resources
- Close connections
- Log statistics

## Running Partial Builds

    # Dry run - stops before write
    markata-go build --dry-run

    # Programmatic (Go):
    m.RunTo(lifecycle.StageRender)  // Stop after render

## Key Source Files

    pkg/lifecycle/stage.go       # Stage definitions
    pkg/lifecycle/manager.go     # Stage orchestration
    pkg/plugins/                 # Stage implementations

## Stage Dependencies

Each stage requires all previous stages to complete:
- Configure must complete before Validate
- Glob must complete before Load
- etc.

## Error Handling by Stage

| Stage | Behavior |
|-------|----------|
| Configure | Fatal - cannot proceed |
| Validate | Fatal - invalid config |
| Glob | Warn if no files found |
| Load | Skip invalid files, continue |
| Transform | Skip post on error |
| Render | Skip post on error |
| Collect | Skip collection on error |
| Write | Log error, continue others |
| Cleanup | Log errors, attempt all |`

const explainTemplates = `# markata-go Template System

Templates use Pongo2 (Django/Jinja2-like syntax).

## Template Basics

    <!DOCTYPE html>
    <html>
    <head>
        <title>{{ post.title }}</title>
    </head>
    <body>
        {{ body|safe }}
    </body>
    </html>

## Available Variables

### In Post Templates

    post.title        # Post title
    post.slug         # URL slug
    post.date         # Publication date
    post.tags         # List of tags
    post.description  # Meta description
    post.content      # Raw markdown
    body              # Rendered HTML (use |safe filter)

    config.title      # Site title
    config.url        # Site URL
    config.author     # Site author

    posts             # All posts (for navigation)

### In Feed Templates

    feed.name         # Feed name
    feed.posts        # Posts in this feed
    feed.title        # Feed title
    feed.description  # Feed description

## Template Inheritance

Base template (base.html):

    <!DOCTYPE html>
    <html>
    <head>{% block head %}{% endblock %}</head>
    <body>{% block content %}{% endblock %}</body>
    </html>

Child template (post.html):

    {% extends "base.html" %}

    {% block head %}
    <title>{{ post.title }} | {{ config.title }}</title>
    {% endblock %}

    {% block content %}
    <article>{{ body|safe }}</article>
    {% endblock %}

## Filters

    {{ post.title|escape }}       # HTML escape
    {{ body|safe }}               # Mark as safe HTML
    {{ post.date|date:"Jan 2, 2006" }}  # Format date
    {{ posts|length }}            # Count items
    {{ post.title|truncatechars:50 }}   # Truncate

## Conditionals

    {% if post.draft %}
    <span class="draft">Draft</span>
    {% endif %}

    {% if post.tags %}
    <ul>
    {% for tag in post.tags %}
        <li>{{ tag }}</li>
    {% endfor %}
    </ul>
    {% endif %}

## Loops

    {% for p in posts %}
    <article>
        <h2>{{ p.title }}</h2>
        <time>{{ p.date }}</time>
    </article>
    {% empty %}
    <p>No posts found.</p>
    {% endfor %}

## Template Configuration

    [markata-go]
    templates_dir = "templates"    # Custom templates
    default_template = "post.html" # Default for posts

Frontmatter override:

    ---
    template: custom.html
    ---

## Key Source Files

    pkg/templates/engine.go       # Template engine
    pkg/templates/filters.go      # Custom filters
    templates/                    # Default templates

## Common Issues

1. **Template not found**
   - Check templates_dir in config
   - Verify template file exists
   - Check frontmatter template value

2. **HTML escaping**
   - Use |safe for pre-rendered HTML
   - Don't use |safe for user input

3. **Variable undefined**
   - Check variable spelling
   - Use {% if var %} for optional vars`

const explainFeeds = `# markata-go Feed System

Feeds are collections of posts with filtering, sorting, and pagination.

## Basic Feed Configuration

    [[markata-go.feeds]]
    name = "all"
    filter = "published == true"
    sort = "-date"

    [[markata-go.feeds]]
    name = "python"
    filter = "'python' in tags and published == true"
    template = "tag-feed.html"

## Feed Properties

| Property | Description | Default |
|----------|-------------|---------|
| name | Feed identifier (used in URL) | required |
| filter | Filter expression | "true" |
| sort | Sort field (prefix - for desc) | "-date" |
| template | Template for feed page | "feed.html" |
| items_per_page | Pagination size | 10 |
| output_formats | ["html", "rss", "atom", "json"] | ["html"] |

## Filter Expressions

Filter syntax supports:

    # Boolean fields
    published == true
    draft == false

    # String comparisons
    slug == "about"
    template == "post.html"

    # Contains (for tags)
    'python' in tags
    'tutorial' in tags

    # Logical operators
    published == true and 'python' in tags
    draft == false or template == "page.html"

    # Parentheses for grouping
    (published == true) and ('python' in tags or 'go' in tags)

## Output Formats

    [[markata-go.feeds]]
    name = "blog"
    output_formats = ["html", "rss", "atom", "json"]

Generates:
- /blog/index.html (paginated)
- /blog/feed.xml (RSS 2.0)
- /blog/atom.xml (Atom 1.0)
- /blog/feed.json (JSON Feed)

## Pagination

    [[markata-go.feeds]]
    name = "blog"
    items_per_page = 10

Generates:
- /blog/index.html (page 1)
- /blog/page/2/index.html
- /blog/page/3/index.html
- etc.

## Feed Templates

Access in templates:

    <h1>{{ feed.title }}</h1>
    <p>{{ feed.description }}</p>

    {% for post in feed.posts %}
    <article>
        <h2><a href="{{ post.url }}">{{ post.title }}</a></h2>
        <time>{{ post.date|date:"Jan 2, 2006" }}</time>
    </article>
    {% endfor %}

    {% if feed.has_prev %}
    <a href="{{ feed.prev_url }}">Previous</a>
    {% endif %}
    {% if feed.has_next %}
    <a href="{{ feed.next_url }}">Next</a>
    {% endif %}

## Built-in Feeds

Without explicit config, these feeds are auto-created:
- "all" - All published posts

## Key Source Files

    pkg/plugins/feeds.go         # Feed generation
    pkg/models/feed.go           # Feed data model
    pkg/filter/                  # Filter expression parser

## Common Issues

1. **Empty feed**
   - Check filter expression syntax
   - Verify posts have published: true
   - Run with -v to see filter results

2. **Feed filter errors**
   - Strings need quotes: 'python' in tags
   - Boolean values: true/false (lowercase)
   - Date comparisons need proper format

3. **Pagination not working**
   - Check items_per_page > 0
   - Verify template handles pagination vars`
