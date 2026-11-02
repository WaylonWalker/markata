---
title: "Slow 3G Harness A/A Validation on WaylonWalker.com"
description: "Real generated-site A/A measurements validating the Slow 3G filmstrip harness."
date: 2026-10-13
published: true
tags:
  - performance
  - benchmarking
  - validation
---

## Slow 3G A/A validation

This report validates PR [#1554](https://github.com/WaylonWalker/markata-go/pull/1554)
against a fresh build of the actual WaylonWalker.com working source. These are
A/A measurements for harness noise and filmstrip inspection. No site
optimization was included, so this report makes no A/B performance claim.

## Pinned inputs

- Site source commit: `82ae5bcfc67674abfd2e5a169ca9551e2255ff57` with local
  working-tree edits and untracked site posts included in an isolated snapshot.
- SHA-256 of the exact copied build-input tree:
  `ba2a97aa9f0b4a870555c19c6b5bcf3e055b59e34f8567c46449feedbbfca346`.
- Markata Go build revision: `0af03c725fecacec934ba1d51649ca7a86155e01`.
- Harness capture revision: `fd66fbd0eb1480a69a5ea2dbebac93d5cb28d9e5`.
- Browser: Chromium `153.0.8010.52`; Playwright from the harness PEP 723
  dependency declaration; Python `3.14.7`.
- Fresh generated output: 3,540 discovered posts, 3,518 rendered, 421 feeds;
  build completed in 45.35 seconds. Output size was 2.4 GiB.

The site checkout contains private encrypted source. Its configured
`MARKATA_GO_ENCRYPTION_KEY_*` values were passed directly from the local
`.env` to the isolated build process. Secret values were not recorded. The
snapshot excluded `.env`, Git metadata, previous output and build caches.

## Capture setup

The same generated output directory was served from ports 8000 and 8001 with
`scripts/slow3g_site_server.py`. Each browser run used a new context, disabled
the browser cache and used the harness defaults: 390×844 viewport, 400 ms
latency, 51,200 bytes/second up and down, and 4× CPU slowdown. Each route ran
five alternating A/B and B/A pairs, with a 16-second observation window and
frames requested at 1, 2, 5, 8, 12 and 15 seconds. The A and B labels refer to
the two server ports; both serve identical bytes.

## Results

Median values below are across five valid runs for each port. Encoded bytes are
completed transfers during the 16-second window, not bytes received at paint.
First-readable times are manual estimates from the saved frames, rounded to
the checkpoint interval.

| Route | Valid runs per port | Median FCP A/B | Median LCP A/B | Median CLS A/B | Median bytes A/B | Median requests A/B | First readable main content |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `/` | 5/5 | 6292/6204 ms | 6292/6204 ms | 0.010765/0.010765 | 613690/613514 | 66/66 | About 8 s; intro/about text is visible. |
| `/archive/` | 5/5 | 5120/5136 ms | 7504/7392 ms | 0.000337/0.000337 | 301620/301582 | 32/32 | About 8 s; first post title and snippet are readable. |
| `/blogging-for-me/` | 5/5 | 4072/4064 ms | 4280/4276 ms | 0.010112/0.010112 | 372229/372229 | 47/47 | By 5 s; article title and opening paragraph are readable. |
| `/shorts/` | 5/5 | 972/972 ms | 2848/2836 ms | 0.006164/0.006164 | 53185/53178 | 8/8 | By 5 s; the image is sharp and its “Bat Cap” title is readable. At 2 s the media is still blurred. |
| `/reader/` | 0/5 | — | — | — | — | — | The 5 s frame shows the Reader heading and an empty-state message, with 0 entries. Runs are invalid as described below. |

Median absolute paired FCP differences for the four valid routes were 84 ms
(home), 88 ms (archive), 24 ms (article) and 8 ms (shorts). Median absolute
paired LCP differences were 84 ms, 224 ms, 40 ms and 36 ms respectively.
These describe same-build measurement variability; they are not performance
gains. Archive request counts varied between 32 and 57 in individual runs,
which is one reason to inspect raw `results.json` instead of relying only on
the medians.

Largest completed resources observed included the gzip archive document at
266,979 bytes, `source-sans-3-prose-core.woff2` at 79,495 bytes,
`caveat-brush-display-core.woff2` at 47,883 bytes, the local `/8bitcc.png` at
56,890 bytes, and the shared palette CSS at 30,787 bytes. The home and shorts
routes also fetched images from `dropper.waylonwalker.com`; those remain
external to the local test server. Per-request encoded and decoded transfer
lengths, URLs, types and completion state are retained in each `results.json`.

### Invalid reader route

All ten reader runs were marked invalid. The freshly generated `reader/index.html`
requests `/_pagefind/pagefind-ui.css` and `.js`, but this output has neither
file; the server correctly returns 404 with HTML MIME. The harness records two
failed requests and two invalid asset responses per run. This is missing output
from the site build, not a reason to suppress the harness checks. Also, the
fresh page contains zero reader entries, so it does not represent a populated
production reader. Reader timings are excluded from the table and from any
noise estimate.

## Reproduction

Use a clean snapshot of the same site source state, including the intended
working-tree content changes. Keep `.env` out of the snapshot and provide the
site's encryption keys only to the build process through environment
variables. Pin the Markata Go and harness revisions above.

```sh
rsync -a --delete \
  --exclude='.git' --exclude='.env' --exclude='.env.*' \
  --exclude='output/' --exclude='cache/' --exclude='.markata/' \
  --exclude='.markata-cache/' --exclude='.steam_cache/' \
  --exclude='.output.markata-serving/' --exclude='node_modules/' \
  "$SITE_CHECKOUT/" "$SITE_SRC/"

git -C "$MARKATA_REPO" worktree add --detach "$MARKATA_BUILD" \
  0af03c725fecacec934ba1d51649ca7a86155e01
go -C "$MARKATA_BUILD" build -o "$MARKATA_BIN" ./cmd/markata-go

"$MARKATA_BIN" build --clean \
  -c markata-go.toml \
  --site-dir "$SITE_SRC" \
  -o "$SITE_OUT"

python scripts/slow3g_site_server.py --root "$SITE_OUT" --port 8000
python scripts/slow3g_site_server.py --root "$SITE_OUT" --port 8001

uv run --script scripts/slow3g_compare.py \
  --chromium /usr/bin/chromium \
  --before http://127.0.0.1:8000/ \
  --after http://127.0.0.1:8001/ \
  --runs 5 --observe-ms 16000 --timeout-ms 26000 \
  --checkpoints 1000 2000 5000 8000 12000 15000 --min-frames 5 \
  --output docs/reports/slow3g-1554/aa-home
```

Repeat the capture command with both route URLs changed to `/archive/`,
`/blogging-for-me/`, `/shorts/` and `/reader/`, and use separate output
directories. Inspect the saved `filmstrip.html` and frames for readability.
Keep failed requests and wrong-MIME responses in the report as invalid runs.

## Preserved artifacts

Each `aa-*` directory contains `results.json`, `summary.json`, `summary.csv`,
`filmstrip.html`, all 60 requested checkpoint frames and the contact sheets
used for the manual readability review. The reader directory preserves its
failed requests and invalid-asset diagnostics as well.
