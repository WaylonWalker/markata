# URL mentions SSR interface rollout

## Intent
Use the merged `pkg/urlmentions` client to render small, quiet provider indicators for canonical articles, archive cards, link posts and reader entries. A user can expand a native HTML `details` disclosure to see source titles and links. No JavaScript is required.

## Reader markup contract (initial work)
The reader templates now conditionally consume `entry.url_mentions`, a list of:
- `provider_name`: escaped visible provider name
- `icon_url`: optional HTTPS provider favicon (untrusted external source; validate)
- `source_url`: validated HTTPS URL to the referencing page
- `title`: escaped source title

These templates deliberately render nothing when the field is absent. **This PR does not yet populate the field.**

## Remaining implementation before release
1. Configurable provider endpoints, disabled by default. Reuse the bounded client; enforce build timeout, offline fallback and deterministic cached snapshots / ETag.
2. Normalize canonical article URLs, outbound link targets and reader item URLs consistently with the provider protocol. Reject malformed references and untrusted URLs.
3. Produce grouped, deduplicated mention view models during collection. Reuse the same index across article, archive, link post and reader rendering.
4. Add shared minimal CSS: compact favicon+count, quiet typography, accessible focus/hover, wrapping on mobile. Ensure long titles and provider names do not overflow.
5. Add article/archive/link-post template components; avoid duplicating per-page network requests.
6. Add fixture-backed render tests for zero/one/multiple mentions, unsafe URLs, escaping, and no-JS disclosure. Run cross-platform CI and screenshot comparison before merging.

## UX behavior
- Zero mentions: no UI.
- One or more: small 'Mentioned by' group, source favicons and count. Native disclosure reveals titles and outbound links.
- All mentions must be static HTML, available offline after a successful build.
- Provider outages must not make a previously complete build silently lose mentions.


## Reader opt-in (experimental)

Set `url_mentions_providers` in the site's extra configuration to a list of public HTTPS v1.0 index endpoints, for example:

```yaml
url_mentions_providers:
  - https://listen.waylonwalker.com/.well-known/url-mentions.json
  - https://brainrot.waylonwalker.com/.well-known/url-mentions.json
```

The blogroll writer fetches these endpoints once per build and adds `url_mentions` to reader entry template data when an entry URL matches a normalized reference. An unavailable provider does not fail the build. The initial integration does not persist an ETag or disk cache; it is opt-in and should not be treated as a complete production rollout. Other article and archive templates are not yet wired.
