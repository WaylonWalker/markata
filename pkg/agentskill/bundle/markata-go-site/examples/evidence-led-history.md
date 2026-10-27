# Evidence-led history: a visual storytelling pattern

Use this optional example for multi-period project histories and retrospectives supported by source records. It explains **what changed, why it mattered, and what the record actually supports**. It is not a new layout, theme, or mandatory template. Prefer `editorial-visual-composition.md` for decision workflows.

## Selection and evidence rules

- Open with an unboxed thesis, scope and date convention. Distinguish reporting dates from deployment dates.
- Choose periods based on meaningful changes, not one chapter per commit. Use short chapter introductions before visual timelines.
- Orient readers with concise linked overview cards. Use `road editorial story` for sourced narrative milestones, not for duplicating the overview.
- Describe effects, uncertainty and source limits. Keep incomplete work and incidents separate from finished outcomes.
- Show a chart only when measurements are comparable, have explicit units, and cite evidence; otherwise omit it.
- Remove all placeholders before publishing. Verify mobile, desktop, light/dark, and no-JavaScript rendering; avoid excessive repeated cards.

## Adaptable source example

Bracketed values and example.com links are placeholders, **not evidence**.

```markdown
---
title: "[Project]: from [before] to [after]"
description: "[The documented change over the stated period.]"
content_width: wide
---

[One or two sentences explaining the project's transformation.]

**Scope:** [Which records and dates are included and excluded.]
Dates mean [source reported / deployed / released].

::: md-grid rail {data-columns="2" data-gap="lg"}

:::: md-card

::::: md-icon {data-size="lg"}
:lucide-rocket:
:::::

### [Recent period](#recent-period)

[One sentence defining the new capability.]

::::

:::: md-card

::::: md-icon {data-size="lg"}
:lucide-blocks:
:::::

### [Earlier period](#earlier-period)

[One sentence defining the foundation.]

::::

:::

## Recent period

[Explain the significance of this period in unboxed prose.]

::: road editorial story

### :lucide-rocket: [Verified outcome]
*[Reported period]*

[What changed, why it mattered, and any uncertainty.]

[Evidence](https://example.com/replace-with-real-source)

:::

## Earlier period

[Describe the earlier constraints and the foundation established.]

::: road editorial story

### :lucide-blocks: [Foundational outcome]
*[Reported period]*

[Evidence-grounded result and its limitations.]

[Evidence](https://example.com/replace-with-real-source)

:::

## Still open in the record

The following should not be presented as completed results.

::: md-grid pending {data-columns="2" data-gap="lg"}

:::: md-card

### :lucide-list-todo: [Unfinished work]

[Open status and next validation needed.]

[Open issue](https://example.com/replace-with-open-issue)

::::

:::: md-card

### :lucide-triangle-alert: [Incident or missing evidence]

[Observed incident or verification gap without an inferred outcome.]

[Source record](https://example.com/replace-with-incident)

::::

:::

## Reading the record

[What these sources do not establish, and what evidence is missing.]
```

## Focused evaluation prompts

1. When given a multi-period project history with source links, select this example over the generic workflow illustration.
2. When a period includes an open issue and a merged result, separate them explicitly.
3. When comparative measurements are unavailable, omit the chart rather than invent numbers.
4. Check reading rhythm: prominent H2s outside cards, distinct overview and narrative roles, generous prose spacing.
5. Verify output is readable with scripting disabled and at narrow viewport widths.
