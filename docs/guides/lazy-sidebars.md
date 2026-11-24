---
title: "Progressive feed sidebars"
description: "How markata-go loads alternate feed navigation only when you need it."
published: true
tags: [performance, navigation, feeds]
---


A post's current feed remains in the HTML for immediate browsing and
no-JavaScript navigation. If the post belongs to several feeds, other feed
post lists are deferred **when a verified public sidebar sidecar exists**.
Choosing another feed loads the matching `/<feed>/_sidebar/NNNN.json` window
and enables normal feed cycling and incremental scrolling.

The browser preserves existing navigation if fetching an alternate feed fails;
select the feed again to retry. Synthetic and private-unavailable feeds without
a safe sidecar are not deferred.

Use `markata-go build --clean` to refresh both HTML and sidebar assets when
testing changes to shared templates. Inspect the new
`output/.markata/page-weight.json` report for changes in embedded JSON bytes.
No full feed is arbitrarily capped: additional navigation is available on
demand.
