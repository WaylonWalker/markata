---
title: "Roads"
description: "Turn a list of links into a scroll-drawn road: timelines, step-by-step guides, roadmaps, cycles, and git histories"
date: 2026-10-20
published: true
slug: /docs/guides/roads/
tags:
  - documentation
  - markdown
  - containers
  - themes
---

# Roads

A road is a path that draws itself as the reader scrolls. Each milestone on the
path is a link from your post. Use roads for roadmaps, release histories,
onboarding paths, project timelines, or any sequence that readers walk through.

Roads are built in to the default theme. You need no template, config, or
plugin. Write a [Markdown container](/docs/guides/markdown-containers/) named
`road`, and put a list of links inside it.

## Quick start

````markdown
::: road

- [Install markata-go](/docs/getting-started/) Step one
- [Configure your site](/docs/guides/configuration/) Step two
- [Pick a theme](/docs/guides/themes/) Step three
- [Publish feeds](/docs/guides/feeds/) Step four
:::
````

That renders as:

::: road

- [Install markata-go](/docs/getting-started/) Step one
- [Configure your site](/docs/guides/configuration/) Step two
- [Pick a theme](/docs/guides/themes/) Step three
- [Publish feeds](/docs/guides/feeds/) Step four
:::

Each list item is one milestone:

- **The link** is the milestone title, and clicking it opens the page.
- **The text after the link** is a small label, such as a date or step number.
  It is also shown in the hover card.

Wikilinks work too. A milestone written as `[[my-post]]` uses that post's
title, and its hover card shows the post's description.

## Set fixed progress

Use fixed progress for a roadmap whose completion does not depend on scrolling.
The normal scroll-driven progress remains the default.

```markdown
::: road timeline {data-road-progress-mode="static" data-road-progress="65"}
- [Design](/design/)
- [Build](/build/)
- [Release](/release/)
:::

::: road steps {data-road-progress-mode="milestones" data-road-current="2"}
- [Design](/design/)
- [Build](/build/)
- [Test](/test/)
- [Release](/release/)
:::
```

`static` accepts a percentage from 0 to 100, with or without a `%` suffix.
`milestones` accepts the number of completed entries. Values are clamped to the
road's range. Invalid values fall back to scroll progress. Fixed progress stays
the same after scroll, resize, and view transitions. The original links remain
available if JavaScript does not run.
At 0, no entries are active. At 100%, or when all milestones are complete, all
entries are active.

## Choose a style

Put the style name after `road`. The default style is `winding`.

| Style | Good for |
|-------|----------|
| `winding` | Journeys, reading paths, story arcs |
| `timeline` | Dated events, changelogs, history |
| `steps` | Short processes of three to six steps |
| `snake` | Long ordered sequences, such as courses or series |
| `cycle` | Processes that repeat, such as sprints or seasons |
| `constellation` | Loosely related ideas or a "see also" map |
| `git` | Branching work, such as release trains or experiments |
| `editorial` | Timelines where each entry has paragraphs, lists, or images |

### Timeline

````markdown
::: road timeline

- [Project announced](/blog/announcement/) Jan 2026
- [Public beta](/blog/beta/) Apr 2026
- [Version 1.0](/blog/v1/) Jul 2026
- [Plugin marketplace](/blog/plugins/) Oct 2026
:::
````

::: road timeline

- [Getting started](/docs/getting-started/) Day 1
- [Markdown features](/docs/guides/markdown/) Day 2
- [Templates](/docs/guides/templates/) Day 3
- [Feeds](/docs/guides/feeds/) Day 4
:::

### Steps

Steps read from left to right. On narrow screens they change to a timeline.

````markdown
::: road steps

- [Draft](/handbook/draft/) write
- [Review](/handbook/review/) edit
- [Publish](/handbook/publish/) ship
- [Share](/handbook/share/) promote
:::
````

::: road steps

- [Install](/docs/getting-started/) setup
- [Configure](/docs/guides/configuration/) config
- [Theme](/docs/guides/themes/) style
- [Deploy](/docs/guides/deployment/) ship
:::

### Snake

The snake style wraps a long list back and forth across rows, like a board
game.

````markdown
::: road snake

- [Variables](/course/variables/) 1
- [Functions](/course/functions/) 2
- [Loops](/course/loops/) 3
- [Structs](/course/structs/) 4
- [Interfaces](/course/interfaces/) 5
- [Errors](/course/errors/) 6
- [Testing](/course/testing/) 7
- [Concurrency](/course/concurrency/) 8
:::
````

### Cycle

````markdown
::: road cycle

- [Plan](/process/plan/) week 1
- [Build](/process/build/) week 2
- [Review](/process/review/) week 3
- [Retro](/process/retro/) week 4
:::
````

::: road cycle

- [Write](/docs/guides/markdown/) draft
- [Build](/docs/reference/cli/) preview
- [Review](/docs/guides/content-diagnostics/) check
- [Publish](/docs/guides/deployment/) ship
:::

### Constellation

The points sit in fixed, scattered positions joined by fine lines. The layout
is the same on every visit.

````markdown
::: road constellation

- [[on-writing-clearly]] essay
- [[notes-on-focus]] note
- [[a-reading-list]] list
- [[why-plain-text]] essay
:::
````

### Git

The git style draws a commit graph with branches and merges. List the commits
in order, and write the branch after each link:

| After the link | Meaning |
|----------------|---------|
| `main` | A commit on `main`. A new branch name starts a branch from the first branch (the trunk). |
| `fix/typo from feat/search` | Starts `fix/typo` from the latest commit on `feat/search`. |
| `main <- feat/search` | Merges `feat/search` into `main`. You can also write `←`. |

````markdown
::: road git

- [Initial commit](/changelog/#init) main
- [Add search index](/changelog/#search) feat/search
- [Update README](/changelog/#readme) main
- [Search results page](/changelog/#results) feat/search
- [Dark mode](/changelog/#dark) feat/dark-mode
- [Merge search](/changelog/#merge-search) main <- feat/search
- [Try a new palette](/changelog/#palette) spike/palette from feat/dark-mode
- [Dark mode polish](/changelog/#polish) feat/dark-mode
- [Keep the palette](/changelog/#keep) feat/dark-mode <- spike/palette
- [Release 1.0](/changelog/#v1) main <- feat/dark-mode
:::
````

The first commit's branch is the trunk. Each branch gets its own lane and
color. A merged branch frees its lane for the next branch.

### Editorial

Use the editorial style when each entry needs more than a title. Write a
heading for each entry, then any Markdown under it: paragraphs, lists,
checklists, quotes, images, or code.

````markdown
::: road editorial

### [Discovery](/projects/redesign/discovery/) March

We interviewed twelve readers and found three common problems.

- [x] interviews
- [x] survey
- [ ] analytics review

### [Prototype](/projects/redesign/prototype/) May

> Smaller pages, bigger type, fewer clicks.

### [Launch](/projects/redesign/launch/) August
:::
````

::: road editorial

#### [Markdown](/docs/guides/markdown/) Write

Write posts in Markdown with frontmatter.

- [x] headings and lists
- [x] admonitions and containers
- [x] wikilinks

#### [Templates](/docs/guides/templates/) Shape

Change the layout with templates, or keep the defaults.

#### [Deployment](/docs/guides/deployment/) Ship

Build static HTML and host it anywhere.
:::

How editorial entries work:

- The text after the heading link (`March`) becomes a small label above the
  title.
- The rail fills in as the reader moves through the entries.
- The headings stay real headings, so they keep their anchor links and appear
  in the table of contents. Use a heading level below your section heading, so
  the table of contents stays tidy. For example, use `###` entries under a `##`
  section.
- A road whose first item is a heading uses the editorial style automatically,
  even when you write `::: road` or `::: road timeline`.

### Editorial stories, icons, and reveal

````markdown
::: road editorial story reveal

### :lucide-key-round: RSA migration and faster job views
*2026 Q3*

A lead paragraph explaining what changed and why it matters.

- Faster job views
- Simpler key rotation
- Clearer audit trails
- Fewer manual steps

### A quieter milestone without an icon
*2026 Q4*

The ordinary rail dot still works.
:::
````

A leading `:lucide-name:` icon becomes a badge on the rail, muted until the
entry is active. An emphasis-only paragraph immediately after a heading becomes
the eyebrow above its title. Both work on all editorial roads, without `story`.
If the heading already has a beside-link label, it takes precedence and the
emphasis paragraph stays in the body. Mixed-content paragraphs are not labels.

Add **`story`** for larger balanced titles, a larger opening paragraph, generous
spacing, and a prose measure of about 64 characters. Direct unordered lists
without task checkboxes become diamond-bullet highlights with a top rule.
Lists of four or more items use two columns at viewport widths of 1100px or
more. Cards, stats, and figures can use the wider entry body. Without `story`,
existing editorial typography stays unchanged.

Add **`reveal`** independently for scroll-triggered entry fades and gentle
upward motion. It works with or without `story`, respects reduced motion, and
does not animate print, anchor targets, or keyboard-focused entries. Entries
stay readable without JavaScript or IntersectionObserver, even if the
enhancement fails: reveal animates arriving content rather than hiding it
until a callback runs.

## More than one road

A post can have any number of roads, in any mix of styles. Each road draws on
its own as it scrolls into view.

## Options

Use the container attribute block for options:

````markdown
::: road {data-road="timeline" data-road-spacing="80"}

- [Alpha](/releases/alpha/) v0.1
- [Beta](/releases/beta/) v0.5
:::
````

| Option | Description |
|--------|-------------|
| `data-road="STYLE"` | Selects the style. This is the same as `::: road STYLE`. |
| `data-road-spacing="N"` | Sets the vertical space for each milestone, in pixels. |

To change the colors, set CSS custom properties in your site CSS:

```css
.road {
  --road-accent: #3fb68b;  /* progress line, active dots, labels */
  --road-muted: #4a4f5a;   /* the path before it is drawn */
}
```

The defaults come from your palette's `--color-primary` and
`--color-border-strong`.

## Using HTML instead of a container

A raw HTML `<div>` also works. Leave a blank line after the opening tag and
before the closing tag, so the list inside is parsed as Markdown:

```html
<div class="road" data-road="timeline">

- [Kickoff](/notes/kickoff/) Jan 2026
- [Launch](/notes/launch/) Jun 2026

</div>
```

## Accessibility and fallbacks

- Without JavaScript, a road shows as the ordinary list or headings that you
  wrote.
- While a road loads, its space is reserved, so the raw list does not flash and
  the page does not jump. If the script fails, the list shows after 2.5 seconds.
- Every milestone stays a real link that works with the keyboard.
- When the reader asks for reduced motion, roads appear fully drawn, with no
  animation.
- Roads keep working when readers move between pages with
  [view transitions](/docs/guides/view-transitions/).

## How roads load

The road CSS and JavaScript load only on pages that contain a road. markata-go
finds roads when it renders the post. CSS purge keeps the road styles in
production builds, with no extra configuration.

## Troubleshooting

**The road shows as a plain list.** Check that the container closes with a
line that contains only colons (`:::`). For a raw `<div>`, check for the blank
lines around the list. After you change shared CSS or templates, do a clean
rebuild.

**Some list items are missing.** Only list items that contain a link become
milestones.

**The git graph joins the wrong branches.** Branch names must match exactly.
A merge (`main <- feat/x`) must come after at least one commit on both
branches.

## See also

- [Markdown containers](/docs/guides/markdown-containers/)
- [Markdown features](/docs/guides/markdown/)
- [Themes](/docs/guides/themes/)
