# Static SVG visualization foundation

This package starts Markata's **backend-rendered visualization** roadmap. It
does not add D2/Vega-Lite/Mermaid dependencies and does not change existing
Chart.js or Mermaid behavior. It is not yet wired to Markata's site builder.

The normative description is [the static visualization specification](../../spec/spec/STATIC_VISUALIZATIONS.md);
the published [developer guide](../../docs/guides/static-visualizations.md)
explains current availability, the Go adapter surface, and remaining integration work.

## Contract

`Renderer.RenderSVG` returns static SVG for a fully specified `Request`.
The `Pipeline` caches successful results in process and deduplicates concurrent
identical requests. Its identity includes engine ID/version, source, font,
palette, effective theme fingerprint, and sorted render options.

`NewSVGAsset` validates a bounded XML document, derives pixel dimensions or
uses its viewBox, and names the output using the SHA-256 of exact SVG bytes.
`Publish` writes that content-addressed file under
`assets/visualizations/` in the site output tree. `FigureHTML` generates
script-free figure markup with explicit width/height and an accessible alt.

The XML restrictions reject active and external-resource constructs from
**trusted renderer outputs**. They are intentionally conservative, and they
are **not** a comprehensive SVG sanitizer. Untrusted authored SVG must not
be passed directly through this API. The current restrictions reject links,
images, animations, and foreignObject; adapters should report incompatible
diagrams rather than silently fall back to browser JavaScript.

## Next steps

1. Connect an optional Go D2 engine adapter. Verify theme/font layout and
   inspect actual output against SVG constraints. Compare build-size, render
   latency, memory, cold/warm cache, and compressed SVG bytes.
2. Integrate published SVG paths with the site-builder and the **additive**
   standalone renderer asset manifest. Preserve schema-v1 clients.
3. Add an optional Vega-Lite adapter using pinned `vl-convert`, with data
   fallbacks and offline output auditing.
4. Reuse this publisher from the separately tested Mermaid backend after
   native renderer evaluation finishes. Preserve current authoring syntax.
5. Add real-site page size, network-request, CLS, and screenshot regression
   checks before enabling static mode as a default.

Current unit tests (no additional dependencies):

```sh
go test ./pkg/visualizations/...
go test -race ./pkg/visualizations/...
```

This work is intentionally independent of renderer-preset issue #1854
and active native Mermaid benchmarking PR #1891.
