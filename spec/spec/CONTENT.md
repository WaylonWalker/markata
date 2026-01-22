# Markdown & Content Processing Specification

This document specifies how markdown content is processed.

## Processing Pipeline

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Source    │ →  │  Frontmatter│ →  │  Pre-render │ →  │   Render    │
│   File      │    │  Extraction │    │  (Jinja)    │    │  (Markdown) │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                                                │
                                                                ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Output    │ ←  │   Template  │ ←  │ Post-render │ ←  │    HTML     │
│   File      │    │   Wrap      │    │  (enhance)  │    │   Content   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

---

## Frontmatter Extraction

### Format

YAML frontmatter between `---` delimiters:

```markdown
---
title: My Post Title
date: 2024-01-15
tags:
  - python
  - tutorial
published: true
custom_field: any value
---

Content starts here...
```

### Parsing Rules

1. File MUST start with `---` on line 1 (or be considered no frontmatter)
2. Find closing `---` on its own line
3. Parse content between as YAML
4. Everything after closing `---` is content

### YAML Features Supported

| Feature | Example | Support |
|---------|---------|---------|
| Strings | `title: Hello` | Required |
| Numbers | `count: 42` | Required |
| Booleans | `published: true` | Required |
| Dates | `date: 2024-01-15` | Required |
| Lists | `tags: [a, b, c]` | Required |
| Block lists | `tags:\n  - a\n  - b` | Required |
| Objects | `author:\n  name: John` | Required |
| Multiline | `desc: \|` | Recommended |
| Null | `subtitle: null` or `subtitle:` | Required |

### Boolean Values

Accept common boolean representations:

| True | False |
|------|-------|
| `true` | `false` |
| `yes` | `no` |
| `on` | `off` |
| `True` | `False` |

### Date Values

Accept common date formats:

| Format | Example |
|--------|---------|
| ISO 8601 | `2024-01-15` |
| ISO 8601 datetime | `2024-01-15T10:30:00` |
| With timezone | `2024-01-15T10:30:00-05:00` |

### Error Handling

| Error | Behavior |
|-------|----------|
| Invalid YAML syntax | Error with file path and line number |
| Unknown field | Ignore (or store in `extra`) |
| Type mismatch | Attempt coercion, error if fails |
| Missing required | Error at validation stage |

---

## Pre-Render Processing

### Template Expressions in Markdown

When `jinja: true` in frontmatter, content is processed as a Jinja template:

```markdown
---
title: All Posts Index
jinja: true
---

# All Posts

Total: {{ core.posts | length }} posts

## Recent

{% for post in core.filter("published == True")[:5] %}
- [{{ post.title }}]({{ post.href }}) - {{ post.date }}
{% endfor %}

## By Tag

{% for tag in all_tags %}
### {{ tag }}
{% for post in core.filter("'" ~ tag ~ "' in tags") %}
- {{ post.title }}
{% endfor %}
{% endfor %}
```

### Template Context

Available variables in content templates:

| Variable | Type | Description |
|----------|------|-------------|
| `post` | Post | Current post being rendered |
| `core` | Core | Core instance with filter/map |
| `config` | Config | Site configuration |
| `today` | date | Current date |
| `now` | datetime | Current datetime |

### Filters

Standard Jinja filters plus:

| Filter | Description | Example |
|--------|-------------|---------|
| `length` | Count items | `{{ posts \| length }}` |
| `first` | First item | `{{ posts \| first }}` |
| `last` | Last item | `{{ posts \| last }}` |
| `sort` | Sort list | `{{ posts \| sort(attribute='date') }}` |
| `reverse` | Reverse list | `{{ posts \| reverse }}` |
| `selectattr` | Filter by attr | `{{ posts \| selectattr('published') }}` |
| `map` | Extract field | `{{ posts \| map(attribute='title') }}` |
| `join` | Join strings | `{{ tags \| join(', ') }}` |
| `default` | Default value | `{{ subtitle \| default('No subtitle') }}` |
| `upper` | Uppercase | `{{ title \| upper }}` |
| `lower` | Lowercase | `{{ title \| lower }}` |
| `title` | Title case | `{{ title \| title }}` |

### Disabling Jinja

Set `jinja: false` to disable template processing:

```markdown
---
title: Jinja Tutorial
jinja: false
---

Here's how to use Jinja: `{{ variable }}`

The above will render literally, not as a template.
```

---

## Markdown Rendering

### Required Features

Every implementation MUST support:

| Feature | Syntax |
|---------|--------|
| Headings | `# H1`, `## H2`, etc. |
| Paragraphs | Blank line separated |
| Emphasis | `*italic*`, `_italic_` |
| Strong | `**bold**`, `__bold__` |
| Links | `[text](url)` |
| Images | `![alt](url)` |
| Code (inline) | `` `code` `` |
| Code (block) | ` ```lang ` |
| Blockquotes | `> quote` |
| Unordered lists | `- item` or `* item` |
| Ordered lists | `1. item` |
| Horizontal rules | `---` or `***` |

### Extended Features

Implementations SHOULD support:

| Feature | Syntax | Output |
|---------|--------|--------|
| Tables | GFM table syntax | `<table>` |
| Strikethrough | `~~text~~` | `<del>text</del>` |
| Task lists | `- [ ] todo` | Checkbox |
| Footnotes | `[^1]` | Footnote |
| Heading IDs | `## Title {#custom-id}` | `<h2 id="custom-id">` |

### Tables

```markdown
| Header 1 | Header 2 | Header 3 |
|----------|:--------:|---------:|
| Left     | Center   | Right    |
| Cell     | Cell     | Cell     |
```

Output:
```html
<table>
  <thead>
    <tr>
      <th>Header 1</th>
      <th style="text-align: center">Header 2</th>
      <th style="text-align: right">Header 3</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Left</td>
      <td style="text-align: center">Center</td>
      <td style="text-align: right">Right</td>
    </tr>
  </tbody>
</table>
```

### Code Blocks

````markdown
```python
def hello():
    print("world")
```
````

Output:
```html
<pre><code class="language-python">def hello():
    print("world")
</code></pre>
```

With syntax highlighting (optional):
```html
<pre><code class="language-python hljs">
<span class="hljs-keyword">def</span> <span class="hljs-title function_">hello</span>():
    <span class="hljs-built_in">print</span>(<span class="hljs-string">"world"</span>)
</code></pre>
```

---

## Admonitions

### Syntax

```markdown
!!! note "Optional Title"
    Admonition content here.
    Can span multiple lines.

!!! warning
    Warning without custom title uses type as title.

!!! danger "Critical Issue"
    Danger admonition.
```

### Types

| Type | Default Title | Typical Color |
|------|---------------|---------------|
| `note` | Note | Blue |
| `info` | Info | Blue |
| `tip` | Tip | Green |
| `hint` | Hint | Green |
| `success` | Success | Green |
| `warning` | Warning | Yellow/Orange |
| `caution` | Caution | Yellow/Orange |
| `danger` | Danger | Red |
| `error` | Error | Red |
| `bug` | Bug | Red |
| `example` | Example | Purple |
| `quote` | Quote | Gray |
| `abstract` | Abstract | Cyan |
| `aside` | (none) | Gray (sidebar) |

### Output

```html
<div class="admonition note">
  <p class="admonition-title">Optional Title</p>
  <p>Admonition content here.
  Can span multiple lines.</p>
</div>
```

### Collapsible (Optional)

```markdown
??? note "Collapsed by default"
    Hidden content.

???+ note "Expanded by default"
    Visible content.
```

### Aside (Sidebar/Marginal Note)

The `aside` admonition creates a sidebar or marginal note that floats to the side of the main content. Useful for supplementary information, definitions, or tangential points.

#### Basic Aside

```markdown
!!! aside
    This is a marginal note that floats to the side.
```

#### Aside with Title

```markdown
!!! aside "Definition"
    A **static site generator** converts source files into static HTML.
```

#### Inline Modifiers

Control positioning with inline modifiers:

```markdown
!!! aside inline
    Floats to the left of content.

!!! aside inline end
    Floats to the right of content (default).
```

#### Output

```html
<aside class="admonition aside aside-inline-end">
  <p class="admonition-title">Definition</p>
  <p>A <strong>static site generator</strong> converts source files into static HTML.</p>
</aside>
```

Asides are best used sparingly for:
- Definitions and glossary terms
- Tangential commentary
- References and citations
- "Did you know?" facts

---

## Conversations (Chat)

The `chat` component renders conversation-style message bubbles, similar to messaging apps. Useful for tutorials showing dialogues, interview transcripts, or example conversations.

### Basic Syntax

```markdown
!!! chat
    !!! chat-left "Alice"
        Hi! How are you?
    
    !!! chat-right "Bob"
        I'm doing great, thanks for asking!
    
    !!! chat-left "Alice"
        Want to grab coffee later?
```

### With Timestamps

```markdown
!!! chat
    !!! chat-left "Alice" "10:30 AM"
        Hi! How are you?
    
    !!! chat-right "Bob" "10:32 AM"
        I'm doing great, thanks for asking!
```

### With Avatars

Configure avatars via frontmatter or inline:

```markdown
---
chat_avatars:
  Alice: /images/alice.jpg
  Bob: /images/bob.jpg
---

!!! chat
    !!! chat-left "Alice"
        Message with avatar from frontmatter config.
```

Or inline:

```markdown
!!! chat-left "Alice" avatar="/images/alice.jpg"
    Message with inline avatar.
```

### System Messages

For system notifications within a conversation:

```markdown
!!! chat
    !!! chat-system
        Alice has joined the chat
    
    !!! chat-left "Alice"
        Hello everyone!
```

### Output

```html
<div class="chat-container">
  <div class="chat-message chat-left">
    <div class="chat-avatar" style="background-image: url('/images/alice.jpg')"></div>
    <div class="chat-bubble">
      <div class="chat-author">Alice</div>
      <div class="chat-content">Hi! How are you?</div>
      <div class="chat-timestamp">10:30 AM</div>
    </div>
  </div>
  
  <div class="chat-message chat-right">
    <div class="chat-bubble">
      <div class="chat-author">Bob</div>
      <div class="chat-content">I'm doing great, thanks for asking!</div>
      <div class="chat-timestamp">10:32 AM</div>
    </div>
    <div class="chat-avatar" style="background-image: url('/images/bob.jpg')"></div>
  </div>
  
  <div class="chat-message chat-system">
    <div class="chat-content">Alice has joined the chat</div>
  </div>
</div>
```

### Use Cases

- **Tutorials**: Show example conversations with APIs or chatbots
- **Documentation**: Illustrate command-line interactions
- **Interviews**: Format Q&A content
- **Stories**: Dialogue-heavy narrative content
- **Support docs**: Example customer support conversations

---

## Internal Links (Wikilinks)

### Syntax

```markdown
Link to another post: [[other-post-slug]]

With custom text: [[other-post-slug|Click here]]
```

### Resolution

1. Find post where `slug == link_target`
2. If found, render as `<a href="{post.href}">{text}</a>`
3. If not found, leave as literal `[[link]]` and warn

### Output

```html
<!-- Found -->
<a href="/other-post-slug/">Other Post Title</a>
<a href="/other-post-slug/">Click here</a>

<!-- Not found -->
[[nonexistent-post]]
```

---

## Post-Render Enhancements

### Heading Anchors

Add anchor links to headings:

```html
<!-- Before -->
<h2>My Section</h2>

<!-- After -->
<h2 id="my-section">
  My Section
  <a href="#my-section" class="heading-anchor">#</a>
</h2>
```

### ID Generation

| Heading | ID |
|---------|-----|
| `## Hello World` | `hello-world` |
| `## What's New?` | `whats-new` |
| `## 2024 Updates` | `2024-updates` |
| `## FAQ` | `faq` |

Handle duplicates by appending numbers:
- `## FAQ` → `faq`
- `## FAQ` (second) → `faq-1`
- `## FAQ` (third) → `faq-2`

### Table of Contents

Generate TOC from headings:

```html
<nav class="toc">
  <ul>
    <li><a href="#introduction">Introduction</a></li>
    <li>
      <a href="#main-content">Main Content</a>
      <ul>
        <li><a href="#subsection">Subsection</a></li>
      </ul>
    </li>
    <li><a href="#conclusion">Conclusion</a></li>
  </ul>
</nav>
```

### Image Processing (Optional)

Add lazy loading:
```html
<img src="photo.jpg" alt="Photo" loading="lazy" />
```

Add dimensions:
```html
<img src="photo.jpg" alt="Photo" width="800" height="600" />
```

---

## Configuration

```toml
[tool-name.markdown]
# Backend library (if supporting multiple)
backend = "markdown-it-py"

# Extensions to enable
extensions = [
    "tables",
    "admonitions",
    "footnotes",
    "syntax_highlight",
]

# Syntax highlighting
[tool-name.markdown.highlight]
enabled = true
theme = "github-dark"
line_numbers = false

# Heading anchors
[tool-name.markdown.anchors]
enabled = true
position = "end"  # "start" or "end"
symbol = "#"

# Admonitions
[tool-name.markdown.admonitions]
enabled = true
collapsible = true
```

---

## Library Recommendations

### Python

| Library | Notes |
|---------|-------|
| `markdown-it-py` | CommonMark, extensible, fast |
| `markdown` | Python-Markdown, many extensions |
| `mistune` | Fast, simple |

### JavaScript

| Library | Notes |
|---------|-------|
| `markdown-it` | CommonMark, pluggable |
| `marked` | Fast, popular |
| `remark` | AST-based, unified ecosystem |

### Go

| Library | Notes |
|---------|-------|
| `goldmark` | CommonMark, extensible |
| `blackfriday` | Fast, feature-rich |

### Rust

| Library | Notes |
|---------|-------|
| `pulldown-cmark` | CommonMark, fast |
| `comrak` | GFM compatible |

---

## Error Handling

| Error | Behavior |
|-------|----------|
| Invalid markdown | Render best-effort, don't error |
| Unclosed code block | Include rest of file in block |
| Invalid admonition | Render as paragraph |
| Broken internal link | Render as literal, warn |
| Unknown language | Render without highlighting |

---

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [THEMES.md](./THEMES.md) - Admonition and code block styling
- [CONFIG.md](./CONFIG.md) - Markdown configuration
- [TEMPLATES.md](./TEMPLATES.md) - Template system
- [DATA_MODEL.md](./DATA_MODEL.md) - Post model
