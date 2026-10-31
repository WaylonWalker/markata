# Markdown Container Specification

This document specifies fenced Markdown containers using `:::` delimiters.

## Opening syntax

A container opener MUST:

- begin with at least three colons
- contain non-whitespace content after the colon marker
- use the number of leading colons as the container depth

Examples:

```markdown
::: card
:::: nested-card
::: card {#intro .featured data-kind="note"}
```

Text after the marker MAY provide classes and an attribute block. A named or attributed marker is always an opener, never a closer.

## Closing syntax

A documented closing marker MUST be a bare line containing only three or more colons.

```markdown
:::
::::
```

These are not documented closing forms:

```markdown
::: card
::: {.card}
:::: nested
```

Authors SHOULD use exact, undecorated colon-only lines as closers rather than depending on whitespace or indentation normalization.

## Nesting semantics

Each container records the number of colons used by its opener.

- A colon-only marker with the same depth closes the current container and consumes the marker.
- Authors SHOULD use a longer delimiter for each nested level and close each level with its matching delimiter.
- Content after a matching inner close MUST remain inside the still-open parent container.
- Authors MUST NOT rely on a shorter parent marker to implicitly unwind deeper open containers; nested containers should be closed from the inside out.

Recommended example:

```markdown
::: outer
Outer before.

:::: inner
Inner content.
::::

Outer after.
:::
```

Expected structure:

```html
<div class="outer">
<p>Outer before.</p>
<div class="inner">
<p>Inner content.</p>
</div>
<p>Outer after.</p>
</div>
```

## Attributes

Opening markers MAY contain:

- one or more whitespace-separated classes
- an ID using `#id` inside an attribute block
- extra classes using `.class` inside an attribute block
- key/value attributes such as `data-kind="note"`

Example:

```markdown
:::: panel wide {#example .accent data-state="ready"}
Content
::::
```

The rendered container MUST retain the declared classes, ID, and valid custom attributes.

The default theme gives some container classes extra behavior. For example,
`::: road STYLE` creates a scroll-drawn road, as defined in [ROADS.md](ROADS.md).
Editorial roads accept independent `story` (heavier typography) and `reveal`
(motion-safe scroll entrance) modifiers: `::: road editorial story reveal`.
Leading Lucide icons and emphasis-only label paragraphs enhance all editorial
roads by default.

### Card heading permalinks

In the default theme, heading anchors anywhere inside an article's `md-card`
MUST remain within the heading's normal text flow, wrapping when needed.
They MUST remain available on hover and keyboard focus without causing
horizontal card scrolling. Headings outside cards retain their existing anchors.
Wide code blocks retain their own scrolling, and wide comparison tables MAY
scroll within the card rather than widening the page.

## Story layout primitives

`::: md-section` MUST render a `section` retaining its classes, ID and attributes.
When its first block is an H2, that heading and the immediately following
paragraph (if present) MUST become `.md-section-head`; remaining blocks become
`.md-section-body`. Without a leading H2, all blocks belong to the body and no
empty head is emitted.

A leading plain-text token of 1–16 non-whitespace characters followed by the
literal separator ` · ` and a nonempty title MUST become a
`span.md-section-kicker` inside the H2. Inline formatting in the title MUST
survive. Heading IDs and anchor URLs MUST still derive from the original full
heading text, before splitting. Other headings MUST remain unchanged.

The default theme MUST provide a ruled, generously spaced section, with a
sticky 16rem head column at viewport widths of at least 1100px, a single column
below that, and static heads in print. A section without a head uses one column.

Frontmatter `content_width: wide` opts a default-theme post into `.post--wide`:
a maximum reading frame of `min(70rem, 100% - 2 * --space-10)`, balanced larger
title, opening lead paragraph and roughly 64ch prose within sections. Mobile
uses the available width. Unset or unknown values MUST preserve normal styling.

`::: md-grid rail` uses the existing modifier-class and attribute syntax.
At least 1000px wide, it defaults to four columns, respects `data-columns`
(2, 3 or 4), connects card tops with a two-pixel gradient rail and ring ticks,
and disables hover lift. Ring ticks MUST sit at `top: -6px` within each card.
Rail cards MUST use visible overflow so the ticks remain visible.
An immediately following `md-section` MUST use
`margin-top: clamp(4.5rem, 8vw, 7.5rem)` to separate the overview from the story.
Below 1000px it stacks along a vertical rail with ring ticks. Stacked cards
use a straight left edge and retain the same hover/focus tick feedback.
Heading links (excluding generated permalinks) stretch across the card, with
visible keyboard focus on the card. Emphasis-only paragraphs are muted uppercase
monospace kickers. `::: md-grid pending` gives cards dashed borders, transparent
backgrounds and muted icons. These features MUST work without JavaScript and
respect reduced-motion preferences.

## Sparklines

`::: md-sparkline` with a `data-values` attribute of 2–500 finite numbers
(separated by whitespace or commas) and a nonempty `data-label` MUST render,
as the first child of its `div`, an inline
`svg.md-sparkline-chart` with `viewBox="0 0 100 24"`,
`preserveAspectRatio="none"`, `role="img"` and `focusable="false"`. It MUST
contain a `polyline.md-sparkline-line` with
`vector-effect="non-scaling-stroke"`. Only with the `area` modifier class
(`::: md-sparkline area`) MUST it also contain a `polygon.md-sparkline-area`
before the line, closing the line down to y=24. Points MUST be evenly spaced from x=0 to
x=100. Values MUST map linearly from the series high (y=2) to the series low
(y=22). A flat series MUST draw at y=12. Coordinates MUST be rounded to two
decimals. Finite input values MUST produce finite coordinates even when the
series range exceeds the largest representable float.

The SVG `aria-label` MUST be `<label>: <n> points from <first> to <last>, low
<low>, high <high>`, HTML-escaped. With invalid, missing or unlabeled data, no
SVG is emitted and the container body MUST render unchanged. Other containers
MUST NOT emit a sparkline, even with the same attributes. The default theme
colors the line with `data-tone` like the other reporting components. It MUST
render without JavaScript.

## Conformance tests

Implementations SHOULD test at minimum:

1. a basic three-colon container
2. nested three/four-colon containers with matching closing markers
3. content after an inner close remains inside the outer container
4. named and attributed `:::` lines are treated as openers
5. bare colon-only markers are the documented closing form
