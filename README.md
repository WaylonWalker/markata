<h1 align=center>
  <br>
  <a href="https://github.com/WaylonWalker/markata"><img src="https://user-images.githubusercontent.com/22648375/167527193-4e66d537-e78f-44dd-be65-2c1c109adf34.png" alt="Markata" width="400"></a>
</h1>

<p align=center>
  <em>
    Markdown to site, plugins all the way down
  </em>
</p>

A static site generator that will give you a great site with many standard web
features like rss, sitemaps, and seo tags, out of the box. Running `markata
build` will get you a that only requires you to write Markdown. If you have
additional features that you want, don't worry, since markata is built
completely on plugins you can develop and install your own plugins to add the
features you want.

> This has been a pet project for me to learn library development, plugin
> driven design, diskcache, and more.  It is the core of what builds my own site [waylonwalker.com](https://waylonwalker.com).

## Disclaimer

Make sure that you pin down what version of markata you want to use.  If you
are starting a new project that's probably the latest version from
[pypi](https://pypi.org/project/markata).  Things are likely to change in major
releases, I do my best to document them, and not to break patches.

## QuickStart

Markata is fully configurable through a `markata.toml` file, but the defaults
allow to build your site right out of the box with nothing more than Markdown.

### Installation

`markata` is hosted on pypi and can be installed using pip.

```bash
python -m pip install markata

# or if pipx is your thing

pipx install markata
```

### Building a Blog with Markata _using a template_

The markata cli includes a `new` command that will present you with questions
to fill in the jinja variables in this repo.

```bash
markata new blog [directory]

markata build
# start the site and watch for changes
markata serve
```

Now if you open localhost:8000, you will be presented with an example site that
will walk you through some features of markata. You can play with it at
your own pace, or drop all the pages and start writing your own content.

![new blog gif](./static/new-blog.gif)

## Motivation

Markata is able to build your site purely from Markdown, allowing you to get
started creating your own content quickly. Out of the box it will cover your
seo tags, rss feeds, sitemap, and og images. Since it is built completely from
plugins you can remove, modify, or add to any of its behavior.

- configurable
- plugins
- seo
- rss
- sitemap
- og-image

## Using Markata

The docs are still a work in progress, but the
[base_cli](https://markata.dev/markata/plugins/base_cli/) walks through how to
effectively use the `markata build` command. At this point Markata is far from
stable and **will change quite a bit**, should you choose to use it on real
site make sure you pin to the version that you want to build from. I will let
you know as it becomes more stable and ready to use without diligence of
pinning to the version you want.

**Honestly** A big motivation for me was wanting to learn and understand how to
create a project that is completely plugin driven. This is highly a learning
project for me, and it has grown into something I use every day.

## Examples Gallary

Markata has a project gallery to show off sites built with markata. Please
[submit](https://github.com/WaylonWalker/markata/issues/78) yours, and check
out the [project-gallery](http://markata.dev/project-gallery/) for inspiration.
