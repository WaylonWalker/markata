# Mentions Plugin Specification

The mentions plugin transforms `@handle` syntax in markdown content into HTML links, resolving handles from blogroll feeds, internal posts, and configured authors.

## Overview

The mentions plugin processes `@handle` patterns in post content and replaces them with HTML anchor links. Handles are resolved from three sources:

1. **Blogroll feeds** - External RSS/Atom feed entries with handles
2. **Internal posts** - Posts matching filter expressions (e.g., contact pages, team pages, author pages)
3. **Authors** - Authors defined in the site configuration are automatically registered as mentionable contacts

## Syntax

```markdown
@handle           # Basic mention
@alice-smith      # Hyphenated handle
@alice123         # Handle with numbers
```

Handle pattern: `@[a-zA-Z][a-zA-Z0-9_.-]*`

### Trailing Punctuation

When a handle includes trailing punctuation (e.g., `@alice.` at end of a sentence), the plugin first tries an exact match. If the exact match fails, trailing punctuation characters (`.,;:!?`) are stripped and the lookup is retried. This preserves domain-style handles like `@simonwillison.net` (which match exactly) while correctly resolving `@alice.` → `@alice` + `.`.

```markdown
Talk to @alice.      <!-- Resolves to @alice, period is preserved as text -->
Visit @simonwillison.net  <!-- Resolves to @simonwillison.net (exact match) -->
Hey @bob, welcome!   <!-- Resolves to @bob, comma is preserved as text -->
```

## Configuration

```toml
[markata-go.mentions]
enabled = true                    # Enable/disable the plugin (default: true)
css_class = "mention"             # CSS class for links (default: "mention")
cache_dir = "cache/mentions"      # Metadata cache directory
cache_duration = "168h"           # Cache TTL for fetched metadata
timeout = 30                      # HTTP timeout in seconds
concurrent_requests = 3           # Max concurrent metadata fetches

# Source handles from internal posts (optional - a default source is provided, see below)
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"  # Filter expression to select posts
handle_field = "handle"           # Frontmatter field for handle (optional, uses slug if not set)
aliases_field = "aliases"         # Frontmatter field for aliases (optional)
avatar_field = "avatar"           # Frontmatter field for avatar URL (optional, auto-detects if not set)
```

### Configuration Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable mentions processing |
| `css_class` | string | `"mention"` | CSS class applied to mention links |
| `cache_dir` | string | `"cache/mentions"` | Directory for mention metadata cache |
| `cache_duration` | string | `"168h"` | How long fetched metadata stays fresh |
| `timeout` | int | `30` | HTTP timeout in seconds |
| `concurrent_requests` | int | `3` | Max concurrent external metadata fetches |
| `from_posts` | array | see below | List of internal post sources |

Mention metadata freshness is determined by the cached `last_fetched` timestamp,
not the cache file mtime, so touching cache files does not accidentally invalidate them.

### from_posts Default

When no `from_posts` sources are configured (or when a `[markata-go.mentions]` section exists without any `[[markata-go.mentions.from_posts]]` entries), a default source is used:

```toml
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"
```

This means any post with `template: contact` in its frontmatter is automatically available as a mentionable contact. To disable this default, set `from_posts` to an explicit source with a different filter.

### from_posts Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `filter` | string | Yes | Filter expression to select posts |
| `handle_field` | string | No | Frontmatter field containing the handle (defaults to post slug) |
| `aliases_field` | string | No | Frontmatter field containing handle aliases |
| `avatar_field` | string | No | Frontmatter field containing avatar URL. If not set, checks `avatar`, `image`, `icon` in order |

## Handle Resolution

### Resolution Order

Handles are registered in this order (first registration wins):

1. **Blogroll feeds** - If blogroll is enabled, handles from feeds are registered first
2. **Authors** - Authors from site configuration are registered (using author ID as handle)
3. **from_posts sources** - Processed in the order they appear in configuration

If a handle or alias is already registered, subsequent registrations are ignored with a warning.

### From Blogroll

When blogroll is enabled, handles are extracted from feed configurations:

```toml
[[markata-go.blogroll.feeds]]
url = "https://example.com/feed.xml"
title = "Example Blog"
handle = "example"                  # Primary handle
site_url = "https://example.com"    # Link target
aliases = ["ex", "exblog"]          # Additional handles
```

The `site_url` becomes the link target for `@example`, `@ex`, and `@exblog`.

When cached blogroll feed data exists, explicit config overrides such as `title`, `description`, `site_url`, `image_url`, `handle`, aliases, category, and tags MUST be reapplied before mentions resolve handles or avatar metadata. Config changes MUST take effect on the next build without requiring cache deletion.

### From Authors

Authors defined in the site configuration (`[markata-go.authors.authors.*]`) are automatically registered as mentionable contacts. Each author is registered using their config key (ID) as the handle. The author's URL (if set) is used as the link target; if no URL is set, the author is not registered.

```toml
[markata-go.authors.authors.waylon]
name = "Waylon Walker"
url = "https://waylonwalker.com"
avatar = "/images/waylon.jpg"
bio = "Python and Go developer"
```

This automatically registers `@waylon` as a mention handle pointing to `https://waylonwalker.com`, with name, avatar, and bio metadata for hovercards.

**Author posts**: Posts with `template: author` are also automatically included as a default `from_posts` source, similar to `template: contact`. This allows author profiles to be defined as regular posts.

### From Internal Posts

Internal posts are filtered using the filter expression, then handles are extracted:

1. **Handle extraction**: Uses `handle_field` frontmatter value, or falls back to post slug
2. **Alias extraction**: Uses `aliases_field` frontmatter value if specified
3. **Link target**: The post's permalink URL

Example post frontmatter:

```yaml
---
title: "Alice Smith"
slug: alice-smith
template: contact
tags:
  - contact
handle: alice
aliases:
  - alices
  - asmith
avatar: /images/alice.jpg
url: https://alicesmith.dev
description: "Software engineer specializing in Go and web development."
---
```

With config:

```toml
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"
```

This registers:
- `@alice` → `/contact/alice-smith/`
- `@alices` → `/contact/alice-smith/`
- `@asmith` → `/contact/alice-smith/`

### Slug Fallback

If `handle_field` is empty or the specified field is not present in frontmatter, the post's slug is used:

```toml
[[markata-go.mentions.from_posts]]
filter = "'project' in tags"
# No handle_field - uses slug
```

A post with `slug: my-project` registers `@my-project`.

## Transform Behavior

### Processing

1. Extract code blocks and store them (to protect from transformation)
2. Find all `@handle` patterns in content
3. For each handle:
   - Look up in handle map
   - If found: replace with anchor link
   - If not found: leave as-is
4. Restore code blocks

### Warm-Build Caching

Implementations MAY restore previously transformed post content from the build cache
when the raw post content and the effective mention-resolution state are unchanged.
The cache key MUST include any data that affects emitted HTML, including:

- the raw post content
- the configured mention CSS class
- resolved handle targets and aliases
- mention metadata used for hovercard attributes
- post author display data used for chat-reply title enrichment

Metadata cache timestamps alone MUST NOT invalidate transformed post content unless
the rendered mention attributes themselves change.

### Email Protection

Email addresses are not transformed. The regex matches `@` only when:
- NOT preceded by word characters (letters, digits, underscore)

```markdown
test@example.com     <!-- Not transformed (email) -->
Follow @example      <!-- Transformed -->
Say hello@example    <!-- Not transformed (no space before @) -->
```

### Code Block Protection

Mentions inside fenced code blocks are preserved:

````markdown
Check out @alice's post.

```
// @alice is not transformed here
const handle = "@alice";
```
````

### Generated HTML

For external mentions (blogroll):

```html
<a href="https://example.com" class="mention" data-name="Example Blog" data-bio="A great blog" data-avatar="https://example.com/avatar.jpg" data-handle="@example">@example</a>
```

For internal mentions (from_posts), the link points to the contact page and includes metadata for hovercards:

```html
<a href="/contact/alice-smith/" class="mention" data-name="Alice Smith" data-bio="Software engineer specializing in Go and web development." data-avatar="/images/alice.jpg" data-handle="@alice">@alice</a>
```

Attributes:
- `href` - Link target URL (contact page for internal, site URL for external)
- `class` - Configured CSS class (default: "mention")
- `data-name` - Display name (from metadata or post title)
- `data-bio` - Bio/description (from metadata or post description)
- `data-avatar` - Avatar URL (from metadata or post frontmatter)
- `data-handle` - The original @handle

The display text is always `@handle` for both internal and external mentions.
The `data-name`, `data-bio`, and `data-avatar` attributes provide hovercard data.

## Lifecycle Stage

The mentions plugin runs in the **Transform** stage with default priority (0).

## Data Model

### MentionsConfig

```go
type MentionsConfig struct {
    Enabled   *bool               // Enable/disable (default: true)
    CSSClass  string              // CSS class (default: "mention")
    FromPosts []MentionPostSource // Internal post sources
}
```

### MentionPostSource

```go
type MentionPostSource struct {
    Filter       string // Filter expression (required)
    HandleField  string // Frontmatter field for handle (optional)
    AliasesField string // Frontmatter field for aliases (optional)
    AvatarField  string // Frontmatter field for avatar URL (optional)
}
```

### Internal Handle Map

```go
// handleMap maps handles to URLs
type handleMap map[string]string

// Example:
// "alice" → "/contact/alice/"
// "alices" → "/contact/alice/"
// "daverupert" → "https://daverupert.com"
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Unknown handle | Left as plain text (`@unknown`) |
| Duplicate handle registration | Warning logged, first registration wins |
| Invalid filter expression | Error logged, source skipped |
| Missing handle_field in frontmatter | Falls back to slug |
| Empty aliases_field | No aliases registered for that post |

## CSS Classes

| Class | Purpose |
|-------|---------|
| `.mention` | Default class for all mention links |

Custom class can be configured via `css_class` option.

## Examples

### Team Directory

```toml
[[markata-go.mentions.from_posts]]
filter = "'team' in tags and published == true"
handle_field = "slack_handle"
aliases_field = "nicknames"
```

### Multiple Sources

```toml
# External blogroll
[markata-go.blogroll]
enabled = true

[[markata-go.blogroll.feeds]]
url = "https://external.com/feed.xml"
handle = "external"
site_url = "https://external.com"

# Internal contacts
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"

# Internal projects
[[markata-go.mentions.from_posts]]
filter = "'project' in tags"
# Uses slug as handle
```

### Minimal Configuration

```toml
# Just enable with defaults - automatically picks up template: contact posts
# and blogroll handles
[markata-go.mentions]
enabled = true
```

This uses the default `from_posts` source (`template == 'contact'` filter with `handle` and `aliases` fields), so any post with `template: contact` frontmatter is automatically mentionable.

## Chat Admonition Integration

When a `chat` or `chat-reply` admonition title starts with `@handle`, the mentions plugin enriches the title with the contact's avatar and a linked name. This creates a conversation-like display where each chat bubble shows who is speaking.

### Syntax

Both quoted and unquoted forms are supported:

```markdown
!!! chat @alice
    Hey, have you seen the new release?

!!! chat-reply @bob
    Yes! The new features look great.

!!! chat "@alice"
    Quoted form also works.
```

`chat-reply` can also omit a handle to render as the post author:

```markdown
!!! chat-reply
    Thanks for the update!
```

When an unquoted `@handle` is used and resolves, the plugin wraps the enriched title in quotes so the admonition parser treats it correctly. When the handle does not resolve, the line is left unchanged.

### Admonition Header Protection

The general `@mention` pass skips all admonition header lines (`!!!`, `???`, `???+`). This prevents `@handles` in admonition titles from being transformed into `<a>` tags, which would break the admonition parser. This protection applies to **all** admonition types, not just `chat`/`chat-reply`.

### Rendered HTML

When the title starts with `@handle` and the handle resolves in the mention map, the admonition title is replaced with:

```html
<span class="chat-contact">
  <img class="chat-contact-avatar" src="/images/alice.jpg" alt="Alice Smith" />
  <a href="/contact/alice/" class="mention" data-name="Alice Smith" data-bio="..." data-avatar="/images/alice.jpg" data-handle="@alice">@alice</a>
</span>
```

If the handle does not resolve, the title remains as plain text (e.g., `@alice`).

### CSS Classes

| Class | Purpose |
|-------|---------|
| `.chat-contact` | Wrapper for avatar + mention link in chat title |
| `.chat-contact-avatar` | Avatar image in chat title (32x32, rounded) |

## Related Features

- **Blogroll** - Source of external handles
- **Wikilinks** (`[[slug]]`) - Internal links by slug
- **Filter Expressions** - Query language for selecting posts


## Quote Attribution Syntax

A leading tilde marks the source of quoted words:

```markdown
If an agent cannot verify its own work, nothing else matters.

~poteto
```

Both `@handle` and `~handle` use the **same blogroll-first** lookup.
`@handle` remains a standard mention link; `~handle` is a quote
attribution and emits `class="mention mention--attribution"` with
the visible `~` prefix. An unresolved tilde attribution is styled as
a noninteractive `span.mention--attribution.mention--unresolved` and
MUST NOT link to a guessed social account. Unresolved `@handle`
continues to render as plain text.

The renderer uses no social-account lookups. Caching must invalidate on
mention syntax changes (cache signature version and/or source changes).

### Offline author warnings

When a site has an enabled blogroll, `markata-go lint` checks both
`@handle` and `~handle` against configured active feed handles,
domain aliases, and explicit aliases. Unknown references are reported
as warnings with their source location and a suggestion to add a handle
or alias to the blogroll. Lint never contacts X, GitHub, Bluesky or other
social sites to guess identity. Default lint preserves its historical
behavior for sites where the blogroll is disabled.

An explicit, interactive discovery tool may later offer candidates in
this priority: configured blogroll, X/Twitter, GitHub, Bluesky. Candidate
links require the author's approval and identity evidence before
becoming blogroll entries; a handle match alone is insufficient.
