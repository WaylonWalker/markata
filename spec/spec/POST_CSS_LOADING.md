# Conditional feed-card CSS on stock reading posts

The stock `post.html` template may omit `cards.css` from initial stylesheet
requests only if the authored post does not include markup or shortcodes for
cards or feeds. Default to retaining it for unknown or custom templates.
Full feed/listing pages retain the stylesheet.

The generator, not a browser bootstrap, must select this behavior.
In the JavaScript navigation lane, conditional-css.js still loads
`cards.css` if a card-bearing page is entered via DOM swap. Native
full-document navigation gets the destination stylesheet naturally.
The initial document still needs a conditional content scan. Feed cards can
contain generated embeds, code blocks, and admonitions whose requirements are
not visible to the post-level template markers. Load card and feed layout CSS
immediately. Defer code, syntax-highlight, and admonition CSS only when all
matching elements are inside feed cards, until those elements approach the
viewport. Rich content outside cards loads its CSS immediately. Keep an
immediate fallback when IntersectionObserver is unavailable. Search stays
lazy until user intent, then reinitializes after SPA swaps.

Do not hide or cap rendered cards. Test custom templates, embedded feeds,
post cards, and normal article text. This is a narrow resource removal, not
general per-page critical CSS extraction.
