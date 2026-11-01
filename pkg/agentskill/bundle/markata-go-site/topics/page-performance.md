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

A proposed static SVG backend publishes content-addressed external images
with intrinsic dimensions, but it is currently an **internal** package,
not an active site-build option. Once integrated, measure SVG file transfer
bytes (including embedded fonts), the SVG element/paint cost, and browser
CLS separately from the HTML-only page-weight report. Do not claim that
image dimensions alone certify zero CLS on a responsive layout.

Regular builds warn without failing. Use `--page-weight=false` only to skip
the scan; `--page-weight-html-kb=96` changes the post-size warning threshold.
Run a clean site rebuild after shared template changes because incremental
output may retain earlier HTML. Both template trees may require edits.


## Alternative feed sidebar hydration

When a post is in several feeds, the generator includes the current feed's
navigation list as normal HTML. Other feeds with verified static
`/<feed>/_sidebar/NNNN.json` sidecars contribute only metadata until someone
selects them. The browser fetches the relevant window on selection, and leaves
existing links intact if a sidecar request fails. Synthetic feeds without a
verified sidecar continue to embed their first window. This avoids trimming
the content a reader can explore and reduces duplicated article HTML.
