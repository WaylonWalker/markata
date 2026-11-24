---
title: "Serve Control Center"
description: "Inspect live builds, warnings, pages, and logs in the terminal or local web dashboard"
date: 2026-09-26
published: true
tags:
  - documentation
  - development
  - cli
---


Start a local development session:

```bash
markata-go serve
```

In a terminal, `serve` opens a dashboard showing server status, recent jobs,
warnings, errors, pages, and logs. The site remains at
`http://localhost:8000`. File changes create rebuild jobs. Each job keeps its
steps and diagnostics, so a failed rebuild remains inspectable after the next
successful one.

For line-oriented output, use `markata-go serve --no-tui`. Markata also chooses
line-oriented output when input or output is redirected, when `--no-input` is
set, or when the terminal is unavailable. This mode works in scripts and CI.
While the interactive dashboard is open, build and plugin messages appear in
its Logs and Problems views. Serve keeps terminal ownership so rebuild output
does not overwrite the dashboard, including while the terminal is resized.

## Local web dashboard

```bash
markata-go serve --admin
# Site:  http://localhost:8000/
# Admin: http://localhost:8000/_markata/
```

`markata-go serve` enables this dashboard by default on a loopback listener.
`markata-go serve --admin` remains supported. Use `--no-admin` to disable the
dashboard. A non-loopback bind, such as `--bind 0.0.0.0`, leaves it disabled;
requesting `--admin` on such a bind fails. The dashboard URL appears in plain
startup output and as an Admin shortcut in the TUI.

`markata-go admin` starts the same local serve session with the web dashboard
and line-oriented terminal output. Both commands show the same jobs,
warnings, errors, and pages that the terminal dashboard sees. The local admin
dashboard only binds to loopback. For remote, release-oriented operations use
the protected [Builder Admin deployment guide](/docs/guides/deployment/builder-admin/).

The web dashboard supports manual builds, job reruns, and rebuilding a selected
page. Jobs keep logs and lifecycle steps. Page details show source diagnostics
and suggested fixes when Markata has enough context. You can open a generated
page from its page entry when it has a URL. The Feeds section uses the same
feed and post metadata as the build, and its post links open the matching page
details.

### Browser navigation and appearance

The browser URL records the current section, selected item, and filter in its
fragment. You can copy a link to a job, problem, page, feed, or filtered log
view. Browser Back and Forward restore earlier selections and sections.

| Key | Action |
| --- | --- |
| `j` / `k`, arrows | Move through the current list |
| `Enter` | Move focus to the selected item's details |
| `Esc` | Clear the current selection or close keyboard help |
| `/` | Focus the current section's filter |
| `g` / `G` | Select the first or last item |
| `w` / `e` | Show warnings or errors |
| `p` / `f` / `l` | Open pages, feeds, or logs |
| `r` | Rerun the selected completed job while Jobs is active |
| `?` | Show keyboard help |

Shortcuts pause while a text field or editable element has focus. At narrow
window widths, the list and details use a single-column layout with a fixed
section bar. Focus indicators remain visible for keyboard use.

The browser colors follow the configured Markata palette through semantic roles
for surfaces, text, links, focus, status, code, and buttons. The shared semantic
token stylesheet is also used by Builder Admin, which keeps its separate
production workflows and page structure. A shared browser shell is tracked in
[#1300](https://github.com/WaylonWalker/markata-go/issues/1300); the local fix
preview and source mutation actions remain available only in local Serve.

## Review and apply fixes

The local Problems view can preview source fixes before applying them. Plans
show their category, safety class, source location, and before/after text.
Use **Review all safe** to review every currently safe fix, a category button
to review one group, or **Select** on individual problems followed by
**Review selected**. The review dialog lists every proposed edit before you
approve the batch. It shows a full-file diff for each source file using bundled
assets, so review works without internet access. If rich rendering fails, the
dialog shows a text diff of the same previewed edits.
`SAFE` fixes are deterministic edits that can be selected together; `REVIEW`
fixes need closer inspection; `MANUAL` diagnostics provide guidance without an
automatic edit. Markata never applies a fix until you review a preview and
choose to apply it. The server authorizes only the exact files, digests, and
selections in that preview, and consumes that authorization after one apply
attempt. The retired single-file preview/apply endpoints are not exposed; all
source mutation uses the grouped batch preview/apply contract.

Each source file is checked against the digest captured by the preview. If the
file changes before apply, Markata skips that file and reports it. In a mixed
batch, still-valid files may be applied while stale files remain untouched.
Markata requests one rebuild after applying any fixes. The build job then
refreshes the current diagnostics projection. The apply response means the
rebuild was queued; the Problems view updates when that job completes.
Ambiguous dates and duplicate frontmatter keys do not receive automatic safe
edits.

## Terminal controls

| Key | Action |
| --- | --- |
| `j` / `k`, arrows | Move selection or scroll |
| `Enter` | Inspect selected job, page, or diagnostic |
| `Esc` | Return to the previous view |
| `t` | Trigger a build |
| `r` | Rerun the selected job |
| `w` / `e` | Open warnings or errors |
| `l` | Open logs |
| `p` / `f` | Open pages or feeds |
| `o` / `v` | Open selected source or local preview |
| `a` | Open Admin in the browser when available |
| `/` | Search or filter the current view |
| `g` / `G` | Jump to first or last item |
| `PgUp` / `PgDn` | Scroll by a page |
| `?` | Show help |
| `q` | Quit |

When you scroll away from the live log tail, new lines do not move your view.
The dashboard shows a count of new lines until you return to the bottom.
Resize the terminal at any time; narrow windows use a single-column layout.

## Reading diagnostics

A diagnostic has a stable code, source path and line when known, the job that
found it, and a suggested fix. The warnings and errors inboxes retain entries
for the session history. The Problems view uses current diagnostics from the
latest page inventory, newest completed job, and any running job; a successful
rebuild removes resolved findings from that view. Their original job records remain available. Selecting
a page shows its current status and diagnostics; selecting a job shows the
diagnostics from that run. Some plugin messages do
not include a source location. Those remain in the job and session inbox so
they cannot disappear into log scrollback.
