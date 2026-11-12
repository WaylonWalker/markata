---
title: "Write your first post"
description: "A practical Markata authoring guide: Markdown, frontmatter, drafts, tags, links, images, and preview."
date: 2026-10-18
published: true
tags: [documentation, authoring]
---


You don't need templates or plugins to start. **Markdown is the source.** Markata turns it into pages.

## Create a post

From your site directory:

```bash title="Terminal"
markata new "Hello, internet"
```

This normally creates `pages/post/hello-internet.md`. Open it in your editor.

You'll see YAML frontmatter between two `---` lines. The body is ordinary Markdown.

```markdown title="pages/post/hello-internet.md"
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

### Create a link post from a URL

```bash title="Terminal"
markata new --template link 'https://example.com/my-article'
```

Markata fetches the page title, description, and preview image for the link post.
Only a usable Open Graph title (`og:title`) becomes the title automatically.
Empty, URL-only, and hostname-only titles do not qualify.
Other suggestions require your confirmation: Twitter metadata, the HTML title, the main heading, the decoded URL slug, then the hostname.

The post keeps the URL you entered, even after redirects.
Relative preview images resolve against the fetched page URL.
A missing description or image does not prevent creation.
Markata does not invent these optional fields.

In the wizard, select the link template, then enter the URL first.
After metadata lookup and any required title confirmation, **Create** is the default action.
Choose **Edit details** to change the title, description, image, directory, tags, privacy, or authors.
Explicit flags take priority over fetched metadata and wizard values.

If the fetch fails, Markata warns and lets you enter details manually.
For unattended use, supply `--title` and `--no-input`:

```bash title="Terminal"
markata new --template link 'https://example.com/my-article' --title 'My article' --draft --no-input
```

The explicit title always wins.
Without a trusted title or `--title`, unattended creation stops with an error and writes no file.
The error includes a copy-pastable command that adds `--title` and preserves your URL and explicit flags.

Link aliases (`bookmark`, `til`, and `stars`) support the same workflow.
Existing title arguments, such as `markata new "My bookmark" --template link`, still work.
Other templates keep their existing creation behavior.
See the [CLI reference](/docs/reference/cli/#new) for more options.

## Preview your words

```bash title="Terminal"
markata serve
```

Visit [http://localhost:8000](http://localhost:8000). Save your post; the development server rebuilds it. The page URL depends on your slug and theme configuration.

## Draft before you publish

```bash title="Terminal"
markata new "An unfinished idea" --draft
```

For an existing post, set `draft: true`. To publish, set `draft: false` and `published: true`, then rebuild. Draft handling also depends on feed filters: do not assume a draft is publicly hidden if you override the default feed configuration.

## Links and pictures

Use ordinary Markdown for external links and images:

```markdown title="Markdown"
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

```bash title="Terminal"
markata config validate
markata build
```

Upload the generated `output/` directory to your static host. Follow [Deployment](/docs/guides/deployment/), or return to the [Quickstart](/docs/quickstart/).
