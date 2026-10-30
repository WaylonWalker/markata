---
title: "Storage cleanup: rich Markdown metrics"
description: "An illustrative before-and-after report with badges, meters, and data tables."
date: 2026-10-20
published: true
slug: layout-metrics
---

# Falcon3 NVMe cleanup

::: md-status {data-tone="success"}

NVMe cleanup verified

:::

**Excellent result.** The containerd prune reclaimed approximately **129 GiB**
of NVMe space in this illustrative report.

::: md-layout

Falcon3 NVMe — before and after cleanup

:::: md-columns

::::: md-stat {data-tone="caution"}
Before

**83 GiB**

18% free

:::::: md-meter {data-value="18" data-max="100" data-label="Before cleanup: 18 percent free" data-tone="caution"}

::::::

:::::

::::: md-stat {data-tone="success"}
After

**212 GiB**

45% free

:::::: md-meter {data-value="45" data-max="100" data-label="After cleanup: 45 percent free" data-tone="success"}

::::::

:::::

::::

:::: md-delta {data-tone="success"}
Space recovered

**+129 GiB**

::::

:::

| Storage component | Before | After | Reduction |
| --- | ---: | ---: | ---: |
| Containerd snapshots | 158 GiB | 76 GiB | 82 GiB |
| Containerd content blobs | 58 GiB | 7.3 GiB | 50.7 GiB |
| Entire K3s directory | 220 GiB | 91 GiB | 129 GiB |
| NVMe used | 383 GiB | 254 GiB | 129 GiB |

The differences among totals reflect rounded measurements taken at slightly
different times. This is test content, not a live monitoring dashboard.

## Example build progress

::: md-progress {data-value="7" data-max="10" data-label="Build steps completed"}

7 of 10 build steps finished.

:::

## Example trend sparkline

::: md-stat {data-tone="success"}
NVMe used

**254 GiB**

Down from a 383 GiB peak after cleanup

:::: md-sparkline {data-values="310 342 361 383 371 254" data-label="NVMe used in GiB, last six checks" data-tone="success"}

Last six checks

::::

:::
