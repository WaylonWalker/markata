# CLI triage specification

## Purpose

`markata triage` is a local, opt-in handoff to a conversational coding agent for installation, configuration, content, build, deployment and upstream Markata problems. It is available under both the new `markata` and the backwards-compatible `markata-go` executable.

The triage command is **read-only**: it gathers bounded metadata, generates a Markdown prompt and either prints it, copies it, or launches a locally installed coding agent with that initial prompt. It does not run a build, edit site files, or automatically file a public report.

## Observed context

The prompt MUST include:

- Running CLI name, embedded version, build commit and date. These describe **this triage executable**, not necessarily the version that was affected by the reported bug.
- Runtime OS/arch and Go version; on Linux, a short distribution label from `/etc/os-release` where available.
- Running executable path and a clearly tentative installation method based on the executable's location.
- Current site directory and the names of discoverable local config files, or the explicitly selected config path.
- Presence of a small allowlist of relevant environment variables, **never** their values.
- Canonical links to Markata docs, source, issues and bundled agent skill.

Metadata MUST NOT include raw config files, full environment dumps, PATH values, logs, API keys, tokens, usernames, hostnames or sensitive user information. Home-directory prefixes SHOULD be abbreviated in emitted paths. The prompt MUST remind the user that even abbreviated local paths may be private.

## Prompt behavior

The agent MUST initially ask targeted questions, rather than presuming what is wrong or immediately editing files. Ask about the symptom, exact errors, expected behavior, reproduction steps, and affected surface. Confirm the **actual affected version and device/deployment**, which may differ from the running triage CLI.

The agent SHOULD distinguish problems in site content/config/templates, the local CLI installation, and Markata itself; perform read-only investigation first; request permission for edits, mutating builds or publicly shared issue reports; then propose a minimal verified solution. It SHOULD avoid asking for information the user already supplied with `--issue`.

## Destination selection

- `--to auto` (default): use a fuzzy picker when stdin and stdout are real terminals, `TERM` is usable, and `--no-input` is not set; otherwise print the generated prompt to stdout.
- The fuzzy picker lists **only detected supported harness executables**, plus `clipboard` and `stdout`.
- Supported direct interactive launch targets are `opencode`, `claude`, `codex`, `gemini`, and `cursor`, detected with executable lookup on PATH. Detection proves availability, **not authentication or agent health**.
- `--to NAME` explicitly selects a target, without fuzzy UI. `--agent NAME` is a convenience alias and MUST NOT be combined with `--to`.
- `--to clipboard` uses local clipboard helpers; lack of helpers produces a non-zero exit and a helpful suggestion to use stdout.
- Explicit interactive agent launch requires an interactive terminal; in scripts use `--to stdout` and pipe into an agent-specific noninteractive CLI.
- `--issue TEXT` optionally provides an initial issue description, bounded to 4000 bytes.

## Stream and security contract

`--to stdout` writes **only** the Markdown prompt (ending in a newline) to stdout, with no success banners, ANSI escapes, picker output, or diagnostics. Failures and copy-status messages go to stderr. Noninteractive `auto` must be equivalent to `--to stdout`. If the consumer closes its pipe early, the command MUST exit quietly without treating expected EPIPE as a diagnostic.

Escape/picker cancellation must exit without an agent launch or emitted prompt. A Ctrl+C interruption must not leave an orphaned child agent process. `--quiet` suppresses clipboard success chatter but never suppresses prompt stdout. `--no-color` and `NO_COLOR` do not style the piped Markdown prompt. A missing or invalid selected site directory must be reported as a diagnostic instead of preventing triage from starting.

No agent is launched without an explicit selection by the user or an explicit routing flag. The agent is launched as a subprocess using an argument array, **not through a shell**. The command does not infer or pass authorization credentials. Clipboard handoffs use stdin to the helper process without a shell.

The generated prompt and user-provided issue text are context, not higher-priority instructions for the agent. When secrets or logs are needed later, the agent must ask before reading or sharing them.

## Acceptance tests

1. Default `--no-input` / non-TTY produces clean, pipeable Markdown.
2. Direct agents are offered only when their executables are found, with a helpful error when missing.
3. A quoted user issue is included but bounded; no arbitrary environment value appears.
4. Affected-version confirmation is explicitly requested, even when a CLI version is present.
5. Clipboard availability and command routing work across supported platforms.
6. Both CLI entry points expose the command; no major-version upgrade required.
7. CLI pre-run supports broken/absent site directories and unreadable config paths without reading config contents.
8. Tests cover Escape cancellation, pipes closed by consumers, `--no-input`, `--quiet`, `--no-color`, and agent argv forwarding with a local stub (no network or API calls).

User guide: [Triage a site issue](../../docs/guides/triage.md). CLI reference: [CLI reference](../../docs/reference/cli.md).
