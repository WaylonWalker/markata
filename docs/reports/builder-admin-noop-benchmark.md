---
title: "Builder Admin Unchanged Build Benchmark"
description: "Measures verified unchanged builds on a local fixture and input scanning on production storage."
date: 2026-11-26
published: true
tags:
  - performance
  - builder-admin
---

# Builder Admin Unchanged Build Benchmark

This report measures the verified unchanged build path for issue #2045. The benchmark uses a local synthetic fixture. It does not use the production Builder or its persistent volume.

## Result

The pre-optimization implementation took a median of **389 ms** across five runs. The optimized implementation took a median of **164 ms** across five runs. This is a **58% reduction**. The optimized runs took 151–173 ms.

The fixture contains 3,341 Markdown posts and 3,689 cache entries from production counts. It also contains 37,655 live-release files, including the `.well-known/time` output. The benchmark created these files before the timed runs.

## Recheck with the well-known time endpoint

On 2026-11-25, the production-shaped fixture was updated to enable `.well-known/time` and include that output in the live release. The verified no-op benchmark checks that the published timestamp bytes remain unchanged, while a full build remains responsible for advancing them.

Five fresh `complete-verified-no-op` samples took 151, 154, 158, 167, and 187 ms. The median was **158 ms**, within the 300 ms local target. Five live-release fingerprint samples took 99, 99, 101, 109, and 113 ms. The fixture still contained 3,341 source posts, 3,689 input-cache entries, and 37,655 release files, one of which was `.well-known/time`.

These results prove the time endpoint does not prevent the fast path on this fixture. They do not prove the same latency on the production Builder volume; that still requires an isolated non-promoting production check.

## Managed cache mount replay

On 2026-11-25, the fixture was extended to use the production-style `.markata` and `.markata-cache` symlinks to a dedicated cache mount. It includes 3,689 files below the managed build cache. Five fresh `complete-verified-no-op` samples took 175, 214, 179, 187, and 180 ms. The median was **180 ms**. Five live-release fingerprint samples took 99, 127, 128, 119, and 110 ms. Managed cache trees are checked for unsafe symlinks but are not treated as semantic inputs; a separate test mutates their files between snapshots and verifies that the no-op fingerprint stays stable.

This exercises the cache-mount path on a local Linux fixture. It does not measure the production volume or prove that the deployed Builder saves and consumes a baseline; the latest production runs still need a post-fix check.

## Read-only production input scan

On 2026-11-26, a temporary probe pod on `falcon3` mounted the production source PVC read-only. The comparison scanned the same `/data/source/public` and `/data/source/static` trees, 54,851 filesystem entries in total. The serial sample uses the previous `DirEntry.Info` traversal; the candidate uses four bounded `Lstat` workers and writes identities in the original sorted walk order. Both paths produced the same fingerprint and entry count.

Ten samples, each averaging three scans, measured a serial median of **395 ms** (356–496 ms) and a candidate median of **329 ms** (299–381 ms), a **17% reduction** for one input scan. Four workers were the best tested setting on this volume. Eight or more workers were slower or less stable.

This measures only the semantic input-tree scan. It does not include the second stability scan, managed-cache validation, the live-release scan, API handling, or a complete Builder Admin no-op. It also does not show that production is running this candidate: all live Builder and editor containers remain on `sha-05000c9`. The candidate probe pod mounted the source PVC read-only, did not enqueue a build, and did not write a baseline or release.

## Method

The benchmark ran on Linux x86_64 with Go 1.27.0. The host exposed 12 processors. The fixture used `tmpfs` storage.

The complete benchmark called `verifiedNoOp` and included build-record persistence. It did not include queue wait, HTTP handling, Kubernetes storage, or production service load. Each timed run used one request after the fixture setup.

The candidate runs file metadata reads in parallel. It writes identity fields as fixed-width bytes. Its final stability pass compares source metadata and input identities without loading the full configuration again.

The baseline sample times were 345, 392, 401, 389, and 302 ms. The candidate sample times were 164, 165, 151, 173, and 157 ms. The medians were 389 ms and 164 ms.

The baseline live-release fingerprint took a median of 300 ms across five runs. The candidate took a median of 109 ms across five runs. The candidate scan still checks every release path and its Linux file identity.

The benchmark command was:

```sh
go test ./pkg/builderadmin -run '^$' \
  -bench '^BenchmarkVerifiedNoOpProductionShape/(live-release-fingerprint|complete-verified-no-op)$' \
  -benchtime=1x -count=5
```

Both versions used the same command and fixture. The benchmark source is `pkg/builderadmin/noop_test.go`.

The baseline came from the pre-optimization candidate on base commit `8e92f37556ba2ee35832a1453c06840d656a8215`. Its tracked patch SHA-256 was `7372f43fe479c9a6f93f6f4074894b6359a99b8a5eb6edca4e975f2785ec564f`. Its original `noop.go` SHA-256 was `bfd81327b55a7ee5029a58b63f031352987801f583cc51edfc70ca80c37f4546`. Its benchmark test SHA-256 was `9d39273526115ce9de28fc3370978a91e6f73e9f3b66ee5a3e5b110ae88dd134`. Treat the reported reduction as a historical comparison. The pre-optimization source is not part of the current benchmark target.

## Limits

This local result meets the 300 ms target on the synthetic fixture. It does not prove the same time on the production Builder or its storage. A production measurement still needs the candidate binary and an isolated no-op request.
