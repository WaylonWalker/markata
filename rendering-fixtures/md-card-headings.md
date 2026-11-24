---
title: "Card heading overflow regression"
description: "Card permalinks, nested headings, and independently scrollable code and comparison tables"
slug: "md-card-headings"
date: 2026-10-26
published: true
---

## An ordinary article heading

::: md-grid

:::: md-card
#### Title

A short card should never need a horizontal scrollbar.

::::

:::: md-card
#### A longer card heading that wraps across several lines on narrow screens

::::: md-stack
##### Nested heading

Permalinks remain available for nested headings too.

:::::

::::

:::: md-card
#### Wide code

```text
This intentionally long code line must remain scrollable inside the code block: abcdefghijklmnopqrstuvwxyz abcdefghijklmnopqrstuvwxyz abcdefghijklmnopqrstuvwxyz abcdefghijklmnopqrstuvwxyz
```

::::

:::

::: md-grid

:::: md-card
#### Wide comparison table

| Component | Before | After |
| --- | --- | --- |
| VeryWideComponentNameWithoutBreakOpportunitiesABCDEFGHIJKLMNOPQRSTUVWXYZ | BeforeValueWithoutBreakOpportunitiesABCDEFGHIJKLMNOPQRSTUVWXYZ | AfterValueWithoutBreakOpportunitiesABCDEFGHIJKLMNOPQRSTUVWXYZ |

::::

:::
