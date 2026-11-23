# Local workspace supervisor (markata up)

Status: foreground-only initial contract for issue #2090; complementary to
agent-first setup in #2095.

## Purpose

`markata up` is an optional convenience for **local** authoring and preview,
not a platform deployer or replacement for deterministic `markata build`.
It reuses the existing browser editor and `serve` watcher/builder unchanged.

## Process lifecycle

- Start editor and builder as independent child processes in the current
  workspace, using the executing Markata binary. This isolates global CLI
  state and builder crashes from the editor.
- Editor binds to `127.0.0.1:8766`, builder live preview to
  `127.0.0.1:8000`. Override via `--editor-port` and
  `--preview-port`. Both services must bind loopback and use distinct ports.
- `--editor-only` and `--builder-only` are mutually exclusive.
- Check default ports before launch. If occupied, choose the first available
  loopback port in the next 20 candidates and tell the user exactly what
  changed; never bind a different interface, kill or attach to unknown
  listeners. The editor and builder ports remain independent.
- An explicitly supplied `--editor-port` or `--preview-port` must not be
  changed. `--strict-ports` enforces that rule on defaults too. Report
  actionable conflicts and candidate-range exhaustion without hiding errors.
- Do not claim ownership of an existing Markata process based only on a port.
  Same-workspace idempotence/reuse requires a separate verifiable identity
  boundary under issue #2108; until then, a second invocation may create a
  separate editor on another port.
- The command prints **starting URLs**, never claims ready until the individual
  service reports readiness. A builder failure is reported and the editor
  remains usable; an editor failure ends the foreground workspace.
- Ctrl+C/SIGTERM cancels the parent supervisor, signals only child processes
  it created, waits for graceful shutdown, then forcibly terminates survivors
  after a bounded timeout. No persistent PID file, daemon/service installation
  or separate status registry is created.
- Draft saving never automatically deploys publicly. Git push and verified
  live deploy are distinct concepts.

## Explicit non-goals

No cloud account, Tailscale enrollment, remote public listener, persistent
daemon, Kubernetes controller, or generic deployment abstraction. These are
user-selected, agent-managed decisions. Agent-readable guidance is provided
through `markata --skill`, independently of `up`.

## Tests

- Planning, port bounds/collisions and mutually exclusive mode validation.
- Subprocess builder failure does not terminate editor.
- Editor failure terminates remaining child; signal shutdown is bounded.
- Cross-platform lint, race and CI checks including Windows compatibility.
