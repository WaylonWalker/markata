# Themes and Customization Specification

Themes control the visual appearance of the generated site. The system supports:
- **Built-in themes** - Ship with usable defaults
- **Theme packages** - Complete themes with templates, CSS, and assets  
- **User customization** - Override any theme file locally

## Design Principles

1. **Zero-config beautiful** - Sites look good with no customization
2. **Progressive customization** - Pick colors → override CSS → replace templates
3. **Theme packages** - Installable, shareable, complete visual identities
4. **Local overrides** - Any theme file can be overridden by placing it in your project

## Integration with Head/Style System

Themes provide the base visual identity. For additional customization, see [HEAD_STYLE.md](./HEAD_STYLE.md):

| System | Purpose | Example |
|--------|---------|---------|
| Theme (`[name.theme]`) | Base templates, CSS, layouts | Selecting "blog" theme |
| Theme Variables (`[name.theme.variables]`) | Quick CSS property overrides | `--color-primary: #8b5cf6` |
| Head Config (`[name.head]`) | Meta tags, scripts, links | Analytics, fonts, OG tags |
| Style Config (`[name.style]`) | Legacy color overrides | `color_bg`, `color_text` |
| Post Overrides (`config_overrides`) | Per-post customization | Custom scripts for one page |

**Resolution order:** Theme defaults → Theme variables → Head/Style config → Post overrides

---

## Theme Structure

A theme is a directory containing templates, CSS, and optionally assets:

```
themes/
└── default/
    ├── theme.toml           # Theme metadata
    ├── templates/
    │   ├── base.html        # Base layout
    │   ├── post.html        # Single post
    │   ├── feed.html        # Feed/index page
    │   ├── card.html        # Post card for feeds
    │   └── partials/
    │       ├── head.html
    │       ├── header.html
    │       ├── footer.html
    │       └── pagination.html
    ├── static/
    │   ├── css/
    │   │   ├── main.css     # Core styles
    │   │   ├── admonitions.css
    │   │   ├── code.css     # Syntax highlighting
    │   │   └── variables.css # CSS custom properties
    │   └── js/
    │       └── main.js      # Optional JavaScript
    └── assets/              # Theme-specific images, fonts
```

### theme.toml

```toml
[theme]
name = "Default"
version = "1.0.0"
description = "Clean, minimal theme with dark mode support"
author = "SSG Team"
license = "MIT"
homepage = "https://github.com/example/ssg-theme-default"

# Minimum SSG version required
min_version = "1.0.0"

# Theme supports these features
[theme.features]
dark_mode = true
responsive = true
syntax_highlighting = true
admonitions = true

# Configurable options this theme exposes
[theme.options]
primary_color = { type = "color", default = "#3b82f6", description = "Primary accent color" }
font_family = { type = "string", default = "system-ui", description = "Body font family" }
max_width = { type = "string", default = "65ch", description = "Content max width" }
show_toc = { type = "boolean", default = true, description = "Show table of contents" }
show_reading_time = { type = "boolean", default = true, description = "Show reading time" }

# Parent theme to inherit from (optional)
[theme.extends]
theme = "base"  # Inherit from another theme
```

---

## Theme Resolution

When looking for a template or static file, the system searches in order:

```
1. Project local:     ./templates/post.html
2. Project theme:     ./themes/[theme]/templates/post.html  
3. Installed theme:   ~/.config/[name]/themes/[theme]/templates/post.html
4. Built-in theme:    [internal]/themes/[theme]/templates/post.html
5. Default theme:     [internal]/themes/default/templates/post.html
```

This allows users to override any theme file by creating it locally.

### Override Examples

**Override just the footer:**
```
my-site/
├── templates/
│   └── partials/
│       └── footer.html    # Your custom footer
├── posts/
└── [name].toml
```

**Override a CSS file:**
```
my-site/
├── static/
│   └── css/
│       └── admonitions.css  # Your custom admonition styles
├── posts/
└── [name].toml
```

---

## Configuration

### Selecting a Theme

```toml
[name.theme]
name = "default"           # Theme name
```

### Customizing Theme Options

Themes expose configurable options:

```toml
[name.theme]
name = "default"

# Theme-specific options (defined in theme.toml)
[name.theme.options]
primary_color = "#8b5cf6"   # Purple instead of blue
font_family = "Inter, system-ui"
max_width = "70ch"
show_toc = false
```

### CSS Variable Overrides

For quick color/font changes without modifying CSS files:

```toml
[name.theme]
name = "default"

# Override CSS custom properties
[name.theme.variables]
"--color-primary" = "#8b5cf6"
"--color-primary-dark" = "#7c3aed"
"--font-body" = "Inter, system-ui"
"--font-heading" = "Inter, system-ui"
"--font-mono" = "JetBrains Mono, monospace"
"--content-width" = "70ch"
"--radius" = "0.5rem"
```

These are injected as a `<style>` block in the `<head>`:

```html
<style>
:root {
  --color-primary: #8b5cf6;
  --color-primary-dark: #7c3aed;
  --font-body: Inter, system-ui;
  /* ... */
}
</style>
```

### Custom CSS File

For more extensive customization:

```toml
[name.theme]
name = "default"
custom_css = "my-styles.css"  # Loaded after theme CSS
```

---

## Built-in Themes

Implementations MUST provide at least the `default` theme.

### Default Theme

A clean, minimal theme with:
- Responsive layout
- Dark mode support (prefers-color-scheme)
- Styled admonitions
- Syntax highlighting
- Typography optimized for reading

### Minimal Theme (Optional)

Bare-bones HTML with minimal styling for:
- Maximum customization flexibility
- Fast loading
- Print-friendly

---

## CSS Custom Properties

Built-in themes SHOULD use CSS custom properties for consistency:

### Colors

```css
:root {
  /* Primary brand color */
  --color-primary: #3b82f6;
  --color-primary-light: #60a5fa;
  --color-primary-dark: #2563eb;
  
  /* Semantic colors */
  --color-text: #1f2937;
  --color-text-muted: #6b7280;
  --color-background: #ffffff;
  --color-surface: #f9fafb;
  --color-border: #e5e7eb;
  
  /* Status colors */
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --color-info: #3b82f6;
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  :root {
    --color-text: #f9fafb;
    --color-text-muted: #9ca3af;
    --color-background: #111827;
    --color-surface: #1f2937;
    --color-border: #374151;
  }
}
```

### Typography

```css
:root {
  /* Font families */
  --font-body: system-ui, -apple-system, sans-serif;
  --font-heading: var(--font-body);
  --font-mono: ui-monospace, 'Cascadia Code', 'Fira Code', monospace;
  
  /* Font sizes (modular scale) */
  --text-xs: 0.75rem;
  --text-sm: 0.875rem;
  --text-base: 1rem;
  --text-lg: 1.125rem;
  --text-xl: 1.25rem;
  --text-2xl: 1.5rem;
  --text-3xl: 1.875rem;
  --text-4xl: 2.25rem;
  
  /* Line heights */
  --leading-tight: 1.25;
  --leading-normal: 1.5;
  --leading-relaxed: 1.75;
}
```

### Spacing

```css
:root {
  /* Spacing scale */
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;
  --space-12: 3rem;
  --space-16: 4rem;
  
  /* Layout */
  --content-width: 65ch;
  --page-width: 1200px;
  --radius: 0.375rem;
  --radius-lg: 0.5rem;
}
```

---

## Admonition Styles

Built-in themes MUST include styles for all admonition types.

### Admonition CSS

```css
/* Base admonition styles */
.admonition {
  margin: var(--space-4) 0;
  padding: var(--space-4);
  border-left: 4px solid var(--admonition-color, var(--color-primary));
  background: var(--admonition-bg, var(--color-surface));
  border-radius: 0 var(--radius) var(--radius) 0;
}

.admonition-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-weight: 600;
  margin-bottom: var(--space-2);
  color: var(--admonition-color, var(--color-primary));
}

.admonition-title::before {
  content: var(--admonition-icon, "");
  font-size: 1.25em;
}

/* Type-specific styles */
.admonition.note {
  --admonition-color: #3b82f6;
  --admonition-bg: #eff6ff;
  --admonition-icon: "ℹ️";
}

.admonition.info {
  --admonition-color: #3b82f6;
  --admonition-bg: #eff6ff;
  --admonition-icon: "ℹ️";
}

.admonition.tip {
  --admonition-color: #10b981;
  --admonition-bg: #ecfdf5;
  --admonition-icon: "💡";
}

.admonition.hint {
  --admonition-color: #10b981;
  --admonition-bg: #ecfdf5;
  --admonition-icon: "💡";
}

.admonition.success {
  --admonition-color: #10b981;
  --admonition-bg: #ecfdf5;
  --admonition-icon: "✅";
}

.admonition.warning {
  --admonition-color: #f59e0b;
  --admonition-bg: #fffbeb;
  --admonition-icon: "⚠️";
}

.admonition.caution {
  --admonition-color: #f59e0b;
  --admonition-bg: #fffbeb;
  --admonition-icon: "⚠️";
}

.admonition.danger {
  --admonition-color: #ef4444;
  --admonition-bg: #fef2f2;
  --admonition-icon: "🚨";
}

.admonition.error {
  --admonition-color: #ef4444;
  --admonition-bg: #fef2f2;
  --admonition-icon: "❌";
}

.admonition.bug {
  --admonition-color: #ef4444;
  --admonition-bg: #fef2f2;
  --admonition-icon: "🐛";
}

.admonition.example {
  --admonition-color: #8b5cf6;
  --admonition-bg: #f5f3ff;
  --admonition-icon: "📝";
}

.admonition.quote {
  --admonition-color: #6b7280;
  --admonition-bg: #f9fafb;
  --admonition-icon: "💬";
}

.admonition.abstract {
  --admonition-color: #06b6d4;
  --admonition-bg: #ecfeff;
  --admonition-icon: "📋";
}

/* Dark mode adjustments */
@media (prefers-color-scheme: dark) {
  .admonition.note { --admonition-bg: #1e3a5f; }
  .admonition.info { --admonition-bg: #1e3a5f; }
  .admonition.tip { --admonition-bg: #064e3b; }
  .admonition.hint { --admonition-bg: #064e3b; }
  .admonition.success { --admonition-bg: #064e3b; }
  .admonition.warning { --admonition-bg: #451a03; }
  .admonition.caution { --admonition-bg: #451a03; }
  .admonition.danger { --admonition-bg: #450a0a; }
  .admonition.error { --admonition-bg: #450a0a; }
  .admonition.bug { --admonition-bg: #450a0a; }
  .admonition.example { --admonition-bg: #2e1065; }
  .admonition.quote { --admonition-bg: #1f2937; }
  .admonition.abstract { --admonition-bg: #083344; }
}
```

### Aside (Sidebar) Styles

```css
/* Aside - Sidebar/Marginal Note */
aside.admonition.aside {
  --admonition-color: #6b7280;
  --admonition-bg: #f9fafb;
  border-left: none;
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  font-size: var(--text-sm);
  max-width: 280px;
}

/* Default: float right */
aside.admonition.aside {
  float: right;
  margin: 0 0 var(--space-4) var(--space-4);
  clear: right;
}

/* Inline left modifier */
aside.admonition.aside.aside-inline {
  float: left;
  margin: 0 var(--space-4) var(--space-4) 0;
  clear: left;
}

/* Inline end (right) modifier - explicit */
aside.admonition.aside.aside-inline-end {
  float: right;
  margin: 0 0 var(--space-4) var(--space-4);
  clear: right;
}

/* Aside title styling */
aside.admonition.aside .admonition-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

aside.admonition.aside .admonition-title::before {
  content: none; /* No icon for asides */
}

/* Responsive: full width on small screens */
@media (max-width: 768px) {
  aside.admonition.aside,
  aside.admonition.aside.aside-inline,
  aside.admonition.aside.aside-inline-end {
    float: none;
    max-width: 100%;
    margin: var(--space-4) 0;
  }
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  aside.admonition.aside {
    --admonition-bg: #1f2937;
  }
}
```

### Chat/Conversation Styles

```css
/* Chat container */
.chat-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin: var(--space-4) 0;
  padding: var(--space-4);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
}

/* Individual message */
.chat-message {
  display: flex;
  align-items: flex-end;
  gap: var(--space-2);
  max-width: 85%;
}

/* Left-aligned messages */
.chat-message.chat-left {
  align-self: flex-start;
  flex-direction: row;
}

/* Right-aligned messages */
.chat-message.chat-right {
  align-self: flex-end;
  flex-direction: row-reverse;
}

/* Avatar */
.chat-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: var(--color-primary);
  background-size: cover;
  background-position: center;
  flex-shrink: 0;
}

/* Message bubble */
.chat-bubble {
  padding: var(--space-3);
  border-radius: var(--radius-lg);
  position: relative;
}

/* Left bubble styling */
.chat-left .chat-bubble {
  background: #e5e7eb;
  color: #1f2937;
  border-bottom-left-radius: var(--radius);
}

/* Right bubble styling */
.chat-right .chat-bubble {
  background: var(--color-primary);
  color: white;
  border-bottom-right-radius: var(--radius);
}

/* Author name */
.chat-author {
  font-size: var(--text-xs);
  font-weight: 600;
  margin-bottom: var(--space-1);
  opacity: 0.8;
}

.chat-left .chat-author {
  color: var(--color-text-muted);
}

.chat-right .chat-author {
  color: rgba(255, 255, 255, 0.8);
}

/* Message content */
.chat-content {
  font-size: var(--text-sm);
  line-height: var(--leading-normal);
}

.chat-content p {
  margin: 0;
}

.chat-content p + p {
  margin-top: var(--space-2);
}

/* Timestamp */
.chat-timestamp {
  font-size: var(--text-xs);
  margin-top: var(--space-1);
  opacity: 0.6;
}

.chat-left .chat-timestamp {
  color: var(--color-text-muted);
}

.chat-right .chat-timestamp {
  color: rgba(255, 255, 255, 0.7);
}

/* System messages */
.chat-message.chat-system {
  align-self: center;
  max-width: 100%;
}

.chat-system .chat-content {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  background: transparent;
  padding: var(--space-2) var(--space-4);
  text-align: center;
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  .chat-container {
    background: var(--color-surface);
  }
  
  .chat-left .chat-bubble {
    background: #374151;
    color: #f9fafb;
  }
  
  .chat-left .chat-author {
    color: #9ca3af;
  }
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .chat-message {
    max-width: 95%;
  }
  
  .chat-avatar {
    width: 28px;
    height: 28px;
  }
}
```

### Collapsible Admonitions

```css
/* Collapsible admonitions use <details>/<summary> */
details.admonition {
  border-left: 4px solid var(--admonition-color, var(--color-primary));
}

details.admonition > summary {
  cursor: pointer;
  list-style: none;
}

details.admonition > summary::before {
  content: "▶";
  display: inline-block;
  margin-right: var(--space-2);
  transition: transform 0.2s;
}

details.admonition[open] > summary::before {
  transform: rotate(90deg);
}
```

---

## Code Block Styles

### Base Code Styles

```css
/* Inline code */
code {
  font-family: var(--font-mono);
  font-size: 0.9em;
  padding: 0.125em 0.375em;
  background: var(--color-surface);
  border-radius: var(--radius);
}

/* Code blocks */
pre {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  line-height: var(--leading-relaxed);
  padding: var(--space-4);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  overflow-x: auto;
}

pre code {
  padding: 0;
  background: none;
}

/* Line numbers (optional) */
pre.line-numbers {
  padding-left: 3.5em;
  position: relative;
}

pre.line-numbers::before {
  content: attr(data-line-numbers);
  position: absolute;
  left: 0;
  top: var(--space-4);
  width: 3em;
  text-align: right;
  color: var(--color-text-muted);
  border-right: 1px solid var(--color-border);
  padding-right: var(--space-2);
  user-select: none;
}

/* Language label */
pre[data-language]::after {
  content: attr(data-language);
  position: absolute;
  top: var(--space-2);
  right: var(--space-2);
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  text-transform: uppercase;
}
```

### Syntax Highlighting Themes

Implementations SHOULD support multiple syntax highlighting themes:

| Theme | Description |
|-------|-------------|
| `github-light` | GitHub's light theme |
| `github-dark` | GitHub's dark theme |
| `monokai` | Classic dark theme |
| `one-dark` | Atom's One Dark |
| `dracula` | Dracula theme |
| `nord` | Nord color palette |
| `solarized-light` | Solarized light |
| `solarized-dark` | Solarized dark |

Configuration:

```toml
[name.markdown.highlight]
theme = "github-dark"
line_numbers = false
```

---

## Template Requirements

### Base Template (`base.html`)

Every theme MUST provide a base template with these blocks:

```jinja2
<!DOCTYPE html>
<html lang="{{ config.lang | default('en') }}">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  
  {% block meta %}
  <title>{% block title %}{{ config.title }}{% endblock %}</title>
  <meta name="description" content="{% block description %}{{ config.description }}{% endblock %}">
  {% endblock %}
  
  {% block head %}
  {# Theme CSS #}
  <link rel="stylesheet" href="{{ 'css/main.css' | theme_asset }}">
  
  {# CSS variable overrides from config #}
  {% if config.theme.variables %}
  <style>:root { {% for k, v in config.theme.variables.items() %}{{ k }}: {{ v }}; {% endfor %} }</style>
  {% endif %}
  
  {# Custom CSS #}
  {% if config.theme.custom_css %}
  <link rel="stylesheet" href="{{ config.theme.custom_css | asset_url }}">
  {% endif %}
  {% endblock %}
</head>
<body>
  {% block header %}
  {% include "partials/header.html" %}
  {% endblock %}
  
  <main>
    {% block content %}{% endblock %}
  </main>
  
  {% block footer %}
  {% include "partials/footer.html" %}
  {% endblock %}
  
  {% block scripts %}{% endblock %}
</body>
</html>
```

### Post Template (`post.html`)

```jinja2
{% extends "base.html" %}

{% block title %}{{ post.title }} | {{ config.title }}{% endblock %}
{% block description %}{{ post.description | default(post.content | truncate(160)) }}{% endblock %}

{% block meta %}
{{ super() }}
<meta property="og:title" content="{{ post.title }}">
<meta property="og:type" content="article">
<meta property="og:url" content="{{ post.absolute_url }}">
{% if post.date %}
<meta property="article:published_time" content="{{ post.date | isoformat }}">
{% endif %}
{% endblock %}

{% block content %}
<article class="post">
  <header class="post-header">
    <h1>{{ post.title }}</h1>
    {% if post.date or post.reading_time %}
    <div class="post-meta">
      {% if post.date %}<time datetime="{{ post.date | isoformat }}">{{ post.date | date }}</time>{% endif %}
      {% if post.reading_time %}<span>{{ post.reading_time }} min read</span>{% endif %}
    </div>
    {% endif %}
  </header>
  
  <div class="post-content">
    {{ body | safe }}
  </div>
  
  {% if post.tags %}
  <footer class="post-footer">
    <div class="tags">
      {% for tag in post.tags %}
      <a href="/tags/{{ tag | slugify }}/" class="tag">{{ tag }}</a>
      {% endfor %}
    </div>
  </footer>
  {% endif %}
</article>

{% if post.prev or post.next %}
<nav class="post-nav">
  {% if post.prev %}<a href="{{ post.prev.href }}" class="prev">← {{ post.prev.title }}</a>{% endif %}
  {% if post.next %}<a href="{{ post.next.href }}" class="next">{{ post.next.title }} →</a>{% endif %}
</nav>
{% endif %}
{% endblock %}
```

### Feed Template (`feed.html`)

```jinja2
{% extends "base.html" %}

{% block title %}{{ feed.title | default(config.title) }}{% endblock %}

{% block content %}
<div class="feed">
  {% if feed.title %}
  <header class="feed-header">
    <h1>{{ feed.title }}</h1>
    {% if feed.description %}<p>{{ feed.description }}</p>{% endif %}
  </header>
  {% endif %}
  
  <div class="posts">
    {% for post in feed.posts %}
    {% include "card.html" %}
    {% endfor %}
  </div>
  
  {% if feed.pagination.total_pages > 1 %}
  {% include "partials/pagination.html" %}
  {% endif %}
</div>
{% endblock %}
```

### Card Template (`card.html`)

```jinja2
<article class="card">
  <h2><a href="{{ post.href }}">{{ post.title }}</a></h2>
  {% if post.description %}
  <p class="card-description">{{ post.description }}</p>
  {% endif %}
  <div class="card-meta">
    {% if post.date %}<time datetime="{{ post.date | isoformat }}">{{ post.date | date }}</time>{% endif %}
    {% if post.reading_time %}<span>{{ post.reading_time }} min</span>{% endif %}
  </div>
</article>
```

---

## Template Filters

Themes should have access to these filters:

| Filter | Description | Example |
|--------|-------------|---------|
| `theme_asset` | URL to theme static file | `{{ 'css/main.css' \| theme_asset }}` |
| `asset_url` | URL to project static file | `{{ 'images/logo.png' \| asset_url }}` |
| `date` | Format date | `{{ post.date \| date('%B %d, %Y') }}` |
| `isoformat` | ISO 8601 date | `{{ post.date \| isoformat }}` |
| `truncate` | Truncate text | `{{ text \| truncate(160) }}` |
| `slugify` | Generate slug | `{{ title \| slugify }}` |
| `safe` | Mark as safe HTML | `{{ body \| safe }}` |

---

## Installing Themes

### From Package Manager

```bash
# Python
pip install [name]-theme-blog

# npm
npm install @[name]/theme-blog

# Direct download
[name] theme install https://github.com/user/theme-blog
```

### Manual Installation

Download theme to `~/.config/[name]/themes/` or `./themes/`:

```bash
mkdir -p themes
git clone https://github.com/user/theme-blog themes/blog
```

### Configuration

```toml
[name.theme]
name = "blog"  # Matches directory name
```

---

## Creating a Theme

### Scaffold a New Theme

```bash
[name] theme new my-theme
```

Creates:
```
themes/my-theme/
├── theme.toml
├── templates/
│   ├── base.html
│   ├── post.html
│   ├── feed.html
│   └── card.html
└── static/
    └── css/
        ├── main.css
        └── variables.css
```

### Extending a Theme

Inherit from an existing theme and override only what you need:

```toml
# themes/my-theme/theme.toml
[theme]
name = "My Theme"

[theme.extends]
theme = "default"
```

Then only include the files you want to override:

```
themes/my-theme/
├── theme.toml
├── templates/
│   └── partials/
│       └── footer.html    # Custom footer only
└── static/
    └── css/
        └── variables.css  # Custom colors only
```

---

## CLI Commands

### `theme list`

List available themes:

```bash
$ [name] theme list

Installed themes:
  default     Clean, minimal theme with dark mode (built-in)
  minimal     Bare-bones HTML for maximum customization (built-in)
  blog        Feature-rich blog theme (~/.config/[name]/themes/blog)
  
Current theme: default
```

### `theme info`

Show theme details:

```bash
$ [name] theme info blog

Name: Blog Theme
Version: 2.1.0
Author: Theme Author
License: MIT
Homepage: https://github.com/user/theme-blog

Features:
  ✓ Dark mode
  ✓ Responsive
  ✓ Syntax highlighting
  ✓ Admonitions

Options:
  primary_color   Color   #3b82f6   Primary accent color
  font_family     String  system-ui Body font family
  show_toc        Boolean true      Show table of contents
```

### `theme install`

Install a theme:

```bash
$ [name] theme install https://github.com/user/theme-blog
Installing theme from https://github.com/user/theme-blog...
Theme 'blog' installed to ~/.config/[name]/themes/blog
```

### `theme new`

Create a new theme:

```bash
$ [name] theme new my-theme
Created theme scaffold in themes/my-theme/
```

---

## See Also

- [TEMPLATES.md](./TEMPLATES.md) - Template system details
- [CONFIG.md](./CONFIG.md) - Theme configuration
- [CONTENT.md](./CONTENT.md) - Admonition syntax
- [SPEC.md](./SPEC.md) - Core specification
