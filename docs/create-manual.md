---
title: Manually create a new site
description: Lets manually create a new website with markata.

---

Markata has some templates that let you get up and running quickly, but you
_can_ make a site with only markdown if you wanted.

### Installation

`markata` is hosted on pypi and can be installed using pip.

```bash
python -m pip install markata

# or if pipx is your thing

pipx install markata
```

### Create Some Content

Make some `.md` files in your current working directory. By default, `markata`
will recursively look in all subdirectories for markdown files `**/*.md`.

```bash
mkdir pages
echo '# My First Post' > first-post.md
echo '# Hello World' > hello-world.md
```

> This example shows how you can build a site from only a single markdown
> file.

### Build your site

Install markata into your virtual environment and run `markata build`. It will
create your site in `./markout`, leave its cache in `./.markata.cache`, and
copy all assets from `./static` into `./markout` by default.

```bash
python -m pip install markata
markata build

# or if pipx is your thing
pipx run markata build
```
