---
title: "Adaptive Client Loading Benchmark"
description: "Baseline and post-change browser measurements for adaptive loading on waylonwalker.com."
date: 2026-10-07
published: true
tags:
  - performance
  - adaptive-loading
  - benchmark
---


Historical measurements from implementation commit `bbb5cb53`. The hardening
report supersedes its policy, autoplay, browser coverage, and readiness findings.
See [Adaptive loading hardening](../adaptive-client-loading-hardening/).

This report records browser measurements before and after adaptive client loading on the real `waylonwalker.com` corpus. The source site was not edited. Both builds used the worktree binary and wrote to temporary output directories.

## Reproduce

The site normally builds with `just build`, which runs `markata-go build -m config/no-tailwind.toml`. The benchmark used the same site and config with a separate output path and cache:

```sh
markata-go build --site-dir ~/git/waylonwalker.com \
  -m ~/git/waylonwalker.com/config/no-tailwind.toml \
  -o /tmp/markata-adaptive/output
```

Serve the output with `python3 -m http.server 8766 --directory /tmp/markata-adaptive/output`. The Chromium harness used Puppeteer already bundled with Lighthouse, 390×844 with 150 ms latency and 192 KiB/s downstream for mobile, and 1365×900 without shaping for desktop. The constrained run selected Save Data in local storage and rapidly scrolled the feed. Counts and byte totals are CDP transferred bytes, including responses captured during that scroll. LCP and CLS are from the same short run and are directional; the rapid scroll and third-party timing make CLS especially noisy. DOM text is the number of body text characters available at `DOMContentLoaded`, not a paint timestamp.

## Representative mobile run

| Page | Before requests / bytes | After requests / bytes | Before / after media bytes | Before / after LCP | Before / after CLS |
|---|---:|---:|---:|---:|---:|
| Text article `/vimgrep-open-buffers/` | 50 / 1,401,880 | 49 / 1,405,312 | 0 / 0 | 5.90 / 5.93 s | .0074 / .0086 |
| Feed `/til/` after rapid scroll | 63 / 1,604,478 | 55 / 1,578,159 | 32,513 / 0 | 5.91 / 5.51 s | .442 / .395 |
| Video article `/tmux-pop-size/` | 57 / 2,843,230 | 55 / 2,348,743 | 501,728 / 0 | 10.78 / 13.85 s | .0074 / .0086 |
| Embed article `/bloatware-is-dying/` | 60 / 2,252,247 | 57 / 2,255,430 | 0 / 0 | 9.74 / 9.78 s | .056 / .0082 |

The text article added about 3.4 KB net and one fewer request. Across the four pages, JavaScript transfer increased by 1.3–1.4 KB; the generated controller is 5,281 bytes uncompressed and 2,031 bytes gzip. It is deferred, non-blocking, and shares the existing page CSS. The video article avoided about 494 KB total, including its 502 KB autoplay transfer, while adding a right-sized poster. The feed avoided 26 KB net and eight requests in this particular rapid-scroll pass. Its remaining images still load as native lazy images when scrolling brings them near the viewport.

The embed page did not transfer media in either CDP run; this run therefore does not support a media-savings claim for that page. It did avoid three requests overall, while its byte total changed by about +3 KB. Its external embed behavior is timing-sensitive. The controller does not make claims about bytes that were not observed.

## Desktop normal run

| Page | Before requests / bytes | After requests / bytes | Before / after media bytes | Before / after LCP |
|---|---:|---:|---:|---:|
| Text article | 50 / 1,402,359 | 49 / 1,405,312 | 0 / 0 | .93 / .59 s |
| Feed | 59 / 1,605,506 | 53 / 1,577,633 | 32,165 / 0 | .80 / .51 s |
| Video article | 57 / 2,843,745 | 55 / 2,348,743 | 501,622 / 0 | .83 / 1.45 s |
| Embed article | 59 / 2,992,113 | 56 / 2,221,185 | 773,597 / 0 | .62 / .52 s |

The video article’s raw HTML previously autoplayed a remote MP4. It now keeps native controls and source markup but defaults to no autoplay and `preload="none"`, so its poster is the initial visual. This accounts for the large stable transfer reduction on both desktop and mobile. The embed page’s desktop reduction is observed in this run; it should be rechecked against a second cache/network state before treating that amount as typical.

## Lighthouse baseline

Mobile Lighthouse baseline reports (12.8.2, simulated throttling) recorded 50 requests / 1.40 MB / 6.37 s LCP for the text page; 63 / 4.32 MB / 7.93 s for the feed (2.74 MB media); 57 / 2.85 MB / 11.86 s for the video page (502 KB media); and 60 / 3.00 MB / 11.42 s for the embed page (774 KB media). These are a separate measurement method from the CDP table above. A comparable post-change Lighthouse run was not captured, so no Lighthouse delta is claimed.

## Stability and UX checks

Chromium mobile and desktop generated pages kept article/card geometry stable when the policy changed. The Save Data video run made zero media requests; the generated HTML contained no `autoplay`, used `preload="none"`, and retained a poster and native controls. A Chromium UI smoke check seeded a confident constrained recommendation with manual Full Quality selected: the small status affordance appeared, opened the settings panel with “Full Quality is on. Your connection appears constrained.”, and “Use Auto” cleared the local override. No full-page screenshot is included because the UI is a small native disclosure/status control rather than a separate screen.

The rendered controller gzip size is 2,031 bytes. The added CSS is part of the existing stylesheet (no extra CSS request); the appended rules gzip to roughly 0.5 KB. The control layer adds one script request on pages that do not already bundle it, but the site’s overall request count was flat or lower because optional player CSS and automatic external video loads were removed/deferred. No long-lived observer or polling loop is used; resource timing samples stop after six observations.

## Browser coverage and limits

| Browser | Result |
|---|---|
| Chromium | Mobile constrained and desktop normal transfer runs; manual mismatch and Auto action smoke-tested. |
| Chrome | Not separately installed; Chromium run covers the shared engine behavior. |
| Brave | Not installed locally. Network Information API is optional and not required. |
| Firefox | Not installed locally. Portable resource timing, online/offline, and native loading paths are used without browser detection. |
| Safari/iOS | Not run. The baseline uses native image/video/HTML behavior and avoids reliance on Chromium-only APIs. |

This is a local corpus benchmark, not a broad device study. Timing variation, browser cache state, third-party responses, and site-level autoplay configuration can affect totals. Sites that explicitly set `[markata-go.md_video] autoplay = true` retain that opt-in behavior. The primary tradeoff is that autoplay is no longer the default; readers press native video controls to start playback. More aggressive image deferral, embed-specific click facades, and a repeat Lighthouse after-run remain follow-up work.
