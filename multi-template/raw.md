---
title: Using Multiple Templates
description: Guide to using multiple templates in markata
jinja: false
template:
  index: post.html
  title: og.html

---

This page will be rendered multiple times.  This is the normal page, check out
[title](./title). It uses the og.html template.  This is achieved through the
use of the template option in the frontmatter.

``` md
---
title: Using Multiple Templates
description: Guide to using multiple templates in markata
jinja: false
template:
  index: post.html
  title: og.html


---
```
