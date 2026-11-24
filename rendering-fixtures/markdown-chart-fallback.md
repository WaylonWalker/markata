---
title: "Chart figures and accessible data"
description: "Rendering fixture for chart captions, disclosure states, and legacy charts"
date: 2026-10-26
published: true
slug: markdown-chart-fallback
---

# Chart figures and accessible data

## Compact horizontal figure

The data starts collapsed, but its summary remains usable without JavaScript.

```chartjs
{
  "title": "Build time",
  "caption": "Seconds per build; lower is better. Full & incremental builds.",
  "data_display": "collapsed",
  "type": "bar",
  "data": {
    "labels": ["Full", "Incremental"],
    "datasets": [{"label": "Seconds", "data": [46, 21]}]
  },
  "options": {"indexAxis": "y", "plugins": {"legend": {"display": false}}}
}
```

## Caption-only figure, default-open data

```chartjs
{
  "caption": "Completed items by month. This longer caption wraps naturally at narrow widths without centering or italic styling.",
  "type": "line",
  "data": {
    "labels": ["January", "February", "March", "April"],
    "datasets": [{"label": "Items", "data": [8, 12, 6, 10]}]
  }
}
```

## Title-only figure with configuration fallback

```chartjs
{
  "title": "Effort versus outcome",
  "data_display": "collapsed",
  "type": "scatter",
  "data": {"datasets": [{"label": "Observations", "data": [{"x": 1, "y": 2}, {"x": 3, "y": 5}]}]}
}
```

## Legacy chart, no figure

This chart retains the original container and expanded no-JS table.

```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Docs", "Tests"],
    "datasets": [{"label": "Items", "data": [8, 12]}]
  }
}
```
