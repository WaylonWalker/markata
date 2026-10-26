// Deferred archive-only CSS should arrive before a rich card nears the viewport.
// Run: node --test scripts/conditional-css-lazy.browser.test.cjs
const {test} = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const conditionalCSS = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/conditional-css.js'));

async function withPage(run, disableObserver = false, visibleRichContent = false, richContentInCards = true) {
  const requests = [];
  const server = http.createServer((req, res) => {
    const pathname = new URL(req.url, 'http://localhost').pathname;
    requests.push(pathname);
    if (pathname === '/') {
      res.setHeader('content-type', 'text/html; charset=utf-8');
      return res.end(`<!doctype html><html><head>
        ${disableObserver ? '<script>window.IntersectionObserver=undefined</script>' : `<script>
          window.__activeConditionalObservers = 0;
          const NativeObserver = window.IntersectionObserver;
          window.IntersectionObserver = class extends NativeObserver {
            constructor(callback, options) { super(callback, options); window.__activeConditionalObservers++; }
            disconnect() { super.disconnect(); window.__activeConditionalObservers--; }
          };
        </script>`}
        <script src="/conditional.js" defer></script></head><body>
        <div class="feed h-feed"><article class="card">Archive card</article></div>
        <div style="height:${visibleRichContent ? 0 : 5000}px"></div>
        ${richContentInCards ? '<article class="card">' : ''}<aside id="far-admonition" class="admonition">A note</aside>
          <pre id="far-code" class="chroma"><code>print('hello')</code></pre>${richContentInCards ? '</article>' : ''}
        </body></html>`);
    }
    if (pathname === '/conditional.js') {
      res.setHeader('content-type', 'application/javascript');
      return res.end(conditionalCSS);
    }
    if (pathname.startsWith('/css/')) {
      res.setHeader('content-type', 'text/css');
      return res.end('/* test stylesheet */');
    }
    res.writeHead(404);
    res.end('not found');
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await chromium.launch({headless:true, executablePath:process.env.CHROMIUM_EXECUTABLE || undefined, args:['--no-sandbox']});
  try {
    const page = await browser.newPage({viewport:{width:412,height:823}});
    await page.goto(`http://127.0.0.1:${server.address().port}/`);
    await page.waitForLoadState('load');
    await run(page, requests);
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
}

test('defers code and admonition CSS until archive cards approach the viewport', async () => {
  await withPage(async (page, requests) => {
    await page.waitForTimeout(50);
    assert.ok(requests.includes('/css/cards.css'), 'visible archive cards load their stylesheet immediately');
    assert.ok(requests.includes('/css/feeds.css'), 'visible feed layout loads immediately');
    for (const name of ['admonitions', 'code', 'chroma']) {
      assert.equal(requests.includes(`/css/${name}.css`), false, `${name} CSS stays off the cold path`);
    }

    await page.evaluate(() => window.scrollTo(0, 2000));
    await page.waitForFunction(() => ['admonitions', 'code', 'chroma'].every(name =>
      [...document.querySelectorAll('link[rel="stylesheet"]')].some(link => link.href.includes(`/css/${name}.css`))));
    for (const name of ['admonitions', 'code', 'chroma']) {
      assert.equal(requests.filter(path => path === `/css/${name}.css`).length, 1, `${name} CSS loads once near content`);
    }
  });
});

test('loads deferred conditional CSS immediately without IntersectionObserver', async () => {
  await withPage(async (page, requests) => {
    await page.waitForFunction(() => ['admonitions', 'code', 'chroma'].every(name =>
      [...document.querySelectorAll('link[rel="stylesheet"]')].some(link => link.href.includes(`/css/${name}.css`))));
    for (const name of ['admonitions', 'code', 'chroma']) {
      assert.equal(requests.filter(path => path === `/css/${name}.css`).length, 1, `${name} CSS fallback loads once`);
    }
  }, true);
});

test('loads deferred rich-content CSS when a rich archive card is visible', async () => {
  await withPage(async (page, requests) => {
    await page.waitForFunction(() => ['admonitions', 'code', 'chroma'].every(name =>
      [...document.querySelectorAll('link[rel="stylesheet"]')].some(link => link.href.includes(`/css/${name}.css`))));
    for (const name of ['admonitions', 'code', 'chroma']) {
      assert.equal(requests.filter(path => path === `/css/${name}.css`).length, 1, `${name} CSS loads once for visible rich content`);
    }
  }, false, true);
});

test('loads article rich-content CSS immediately outside archive cards', async () => {
  await withPage(async (page, requests) => {
    await page.waitForFunction(() => ['admonitions', 'code', 'chroma'].every(name =>
      [...document.querySelectorAll('link[rel="stylesheet"]')].some(link => link.href.includes(`/css/${name}.css`))));
    for (const name of ['admonitions', 'code', 'chroma']) {
      assert.equal(requests.filter(path => path === `/css/${name}.css`).length, 1, `${name} CSS is immediate for article content`);
    }
  }, false, false, false);
});

test('replaces pending observers across repeated view transitions', async () => {
  await withPage(async (page, requests) => {
    await page.waitForFunction(() => window.__activeConditionalObservers === 3);
    await page.evaluate(() => {
      dispatchEvent(new Event('view-transition-complete'));
      dispatchEvent(new Event('view-transition-complete'));
    });
    await page.waitForFunction(() => window.__activeConditionalObservers === 3);
    assert.equal(requests.filter(path => path.endsWith('.css')).length, 2, 'only immediate cards and feeds CSS loaded while content remains far away');
    await page.evaluate(() => window.scrollTo(0, 2000));
    await page.waitForFunction(() => ['admonitions', 'code', 'chroma'].every(name =>
      [...document.querySelectorAll('link[rel="stylesheet"]')].some(link => link.href.includes(`/css/${name}.css`))));
    for (const name of ['admonitions', 'code', 'chroma']) {
      assert.equal(requests.filter(path => path === `/css/${name}.css`).length, 1, `${name} loads once after repeated scans`);
    }
  });
});

test('disconnects deferred observers when a transition removes rich card content', async () => {
  await withPage(async (page, requests) => {
    await page.waitForFunction(() => window.__activeConditionalObservers === 3);
    await page.evaluate(() => {
      document.querySelector('#far-admonition').closest('.card').remove();
      dispatchEvent(new Event('view-transition-complete'));
    });
    await page.waitForFunction(() => window.__activeConditionalObservers === 0);
    assert.equal(requests.some(path => ['/css/admonitions.css', '/css/code.css', '/css/chroma.css'].includes(path)), false,
      'removed rich content does not trigger deferred stylesheet requests');
  });
});
