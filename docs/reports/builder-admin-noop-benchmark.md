---
title: "Builder Admin Unchanged Build Benchmark"
description: "Measures the verified unchanged build path on a production-shaped local fixture."
date: 2026-11-25
published: true
tags:
  - performance
  - builder-admin
---

# Builder Admin Unchanged Build Benchmark

This report measures the verified unchanged build path for issue #2045. The benchmark uses a local synthetic fixture. It does not use the production Builder or its persistent volume.

## Result

The pre-optimization implementation took a median of **389 ms** across five runs. The optimized implementation took a median of **164 ms** across five runs. This is a **58% reduction**. The optimized runs took 151–173 ms.

The fixture contains 3,341 Markdown posts and 3,689 cache entries from production counts. It also contains 37,655 synthetic live-release files. The benchmark created these files before the timed runs.

## Recheck with the well-known time endpoint

On 2026-11-25, the production-shaped fixture was updated to enable `.well-known/time`. The verified no-op kept the published `time.txt` bytes unchanged, while a full build remains responsible for advancing that timestamp.

Five fresh `complete-verified-no-op` samples took 155, 159, 159, 166, and 172 ms. The median was **159 ms**, within the 300 ms local target. Five live-release fingerprint samples took 108, 109, 112, 113, and 122 ms. The fixture still contained 3,341 source posts, 3,689 input-cache entries, and 37,655 release files.

These results prove the time endpoint does not prevent the fast path on this fixture. They do not prove the same latency on the production Builder volume; that still requires an isolated non-promoting production check.

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
