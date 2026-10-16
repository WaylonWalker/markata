# Deferred alternate feed sidebar windows

Posts must remain readable and navigable without JavaScript. The current
feed's sidebar is server-rendered, and its metadata retains the initial post
window. When a post appears in multiple configured public feeds, alternate
windows may be omitted from `sidebar_feeds_json` **only** if their
`chunkBase` and `chunkSize` refer to validated emitted public sidecars.

On selection, the client fetches the sidecar chunk containing `windowStart`,
renders after the request succeeds and keeps the previous list on failures.
The UI must permit retries and ignore stale requests following a new document
or feed selection. Alternate feeds without validated sidecars retain normal
embedded windows.

The current feed, prev/next metadata, feed picker, rotation order and URL
parameters remain unaffected; browser navigation without JavaScript still
works on the current feed.
