---
title: "Start the Markata workspace"
description: "Run the browser editor and live-reloading site builder together."
date: 2026-11-22
published: true
tags: [documentation, editor, cli]
---

# Start your writing workspace

After creating a site with `markata config init` or working in an existing Markata
site directory, launch the editor and builder in a single terminal:

```sh
markata up
```

The command starts **two independently managed Go processes**:

- **Editor** — `http://127.0.0.1:8766/` (Markdown writing and preview)
- **Builder** — `http://127.0.0.1:8000/` (static site, live preview, file watching)

The URLs printed at startup are **expected addresses**, not confirmation of
HTTP readiness. Wait for each component's own ready message. The editor can
still be used if the builder fails; a failed builder is reported in the
terminal, never represented as a successful build. Draft saving does **not**
publish the public website.

## Change ports or run one component

```sh
markata up --editor-port 8767 --preview-port 8001
markata up --editor-only
markata up --builder-only
```

The editor and builder bind to **loopback only**. A port conflict stops startup
rather than exposing a service on another interface. The parent forwards
Ctrl+C to its child processes, gives them an opportunity to shut down
gracefully, and cleans them up before exiting.

For separate processes and advanced options, you can still run
`markata editor` and `markata serve` independently.

## Remote and background access

This initial command only manages a **foreground** workspace. It does not
create a background service, pairing token, or public editor endpoint.
Do **not** open the editor port directly to the internet or expose the
builder's local administration endpoints.

Background lifecycle, `markata status` / `down` / `doctor`, private
Tailscale pairing, and mobile PWA onboarding are tracked in
[the onboarding epic](https://github.com/WaylonWalker/markata/issues/2087).
