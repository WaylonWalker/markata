---
title: "Try native post navigation"
description: "Experiment with browser-native cross-document transitions and conservative prefetch."
published: true
tags: [performance, navigation, experiment]
---

## Try browser-native navigation

For a staging or local site, add:

```toml
[markata-go.view_transitions]
enabled = true
native = true
```

Run a **clean build**. This replaces markata-go's JavaScript SPA interception
and DOM swapping with ordinary links plus native cross-document view
transitions on browsers that support them. Older browsers keep regular links.

In native mode the generated page can prefetch post/sidebar navigation
targets at conservative click intent using the Speculation Rules API.
This is deliberately **not prerendering** and does not fetch arbitrary
content at startup. Visitors who prefer reduced motion receive no
view-transition animation.

To return to established behavior, set `native = false` or omit it.

The default remains the existing navigation system while comparing
full-document first paint, return navigation, LCP, CLS, theme controls,
scroll positions, focus, and keyboard behaviors on multiple browsers.
Browsers lacking either feature must still navigate normally. This option
is experimental and does not guarantee better Lighthouse scores.
