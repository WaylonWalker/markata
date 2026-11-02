---
title: "Interactive Markdown checklists"
description: "Ordinary GFM task lists in Markata remain static, disabled, readable checkboxes."
date: 2026-10-21
published: true
tags:
  - documentation
---


Ordinary GFM task lists in Markata remain static, disabled, readable checkboxes.
To let a reader check off tasks, wrap one or more task lists with the
`md-checklist` container:

````markdown
::: md-checklist {data-id="deployment"}

- [x] Back up the configuration
- [ ] Run the build
- [ ] Verify the live site

:::
````

The built HTML remains ordinary semantic Markdown list items with disabled
checkboxes. On pages containing an `md-checklist`, Markata conditionally
loads a small JavaScript enhancer that enables just those checkboxes. There is
no global checklist runtime or JavaScript requirement for ordinary task lists.

`data-id` is optional. Supply a **stable and unique ID per list in a page**
to save progress in the browser's local storage under that page's pathname.
Without `data-id`, the list is interactive in the current session but does
not persist. If storage is unavailable, checking tasks still works and
nothing is uploaded to a server.

The task labels and original check states remain in the published Markdown.
Local progress does not rewrite the site, user account, feeds or shared copy.
The browser's native checkbox provides keyboard and screen-reader support;
list text is attached to each checkbox as an accessible label.

## Graceful fallbacks

- JavaScript disabled or blocked: boxes stay disabled and the list is readable.
- JavaScript enabled: only `md-checklist` containers become interactive.
- Partial-page navigation: newly loaded checklists initialize once.
- Content changes: if the count of tasks changes, outdated saved state is
  ignored rather than applied to the wrong items.
- Privacy: only a boolean completion array is stored locally; no task text
  is persisted and no network request is made.
