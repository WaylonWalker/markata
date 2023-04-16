---
title: Getting Started with Markata
description: Guide to get going with Markata
---

Markata is a fully plugins all the way down static site generator for
folks who just want their site to get started quickly and build great
content, with the ability to tinker with everything under the hood if
they want to.

## QuickStart

Markata is fully configurable through a `markata.toml` file, but the defaults
allow to build your site right out of the box with nothing more than markdown.

### Create Some Content

```
mkdir pages
cd pages
echo '# My First Post' > first-post.md
echo '# Hello World' > hello-world.md
```

### Build your site

```bash
pip install markata
markata build

# or if pipx is your thing
pipx run markata build
```

### Frontmatter

You will likely want to set things like `title`, `date`, `description`,
`published`, or `template` per post, this can all be done inside yaml frontmatter.

```markdown
---
templateKey: blog-post
tags: ["python"]
title: My Awesome Post
date: 2022-01-21T16:40:34
published: False
---

This is my awesome post.
```

> Frontmatter is not required, but definitely gives you more control over your site.

## Next steps

_blog starter_

The [blog-starter](https://blog-starter.markata.dev/) has a really great write
up on how to use markata. You can see it in your brower at the
[link](https://blog-starter.markata.dev/) or run it yourself `pipx run markata
new blog`.

## Examples Gallary

Markata has a project gallery to show off sites built with markata. Please
[submit](https://github.com/WaylonWalker/markata/issues/78) yours, and check
out the [project-gallery](http://markata.dev/project-gallery/) for inspiration.

## Deploying to a sub route

_gh pages_

To deploy a subroute, add a markata.path_prifix to your config (markata.toml).

```toml
[markata]
path_prefix='my-sub-route'
```

## Source Code

If you have comments, questions, issues, or like it enough to give a ⭐
check out
[waylonwalker/markata](https://github.com/WaylonWalker/markata)

## Markata Docs

Not much is documented yet, lots of work to do on the docs. Checkout
[LifeCycle](https://markata.dev/markata/lifecycle/) to see what a more
finished one looks like.

UPDATE - the
[`base_cli`](https://markata.dev/markata/plugins/base_cli/) is also up to
date and includes a lot of examples of how to use the markata cli.

> **Yes** this library generates it's own docs

- [All Modules](https://markata.dev/autodoc/)
- [Core Modules](https://markata.dev/core_modules/)
- [Plugins](https://markata.dev/plugins/)
- [color theme](https://markata.dev/color-theme/)
