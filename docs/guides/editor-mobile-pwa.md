# Mobile Markata Editor: PWA and inbound sharing

The Go workbench now offers an installation manifest, icons, a service worker share target and a **device-local share inbox**. This is the first implementation slice, not a hosted publishing service.

## Android (Chrome)

The editor must be accessible over HTTPS (or localhost for same-device testing). Open Chrome and choose **Install app** (not just a browser shortcut). From another app, share a link, text or PNG/JPEG/GIF/WebP photo to **Markata**. The installed PWA stores the item locally in IndexedDB and opens the editor.

Choose **Create draft** in the inbox. The existing Link, Photo or Note starter is prefilled. Review the title, metadata, destination, preview and image; choose **Create file** only when ready. Shared content gets published: false frontmatter. It does not publish automatically. Saved inbox items survive a reload until discarded or successfully created.

Sharing several photos queues one unpublished draft per photo; none are silently discarded. The limit is six photos per share, 15 MiB each and 40 MiB total. If storage is unavailable, the editor displays an error rather than claiming the share was retained.

## Security and deployment

- The share service worker intercepts only POST /share-target, persists local share payloads and redirects to the editor. It does not cache authenticated pages, workspace documents, session data or API responses.
- The existing /workspace API is **restricted to loopback**. This patch does not bypass or loosen that requirement. Creating and publishing posts from remote phones requires the authenticated hosted-editor, per-site permissions and sync infrastructure tracked under #1949, #1980 and #1981.
- The app can receive shared content offline when its service worker and browser storage are available. Creating the post still requires a reachable and authorized backend. Avoid exposing an unauthenticated local editor to the public Internet.
- iOS Safari can install home-screen web apps but does not currently register PWAs as native iOS Share Sheet targets. iPhone users can open the editor and paste a URL or choose a file. Native inbound Share Sheet support on iOS requires a compatible platform API or native wrapper.
- The queue lives on the device and browser profile where the item was shared; it does not sync across devices. A browser reset or storage eviction can remove it. HEIC photos are not currently accepted by the existing media pipeline.

## Verification

Run:

    node --test pkg/workbench/web/share-target-utils.test.cjs
    go test ./pkg/workbench

Then test Chrome Android real-device installation and OS Share Sheet handling for links, text, photo, multi-photo batches, offline receipt, reload, discard, and creation through an authorized editor. Before exposing a hosted build, add end-to-end authenticated remote upload tests and confirm the browser service worker cannot bypass authorization.
