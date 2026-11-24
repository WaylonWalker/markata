# Reading the numbers

*Illustrative values only · Four reporting periods · No live data source*

The first chart answers **what is largest**. The second answers **what
changed**. They use the same underlying activity as an example, but
present different questions instead of repeating one visual twice.

## At a glance

::: md-columns {data-gap="lg"}

:::: md-card

### 24 items

Total across the four illustrative periods.

::::

:::: md-card

### 9 in the latest period

The largest single-period value in this example.

::::

:::

## Compare categories

The bars make differences between groups easy to compare.

```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Guides", "Notes", "Tutorials", "Reference"],
    "datasets": [
      {"label": "Items", "data": [9, 6, 5, 4]}
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {
      "legend": {"display": false}
    },
    "scales": {
      "y": {
        "beginAtZero": true,
        "title": {"display": true, "text": "Items"}
      }
    }
  }
}
```

Guides are the largest category in this invented dataset. This
comparison is easier to read as bars than as a list of percentages.

## Follow the trend

The line emphasizes the sequence of values over four reporting periods.
Give it the full reading width, rather than fitting it inside a small
stat card.

```chartjs
{
  "type": "line",
  "data": {
    "labels": ["Period 1", "Period 2", "Period 3", "Period 4"],
    "datasets": [
      {
        "label": "Items",
        "data": [3, 5, 7, 9],
        "fill": false,
        "tension": 0.2
      }
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {
      "legend": {"position": "bottom"}
    },
    "scales": {
      "y": {
        "beginAtZero": true,
        "title": {"display": true, "text": "Items"}
      }
    }
  }
}
```

The illustrated value rises each period. We can describe the increase,
but not infer why it happened from the chart alone.

## Show composition

A doughnut emphasizes proportions of one **known whole**. Avoid using
it when there are many small categories or readers need exact rankings.

```chartjs
{
  "type": "doughnut",
  "data": {
    "labels": ["Guides", "Notes", "Tutorials", "Reference"],
    "datasets": [
      {
        "label": "Items",
        "data": [9, 6, 5, 4]
      }
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {
      "legend": {"position": "bottom"}
    }
  }
}
```

The values are the same as the category bars. In a **real** report,
normally choose the bar chart **or** the doughnut rather than placing
both adjacent: this example includes both solely to demonstrate the
available chart types.

## Read the exact values

Chart.js's static-first Markata renderer creates a `Chart data`
disclosure from these aligned label/data arrays at build time. The
table stays available with JavaScript off or when the chart can't load.

Do not mix actual measurements with this example. Replace all values,
labels, period descriptions and interpretation together when adapting
it to a real article. Check dark/light palette contrast and narrow
screen axes before publishing.

---

*Authoring note: When the site template already renders the post title,
don't repeat this source H1 in the body.*
