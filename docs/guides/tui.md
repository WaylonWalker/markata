---
title: "Terminal UI"
description: "Browse posts and feeds with markata-go's terminal UI"
date: 2026-08-02
published: true
tags:
  - documentation
  - tui
  - navigation
---


Launch the terminal UI with `markata-go tui`.

The interface appears immediately with a loading state while site data is
prepared in the background. You can quit during loading; navigation becomes
available when the initial data arrives.

The TUI reuses `.markata/cache/list.json` between launches. It checks file
metadata and the configured globs, then refreshes only changed posts. Ignored
directories are skipped during discovery, so generated output and dependency
trees do not need to be scanned. On a cold cache, changed posts are parsed in
parallel, including source-encrypted posts.

## Filter the posts you can see

Press `/` in the posts or feeds view and start typing. Results update
immediately and are limited to the items currently in the list. Post matching
checks the title, path, description, tags, and content; feed matching checks
the feed name, title, and path. Press `Enter` to keep the filtered list, or
`Esc` to cancel and restore the list.

This is intentionally a simple list filter; it does not parse filter
expressions. Use the regular CLI list command when you need structured filters.

## Feed inventory

Press `f` to open the feed inventory. It includes every configured feed,
including private-enabled, sidebar-hidden, and automatically generated feeds.
Their titles include explicit labels such as `[private]`, `[hidden]`, and
`[auto]`; the labels are also color-coded when the terminal supports color.

Press `Enter` on any feed to view its posts.
