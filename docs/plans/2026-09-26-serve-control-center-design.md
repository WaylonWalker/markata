---
title: "Serve Control Center design"
description: "Architecture and migration plan for shared serve jobs, diagnostics, terminal, and web clients"
date: 2026-09-26
published: true
tags:
  - documentation
  - development
---


## Existing Builder Admin

The `markata-go builder-admin` command launches `pkg/builderadmin.Service` from
the same executable as the normal CLI. Helm launches this command in the builder
image. The service has a leader-locked, serialized queue for builds, refresh
tasks, and rollback. It keeps queued and running operations, completed build and
refresh records, trigger identity, timestamps, changed paths, release metadata,
and log file paths. It persists `state.json` and per-job text logs under its
history directory. The UI polls `/api/state`; users can enqueue a build or
refresh, inspect logs, and promote a retained release. Build jobs execute a
child `markata-go build` process before release promotion. Proxy identity,
HTTPS origin validation, and CSRF protect the production UI.

Builder Admin does not currently hold page status, warnings, structured
diagnostics, step records, or structured log events. Its local command binds to
loopback by default, but the protected UI cannot be used without the proxy
contract. The existing release and queue workflow stays in place. The local
serve control view is mounted on the serve listener and consumes the shared
session runtime. This avoids weakening production authentication.

## Shared state

A session runtime stores jobs, steps, logs, diagnostics, pages, and server
status. It sends immutable snapshots to subscribers. Build integration creates
an initial job and a job for each rebuild, with one step per lifecycle stage.
The existing content diagnostics ledger supplies page dispositions and stable
codes. Errors that stop a build early create session diagnostics. Text logs are
retained within a bounded buffer, while job and diagnostic history remains
accessible after a successful rebuild.

The terminal and local web clients render runtime snapshots. They keep only
their own view state, such as selection, search, and scroll position. Actions
return to the runtime and use one validated handler. Build state never comes
from parsing displayed log text.

## Operational behavior

The serve listener reports ready only after binding. Interactive terminals use
Bubble Tea. Redirected and CI output stays line-oriented, and `--no-tui`
forces it. The local web view is `/_markata/` on the serve listener when
`--admin` is set; `markata-go admin` starts the same serve runtime with that
view enabled. Its default host is loopback. Fatal errors remain visible after
the TUI exits. Shutdown closes the server, watcher, and subscribers.

## Validation

Unit tests cover runtime lifecycle, history, diagnostic attachment, subscriber
behavior, and actions. Client tests render synthetic snapshots at normal and
narrow widths and exercise navigation, resize, search, and log scrollback.
Command tests cover TTY selection and local web launch. Manual terminal and
browser checks verify the integrated path.
