# Mobile Markata Editor: PWA and inbound sharing

The Go workbench now offers an installation manifest, icons, a service worker share target and a **device-local share inbox**. This is the first implementation slice, not a hosted publishing service.

## Android (Chrome)

The editor must be accessible over HTTPS (or localhost for same-device testing). Open Chrome and choose **Install app** (not just a browser shortcut). From another app, share a link, text or PNG/JPEG/GIF/WebP photo to **Markata**. The installed PWA stores the item locally in IndexedDB and opens the editor.

Choose **Create draft** in the inbox. The existing Link, Photo or Note starter is prefilled. Review the title, metadata, destination, preview and image; choose **Create file** only when ready. Shared content gets published: false frontmatter. It does not publish automatically. Saved inbox items survive a reload until discarded or successfully created.

Sharing several photos queues one unpublished draft per photo; none are silently discarded. The limit is six photos per share, 15 MiB each and 40 MiB total. If storage is unavailable, the editor displays an error rather than claiming the share was retained.

## Quick actions, camera capture and updates

On Android, the installed Markata icon offers shortcuts for **New Note**, **New Photo**, **New Link** and **Share Inbox**. Quick-create shortcuts open the existing site templates; they still require a reachable, authorized workspace. They never publish on their own.

**Take photo** in the editor header opens the camera (or an image picker on devices without capture support) and saves the selected JPEG/PNG/GIF/WebP file into the same IndexedDB share inbox before a remote workspace is involved. This allows camera capture while offline. A cold offline launch falls back to a small, public-only offline capture screen for photos, links and notes. The photo waits for an explicit **Create draft** action after reopening the online editor; it does not become a server-side post until the workspace is available. The post dialog also offers **Take photo**.

When pending shares exist, **Protect storage** asks the browser to retain local storage. The browser may decline; no web app can guarantee immunity to private browsing resets or user-initiated data clearing. Shared photos remain private device/browser-profile data, **not encrypted**, and are not currently isolated between signed-in accounts. Do not enable this shared inbox in a multi-user hosted deployment before account-bound locking/encryption is implemented.

Service-worker updates do not forcibly take over open editing sessions. When an update is waiting, the header offers **Update app** with an explicit save-and-reload confirmation; avoid updating while edits are still saving.

## Security and deployment

- The share service worker handles POST /share-target and network-first GET / navigations. It caches **only** the public static /editor/offline-capture.html fallback for cold offline launches; authenticated pages, workspace documents, session data and API responses are never cached.
- The existing /workspace API is **restricted to loopback**. This patch does not bypass or loosen that requirement. Creating and publishing posts from remote phones requires the authenticated hosted-editor, per-site permissions and sync infrastructure tracked under #1949, #1980 and #1981.
- The app can receive shared content and capture new links, notes and camera photos offline when its installed service worker and browser storage are available. The offline capture screen is **not** an offline editor or preview; creating and publishing the post still requires a reachable and authorized backend. Avoid exposing an unauthenticated local editor to the public Internet.
- iOS Safari can install home-screen web apps but does not currently register PWAs as native iOS Share Sheet targets. iPhone users can open the editor and paste a URL or choose a file. Native inbound Share Sheet support on iOS requires a compatible platform API or native wrapper.
- The queue lives on the device and browser profile where the item was shared; it does not sync across devices. A browser reset or storage eviction can remove it. HEIC photos are not currently accepted by the existing media pipeline.

## Verification

Run:

    node --test pkg/workbench/web/share-target-utils.test.cjs
    go test ./pkg/workbench

Then test Android home-screen shortcuts (including offline photo capture), storage protection refusal, manual update activation, and Chrome Android real-device installation and OS Share Sheet handling for links, text, photo, multi-photo batches, offline receipt, reload, discard, and creation through an authorized editor. Before exposing a hosted build, add end-to-end authenticated remote upload tests and confirm the browser service worker cannot bypass authorization.
