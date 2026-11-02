---
title: "PostIndex Correctness and Consolidation"
description: "Contract and audit for the shared Markata-Go post lookup index"
date: 2026-09-03
published: true
tags:
  - documentation
  - plan
  - lifecycle
---


## Goal

Keep `lifecycle.Manager.PostIndex()` correct across manager resets and use it
where an existing production lookup has the same source, key semantics, and
lifecycle. Do not replace filtered collections, cache snapshots, validation
indexes, or aggregation maps.

## PostIndex contract

`PostIndex()` builds lazily from the manager's current post slice. It is not
restricted to the end of `Load`; any plugin can request it at any stage. The
manager invalidates the cached index after `SetPosts`, `AddPost`, and `Reset`.

| Key | Behavior |
| --- | --- |
| Slug | `BySlug` uses lowercase `Post.Slug`; `LookupBySlug` tries this first |
| Slugified | `BySlugified` uses `models.Slugify`; lookup tries it second, then the slugified key in `BySlug` |
| Href | `ByHref` uses the exact non-empty `Post.Href` |
| Path | `ByPath` uses the exact non-empty `Post.Path` |
| Aliases | String aliases use lowercase and slugified keys; existing alias entries are not overwritten |
| Duplicates | Main slug, slugified slug, href, and path entries are last-post-wins |
| Nil | Nil index receivers and misses return nil from `LookupBySlug`; nil posts are not valid index input |

The maps are shared pointers to the posts, not copies of post values. A direct
mutation of `Slug`, `Href`, `Path`, or alias data on an existing post does not
invalidate the index. Such a mutation requires `PostIndex().Refresh(m)` before
the next lookup. `Refresh` is used once by the wikilinks transform because
other plugins may have changed post fields before that consumer runs.

## Lookup inventory

| Component | Lookup key | Current implementation | Lifetime | Source of truth | Equivalent to PostIndex? | Can replace safely? |
| --- | --- | --- | --- | --- | --- | --- |
| Embeds | slug | Shared `PostIndex.LookupBySlug` | Transform | Manager posts | Yes | Already uses it |
| Wikilinks | slug | Shared `PostIndex.LookupBySlug` plus one refresh | Transform | Manager posts | Yes | Already uses it |
| Wikilink hover | href with slash variants | Shared `PostIndex.ByHref` | Render | Manager posts | Yes, with consumer normalization | Already uses it |
| Link collector | slug/href with slash normalization | Shared `PostIndex` maps | Render | Manager posts | Yes, with consumer normalization | Already uses it |
| `include_post` template function | exact trimmed slug | Linear `Manager.Posts()` scan | Render/template call | Manager posts | No: case, slugification, alias, and duplicate semantics differ | No |
| Post service `Get` | exact path | Linear `Manager.Posts()` scan, first match | Service lifetime | Manager posts | No: `ByPath` is last match on duplicates | No |
| Prev/next | post slug to feed membership | `map[string][]*Feed` | Collect call | Selected feed posts | No: membership and order are the result | No |
| Templates/sidebar | current slug in selected feeds | Linear membership/position scans | Render/collect call | Selected feed posts | No | No |
| Mentions | configured post filter | Collection scan and filter | Transform | Manager posts | No: aggregation/filtering | No |
| Slug conflicts | normalized slug to all paths | `map[string][]string` | Validation call | Publishable posts | No: duplicate reporting is intentional | No |
| Search | path to visible post | Filtered `search.PostsByPath` map | Search request/index | Search-visible posts | No: drafts/skips are excluded | No |
| Search API | path to visible post | Handler-owned snapshot map | Handler lifetime | Handler update snapshot | No: independent lifecycle | No |
| List cache | source path | Cache reconstruction map | Load/cache operation | Cached files and feeds | No: needed before the manager collection is final | No |
| Incremental serve/load | source path | Cached post maps | Build/serve cycle | Serve or build cache | No: cache state | No |
| Publish HTML | path to slug | Transient invalidation map | One write operation | Current post paths and cache | No: reverse projection | No |
| Content index | path eligibility | Filtered eligibility map | Content-index operation | Content-index policy | No: eligibility is the value | No |
| LSP | slug/URI to `PostInfo` | Separate locked workspace index | LSP session | Workspace files | No: separate model and lifecycle | No |

No local post lookup map met all four replacement conditions: identical source
posts, normalization, lifecycle visibility, and duplicate/result semantics.

## Lifecycle decision

Synthetic posts are added with `AddPost` during `Configure` or late `Load`, so
they invalidate the index before Transform consumers run. A consumer that
needs a filtered subset or a cache-owned snapshot must keep its local lookup.
No plugin is changed to use the shared index solely to reduce a small loop when
doing so would change normalization or duplicate behavior.

## Verification

The regression sequence is:

1. Set post A and build the index.
2. Reset the manager.
3. Confirm the manager has no posts and its next index cannot find A.

Additional tests cover `SetPosts`, `AddPost`, slug case and slugification,
href/path keys, duplicate behavior, direct field mutation followed by explicit
`Refresh`, and the nil lookup receiver.
