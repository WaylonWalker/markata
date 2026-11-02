---
title: Retained builder workspace and incremental publication
description: Preserve immutable releases while avoiding full-tree workspace copies
date: 2026-10-03
published: false
---

## Retained workspace and incremental publication

Issue #1421 tracks repeated full-tree preparation. Production spends approximately
169 seconds in preparation and 30 seconds in the engine. Same-filesystem promotion
consumes the workspace, so the existing reuse marker cannot avoid the next copy.

Retain the independent workspace and stage each release. Compare regular files
with the previous immutable release. Link matching bytes and permissions from
that release, and independently copy changed files. Never link workspace files
into releases. Build a fresh directory tree so deleted paths remain absent.

Direct comparisons avoid a new manifest trust contract. They still read unchanged
output. Logs report linked files, copied files, copied bytes, and compared bytes.
Hard-link failures fall back to independent copies. Constrained filesystem roots
prevent baseline symlinks from redirecting publication outside the release tree.

The existing marker handles successful reuse, interruption, failure, and rollback.
A failed stage never becomes current. Tests cover unchanged and changed bytes,
permissions, symlinks, deletion, unsupported links, isolation, and recovery.
A synthetic benchmark compares full copying with incremental publication.
Production validation measures complete jobs and verifies old-release bytes.

Reflinks alone cannot serve the production ext4 volume. A manifest cache remains
a later optimization after direct comparison provides deployment measurements.

## Production follow-up

The first production publication linked 34,030 files and copied 14, but required
273 seconds. A concurrent storage probe linked 500 files in 0.04 seconds, while
reading their paired contents took 2.83 seconds for 5.3 MB. These observations
identify serial storage latency as a remaining cost. A bounded pool overlaps
up to eight independent file operations. Workers own buffers and counters.
Publication waits for every worker before applying directory metadata or
exposing the release. Exact comparisons and file synchronization remain intact.

### Warm publication digest cache

Builder Admin stores a private SHA-256 manifest beside its mutable workspace.
On Linux, unchanged file identity and change timestamps allow cached digest reuse.
Files with recent or changed timestamps require content hashing. Missing or corrupt
cache records fall back to hashing. Other platforms always hash source files.

The manifest binds the workspace and previous release. Published releases must
stay immutable. Do not edit retained releases in place. Delete the workspace's
`.publication.json` sibling to force content verification on the next publication.
Manifest write failures leave publication usable but disable cache reuse. Logs
report `cached_source_hashes`, `cached_release_hashes`, and `manifest_saved`.
The manifest contains file identities and digests and never appears in a release.
