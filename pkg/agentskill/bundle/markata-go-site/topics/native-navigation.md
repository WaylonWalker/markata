# Native navigation performance comparison

The site can opt into full-document navigation with
`[markata-go.view_transitions] native = true`, leaving the default SPA
behavior unchanged. In this lane, the custom DOM-swap script is absent;
CSS `@view-transition` and conservative Speculation Rules enhance normal
links only in browsers that support them.

Before changing a site's default, test LCP, CLS, transfer, keyboard shortcuts,
search, theme persistence, back/forward history and reduced motion in major
browsers. This is an experiment, not a guarantee of Lighthouse 100.
