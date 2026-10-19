---
title: "Closed Drawer Rendering"
description: "Defer closed navigation rendering without changing drawer behavior."
date: 2026-10-19
published: false
tags:
  - performance
  - themes
---

# Closed Drawer Rendering

Refs #1719.

The archive still spends more time laying out and styling content than running
JavaScript. Closed feed and table-of-contents drawers retain substantial
navigation trees even though their contents are outside the viewport.

Use native `content-visibility: hidden` on the contents of closed drawers only
when both content visibility and discrete transitions are supported. Restore
visible content immediately on opening so focus and active-item measurements
work. Delay hiding until the closing slide finishes. Preserve reduced motion,
saved-open desktop state, handles, typography, and all navigation links.
Unsupported browsers retain their existing rendering behavior.

Mobile flow layout and per-card containment prototypes did not meaningfully
reduce layout work. The drawer prototype reduced median archive layout from
339.849 ms to 288.009 ms across five alternating throttled runs, and style
recalculation from 276.861 ms to 230.836 ms. Five paired mobile Lighthouse runs
kept median archive performance at 91 and article performance at 93, with all
accessibility, best-practices, and SEO scores at 100. This is a rendering CPU
reduction, not evidence of a Lighthouse score improvement. Text/font rendering
remains the next investigation.

Browser regression coverage checks mobile and desktop opening, keyboard focus,
closing-animation visibility, reduced motion, and saved-open restoration.
Run the existing startup regressions and build the full reference site before
release. Earlier drawer experiments against the pre-font-lifecycle baseline
regressed article timing and were not shipped; these measurements use the
current baseline containing #1734.

No bundled agent skill update is needed: site-author configuration, commands,
and theme customization workflows are unchanged.
