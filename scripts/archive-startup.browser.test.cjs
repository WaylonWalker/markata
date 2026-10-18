// Behavioral checks for deferred archive layout and intent-based pagination.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const http = require('node:http');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const root = path.join(__dirname, '../pkg/themes/default/static');

test('archive leaves adjacent pages idle and preserves deferred weeks', async () => {
  const requests = [];
  const server = http.createServer((req, res) => {
    requests.push(req.url);
    if (req.url === '/navigation.js') {
      res.setHeader('content-type', 'text/javascript');
      return res.end(fs.readFileSync(path.join(root, 'js/navigation-shortcuts.js')));
    }
    if (req.url === '/feeds.css') {
      res.setHeader('content-type', 'text/css');
      return res.end(fs.readFileSync(path.join(root, 'css/feeds.css')));
    }
    res.setHeader('content-type', 'text/html');
    res.end(`<!doctype html><html><head><link rel="stylesheet" href="/feeds.css"></head>
      <body><nav><a class="pagination-next" href="/older/">Older</a></nav>
      <main class="feed feed--editorial"><div class="editorial-feed">
      ${Array.from({length: 7}, (_, i) => `<section class="editorial-week" id="week-${i}">
        <h2>Week ${i}</h2><div class="editorial-week-entries">
        <div class="editorial-entry" style="height:1800px"><a href="/post-${i}/">Story ${i}</a></div>
        <div class="editorial-entry" style="height:1800px"><a href="/extra-${i}/">Extra ${i}</a></div>
        </div></section>`).join('')}
      </div></main><script>
      window.shortcutsRegistry={register(){}};
      window.prefetchViewTransitionUrl=url=>fetch(url);
      </script><script src="/navigation.js"></script></body></html>`);
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
  try {
    const origin = `http://127.0.0.1:${server.address().port}`;
    const page = await browser.newPage({viewport: {width: 412, height: 823}});
    await page.goto(origin);
    await page.waitForFunction(() => typeof window.initNavigationShortcuts === 'function');
    await page.evaluate(() => window.initNavigationShortcuts());
    await page.waitForTimeout(100);
    assert.equal(requests.includes('/older/'), false, 'setup must not fetch another page');
    assert.equal(await page.locator('#week-0').evaluate(e => getComputedStyle(e).contentVisibility), 'visible');
    assert.equal(await page.locator('#week-6').evaluate(e => getComputedStyle(e).contentVisibility), 'auto');
    // Finding skipped content can schedule its layout for a subsequent frame.
    // Wait for a real selection instead of requiring synchronous discovery.
    await page.waitForFunction(() => window.find('Extra 6', false, false, true), null, {timeout: 5000});
    assert.equal(await page.evaluate(() => getSelection().toString()), 'Extra 6');
    // Test focus from a fresh document with skipped weeks, independently of
    // Find's selection and scroll anchoring. Allow one CSS pixel of rounding.
    await page.goto(origin + '/?focus-check');
    await page.locator('#week-6 a').first().focus();
    await page.waitForFunction(() => {
      const link = document.querySelector('#week-6 a');
      const rect = link.getBoundingClientRect();
      const visible = document.activeElement === link && rect.top >= -1 && rect.bottom <= innerHeight + 1;
      window.focusVisibleFrames = visible ? (window.focusVisibleFrames || 0) + 1 : 0;
      return window.focusVisibleFrames >= 3;
    }, null, {timeout: 5000});
    const focused = await page.locator('#week-6 a').first().evaluate(e => {
      const rect = e.getBoundingClientRect();
      return {top: rect.top, bottom: rect.bottom, viewport: innerHeight};
    });
    assert.ok(focused.top >= -1 && focused.bottom <= focused.viewport + 1,
      'keyboard focus must keep stories visible: ' + JSON.stringify(focused));
    await page.emulateMedia({media: 'print'});
    assert.equal(await page.locator('#week-6').evaluate(e => getComputedStyle(e).contentVisibility), 'visible');
    await page.evaluate(() => window.preloadNavigationShortcuts());
    await page.waitForFunction(() => performance.getEntriesByType('resource').some(e => e.name.endsWith('/older/')));
    assert.equal(requests.includes('/older/'), true, 'explicit navigation intent may prefetch');
    const noJS = await browser.newContext({javaScriptEnabled: false});
    const plain = await noJS.newPage();
    await plain.goto(origin);
    assert.equal(await plain.locator('.editorial-week a').count(), 14);
    await plain.locator('#week-6 a').first().scrollIntoViewIfNeeded();
    assert.equal(await plain.locator('#week-6 a').first().isVisible(), true);
    await noJS.close();
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
});
