<h1 align=center>
  <br>
  <a href="https://github.com/WaylonWalker/markata"><img src="https://user-images.githubusercontent.com/22648375/120665701-a126c600-c451-11eb-9dc5-1e0ac4eeeb2f.png" alt="Markata" width="200"></a>
</h1>

<p align=center>
  <em>
    Markdown and Data, plugins all the way down
  </em>
</p>


## Example

Here is markata creating my personal site of several hundred pages in ~16s.

<p align="center">
<a href='https://user-images.githubusercontent.com/22648375/116888181-c32dee00-abf0-11eb-9eba-8c6997d3c888.mp4' >
   <img src='https://user-images.githubusercontent.com/22648375/116888173-c1642a80-abf0-11eb-8647-25a47aacc1fc.gif' align=center>
</a>
</p>


Here is markata doing a near instant with a fresh cache.

<p align="center">
<a href='https://user-images.githubusercontent.com/22648375/116885339-654bd700-abed-11eb-8e65-3202bcce1773.mp4' >
   <img src='https://user-images.githubusercontent.com/22648375/116885394-7563b680-abed-11eb-8649-b8d3fbc728b9.gif' align=center>
</a>
</p>



I am still learning what this will be.  Code is completely in develop branch and likely to change significantly.

## Motivation

I want a simple static site generator built off of plugins that are super easy to create, and intuitive to use.  I want to create content in markdown and not have to put much thought into seo meta tags or generating og images.  The final site should be lightweight, and not weighed down heavily with MB's of unused css/js.

⚡ Fast enough site generation

🔌 Pluggable

🧠 Intuitive

🐍 Familiar language

🖼 OG:image out of the box

🎯 SEO out of the box


**Honestly**  A big motivation for me was wanting to learn and understand how to create a project that is completely plugin driven.  This is highly a learning project for me, and it has grown into something I use each and every day.

## Development

Currently everything is on the [develop](https://github.com/WaylonWalker/markata/tree/develop) branch.  As its still a heavy work in progress and likely to change significantly.


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

### [images.WaylonWalker.com](https://images.waylonwalker.com)

Waylonwalker.com currently has the built in cover image pluugin disabled for quick builds as it it a constantly evolving site with a lot of posts.  The cover images are generated in a second repo by loading article data in from [markata.json](https://waylonwalker.com/markata.json) and running the covers plugin.


<p align=center>
  <img src='https://user-images.githubusercontent.com/22648375/116886610-e6f03480-abee-11eb-92c8-f883314fd09a.png' width=400px>
</p>
