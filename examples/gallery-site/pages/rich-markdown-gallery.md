---
title: "Rich Markdown component gallery"
description: "Static-first responsive Markdown boxes, tables, status, code, charts, and checklists."
date: 2026-10-21
published: true
slug: rich-markdown-gallery
---

# Rich Markdown component gallery

Every component below keeps the article readable without JavaScript.
Interaction is progressively enhanced when supported, and examples here
use illustrative values rather than live operational data.

## Responsive architecture boxes

::: md-layout

:::: md-box {data-align="center"}
### K3s site publishing

Two deployments share a published site volume.

::::

:::: md-columns {data-gap="sm"}

::::: md-card
### Website

Nginx serves the published release.

:::::

::::: md-card
### Builder

Prepares and publishes the next release.

:::::

::::

:::: md-connector {data-direction="both"}

::::

:::: md-box {data-align="center"}
### Shared persistent volume

The current static site is available to both workloads.

::::

:::

## Lucide icons in cards and feature rows

The `:lucide-…:` shortcode resolves at build time. It does not require
Web Awesome or a browser-time icon library.

::: md-columns

:::: md-card {data-align="center"}

::::: md-icon {data-size="lg"}
:lucide-server:

:::::

### Falcon2
K3s node

::::

:::: md-card {data-align="center"}

::::: md-icon {data-size="lg"}
:lucide-database:

:::::

### Falcon3
NAS + ZFS

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-rocket:

::::

:::: md-feature-body
### Verify production
Confirm the new image is deployed and check live pages without JavaScript.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-gauge:

::::

:::: md-feature-body
### Protect performance
Test narrow screens and simulated slow connections before rollout.

::::

:::

## Independent side-by-side tables

::: md-columns

:::: md-card
### Before

| Node | Capacity |
| --- | ---: |
| Falcon2 | 3 TB |
| Falcon3 | 4 TB |

::::

:::: md-card
### After

| Node | Capacity |
| --- | ---: |
| Falcon2 | 3 TB |
| Falcon3 | 14 TB |

::::

:::

## Status and before/after metrics

::: md-status {data-tone="success"}

Cleanup completed

:::

::: md-layout

:::: md-columns

::::: md-stat
Before

**18%**

Available space

:::::

::::: md-stat
After

**45%**

Available space

:::::

::::

:::: md-delta {data-tone="success"}
Capacity recovered

**+27 percentage points**

::::

:::: md-meter {data-value="45" data-max="100" data-label="Available space after cleanup: 45 percent" data-tone="success"}

45 percent available

::::

:::

## Code with progressive copy controls

~~~bash
kubectl get nodes
kubectl get pods -A
~~~

Without JavaScript, the commands are still ordinary selectable code.

## Theme-aware Mermaid diagram

The diagram keeps its node, label, edge, and cluster colors in step with the
active site palette.

~~~mermaid
flowchart TD
  subgraph input[Input]
    A[Raw water]
  end
  A -->|clean| B{Ready?}
  B -->|yes| C[Storage]
  B -->|no| D[Filter]
  D --> B
~~~

## Native rendered/source disclosure

### Preview

This sentence has **emphasis** and a [link](https://example.org).

::: md-source {data-label="Show the Markdown"}

~~~markdown
This sentence has **emphasis** and a [link](https://example.org).
~~~

:::

The disclosure opens with the keyboard or pointer without any JavaScript.

## Interactive checklist

::: md-checklist {data-id="gallery-release-steps"}

- [x] Build documentation
- [ ] Check tests
- [ ] Verify the deployment

:::

With JavaScript disabled, the list remains readable with disabled boxes.

## Chart.js with a static data fallback

~~~chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Full", "Incremental"],
    "datasets": [{"label": "Seconds", "data": [46, 21]}]
  }
}
~~~

The table of exact values remains available if Chart.js fails or JavaScript
is blocked. Interactive charts remain an opt-in enhancement.

## Native disclosure

<details>
<summary>Show verification details</summary>

These HTML disclosure controls work without a JavaScript runtime.

</details>

## Native Markdown tabs

Tabs render fully on the server, even if JavaScript is disabled.

::: md-tabs {style="boxed" data-label="Install options"}

:::: md-tab {label="curl"}

```bash
curl -fsSL https://example.com/install.sh | bash
```

::::

:::: md-tab {label="mise"}

```bash
mise use -g github:WaylonWalker/markata
```

::::

:::: md-tab {label="Go"}

```bash
go install github.com/WaylonWalker/markata/cmd/markata@main
```

::::

:::
