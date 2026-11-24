---
title: "View Transitions Configuration"
description: "Complete configuration reference for View Transitions API settings"
date: 2026-02-02
published: true
tags:
  - view-transitions
  - configuration
  - reference
---


Complete reference for configuring the View Transitions API behavior in markata-go.

## Overview

View Transitions can be configured via your `markata.toml` or YAML config file. All settings are optional - sensible defaults are provided.

## Configuration Options

Add a `[markata-go.view_transitions]` section to your config:

```toml
[markata-go.view_transitions]
enabled = true              # Enable/disable view transitions globally
debug = false               # Log debug messages to console
duration = 120              # Default transition duration in ms
update_meta = true          # Update meta tags on navigation
scroll_to_top = true        # Scroll to top on navigation (unless hash present)
skip_classes = []           # Additional CSS classes to skip
skip_selectors = []         # Additional selectors to skip
```

In TOML, nest this under `[markata-go.view_transitions]` because markata-go loads site config from the top-level `[markata-go]` table.

## Quick Examples

### TOML (`markata.toml`)

```toml
[markata-go.view_transitions]
enabled = true
debug = false
duration = 120
update_meta = true
scroll_to_top = true
skip_classes = ["no-transition", "external-link"]
skip_selectors = ["[data-skip-transition]", ".modal a"]
```

### YAML (`markata.yaml`)

```yaml
markata-go:
  view_transitions:
    enabled: true
    debug: false
    duration: 120
    update_meta: true
    scroll_to_top: true
    skip_classes:
      - no-transition
      - external-link
    skip_selectors:
      - "[data-skip-transition]"
      - ".modal a"
```

## Option Reference

### `enabled`

**Type**: Boolean  
**Default**: `true`

Enable or disable view transitions globally.

```toml
[markata-go.view_transitions]
enabled = false  # Disable view transitions completely
```

When disabled, all navigation uses standard browser behavior (no transitions).

**Use case**: Temporarily disable for testing or performance debugging.

---

### `debug`

**Type**: Boolean  
**Default**: `false`

Enable debug logging to browser console.

```toml
[markata-go.view_transitions]
debug = true
```

**Console output when enabled:**

```
View Transitions API initialized
Starting view transition to: /blog/my-post/
Prefetching URL: /blog/next-post/
\[view-transitions\] timing { fetchMs: 18.4, parseMs: 4.2, swapMs: 7.3, reinitMs: 5.1, totalMs: 42.6 }
View transition completed
Skipping link with class: no-transition
```

**Useful for:**
- Debugging transition issues
- Seeing which links are being intercepted
- Inspecting fetch/parse/swap/reinit timing in the browser console
- Understanding why certain links are skipped

---

### `duration`

**Type**: Integer (milliseconds)  
**Default**: `120`

Default transition duration in milliseconds. Used as a reference value.

```toml
[markata-go.view_transitions]
duration = 220  # Slower transitions
```

**Note**: The actual animation duration is controlled by CSS:

```css
::view-transition-old(main-content),
::view-transition-new(main-content) {
  animation-duration: 0.5s;  /* Match config.duration / 1000 */
}
```

---

### `update_meta`

**Type**: Boolean  
**Default**: `true`

Update meta tags (description, Open Graph, Twitter cards, canonical links) when navigating.

```toml
[view_transitions]
update_meta = false  # Skip meta tag updates
```

**When `true`**, updates:
- `<meta name="description">`
- `<meta property="og:*">` (Open Graph tags)
- `<meta name="twitter:*">` (Twitter Card tags)
- `<link rel="canonical">`

**When `false`**, only updates:
- Page title
- Body content

**Use case**: Set to `false` for slightly better performance if you don't rely on meta tags changing.

---

### `scroll_to_top`

**Type**: Boolean  
**Default**: `true`

Automatically scroll to top of page after navigation (unless URL has a hash).

```toml
[view_transitions]
scroll_to_top = false  # Keep scroll position
```

**Behavior:**
- When `true`: Scrolls to top (`window.scrollTo(0, 0)`)
- If URL has hash (`#section`): Smooth scrolls to that element (ignores this setting)
- When `false`: Maintains current scroll position

**Use case**: Set to `false` for documentation sites where users may want to maintain context.

---

### `skip_classes`

**Type**: Array of strings  
**Default**: `[]`

Additional CSS classes to skip (won't trigger view transitions).

```toml
[view_transitions]
skip_classes = ["no-transition", "instant-nav", "legacy-link"]
```

**Usage in HTML:**

```html
<a href="/page/" class="no-transition">Skip transition</a>
<a href="/instant/" class="instant-nav">Instant navigation</a>
```

**Built-in skip classes** (always skipped):
- `.toc-link` - Table of Contents links (use smooth scroll instead)

---

### `skip_selectors`

**Type**: Array of strings  
**Default**: `[]`

Additional CSS selectors to skip using `element.matches()`.

```toml
[view_transitions]
skip_selectors = [
  "[data-no-transition]",
  ".modal a",
  "#sidebar nav a",
  "[rel~='nofollow']"
]
```

**Usage in HTML:**

```html
<a href="/page/" data-no-transition>No transition</a>

<div class="modal">
  <a href="/help/">Won't transition (inside modal)</a>
</div>

<nav id="sidebar">
  <a href="/nav/">Won't transition (in sidebar)</a>
</nav>
```

**Built-in skip selectors** (always skipped):
- `a[target="_blank"]` - External links
- `a[download]` - Download links
- `a[rel~="external"]` - External rel links
- `a[hx-get]`, `a[hx-post]` - HTMX-managed links
- `a[data-no-transition]` - Explicit opt-out

**Built-in skip behavior** (always skipped):
- Non-HTML resources (any URL with a non-HTML file extension), including `.md`, `.txt`, `.xml`, `.json`

---

## Common Use Cases

### Development Mode

Enable detailed logging during development:

```toml
[view_transitions]
enabled = true
debug = true
```

### Skip Modal/Overlay Links

Prevent transitions for links inside modals:

```toml
[view_transitions]
skip_selectors = [
  ".modal a",
  ".overlay a",
  "[role='dialog'] a"
]
```

### Disable Meta Updates for Performance

Skip meta tag updates for slightly faster transitions:

```toml
[view_transitions]
update_meta = false
```

### Keep Scroll Position (Documentation Sites)

Maintain scroll position when navigating docs:

```toml
[view_transitions]
scroll_to_top = false
```

### Skip Transitions for External-Looking Links

```toml
[view_transitions]
skip_classes = ["external", "outbound"]
skip_selectors = ["[rel~='external']", "[target='_blank']"]
```

## Runtime Configuration

Inspect or modify configuration at runtime via browser console:

```javascript
// View current config
console.log(window.VIEW_TRANSITIONS_CONFIG);

// Temporarily disable
window.VIEW_TRANSITIONS_CONFIG.enabled = false;

// Enable debug mode
window.VIEW_TRANSITIONS_CONFIG.debug = true;

// Add skip class
window.VIEW_TRANSITIONS_CONFIG.skipClasses.push('my-class');
```

**Note**: Runtime changes only affect the current page. Reload to reset.

## Per-Link Opt-Out

Disable transitions for individual links without changing config:

```html
<!-- Using data attribute (always works) -->
<a href="/page/" data-no-transition>No transition</a>

<!-- Using CSS class (requires config) -->
<a href="/page/" class="no-transition">No transition</a>
```

## Disabling Temporarily

For testing or debugging, disable in browser console:

```javascript
// Option 1: Disable via config
window.VIEW_TRANSITIONS_CONFIG.enabled = false;

// Option 2: Undefine the API
document.startViewTransition = undefined;

// Option 3: Reload page
location.reload();
```

## Default Behavior

If you don't add `[view_transitions]` to your config, these defaults apply:

```toml
[markata-go.view_transitions]
enabled = true              # View transitions ON by default
debug = false               # No console logging
duration = 120              # 120ms reference duration
update_meta = true          # Update meta tags
scroll_to_top = true        # Scroll to top on navigation
skip_classes = []           # No custom skip classes
skip_selectors = []         # No custom skip selectors
```

## Minimal Configuration

The simplest valid config (uses all defaults):

```toml
[view_transitions]
enabled = true
```

Or just omit the section entirely - view transitions are enabled by default!

## Troubleshooting

### "View Transitions not working"

1. Check browser support (Chrome 111+, Safari 18+, Edge 111+)
2. Check console for errors with `debug = true`
3. Verify `enabled = true` in config
4. Check `window.VIEW_TRANSITIONS_CONFIG` in console

### "Specific links not transitioning"

1. Check if link matches `skip_selectors` or `skip_classes`
2. Check if link is external (`target="_blank"`)
3. Check if link has HTMX attributes (`hx-get`, etc.)
4. Enable `debug = true` to see why link is skipped

### "Navigation feels sluggish"

1. Enable `debug = true` and watch the `fetchMs`, `parseMs`, `swapMs`, and `reinitMs` timing log
2. Check whether the transition was prefetched (`prefetched: true`)
3. Use the browser Network tab to confirm HTML responses are served from cache when expected
4. Profile re-initialization listeners after `view-transition-complete` if `reinitMs` is consistently high
5. If `swapMs` is high, check whether the page is replacing more DOM than necessary

### "Scripts not working after transition"

Make sure scripts expose init functions and listen for `view-transition-complete`:

```javascript
function init() {
  // Your initialization
}

init(); // Initial load
window.addEventListener('view-transition-complete', init); // After transitions
```

## CSS Animation Customization

While config controls behavior, CSS controls visual appearance:

```css
/* Faster transitions */
::view-transition-old(main-content),
::view-transition-new(main-content) {
  animation-duration: 0.2s;
}

/* Different animation per element */
::view-transition-new(main-content) {
  animation: fade-in 0.3s, slide-up 0.3s;
}

::view-transition-new(site-nav) {
  animation: none;  /* Nav doesn't animate */
}
```

See [[view-transitions|View Transitions Guide]] for CSS customization examples.

## Related Documentation

- [[view-transitions|View Transitions Guide]] - Main documentation
- [[performance|Performance]] - Performance tips
- [[configuration|Configuration]] - General markata-go configuration

## Example Configurations

### Minimal (All Defaults)

```toml
# No config needed - view transitions enabled by default
```

### Development

```toml
[view_transitions]
enabled = true
debug = true
```

### Production (Performance Optimized)

```toml
[view_transitions]
enabled = true
debug = false
update_meta = false
duration = 250
```

### Documentation Site

```toml
[view_transitions]
enabled = true
scroll_to_top = false
skip_selectors = [".toc a", ".sidebar a"]
```

### With Skip Rules

```toml
[view_transitions]
enabled = true
skip_classes = ["instant", "no-anim"]
skip_selectors = [
  ".modal a",
  "[data-external]",
  "#admin-nav a"
]
```

### Disabled

```toml
[view_transitions]
enabled = false
```
