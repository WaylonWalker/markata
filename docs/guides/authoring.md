---
title: "Write your first post"
description: "A practical Markata authoring guide: Markdown, frontmatter, drafts, tags, links, images, and preview."
date: 2026-10-18
published: true
tags: [documentation, authoring]
---

# Write your first post

You don't need templates or plugins to start. **Markdown is the source.** Markata turns it into pages.

## Create a post

From your site directory:

```bash
markata new "Hello, internet"
```

This normally creates `pages/post/hello-internet.md`. Open it in your editor.

You'll see YAML frontmatter between two `---` lines. The body is ordinary Markdown.

```markdown
---
title: "Hello, internet"
date: 2026-10-18
published: true
draft: false
tags: ["notes", "web"]
description: "My first post with Markata."
---

Hello, internet. I'm building a little corner of the web.

## Why this exists

I want somewhere to write and share what I'm learning.

Here's **bold text**, a [link](https://example.com), and a list:

- Make something small
- Publish it
- Improve it later
```

Keep the frontmatter simple: `title`, `date`, `published`, `draft`, and optional `tags` are enough to begin. Markata can infer a slug, and additional frontmatter fields are available to templates.

## Preview your words

```bash
markata serve
```

Visit [http://localhost:8000](http://localhost:8000). Save your post; the development server rebuilds it. The page URL depends on your slug and theme configuration.

## Draft before you publish

```bash
markata new "An unfinished idea" --draft
```

For an existing post, set `draft: true`. To publish, set `draft: false` and `published: true`, then rebuild. Draft handling also depends on feed filters: do not assume a draft is publicly hidden if you override the default feed configuration.

## Links and pictures

Use ordinary Markdown for external links and images:

```markdown
[My favorite reference](https://example.com)
![A description of the photograph](/images/my-photo.jpg)
```

Place site assets in `static/` according to your configured assets directory. Check the emitted URL in the browser.

Markata also supports **wikilinks** for your own posts:

```markdown
I wrote about this in [[hello-internet]].
Read [[hello-internet|my first post]].
```

Use a matching target slug. Prefer a normal link when sending someone to an external site.

## Organize later

Add tags when they're useful. Use [Feeds](/docs/guides/feeds/) for collections, archives, and RSS; [Themes](/docs/guides/themes/) for visual changes; and [Templates](/docs/guides/templates/) for layout overrides.

Don't start by writing a plugin. First ask whether frontmatter, config, an existing feed, or a template already does what you need.

## Writing with an agent

After [installing the skill](/docs/guides/agent-skills/), a safe first request is:

> Make a draft in the existing posts directory from these notes. Keep my voice, use simple frontmatter, and don't publish anything yet. Validate config and run `markata build --fast` when done.

Review the result yourself before changing `draft` or deploying.

## Ship it

```bash
markata config validate
markata build
```

Upload the generated `output/` directory to your static host. Follow [Deployment](/docs/guides/deployment/), or return to the [Quickstart](/docs/quickstart/).
