---
title: "Adaptive Loading Hardening"
description: "Auto detection, browser correctness, repeat transfer measurements, and PR readiness."
date: 2026-10-08
published: true
tags:
  - performance
  - adaptive-loading
  - benchmark
---


An audit of the existing full-document View Transition navigator found that it
could reset loading attributes to server defaults. The final follow-up preserves
those attributes, initializes new-route authored media, and restarts a bounded
timing window on its existing completion event. Chromium integration exercises
that actual navigator for Auto and both manual modes. The CDP/Lighthouse tables
below measure checkpoint `aa86d8ec`; this route-entry follow-up adds about 68 gzip
bytes to the controller and runs on subsequent route completion, not initial
paint. Initial-load policies and representative media markup are unchanged.

Issue [#1524](https://github.com/WaylonWalker/markata-go/issues/1524), existing
`feat/adaptive-client-loading` branch. This continues the design and implementation
commits rather than replacing their architecture.

## Benefit

The final controller is 12,072 source bytes, 6,408 generated/minified bytes,
2,393 bytes gzip (level 6), and 2,125 bytes Brotli.
The generated main stylesheet delta is 2,134 bytes uncompressed and 433 bytes
gzip. One deferred controller request is added; CSS uses the existing request.
Overall route request counts also reflect removed optional player assets.

Real HTTP transfer tests prove deferred startup before `load`, timing collection
with PerformanceObserver and Network Information disabled, constrained Auto,
four-good recovery, offline/online hysteresis, manual precedence, mismatch
Keep/Use Auto actions, unchanged loaded image sources, stable video geometry,
and no playback restart. Cached images cannot masquerade as fast network
measurements. A page stops after twelve votes; another page can provide new
measurements. Generic healthy 4g hints cannot erase measured poor confidence;
browser Save Data intent keeps Auto constrained even with fast measurements.

### CDP repeat measurements

Each cell is a median of three fresh contexts. Network Information is disabled;
Auto has no manual preference or seeded confidence. Mobile is 390×844, with
150 ms latency and 192 KiB/s downstream for constrained cases. Desktop is
1365×900. Feed/image routes receive the same four 700-pixel scrolls. LCP and
CLS are captured before scrolling; transfer totals include exploration.

Baseline is `73c4f606`. Final candidate after-runs were repeated after reserving
poster geometry. All runs use gzip responses from the same local server and
live external media, with a 2.5-second post-load observation window. Baseline
and candidate are built from an isolated copy of the current real site's inputs,
including its local edits, with separate cache/output directories. The source
site is unchanged. Encryption is disabled equally for both builds; Searchcraft
index publication is disabled. These totals are not directly comparable to the
older report's server/compression setup. First-party inputs and build output
are deterministic; external media and server redirect latency still vary.

Raw per-run metrics, paint attribution, confidence, and aborted requests are in
`docs/reports/adaptive-loading-hardening-results.json`. Aborted baseline video
requests are included in request counts; totals include completed CDP transfers.
This is a bounded observation window, not the lifetime download volume.

| Mode | Route | Bytes before → after | Requests | Media bytes | LCP seconds | CLS | Final policy |
|---|---|---:|---:|---:|---:|---:|---|
| auto-constrained | /vimgrep-open-buffers/ | 652,982 → 649,995 | 50 → 49 | 0 → 0 | 1.344 → 1.812 | 0.0266 → 0.0276 | constrained |
| auto-constrained | /til/ | 723,982 → 688,951 | 63 → 56 | 32,520 → 0 | 1.420 → 1.412 | 0.0674 → 0.0681 | normal |
| auto-constrained | /2025-nas/ | 1,334,700 → 1,331,512 | 60 → 59 | 0 → 0 | 3.996 → 3.952 | 0.0073 → 0.0075 | constrained |
| auto-constrained | /tmux-pop-size/ | 1,387,406 → 886,531 | 57 → 57 | 502,099 → 0 | 4.044 → 4.036 | 0.0073 → 0.0075 | constrained |
| auto-constrained | /bloatware-is-dying/ | 766,187 → 763,718 | 60 → 57 | 0 → 0 | 2.844 → 2.948 | 0.0109 → 0.0111 | constrained |
| auto-normal | /vimgrep-open-buffers/ | 652,978 → 649,995 | 50 → 49 | 0 → 0 | 0.672 → 0.428 | 0.0072 → 0.0071 | normal |
| auto-normal | /til/ | 952,811 → 688,694 | 60 → 56 | 261,857 → 0 | 0.732 → 0.560 | 0.0618 → 0.0071 | normal |
| auto-normal | /2025-nas/ | 1,334,030 → 1,331,506 | 60 → 59 | 0 → 0 | 0.760 → 0.444 | 0.0072 → 0.0071 | normal |
| auto-normal | /tmux-pop-size/ | 1,387,390 → 886,532 | 57 → 57 | 502,099 → 0 | 0.776 → 0.472 | 0.0072 → 0.0071 | normal |
| auto-normal | /bloatware-is-dying/ | 1,539,251 → 763,718 | 60 → 57 | 773,075 → 0 | 0.992 → 0.520 | 0.0072 → 0.0071 | normal |
| save-data-constrained | /vimgrep-open-buffers/ | 652,998 → 649,995 | 50 → 49 | 0 → 0 | 1.632 → 1.536 | 0.0266 → 0.0279 | constrained |
| save-data-constrained | /til/ | 723,756 → 688,683 | 63 → 56 | 32,512 → 0 | 1.520 → 1.400 | 0.0674 → 0.0702 | constrained |
| save-data-constrained | /2025-nas/ | 1,334,510 → 1,331,521 | 60 → 59 | 0 → 0 | 3.932 → 3.972 | 0.0073 → 0.0084 | constrained |
| save-data-constrained | /tmux-pop-size/ | 1,387,402 → 886,536 | 57 → 57 | 502,095 → 0 | 4.044 → 4.040 | 0.0073 → 0.0084 | constrained |
| save-data-constrained | /bloatware-is-dying/ | 766,187 → 763,718 | 60 → 57 | 0 → 0 | 2.772 → 2.988 | 0.0109 → 0.0141 | constrained |
| desktop-normal | /vimgrep-open-buffers/ | 652,986 → 649,995 | 50 → 49 | 0 → 0 | 0.724 → 0.548 | 0.0000 → 0.0000 | normal |
| desktop-normal | /til/ | 953,124 → 688,679 | 61 → 56 | 262,085 → 0 | 0.848 → 0.904 | 0.0342 → 0.0052 | normal |
| desktop-normal | /2025-nas/ | 1,334,502 → 1,331,513 | 60 → 59 | 0 → 0 | 0.812 → 0.484 | 0.0000 → 0.0000 | normal |
| desktop-normal | /tmux-pop-size/ | 1,387,905 → 886,530 | 58 → 57 | 502,096 → 0 | 1.504 → 0.960 | 0.0138 → 0.0000 | normal |
| desktop-normal | /bloatware-is-dying/ | 1,539,761 → 763,718 | 61 → 57 | 773,093 → 0 | 1.012 → 0.592 | 0.0052 → 0.0058 | normal |

Auto independently detects the constrained profile on the text, image-heavy,
video, and embed routes. The feed has only one usable poor sample and remains
Normal: many remote media responses omit Timing-Allow-Origin. It still saves
about 35 KB and seven requests through the conservative static policy. This
is not a claim that Auto classified the feed as constrained. Auto normal
also stays Normal when there are fewer than four actual good media/navigation
samples; generic 4g alone does not grant autoplay.

Constrained video avoids approximately 501 KB total and 502 KB media transfer.
Constrained embed transfer is essentially unchanged (about 2.5 KB less), with
three fewer requests; that baseline did not transfer embed media in the window.
Normal embed runs do show substantial deferred media savings. The controller
cost and media savings are separate measurements.

### Lighthouse (separate method)

Lighthouse 12.8.2 mobile simulated throttling, performance category, Chromium.
Comparable before/after configurations; these are snapshots for four routes,
and three-run medians for video. Never subtract these from CDP measurements.
Totals and counts come from Lighthouse network-requests, including redirects.

| Route | LCP seconds | CLS | Total bytes | Media bytes | Requests |
|---|---:|---:|---:|---:|---:|
| text | 4.956 → 5.105 | 0.0000 → 0.0000 | 652,988 → 649,995 | 0 → 0 | 50 → 49 |
| til | 5.482 → 5.189 | 0.0000 → 0.0647 | 5,214,156 → 689,525 | 4,522,564 → 0 | 60 → 56 |
| 2025-nas | 7.385 → 6.382 | 0.0000 → 0.0075 | 1,334,327 → 1,331,554 | 0 → 0 | 60 → 59 |
| video | 6.157 → 6.229 | 0.0000 → 0.0086 | 1,387,879 → 887,168 | 502,101 → 0 | 57 → 57 |
| bloatware-is-dying | 5.855 → 5.708 | 0.0000 → 0.0075 | 1,540,179 → 763,718 | 773,608 → 0 | 60 → 57 |

### Video LCP and stability

The original 10.78 → 13.85 s constrained regression did not reproduce. Final
constrained CDP video medians are 4.044 → 4.036 s; desktop is
1.504 → 0.960 s. Absolute times changed with response compression and harness
conditions, so this establishes a comparable local delta, not a universal
speedup. Mobile's measured LCP is a paragraph, while desktop may count a
video frame or its new poster. First candidate desktop medians were
1.50 → 1.56 s, with substantial per-run variance.

A Lighthouse trace shows the poster API request returning a 307 after about
629 ms, followed by a 78 ms, 3.3 KB thumbnail response. Both poster and video
requests were Low priority; poster work starts earlier, not later. This
redirect/server latency explains why a tiny poster can still be a late desktop
paint. No evidence attributes the old regression to controller initialization.
The historical run lacks a trace, so its exact cause cannot be established.

Layout-shift attribution identified content below the unsized video moving
when the poster arrives. Reserving a 16:9 frame for derived 1200×675 posters
removed that shift in final desktop video runs (CLS about .00002). Raw authored
height/inline styles remain authoritative. Loaded images are not replaced and
policy updates do not pause, rewind, reload, or replace actively requested media.
The controller does not enable animations in response to improved connectivity.
Small mobile font/layout shifts remain site-level baseline behavior.

## Tradeoffs

Generic network hints bootstrap conservatively rather than declaring a fast
connection. Timing needs useful same-origin/Timing-Allow-Origin transfers;
opaque, cached, tiny, and neutral-rate resources do not vote. Observations end
after thirty seconds or twelve votes, so later changes may require another page.
No continuous polling or framework is added. Session confidence expires after
ten minutes. The new observer adds bounded lifecycle work for late resources.

Explicit raw autoplay and configured Markdown autoplay survive as
`data-authored-autoplay="true"`. Full Quality or measured fast Auto at page entry can start
unrequested authored media. Fresh fast evidence during a page view affects
subsequent pages and does not suddenly animate an existing poster. Save Data, unknown/Normal/constrained Auto, and
reduced motion suppress automatic playback. Raw video without autoplay never
inherits it merely from site-wide config. Native controls and sources remain.
With JavaScript disabled, authored autoplay also requires interaction: this
is the deliberate conservative-baseline tradeoff. Browser autoplay permissions
can still reject playback; the controller does not retry or restart it.

Chromium, Firefox (installed Playwright Firefox 1511), and Brave integration
passed with missing Network Information, real timing fallback, all user modes,
mismatch UI, video interaction, recovery, offline/online, sample caps, and cache
exclusion. Brave's Network Information is deliberately disabled in the portable
path tests; this is not a claim about every privacy setting. Safari/iOS were
not run. No broader device study is claimed.

## Readiness

Adaptive loading: **READY FOR PR** after the recorded validation and clean commit.

- `go test ./...` passes using disk-backed TMPDIR.
- `node --test scripts/adaptive-loading.test.cjs`: 13 pass.
- `node --test scripts/adaptive-loading.browser.test.cjs`: Chromium/Firefox/Brave pass.
- `just lint-new`, `go fmt ./...`, and `git diff --check` pass.
- Generated controller stays below the 3 KB gzip ceiling.
- Spec, configuration/plugin/user docs, and bundled site-agent guidance updated.

Initial test/build attempts hit the temporary-filesystem quota, and a moved
output root was rejected as a symlink. Those incomplete builds and failing
runs are discarded. Final corpus builds use direct disk paths; final checks
use disk-backed temporary files. No unrelated heading_anchors change was made.

## Next steps

Review the focused adaptive branch without merging. Keep the poster redirect
as a measured deployment/media-service follow-up rather than expanding this
controller. Recheck on Safari/iOS when a real device is available. Investigate
partial navigation separately only after this readiness gate.

## Reproduction

`adaptive-loading.benchmark.cjs BEFORE_URL AFTER_URL RESULTS.json` runs three
fresh contexts per mode/route. Serve both output trees with identical gzip
handling; the historical Python plain-file server is not the same measurement
configuration. Optional `BENCHMARK_MODES`, `BENCHMARK_ROUTES`, and
`BENCHMARK_VERSIONS` limit reruns. The browser integration test uses Playwright
bundled with agent-browser, or `PLAYWRIGHT_MODULE` for another installation.
Set `ADAPTIVE_FIREFOX_EXECUTABLE` when the installed Firefox path differs from
the Playwright version's default. Test SVGs use real delayed server responses,
not injected vote/state samples. The tiny WebM fixture is a three-second,
160×90 blue frame generated with FFmpeg/libvpx-vp9; it tests actual media playback.


Keep native controls enabled (or supply an accessible play action) for deferred
videos. A video deliberately authored without controls needs its own interaction
path in Save Data and no-JS mode; autoplay intent alone is not a play control.
