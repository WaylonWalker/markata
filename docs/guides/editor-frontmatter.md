---
title: Editing page properties
slug: /docs/guides/editor-frontmatter/
published: true
tags: [documentation, editor, frontmatter]
---


In the Markata workbench, **Properties** (beside the current file name) opens the
visual frontmatter editor. You can also use the command palette
(**Page: Edit Properties**) or **Ctrl/Cmd+Alt+Shift+M**.

Properties start **closed** and the underlying YAML header is folded away, so
writing feels like editing a document rather than configuring a build.

## Common editing tasks

- **Title, description, slug, date, tags:** edit the corresponding form field.
- **Published and Draft:** separate switches, since these have distinct meaning.
  \`draft: true\` excludes a page from output; \`published: false\` excludes it
  from public listings but can leave a directly accessible shadow page.
  The form warns when both \`published\` and \`draft\` are true.
- **Featured image, authors, layout, templates, series, canonical URL:** add these
  from the property picker, or edit them when already present.
- **Other metadata:** enter a custom property name and select its value type.
  Markata permits plugin/site-specific extra properties, so the form does not
  reject an unrecognized but safe YAML property name.
- **Nested maps and lists:** expand the object/list editor, add or remove
  entries, and set the type of nested values.

**Edit YAML** reveals the original source header in the CodeMirror editor
without creating a second copy of frontmatter. Use this for YAML aliases,
anchors, merge keys, unusual custom keys, or when you prefer the full syntax.
**Hide YAML** folds the header again. Both modes modify the same Markdown file.

## What gets written

The form patches the single changed top-level field using Markata's Go YAML
parser. Other frontmatter fields keep their existing syntax, comments,
order, and quotation style. The Markdown body is not re-rendered into source or
rewritten. An edited field may be normalized to standard YAML (including its
nested value), so use raw YAML if you need exact formatting within that field.

The editor sends changes through the **same autosave, preview, diagnostics and
Git publish-review** flow as normal typing. Opening the form does not publish
anything, and no new frontmatter defaults are silently added. A missing
property stays missing until you explicitly add it.

Malformed YAML is never silently replaced; fix it in source mode before using
the form. Workbench frontmatter requests are restricted to the local
same-origin workspace, have source-size limits, and do not write files directly.

See the [frontmatter field guide](/docs/guides/frontmatter/) for the detailed
Markata semantics and accepted date/metadata formats.
