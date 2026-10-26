# URL mentions provider integration (initial foundation)

The `pkg/urlmentions` package consumes version 1.0 of the public URL mentions protocol.

Provider indexes:
- https://listen.waylonwalker.com/.well-known/url-mentions.json
- https://brainrot.waylonwalker.com/.well-known/url-mentions.json

`urlmentions.Client.Fetch(ctx, endpoint)` retrieves all pages and returns a map
from provider-supplied `references[].normalized_url` to source entries.
It bounds response sizes and page counts, rejects non-HTTPS starting URLs,
and disallows cross-origin pagination. The caller should use an HTTP client
with redirects restricted to the same origin and enforce a build-time budget.

This is **phase one only**. The package is not registered as a build plugin,
does not fetch either provider automatically, and does not render UI yet.
Follow-up work: align canonical URL normalization with the provider protocol;
cache provider pages/ETags; configure providers; map post canonical, external
target, and reader URLs; SSR accessible favicon badges with no-JS fallbacks;
graceful offline builds and snapshot tests.

Do not treat this PR as production feature-complete.
