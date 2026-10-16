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


## Browser-native navigation experiment

To compare the existing SPA transition engine with cross-document navigation,
set `[markata-go.view_transitions] native = true` and keep `enabled = true`.
The generated page opts into CSS `@view-transition { navigation: auto; }`,
emits a conservative intent-only Speculation Rules prefetch for sidebar and
post-navigation links, and **does not ship the custom view-transitions.js
DOM-swap runtime or its conditional-css.js companion**. Browsers without
native support use ordinary full-document link navigation. This is an opt-in
experiment until cold-mobile Lighthouse, browser-navigation and theme tests
show it is appropriate as the default. Setting `enabled = false` disables
both modes.
