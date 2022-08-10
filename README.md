<h1 align=center>
  <br>
  <a href="https://github.com/WaylonWalker/markata"><img src="https://user-images.githubusercontent.com/22648375/167527193-4e66d537-e78f-44dd-be65-2c1c109adf34.png" alt="Markata" width="400"></a>
</h1>

<h1 align=center>
  <br>
  <a href="https://github.com/WaylonWalker/markata"><img src="https://user-images.githubusercontent.com/22648375/183960585-a8de2452-b326-4e49-912c-891b686270c4.png" alt="Markata" width="400"></a>
</h1>

<p align=center>
  <em>
    Markdown to site, plugins all the way down
  </em>
</p>

A static site generator that will give you a great site with many standard web
features like rss, sitemaps, and seo tags, out of the box.  Running `markata
build` will get you a that only requires you to write markdown.  If you have
additional features that you want, don't worry, since markata is built
completely on plugins you can develop and install your own plugins to add the
features you want.


## QuickStart

Markata is fully configurable through a `markata.toml` file, but the defaults
allow to build your site right out of the box with nothing more than markdown.

### Create Some Content

Make some `.md` files in your current working directory.  By default `markata`
will recursively look in all subdirectories for markdown files `**/*.md`.

```
mkdir pages
echo '# My First Post' > first-post.md
echo '# Hello World' > hello-world.md
```

### Build your site

Install markata into your virtual environment and run `markata build`.  It will
create your site in `./markout`, leave its cache in `./.markata.cache`, and
copy all assets from  `./static` into `./markout` by default.

``` bash
python -m pip install markata
markata build

# or if pipx is your thing
pix run markata build
```

## Motivation

Markata is able to build your site purely from markdown, allowing you to get
started creating your own content quickly.  Out of the box it will cover your
seo tags, rss feeds, sitemap, and og images.  Since it is built completely from
plugins you are able to remove, modify, or add to any of its behavior.

* configurable
* plugins
* seo
* rss
* sitemap
* og-image

## Using Markata

The docs are still a work in progress, but the
[base_cli](https://markata.dev/markata/plugins/base_cli/) walks through how to
effectively use the `markata build` command. At this point Markata is far from
stable and **will change quite a bit**, should you choose to use it on  real
site make sure you pin to the version that you want to build from.  I will let
you know as it becomes more stable and ready to use without diligence of
pinning to the version you want.

**Honestly**  A big motivation for me was wanting to learn and understand how
to create a project that is completely plugin driven.  This is highly a
learning project for me, and it has grown into something I use each and every
day.

## Examples Gallary

### [Markata.dev](https://markata.dev)

Yes, markata builds its own docs

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/151674260-d5b1a073-ba68-4274-aac1-3b891a31e3ed.png' width=400px>
</p>

> Home page, created with index.md

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/151674334-a5fb0205-2631-4057-8ecb-f8ba1e7ebaf9.png' width=400px>
</p>

> [base_cli plugin](https://markata.dev/markata/plugins/base_cli/) documentation generated with the [docs plugin](https://markata.dev/markata/plugins/docs/)

### [WaylonWalker.com](https://waylonwalker.com)

Waylonwalker.com is created completely through markata

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/151674183-8c36cab2-bccd-4733-b78b-99384e257b00.png' width=400px>
</p>

> Post Page

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/151674204-264d549a-fc33-4373-a675-4f5a31daaf9f.png' width=400px>
</p>

> archive page created through custom plugin

### [techstructive Blog](https://mr-destructive.github.io/techstructive-blog/series/)

The very first adopter of markata, [meet gor](https://twitter.com/MeetGor21) writes about golang, Django, and Bash Scripting.

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/168610225-9a125354-56c9-4d30-8d5d-6d45a8bb6ac1.png' width=400px>
</p>

He has even created a custom plugin for dynamically adding series's to his site!

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/168610778-be626d3b-100d-44e0-af46-10ab0255da07.png' width=400px>
</p>

### [images.WaylonWalker.com](https://images.waylonwalker.com)

Waylonwalker.com currently has the built in cover image pluugin disabled for quick builds as it it a constantly evolving site with a lot of posts.  The cover images are generated in a second repo by loading article data in from [markata.json](https://waylonwalker.com/markata.json) and running the covers plugin.


<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/116886610-e6f03480-abee-11eb-92c8-f883314fd09a.png' width=400px>
</p>

### [pype.dev](https://pype.dev)

Pype.dev is a menta data lake of all things python, linux, and homelab.

<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/168611555-b7c918fa-836e-4334-9cbd-cbb07b8ac350.png' width=400px>
</p>

