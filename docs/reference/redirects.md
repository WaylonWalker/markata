---
title: Redirects Reference
description: Configuration and generated outputs for the built-in redirects plugin.
---


The built-in `redirects` plugin reads a `_redirects` source file and generates nginx-native redirects plus optional HTML fallback pages.

For deployment guidance and nginx reload behavior, see [Nginx-native Redirects](/docs/guides/nginx-redirects/).

## Configuration

```toml
[markata-go.redirects]
redirects_file = "static/_redirects"
redirect_template = "templates/redirect.html" # optional
html_fallback = true                            # default
```

| Option | Default | Description |
|---|---|---|
| `redirects_file` | `static/_redirects` | Redirect source file |
| `redirect_template` | built-in template | Optional custom HTML fallback template |
| `html_fallback` | `true` | Generate per-source HTML redirect pages in addition to `redirects.conf` |

## Source format

```text
/old-path /new-path
/go/docs https://docs.example.com/start
```

Supported rules require an absolute source path and a destination beginning with `/`, `http://`, or `https://`. Wildcards and malformed rules are skipped. Extra fields such as provider-specific status codes are currently ignored.

## Generated outputs

A normal build with a redirects source writes:

- `<output_dir>/redirects.conf` with exact-match nginx `301` rules
- `<output_dir>/<source>/index.html` when `html_fallback = true`

An empty redirects source writes a header-only `redirects.conf` so old native rules are cleared. Removing the source removes a stale `redirects.conf` only when Markata recognizes the file as generated output.

If an unowned `redirects.conf` already exists while `_redirects` is present, the build fails instead of overwriting the user-managed file.

## Fast builds

Redirect generation is skipped in `--fast` mode. Use a normal build before validating or deploying native redirect configuration.
