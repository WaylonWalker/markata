# Mermaid Pre-Rendering Specification

This document specifies the Mermaid diagram pre-rendering feature, which allows rendering Mermaid diagrams to static SVGs at build time instead of relying on client-side JavaScript.

## Overview

The Mermaid plugin supports three rendering modes:

| Mode | Rendering | Setup | Speed | Dependencies |
|------|-----------|-------|-------|--------------|
| **client** | Browser (Mermaid.js) | None | Fast | CDN (or offline) |
| **cli** | Build-time (mmdc CLI) | npm install | Moderate | Node.js + npm |
| **chromium** | Build-time (Chrome) | Install browser | Fast | Chrome/Chromium |

## Configuration

All three modes are configured under `[markata-go.mermaid]`:

```toml
[markata-go.mermaid]
enabled = true
mode = "client"                    # "client", "cli", or "chromium"
theme = "default"                  # default, dark, forest, neutral
use_css_variables = true           # Apply site palette to diagrams
lightbox = true                    # Enable click-to-zoom lightbox
lightbox_selector = ".glightbox-mermaid"

# CLI mode settings (npm mmdc)
[markata-go.mermaid.cli]
mmdc_path = ""                     # Auto-detect if empty

# Chromium mode settings (chromedp)
[markata-go.mermaid.chromium]
browser_path = ""                  # Auto-detect if empty
timeout = 30                       # Seconds per diagram
max_concurrent = 4                 # Parallel diagram renderers
no_sandbox = false                 # Disable Chromium sandbox (required in containers)
```

## Mode Details

### Client Mode (Default)

**Rendering:** Browser-side via Mermaid.js CDN  
**Stage:** Render (post-processing)  
**Requirements:** None

Mermaid code blocks are converted to `<pre class="mermaid">` and Mermaid.js is loaded via CDN. Diagrams render in the browser.

```toml
[markata-go.mermaid]
mode = "client"
cdn_url = "/assets/vendor/mermaid/mermaid.min.js"
```

**Behavior:**
1. Find `<pre><code class="language-mermaid">` blocks
2. Replace with `<pre class="mermaid">{diagram code}</pre>`
3. Inject Mermaid.js module initialization script
4. Browser renders diagrams on load

**Pros:**
- No dependencies
- Simple setup
- Works offline after first load
- Can change themes dynamically

**Cons:**
- Requires JavaScript
- Not feed-friendly (RSS/Atom)
- Slower initial page load
- Browser CPU usage

---

### CLI Mode (npm mmdc)

**Rendering:** Build-time via mermaid-cli  
**Stage:** Render (post-processing)  
**Requirements:** Node.js, npm, @mermaid-js/mermaid-cli

The `mmdc` (Mermaid CLI) tool converts Mermaid diagrams to SVGs during the build process.

```toml
[markata-go.mermaid]
mode = "cli"

[markata-go.mermaid.cli]
# mmdc_path = "/usr/local/bin/mmdc"  # Optional: specify path
```

**Setup:**

```bash
# Install Node.js v14+
node --version  # Should be v14+

# Install mermaid-cli globally
npm install -g @mermaid-js/mermaid-cli

# Verify
mmdc --version
```

**Behavior:**
1. Find `<pre><code class="language-mermaid">` blocks
2. Write diagram code to temporary `.mmd` file
3. Call `mmdc -i diagram.mmd -o diagram.svg`
4. Read SVG output and embed in HTML
5. Remove temporary files

**Output:**
```html
<pre class="mermaid">
  <svg><!-- embedded SVG content --></svg>
</pre>
```

**Pros:**
- Pre-rendered to static SVGs
- No client-side JavaScript needed
- Works in feeds (RSS/Atom)
- Faster page load
- Deterministic rendering

**Cons:**
- Requires Node.js installation
- Slower builds (one process per diagram)
- Version locked to build time
- CLI tool maintenance dependency

---

### Chromium Mode (chromedp)

**Rendering:** Build-time via Chrome DevTools Protocol  
**Stage:** Render (post-processing)  
**Requirements:** Chrome/Chromium browser binary

Uses the project's chromedp renderer to render diagrams in a headless Chrome browser. It reuses one browser connection for all diagrams. The renderer applies configured palette colors when available and uses the configured Mermaid theme otherwise.

```toml
[markata-go.mermaid]
mode = "chromium"

[markata-go.mermaid.chromium]
# browser_path = "/usr/bin/chromium"  # Optional: specify path
timeout = 30
max_concurrent = 4
# no_sandbox = true                   # Required in containers (Docker, Distrobox, etc.)
```

**Container setup:**

When running inside Docker, Distrobox, or other containerized environments, you must
set `no_sandbox = true` because Chromium's sandbox requires kernel features
typically unavailable in containers.

```toml
[markata-go.mermaid.chromium]
no_sandbox = true
```

**Browser auto-detection:** The plugin searches for browser binaries in this order:
`headless-shell`, `headless_shell`, `google-chrome`, `chromium-browser`, `chromium`,
`google-chrome-stable`. It also verifies the binary is functional by running
`--version`, which catches Ubuntu snap stubs that are not real browsers.

**MermaidJS source caching:** The MermaidJS library (~3.3 MB) is downloaded once
and cached at `~/.cache/markata-go/mermaid/mermaid-v{version}.min.js` following
XDG conventions. Subsequent builds load from cache, eliminating the network
download. The cache persists across projects and `--clean` / `--clean-all`.

**Rendered SVG caching:** In `cli` and `chromium` modes, each rendered SVG is
cached in the build cache and keyed by the diagram source plus rendering inputs:
mode, theme, CSS-variable usage, lightbox setting, renderer configuration, and
resolved palette theme variables. Hot builds reuse matching SVGs without
invoking the external renderer. The cache MUST invalidate when any keyed input
changes, and cache failures MUST NOT fail the build after rendering succeeds.

**Setup:**

```bash
# Linux (Debian/Ubuntu)
sudo apt-get install chromium-browser

# Linux (Fedora/RHEL)
sudo dnf install chromium

# macOS
brew install chromium

# Windows
choco install chromium

# Or download from: https://www.chromium.org/getting-involved/download-chromium

# Lightweight alternative (no GUI dependencies):
# Download chrome-headless-shell from Chrome for Testing
# https://googlechromelabs.github.io/chrome-for-testing/
# Place the binary in your PATH as "headless-shell"
```

**Behavior:**
1. Start single Chrome instance (reused across all diagrams)
2. Find `<pre><code class="language-mermaid">` blocks
3. For each diagram (up to `max_concurrent` in parallel):
   - Reuse cached SVG output when the diagram and rendering inputs match
   - Otherwise render via Chrome DevTools Protocol
   - Extract SVG output and store it in the build cache
   - Embed the SVG in HTML
4. Close Chrome connection

**Output:**
```html
<pre class="mermaid">
  <svg><!-- embedded SVG content --></svg>
</pre>
```

**Pros:**
- Pre-rendered to static SVGs
- No JavaScript needed
- Fast builds (Chrome reused, parallel rendering)
- Pure Go integration
- Works in feeds (RSS/Atom)
- Best performance for large diagram counts

**Cons:**
- Requires Chrome/Chromium binary
- Only supports Mermaid syntax that Chrome supports
- More complex than CLI mode

---

## Error Handling

The plugin validates dependencies at configuration time and provides helpful error messages with installation instructions.

### Dependency Not Found

If the selected mode's dependencies are missing, the build fails with a detailed error:

```
error: mermaid render error in posts/example.md (cli mode): mmdc binary not found

Suggestion:
Missing dependency: @mermaid-js/mermaid-cli

Installation instructions:

1. Install Node.js v14+ from https://nodejs.org/
   (Check: node --version)

2. Install mermaid-cli globally:
   npm install -g @mermaid-js/mermaid-cli

3. Verify installation:
   mmdc --version

Or specify the path explicitly in your config:
  [markata-go.mermaid.cli]
  mmdc_path = "/path/to/mmdc"

To use client-side rendering instead, change your config:
  [markata-go.mermaid]
  mode = "client"
```

### Error Recovery

1. **Fix the dependency** - Install the required tool
2. **Specify the path** - If tool is installed at non-standard location
3. **Switch modes** - Change to `mode = "client"` for no-dependency fallback

## Feature Behavior

### Common Behavior (All Modes)

**Detection:** Posts with `<pre><code class="language-mermaid">` or `class="mermaid"` blocks

**Theme Application:** When `use_css_variables = true`:
- Extract CSS custom properties: `--color-background`, `--color-text`, `--color-primary`, `--color-code-bg`, `--color-surface`
- Pass as Mermaid `themeVariables`
- Fallback to hardcoded defaults if CSS props missing

**Lightbox:** When `lightbox = true`:
- Wrap SVG in clickable container
- Load svg-pan-zoom on first click
- Enable pan/zoom interaction
- Show toolbar (Fit/+/- buttons)

### Mode-Specific Behavior

**Client mode:** CSS variables applied in browser via JavaScript

**CLI/Chromium modes:** CSS variables applied at build time; SVG is static

## Configuration Validation

The plugin validates configuration during the `Configure` stage:

1. Check `mode` is valid: "client", "cli", or "chromium"
2. Check dependencies exist:
   - **client** - No validation (CDN-based)
   - **cli** - Verify `mmdc` binary found or provide path
   - **chromium** - Verify Chrome/Chromium binary found or provide path
3. Check paths if explicitly specified
4. Fail fast with helpful error message if validation fails

## Performance Characteristics

### Build Time

| Mode | Time per Diagram | Parallel |
|------|------------------|----------|
| **client** | ~0ms (skipped) | N/A |
| **cli** | 100-500ms | Sequential |
| **chromium** | 50-200ms | Yes (max_concurrent) |

### Output Size

| Mode | HTML Size | Requires JS |
|------|-----------|-------------|
| **client** | Smaller (~1KB) | Yes |
| **cli/chromium** | Larger (~5-20KB per SVG) | No |

### Feed Compatibility

| Mode | Works in Feeds | Notes |
|------|----------------|-------|
| **client** | ❌ No | Feeds get raw code block |
| **cli** | ✅ Yes | SVG embedded directly |
| **chromium** | ✅ Yes | SVG embedded directly |

## Implementation Details

### Configuration Structure

```go
type MermaidConfig struct {
    Enabled            bool
    Mode               string                    // "client", "cli", "chromium"
    Theme              string
    UseCSSVariables    bool
    Lightbox           bool
    LightboxSelector   string
    CDNURL             string                    // client mode only
    CLIConfig          *CLIRendererConfig
    ChromiumConfig     *ChromiumRendererConfig
}

type CLIRendererConfig struct {
    MMDCPath  string
    ExtraArgs string
}

type ChromiumRendererConfig struct {
    BrowserPath   string
    Timeout       time.Duration
    MaxConcurrent int
    NoSandbox     bool
}
```

### Plugin Interface

The plugin must implement:
- `Plugin` - Basic interface with `Name()` method
- `ConfigurePlugin` - Read configuration and validate dependencies
- `RenderPlugin` - Process posts and render diagrams
- `PriorityPlugin` - Ensure execution after markdown rendering

### Error Types

New error type for mermaid-specific failures:

```go
type MermaidRenderError struct {
    Path       string  // Post file path
    DiagramID  string  // Optional diagram identifier
    Mode       string  // "client", "cli", "chromium"
    Message    string
    Err        error
    Suggestion string  // Installation instructions
}
```

## Backward Compatibility

The default mode is `"client"` (existing behavior). Users on older configs without the `mode` field will continue using client-side rendering without changes.

## Future Considerations

1. **Server-side rendering to static files** - Generate SVG files separate from HTML
2. **Mermaid theme library** - Pre-defined themes matching popular site designs
3. **Diagram caching** - Cache rendered SVGs to skip re-rendering unchanged diagrams
4. **Format support** - Output as PNG/PDF in addition to SVG
