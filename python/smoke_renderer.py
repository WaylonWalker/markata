"""Run against an installed official binary wheel on each supported CI OS."""
import asyncio
from html.parser import HTMLParser
from markata import render, render_async, RendererPool

source = '# Hello\n\n```mermaid\ngraph TD\nA-->B\n```\n\n```chartjs\n{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[1]}]}}\n```'
result = render(source, inline_assets=True)
assert '<h1' in result.html
assert result.features['mermaid'] and result.features['chartjs']
assert result.external_urls == ()
class AssetParser(HTMLParser):
    def handle_starttag(self, tag, attrs):
        assert tag != "link", "inline output contains external CSS"
        assert tag != "script" or "src" not in dict(attrs), "inline output contains external JS"
AssetParser().feed(result.head_html + result.body_end_html)
assert '<style>' in result.head_html
assert '<!doctype html>' in result.document()
async_result = asyncio.run(render_async('# Async'))
assert 'Async' in async_result.html
async def smoke_pool():
    async with RendererPool(workers=2) as pool:
        results = await asyncio.gather(
            pool.render("# Warm one"),
            pool.render("# Warm two"),
            pool.render("# Reused worker"),
        )
        assert all("Warm" in item.html or "Reused" in item.html for item in results)

asyncio.run(smoke_pool())
print('Python sync/async renderer, warm pool, and inline dependencies verified.')
