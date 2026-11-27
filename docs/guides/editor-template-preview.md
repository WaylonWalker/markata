# Template-aware draft preview

The Go editor renders unsaved Markdown and frontmatter. The production HTML
template selector resolves per-format overrides, presets, explicit templates,
layout rules and global defaults. Shot/photo and link metadata are displayed
even with an empty Markdown body. Local frontmatter images use the workspace
image route rather than synthetic Markdown source lines.

For a site-specific layout, add a companion template in
`<templates_dir>/editor/<resolved-html-template>`. For example, a post with
`template: project` selects `project.html`; its companion is
`templates/editor/project.html`:

```html
<header class="project-header">
  <h1>{{ post.title }}</h1>
  <p>Status: {{ post.Extra.status }}</p>
</header>
```

The companion replaces the basic metadata prefix, not the Markdown body. The
editor renders the body once beneath it and keeps cursor and scroll mapping.
Use the same partials in your published and editor templates where appropriate.
Includes are restricted to the configured template directory, including symlinks.
The companion receives the production post projection (including custom
frontmatter and `post.Extra`), template filters, and public site title, URL,
description, author and theme. Build-only functions and arbitrary configuration
extras are not available. Do not add scripts or full-page navigation to a
companion. Custom full-page layouts without a companion fall back to the basic
draft preview with a notice rather than claiming published-page fidelity.

Preview does not save, publish or build the site. Changes to unsaved metadata
are reflected on each Go render; site-dependent frontmatter does not use the
plain-text WASM fast path.
