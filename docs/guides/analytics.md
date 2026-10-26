---
title: "Analytics & Visualizations"
description: "Guide to analytics plugins for visualizing data and tracking activity in markata-go"
date: 2024-01-15
published: true
slug: /docs/guides/analytics/
tags:
  - documentation
  - guides
  - analytics
  - visualization
---

# Analytics & Visualizations

markata-go provides several plugins for visualizing data directly in your markdown content. This guide covers the contribution graph and chart visualization plugins.

## Contribution Graph

The `contribution_graph` plugin renders GitHub-style calendar heatmaps showing activity over time. It uses the [Cal-Heatmap](https://cal-heatmap.com/) library.

### Configuration

Enable the plugin in your configuration:

```toml
[markata-go]
hooks = ["default", "contribution_graph"]

[markata-go.contribution_graph]
enabled = true
cdn_url = "/assets/vendor/cal-heatmap"                # Cal-Heatmap base URL (local by default)
container_class = "contribution-graph-container"
theme = "light"
scale_max_percentile = 0                               # Optional global outlier cap
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable the plugin |
| `cdn_url` | string | `/assets/vendor/cal-heatmap` | Cal-Heatmap base URL (local by default) |
| `container_class` | string | `contribution-graph-container` | CSS class for the container |
| `theme` | string | `light` | Color theme (`light` or `dark`) |
| `scale_max_percentile` | number | `0` | Optional global percentile cap for the color scale |

### Basic Usage

Add a contribution graph to your markdown using a fenced code block:

````markdown
```contribution-graph
{
  "data": [
    {"date": "2024-01-01", "value": 5},
    {"date": "2024-01-02", "value": 3},
    {"date": "2024-01-03", "value": 8},
    {"date": "2024-01-04", "value": 2},
    {"date": "2024-01-05", "value": 10}
  ],
  "options": {
    "domain": "year",
    "subDomain": "day"
  }
}
```
````

### Data Format

The `data` array contains objects with:
- `date`: Date string in ISO format (YYYY-MM-DD)
- `value`: Numeric value that determines cell color intensity

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `domain` | string | `year` | Time domain: `year`, `month`, `week`, `day` |
| `subDomain` | string | `day` | Sub-domain: `day`, `hour`, `minute` |
| `cellSize` | number | `10` | Size of each cell in pixels |
| `range` | number | `1` | Number of domain units to display |
| `maxValue` | number | unset | Explicit color-scale maximum |
| `maxPercentile` | number | unset | Per-graph percentile cap for the color scale |

### Examples

#### Year View

Show a full year of activity:

````markdown
```contribution-graph
{
  "data": [
    {"date": "2024-01-15", "value": 5},
    {"date": "2024-02-20", "value": 8},
    {"date": "2024-03-10", "value": 3},
    {"date": "2024-06-01", "value": 12}
  ],
  "options": {
    "domain": "year",
    "subDomain": "day",
    "range": 1
  }
}
```
````

#### Monthly View

Show activity by month:

````markdown
```contribution-graph
{
  "data": [
    {"date": "2024-01-01", "value": 5},
    {"date": "2024-01-15", "value": 8},
    {"date": "2024-02-01", "value": 3}
  ],
  "options": {
    "domain": "month",
    "subDomain": "day",
    "range": 3
  }
}
```
````

#### Outlier Control

When one or two very large days wash out the rest of the graph, cap the color scale at a percentile instead of the raw max:

````markdown
```contribution-graph
{
  "data": [
    {"date": "2024-01-01", "value": 1},
    {"date": "2024-01-02", "value": 2},
    {"date": "2024-01-03", "value": 3},
    {"date": "2024-01-04", "value": 25}
  ],
  "options": {
    "domain": "year",
    "subDomain": "day",
    "range": 1,
    "maxPercentile": 95
  }
}
```
````

Set `scale_max_percentile` in config to apply the same percentile cap to every contribution graph by default.

#### Tracking Blog Post Frequency

Use Jinja templating to automatically generate data from your posts:

````markdown
---
title: My Writing Activity
jinja: true
---

```contribution-graph
{
  "data": [
    {% for post in filter("published == true") %}
    {"date": "{{ post.Date.Format "2006-01-02" }}", "value": 1}{% if not loop.last %},{% endif %}
    {% endfor %}
  ],
  "options": {
    "domain": "year",
    "subDomain": "day"
  }
}
```
````

### Styling

Add custom CSS to style your contribution graphs:

```css
.contribution-graph-container {
  margin: 2rem 0;
  overflow: hidden;
}

/* Cal-Heatmap specific overrides */
.ch-domain-text {
  fill: var(--color-text-muted);
}

.ch-subdomain-bg {
  fill: var(--color-surface);
  stroke: var(--color-border);
}
```

---

## Chart.js Charts

The `chartjs` plugin lets you embed interactive charts using [Chart.js](https://www.chartjs.org/).

### Configuration

```toml
[markata-go]
hooks = ["default", "chartjs"]

[markata-go.chartjs]
enabled = true
cdn_url = "/assets/vendor/chartjs/chart.min.js"
container_class = "chartjs-container"
data_display = "open" # or "collapsed"; per-chart JSON can override
```

For built-in figure cards, add top-level plain-text `title` and/or `caption`
keys to the chart JSON. Use `data_display` to choose its initial data disclosure
state. These keys are removed before Chart.js receives the configuration.
See [chart figures and accessible fallbacks](markdown-chart-fallback.md) for
complete examples and no-JavaScript behavior.

### Basic Usage

Create charts with JSON configuration:

````markdown
```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["January", "February", "March", "April", "May"],
    "datasets": [{
      "label": "Posts Published",
      "data": [3, 5, 2, 8, 4],
      "backgroundColor": "rgba(54, 162, 235, 0.7)"
    }]
  }
}
```
````

### Supported Chart Types

| Type | Description |
|------|-------------|
| `bar` | Bar chart (vertical or horizontal) |
| `line` | Line chart with optional fill |
| `pie` | Pie chart |
| `doughnut` | Doughnut chart |
| `radar` | Radar/spider chart |
| `polarArea` | Polar area chart |
| `bubble` | Bubble chart |
| `scatter` | Scatter plot |

### Examples

#### Line Chart

````markdown
```chartjs
{
  "type": "line",
  "data": {
    "labels": ["Week 1", "Week 2", "Week 3", "Week 4"],
    "datasets": [{
      "label": "Page Views",
      "data": [1200, 1900, 1500, 2400],
      "borderColor": "rgb(75, 192, 192)",
      "fill": false,
      "tension": 0.3
    }]
  }
}
```
````

#### Pie Chart

````markdown
```chartjs
{
  "type": "pie",
  "data": {
    "labels": ["Tutorials", "Guides", "Reference", "Blog"],
    "datasets": [{
      "data": [35, 25, 20, 20],
      "backgroundColor": [
        "#ff6384",
        "#36a2eb",
        "#ffce56",
        "#4bc0c0"
      ]
    }]
  }
}
```
````

### Advanced Options

Chart.js supports extensive configuration. Add an `options` key for customization:

````markdown
```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Q1", "Q2", "Q3", "Q4"],
    "datasets": [{
      "label": "Revenue",
      "data": [12000, 19000, 15000, 22000]
    }]
  },
  "options": {
    "responsive": true,
    "plugins": {
      "title": {
        "display": true,
        "text": "Quarterly Revenue"
      },
      "legend": {
        "position": "bottom"
      }
    },
    "scales": {
      "y": {
        "beginAtZero": true
      }
    }
  }
}
```
````

---

## Use Cases

### Blog Statistics Dashboard

Create a dedicated page showing your blog's activity:

```markdown
---
title: Blog Statistics
template: page.html
jinja: true
---

## Publishing Activity

```contribution-graph
{
  "data": [...generated from posts...],
  "options": {"domain": "year", "subDomain": "day"}
}
```

## Content by Category

```chartjs
{
  "type": "pie",
  "data": {
    "labels": ["Go", "Python", "DevOps", "Other"],
    "datasets": [{"data": [15, 12, 8, 5]}]
  }
}
```
```

### Project Documentation

Track development activity or feature completion:

```markdown
## Development Progress

```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Features", "Bug Fixes", "Tests", "Docs"],
    "datasets": [{
      "label": "Completed",
      "data": [12, 8, 25, 15],
      "backgroundColor": "#4CAF50"
    }, {
      "label": "In Progress",
      "data": [3, 2, 5, 2],
      "backgroundColor": "#FFC107"
    }]
  },
  "options": {
    "indexAxis": "y"
  }
}
```
```

---

## Related

- [[configuration-guide|Configuration Guide]] - Full configuration reference
- [[markdown|Markdown Features]] - Other markdown extensions
- [[templates|Templates]] - Using Jinja in markdown
