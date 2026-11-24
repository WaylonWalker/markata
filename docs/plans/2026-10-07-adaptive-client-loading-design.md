---
title: "Adaptive Client Loading Design"
description: "Design for connection-aware media loading with a stable zero-JavaScript baseline."
date: 2026-10-07
published: false
tags:
  - performance
  - design
---

## Adaptive Client Loading

Issue: #1524

## Context

Markata-Go already emits native image and video elements, supports image
variants, and conditionally loads theme scripts. The generated sites also use
resource hints and video previews. The feature should improve mobile behavior
without turning every page into a JavaScript application or mutating content
that a reader has already seen.

Three approaches were considered:

1. **Static-only defaults.** This costs no runtime bytes, but cannot provide
   user-selected modes or learn from actual connection behavior.
2. **A small global controller.** This supports mode choice, conservative
   session learning, and future-resource decisions. It adds one small script
   request to each page.
3. **Separate controller per media type.** This keeps code conditional, but
   duplicates policy and makes settings inconsistent across components.

Use the second approach. One controller owns preference, confidence,
hysteresis, and the effective policy. Existing conditional loaders remain
responsible for feature-specific code.

## Behavior

The public modes are Auto, Save Data, and Full Quality. Auto is the default.
Explicit user choice always determines the effective policy and is stored in
local storage. A small session record may retain recent Auto confidence across
page navigations, with a short lifetime. Network Information API values are
optional hints; navigation/resource timing and online state work without them.

The controller requires repeated observations to change Auto state. Poor
signals enter constrained mode sooner than good signals restore full quality.
Manual/Auto mismatch prompts need stronger confidence than internal policy
changes. A dismissal remains quiet for the session. Manual Full Quality on a
confidently constrained connection may show one restrained prompt sooner.

The controller sets a stable data attribute and exposes the effective mode to
optional components. It does not poll or scan the document repeatedly. It
collects a small bounded sample of meaningful resource timings, then settles.
Offline state is immediate; online recovery returns through normal hysteresis.

## Media and stability

Static HTML remains useful without JavaScript. Images keep native lazy loading,
intrinsic dimensions, and responsive sources. Video elements keep a poster,
controls where appropriate, and native sources, but previews do not autoplay or
preload bytes. Article video uses `preload="none"` by default so browser-native
playback still works on interaction. Existing author configuration can request
autoplay explicitly.

Cards and embeds reserve their existing geometry. The controller never replaces
an image source, changes dimensions, restarts playback, or inserts a large
placeholder after initial render. Optional external players load only after
reader interaction where practical. Static resource hints remain limited to
critical assets; speculative third-party connection work must not be triggered
by mere mention of a URL.

## UI and implementation

Add a compact loading preference control to the existing header controls, with
a native disclosure-style panel and a polite mismatch status. It must work on
small screens, with keyboard input, and without relying on color alone. The
controller is a standalone theme asset, measured after minification and gzip.
The desired budget is 2 KB compressed; 3 KB is the hard ceiling.

Add pure policy tests for signal classification, hysteresis, overrides,
persistence decisions, and offline/online transitions. Add template tests for
native media defaults and verify both active template trees. Use the
waylonwalker.com site as the real corpus and compare desktop plus constrained
mobile loads in Chromium. Report browser automation limits for Firefox and
Brave when unavailable locally.

## Hardening decisions (2026-10-08)

Keep the existing controller and static media architecture. Generic connection
hints initialize a recommendation without contributing measured confidence.
Browser Save Data is explicit constrained intent. Real timing observations take
priority, with two poor votes to constrain and four consecutive good votes to
recover; manual preferences remain absolute. Preserve constrained recommendation
in session storage during recovery, reset fast confidence offline, and allow new
measurements on each page rather than consuming the entire session's budget.

Collect timings at load when deferred startup precedes load, plus a resource
observer for late work. Deduplicate the load snapshot and observer entries. Stop
at twelve measured samples per page or thirty seconds. Ignore cache hits and
opaque/insignificant transfers. The 2 Mbps/4 Mbps thresholds leave a neutral band
and cover the benchmark's 192 KiB/s profile without requiring browser detection.

Preserve raw HTML's own autoplay intent and Markdown's configured intent in a
data attribute. Initial markup uses preload none without native autoplay. Only
Full Quality or measured fast Auto may start unrequested authored media;
reduced motion suppresses automatic playback. Policy changes never rewind,
pause, or replace media already loading or playing. The no-JS tradeoff is
explicit: native controls work, but authored autoplay requires interaction.

Validation uses real delayed HTTP image responses in Chromium, Firefox, and
Brave, including a disabled PerformanceObserver to prove the load fallback.
Corpus comparisons use isolated copies of the real site's current inputs,
fresh browser contexts, three repetitions, and medians. CDP and Lighthouse
results remain separate. Video LCP attribution includes the actual paint element
and poster resource timing, rather than assuming transfer savings improve LCP.


Auto honors authored autoplay when measured confidence is already fast at page
entry. Learning that the network is fast during a page view affects subsequent
pages; it does not suddenly animate an existing poster. An explicit Full Quality
choice may start unrequested authored video on the current page. Active media
remains untouched in either case.
