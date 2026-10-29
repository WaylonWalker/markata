# CLI UX Specification

## Overview

The `markata-go` CLI follows a human-first, script-friendly interaction model.
This specification defines shared user experience rules for command output,
color, prompts, help text, and error handling.

The design reference is [clig.dev](https://clig.dev).

## Shared Principles

- Commands MUST be pleasant for humans to run directly.
- Commands MUST remain composable in shell pipelines.
- Defaults SHOULD favor interactive human use.
- Machine-readable modes SHOULD be explicit.
- Help and errors SHOULD teach the next step.

## Output Streams

### Serve sessions

`serve` uses the Serve Control Center TUI when both input and output are real
terminals and interactivity is enabled. `--no-tui` forces line-oriented output.
Redirected input/output, `--no-input`, and dumb terminals also use the plain
renderer. The TUI receives structured runtime snapshots; it does not parse log
text for build state. `serve` on loopback and `admin` expose the same session in
a loopback-only web dashboard at `/_markata/`. `serve --no-admin` disables it;
`serve --admin` explicitly enables it on loopback. The `a` alias selects `admin`.

When the interactive TUI is active, it owns the terminal. Build and plugin
logs MUST pass through the Serve logger observer into runtime logs and
diagnostics; they MUST NOT write directly to stdout or stderr. This policy is
reinstalled after each manager is configured, including every rebuild.
Plugins that run child processes MUST capture their output and forward it
through the configured logger so the TUI can display it without terminal
corruption.

Warnings and errors observed during a serve session are retained as runtime
diagnostics. A fatal TUI error is reported after the alternate screen exits.

### `stdout`

Commands MUST send primary results to `stdout`.

Examples:
- `config show`
- `config get`
- `list posts --format json`
- `version`
- `explain`
- generated passwords or other command results intended for piping

### `stderr`

Commands MUST send operational messages to `stderr`.

Examples:
- progress updates
- verbose logs
- warnings
- prompt text
- validation diagnostics that are not part of the command's primary result
- user-facing error messages emitted by the command runner

### Writer Discipline

Command handlers SHOULD use Cobra command writers instead of hard-coded
`os.Stdout` and `os.Stderr`.

- Use `cmd.OutOrStdout()` for result output.
- Use `cmd.ErrOrStderr()` for progress, warnings, and errors.

This makes commands testable, consistent under redirection, and safe to compose.

## Color and Terminal Detection

Color MAY be used to increase scannability, but MUST remain optional.

By default, color MUST be enabled for a TTY stream and disabled for a redirected
stream. `NO_COLOR` MUST disable automatic color when it is set to any non-empty
value, including whitespace. `TERM=dumb` MUST also disable automatic color.

Explicit options have this precedence:

- `--no-color` and plain log format disable color
- `--color` or a non-empty `FORCE_COLOR` value other than `0` force color,
  including for redirected output and when `NO_COLOR` or `TERM=dumb` is set
- otherwise, color follows the TTY and environment rules above

Commands MAY also offer explicit formatting controls for operational logs.
When they do, the supported modes SHOULD include:

- `plain` - unstyled, stable text for copy/paste and scripting
- `rich` - colored, more scannable human-oriented logs
- `auto` - choose `rich` for interactive terminals and `plain` otherwise

If both force-enable and force-disable color flags exist, commands MUST reject
conflicting combinations such as `--color` with `--no-color`.

Rich log format MAY use structured logging metadata to drive presentation.
When phase metadata is available, lifecycle-oriented logs SHOULD colorize by
phase rather than by component name alone. If a site palette is configured,
rich logs SHOULD prefer palette-derived colors when they remain readable.

All CLI surfaces MUST resolve their semantic colors from the same terminal
theme. When a site palette is configured, logs, help, prompts, and TUIs MUST use
that palette. When no palette is configured, they MUST use the shared CLI
fallback theme. Palette pickers MAY show each candidate palette in its preview,
but their controls and labels MUST use the active CLI theme.

Commands SHOULD make color decisions per stream. For example, `stderr` MAY stay
colored when `stdout` is piped.

When color is enabled for an interactive stream, human-facing summaries SHOULD
use color by default to improve scannability.

When a site configuration and palette are available, human-facing CLI headings,
statuses, setup snippets, and help text MUST use palette-derived semantic
colors. Commands MUST fall back to the default CLI theme if configuration or
palette resolution fails. Primary structured output such as JSON, YAML, TOML,
and line-oriented data MUST remain uncolored.

Commands MUST preserve readable plain output when color is disabled.

## Output Modes

### Human-Friendly Default

Commands SHOULD default to concise, readable output for humans.

### `--json`

Commands that expose structured records SHOULD offer `--json` or an equivalent
structured format flag.

### `--plain`

Commands SHOULD offer `--plain` when their default human output adds styling,
wrapping, or layout choices that make scripting less reliable.

### `--quiet`

Commands with non-essential progress or status output SHOULD support
`-q, --quiet`.

- Quiet mode MUST suppress non-essential chatter.
- Quiet mode MUST NOT suppress primary command results.
- Quiet mode MUST NOT suppress fatal errors.

## Interactive Input

Commands MUST prompt only when `stdin` is a TTY.

Commands that may prompt SHOULD support `--no-input`.

- If `--no-input` is passed, the command MUST not prompt.
- If required values are missing, the command MUST return a clear error telling
  the user which argument or flag to pass.

Commands MAY offer richer TUI flows when `stdin` and `stdout` are terminals.
Those flows MUST degrade to plain text or fail clearly in non-interactive use.

## Help Text

Top-level help MUST include:

- a short description
- common example commands
- the most important global flags
- how to get command-specific help
- a documentation path
- an issue-reporting path

Subcommand help SHOULD lead with the common workflow and examples before edge
cases.

Running `markata-go` with no arguments MUST display help.

Interactive help MUST use the resolved CLI theme for headings and command
names. Help rendering MAY load the nearest site configuration to resolve the
palette, but MUST fall back cleanly when it cannot be loaded.
Clear section labels in long-form help text (for example, `Examples:` and
`Documentation:`) MUST render as themed Unicode section separators.

## Section Separators

Human-facing section separators MUST use Unicode box-drawing characters rather
than ASCII `-` or `=` rules. A separator SHOULD be 80 columns when space allows
and MUST not exceed the current terminal width. On terminals narrower than 32
columns, commands MUST use a compact labeled separator instead of a padded
rule. Redirected output retains the Unicode separator but has no ANSI color.

## Error Handling

The north star is that a user should almost never need to paste a Markata
error or warning into a search engine or AI merely to discover the obvious
next action. Markata SHOULD provide that action whenever it has enough context
to make a useful recommendation.

Commands SHOULD return errors instead of calling `os.Exit()` directly.

Command invocation and usage errors MUST return exit code `2`.

Examples:

- unknown commands
- unknown flags
- conflicting flags
- missing required positional arguments

Error messages SHOULD:

- explain what went wrong
- name the relevant flag, argument, path, or topic when possible
- suggest the next command or action when there is a clear fix

Unknown commands, long flags, and configuration keys MUST offer useful,
deterministic suggestions when context supports them. Suggestions MUST be
labeled as possibilities; the CLI MUST never run or write a guessed command or
key. When Markata cannot confidently infer the intended input, it MUST NOT
silently execute a guess. It SHOULD still present the most useful
context-specific recovery actions available. A command typo MUST NOT fall
through to the root single-file build shortcut. That shortcut accepts only
Markdown paths.

Command intent has three levels: an explicit alias executes the command; a
semantic synonym suggests the command; a fuzzy typo or prefix suggests the
command. Suggestions never execute. `serve` accepts `s` and `serv` as explicit
command aliases and suggests itself for `dev` and `preview`. Its `--bind` flag is an alias for
`--host`; both set the same listening address, and conflicting values are a
usage error.

Configuration key suggestions apply to `config get` and `config set`. Missing
file keys remain errors even when a similarly named key exists. Validation of
closed-choice configuration values SHOULD show the closest permitted value and
the set of choices. Diagnostics MUST never echo secret values.
Config file loading MUST preserve unknown plugin sections. It MAY reject a
likely typo of a built-in key within a known section or a close spelling of a
multiword top-level key; it MUST name the file and suggested key.

Fatal errors are printed once to `stderr`. An interactive `stderr` uses themed
color to distinguish the error, candidate commands, and next step. Redirected
output, `--no-color`, `NO_COLOR`, and plain log mode receive stable text without
ANSI escapes. Ordinary runtime failures do not dump usage text.

Unexpected diagnostic detail belongs in verbose/debug modes, not normal output.

## Core Command Expectations

- `build` and `serve` print progress and warnings to `stderr`, with summaries and
  explicit results remaining readable in both terminal and redirected use.
- `build` SHOULD include a concise benchmark summary in its final result output,
  including estimated wall-time spent on CPU work, network wait, disk read wait,
  disk write wait, and idle time, plus the slowest lifecycle hotspots.
- When outbound HTTP requests occur during a build, the benchmark summary SHOULD
  also include the slowest requests with enough context to identify the owning
  stage/plugin and remote endpoint without exposing query-string secrets.
- `build` SHOULD offer an explicit machine-readable benchmark mode such as
  `--benchmark-json` for scripting and regression analysis.
- Detailed per-stage benchmark summaries SHOULD be opt-in via verbose mode or a
  dedicated benchmark detail flag, rather than always-on in normal output.
- `config`, `list`, `version`, and `explain` keep their primary output on
  `stdout`.
- When an explicit site directory is selected, `list` and `search` MUST return
  absolute source paths so automation can edit results without relying on its
  current directory.
- `new` and `init` support interactive prompts by default and honor
  `--no-input`.
- `lint` MAY style severities, but the output MUST remain readable without
  color.
