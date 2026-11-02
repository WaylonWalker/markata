---
title: "Content Diagnostics Ledger Design"
description: "Design for deterministic content dispositions and frontmatter diagnostics"
date: 2026-09-10
published: true
tags:
  - development
  - diagnostics
  - content
---


## Decision

Use one manager-owned `diagnostics.ContentLedger`. The ledger stores per-source
state and derives the build summary from those entries. The manager already
owns discovered files, posts, feeds, stage execution, and reset behavior, so it
is the smallest place that can connect all lifecycle boundaries without a new
parallel content inventory.

The existing `diagnostics.Issue` remains the shared issue shape for source
positions and severities. Frontmatter analysis lives beside those checks and
is reused by the loader. The loader continues to be strict for direct parser
callers, but a content parse failure is recorded for that candidate and does
not prevent valid sibling files from reaching later stages.

## State Model

Each discovered path records candidate status, source load status,
frontmatter presence and validity, post creation, public eligibility, render,
emission, diagnostics, reasons, and feed-level selection. Final disposition is
derived from those fields as `emitted`, `shadow`, `excluded`, or
`not_candidate`.

The snapshot sorts paths, reasons, feed names, and diagnostics. Ledger mutation
is synchronized because rendering and writing use worker pools. Duplicate
issues are ignored by their stable source/code/message identity.

## Frontmatter Policy

The parser recognizes an exact first-line `---`, after an optional UTF-8 BOM.
A longer hyphen run, such as `----` or `-----`, is only called suspicious when
the following lines look like a YAML metadata block and an exact closing
delimiter exists. This avoids warnings for ordinary Markdown horizontal rules.
Leading whitespace before a frontmatter-like opening is diagnosed but is not
silently normalized. Missing and malformed closing delimiters are diagnosed.
No-frontmatter Markdown remains valid and quiet.

## Reporting and Testing

Build summaries show counts for discovery, candidates, loading, valid
frontmatter, posts, eligibility, rendering, emission, exclusion, warnings, and
errors. Normal warnings include a relative path and line. `build -v` explains
each excluded or shadow source and its reasons. The existing benchmark JSON
surface carries the sanitized ledger snapshot for later consumers.

The regression suite uses temporary real sites and runs the actual glob, load,
render, feed, and write plugins. It covers valid frontmatter, no frontmatter,
the production `----` case, `-----`, malformed and missing closing delimiters,
leading whitespace, BOM, invalid YAML, duplicate keys, invalid metadata types,
invalid dates, empty frontmatter, and ordinary Markdown controls.
