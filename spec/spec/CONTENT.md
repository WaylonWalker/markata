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

## Expressive Inline Titles and Headings

Markdown inline semantics are supported in body headings and frontmatter
titles. The authored `Post.Title` remains the source string. The build derives
two explicit representations:

- `TitleHTML` is safe inline HTML for the visible page title.
- `TitleText` is semantic plain text for browser titles, feeds, search,
  structured data, social metadata, slugs, and other non-visual consumers.

HTML character references in authored titles are decoded before the title is
exposed as plain text. For example, `It&rsquo;s just the Carpet&hellip;` becomes
`It’s just the Carpet…` in `TitleText`, rather than exposing the entity names
to feeds and metadata. References are decoded after Markdown parsing so code
spans and escaped Markdown syntax keep their authored semantics.

The shared inline renderer uses the configured Goldmark extensions, including
GFM strikethrough, superscript/subscript, and Markata's `==mark==` extension. Mark is a delimiter
container, so nested constructs such as `==**strong**==`, ``==`code`==``, and
`==[link](/path)==` retain their child semantics. Raw HTML is not enabled in
title rendering, and title links use the same safe destination policy as body
Markdown.

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

### Frontmatter Aliases

The frontmatter loader normalizes alias fields for compatibility:

| Canonical | Accepted Aliases |
|-----------|------------------|
| `aliases` | `alias`, `handle`, `handles` |

All values are merged into the canonical `aliases` key as a list of strings. Existing keys are preserved.

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

### Heading Best Practices

**Avoid H1 headings in content.** Templates automatically generate an H1 from the frontmatter `title` field. Using H1 (`# Heading`) in markdown content creates duplicate H1 tags, which:

- Harms SEO (search engines expect a single H1 per page)
- Reduces accessibility (screen readers use H1 for page identification)
- Violates HTML document outline best practices

**Start content with H2 (`##`) or deeper headings.** The linter will warn if H1 headings are detected in content.

```markdown
---
title: My Page Title  # This becomes the H1
---

## First Section       # Start with H2, not H1

Content here...

### Subsection         # H3 for nested sections
```

### Extended Features

Implementations SHOULD support:

| Feature | Syntax | Output |
|---------|--------|--------|
| Tables | GFM table syntax | `<table>` |
| Strikethrough | `~~text~~` | `<del>text</del>` |
| Task lists | `- [ ] todo` | Checkbox |
| Footnotes | `[^1]` | Footnote |
| Heading IDs | `## Title {#custom-id}` | `<h2 id="custom-id">` |
| Attributes | `{.class}`, `{#id}` | Element with class/id |
| Smart Quotes | `"text"` | `"text"` (curly) |
| Definition Lists | `Term\n:   Def` | `<dl>`, `<dt>`, `<dd>` |

GFM tables with a header row MUST expose sortable header controls in rendered
post pages. Sorting MUST toggle between ascending and descending order, preserve
stable order for equal values, sort numeric columns numerically, and leave empty
values at the end. The sortable-table script MUST be deferred and included only
on pages whose rendered body contains a table with a header.

### Attribute Syntax

The attribute syntax `{...}` allows adding CSS classes, IDs, and other attributes to elements.

#### Block Elements (Headings, Paragraphs, Lists, Blockquotes, Tables)

```markdown
## My Section {.highlighted}

## Installation {#install}

## Features {#features .important}

Paragraph text
{.lead}

> Blockquote text
{.callout}

- List item
  {.item}

| Col |
|---|
| Cell {.cell} |
```

Output:
```html
<h2 class="highlighted" id="my-section">My Section</h2>
<h2 id="install">Installation</h2>
<h2 id="features" class="important">Features</h2>
<p class="lead">Paragraph text</p>
<blockquote class="callout">
  <p>Blockquote text</p>
</blockquote>
<ul>
  <li class="item">List item</li>
</ul>
<table>
  <tbody>
    <tr>
      <td class="cell">Cell</td>
    </tr>
  </tbody>
</table>
```

#### Inline Elements (Images, Links, Emphasis, Code, Strikethrough, Wikilinks)

```markdown
![alt text](image.webp){.more-cinematic}

![photo](photo.jpg){.shadow .bordered}

![hero](hero.png){#hero-image}

[Read more](url){.external-link}

*note*{.accent}

`code`{.inline-code}

~~deprecated~~{.deprecated}

[[hello-world]]{.wikilink-pill}
```

Output:
```html
<img src="image.webp" alt="alt text" class="more-cinematic">
<img src="photo.jpg" alt="photo" class="shadow bordered">
<img src="hero.png" alt="hero" id="hero-image">
<a href="url" class="external-link">Read more</a>
<em class="accent">note</em>
<code class="inline-code">code</code>
<del class="deprecated">deprecated</del>
<a href="/hello-world/" class="wikilink wikilink-pill">Hello World</a>
```

#### Supported Attribute Formats

| Syntax | Description | Example |
|--------|-------------|---------|
| `{.classname}` | CSS class | `{.highlight}` → `class="highlight"` |
| `{.class1 .class2}` | Multiple classes | `{.shadow .rounded}` → `class="shadow rounded"` |
| `{#idname}` | ID attribute | `{#hero}` → `id="hero"` |
| `{#id .class}` | Combined | `{#main .featured}` → `id="main" class="featured"` |
| `{key=value}` | Custom attribute | `{data-size=large}` → `data-size="large"` |

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

### Smart Quotes (Typographer)

Automatically converts straight quotes to typographic (curly) quotes and other punctuation:

| Input | Output | Description |
|-------|--------|-------------|
| `"text"` | `"text"` | Double curly quotes |
| `'text'` | `'text'` | Single curly quotes |
| `It's` | `It's` | Apostrophe |
| `9--5` | `9–5` | En dash |
| `hello---world` | `hello—world` | Em dash |
| `wait...` | `wait…` | Ellipsis |

**Note:** These are HTML entities (`&ldquo;`, `&rdquo;`, etc.) that render as proper typographic characters.

**Configuration:**
```toml
[markdown.extensions]
typographer = true  # Enable smart quotes (default: true)
```

### Definition Lists

PHP Markdown Extra style definition lists:

```markdown
Term 1
:   Definition 1

Term 2
:   Definition 2a
:   Definition 2b
```

Output:
```html
<dl>
  <dt>Term 1</dt>
  <dd>Definition 1</dd>
  <dt>Term 2</dt>
  <dd>Definition 2a</dd>
  <dd>Definition 2b</dd>
</dl>
```

**Configuration:**
```toml
[markdown.extensions]
definition_list = true  # Enable definition lists (default: true)
```

### Footnotes

PHP Markdown Extra style footnotes:

```markdown
Here's a sentence with a footnote.[^1]

[^1]: This is the footnote content.
```

Output:
```html
<p>Here's a sentence with a footnote.<sup><a href="#fn:1">1</a></sup></p>
<!-- ... later in document ... -->
<section class="footnotes">
  <ol>
    <li id="fn:1">
      <p>This is the footnote content. <a href="#fnref:1">↩</a></p>
    </li>
  </ol>
</section>
```

**Configuration:**
```toml
[markdown.extensions]
footnote = true  # Enable footnotes (default: true)
```

### CJK Line Breaks

Enable proper line breaking for Chinese, Japanese, and Korean text.

```markdown
これはテストです。これはテストです。
```

Output:
```html
<p>これはテストです。<br>これはテストです。</p>
```

**Configuration:**
```toml
[markdown.extensions]
cjk = true  # Enable CJK line breaks (default: true)
```

### Figures

Convert images with following paragraphs into `<figure>` elements with `<figcaption>`. A single-paragraph blockquote immediately following a figure SHOULD also be treated as caption content and merged into the figure.

```markdown
![Alt text](image.jpg)
This is the caption.
```

```markdown
![Alt text](image.jpg)
> This is also treated as the caption.
```

Output:
```html
<figure>
  <img src="image.jpg" alt="Alt text">
  <figcaption>This is the caption.</figcaption>
</figure>
```

**Configuration:**
```toml
[markdown.extensions]
figure = true  # Enable figures (default: true)
```

### Blockquote Attribution

Blockquotes followed by attribution text (no blank line between them) SHOULD be converted into `<figure class="quote">` elements with the attribution rendered as `<figcaption>`.

Supported attribution patterns:
- `@mention` — `@username`
- `[[wikilink]]` — `[[Page Name]]`
- `[text](url)` — Markdown links
- `— text` or `-- text` — Dash-prefixed text
- Short plain text (≤ 60 characters) — names or short attributions

The attribution text MUST directly follow the blockquote with no blank line. If a blank line separates them, the text becomes a regular paragraph and is NOT treated as attribution.

```markdown
> The only way to do great work is to love what you do.
Steve Jobs
```

```markdown
> Stay hungry, stay foolish.
[[Steve Jobs]]
```

Output:
```html
<figure class="quote">
  <blockquote>
    <p>The only way to do great work is to love what you do.</p>
  </blockquote>
  <figcaption>Steve Jobs</figcaption>
</figure>
```

**Edge cases:**
- Multi-paragraph blockquotes followed by attribution: attribution applies to the entire blockquote
- Long continuation lines (> 60 chars) after a blockquote are NOT treated as attribution
- This behavior is always enabled and does not require configuration

### Heading Anchors

Add clickable permalink anchors to headings.

```markdown
## My Heading
```

Output:
```html
<h2 id="my-heading">My Heading <a class="anchor" href="#my-heading">¶</a></h2>
```

**Configuration:**
```toml
[markdown.extensions]
anchor = false  # Enable goldmark heading anchors (default: false; use heading_anchors plugin instead)
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

!!! note

    A blank line after the admonition header is allowed.
    Content must still be indented.

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

Link to a specific section: [[other-post-slug#section-id]]

Section link with custom text: [[other-post-slug#section-id|See this section]]
```

### Resolution

1. Split `link_target` on the first `#` into `slug` and optional `fragment`
2. Find post where `slug == slug` (or alias match, case-insensitive)
3. If found, render as `<a href="{post.href}#{fragment}">{text}</a>` (fragment omitted if empty)
4. If not found, leave as literal `[[link]]` and warn

### Output

```html
<!-- Found -->
<a href="/other-post-slug/">Other Post Title</a>
<a href="/other-post-slug/">Click here</a>

<!-- Found with fragment -->
<a href="/other-post-slug/#section-id">Other Post Title</a>
<a href="/other-post-slug/#section-id">See this section</a>

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

# Extension configuration (markata-go specific)
[markdown.extensions]
typographer = true       # Smart quotes, dashes, ellipses (default: true)
definition_list = true   # PHP Markdown Extra definition lists (default: true)
footnote = true          # PHP Markdown Extra footnotes (default: true)
```

### Disabling Extensions

Individual extensions can be disabled if needed:

```toml
# Disable all optional extensions
[markdown.extensions]
typographer = false
definition_list = false
footnote = false

# Or disable just one
[markdown.extensions]
typographer = false  # Keep straight quotes
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

## Legacy homepage format outputs

An empty-slug homepage MUST publish Markdown and text as root files (`index.md` and `index.txt`). When incremental output contains a legacy directory at one of these paths, the publisher MUST replace it only if it contains exactly one regular `index.html` file matching the generated redirect to `/.md` or `/.txt`. Unexpected files, directories, or symlinks MUST be preserved and the write MUST fail. Migration MUST affect staging output only, preserving previous releases.

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [THEMES.md](./THEMES.md) - Admonition and code block styling
- [CONFIG.md](./CONFIG.md) - Markdown configuration
- [TEMPLATES.md](./TEMPLATES.md) - Template system
- [DATA_MODEL.md](./DATA_MODEL.md) - Post model
