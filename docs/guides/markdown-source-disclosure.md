# Markdown example preview and source

A source disclosure can follow any rendered Markdown example. Markata renders
this as native HTML `<details>`, so it supports mouse, keyboard, and
screen readers **without JavaScript**. No new JS bundle is necessary.

The source is explicitly authored: this feature does not automatically
duplicate Markdown, execute snippets, or copy untrusted HTML into the page.

## Basic example

````markdown
::: md-layout

:::: md-card
### Preview

This is **bold** and this is `inline code`.

::::

:::: md-source {data-label="See how it is written"}

```markdown
### Preview

This is **bold** and this is `inline code`.
```

::::

:::
````

Use a longer backtick fence (four backticks or more) around an example
containing nested triple-backtick code fences, as in other Markdown guides.

The label is optional and defaults to **Show Markdown source**. Any authored
label is escaped before being inserted in the disclosure summary.

## Output and compatibility

- Static semantic HTML only; no JavaScript is required to open the source.
- Source code uses the existing fenced-code highlighter and copy controls
  when those optional enhancements are present.
- Source text is explicitly duplicated by the author, not extracted
  automatically from the preview. This helps avoid hidden dynamically
  generated markup and makes it easy to inspect what a post actually teaches.
- The source disclosure is marked `data-markata-ui="source-preview"`
  so the rich post-copy command can omit its UI and avoid duplication there.
- Search indexes and ordinary feed readers may still include explicitly
  authored source content; use this for deliberate educational examples,
  not every sentence of an article.
