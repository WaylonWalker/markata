---
title: "Garden View"
description: "Explore your content as a knowledge graph with tag clusters and relationship visualization"
date: 2024-01-15
published: true
slug: /docs/guides/garden/
tags:
  - documentation
  - garden
  - knowledge-graph
---


The garden view plugin generates a knowledge graph of your site's content, showing how posts relate through tags and internal links. It exports the graph as JSON and renders an interactive garden page.

## Quick Start

The garden view is enabled by default. After building your site, you'll find:

- `/garden/graph.json` -- the full knowledge graph as JSON
- `/garden/index.html` -- an HTML page showing tag clusters

No configuration is needed for the default behavior.

## Configuration

Add a `[markata-go.garden]` section to your config file to customize the garden view:

```toml
[markata-go.garden]
enabled = true              # Enable the garden view (default: true)
path = "garden"             # Output path prefix (default: "garden")
export_json = true          # Emit graph.json (default: true)
render_page = true          # Generate garden/index.html (default: true)
include_tags = true         # Include tag nodes in graph (default: true)
include_posts = true        # Include post nodes in graph (default: true)
max_nodes = 2000            # Maximum number of nodes (default: 2000)
exclude_tags = []           # Tags to exclude from graph (default: [])
template = "garden.html"    # Template for the garden page (default: "garden.html")
title = "Garden"            # Page title (default: "Garden")
description = ""            # Page description (default: "")
```

## Graph JSON

The `graph.json` file contains the full knowledge graph with nodes and edges:

```json
{
  "nodes": [
    {
      "id": "post:my-post",
      "type": "post",
      "label": "My Post Title",
      "href": "/my-post/",
      "tags": ["go", "tutorial"],
      "date": "2024-01-15T00:00:00Z",
      "description": "A short description"
    },
    {
      "id": "tag:go",
      "type": "tag",
      "label": "go",
      "href": "/tags/go/",
      "count": 12
    }
  ],
  "edges": [
    {
      "source": "post:my-post",
      "target": "post:another-post",
      "type": "link"
    },
    {
      "source": "post:my-post",
      "target": "tag:go",
      "type": "tag"
    },
    {
      "source": "tag:go",
      "target": "tag:tutorial",
      "type": "co-occurrence",
      "weight": 5
    }
  ]
}
```

### Node Types

| Type   | ID Format        | Description                    |
|--------|------------------|--------------------------------|
| `post` | `post:<slug>`    | A published post               |
| `tag`  | `tag:<name>`     | A tag used by one or more posts|

### Edge Types

| Type            | Source -> Target | Description                              |
|-----------------|------------------|------------------------------------------|
| `link`          | post -> post     | Internal link from one post to another   |
| `tag`           | post -> tag      | Post belongs to tag                      |
| `co-occurrence` | tag -> tag       | Tags appear together on posts (weighted) |

## Tag Clusters

The garden page shows tags grouped with their related tags. Two tags are "related" when they appear on the same post (co-occurrence). The more posts they share, the stronger the relationship.

Tags are sorted by post count (most used first), and related tags are listed by co-occurrence weight (highest first).

The default garden page also includes a graph preview that uses `graph.json` to show top tags, their strongest relationships, and (optionally) recent posts. Post templates can add a compact preview that focuses on a single post and its connections.

### Post Graph Preview

Add the post-level graph component to show a small connections graph for the current entry. The component uses `graph_json` plus `post.href` to fetch and filter the node set down to direct connections (the post itself and its related tags/posts).

In your post template, include:

```html
{% include "components/post_graph.html" %}
```

The preview automatically hides when the post is not present in `graph.json` or has no connections.

When the preview is not rendered, its JavaScript is not loaded either, which keeps ordinary article pages lighter.

When view transitions are enabled, the preview also re-initializes after client-side page swaps, so navigating between posts does not require a hard refresh to render the graph.

## Filtering

Posts are included in the garden graph only if all of the following are true:

- `published` is `true`
- `draft` is `false`
- `private` is `false`
- `skip` is `false`
- None of the post's tags are in the `exclude_tags` list

## Node Limit

If the total number of nodes exceeds `max_nodes`, the plugin removes the least-used tags first. Posts are never removed by the node limit. Edges referencing removed nodes are also pruned.

## Excluding Tags

Use `exclude_tags` to remove specific tags from the graph entirely. Posts with excluded tags are also excluded:

```toml
[markata-go.garden]
exclude_tags = ["draft-ideas", "internal"]
```

## Custom Templates

Override the garden page by creating a `garden.html` in your templates directory. The template receives these context variables:

| Variable       | Type          | Description                       |
|----------------|---------------|-----------------------------------|
| `title`        | string        | Page title                        |
| `description`  | string        | Page description                  |
| `graph_json`   | string        | URL path to graph.json            |
| `tag_clusters` | []TagCluster  | Tag groups with related tags      |
| `total_posts`  | int           | Number of post nodes in the graph |
| `total_tags`   | int           | Number of tag nodes in the graph  |
| `total_edges`  | int           | Number of edges in the graph      |

Each `TagCluster` has:

| Field     | Type          | Description                         |
|-----------|---------------|-------------------------------------|
| `Name`    | string        | Tag name                            |
| `Count`   | int           | Number of posts with this tag       |
| `Href`    | string        | URL to the tag listing page         |
| `Related` | []TagRelation | Related tags with co-occurrence     |

Each `TagRelation` has:

| Field   | Type   | Description                                |
|---------|--------|--------------------------------------------|
| `Name`  | string | Related tag name                            |
| `Count` | int    | Number of shared posts (edge weight)       |
| `Href`  | string | URL to the related tag listing page        |

## Using graph.json with Visualization Libraries

The graph.json file is designed to work with JavaScript graph visualization libraries. Here is an example using D3.js:

```html
<script src="https://d3js.org/d3.v7.min.js"></script>
<script>
fetch('/garden/graph.json')
  .then(r => r.json())
  .then(data => {
    // data.nodes and data.edges are ready for D3 force layout
    const simulation = d3.forceSimulation(data.nodes)
      .force('link', d3.forceLink(data.edges).id(d => d.id))
      .force('charge', d3.forceManyBody())
      .force('center', d3.forceCenter(width / 2, height / 2));
  });
</script>
```

## Disabling the Garden

To disable the garden view entirely:

```toml
[markata-go.garden]
enabled = false
```

To keep the JSON export but skip the HTML page:

```toml
[markata-go.garden]
render_page = false
```
