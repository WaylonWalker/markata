---
title: "@Mentions"
description: "Resolve @handles to links from blogroll feeds and internal posts"
date: 2025-01-26
published: true
tags:
  - documentation
  - mentions
  - social
---

# @Mentions

The mentions plugin transforms `@handle` syntax in your markdown content into HTML links. It resolves handles from three sources:

1. **Blogroll feeds** - External RSS/Atom feeds from your blogroll configuration
2. **Authors** - Authors defined in your site configuration (automatically registered)
3. **Internal posts** - Posts in your site matching filter expressions (like contact pages)

## Basic Usage

Write mentions in your markdown using `@handle` syntax:

```markdown
I was reading @alice's post about static site generators.
Also check out @bob for more great content.
```

When handles are resolved, they become clickable links:

```html
I was reading <a href="/contact/alice/" class="mention">@alice</a>'s post...
```

## Configuration

Configure mentions in your `markata-go.toml`:

```toml
[markata-go.mentions]
enabled = true
css_class = "mention"  # CSS class for links (default: "mention")
```

That's it for most use cases! By default, mentions resolves handles from:

1. **Blogroll feeds** (if blogroll is enabled)
2. **Authors** from your site configuration (if they have a `url`)
3. **Contact posts** - any post with `template: contact` in frontmatter
4. **Author posts** - any post with `template: author` in frontmatter

The default `from_posts` sources are equivalent to:

```toml
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"

[[markata-go.mentions.from_posts]]
filter = "template == 'author'"
handle_field = "handle"
aliases_field = "aliases"
```

You can override this by specifying your own `from_posts` sources:

```toml
# Custom sources replace the default
[[markata-go.mentions.from_posts]]
filter = "'team' in tags"
handle_field = "github_handle"
aliases_field = "nicknames"
```

## Resolving from Internal Posts

The `from_posts` configuration lets you resolve `@handles` from posts in your site. This is useful for:

- **Contact pages** - Team member or author profiles
- **Contributor pages** - Guest authors or collaborators  
- **Partner pages** - Companies or organizations you work with

### Example: Contact Pages

Create contact pages with handle information in frontmatter:

```yaml
# pages/contact/alice-smith.md
---
title: "Alice Smith"
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

Alice is a software engineer specializing in Go and web development.
```

Configure the mentions plugin to use these posts (or rely on the default which already matches `template: contact`):

```toml
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"
```

Now `@alice`, `@alices`, and `@asmith` all link to `/contact/alice-smith/` and display `@alice` as the link text.

The mention hovercard will show the avatar (from the `avatar`, `image`, or `icon` frontmatter field), the post title as the name, and the description as the bio. You can specify which field to use for the avatar with `avatar_field`:

```toml
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"
avatar_field = "icon"   # Use 'icon' field instead of default lookup order
```

### Multiple Sources

You can define multiple `from_posts` sources for different content types:

```toml
# Contact pages (this is the default, shown here for clarity)
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"
aliases_field = "aliases"

# Team members
[[markata-go.mentions.from_posts]]
filter = "template == 'team-member.html'"
handle_field = "github_handle"

# Projects (no explicit handle - uses slug)
[[markata-go.mentions.from_posts]]
filter = "'project' in tags"
```

### Fallback to Slug

If `handle_field` is not specified or the frontmatter field is empty, the post's slug is used as the handle.

```toml
[[markata-go.mentions.from_posts]]
filter = "'contributor' in tags"
# No handle_field - uses slug
```

A post at `contributors/jane-doe.md` becomes `@jane-doe`.

## Resolving from Blogroll

If you have a [blogroll](/guides/blogroll/) configured, handles are automatically extracted from your RSS/Atom feeds:

```toml
[markata-go.blogroll]
enabled = true

[[markata-go.blogroll.feeds]]
url = "https://daverupert.com/feed.xml"
title = "Dave Rupert"
handle = "daverupert"
site_url = "https://daverupert.com"
aliases = ["dave", "davatron"]
```

Now `@daverupert`, `@dave`, and `@davatron` all link to `https://daverupert.com`.

## Authors as Contacts

Authors defined in your site configuration are automatically registered as mentionable contacts. If an author has a `url` field, they become a valid `@handle` target.

```toml
[markata-go.authors.authors.waylon]
name = "Waylon Walker"
url = "https://waylonwalker.com"
avatar = "/images/waylon.jpg"
bio = "Python and Go developer"
```

With this configuration, `@waylon` becomes a valid mention that links to `https://waylonwalker.com`. The author's avatar, name, and bio are used for hovercards.

Authors are registered after blogroll feeds but before `from_posts` sources. If a blogroll feed already registered the same handle, the blogroll entry wins (first registration wins).

**No extra configuration needed** -- as long as your authors have a `url` field, they are automatically available as mentions.

## Authors in from_posts

By default, the mentions plugin includes posts with `template: author` as a post source, in addition to `template: contact`:

```toml
# These are the defaults (you don't need to add them):
[[markata-go.mentions.from_posts]]
filter = "template == 'contact'"
handle_field = "handle"

[[markata-go.mentions.from_posts]]
filter = "template == 'author'"
handle_field = "handle"
aliases_field = "aliases"
```

This means if you have author profile pages with `template: author`, they are automatically included as mention sources.

## Trailing Punctuation

Mentions handle trailing punctuation gracefully. When you write `@alice.` or `@bob,` at the end of a sentence, the plugin strips the punctuation and resolves the handle correctly:

```markdown
I was talking to @alice.
Check out @bob, they have great content.
Thanks @charlie!
```

Produces:

```html
I was talking to <a class="mention" href="...">@alice</a>.
Check out <a class="mention" href="...">@bob</a>, they have great content.
Thanks <a class="mention" href="...">@charlie</a>!
```

The following trailing punctuation characters are stripped: `. , ; : ! ?`

**Domain-style handles are preserved:** If the exact handle (including dots) matches a known contact, it is used as-is. For example, `@simonwillison.net` resolves to the domain handle without stripping the `.net`. Only when the exact handle fails to match does the plugin try stripping trailing punctuation.

## Chat Admonitions

The mentions plugin enriches chat admonition titles with contact information. When you use `!!! chat` or `!!! chat-reply` with an `@handle`, the title is replaced with the contact's avatar and a linked mention.

Both quoted and unquoted forms work:

```markdown
!!! chat @alice

    Hey, have you seen the new release?

!!! chat-reply @bob

    Yes! It looks great.

!!! chat "@alice"

    Quoted form works too.
```

`chat-reply` can also omit a handle to render as the post author:

```markdown
!!! chat-reply

    Thanks for the update!
```

This renders with avatar images and linked names in the admonition title, creating a conversation-style display.

Collapsible chat admonitions work too:

```markdown
??? chat @alice

    This content is collapsible.
```

If the handle is not found in the contact map, the title is left unchanged.

### Admonition Header Protection

The mentions plugin automatically protects `@handles` on admonition header lines from being transformed into links. This applies to all admonition types (`!!!`, `???`, `???+`), not just chat. Mentions in the body text and outside admonitions are still processed normally.

```markdown
!!! note @alice
    This note title is preserved as-is.

But @alice in regular text is still linked.
```

### Styling Chat Admonitions

Chat admonitions include CSS classes for styling:

```css
.chat-contact {
  display: inline-flex;
  align-items: center;
  gap: 0.5em;
}

.chat-contact-avatar {
  width: 1.5em;
  height: 1.5em;
  border-radius: 50%;
  vertical-align: middle;
}

.chat-contact-name {
  font-weight: 600;
}
```

These styles are included in the default theme. Customize them to match your site's design.

## Handle Resolution Priority

When building the handle map, sources are processed in this order:

1. **Blogroll feeds** (if enabled) - First registered wins
2. **Authors** (from site config) - Authors with a `url` field
3. **From posts** sources - Processed in order they appear in config

If the same handle is registered multiple times, the first registration wins and subsequent duplicates are logged as warnings.

## Styling Mentions

Mention links have the configured CSS class (default: `mention`). Add styles to your CSS:

```css
a.mention {
  color: var(--color-primary);
  font-weight: 500;
}

a.mention::before {
  content: "";  /* Remove @ if you want */
}

a.mention:hover {
  text-decoration: underline;
}
```

## Code Block Protection

Mentions inside fenced code blocks are preserved and not transformed:

````markdown
Check out @alice's post.

```
// @alice is not transformed here
const handle = "@alice";
```
````

## Email Protection

Email addresses are not transformed as mentions. The plugin detects the `@` preceded by word characters and ignores them:

```markdown
Contact me at test@example.com  <!-- Not transformed -->
Follow @example on social       <!-- Transformed -->
```

## Configuration Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable mentions processing |
| `css_class` | string | `"mention"` | CSS class for mention links |
| `cache_dir` | string | `"cache/mentions"` | Directory for cached external metadata |
| `cache_duration` | string | `"168h"` | How long external metadata stays fresh |
| `timeout` | int | `30` | HTTP timeout in seconds |
| `concurrent_requests` | int | `3` | Max concurrent external metadata fetches |
| `from_posts` | array | see below | List of internal post sources. Default: two sources with `filter = "template == 'contact'"` and `filter = "template == 'author'"`, both with `handle_field = "handle"` |

### from_posts Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `filter` | string | required | Filter expression to select posts |
| `handle_field` | string | (slug) | Frontmatter field for handle |
| `aliases_field` | string | - | Frontmatter field for aliases |
| `avatar_field` | string | - | Frontmatter field for avatar URL. If not set, checks `avatar`, `image`, `icon` in order |

## Examples

### Team Directory

```toml
[[markata-go.mentions.from_posts]]
filter = "'team' in tags and published == true"
handle_field = "slack_handle"
aliases_field = "nicknames"
```

### Author Profiles

```toml
[[markata-go.mentions.from_posts]]
filter = "template == 'author.html'"
handle_field = "username"
```

### Combined Sources

```toml
[markata-go.blogroll]
enabled = true

[[markata-go.blogroll.feeds]]
url = "https://external-blog.com/rss"
handle = "external"

[[markata-go.mentions.from_posts]]
filter = "'internal' in tags"
handle_field = "handle"
```

Now you can use `@external` for external links and `@internal-user` for internal profile links in the same post.


## Quote attributions with a tilde

Use `~handle` to credit the source of a short quote without turning it
into a prominent callout:

```markdown
if an agent cant verify its own work

==nothing else matters==

~poteto
```

When `poteto` is in your configured blogroll, Markata links it to the
canonical `site_url` and displays it with a quiet, theme-aware
`mention--attribution` style. Ordinary `@poteto` references keep the
normal mention appearance. You can override the styling using
`.mention--attribution` in your site CSS.

If the handle is not yet in the blogroll, `~poteto` stays visible and
styled, but is not made into a link. Markata does **not** assume that an
identically named GitHub/X/Bluesky account belongs to the person you
quoted.

### Lint missing mentions locally

```sh
markata-go lint pages/ping/ping-84.md
```

With blogroll enabled, both `@handle` and `~handle` produce
`unknown-mention` warnings when they cannot be resolved to an active
configured feed handle or alias. The warning suggests adding the
person's handle or alias to the blogroll; warnings do not fail the
build. This check is deterministic and **does not search the network**.

For new people, first check your blogroll. If absent, manually
investigate X/Twitter, then GitHub, then Bluesky for a verified match,
and add the verified canonical link plus an actual feed to the
blogroll. Never automatically insert a social profile based only on a
username collision.
