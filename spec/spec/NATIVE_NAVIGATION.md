# Optional browser-native document navigation

The `view_transitions.native` setting is a nullable boolean with default
false. Existing installations retain JavaScript-driven SPA-like view
transitions. When `enabled=true` and `native=true`, the default theme
must omit both custom `view-transitions.js` and `conditional-css.js`
navigation-enhancement scripts, and emit `@view-transition { navigation:
auto; }`. Only same-origin cross-document navigations are supported by
browser-native transitions. Unsupported browsers navigate with normal links.

The native lane may emit conservative Speculation Rules `prefetch` for
sidebar and post-navigation links; it must not prerender, prefetch every
article, or intervene in form/history behavior. Respect reduced motion.
`enabled=false` always disables transitions.

This is an opt-in comparison lane: no claims about production performance
or Lighthouse scores without measurement. Keep regular SPA navigation as
the default until browser compatibility, copy controls, search, theme,
back/forward navigation, anchor linking and mobile metrics are validated.
