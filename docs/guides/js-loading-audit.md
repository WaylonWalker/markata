---
title: Audit JavaScript requested by real pages
description: Compare actual cold and repeat JavaScript resource requests before changing Markata's asset pipeline.
---

# JavaScript loading audit

The first phase of [issue #2071](https://github.com/WaylonWalker/markata/issues/2071) measures **what browsers request**, not everything copied into the generated site's static asset directory.

Serve a built Markata site or use a deployed test site. Provide representative URLs for a plain post, a post with code, a feed/archive, a docs page with TOC, and a gallery.

    node --test scripts/js-loading-audit.test.cjs

    # Use the Playwright installation from existing browser tests, or install Playwright.
    PLAYWRIGHT_MODULE=/path/to/playwright-core \
    CHROMIUM_EXECUTABLE=/usr/bin/chromium \
    node scripts/js-loading-audit.cjs \
      --url http://localhost:8000/example-post/ \
      --url http://localhost:8000/archive/ \
      --out test-results/js-loading-audit.json

The CLI requires Playwright (playwright or playwright-core; it also discovers an agent-browser global installation) and a working Chromium. No Playwright install is required for normal Markata builds.

For each URL it opens a fresh browser context, records a cold **full page navigation**, then navigates to the same URL again in the same context. JSON contains script resource URLs, script-request counts, observed transfer/encoded/decoded bytes, and executable inline-script character counts. Check the cold and repeat records; run this against the same site and hardware when comparing candidate changes.

## What it does not claim

- **0 bytes transferred does not prove a cache hit.** Cross-origin resources lacking Timing-Allow-Origin can report zero sizes. The report explicitly counts zero-transfer entries rather than categorizing them as cached.
- Resource Timing counts are not equivalent to total network bytes or JavaScript parse/evaluation CPU time. Script initiators and `.js`/`.mjs` paths are included, but extensionless fetch/module resources without a script initiator may be missed. Failed HTTP page navigations abort the audit instead of measuring an error page.
- Full-page repeat navigation does not measure Markata's custom SPA View Transitions. Add a separate interaction/navigation test for that.
- Scripts loaded after hovering, opening a theme menu, expanding search, or opening a graph are not measured unless triggered in the settle window.
- Bytes published to disk are not automatically downloaded by readers.

**Next work:** add route-specific interaction scenarios and trace analysis, then set a per-route JavaScript budget. Make lazy-loading and conditional-publication changes in separate reviewable PRs; do not rewrite benchmark fixture content to manufacture build speedups.
