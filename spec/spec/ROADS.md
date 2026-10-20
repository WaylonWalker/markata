# Roads

A road is a scroll-drawn path through a sequence of linked milestones. Authors
write roads in Markdown. The default theme enhances them in the browser.

## Authoring

A road is any element with the class `road` or the attribute `data-road`. The
normal way to create one is a [Markdown container](MARKDOWN_CONTAINERS.md):

```markdown
::: road timeline

- [Kickoff](/notes/kickoff/) Jan 2026
- [Public beta](/blog/beta/) Apr 2026
- [[launch-retrospective]] Jul 2026
:::
```

A raw `<div class="road" data-road="STYLE">` MUST behave the same as the
container. Markdown inside a raw `<div>` needs blank lines around it.

The style comes from the first match of:

1. the `data-road` attribute, when its value is a known style;
2. a container class that is a known style, such as `::: road snake`;
3. `editorial`, when the road's first child element is a heading;
4. `winding`.

### List roads

All styles except `editorial` use the first list inside the road. Each list
item that contains a link is one point, in document order. List items without
links are ignored.

- The link text is the point's title, and the link is the point's target.
- The remaining text of the list item, without nested lists, is the label.
  Leading separators (`-`, `–`, `—`, `:`, `·`, `|`) are removed.
- Wikilinks are valid links and keep their own hover data.

### Editorial roads

An editorial road splits its children into entries at each heading. The entries
use the level of the first heading. Content before the first heading is
outside all entries.

- The first link in a heading is the entry title.
- The other text in the heading is shown as a label above the title.
- All content until the next heading of the same level is the entry body. Any
  Markdown is allowed in the body.
- Headings MUST stay heading elements and keep their IDs, so anchors and the
  table of contents keep working.

### Attributes

| Attribute | Meaning |
|-----------|---------|
| `data-road="STYLE"` | Selects a style. |
| `data-road-spacing="N"` | Vertical pixels per point for list styles with vertical spacing. |

The `--road-accent` and `--road-muted` CSS custom properties set the colors.
They default to `--color-primary` and `--color-border-strong`.

## Styles

| Style | Layout |
|-------|--------|
| `winding` | A curved path that swings left and right. Cards alternate sides. |
| `timeline` | A straight vertical rail with cards on the right. |
| `steps` | Numbered stages on a horizontal line. It uses the `timeline` layout when each column is narrower than 110px. |
| `snake` | Rows that alternate direction, joined by curves. |
| `cycle` | Points spaced evenly on a closed circle, starting at the top. |
| `constellation` | Points at fixed pseudo-random positions, joined by straight lines. The positions MUST be the same on every page load. |
| `git` | A commit graph with one lane per active branch. |
| `editorial` | A vertical rail beside headed entries with prose bodies. |

### Git labels

In the `git` style the label names the commit's branch:

| Label | Meaning |
|-------|---------|
| `BRANCH` | A commit on `BRANCH`. If the branch is new, it forks from the trunk. |
| `BRANCH from PARENT` | A new branch that forks from the latest commit on `PARENT`. |
| `BRANCH <- OTHER` | A merge of `OTHER` into `BRANCH`. `←` MAY replace `<-`. |

- The trunk is the branch of the first commit. A commit with no label is on the
  trunk.
- A new branch MUST take the lowest free lane, and a merged branch MUST free its
  lane.
- Each branch MUST have a distinct color from a fixed palette. The first color
  is `--road-accent`.

## Behavior

- Without JavaScript, a road MUST remain its original readable Markdown.
- When scripting is enabled, a road that is not built yet (no
  `data-road-ready`) MUST be hidden, and its CSS MUST reserve about the height
  of the built road, so the raw list does not flash and the page does not jump.
  If the road is still not built after 2.5 seconds, CSS MUST show the raw
  content. A built road SHOULD fade in, except with reduced motion.
- A page MAY contain any number of roads, in any mix of styles.
- The drawn progress MUST follow the reader's scroll position. A point or entry
  MUST be marked `active` when the progress reaches it.
- When the reader prefers reduced motion, roads MUST render fully drawn without
  animation.
- Roads MUST lay out again when their width changes.
- Point links MUST stay focusable and navigable. When a link has no hover data
  of its own, it MUST get `data-hover-title` from its title and
  `data-hover-description` from its label, so the default tooltip shows them.
- Every class that the enhancement adds MUST start with `road-`, except the
  shared `active` state class.

## Asset loading

`render_markdown` MUST set `post.Extra["needs_roads"]` when the rendered HTML
has a class attribute with the `road` token, or a `data-road` attribute. Other
class names that contain `road`, such as `roadmap`, MUST NOT set it.

When `needs_roads` is true, `base.html` MUST load `css/roads.css` and
`js/roads.js`. The script MUST be in `<head>` so that view transitions run it.
It MUST be safe to load more than once. It MUST build new roads on
`view-transition-complete` and `pageshow`, and forget roads that are no longer
in the document. `conditional-css.js` MUST load both assets after a view
transition to a page that contains a road.

CSS purge MUST keep `road-*` classes and `data-road` attribute selectors by
default.
