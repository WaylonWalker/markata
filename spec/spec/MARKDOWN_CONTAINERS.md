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

## Conformance tests

Implementations SHOULD test at minimum:

1. a basic three-colon container
2. nested three/four-colon containers with matching closing markers
3. content after an inner close remains inside the outer container
4. named and attributed `:::` lines are treated as openers
5. bare colon-only markers are the documented closing form
