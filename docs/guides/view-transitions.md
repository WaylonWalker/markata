---
title: "View Transitions API"
description: "Smooth, animated page transitions for improved user experience and visual continuity"
date: 2026-02-02
published: true
tags:
  - view-transitions
  - animations
  - javascript
  - performance
---


Smooth, animated transitions between pages and views using the modern View Transitions API. This implementation provides a polished user experience with automatic transitions for all internal navigation.

## Quick Start

View Transitions are **enabled by default**. No configuration required!

To customize behavior, add to your `markata.toml`:

```toml
[view_transitions]
enabled = true
debug = false
```

See [[view-transitions-config|Configuration Reference]] for all options.

## What It Does

The View Transitions API automatically adds smooth animations when users navigate between pages:

- **Wikilinks** (`[[slug]]`) - Smooth transitions between related content
- **Card clicks** - Polished navigation from feeds to full posts
- **Navigation links** - Stable header/footer during transitions
- **Post navigation** - Smooth prev/next post transitions
- **Breadcrumbs** - Contextual navigation animations
- **All internal links** - Any `<a href>` to same-origin pages

### Smart Exclusions

These link types are automatically skipped (use normal navigation):

- External links (`target="_blank"`)
- Download links
- TOC/anchor links (use smooth scroll instead)
- HTMX links (handle their own updates)
- Links with `data-no-transition` attribute

## How It Works

### 1. Link Interception

A global event listener intercepts clicks on internal links:

```javascript
// User clicks an internal link
document.addEventListener('click', (e) => {
  const link = e.target.closest('a');
  if (shouldTransition(link)) {
    e.preventDefault();
    // Start view transition...
  }
});
```

### 2. Content Fetching

The new page is fetched and parsed before the transition starts so the current view stays responsive while network work happens:

```javascript
const response = await fetch(url);
const html = await response.text();
const newDoc = new DOMParser().parseFromString(html, 'text/html');
```

### 3. Transition Animation

The browser animates between old and new content:

```javascript
document.startViewTransition(() => {
  // Update only the changing layout regions.
  document.querySelector('#view-transition-page').replaceWith(
    newDoc.querySelector('#view-transition-page').cloneNode(true),
  );
  document.title = newDoc.title;
  history.pushState(null, '', url);
});
```

### 4. Script Re-initialization

After transition, page scripts are re-initialized:

```javascript
// Dispatch event for scripts to listen to
window.dispatchEvent(new CustomEvent('view-transition-complete'));

// Scripts automatically re-initialize
if (window.initTooltips) window.initTooltips();
if (window.initScrollSpy) window.initScrollSpy();
```

## Browser Support

- ✅ **Chrome/Edge 111+** - Full support
- ✅ **Safari 18+** (macOS 15+, iOS 18+) - Full support
- ✅ **Opera 97+** - Full support
- ⚠️ **Firefox** - In development

**Graceful degradation:** On unsupported browsers, navigation works normally without transitions.

## Default Animation

The default transition is a **fade + slide-up** effect:

- **Old content**: Fades out (250ms)
- **New content**: Fades in + slides up from 20px below (300ms)
- **Navigation**: Stays stable (200ms minimal animation)

### CSS Implementation

```css
/* Assign transition names */
.post-content,
.posts-list {
  view-transition-name: main-content;
}

/* Animate old content out */
::view-transition-old(main-content) {
  animation: fade-out 0.25s ease-out;
}

/* Animate new content in */
::view-transition-new(main-content) {
  animation: fade-in 0.3s ease-in, slide-up 0.3s ease-out;
}

@keyframes fade-out {
  to { opacity: 0; }
}

@keyframes fade-in {
  from { opacity: 0; }
}

@keyframes slide-up {
  from { transform: translateY(20px); }
}
```

## Customization

See [[view-transitions-config|Configuration Reference]] for complete customization options.

### Quick Examples

**Enable debug logging:**

```toml
[view_transitions]
debug = true
```

**Skip transitions for specific classes:**

```toml
[view_transitions]
skip_classes = ["instant-nav", "no-animation"]
```

**Disable globally:**

```toml
[view_transitions]
enabled = false
```

**Per-link opt-out:**

```html
<a href="/page/" data-no-transition>Skip transition</a>
```

## Custom CSS Animations

Override the default animations in your custom CSS:

### Different Animation Style

```css
/* Slide from right instead of bottom */
::view-transition-new(main-content) {
  animation: fade-in 0.3s, slide-from-right 0.3s;
}

@keyframes slide-from-right {
  from {
    opacity: 0;
    transform: translateX(30px);
  }
}
```

### Faster Transitions

```css
::view-transition-old(main-content),
::view-transition-new(main-content) {
  animation-duration: 0.2s;
}
```

### Per-Element Transitions

```css
/* Different animation for post titles */
.post-title {
  view-transition-name: post-title;
}

::view-transition-new(post-title) {
  animation: fade-in 0.3s, scale-up 0.3s;
}

@keyframes scale-up {
  from { transform: scale(0.95); }
}
```

### Individual Card Morphing

Give each card a unique transition name:

```html
<article class="card" style="view-transition-name: card-{{ post.slug }};">
  ...
</article>
```

Now cards will morph smoothly when navigating!

## Performance

- **Bundle size**: ~8KB unminified JavaScript
- **Overhead**: Near-zero (exits early on unsupported browsers)
- **Animation**: GPU-accelerated
- **Network**: Uses standard `fetch()` API
- **Memory**: One active navigation at a time with a small bounded prefetch cache

### Performance Tips

1. **Keep pages small** - Large HTML takes longer to parse
2. **Optimize images** - Images in new content delay transition
3. **Minimize inline scripts** - Scripts need re-initialization
4. **Use caching** - Browser caches reduce fetch time
5. **Swap smaller regions** - Replacing the page wrapper is usually cheaper than replacing the entire `body`
6. **Avoid broad named snapshots** - Prefer shared element names on small cards, titles, or headers rather than whole articles or sidebars

## Accessibility

View Transitions respect user preferences:

### Reduced Motion

Users who prefer reduced motion automatically get instant transitions:

```css
@media (prefers-reduced-motion: reduce) {
  ::view-transition-group(*),
  ::view-transition-old(*),
  ::view-transition-new(*) {
    animation: none !important;
  }
}
```

### Screen Readers

- Page title updates immediately
- Main content landmark is preserved
- Focus management handled automatically
- No ARIA changes needed

## Troubleshooting

### Transitions Not Working

1. **Check browser support** - Only Chrome 111+, Safari 18+, Opera 97+
2. **Check console** - Look for errors or warnings
3. **Enable debug mode** - Set `debug = true` in config
4. **Verify enabled** - Check `window.VIEW_TRANSITIONS_CONFIG.enabled`

### Specific Links Not Transitioning

1. **External link?** - Links with `target="_blank"` are skipped
2. **HTMX link?** - Links with `hx-get` are skipped
3. **TOC link?** - Anchor links use smooth scroll
4. **Custom skip rule?** - Check `skip_classes` and `skip_selectors`
5. **Non-HTML file?** - URLs like `.md`, `.txt`, `.xml`, `.json` use native browser navigation

Enable debug mode to see why links are skipped:

```toml
[view_transitions]
debug = true
```

Then check browser console when clicking links.

### Scripts Not Working After Transition

Make sure your custom scripts listen for the re-initialization event:

```javascript
function myInit() {
  // Your initialization code
}

// Run on initial page load
myInit();

// Re-run after view transitions
window.addEventListener('view-transition-complete', myInit);
```

markata-go dispatches this event after the browser transition finishes. That keeps expensive script hydration out of the critical transition snapshot and makes animations less likely to stutter.

### Janky Transitions

If transitions stutter, enable debug mode and inspect the timing payload logged to the console:

```toml
[view_transitions]
debug = true
```

High `fetchMs` or `parseMs` points to network or page size. High `swapMs` points to expensive DOM/head changes. High `reinitMs` points to custom scripts doing too much work after navigation.

### Transition Too Fast/Slow

Adjust animation duration in CSS:

```css
::view-transition-old(main-content),
::view-transition-new(main-content) {
  animation-duration: 0.5s; /* Slower */
}
```

## Advanced Usage

### Transition Types

Add custom data attributes to control transitions per-route:

```javascript
// Future enhancement - not yet implemented
link.dataset.transitionType = 'slide';
```

### Prefetching

markata-go prefetches transition-managed links on hover/focus and can eagerly warm nearby pagination targets so `[` and `]` navigation feels more immediate. Prefetching is debounced and bounded to avoid excessive fetch/parse work when a user moves across a dense sidebar or feed.

Prefetched documents are reused for five seconds after they finish loading. Expired entries are evicted before the bounded cache admits new work; navigating to an older entry fetches it again with HTTP cache revalidation, so a published update can replace an older in-memory page. HTTP caches can use retained ETag and Last-Modified validators; servers that mark responses `no-store` may require a full response. If revalidation fails, navigation falls back to the ordinary static page request.

```javascript
document.addEventListener('mouseover', (event) => {
  const link = event.target.closest('a');
  if (link) {
    window.prefetchViewTransitionUrl(link.href);
  }
});
```

### Loading States

Show loading indicator during fetch:

```javascript
// Future enhancement - not yet implemented
document.addEventListener('htmx:beforeRequest', () => {
  showLoadingSpinner();
});
```

## Examples

### Simple Fade

```css
::view-transition-old(main-content) {
  animation: fade-out 0.3s;
}

::view-transition-new(main-content) {
  animation: fade-in 0.3s;
}
```

### Scale Transition

```css
::view-transition-old(main-content) {
  animation: scale-down 0.3s;
}

::view-transition-new(main-content) {
  animation: scale-up 0.3s;
}

@keyframes scale-down {
  to { transform: scale(0.95); opacity: 0; }
}

@keyframes scale-up {
  from { transform: scale(0.95); opacity: 0; }
}
```

### Slide Transition

```css
::view-transition-old(main-content) {
  animation: slide-out-left 0.3s;
}

::view-transition-new(main-content) {
  animation: slide-in-right 0.3s;
}

@keyframes slide-out-left {
  to { transform: translateX(-100%); opacity: 0; }
}

@keyframes slide-in-right {
  from { transform: translateX(100%); opacity: 0; }
}
```

## Related Documentation

- [[view-transitions-config|Configuration Reference]] - All configuration options
- [[performance|Performance]] - Performance optimization guide
- [[keyboard-navigation|Keyboard Navigation]] - Accessibility features
- [[themes|Themes]] - Customizing appearance

## Resources

- [MDN: View Transitions API](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API)
- [Chrome Developers: View Transitions](https://developer.chrome.com/docs/web-platform/view-transitions/)
- [View Transitions Demos](https://view-transitions.chrome.dev/)

## Implementation Files

- **JavaScript**: `pkg/themes/default/static/js/view-transitions.js`
- **CSS**: `pkg/themes/default/static/css/components.css`
- **Template**: `pkg/themes/default/templates/base.html`
