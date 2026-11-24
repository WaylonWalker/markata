---
title: "Nginx-native redirects"
description: "The built-in redirects plugin generates nginx-native redirect rules from the same static/redirects file used for static HTML fallbacks."
date: 2026-09-27
published: true
tags:
  - documentation
---


The built-in redirects plugin generates nginx-native redirect rules from the same `static/_redirects` file used for static HTML fallbacks.

## Default behavior

Given:

```text
/old-post /new-post
/go/docs https://docs.example.com/start
```

A normal build writes under the configured `output_dir`:

```text
<output_dir>/
├── redirects.conf
├── old-post/index.html
└── go/docs/index.html
```

`redirects.conf` contains exact-match nginx rules:

```nginx
location = "/old-post" {
    return 301 "/new-post";
}

location = "/go/docs" {
    return 301 "https://docs.example.com/start";
}
```

HTML fallback pages remain enabled by default so the same build artifact continues to work on static hosts that do not consume nginx configuration.

## Configuration

```toml
[markata-go.redirects]
redirects_file = "static/_redirects"
redirect_template = "templates/redirect.html" # optional
html_fallback = true
```

Set `html_fallback = false` when the deployment is guaranteed to use the generated nginx rules and you do not want new per-path HTML redirect pages:

```toml
[markata-go.redirects]
html_fallback = false
```

When changing this setting on an existing incremental output directory, use a clean build if you also want old HTML fallback files removed. Markata-go deliberately does not delete arbitrary `index.html` files because redirect paths can overlap normal site content and custom templates do not provide a safe ownership marker.

## Keeping `redirects.conf` current

On normal builds, nginx output is refreshed whenever the configured `_redirects` file exists. If that file is empty, markata-go writes a valid header-only `redirects.conf`, so previously generated native rules cannot remain active.

If the `_redirects` source is deleted, markata-go removes an existing `redirects.conf` only when it recognizes the file as its own generated artifact. A user-managed `redirects.conf` is left untouched.

When `_redirects` exists, `redirects.conf` is the plugin's generated output path. If the output directory already contains a file at that path that is not marked as Markata-generated—for example `static/redirects.conf` copied by the static-assets plugin—the build fails with a clear conflict instead of overwriting it. Rename or remove the manually managed file before enabling generated nginx redirects.

## Include from nginx

Include the generated file from inside the nginx `server` block that serves the built site:

```nginx
server {
    listen 80;
    server_name example.com;
    root /usr/share/nginx/html;

    include /usr/share/nginx/html/redirects.conf;

    location / {
        try_files $uri $uri/ /index.html =404;
    }
}
```

The include path must match the deployed `output_dir`.

### Release directories and `current` symlinks

If a deployment serves releases through a `current` symlink, point the include at the active release only when the nginx process will be reloaded as part of the release switch. Nginx reads included configuration when it loads or reloads its configuration; changing the symlink alone does not make an already-running nginx process reread `redirects.conf`.

For builder-admin or other in-place release switching, either:

- reload nginx after activating a release whose redirect config changed, or
- keep nginx-native redirects out of that live-switch path until a reload hook is part of the deployment.

HTML fallbacks continue to work without an nginx reload because they are served as ordinary files from the active release.

## Supported rules

The native output intentionally follows the redirects plugin's existing parser:

- source paths must start with `/`
- destinations must be `/absolute-path`, `http://...`, or `https://...`
- comments and blank lines are ignored
- malformed rules are ignored
- wildcard rules containing `*` are skipped
- extra fields such as platform-specific status codes are currently ignored
- generated nginx redirects use HTTP `301`

This keeps nginx output and HTML fallback behavior aligned. Use hand-written nginx configuration or the hosting provider's native redirect system for wildcard, conditional, or status-code-specific rules.

## Fast builds

The redirects plugin is skipped in `--fast` mode. Run a normal build before publishing, validating, or reloading nginx with `redirects.conf`.
