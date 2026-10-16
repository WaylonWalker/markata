# Investigating heavy generated pages

Run `markata-go build --clean`, then inspect
`output/.markata/page-weight.json` (honor a configured `output_dir`).
The array `pages` is sorted largest first. Focus on `kind: post`,
`findings`, and the byte breakdown rather than guessing from source CSS sizes.
Use the stable finding `code`, measured `bytes`, `budget_bytes`, and
`advice` to pick actionable work. Examples include duplicated post copies,
sidebar feed JSON, and optional inline application scripts.

The JSON measures **uncompressed emitted HTML source**, not actual network
transfer, render time, total asset bytes, or Lighthouse score. Verify
network/LCP/TBT/CLS separately. Prefer fixes in the static generator over
browser runtime work. Do not cut useful content to meet the budget.

Regular builds warn without failing. Use `--page-weight=false` only to skip
the scan; `--page-weight-html-kb=96` changes the post-size warning threshold.
Run a clean site rebuild after shared template changes because incremental
output may retain earlier HTML. Both template trees may require edits.


## Plain post card stylesheet

Stock `post.html` reading pages with no feed/card markup omit
`css/cards.css` on initial paint. Embedded feeds/cards, custom templates,
and unknown markup must preserve their styles. Inspect generated HTML
`<link rel="stylesheet">` tags and check real navigations from plain posts
to card-heavy feeds; conditional-css.js handles the existing SPA lane.
Run clean builds and compare actual Lighthouse LCP rather than assuming
that one skipped stylesheet guarantees a score improvement.
