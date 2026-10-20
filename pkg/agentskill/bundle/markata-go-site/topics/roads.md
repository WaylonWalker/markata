# Roads

Use this topic when a user wants a timeline, roadmap, step-by-step path, release
history, course outline, recurring process, or branching history in a post.

Roads are built in to the default theme. You need no template, config, plugin,
or frontmatter. Do not create a custom template or a JavaScript file for this.

## Authoring

Write a `road` container with a list of links. The text after each link is a
small label, such as a date or step number.

```markdown
::: road timeline

- [Kickoff](/notes/kickoff/) Jan 2026
- [[public-beta]] Apr 2026
- [Launch](/notes/launch/) Jul 2026
:::
```

Rules:

- Close the container with a line that contains only colons (`:::`).
- Only list items that contain a link become milestones.
- Wikilinks are good milestones. Their hover cards show the target's
  description.
- A post can have many roads, in any mix of styles.

## Choosing a style

| Request | Style |
|---------|-------|
| Journey, story arc, reading path | `winding` (default) |
| Dated events, history, changelog | `timeline` |
| Three to six process steps | `steps` |
| Long ordered series or course | `snake` |
| Repeating process | `cycle` |
| Loosely related ideas | `constellation` |
| Branches and merges | `git` |
| Entries that need paragraphs, lists, or images | `editorial` |

## Editorial

Write a heading per entry. Put any Markdown under each heading. Use a heading
level below the section heading, so the table of contents stays tidy.

```markdown
## Project history

::: road editorial

### [Discovery](/projects/redesign/discovery/) March

We interviewed twelve readers.

- [x] interviews
- [ ] analytics review

### [Launch](/projects/redesign/launch/) August
:::
```

A road whose first child is a heading always uses the editorial style.

## Git

List commits in order. After each link, write the branch name:

- `feat/x`: a commit on `feat/x`. A new name starts a branch from the trunk,
  which is the first commit's branch.
- `fix/y from feat/x`: starts `fix/y` from the latest commit on `feat/x`.
- `main <- feat/x`: merges `feat/x` into `main`.

A merge must come after at least one commit on both branches.

## Options

- `::: road {data-road="timeline" data-road-spacing="80"}` sets the vertical
  space for each milestone, in pixels.
- Set the `--road-accent` and `--road-muted` CSS custom properties on `.road` in
  the site CSS to change colors.
- A raw `<div class="road" data-road="STYLE">` also works. Put blank lines
  around the list inside it.

## Verify

Run `markata-go build --fast` and open the page. If a road shows as a plain
list, check the closing `:::` line. Then remove `output/` and `.markata-cache/`
and rebuild.
