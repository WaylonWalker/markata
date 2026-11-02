---
title: "Actionable CLI Error Design"
description: "Aliases and clear suggestions for command, flag, and configuration mistakes"
date: 2026-09-26
published: false
tags:
  - design
  - cli
---

## Actionable CLI errors

Issue: #1288

## Problem

The root command accepts a positional Markdown path. A mistyped command such
as `markata-go s` can reach the build path. `serve --bind` fails even though
the existing `--host` flag performs that action. Missing configuration keys
often produce a bare `key not found` error, and misspelled keys in config files
can be silently ignored.

## Design

Keep Cobra as the parser and keep canonical command names. Add memorable
aliases for `build`, `serve`, and `list`, plus `serve --bind` as an alias for
`--host`. Reject conflicting `--bind` and `--host` values. Limit the root
single-file shortcut to `.md` and `.markdown` paths.

Use one small edit-distance matcher for command names, long flags, and config
keys. Rank candidates deterministically and show at most three close names.
Show suggestions as possibilities and never run a guessed command. If no name
is close, state that the CLI cannot infer the intended input and point to the
exact help command. Keep diagnostics on `stderr`; color the status and next
step only when the error stream supports color.

For config files, detect close spellings of known settings before conversion
to typed configuration. Preserve extension fields and plugin sections. A
top-level unknown key is diagnosed only when it closely resembles a built-in
multiword key; nested keys are checked within their known
parent section. Validation errors for finite choice values list the choices
and suggest a close value.

## Verification

Test alias resolution, command and flag suggestions, no-match fallback, config
key typo diagnostics, preservation of plugin sections, and plain error output.
Exercise actual CLI success and failure paths, then run relevant Go tests and
lint.
