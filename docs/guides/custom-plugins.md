---
title: "Custom plugins with your own agent"
description: "Extend the Markata CLI with a site-owned, version-pinned Go plugin."
date: 2026-11-22
published: true
tags: [documentation, plugins, agents]
---

# Site-owned Markata plugins

Markata's stock CLI **does not dynamically compile or load arbitrary Go files
from a site's `plugins/` directory**. A custom Go plugin needs an executable
that registers it. The tested starter at
[`examples/custom-plugin/main.go`](https://github.com/WaylonWalker/markata/blob/main/examples/custom-plugin/main.go)
shows how to do this *without forking Markata, adding Node, or building a
hosting control plane*.

The starter uses existing public Markata imports:
`github.com/WaylonWalker/markata/cmd/markata-go/cmd`,
`pkg/plugins`, `pkg/lifecycle`, and `pkg/models`.

## Use the starter in a site you own

Ask your agent to copy the starter `main.go` into your chosen site directory,
then run these commands **from that directory**:

```sh
go mod init example.com/my-markata-site
go get github.com/WaylonWalker/markata@v0.15.0
```

Replace the example module name with your site's own identity. Keep
`go.mod` and `go.sum` in version control. Pin a released Markata version
rather than depending on `@latest` or mutable `main`. Make sure the Go
toolchain meets the installed Markata release requirement.

Add the plugin name to the site's `markata-go.toml`, preserving the site's
existing hooks:

```toml
[markata-go]
hooks = ["default", "greeting"]
```

Now build using the **site-owned executable**, not the stock `markata`:

```sh
go test ./...
go run . build --clean
```

The `greeting` hook sets `post.Extra.greeting` on non-skipped posts during
Transform. An overridden template can render it using
`{{ post.Extra.greeting }}`. Test the resulting HTML before publishing.

The wrapper calls Markata's existing CLI entry point. Other Markata commands,
including `serve` and `version`, work through this same executable.
Production deployment remains the user's choice: run the pinned custom binary
in the chosen CI environment, then upload the generated static output to the
hosting provider. A local build is not permission to publish.

## How to customize safely

An agent should inspect existing site config, themes and templates *before*
creating a Go plugin. Change one stage at a time, avoid network/time/random
inputs in deterministic transforms, and add a unit test for both transformed
and skipped content. Use a full production build for release validation, not
`--fast`. To roll back, remove `greeting` from the configured hooks and use
the stock `markata build` or revert the wrapper to a pinned known-good commit.

For portable agent instructions use `markata --skill`. To install the full
reference bundle for your agent, use `markata agent install`.

## Limits

This is **compile-time Go customization**. It is not runtime loading of
untrusted site-local code, WASM, or a subprocess plugin protocol. It carries
the Go dependency/build cost of a custom executable but does not add any
runtime service. All custom code still runs with the builder's filesystem and
network privileges, so inspect third-party plugin sources before building.
