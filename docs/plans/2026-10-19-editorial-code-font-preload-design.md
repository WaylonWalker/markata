---
title: "Editorial Code Font Preload"
description: "Fetch the already-used weekly label font before layout requests it."
date: 2026-10-19
published: false
tags:
  - performance
  - fonts
---


Refs #1737.

The first editorial week label uses `--font-mono`, which bundled font packs map
to their code role. That face is currently discovered during layout through the
shared font stylesheet, later than the preloaded body and heading faces.

Three alternatives were considered: preload code everywhere, detect arbitrary
code markup per page, or preload only where the default editorial layout
necessarily renders monospace metadata. Use the last option. It has a concrete
font consumer and avoids adding downloads to code-free pages or guessing about
custom hidden markup. Article preload remains a separate potential improvement.

Cache two per-pack URL lists: the ordinary body/display/heading list and an
editorial list that additionally includes the code role. Reuse the existing
emitted-tier selection and source deduplication. Validate both lists and include
both in the HTML cache identity. Advance the cache schema version.

Both template trees choose the editorial list only for editorial pages with
nonempty weekly sections. With a picker, preload only the visitor's active pack;
without one, emit the default pack's editorial list. System code stacks and code
sources already selected by another role add no request. Empty and other layouts
keep their existing preload behavior. Font files, font CSS, glyphs, widths, fallback
metrics, and font-display behavior are unchanged.

Three alternating fixed-snapshot mobile Lighthouse pairs improved archive
Performance from median91 to94, and median LCP by about300ms. All A/BP/SEO scores
were100; candidate TBT was0 in every run. These prototype results require exact
implementation validation and a full site rebuild before release. The article
prototype improved median93 to94 but is outside this change's implementation.

Tests must cover emitted tier selection, shared-source deduplication, system
stacks, invalid/missing cache metadata and files, HTML cache invalidation,
picker selection, empty and classic feeds, and both template trees. The rendered
reference site must request the same font bytes with an earlier code request.

The bundled site skill was reviewed. No command, configuration, customization,
or site-author workflow changes require updating it.

Exact generated-output pairs subsequently scored92→96 with LCP3236→2709ms.
Earlier paint exposed a pre-existing toolbar enhancement shift: the initially
hidden Calendar view switcher appears after initialization, moving pagination.
For JavaScript-enabled editorial calendars, reserve the switcher's own dimensions
with visibility:hidden until initialization removes its hidden attribute.
A synchronous bootstrap marker limits this reservation to JavaScript-enabled
pages. Keep pending controls inaccessible, and preserve no-JavaScript geometry
and the final toolbar. Validate at narrow mobile, mobile, and desktop widths.

Final rebuilt artifact validation after the reservation fix: three alternating
mobile Lighthouse pairs scored 92/92/92 before and 96/96/96 after. Median LCP
improved from 3234ms to 2709ms. Candidate CLS and TBT were zero in every run;
accessibility, best practices, and SEO were 100 throughout. All six runs fetched
the same four font URLs and 143,839 transferred font bytes with no failed HTTP
requests. These are controlled local measurements, pending production rollout.
