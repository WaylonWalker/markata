# Plugin Creation

Use this topic only after checking whether the task can be solved in config, content, templates, CSS, feeds, or an existing built-in feature.

## Important Constraint

The stock `markata-go` CLI does not load arbitrary Go plugins from a site directory at runtime.

Today, a custom plugin means one of these paths:

1. add the plugin upstream in the `markata-go` codebase
2. build a custom Go binary or wrapper binary that registers the plugin constructor and then runs markata-go lifecycle code

Do not assume a site-local `plugins/` directory is automatically discovered by the stock CLI.

## Decision Rule

Reach for a plugin when the change needs new build-stage behavior, derived content, external data, or output generation that templates alone cannot express.

## Lifecycle Stages

Markata-go runs plugins through these stages:

- Configure
- Validate
- Glob
- Load
- Transform
- Render
- Collect
- Write
- Cleanup

## Matching Work To A Stage

- config parsing or plugin setup -> Configure
- config sanity checks -> Validate
- discover files -> Glob
- parse new content types -> Load
- modify markdown or derived fields before render -> Transform
- generate HTML or wrap rendered content -> Render
- build feeds, archives, or aggregate data -> Collect
- write output files -> Write
- flush caches or cleanup -> Cleanup

## Guidance

- Prefer matching the work to one lifecycle stage instead of spreading logic across many stages.
- Follow the existing plugin style in the current codebase.
- Update spec and user docs alongside plugin behavior.
- Validate whether the task belongs in a site-local customization or in `markata-go` upstream.
- If the task is site-only and mostly presentational, templates are usually cheaper than plugin code.

## Site Repo Vs Upstream Repo

Use this rule:

- if the task changes only one site's content model, templates, output choices, or local integrations, it probably belongs in the site repo
- if the task adds generally useful build behavior, new CLI support, reusable config, or a new built-in plugin, it probably belongs in `markata-go` upstream

If the user asks for a plugin inside a standalone site repo, first verify whether they also want a custom binary or whether the work should move upstream.

## Before Writing A Plugin

Check these lower-cost options first:

1. config change
2. frontmatter field
3. feed definition
4. layout config
5. template or partial override
6. CSS change

If all of those fail to express the behavior, then plugin work is justified.

## Minimal Working Plugin Example

This is the minimum shape of a transform-stage plugin:

```go
package plugins

import (
    "github.com/WaylonWalker/markata/pkg/lifecycle"
    "github.com/WaylonWalker/markata/pkg/models"
)

type GreetingPlugin struct{}

func NewGreetingPlugin() *GreetingPlugin {
    return &GreetingPlugin{}
}

func (p *GreetingPlugin) Name() string {
    return "greeting"
}

func (p *GreetingPlugin) Transform(m *lifecycle.Manager) error {
    return m.ProcessPostsConcurrently(func(post *models.Post) error {
        if post.Skip {
            return nil
        }
        post.Set("greeting", "Hello from a custom plugin")
        return nil
    })
}

var (
    _ lifecycle.Plugin          = (*GreetingPlugin)(nil)
    _ lifecycle.TransformPlugin = (*GreetingPlugin)(nil)
)
```

## Site-Owned Custom Binary (Supported Starter)

The repository now includes `examples/custom-plugin/main.go`: a minimal
working wrapper that calls `plugins.RegisterPluginConstructor` and
`cmd.Execute()`. Copy it into a site-owned Go module, pin a released Markata
version (for example `github.com/WaylonWalker/markata@v0.15.0`), and add
`greeting` to the site's `hooks` configuration. Build with `go run . build`
or a compiled site binary, not the unmodified stock `markata`. The matching
test file demonstrates deterministic metadata and skipped-post handling.

Read [Custom Plugins](/docs/guides/custom-plugins/) for the complete setup,
validation, security and rollback steps. No custom runtime loader is implied.

## Registration Pattern

Custom plugins are used by registering a constructor in Go:

```go
plugins.RegisterPluginConstructor("greeting", func() lifecycle.Plugin {
    return NewGreetingPlugin()
})
```

Then the manager can load it by name, or a custom binary can register it directly with the manager.

## Realistic Loading Model

To actually use a custom plugin, one of these must happen:

### Upstream path

- add the plugin to the markata-go codebase
- register its constructor
- include it in the chosen plugin set or config-driven name resolution

### Custom binary path

- create a small Go program that imports markata-go packages
- register the custom plugin constructor
- create and run a lifecycle manager with default plugins plus the custom plugin

If neither path is available, the plugin cannot be used by the stock `markata-go build` command.

## Manager APIs Agents Usually Need

Useful manager methods and data:

- `m.Posts()`
- `m.AddPost(...)`
- `m.SetPosts(...)`
- `m.Files()`
- `m.SetFiles(...)`
- `m.Config()`
- `m.Feeds()`
- `m.AddFeed(...)`
- `m.Cache()`
- `m.ProcessPostsConcurrently(...)`

Useful post operations:

- `post.Skip`
- `post.Content`
- `post.ArticleHTML`
- `post.HTML`
- `post.Tags`
- `post.Set("key", value)`
- `post.Get("key")`
- `post.Has("key")`

## Config Pattern For Plugin Behavior

Plugin-specific settings are typically read from `m.Config().Extra`.

Example:

```go
func (p *GreetingPlugin) Configure(m *lifecycle.Manager) error {
    if enabled, ok := m.Config().Extra["greeting_enabled"].(bool); ok && !enabled {
        return nil
    }
    return nil
}
```

## Priority And Ordering

If ordering matters, implement `Priority(stage lifecycle.Stage) int`.

Typical constants are:

- `lifecycle.PriorityFirst`
- `lifecycle.PriorityEarly`
- `lifecycle.PriorityDefault`
- `lifecycle.PriorityLate`
- `lifecycle.PriorityLast`

Use this when another plugin depends on fields your plugin computes.

## Validation Workflow

For plugin work, the minimum validation loop is:

1. `go test ./...` or focused package tests for the plugin code
2. `markata-go build -v` using the binary that actually includes the plugin
3. inspect generated output or derived post fields
4. if ordering looks wrong, inspect plugin stage and priority before changing code shape

## Decision Shortcut For Agents

If the repo is only a content site and does not include Go code or a custom markata-go wrapper binary, do not promise a working custom plugin through the stock CLI.

Instead:

- propose a config/template/feed solution if possible, or
- tell the user the plugin must be added upstream or via a custom binary

## Start With

- `markata-go explain plugins`
- plugin development docs for the current version of markata-go
- the closest built-in plugin to the behavior you need
- current `hooks` and `disabled_hooks` settings in site config

## Common Plugin Use Cases

- fetch and cache remote data
- add derived post fields
- create new output files or collection pages
- transform markdown or rendered HTML before publishing

## Operational Debugging Tips

- if a plugin seems missing, check `hooks = ["default"]` and `disabled_hooks`
- use `markata-go build -v` when you need plugin-stage error detail
- isolate plugin issues by temporarily disabling the suspected hook rather than changing templates first
- if a custom plugin never appears, verify the binary actually registers the constructor; config alone is not enough
