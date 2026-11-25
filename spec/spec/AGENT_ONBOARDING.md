# Agent-first installation handoff (`markata --skill`)

Status: initial read-only contract for issue #2095, parent epic #2087.

## Intent

Installing Markata yields a fast static-site builder, not a hosted account,
a deployment agent, or a provider-specific setup wizard. People can supply
their own coding agent and choose their own templates, plugins, and deployment
target. An agent can start without writing any files.

## CLI contract

- `markata --skill` prints a standalone Markdown onboarding guide to standard
  output, ending in a newline. No ANSI color, interactive prompt, logs, JSON
  wrapper, or environment-specific values precede it.
- The guide is embedded in and therefore version-aligned with the executing
  binary. This is the same embedded skill bundle that the existing
  `markata agent install` command uses. It is available offline.
- The command does not need a workspace, a configuration file, agent executable,
  Tailscale, a server, a network connection, or elevated privileges.
- Printing the guide is read-only; no project files, config, credentials,
  editor sockets, or agent-specific installations are created or modified.
- When `--skill` is combined with a positional Markdown filename, fail with
  an actionable error instead of building that file.
- Existing `markata agent install`, `agent update`, and `agent doctor`
  remain the opt-in interface for full, topic-based skill installations.
- The legacy `markata-go` binary continues to share the canonical CLI
  behavior; existing build and serve commands are unaffected.

## Agent instruction and safety contract

The handoff instructs agents to inspect the chosen workspace before changes,
prefer native configuration and templates before custom plugins, and run the
complete production build for deployment verification. A fast development
build is not a release-build substitute. A build is not a deploy: agent-driven
hosting, credential access, public publication, and destructive operations
require user approval and a selected destination. Do not invent credentials,
publish secrets, or assume a proprietary agent/runtime.

A reproducibility claim requires pinned inputs and independently compared
artifacts; the guide itself cannot promise all plugin-generated content is
deterministic.

## Tests

- Verify exact, stable, newline-terminated text from repeated calls.
- Verify a read-only invocation without site configuration.
- Verify output errors propagate; reject file positional arguments.
- Run CLI build, release packaging, documentation-site and cross-platform CI.
