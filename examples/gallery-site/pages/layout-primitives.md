---
title: "Rich Markdown layout primitives"
description: "Composed boxes, diagrams, and side-by-side Markdown tables."
date: 2026-10-20
published: true
slug: layout-primitives
---

# Layout primitives

All examples below render as static HTML with no JavaScript.

## Architecture diagram

::: md-layout

:::: md-box {data-align="center"}
### Falcon3 — Kubernetes
Two deployments, shared published site volume.
::::

:::: md-columns {data-gap="sm"}

::::: md-card {data-align="center"}
### Website
Nginx serves current release.
:::::

::::: md-card {data-align="center"}
### Builder Admin
Builds and publishes releases.
:::::

::::

:::: md-connector {data-direction="both"}
::::

:::: md-box {data-align="center"}
### Shared `/data/site` Longhorn PVC
Replicated storage: NVMe + HDD.
::::

:::

## Side-by-side tables

::: md-columns

:::: md-card
### Current

| Node | Capacity |
| --- | ---: |
| Falcon2 | 3 TB |
| Falcon3 | 4 TB |
::::

:::: md-card
### Proposed

| Node | Capacity |
| --- | ---: |
| Falcon2 | 3 TB |
| Falcon3 | 14 TB |
::::

:::

## Three-column grid

::: md-grid {data-columns="3" data-gap="lg"}

:::: md-box
### Fast
Static first.
::::

:::: md-box
### Accessible
Normal Markdown remains readable.
::::

:::: md-box
### Themed
Inherits active palette variables.
::::

:::

## Stack and divider

::: md-stack
:::: md-card
First step: deploy a change.
::::

:::: md-divider
::::

:::: md-card
Second step: verify the result.
::::
:::
