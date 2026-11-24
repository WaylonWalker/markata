---
title: "Plugin Development"
description: "Guide to creating custom plugins for markata-go using the 9-stage lifecycle"
date: 2024-01-15
published: true
tags:
  - documentation
  - plugins
  - development
---

## Plugin Development Guide

This guide explains how to create plugins for markata-go, a comprehensive, extendable static site generator written in Go.

## Philosophy

Plugins are the primary extension mechanism in markata-go. A plugin can:

- Hook into any of the 9 lifecycle stages
- Extend configuration with new fields
- Add computed fields to posts
- Control execution order via priorities
- Process content concurrently for performance

## Plugin Architecture Overview

markata-go uses a **Standard 9-stage lifecycle** that provides a good balance between flexibility and simplicity. Each stage is a hook point where plugins can participate by implementing the corresponding interface.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CONFIGURATION PHASE                          │
├─────────────────────────────────────────────────────────────────────┤
│  configure → validate                                               │
│       │           │                                                 │
│       ▼           ▼                                                 │
│  [load config, [validate                                            │
│   init plugins] config]                                             │
├─────────────────────────────────────────────────────────────────────┤
│                         CONTENT PHASE                               │
├─────────────────────────────────────────────────────────────────────┤
│  glob → load → transform → render → collect                         │
│    │      │        │          │         │                           │
│    ▼      ▼        ▼          ▼         ▼                           │
│  [find  [parse  [pre-proc  [markdown [feeds,                        │
│   files] posts]  content]   → HTML]   nav]                          │
├─────────────────────────────────────────────────────────────────────┤
│                         OUTPUT PHASE                                │
├─────────────────────────────────────────────────────────────────────┤
│  write → cleanup                                                    │
│    │        │                                                       │
│    ▼        ▼                                                       │
│  [output  [close                                                    │
│   files]   resources]                                               │
└─────────────────────────────────────────────────────────────────────┘
```

## The 9-Stage Lifecycle

### Stage 1: Configure

**Purpose:** Load configuration and initialize plugin state.

**When:** First stage, runs before any content processing.

**What plugins do:**
- Read configuration options from `m.Config()`
- Initialize internal state (parsers, clients, caches)
- Set up external connections

### Stage 2: Validate

**Purpose:** Validate configuration after all plugins have configured.

**When:** After all plugins have initialized.

**What plugins do:**
- Validate plugin-specific configuration
- Check for required fields
- Return errors for invalid configurations

### Stage 3: Glob

**Purpose:** Discover content files.

**When:** First content phase stage.

**What plugins do:**
- Find files matching patterns
- Filter based on gitignore
- Add files to the manager via `m.SetFiles()`

### Stage 4: Load

**Purpose:** Parse files into post objects.

**When:** After files are discovered.

**What plugins do:**
- Read file contents
- Parse frontmatter
- Create `*models.Post` objects
- Add posts via `m.AddPost()`

### Stage 5: Transform

**Purpose:** Pre-render content processing.

**When:** After content is loaded, before rendering.

**What plugins do:**
- Expand template expressions in content
- Calculate derived fields (reading time, descriptions)
- Process wikilinks, shortcodes
- Set up prev/next links

### Stage 6: Render

**Purpose:** Convert markdown to HTML.

**When:** After transform processing.

**What plugins do:**
- Convert markdown to HTML
- Apply syntax highlighting
- Process admonitions
- Apply templates

### Stage 7: Collect

**Purpose:** Build aggregated content.

**When:** After rendering is complete.

**What plugins do:**
- Build RSS/Atom/JSON feeds
- Generate sitemaps
- Create tag/category pages
- Build navigation structures

### Stage 8: Write

**Purpose:** Write output files to disk.

**When:** After all content is processed.

**What plugins do:**
- Write HTML files
- Write feed files
- Copy static assets
- Generate any final output

### Stage 9: Cleanup

**Purpose:** Release resources.

**When:** Final stage.

**What plugins do:**
- Close database connections
- Flush caches
- Clean up temporary files
- Log final statistics

## Plugin Interfaces

Every plugin must implement the base `Plugin` interface. Additional interfaces determine which stages the plugin participates in.

### Base Interface (Required)

```go
// Plugin is the base interface that all plugins must implement.
type Plugin interface {
    // Name returns the unique name of the plugin.
    Name() string
}
```

### Stage Interfaces

Implement these interfaces to participate in specific stages:

```go
// ConfigurePlugin participates in the configure stage.
type ConfigurePlugin interface {
    Plugin
    Configure(m *Manager) error
}

// ValidatePlugin participates in the validate stage.
type ValidatePlugin interface {
    Plugin
    Validate(m *Manager) error
}

// GlobPlugin participates in the glob stage.
type GlobPlugin interface {
    Plugin
    Glob(m *Manager) error
}

// LoadPlugin participates in the load stage.
type LoadPlugin interface {
    Plugin
    Load(m *Manager) error
}

// TransformPlugin participates in the transform stage.
type TransformPlugin interface {
    Plugin
    Transform(m *Manager) error
}

// RenderPlugin participates in the render stage.
type RenderPlugin interface {
    Plugin
    Render(m *Manager) error
}

// CollectPlugin participates in the collect stage.
type CollectPlugin interface {
    Plugin
    Collect(m *Manager) error
}

// WritePlugin participates in the write stage.
type WritePlugin interface {
    Plugin
    Write(m *Manager) error
}

// CleanupPlugin participates in the cleanup stage.
type CleanupPlugin interface {
    Plugin
    Cleanup(m *Manager) error
}
```

## Creating a Basic Plugin

Here's a minimal plugin that adds a custom field to each post:

```go
package plugins

import (
    "github.com/example/markata-go/pkg/lifecycle"
    "github.com/example/markata-go/pkg/models"
)

// HelloPlugin is a minimal example plugin.
type HelloPlugin struct{}

// NewHelloPlugin creates a new HelloPlugin.
func NewHelloPlugin() *HelloPlugin {
    return &HelloPlugin{}
}

// Name returns the unique plugin identifier.
func (p *HelloPlugin) Name() string {
    return "hello"
}

// Transform adds a greeting to each post.
func (p *HelloPlugin) Transform(m *lifecycle.Manager) error {
    for _, post := range m.Posts() {
        if post.Skip {
            continue
        }
        post.Set("greeting", "Hello from markata-go!")
    }
    return nil
}

// Ensure HelloPlugin implements the required interfaces.
var (
    _ lifecycle.Plugin          = (*HelloPlugin)(nil)
    _ lifecycle.TransformPlugin = (*HelloPlugin)(nil)
)
```

## Plugin Priority and Ordering

Control when your plugin runs relative to other plugins within the same stage using the `PriorityPlugin` interface.

### Priority Constants

```go
const (
    // PriorityFirst ensures a plugin runs before most others.
    PriorityFirst = -1000

    // PriorityEarly ensures a plugin runs early in the stage.
    PriorityEarly = -100

    // PriorityDefault is the default priority.
    PriorityDefault = 0

    // PriorityLate ensures a plugin runs late in the stage.
    PriorityLate = 100

    // PriorityLast ensures a plugin runs after most others.
    PriorityLast = 1000
)
```

### Implementing Priority

```go
// PriorityPlugin can be implemented to control execution order.
type PriorityPlugin interface {
    Plugin
    // Priority returns the plugin's priority for a given stage.
    // Lower values run first.
    Priority(stage Stage) int
}
```

Example implementation:

```go
// Priority returns the priority for the given stage.
// Description should run early in transform to make descriptions
// available for other plugins.
func (p *DescriptionPlugin) Priority(stage lifecycle.Stage) int {
    if stage == lifecycle.StageTransform {
        return lifecycle.PriorityEarly
    }
    return lifecycle.PriorityDefault
}
```

Plugins without the `PriorityPlugin` interface use `PriorityDefault` (0). Within the same priority level, plugins run in registration order.

## Accessing the Manager

The `*lifecycle.Manager` is passed to all hook methods and provides access to:

### Posts

```go
// Get all posts
posts := m.Posts()

// Add a new post
m.AddPost(post)

// Replace all posts
m.SetPosts(posts)
```

### Post lookups

Use the manager's shared index instead of building a slug, href, or path map
inside a plugin:

```go
index := m.PostIndex()
post := index.LookupBySlug("other-post")
postByHref := index.ByHref["/other-post/"]
postByPath := index.ByPath["posts/other-post.md"]
```

`LookupBySlug` accepts case-insensitive slugs, slugified slugs, and aliases.
Href and path keys use exact strings. `SetPosts`, `AddPost`, and `Reset`
invalidate the cached index automatically. If a plugin changes `Slug`, `Href`,
`Path`, or aliases on an existing post, call `index.Refresh(m)` before the
lookup. Keep a separate map when the plugin needs a filtered subset, cache
snapshot, validation of duplicate keys, or different lookup semantics.

### Files

```go
// Get discovered file paths
files := m.Files()

// Set file paths (typically done by glob plugins)
m.SetFiles(files)

// Add a single file
m.AddFile("path/to/file.md")
```

### Configuration

```go
// Get the configuration
config := m.Config()

// Access standard config fields
contentDir := config.ContentDir  // Source directory
outputDir := config.OutputDir    // Output directory
patterns := config.GlobPatterns  // Glob patterns for file discovery

// Access custom config via Extra map
if val, ok := config.Extra["my_setting"].(string); ok {
    // Use val
}
```

### Feeds

```go
// Get all feeds
feeds := m.Feeds()

// Add a feed
m.AddFeed(&lifecycle.Feed{
    Name:    "main",
    Title:   "My Blog",
    Posts:   posts,
    Content: feedXML,
    Path:    "feed.xml",
})
```

### Cache

```go
// Get the cache
cache := m.Cache()

// Store a value
cache.Set("key", value)

// Retrieve a value
if val, ok := cache.Get("key"); ok {
    // Use val
}

// Delete a value
cache.Delete("key")

// Clear all cached data
cache.Clear()
```

### Filtering Posts

```go
// Filter posts using expressions
published, err := m.Filter("published==true")
drafts, err := m.Filter("draft==true")
tagged, err := m.Filter("tags contains golang")

// Complex filters with AND/OR
recent, err := m.Filter("published==true and draft!=true")
```

### Concurrent Processing

For performance, use the built-in concurrent processor:

```go
func (p *MyPlugin) Transform(m *lifecycle.Manager) error {
    return m.ProcessPostsConcurrently(func(post *models.Post) error {
        if post.Skip {
            return nil
        }
        // Process the post (runs in parallel)
        post.Set("computed_field", computeValue(post))
        return nil
    })
}
```

## Extending Configuration

Read custom configuration from the `Extra` map:

```go
func (p *MyPlugin) Configure(m *lifecycle.Manager) error {
    config := m.Config()

    // Read configuration with defaults
    if config.Extra != nil {
        if enabled, ok := config.Extra["my_plugin_enabled"].(bool); ok {
            p.enabled = enabled
        }
        if threshold, ok := config.Extra["my_plugin_threshold"].(int); ok && threshold > 0 {
            p.threshold = threshold
        }
    }

    return nil
}
```

Users configure in their `markata.toml`:

```toml
[markata]
my_plugin_enabled = true
my_plugin_threshold = 100
```

## Working with Posts

### Post Structure

```go
type Post struct {
    Path        string                 // Source file path
    Content     string                 // Raw markdown content
    Slug        string                 // URL-safe identifier
    Href        string                 // Relative URL path (e.g., /my-post/)
    Title       *string                // Optional title
    Date        *time.Time             // Optional publication date
    Published   bool                   // Is the post published?
    Draft       bool                   // Is it a draft?
    Skip        bool                   // Should it be skipped?
    Tags        []string               // Associated tags
    Description *string                // Meta description
    Template    string                 // Template file (default: "post.html")
    HTML        string                 // Final rendered HTML
    ArticleHTML string                 // Content HTML without template
    Extra       map[string]interface{} // Dynamic/custom fields
}
```

### Getting and Setting Custom Fields

```go
// Get a custom field (returns nil if not found)
value := post.Get("custom_field")

// Set a custom field
post.Set("reading_time", 5)
post.Set("word_count", 1200)

// Check if a field exists
if post.Has("custom_field") {
    // ...
}
```

### Respecting Skip Flag

Always check the `Skip` flag before processing:

```go
for _, post := range m.Posts() {
    if post.Skip {
        continue
    }
    // Process the post
}
```

## Complete Example Plugin

Here's a complete plugin that calculates reading time for posts:

```go
package plugins

import (
    "fmt"
    "math"
    "regexp"
    "strings"
    "unicode"

    "github.com/example/markata-go/pkg/lifecycle"
    "github.com/example/markata-go/pkg/models"
)

// ReadingTimePlugin calculates word count and estimated reading time
// for each post during the transform stage.
type ReadingTimePlugin struct {
    // wordsPerMinute is the average reading speed (default: 200)
    wordsPerMinute int
}

// NewReadingTimePlugin creates a new ReadingTimePlugin with default settings.
func NewReadingTimePlugin() *ReadingTimePlugin {
    return &ReadingTimePlugin{
        wordsPerMinute: 200,
    }
}

// Name returns the unique name of the plugin.
func (p *ReadingTimePlugin) Name() string {
    return "reading_time"
}

// Configure reads configuration options for the plugin.
func (p *ReadingTimePlugin) Configure(m *lifecycle.Manager) error {
    config := m.Config()
    if config.Extra != nil {
        if wpm, ok := config.Extra["words_per_minute"].(int); ok && wpm > 0 {
            p.wordsPerMinute = wpm
        }
    }
    return nil
}

// Transform calculates word count and reading time for each post.
func (p *ReadingTimePlugin) Transform(m *lifecycle.Manager) error {
    return m.ProcessPostsConcurrently(func(post *models.Post) error {
        if post.Skip || post.Content == "" {
            return nil
        }

        // Count words
        wordCount := p.countWords(post.Content)
        post.Set("word_count", wordCount)

        // Calculate reading time in minutes
        readingTime := p.calculateReadingTime(wordCount)
        post.Set("reading_time", readingTime)

        // Also store a formatted string
        post.Set("reading_time_text", p.formatReadingTime(readingTime))

        return nil
    })
}

// Regex pattern to match code blocks
var codeBlockPattern = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~|`[^`]+`")

// countWords counts the number of words in markdown content.
func (p *ReadingTimePlugin) countWords(content string) int {
    // Remove code blocks
    text := codeBlockPattern.ReplaceAllString(content, " ")

    // Count words
    words := 0
    inWord := false

    for _, r := range text {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            if !inWord {
                words++
                inWord = true
            }
        } else {
            inWord = false
        }
    }

    return words
}

// calculateReadingTime estimates reading time in minutes.
func (p *ReadingTimePlugin) calculateReadingTime(wordCount int) int {
    if wordCount == 0 {
        return 0
    }

    minutes := float64(wordCount) / float64(p.wordsPerMinute)
    return int(math.Ceil(minutes))
}

// formatReadingTime creates a human-readable reading time string.
func (p *ReadingTimePlugin) formatReadingTime(minutes int) string {
    if minutes == 0 {
        return "< 1 min read"
    }
    if minutes == 1 {
        return "1 min read"
    }
    return fmt.Sprintf("%d min read", minutes)
}

// Ensure ReadingTimePlugin implements the required interfaces.
var (
    _ lifecycle.Plugin          = (*ReadingTimePlugin)(nil)
    _ lifecycle.ConfigurePlugin = (*ReadingTimePlugin)(nil)
    _ lifecycle.TransformPlugin = (*ReadingTimePlugin)(nil)
)
```

## Testing Plugins

### Unit Test Example

```go
package plugins

import (
    "testing"

    "github.com/example/markata-go/pkg/lifecycle"
    "github.com/example/markata-go/pkg/models"
)

func TestReadingTimePlugin_Transform(t *testing.T) {
    tests := []struct {
        name            string
        content         string
        wordsPerMinute  int
        wantWordCount   int
        wantReadingTime int
    }{
        {
            name:            "short post",
            content:         "Hello world this is a test.",
            wordsPerMinute:  200,
            wantWordCount:   6,
            wantReadingTime: 1, // Minimum 1 minute
        },
        {
            name:            "longer post",
            content:         strings.Repeat("word ", 400),
            wordsPerMinute:  200,
            wantReadingTime: 2,
        },
        {
            name:            "custom WPM",
            content:         strings.Repeat("word ", 100),
            wordsPerMinute:  100,
            wantReadingTime: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create plugin with custom WPM
            plugin := NewReadingTimePlugin()
            plugin.wordsPerMinute = tt.wordsPerMinute

            // Create manager with a test post
            m := lifecycle.NewManager()
            post := models.NewPost("test.md")
            post.Content = tt.content
            m.AddPost(post)

            // Run transform
            err := plugin.Transform(m)
            if err != nil {
                t.Fatalf("Transform() error = %v", err)
            }

            // Check results
            posts := m.Posts()
            if len(posts) != 1 {
                t.Fatalf("Expected 1 post, got %d", len(posts))
            }

            readingTime, ok := posts[0].Get("reading_time").(int)
            if !ok {
                t.Fatal("reading_time not set")
            }
            if readingTime != tt.wantReadingTime {
                t.Errorf("reading_time = %d, want %d", readingTime, tt.wantReadingTime)
            }
        })
    }
}

func TestReadingTimePlugin_SkipsSkippedPosts(t *testing.T) {
    plugin := NewReadingTimePlugin()

    m := lifecycle.NewManager()
    post := models.NewPost("test.md")
    post.Content = "Some content here"
    post.Skip = true
    m.AddPost(post)

    err := plugin.Transform(m)
    if err != nil {
        t.Fatalf("Transform() error = %v", err)
    }

    // Verify reading_time was NOT set
    if m.Posts()[0].Has("reading_time") {
        t.Error("reading_time should not be set for skipped posts")
    }
}
```

### Integration Test Example

```go
func TestReadingTimePlugin_Integration(t *testing.T) {
    // Create a temporary directory with test files
    tmpDir := t.TempDir()

    // Write a test markdown file
    content := `---
title: Test Post
published: true
---

This is test content with enough words to calculate reading time.
` + strings.Repeat("word ", 200)

    err := os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte(content), 0644)
    if err != nil {
        t.Fatal(err)
    }

    // Setup manager
    m := lifecycle.NewManager()
    cfg := m.Config()
    cfg.ContentDir = tmpDir
    cfg.OutputDir = filepath.Join(tmpDir, "output")

    // Register plugins
    m.RegisterPlugins(
        NewGlobPlugin(),
        NewLoadPlugin(),
        NewReadingTimePlugin(),
    )

    // Run through transform stage
    err = m.RunTo(lifecycle.StageTransform)
    if err != nil {
        t.Fatalf("RunTo() error = %v", err)
    }

    // Verify reading time was calculated
    posts := m.Posts()
    if len(posts) != 1 {
        t.Fatalf("Expected 1 post, got %d", len(posts))
    }

    readingTime, ok := posts[0].Get("reading_time").(int)
    if !ok {
        t.Fatal("reading_time not set")
    }
    if readingTime < 1 {
        t.Errorf("reading_time = %d, expected >= 1", readingTime)
    }
}
```

## Registering Plugins

### Built-in Registry

markata-go includes a plugin registry for managing plugins by name:

```go
// Register a plugin constructor
plugins.RegisterPluginConstructor("my_plugin", func() lifecycle.Plugin {
    return NewMyPlugin()
})

// Get a plugin by name
plugin, ok := plugins.PluginByName("my_plugin")
if !ok {
    log.Fatal("Plugin not found")
}

// List all registered plugins
names := plugins.RegisteredPlugins()
```

### Using Default Plugins

```go
// Get all default plugins
m := lifecycle.NewManager()
m.RegisterPlugins(plugins.DefaultPlugins()...)

// Or use minimal set
m.RegisterPlugins(plugins.MinimalPlugins()...)
```

### Manual Registration

```go
m := lifecycle.NewManager()

// Register individual plugins
m.RegisterPlugin(plugins.NewGlobPlugin())
m.RegisterPlugin(plugins.NewLoadPlugin())
m.RegisterPlugin(NewMyCustomPlugin())

// Or register multiple at once
m.RegisterPlugins(
    plugins.NewGlobPlugin(),
    plugins.NewLoadPlugin(),
    plugins.NewRenderMarkdownPlugin(),
    NewMyCustomPlugin(),
)
```

### Plugin Loading from Config

```go
// Load plugins by name from configuration
pluginNames := []string{"glob", "load", "render_markdown", "templates"}
loadedPlugins, warnings := plugins.PluginsByNames(pluginNames)

for _, w := range warnings {
    log.Printf("Warning: %s", w)
}

m.RegisterPlugins(loadedPlugins...)
```

## Error Handling

### Returning Errors

Return errors from hook methods to signal failure:

```go
func (p *MyPlugin) Transform(m *lifecycle.Manager) error {
    for _, post := range m.Posts() {
        if err := p.processPost(post); err != nil {
            return fmt.Errorf("processing %s: %w", post.Path, err)
        }
    }
    return nil
}
```

### Critical vs Non-Critical Stages

Some stages are critical and halt execution on error:
- **Critical:** Configure, Validate, Glob, Load
- **Non-Critical:** Transform, Render, Collect, Write, Cleanup

For non-critical stages, individual post errors can be logged while allowing the build to continue.

### Accessing Warnings

```go
// After running stages
err := m.Run()

// Check warnings even if no error
for _, warning := range m.Warnings() {
    log.Printf("Warning: %s", warning)
}
```

## Best Practices

### 1. Always Check Configuration

```go
func (p *MyPlugin) Configure(m *lifecycle.Manager) error {
    config := m.Config()

    // Provide sensible defaults
    p.enabled = true
    p.threshold = 100

    // Override with config values if present
    if config.Extra != nil {
        if enabled, ok := config.Extra["my_plugin_enabled"].(bool); ok {
            p.enabled = enabled
        }
    }

    return nil
}
```

### 2. Respect the Skip Flag

```go
for _, post := range m.Posts() {
    if post.Skip {
        continue  // Always skip posts marked for skipping
    }
    // Process the post
}
```

### 3. Use Concurrent Processing for Performance

```go
func (p *MyPlugin) Transform(m *lifecycle.Manager) error {
    return m.ProcessPostsConcurrently(func(post *models.Post) error {
        // This runs in parallel across posts
        return p.processPost(post)
    })
}
```

### 4. Implement Interface Verification

At the end of your plugin file, verify interface implementation:

```go
// Ensure MyPlugin implements the required interfaces at compile time.
var (
    _ lifecycle.Plugin          = (*MyPlugin)(nil)
    _ lifecycle.ConfigurePlugin = (*MyPlugin)(nil)
    _ lifecycle.TransformPlugin = (*MyPlugin)(nil)
)
```

### 5. Handle Missing Fields Gracefully

```go
// Check for nil pointers
if post.Title != nil {
    title = *post.Title
} else {
    title = post.Slug
}

// Use Get() with type assertion
if val, ok := post.Get("custom_field").(string); ok {
    // Use val
}
```

### 6. Document Your Plugin

Include clear documentation at the top of your plugin file:

```go
// Package plugins provides lifecycle plugins for markata-go.
//
// MyPlugin does X, Y, and Z.
//
// Configuration:
//     [markata]
//     my_plugin_enabled = true
//     my_plugin_threshold = 100
//
// Post fields set:
//     - my_computed_field: Description of the field
//
// Usage in templates:
//     {{ post.my_computed_field }}
package plugins
```

## See Also

- [Lifecycle Stages Specification](../../spec/spec/LIFECYCLE.md) - Detailed stage documentation
- [Plugin Specification](../../spec/spec/PLUGINS.md) - Full plugin development specification
- [Built-in Plugins](../plugins/) - Documentation for built-in plugins
