---
title: "Progressive article styles"
description: "Skip unnecessary feed-card CSS on ordinary reading posts."
published: true
tags: [performance, css, posts]
---


Markata-go avoids loading the feed-card stylesheet on an ordinary `post.html`
reading page when its article has no embedded feed or card markup.
Feed archives and posts embedding cards keep their styles, as do custom
templates. No content is hidden and the article text remains server-rendered.

This is an initial optimization, not a fully generated per-page critical CSS
system. If a custom post relies on card styling without markup the generator
can recognize, use a custom template and the regular stylesheet stack.

After a shared template change, run `markata-go build --clean` and inspect
the generated post `<head>`. Compare a representative text post, card-heavy
article, feed archive, and navigation between page types in real browsers.
