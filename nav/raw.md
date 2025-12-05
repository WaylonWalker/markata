---
title: Creating your Navbar
description: Guide to creating a navbar in markata using the default template.
jinja: false
---

Creating navbar links with the default markata templates is done by adding
links in your `markata.toml` configuration within a `markata.nav` block.

## Example

The following example will create two links, one to the root of the site, with
the text `markata` and one to the github repo for markata with the text of
`GitHub`.

```toml
[markata.nav]
'markata'='/'
'GitHub'='https://github.com/WaylonWalker/markata'
```

### Result

The resulting navbar would look something like this.

---

<nav>
   <a href="/">
    markata
   </a>
   <a href="https://github.com/WaylonWalker/markata">
    GitHub
   </a>
</nav>

---

## In your own template

If you want to continue using this method of maintaining your nav links with a
custom template, add this block to your template where you want your nav to
appear.

```html
<nav>
  {% for text, link in markata.config.nav.items() %}
  <a href="{{link}}">{{text}}</a>
  {% endfor %}
</nav>
```
