---
title: Homepage Lighthouse improvements
description: Local homepage performance design and validation approach.
date: 2026-10-15
published: false
tags: [performance]
---


The local waylonwalker.com homepage already scores 100 for accessibility,
best practices, and SEO. Its mobile render is delayed by the shared stylesheet
stack. Preserve its layout, theme controls, and content while reducing the
initial stylesheet cost.

Keep the complete palette stylesheet for compatibility. Generate a smaller
initial stylesheet containing the configured light/dark pair and picker
metadata. Restore saved or seasonal choices before paint by loading the full
stylesheet only when required. Load all preview colors when opening the picker.
Keep keyboard selection, reset, and navigation working.

Removing the picker would lose functionality. Inlining all theme CSS would
increase every HTML response. Splitting its stylesheet keeps the feature and
preserves browser caching. Use the baseline trace to select further changes.
Avoid editing post copy or reducing feed content without evidence.

The trace also shows full-size photo downloads and eager media inside inline
TIL cards. Add responsive candidates for resize-capable hosts, bounded TIL
summaries, and deferred rendering for the offscreen feed section. Consolidate
the homepage's stylesheets in cascade order. Load image zoom assets only when
the visitor opens an image. Keep the original posts and all eight TIL links.

Build the engine and complete site before and after. Serve both outputs with
the same gzip/cache policy. Run three cold mobile and desktop Lighthouse audits
using the same binaries. Keep uncompressed-server reports separately. Record
scores, FCP, LCP, TBT, CLS, screenshots, and exact build revisions.

Test palette artifacts, saved preferences, picker controls, dark/light mode,
and navigation. Review desktop and mobile screenshots. Keep the baseline
output available for local preview.

Refs #980.
