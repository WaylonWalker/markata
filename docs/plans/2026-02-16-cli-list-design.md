---
title: "CLI List Command Design"
description: "Design for list subcommands and formats for data exploration"
date: 2026-02-16
published: true
tags:
  - documentation
  - plan
  - cli
---


## Goal

Add a `list` command with `posts`, `tags`, and `feeds` subcommands for fast, scriptable data exploration. Output must support a Unix-friendly path-only format for piping.

## Command Shape

- `markata-go list posts`
- `markata-go list tags`
- `markata-go list feeds`

Shared flags:

- `--format` (`table` default, `json`, `csv`, `path`)
- `--sort`
- `--order` (`asc` or `desc`, default `desc`)
- `--filter` (posts only, passed to filter expression engine)

## Output Formats

- `table`: fixed-width ASCII table with headers, no color
- `json`: array of objects with raw values
- `csv`: header row + values, same column order as table
- `path`: one value per line, optimized for piping

Path format fields:

- posts: `path`
- tags: `name`
- feeds: `path` (fallback to `name` if missing)

## Cache

- Cache file: `.markata/cache/list.json`
- Reused when config hash matches and file metadata is unchanged
- Partial refresh re-parses only changed files and rebuilds tags/feeds from cached posts

## Data Sources

- Use `services.App` and `Build.LoadForTUI` for consistent stats and feed data.
- Posts: `Posts.List` with `ListOptions` for sort/filter.
- Tags: `Tags.List`, sorted in CLI when custom sort requested.
- Feeds: `Feeds.List`, sorted in CLI when custom sort requested.

## Sorting

Posts: `date`, `title`, `words`, `path`, `reading_time`, `tags` (service-level).

Tags: `name`, `count`, `words`, `reading_time`.

Feeds: `name`, `posts`, `words`, `reading_time`, `avg_reading_time`.

## Errors

- Invalid format/sort/order: print a clear error and exit code 1.
- Load errors: return the underlying error.

## Docs + Spec

- New spec file for `list` behavior.
- Update CLI reference and add a short guide with examples.
